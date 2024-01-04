package relayer

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	MuSigner struct {
		cfg            muSignerConfig
		cancelHanlders context.CancelFunc

		txBuffer           chan *txMessage
		unsignedTxMap      map[chainhash.Hash]*txMessage
		unsignedTxMapMutex sync.RWMutex

		nonceBuffer     chan *nonceMessage
		nouncesMap      map[chainhash.Hash]map[pubKey][]*types.MusigNonce
		nouncesMapMutex sync.RWMutex

		sigBuffer    chan *signatureMessage
		sigsMap      map[chainhash.Hash]map[pubKey][]*types.MusigPartialSignature
		sigsMapMutex sync.RWMutex

		signedTxMap *lru.Cache[chainhash.Hash, struct{}]
	}

	sinkSignedTransactionFunc func(sigs []*wire.MsgTx) error

	muSignerConfig struct {
		btcClient           *rpcclient.Client
		sinkSignedTXHandler sinkSignedTransactionFunc
		sinkInterval        time.Duration
		witnessPubKeys      []*btcec.PublicKey
	}

	txMessage struct {
		tx              *wire.MsgTx
		transferMapping map[uint64]common.Hash
	}

	nonceMessage struct {
		sender pubKey
		msg    *types.TransactionNonces
	}

	signatureMessage struct {
		sender pubKey
		msg    *types.TransactionSignatures
	}

	pubKey [secp256k1.PubKeyBytesLenCompressed]byte
)

const (
	_defaultBufferSize = 3000
)

func NewMuSigner(cfg muSignerConfig) *MuSigner {
	signedTxMapLRU, err := lru.New[chainhash.Hash, struct{}](_defaultBufferSize)
	if err != nil {
		log.Fatalln(err)
	}
	return &MuSigner{
		txBuffer:      make(chan *txMessage, _defaultBufferSize),
		unsignedTxMap: make(map[chainhash.Hash]*txMessage),
		nonceBuffer:   make(chan *nonceMessage, _defaultBufferSize),
		nouncesMap:    make(map[chainhash.Hash]map[pubKey][]*types.MusigNonce),
		sigBuffer:     make(chan *signatureMessage, _defaultBufferSize),
		sigsMap:       make(map[chainhash.Hash]map[pubKey][]*types.MusigPartialSignature),
		signedTxMap:   signedTxMapLRU,
		cfg:           cfg,
	}
}

func (ms *MuSigner) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	go ms.run(ctx, ms.cfg.sinkInterval)
	ms.cancelHanlders = cancel
	return nil
}

func (ms *MuSigner) Stop() error {
	ms.cancelHanlders()
	return nil
}

func (m *MuSigner) AddTX(msg *txMessage) error {
	select {
	case m.txBuffer <- msg:
		return nil
	default:
		return errors.New("the msg buffer of writer is full")
	}
}

func (m *MuSigner) UnsignedTXWithoutNonces(excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithoutNoncesResponse, error) {
	m.unsignedTxMapMutex.RLock()
	defer m.unsignedTxMapMutex.RUnlock()

	excludedTxsMap := make(map[chainhash.Hash]struct{})
	for _, txHash := range excludedTxs.GetTxHash() {
		hh, err := chainhash.NewHashFromStr(txHash)
		if err != nil {
			return nil, err
		}
		excludedTxsMap[*hh] = struct{}{}
	}

	txs := make([]*types.BTCTransaction, 0, len(m.unsignedTxMap))
	for _, txMsg := range m.unsignedTxMap {
		if _, excluded := excludedTxsMap[txMsg.tx.TxHash()]; excluded {
			continue
		}
		var b bytes.Buffer
		if err := txMsg.tx.SerializeNoWitness(&b); err != nil {
			return nil, err
		}
		txs = append(txs, &types.BTCTransaction{
			Data: b.Bytes(),
		})
	}
	return &services.ListUnsignedBTCTXWithoutNoncesResponse{Transactions: txs}, nil
}

func (m *MuSigner) AddMusigNonces(req *services.MusigNonceMessage) error {
	hh, err := util.HashProto(req.GetMessage())
	if err != nil {
		return err
	}
	pbk, err := btcec.ParsePubKey(req.Pubkey)
	if err != nil {
		return err
	}
	sig, err := schnorr.ParseSignature(req.Signature)
	if err != nil {
		return err
	}
	if !m.validateSource(hh[:], pbk, sig) {
		return errors.New("invalid signature when adding nonces")
	}
	sender := pubKey{}
	copy(sender[:], pbk.SerializeCompressed()[:])
	for _, nonce := range req.GetMessage().GetNonces() {
		select {
		case m.nonceBuffer <- &nonceMessage{sender: sender, msg: nonce}:
			continue
		default:
			return errors.New("the nonce buffer of writer is full")
		}
	}
	return nil
}

