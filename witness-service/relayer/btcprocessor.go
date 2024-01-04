package relayer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/btcsuite/btcwallet/wtxmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/util"
)

const (
	_limitSize      = 50
	_retryThreshold = 5
	_utxoQuery      = "https://mempool.space/testnet/api/address/%s/utxo"
	_targetBlocks   = 10
)

type BTCProcessor struct {
	signer         *MuSigner
	btcRecorder    *BTCRecorder
	recorder       *Recorder
	client         *rpcclient.Client
	chain          *chaincfg.Params
	runner         dispatcher.Runner
	witnessPubKeys []*btcec.PublicKey

	lockedOutpoints    map[wire.OutPoint]struct{}
	lockedOutpointsMtx sync.RWMutex

	ratelimiter *util.DailyResetCounter
}

func NewBTCProcessor(
	btcRecorder *BTCRecorder,
	recorder *Recorder,
	client *rpcclient.Client,
	chain *chaincfg.Params,
	interval time.Duration,
	witnessPubKeys []string,
	ratelimiter *util.DailyResetCounter,
) *BTCProcessor {
	keys := []*btcec.PublicKey{}
	for _, pkHex := range witnessPubKeys {
		pk, err := util.HexToPubkey(pkHex)
		if err != nil {
			log.Fatalln(err)
		}
		keys = append(keys, pk)
	}
	musigCfg := muSignerConfig{
		btcClient: client,
		sinkSignedTXHandler: func(txs []*wire.MsgTx) error {
			var err error
			for _, tx := range txs {
				var buf bytes.Buffer
				if err := tx.Serialize(&buf); err != nil {
					log.Printf("failed to serialize transaction %s\n", tx.TxHash().String())
					continue
				}
				if err = btcRecorder.MarkAsSigned(tx.TxHash(), buf.Bytes()); err != nil {
					log.Printf("failed to mark transaction %s as signed\n", tx.TxHash().String())
					continue
				}
			}
			return err
		},
		sinkInterval:   interval,
		witnessPubKeys: keys,
	}

	p := &BTCProcessor{
		signer:          NewMuSigner(musigCfg),
		btcRecorder:     btcRecorder,
		recorder:        recorder,
		client:          client,
		chain:           chain,
		witnessPubKeys:  keys,
		lockedOutpoints: make(map[wire.OutPoint]struct{}),
		ratelimiter:     ratelimiter,
	}

	var err error
	p.runner, err = dispatcher.NewRunner(interval, p.process)
	if err != nil {
		log.Fatalln(err)
	}
	return p
}

func (p *BTCProcessor) Start() error {
	p.ratelimiter.Start()

	if err := p.btcRecorder.Start(context.Background()); err != nil {
		return err
	}
	if err := p.signer.Start(); err != nil {
		return err
	}

	// restore the locked outpoints and resend new txs to signer
	newTXs, err := p.btcRecorder.Transfers(0, 0, false, false, StatusQueryOption(WaitingForWitnesses))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers")
	}
	for _, tx := range newTXs {
		var msgTx wire.MsgTx
		if err := msgTx.DeserializeNoWitness(bytes.NewReader(tx.txSerialized)); err != nil {
			return err
		}
		for _, in := range msgTx.TxIn {
			p.LockOutpoint(in.PreviousOutPoint)
		}
		p.signer.AddTX(&txMessage{
			tx:              &msgTx,
			transferMapping: tx.transferID,
		})
	}
	signedTXs, err := p.btcRecorder.Transfers(0, 0, false, false, StatusQueryOption(TransferSigned))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers")
	}
	submittedTXs, err := p.btcRecorder.Transfers(0, 0, false, false, StatusQueryOption(TransferSettled))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers")
	}
	for _, tx := range append(signedTXs, submittedTXs...) {
		var msgTx wire.MsgTx
		if err := msgTx.Deserialize(bytes.NewReader(tx.txSerialized)); err != nil {
			return err
		}
		for _, in := range msgTx.TxIn {
			p.LockOutpoint(in.PreviousOutPoint)
		}
	}

	return p.runner.Start()
}

