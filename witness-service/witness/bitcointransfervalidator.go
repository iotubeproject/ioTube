package witness

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type BTCValidator struct {
	privateKey *secp256k1.PrivateKey
	relayerURL string
	btcClient  *rpcclient.Client
	recorder   *Recorder
	processor  dispatcher.Runner

	noncesSet *lru.Cache[chainhash.Hash, map[uint64]*musig2.Nonces]
	sigsCache *lru.Cache[chainhash.Hash, map[uint64]*musig2.PartialSignature]
}

const (
	lruSize = 500
)

func NewBTCTransferValidator(pvk *ecdsa.PrivateKey, btcClient *rpcclient.Client,
	processInterval time.Duration, relayerURL string, recorder *Recorder) *BTCValidator {
	secpPvk, err := util.SecpPvkFromEcdsa(pvk)
	if err != nil {
		log.Fatalln(err)
	}

	noncesSetLRU, err := lru.New[chainhash.Hash, map[uint64]*musig2.Nonces](lruSize)
	if err != nil {
		log.Fatalln(err)
	}
	sigsCacheLRU, err := lru.New[chainhash.Hash, map[uint64]*musig2.PartialSignature](lruSize)
	if err != nil {
		log.Fatalln(err)
	}

	bv := &BTCValidator{
		privateKey: secpPvk,
		relayerURL: relayerURL,
		btcClient:  btcClient,
		recorder:   recorder,
		noncesSet:  noncesSetLRU,
		sigsCache:  sigsCacheLRU,
	}

	if bv.processor, err = dispatcher.NewRunner(processInterval, bv.process); err != nil {
		log.Fatalln(err)
	}

	return bv
}

func (b *BTCValidator) Start(ctx context.Context) error {
	return b.processor.Start()
}

func (b *BTCValidator) Stop(ctx context.Context) error {
	return b.processor.Close()
}

func (b *BTCValidator) process() error {
	if err := b.SyncMusigNonces(); err != nil {
		util.LogErr(err)
	}

	if err := b.SignBTCTransactions(); err != nil {
		util.LogErr(err)
	}
	return nil
}

func (b *BTCValidator) SyncMusigNonces() error {
	log.Println("BTC Validator Syncing MusigNonces...")
	txs, err := b.fetchUnsignedBTCTransactionsWithoutNonces()
	if err != nil {
		return err
	}
	nonces, err := b.generateMusigNonces(txs)
	if err != nil {
		return err
	}
	return b.submitMusigNonces(nonces)
}

func (b *BTCValidator) fetchUnsignedBTCTransactionsWithoutNonces() ([]*wire.MsgTx, error) {
	txHashes := make([]string, 0)
	for _, tx := range b.noncesSet.Keys() {
		txHashes = append(txHashes, tx.String())
	}

	conn, err := grpc.Dial(b.relayerURL, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to relayer")
	}
	defer conn.Close()
	resp, err := services.NewRelayServiceClient(conn).ListUnsignedBTCTXWithoutNonces(context.Background(),
		&services.ExcludedTransactions{
			TxHash: txHashes,
		})
	if err != nil {
		return nil, err
	}
	txs := make([]*wire.MsgTx, 0)
	for _, v := range resp.GetTransactions() {
		serializedTx := v.GetData()
		var msgTx wire.MsgTx
		if err := msgTx.DeserializeNoWitness(bytes.NewReader(serializedTx)); err != nil {
			return nil, err
		}
		txs = append(txs, &msgTx)
	}
	return txs, nil
}