func (m *MuSigner) validateSource(msgHash []byte, pk *btcec.PublicKey, signature *schnorr.Signature) bool {
	for _, wp := range m.cfg.witnessPubKeys {
		if wp.IsEqual(pk) {
			return signature.Verify(msgHash, pk)
		}
	}
	log.Println("no matching pubkey found")
	return false
}

func (m *MuSigner) UnsignedTXWithNonces(excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithNoncesResponse, error) {
	type txWithNonces struct {
		hash            chainhash.Hash
		tx              *wire.MsgTx
		noncesFromWit   map[util.NonceID]*pubNonceGroup
		transferMapping map[uint64]common.Hash
	}

	excludedTxsMap := make(map[chainhash.Hash]struct{})
	for _, txHash := range excludedTxs.GetTxHash() {
		hh, err := chainhash.NewHashFromStr(txHash)
		if err != nil {
			return nil, err
		}
		excludedTxsMap[*hh] = struct{}{}
	}

	txs := []*txWithNonces{}
	m.nouncesMapMutex.RLock()
	for txHash, noncesFromWit := range m.nouncesMap {
		if _, excluded := excludedTxsMap[txHash]; excluded {
			continue
		}
		if len(noncesFromWit) == len(m.cfg.witnessPubKeys) {
			tx := &txWithNonces{
				hash:          txHash,
				noncesFromWit: make(map[util.NonceID]*pubNonceGroup),
			}
			for _, nonces := range noncesFromWit {
				for _, nonce := range nonces {
					nonceID := util.NonceIDFromTx(txHash, nonce.GetTxinIndex())
					if _, exist := tx.noncesFromWit[nonceID]; !exist {
						tx.noncesFromWit[nonceID] = newPubNonceGroup()
					}
					tx.noncesFromWit[nonceID].add(nonce.GetPubNonce())
				}
			}
			txs = append(txs, tx)
		}
	}
	m.nouncesMapMutex.RUnlock()

	m.unsignedTxMapMutex.RLock()
	for _, txInfo := range txs {
		txMsg, exist := m.unsignedTxMap[txInfo.hash]
		if !exist {
			m.unsignedTxMapMutex.RUnlock()
			return nil, errors.New("unsigned tx not found")
		}
		txInfo.tx = txMsg.tx
		txInfo.transferMapping = txMsg.transferMapping
	}
	m.unsignedTxMapMutex.RUnlock()

	combinedNonces := make([]*types.TransactionWithCombinedNonces, 0, len(txs))
	for _, tx := range txs {
		var serializedTx bytes.Buffer
		if err := tx.tx.SerializeNoWitness(&serializedTx); err != nil {
			return nil, err
		}
		txWithCombinedNonce := &types.TransactionWithCombinedNonces{
			Tx: &types.BTCTransaction{
				Data: serializedTx.Bytes(),
			},
			CombinedNonces: make([][]byte, 0, len(tx.tx.TxIn)),
			TransferMap:    make([]*types.TransferMap, 0, len(tx.transferMapping)),
		}
		for vout, transferID := range tx.transferMapping {
			txWithCombinedNonce.TransferMap = append(txWithCombinedNonce.TransferMap,
				&types.TransferMap{
					Vout:       vout,
					TransferId: transferID.String(),
				})
		}
		for i := range tx.tx.TxIn {
			nonceID := util.NonceIDFromTx(tx.hash, uint64(i))
			combinedNonce, err := tx.noncesFromWit[nonceID].combine()
			if err != nil {
				log.Printf("failed to combine nonces for tx %s, index %d", tx.hash.String(), i)
				log.Println(err)
				return nil, err
			}
			txWithCombinedNonce.CombinedNonces = append(txWithCombinedNonce.CombinedNonces, combinedNonce[:])
		}
		combinedNonces = append(combinedNonces, txWithCombinedNonce)
	}

	serailizedPubkeys := make([][]byte, 0, len(m.cfg.witnessPubKeys))
	for _, pk := range m.cfg.witnessPubKeys {
		serailizedPubkeys = append(serailizedPubkeys, pk.SerializeCompressed())
	}

	return &services.ListUnsignedBTCTXWithNoncesResponse{
		Pubkeys:      serailizedPubkeys,
		Transactions: combinedNonces,
	}, nil
}

func (m *MuSigner) AddMusigSignatures(req *services.MusigSignatureMessage) error {
	hh, err := util.HashProto(req.GetMessage())
	if err != nil {
		return err
	}
	pbk, err := btcec.ParsePubKey(req.Pubkey)
	if err != nil {
		return err
	}
	sig, err := schnorr.ParseSignature(req.Signature)
	if err != nil {
		return err
	}
	if !m.validateSource(hh[:], pbk, sig) {
		return errors.New("invalid signature")
	}
	sender := pubKey{}
	copy(sender[:], pbk.SerializeCompressed()[:])
	for _, sigs := range req.GetMessage().GetSignatures() {
		select {
		case m.sigBuffer <- &signatureMessage{sender: sender, msg: sigs}:
			continue
		default:
			return errors.New("the nonce buffer of writer is full")
		}
	}
	return nil
}