func (p *BTCProcessor) Close() error {
	p.ratelimiter.Stop()

	if err := p.runner.Close(); err != nil {
		return err
	}
	if err := p.signer.Stop(); err != nil {
		return err
	}
	return p.btcRecorder.Stop(context.Background())
}

func (p *BTCProcessor) process() error {
	fmt.Println("btc process fired!")
	if err := p.ConfirmTransfers(); err != nil {
		util.LogErr(err)
	}
	if err := p.SubmitTransfers(); err != nil {
		util.LogErr(err)
	}
	return nil
}

// Confirm BTC transactions and update the status of Transfers in DB and
// upsert UTXOs in memory
func (p *BTCProcessor) ConfirmTransfers() error {
	submittedTXs, err := p.btcRecorder.Transfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	for _, rawTX := range submittedTXs {
		tx, err := p.client.GetRawTransactionVerbose(&rawTX.txHash)
		if err != nil {
			log.Printf("failed to check transaction %s, %+v\n", rawTX.txHash.String(), err)
			if err := p.resetTransaction(rawTX); err != nil {
				log.Printf("failed to reset transaction for %s\n", rawTX.txHash.String())
			}
			continue
		}
		if !isTxConfirmed(tx) {
			continue
		}
		if err := p.settleTransaction(rawTX); err != nil {
			log.Printf("failed to settle transaction %s, %+v\n", rawTX.txHash.String(), err)
			continue
		}
		log.Printf("transaction %s is settled\n", rawTX.txHash.String())
	}
	return nil
}

func (p *BTCProcessor) resetTransaction(tx *BTCRawTransaction) error {
	if tx.retryTimes <= _retryThreshold {
		return p.btcRecorder.AddRetry(tx)
	}
	if err := p.btcRecorder.MarkAsFailed(tx.txHash); err != nil {
		return err
	}
	for _, tsfID := range tx.transferID {
		if err := p.recorder.ResetFailedTransfer(tsfID); err != nil {
			return err
		}
	}
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(tx.txSerialized)); err != nil {
		return err
	}
	for _, in := range msgTx.TxIn {
		p.UnlockOutpoint(in.PreviousOutPoint)
	}
	return nil
}

func isTxConfirmed(tx *btcjson.TxRawResult) bool {
	return len(tx.BlockHash) > 0 && tx.Confirmations > 0 && tx.Blocktime > 0
}

func (p *BTCProcessor) settleTransaction(tx *BTCRawTransaction) error {
	if err := p.btcRecorder.MarkAsSettled(tx.txHash); err != nil {
		return err
	}
	for _, tsfID := range tx.transferID {
		if err := p.recorder.MarkAsSettled(tsfID, 0, time.Now()); err != nil {
			return err
		}
	}
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(tx.txSerialized)); err != nil {
		return err
	}
	for _, in := range msgTx.TxIn {
		p.UnlockOutpoint(in.PreviousOutPoint)
	}
	return nil
}

// Build BTC transactions from witnessed Transfers and send them to bitcoin network
func (p *BTCProcessor) SubmitTransfers() error {
	if err := p.buildBTCTransactions(); err != nil {
		util.LogErr(err)
	}
	if err := p.submitSignedBTCTransactions(); err != nil {
		util.LogErr(err)
	}
	return nil
}

// build BTC transactions from transfers and distribute out the txs
// to be signed by multi signers(see Musig2)
func (p *BTCProcessor) buildBTCTransactions() error {
	tsfs, err := p.collectNewTransfers()
	if err != nil {
		return err
	}
	// ratelimit when testing
	if p.ratelimiter != nil {
		tsfs = p.ratelimit(tsfs)
	}
	return p.buildBTCTXs(tsfs)
}