func (b *BTCValidator) generateMusigNonces(txs []*wire.MsgTx) ([]*types.TransactionNonces, error) {
	res := make([]*types.TransactionNonces, 0)
	for _, msgTx := range txs {
		txHash := msgTx.TxHash()
		if !b.noncesSet.Contains(txHash) {
			b.noncesSet.Add(txHash, make(map[uint64]*musig2.Nonces))
		}
		nonceMap, _ := b.noncesSet.Get(txHash)
		txWithNonces := []*types.MusigNonce{}
		for j := range msgTx.TxIn {
			if _, exist := nonceMap[uint64(j)]; !exist {
				var err error
				log.Printf("generating nonce for transaction %s idx %d\n", txHash.String(), j)
				nonce, err := musig2.GenNonces(musig2.WithPublicKey(b.privateKey.PubKey()), musig2.WithNonceSecretKeyAux(b.privateKey))
				if err != nil {
					return nil, err
				}
				nonceMap[uint64(j)] = nonce
			}
			nonce := nonceMap[uint64(j)]
			txWithNonces = append(txWithNonces, &types.MusigNonce{
				TxHash:    txHash.String(),
				TxinIndex: uint64(j),
				PubNonce:  nonce.PubNonce[:],
			})
		}
		res = append(res, &types.TransactionNonces{
			TxHash: txHash.String(),
			Nonces: txWithNonces,
		})
	}
	return res, nil
}

func (b *BTCValidator) submitMusigNonces(nonces []*types.TransactionNonces) error {
	if len(nonces) == 0 {
		return nil
	}

	msg := &services.MusigNonces{
		Nonces: nonces,
	}

	pbk, sig, err := b.signProtoMessage(msg)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(b.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to connect to relayer")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	_, err = relayer.SubmitMusigNonces(context.Background(), &services.MusigNonceMessage{
		Message:   msg,
		Pubkey:    pbk,
		Signature: sig,
	})
	return err
}

func (b *BTCValidator) signProtoMessage(msg protoreflect.ProtoMessage) ([]byte, []byte, error) {
	hh, err := util.HashProto(msg)
	if err != nil {
		return nil, nil, err
	}

	sig, err := schnorr.Sign(b.privateKey, hh[:])
	if err != nil {
		return nil, nil, err
	}

	return b.privateKey.PubKey().SerializeCompressed(), sig.Serialize(), nil
}

func (b *BTCValidator) SignBTCTransactions() error {
	txs, nonceMap, transferMap, pubkeys, err := b.fetchUnsignedBTCTransactionsWithNonces()
	if err != nil {
		return err
	}
	log.Printf("BTC Validator Signing %d btc transactions\n", len(txs))
	sigs, err := b.muSig2(txs, nonceMap, transferMap, pubkeys)
	if err != nil {
		return err
	}
	return b.submitMusigSigs(sigs)
}

func (b *BTCValidator) fetchUnsignedBTCTransactionsWithNonces() (
	[]*wire.MsgTx, map[util.NonceID][musig2.PubNonceSize]byte,
	map[util.NonceID]common.Hash, []*btcec.PublicKey, error) {
	txHashes := make([]string, 0)
	for _, tx := range b.sigsCache.Keys() {
		txHashes = append(txHashes, tx.String())
	}

	conn, err := grpc.Dial(b.relayerURL, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "failed to connect to relayer")
	}
	defer conn.Close()
	resp, err := services.NewRelayServiceClient(conn).ListUnsignedBTCTXWithNonces(context.Background(),
		&services.ExcludedTransactions{
			TxHash: txHashes,
		})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	txs := make([]*wire.MsgTx, 0)
	nonceMap := make(map[util.NonceID][musig2.PubNonceSize]byte)
	transferMap := make(map[util.NonceID]common.Hash)
	for _, v := range resp.GetTransactions() {
		serializedTx := v.GetTx().GetData()
		var msgTx wire.MsgTx
		if err := msgTx.DeserializeNoWitness(bytes.NewReader(serializedTx)); err != nil {
			return nil, nil, nil, nil, err
		}
		txs = append(txs, &msgTx)
		for i := range v.GetCombinedNonces() {
			if len(v.GetCombinedNonces()[i]) != musig2.PubNonceSize {
				return nil, nil, nil, nil, errors.New("invalid nonce size")
			}
			nonce := [musig2.PubNonceSize]byte{}
			copy(nonce[:], v.GetCombinedNonces()[i])
			nonceID := util.NonceIDFromTx(msgTx.TxHash(), uint64(i))
			nonceMap[nonceID] = nonce
		}

		for _, tsf := range v.GetTransferMap() {
			nonceID := util.NonceIDFromTx(msgTx.TxHash(), uint64(tsf.GetVout()))
			transferMap[nonceID] = common.HexToHash(tsf.GetTransferId())
		}
	}
	pubkeys := make([]*btcec.PublicKey, 0)
	for _, v := range resp.GetPubkeys() {
		pbk, err := btcec.ParsePubKey(v)
		if err != nil {
			return nil, nil, nil, nil, errors.New("invalid pubkey")
		}
		pubkeys = append(pubkeys, pbk)
	}
	return txs, nonceMap, transferMap, pubkeys, nil
}