func (m *MuSigner) run(ctx context.Context, interval time.Duration) {
	assembleTicker := time.NewTicker(interval)
	for {
		select {
		case txMsg := <-m.txBuffer:
			txHash := txMsg.tx.TxHash()
			log.Printf("Musigner received tx %s\n", txHash.String())
			if m.signedTxMap.Contains(txHash) {
				continue
			}
			if _, exist := m.unsignedTxMap[txHash]; !exist {
				m.unsignedTxMapMutex.Lock()
				m.unsignedTxMap[txHash] = txMsg
				log.Printf("Musigner added tx %s into signing list(%d txs)\n", txHash.String(), len(m.unsignedTxMap))
				m.unsignedTxMapMutex.Unlock()
			}
		case nonceMsg := <-m.nonceBuffer:
			log.Printf("Musigner received the nonces of transaction %s from %s\n", nonceMsg.msg.TxHash, hex.EncodeToString(nonceMsg.sender[:]))
			hh, err := chainhash.NewHashFromStr(nonceMsg.msg.TxHash)
			if err != nil {
				log.Printf("failed to parse tx hash %s\n", nonceMsg.msg.TxHash)
				continue
			}
			if m.signedTxMap.Contains(*hh) {
				continue
			}
			m.nouncesMapMutex.Lock()
			if _, exist := m.nouncesMap[*hh]; !exist {
				m.nouncesMap[*hh] = make(map[pubKey][]*types.MusigNonce)
			}
			m.nouncesMap[*hh][nonceMsg.sender] = nonceMsg.msg.GetNonces()
			log.Printf("Musigner added tx %s into nonce list(%d txs)\n", nonceMsg.msg.TxHash, len(m.nouncesMap))
			m.nouncesMapMutex.Unlock()
		case sigMessage := <-m.sigBuffer:
			log.Printf("Musigner received the signatures of transaction %s from %s\n", sigMessage.msg.TxHash, hex.EncodeToString(sigMessage.sender[:]))
			hh, err := chainhash.NewHashFromStr(sigMessage.msg.TxHash)
			if err != nil {
				log.Println(err)
				continue
			}
			if m.signedTxMap.Contains(*hh) {
				continue
			}
			m.sigsMapMutex.Lock()
			if _, exist := m.sigsMap[*hh]; !exist {
				m.sigsMap[*hh] = make(map[pubKey][]*types.MusigPartialSignature)
			}
			m.sigsMap[*hh][sigMessage.sender] = sigMessage.msg.GetSignatures()
			log.Printf("Musigner added tx %s into sig list(%d txs)\n", sigMessage.msg.TxHash, len(m.sigsMap))
			m.sigsMapMutex.Unlock()
		case <-assembleTicker.C:
			signedBTCTXs, signedBTCTXHashes, err := m.sythesizeSignedTXs()
			if err != nil {
				log.Println(err)
			}
			if len(signedBTCTXs) == 0 {
				// Skip if no txs are sythesized
				continue
			}
			log.Printf("Musigner sythesized %d signed txs: ", len(signedBTCTXs))
			for _, hh := range signedBTCTXHashes {
				log.Printf("%s ", hh.String())
			}
			log.Printf("\n")

			if err := m.cfg.sinkSignedTXHandler(signedBTCTXs); err != nil {
				log.Println(err)
				continue
			}

			for _, txHash := range signedBTCTXHashes {
				m.signedTxMap.Add(txHash, struct{}{})
			}
			// Remove sinked txs from maps
			m.sigsMapMutex.Lock()
			for _, txHash := range signedBTCTXHashes {
				delete(m.sigsMap, txHash)
			}
			m.sigsMapMutex.Unlock()
			m.nouncesMapMutex.Lock()
			for _, txHash := range signedBTCTXHashes {
				delete(m.nouncesMap, txHash)
			}
			m.nouncesMapMutex.Unlock()
			m.unsignedTxMapMutex.Lock()
			for _, txHash := range signedBTCTXHashes {
				delete(m.unsignedTxMap, txHash)
			}
			m.unsignedTxMapMutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (m *MuSigner) sythesizeSignedTXs() ([]*wire.MsgTx, []chainhash.Hash, error) {
	type txWithSigs struct {
		hash        chainhash.Hash
		tx          *wire.MsgTx
		sigsFromWit map[util.NonceID]*signatureGroup
	}
	txs := []*txWithSigs{}
	m.sigsMapMutex.RLock()
	for txHash, sigsFromWit := range m.sigsMap {
		if len(sigsFromWit) == len(m.cfg.witnessPubKeys) {
			tx := &txWithSigs{
				hash:        txHash,
				sigsFromWit: make(map[util.NonceID]*signatureGroup),
			}
			for _, sigs := range sigsFromWit {
				for _, sig := range sigs {
					nonceID := util.NonceIDFromTx(txHash, sig.GetTxinIndex())
					if _, exist := tx.sigsFromWit[nonceID]; !exist {
						tx.sigsFromWit[nonceID] = newSignatureGroup()
					}
					tx.sigsFromWit[nonceID].add(sig.GetCombinedNonce(), sig.GetSignature())
				}
			}
			txs = append(txs, tx)
		}
	}
	m.sigsMapMutex.RUnlock()

	m.unsignedTxMapMutex.RLock()
	for _, txInfo := range txs {
		txMsg, exist := m.unsignedTxMap[txInfo.hash]
		if !exist {
			m.unsignedTxMapMutex.RUnlock()
			return nil, nil, errors.New("unsigned tx not found")
		}
		txInfo.tx = txMsg.tx
	}
	m.unsignedTxMapMutex.RUnlock()

	signedTxs := make([]*wire.MsgTx, 0)
	signedTxHashes := make([]chainhash.Hash, 0)
	for _, tx := range txs {
		signedTx := tx.tx.Copy()
		prevOutputFetcher, err := util.NewPrevOutFetcher(signedTx, m.cfg.btcClient)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}
		sigHashes := txscript.NewTxSigHashes(signedTx, prevOutputFetcher)
		failedFlag := false
		for i := range signedTx.TxIn {
			msg, err := util.TaprootSignatureHash(signedTx, i, txscript.SigHashDefault, sigHashes, prevOutputFetcher)
			if err != nil {
				log.Println(err)
				failedFlag = true
				break
			}
			nonceID := util.NonceIDFromTx(tx.hash, uint64(i))
			sig, err := tx.sigsFromWit[nonceID].combine(msg, m.cfg.witnessPubKeys)
			if err != nil {
				log.Println(err)
				failedFlag = true
				break
			}
			signedTx.TxIn[i].Witness = wire.TxWitness{sig}
		}
		if !failedFlag {
			log.Printf("Signed tx %s!\n", signedTx.TxHash().String())
			signedTxs = append(signedTxs, signedTx)
			signedTxHashes = append(signedTxHashes, signedTx.TxHash())
		}
	}
	var err error
	if len(signedTxs) != len(txs) {
		err = errors.New("part of txs' signatures are invalid")
	}
	return signedTxs, signedTxHashes, err
}

type pubNonceGroup struct {
	nonces [][musig2.PubNonceSize]byte
}

func newPubNonceGroup() *pubNonceGroup {
	return &pubNonceGroup{
		nonces: make([][musig2.PubNonceSize]byte, 0),
	}
}

func (p *pubNonceGroup) add(pubNonce []byte) error {
	if len(pubNonce) != musig2.PubNonceSize {
		return errors.Errorf("invalid pub nonce size %d", len(pubNonce))
	}
	var pn [musig2.PubNonceSize]byte
	copy(pn[:], pubNonce)
	p.nonces = append(p.nonces, pn)
	return nil
}

func (p *pubNonceGroup) combine() ([musig2.PubNonceSize]byte, error) {
	return musig2.AggregateNonces(p.nonces)
}

type signatureGroup struct {
	sigs []*musig2.PartialSignature
}

func newSignatureGroup() *signatureGroup {
	return &signatureGroup{
		sigs: make([]*musig2.PartialSignature, 0),
	}
}

func (p *signatureGroup) add(R []byte, S []byte) error {
	sig := &musig2.PartialSignature{}
	var err error
	sig.R, err = btcec.ParsePubKey(R)
	if err != nil {
		return err
	}
	if err := sig.Decode(bytes.NewReader(S)); err != nil {
		return err
	}
	p.sigs = append(p.sigs, sig)
	return nil
}

func (p *signatureGroup) combine(msg [chainhash.HashSize]byte, pubkeySet []*btcec.PublicKey) ([]byte, error) {
	if len(p.sigs) == 0 {
		return nil, errors.New("no sigs")
	}

	finalNonce := p.sigs[0].R
	shouldSort := true
	finalSig := musig2.CombineSigs(finalNonce, p.sigs, musig2.WithBip86TweakedCombine(msg, pubkeySet, shouldSort))

	pbk, err := util.Musig2CombinedPubkey(pubkeySet)
	if err != nil {
		return nil, err
	}
	finalSig.Verify(msg[:], pbk)
	return finalSig.Serialize(), nil
}