func (p *BTCProcessor) ratelimit(tsfs []*Transfer) []*Transfer {
	ret := make([]*Transfer, 0)
	for _, tsf := range tsfs {
		if !p.ratelimiter.Add(tsf.amount.Uint64()) {
			log.Printf("amount %s is too large\n", tsf.amount.String())
			continue
		}
		ret = append(ret, tsf)
	}
	return ret
}

// collect new transfers which have been witnessed by all witnesses
// but not processed by the processor yet
func (p *BTCProcessor) collectNewTransfers() ([]*Transfer, error) {
	newTransfers, err := p.recorder.Transfers(0, uint8(_limitSize), true, false,
		StatusQueryOption(WaitingForWitnesses))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read transfers to confirm")
	}

	tsfIDs := make([]common.Hash, 0, len(newTransfers))
	for _, transfer := range newTransfers {
		tsfIDs = append(tsfIDs, transfer.ID())
	}
	witnesses, err := p.recorder.Witnesses(tsfIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch witness")
	}

	ret := make([]*Transfer, 0)
	for _, transfer := range newTransfers {
		if p.checkReadiness(witnesses[transfer.ID()]) {
			ret = append(ret, transfer)
		}
	}
	return ret, nil
}

func (p *BTCProcessor) checkReadiness(witnesses []*Witness) bool {
	if len(witnesses) < len(p.witnessPubKeys) {
		return false
	}
	witnessAddrsMap := map[common.Address]struct{}{}
	for _, wtn := range witnesses {
		witnessAddrsMap[wtn.addr] = struct{}{}
	}

	for _, pbk := range p.witnessPubKeys {
		if _, existed := witnessAddrsMap[crypto.PubkeyToAddress(*pbk.ToECDSA())]; !existed {
			return false
		}
	}
	return true
}

// build unsigned BTC txs from UTXOs and send them to signer and lock the used UTXOs and mark transfer as InProcess
func (p *BTCProcessor) buildBTCTXs(tsfs []*Transfer) error {
	if len(tsfs) == 0 {
		return nil
	}
	eligible, err := p.findEligibleOutputs()
	if err != nil {
		return err
	}

	musigAddr, err := util.TaprootAddrFromCombinedPubkeys(p.witnessPubKeys, p.chain)
	if err != nil {
		return err
	}
	changeSource := makeChangeScriptSource(musigAddr)
	feeResp, err := p.client.EstimateSmartFee(_targetBlocks, &btcjson.EstimateModeEconomical)
	if err != nil || feeResp.FeeRate == nil {
		return errors.Wrap(err, "failed to estimate fee rate")
	}
	estimatedFeeRate, err := btcutil.NewAmount(*feeResp.FeeRate)
	if err != nil {
		return err
	}
	for _, tsf := range tsfs {
		out, err := p.txOutFromTransfer(tsf)
		if err != nil {
			log.Printf("failed to create txout for %s, %+v\n", tsf.ID().String(), err)
			continue
		}
		outs := []*wire.TxOut{out}
		inputSource := p.makeInputSource(eligible)
		tx, err := txauthor.NewUnsignedTransaction(outs, estimatedFeeRate, inputSource, changeSource)
		if err != nil {
			log.Printf("failed to create unsigned transaction for %s, %+v\n", tsf.ID().String(), err)
			for _, o := range outs {
				log.Printf("outs: %+v ", *o)
			}
			log.Printf("\n")
			continue
		}
		unsignedTX := tx.Tx
		if err := p.signer.AddTX(
			&txMessage{
				tx:              unsignedTX,
				transferMapping: map[uint64]common.Hash{0: tsf.ID()},
			},
		); err != nil {
			log.Printf("failed to add unsigned transaction for %s, %+v\n", tsf.ID().String(), err)
			return err
		}

		// add unsignedTx indo btc db
		var buf bytes.Buffer
		if err := unsignedTX.SerializeNoWitness(&buf); err != nil {
			return err
		}
		if err := p.btcRecorder.AddTransaction(
			&BTCRawTransaction{
				txHash:       unsignedTX.TxHash(),
				txSerialized: buf.Bytes(),
				status:       WaitingForWitnesses,
				transferID:   map[uint64]common.Hash{0: tsf.ID()},
				retryTimes:   0,
			},
		); err != nil {
			log.Printf("failed to add unsigned transaction for %s, %+v\n", tsf.ID().String(), err)
			continue
		}

		// lock UTXOs in the unsignedTx
		for _, in := range unsignedTX.TxIn {
			p.LockOutpoint(in.PreviousOutPoint)
		}

		// Mark the transfer as InProcess in Recorder
		if err := p.recorder.MarkAsProcessing(tsf.ID()); err != nil {
			log.Printf("failed to mark transfer %s as processing\n", tsf.ID().String())
		}
	}
	return nil
}