func (b *BTCValidator) muSig2(txs []*wire.MsgTx,
	combinedNonceMap map[util.NonceID][musig2.PubNonceSize]byte,
	transferMap map[util.NonceID]common.Hash, pubKeys []*btcec.PublicKey,
) ([]*types.TransactionSignatures, error) {
	res := make([]*types.TransactionSignatures, 0)
	for _, msgTx := range txs {
		prevOutputFetcher, err := util.NewPrevOutFetcher(msgTx, b.btcClient)
		if err != nil {
			log.Printf("failed to create prev output fetcher for transaction %s, error %+v\n", msgTx.TxHash().String(), err)
			return nil, err
		}

		if !b.validateTx(msgTx, prevOutputFetcher, transferMap) {
			log.Printf("invalid transaction %s\n", msgTx.TxHash().String())
			continue
		}

		var (
			sigHashes  = txscript.NewTxSigHashes(msgTx, prevOutputFetcher)
			txHash     = msgTx.TxHash()
			sigs       = []*types.MusigPartialSignature{}
			failedFlag = false
		)

		if !b.sigsCache.Contains(txHash) {
			b.sigsCache.Add(txHash, make(map[uint64]*musig2.PartialSignature))
		}
		sigMap, _ := b.sigsCache.Get(txHash)
		for j := range msgTx.TxIn {
			nonceID := util.NonceIDFromTx(txHash, uint64(j))
			if _, exist := sigMap[uint64(j)]; !exist {
				nonceMap, exist := b.noncesSet.Get(txHash)
				if !exist {
					log.Printf("nonce not found for transaction %s %d, nonceID %d\n", txHash.String(), j, nonceID)
					failedFlag = true
					break
				}
				nonce, exist := nonceMap[uint64(j)]
				if !exist {
					log.Printf("nonce not found for transaction %s %d, nonceID %d\n", txHash.String(), j, nonceID)
					failedFlag = true
					break
				}
				combinedNonce, exist := combinedNonceMap[nonceID]
				if !exist {
					log.Printf("combined nonce not found for transaction %s %d, nonceID %d\n", txHash.String(), j, nonceID)
					failedFlag = true
					break
				}

				msg, err := util.TaprootSignatureHash(msgTx, j, txscript.SigHashDefault, sigHashes, prevOutputFetcher)
				if err != nil {
					log.Printf("failed to get signature hash for transaction %s %d, nonceID %d, error %+v\n", txHash.String(), j, nonceID, err)
					failedFlag = true
					break
				}
				sig, err := musig2.Sign(nonce.SecNonce, b.privateKey, combinedNonce, pubKeys, msg, musig2.WithBip86SignTweak())
				if err != nil {
					log.Printf("failed to sign transaction %s %d, nonceID %d, error %+v\n", txHash.String(), j, nonceID, err)
					failedFlag = true
					break
				}
				sigMap[uint64(j)] = sig
			}
			sig := sigMap[uint64(j)]
			var buf bytes.Buffer
			if err := sig.Encode(&buf); err != nil {
				return nil, err
			}
			sigs = append(sigs, &types.MusigPartialSignature{
				TxHash:        txHash.String(),
				TxinIndex:     uint64(j),
				Signature:     buf.Bytes(),
				CombinedNonce: sig.R.SerializeUncompressed(),
			})
		}
		if !failedFlag {
			log.Printf("transaction %s is signed\n", txHash.String())
			res = append(res, &types.TransactionSignatures{
				TxHash:     txHash.String(),
				Signatures: sigs,
			})
		}
	}
	return res, nil
}