func (p *BTCProcessor) txOutFromTransfer(tsf *Transfer) (*wire.TxOut, error) {
	addr, ok := tsf.recipient.Address().(btcutil.Address)
	if !ok {
		return nil, errors.New("invalid recipient address")
	}
	pkscript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}
	return wire.NewTxOut(tsf.amount.Int64(), pkscript), nil
}

func (p *BTCProcessor) findEligibleOutputs() ([]wtxmgr.Credit, error) {
	musigAddr, err := util.TaprootAddrFromCombinedPubkeys(p.witnessPubKeys, p.chain)
	if err != nil {
		return nil, err
	}
	// query UTXO on Chain
	response, err := http.Get(fmt.Sprintf(_utxoQuery, musigAddr.EncodeAddress()))
	if err != nil {
		log.Printf("failed to query utxo for %s, %+v\n", musigAddr.EncodeAddress(), err)
		return nil, errors.Wrapf(err, "failed to query utxo for %s", musigAddr.EncodeAddress())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	data := string(body)
	if !gjson.Valid(data) {
		return nil, errors.New("request json format is not valid")
	}
	ret := []wtxmgr.Credit{}
	for _, req := range gjson.Parse(data).Array() {
		txid := req.Get("txid")
		vout := req.Get("vout")
		value := req.Get("value")
		confirmed := req.Get("status.confirmed")
		if !txid.Exists() || !vout.Exists() || !value.Exists() || !confirmed.Exists() {
			return nil, errors.New("request field is incomplete")
		}
		if !confirmed.Bool() {
			continue
		}

		hh, err := chainhash.NewHashFromStr(txid.String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse tx hash")
		}

		resp, err := p.client.GetRawTransaction(hh)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get raw transaction %s", hh.String())
		}
		tx := resp.MsgTx()

		ret = append(ret, wtxmgr.Credit{
			OutPoint: *wire.NewOutPoint(hh, uint32(vout.Uint())),
			Amount:   btcutil.Amount(value.Uint()),
			PkScript: tx.TxOut[vout.Uint()].PkScript,
		})
	}

	return ret, nil
}

func (p *BTCProcessor) makeInputSource(eligible []wtxmgr.Credit) txauthor.InputSource {
	// exclude locked outpoints
	unlockedCredits := make([]wtxmgr.Credit, 0, len(eligible))
	for _, credit := range eligible {
		if !p.LockedOutpoint(credit.OutPoint) {
			unlockedCredits = append(unlockedCredits, credit)
		}
	}

	// Pick largest outputs first: wallet.CoinSelectionLargest
	sort.Sort(sort.Reverse(byAmount(unlockedCredits)))

	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, len(eligible))
	currentScripts := make([][]byte, 0, len(eligible))
	currentInputValues := make([]btcutil.Amount, 0, len(eligible))

	return func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn,
		[]btcutil.Amount, [][]byte, error) {

		for currentTotal < target && len(unlockedCredits) != 0 {
			nextCredit := &unlockedCredits[0]
			unlockedCredits = unlockedCredits[1:]
			nextInput := wire.NewTxIn(&nextCredit.OutPoint, nil, nil)
			currentTotal += nextCredit.Amount
			currentInputs = append(currentInputs, nextInput)
			currentScripts = append(currentScripts, nextCredit.PkScript)
			currentInputValues = append(currentInputValues, nextCredit.Amount)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
}

// LockedOutpoint returns whether an outpoint has been marked as locked and
// should not be used as an input for created transactions.
func (p *BTCProcessor) LockedOutpoint(op wire.OutPoint) bool {
	p.lockedOutpointsMtx.RLock()
	defer p.lockedOutpointsMtx.RUnlock()

	_, locked := p.lockedOutpoints[op]
	return locked
}

// LockOutpoint marks an outpoint as locked, that is, it should not be used as
// an input for newly created transactions.
func (p *BTCProcessor) LockOutpoint(op wire.OutPoint) {
	p.lockedOutpointsMtx.Lock()
	defer p.lockedOutpointsMtx.Unlock()

	p.lockedOutpoints[op] = struct{}{}
	log.Println("Lock outpoint: ", op.String())
}

// UnlockOutpoint marks an outpoint as unlocked, that is, it may be used as an
// input for newly created transactions.
func (p *BTCProcessor) UnlockOutpoint(op wire.OutPoint) {
	p.lockedOutpointsMtx.Lock()
	defer p.lockedOutpointsMtx.Unlock()

	delete(p.lockedOutpoints, op)
}

// ResetLockedOutpoints resets the set of locked outpoints so all may be used
// as inputs for new transactions.
func (p *BTCProcessor) ResetLockedOutpoints() {
	p.lockedOutpointsMtx.Lock()
	defer p.lockedOutpointsMtx.Unlock()

	p.lockedOutpoints = map[wire.OutPoint]struct{}{}
}

func (p *BTCProcessor) ListUnsignedBTCTXWithoutNonces(excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithoutNoncesResponse, error) {
	return p.signer.UnsignedTXWithoutNonces(excludedTxs)
}

func (p *BTCProcessor) SubmitMusigNonces(req *services.MusigNonceMessage) error {
	return p.signer.AddMusigNonces(req)
}

func (p *BTCProcessor) ListUnsignedBTCTXWithNonces(excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithNoncesResponse, error) {
	return p.signer.UnsignedTXWithNonces(excludedTxs)
}

func (p *BTCProcessor) SubmitMusigSignatures(req *services.MusigSignatureMessage) error {
	return p.signer.AddMusigSignatures(req)
}

// send the raw transactions to bitcoin network
func (p *BTCProcessor) submitSignedBTCTransactions() error {
	signedTXs, err := p.btcRecorder.Transfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(TransferSigned))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	for _, signedTX := range signedTXs {
		var msgTx wire.MsgTx
		if err := msgTx.Deserialize(bytes.NewReader(signedTX.txSerialized)); err != nil {
			log.Printf("failed to deserialize transaction %s\n", signedTX.txHash.String())
			if err := p.resetTransaction(signedTX); err != nil {
				log.Printf("failed to reset transaction for %s\n", signedTX.txHash.String())
			}
			return err
		}
		_, err := p.client.SendRawTransaction(&msgTx, false)
		if err != nil {
			log.Printf("failed to submit transaction %s, %+v\n", signedTX.txHash.String(), err)
			if err := p.resetTransaction(signedTX); err != nil {
				log.Printf("failed to reset transaction for %s\n", signedTX.txHash.String())
			}
			continue
		}
		if err := p.btcRecorder.MarkAsValidated(signedTX.txHash); err != nil {
			log.Printf("failed to update validated")
		}
		for _, tsfID := range signedTX.transferID {
			if err := p.recorder.MarkAsValidated(tsfID, common.Hash(msgTx.TxHash()), common.Address{}, 0, new(big.Int).SetInt64(0)); err != nil {
				log.Printf("failed to update validated")
			}
			log.Printf("transaction %s (transfer %s) is submitted\n", signedTX.txHash.String(), tsfID.String())
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