// TODO: move to config
var (
	_maxTxFeeRate, _ = btcutil.NewAmount(0.001)
)

func (b *BTCValidator) validateTx(
	tx *wire.MsgTx,
	prevOutputFetcher *txscript.MultiPrevOutFetcher,
	transferMap map[util.NonceID]common.Hash,
) bool {
	if len(tx.TxOut) != len(transferMap) && len(tx.TxOut) != len(transferMap)+1 {
		log.Printf("invalid transaction %s, txout length %d, transferMap length %d\n", tx.TxHash().String(), len(tx.TxOut), len(transferMap))
		return false
	}
	var totalSatoshiOut uint64
	for idx, txOut := range tx.TxOut {
		nonceID := util.NonceIDFromTx(tx.TxHash(), uint64(idx))
		if transferID, exist := transferMap[nonceID]; exist {
			// validate amount matching transfers in the local db
			tsf, err := b.recorder.Transfer(transferID)
			if err != nil {
				log.Printf("failed to get transfer %s, error %+v\n", transferID.String(), err)
				return false
			}
			if tsf.Amount().Cmp(big.NewInt(txOut.Value)) != 0 {
				log.Printf("invalid transaction %s, txout value %d, transfer amount %d\n", tx.TxHash().String(), txOut.Value, tsf.Amount())
				return false
			}
			// validate dest matching transfers in the local db
			addr, ok := tsf.Recipient().Address().(btcutil.Address)
			if !ok {
				log.Printf("invalid address type %T\n", tsf.Recipient().Address())
				return false
			}
			pkscript, err := txscript.PayToAddrScript(addr)
			if err != nil {
				log.Printf("failed to get pkscript for address %s, error %+v\n", addr.String(), err)
				return false
			}
			if !bytes.Equal(pkscript, txOut.PkScript) {
				log.Printf("invalid transaction %s, txout pkscript %s, transfer pkscript %s\n", tx.TxHash().String(), txOut.PkScript, pkscript)
				return false
			}
		} else {
			// validate the change back to the sender
			if len(tx.TxIn) < 1 {
				return false
			}
			prevTxOut := prevOutputFetcher.FetchPrevOutput(tx.TxIn[0].PreviousOutPoint)
			if !bytes.Equal(prevTxOut.PkScript, txOut.PkScript) {
				log.Printf("invalid transaction %s, txout pkscript %s, prevTxOut pkscript %s\n", tx.TxHash().String(), txOut.PkScript, prevTxOut.PkScript)
				return false
			}
		}
		totalSatoshiOut += uint64(txOut.Value)
	}
	// validate the reasonable miner fee
	var totalSatoshiIn uint64
	for _, txIn := range tx.TxIn {
		prevTxOut := prevOutputFetcher.FetchPrevOutput(txIn.PreviousOutPoint)
		totalSatoshiIn += uint64(prevTxOut.Value)
	}
	txFeeInSatoshi := totalSatoshiIn - totalSatoshiOut
	FeePerKB := txFeeInSatoshi * 1000 / uint64(mempool.GetTxVirtualSize(btcutil.NewTx(tx)))
	if FeePerKB > uint64(_maxTxFeeRate) {
		log.Printf("invalid transaction %s, fee per kb %d\n", tx.TxHash().String(), FeePerKB)
		return false
	}

	return true
}

func (b *BTCValidator) submitMusigSigs(sigs []*types.TransactionSignatures) error {
	if len(sigs) == 0 {
		return nil
	}
	msg := &services.MusigSignatures{
		Signatures: sigs,
	}

	pbk, sig, err := b.signProtoMessage(msg)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(b.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to connect to relayer")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	_, err = relayer.SubmitMusigSignatures(context.Background(), &services.MusigSignatureMessage{
		Message:   msg,
		Pubkey:    pbk,
		Signature: sig,
	})
	return err
}
