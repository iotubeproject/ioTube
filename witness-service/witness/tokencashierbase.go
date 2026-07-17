// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_ReceiptEventTopic, _TransferEventTopic common.Hash
	_ReceiptEventTopicWithPayload           common.Hash
	_ZeroHash                               = common.Hash{}
)

func init() {
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		log.Panicf("failed to decode token cashier abi, %+v", err)
	}
	_ReceiptEventTopic = tokenCashierABI.Events["Receipt"].ID
	tokenCashierWithPayloadABI, err := abi.JSON(strings.NewReader(contract.TokenCashierWithPayloadABI))
	if err != nil {
		log.Panicf("failed to decode token cashier abi, %+v", err)
	}
	_ReceiptEventTopicWithPayload = tokenCashierWithPayloadABI.Events["Receipt"].ID
	erc20ABI, err := abi.JSON(strings.NewReader(contract.CrosschainERC20ABI))
	if err != nil {
		log.Panicf("failed to decode erc20 abi, %+v", err)
	}
	_TransferEventTopic = erc20ABI.Events["Transfer"].ID
}

type (
	tokenCashierBase struct {
		id                     string
		cashierContractAddr    util.Address
		previousCashierAddr    util.Address
		recorder               AbstractRecorder
		relayerURL             string
		validatorContractAddr  []byte
		startBlockHeight       uint64
		lastProcessBlockHeight uint64
		lastPatrolBlockHeight  uint64
		lastPullTimestamp      time.Time
		calcConfirmHeight      calcConfirmHeightFunc
		pullTransfers          pullTransfersFunc
		signHandler            SignHandler
		signerAddr             []byte
		hasEnoughBalance       hasEnoughBalanceFunc
		start                  startStopFunc
		stop                   startStopFunc
		disablePull            bool
		useFinalizedBlock      bool
	}
	calcConfirmHeightFunc func(startHeight uint64, count uint16) (uint64, uint64, error)
	pullTransfersFunc     func(startHeight uint64, endHeight uint64) ([]AbstractTransfer, error)
	hasEnoughBalanceFunc  func(token util.Address, amount *big.Int) bool
	startStopFunc         func(context.Context) error
)

func newTokenCashierBase(
	id string,
	cashierContractAddr util.Address,
	previousCashierAddr util.Address,
	recorder AbstractRecorder,
	relayerURL string,
	validatorContractAddr []byte,
	startBlockHeight uint64,
	calcConfirmHeight calcConfirmHeightFunc,
	pullTransfers pullTransfersFunc,
	signHandler SignHandler,
	hasEnoughBalance hasEnoughBalanceFunc,
	start startStopFunc,
	stop startStopFunc,
	disablePull bool,
	useFinalizedBlock bool,
) TokenCashier {
	// derive the witness's own address once (independent of any transfer) so the
	// liveness heartbeat can report it even when there is no traffic to sign.
	var signerAddr []byte
	if signHandler != nil {
		if _, addr, _, err := signHandler(nil, validatorContractAddr); err == nil {
			signerAddr = addr
		}
	}
	return &tokenCashierBase{
		id:                     id,
		cashierContractAddr:    cashierContractAddr,
		previousCashierAddr:    previousCashierAddr,
		recorder:               recorder,
		relayerURL:             relayerURL,
		startBlockHeight:       startBlockHeight,
		lastProcessBlockHeight: startBlockHeight,
		validatorContractAddr:  validatorContractAddr,
		calcConfirmHeight:      calcConfirmHeight,
		pullTransfers:          pullTransfers,
		signHandler:            signHandler,
		signerAddr:             signerAddr,
		hasEnoughBalance:       hasEnoughBalance,
		lastPullTimestamp:      time.Now(),
		start:                  start,
		stop:                   stop,
		disablePull:            disablePull,
		useFinalizedBlock:      useFinalizedBlock,
	}
}

func (tc *tokenCashierBase) Start(ctx context.Context) error {
	if err := tc.recorder.Start(ctx); err != nil {
		return err
	}
	if tc.useFinalizedBlock {
		// Self-heal the DB before pulling: if a previous latest-minus-N run
		// advanced the sync height past the finalized confirmed tip, roll it
		// back so the witness does not fail closed. Runs after recorder.Start
		// (DB open + tables created) and before the pull loop.
		if err := tc.reconcileConfirmedTip(); err != nil {
			return err
		}
	}
	return tc.start(ctx)
}

func (tc *tokenCashierBase) Stop(ctx context.Context) error {
	if err := tc.stop(ctx); err != nil {
		return err
	}
	return tc.recorder.Stop(ctx)
}

func (tc *tokenCashierBase) ID() string {
	return tc.id
}

func (tc *tokenCashierBase) GetRecorder() AbstractRecorder {
	return tc.recorder
}

func (tc *tokenCashierBase) PullTransfersByHeight(height uint64) error {
	transfers, err := tc.pullTransfers(height, height)
	if err != nil {
		return errors.Wrapf(err, "failed to pull transfers for %d", height)
	}
	tip, err := tc.recorder.TipHeight(tc.id)
	if err != nil {
		return err
	}
	if tip < height {
		return errors.Errorf("invalid height %d is larger than tip %d", height, tip)
	}
	for _, transfer := range transfers {
		if err := tc.recorder.UpsertTransfer(transfer); err != nil {
			return errors.Wrap(err, "failed to add transfer")
		}
	}
	return nil
}

// reconcileConfirmedTip self-heals the DB when finalized mode is enabled on a
// witness whose sync height was already advanced past the finalized confirmed
// tip by a previous (probabilistic latest-minus-N) run. Without it, PullTransfers
// would return a regression error every cycle and the witness would fail closed
// until an operator intervened. It rolls the sync height back to the confirmed
// tip and drops the not-yet-signed rows above it so they are re-derived from
// the finalized chain on the next scan. If any committed or possibly-signed
// rows (confirmed, settled, or ready — a ready row may already be signed and at
// the relayer) sit above the confirmed tip, it refuses to proceed and returns
// an error for manual review, because dropping them could desync the relayer or
// mask a real reorg. It is a no-op in steady state (sync height <= confirmed
// tip) and for non-ethereum recorders.
// Called from Start (after recorder.Start) only when useFinalizedBlock is set.
func (tc *tokenCashierBase) reconcileConfirmedTip() error {
	rec, ok := tc.recorder.(*Recorder)
	if !ok {
		// Only the ethereum recorder participates in finalized mode.
		return nil
	}
	// confirmHeight in finalized mode is finalized-minus-margin and does not
	// depend on the start height, so any start height yields the same value.
	confirmHeight, _, err := tc.calcConfirmHeight(0, 1)
	if err != nil {
		return errors.Wrap(err, "failed to compute confirmed height for reconciliation")
	}
	// Sync height lives in cashier_meta, keyed by the route id.
	syncHeight, err := tc.recorder.TipHeight(tc.id)
	if err != nil {
		return errors.Wrap(err, "failed to read sync height for reconciliation")
	}
	if syncHeight <= confirmHeight {
		return nil
	}
	// Transfer rows are keyed by the cashier contract address(es) (current and,
	// if set, the previous cashier), not by the route id.
	cashiers := []string{tc.cashierContractAddr.String()}
	if tc.previousCashierAddr != nil {
		cashiers = append(cashiers, tc.previousCashierAddr.String())
	}
	committed, minHeight, maxHeight, err := rec.committedTransfersAboveHeight(cashiers, confirmHeight)
	if err != nil {
		return errors.Wrap(err, "failed to inspect committed transfers above confirmed height")
	}
	if committed > 0 {
		return errors.Errorf(
			"cannot enable finalized mode for %s: %d committed or possibly-signed transfer(s) at blocks %d-%d sit above the finalized confirmed height %d; "+
				"they were processed under weaker confirmation settings and must be reviewed manually before enabling finalized mode",
			tc.id, committed, minHeight, maxHeight, confirmHeight,
		)
	}
	deleted, err := rec.rollbackToConfirmedTip(cashiers, tc.id, confirmHeight)
	if err != nil {
		return errors.Wrap(err, "failed to roll back to confirmed height")
	}
	// Also rewind the in-memory floor, otherwise PullTransfers uses
	// max(TipHeight, lastProcessBlockHeight) and would start above confirmHeight
	// again (e.g. when the route's startBlockHeight was set above the tip),
	// keeping the regression error despite the DB rollback.
	tc.lastProcessBlockHeight = confirmHeight
	log.Printf("reconciled %s for finalized mode: rolled sync height %d -> %d, dropped %d uncommitted transfer(s) above the confirmed tip\n",
		tc.id, syncHeight, confirmHeight, deleted)
	return nil
}

func (tc *tokenCashierBase) PullTransfers(count uint16) error {
	if tc.disablePull {
		return tc.fetchTransfers(count)
	}
	startHeight, err := tc.recorder.TipHeight(tc.id)
	if err != nil {
		return err
	}
	if startHeight < tc.lastProcessBlockHeight {
		startHeight = tc.lastProcessBlockHeight
	}
	if count == 0 {
		count = 1
	}
	patrolSize := uint64(count) * 6
	if patrolSize > 1000 {
		patrolSize = 1000
	}
	if tc.lastPatrolBlockHeight == 0 && startHeight > patrolSize {
		tc.lastPatrolBlockHeight = startHeight - patrolSize
		if tc.lastPatrolBlockHeight < tc.startBlockHeight {
			tc.lastPatrolBlockHeight = tc.startBlockHeight
		}
	}
	startHeight = startHeight + 1
	confirmHeight, endHeight, err := tc.calcConfirmHeight(startHeight, count)
	if err != nil {
		if tc.lastPullTimestamp.Add(3 * time.Minute).After(time.Now()) {
			log.Printf("failed to get end height with start height %d, count %d: %+v\n", startHeight, confirmHeight, err)
			return nil
		}
		return errors.Wrapf(err, "failed to get end height and tip height with start height %d, count %d", startHeight, count)
	}
	if confirmHeight+1 == startHeight {
		// Exactly caught up: the confirmed tip equals the last synced height, so there
		// are no new confirmed blocks this cycle (e.g. the chain's finalized height has
		// not advanced / finality is stalled). This is not a failure: refresh
		// lastPullTimestamp (a healthy poll, so a later transient RPC error still gets
		// the grace window above) and return nil so the caller still runs
		// SubmitTransfers/CheckTransfers for already-confirmed transfers instead of
		// skipping them (service.process skips those whenever PullTransfers errors).
		tc.lastPullTimestamp = time.Now()
		return nil
	}
	if confirmHeight < startHeight {
		// The confirmed tip has regressed below the last synced height by more
		// than routine RPC jitter (finalized mode clamps that to a monotonic
		// high-water mark upstream). The remaining causes are a genuine, large
		// regression: the DB was synced past the current confirmed tip (e.g.
		// switching to finalized mode after syncing further under latest-minus-N),
		// or a serious RPC fault. Do NOT treat this as a benign no-op — return an
		// error so service.process skips submissions this cycle rather than
		// signing already-ready rows that now sit above a confirmed tip that
		// regressed.
		return errors.Errorf("confirm height %d regressed below sync height %d", confirmHeight, startHeight-1)
	}
	var transfers []AbstractTransfer
	tc.lastPullTimestamp = time.Now()
	if startHeight > tc.lastPatrolBlockHeight+patrolSize {
		log.Printf("fetching events from block %d to %d for %s with patrol\n", tc.lastPatrolBlockHeight, startHeight, tc.id)
		transfers, err = tc.pullTransfers(tc.lastPatrolBlockHeight, startHeight)
		if err != nil {
			return errors.Wrapf(err, "failed to pull transfers from %d to %d with patrol", tc.lastPatrolBlockHeight, startHeight)
		}
		tc.lastPatrolBlockHeight = startHeight
		endHeight = startHeight
	} else {
		// log.Printf("fetching events from block %d to %d for %s\n", startHeight, endHeight, tc.id)
		transfers, err = tc.pullTransfers(startHeight, endHeight)
		if err != nil {
			return errors.Wrapf(err, "failed to pull transfers from %d to %d", startHeight, endHeight)
		}
	}
	for _, transfer := range transfers {
		if transfer.BlockHeight() <= confirmHeight {
			if err := tc.recorder.UpsertTransfer(transfer); err != nil {
				return errors.Wrap(err, "failed to upsert transfer")
			}
		} else {
			status := TransferNew
			if transfer.Amount().Sign() != 1 {
				status = TransferInvalid
				log.Printf("amount %d should be larger than 0 for new transfer %s\n", transfer.Amount(), transfer.ID())
			}
			if err := tc.recorder.AddTransfer(transfer, status); err != nil {
				return errors.Wrap(err, "failed to add transfer")
			}
		}
	}
	if confirmHeight < endHeight {
		endHeight = confirmHeight
	}
	tc.lastProcessBlockHeight = endHeight

	if err := tc.recorder.UpdateSyncHeight(tc.id, endHeight); err != nil {
		return errors.Wrap(err, "failed to update sync height")
	}
	return nil
}

func (tc *tokenCashierBase) SubmitTransfers() error {
	if tc.signHandler == nil {
		return nil
	}
	transfersToSubmit, err := tc.recorder.TransfersToSubmit(tc.cashierContractAddr.String())
	if err != nil {
		return err
	}
	if tc.previousCashierAddr != nil {
		transfersFromPreviousCashier, err := tc.recorder.TransfersToSubmit(tc.previousCashierAddr.String())
		if err != nil {
			return err
		}
		transfersToSubmit = append(transfersToSubmit, transfersFromPreviousCashier...)
	}
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	for _, transfer := range transfersToSubmit {
		if !tc.hasEnoughBalance(transfer.Token(), transfer.Amount()) {
			return errors.Errorf("not enough balance for token %s", transfer.Token())
		}
		id, pubkey, signature, err := tc.signHandler(transfer, tc.validatorContractAddr)
		if err != nil {
			return err
		}
		transfer.SetID(id)
		if signature == nil {
			continue
		}
		var witness *types.Witness
		if transfer.Status() == TransferReady {
			witness = &types.Witness{
				Transfer:  transfer.ToTypesTransfer(),
				Address:   pubkey,
				Signature: signature,
			}
		} else {
			witness = &types.Witness{
				Transfer:  transfer.ToTypesTransfer(),
				Address:   pubkey,
				Signature: []byte{},
			}
		}
		response, err := relayer.Submit(context.Background(), witness)
		if err != nil {
			return err
		}
		if !response.Success {
			log.Printf("something went wrong when submitting transfer (%s, %s, %s) for %s\n", transfer.Cashier(), transfer.Token(), transfer.Index().String(), id)
			continue
		}
		if transfer.Status() == TransferReady {
			if err := tc.recorder.ConfirmTransfer(transfer); err != nil {
				return err
			}
		} else {
			if err := tc.recorder.MarkTransferAsPending(transfer); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tc *tokenCashierBase) ProcessStales() error {
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	// tipHeight is the heartbeat's reported scan tip. For disable-pull cashiers
	// PullTransfers never advances lastProcessBlockHeight, so reporting it would
	// pin a misleading fixed height forever — send 0 (no scan tip) in that case;
	// the witness address alone still carries the liveness signal.
	var tipHeight uint64
	if !tc.disablePull {
		tipHeight = tc.lastProcessBlockHeight
	}
	response, err := relayer.StaleHeights(context.Background(), &services.StaleHeightsRequest{
		Cashier:     tc.cashierContractAddr.Bytes(),
		WitnessAddr: tc.signerAddr,
		TipHeight:   tipHeight,
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch stale heights")
	}
	if tc.previousCashierAddr != nil {
		previousResponse, err := relayer.StaleHeights(context.Background(), &services.StaleHeightsRequest{
			Cashier:     tc.previousCashierAddr.Bytes(),
			WitnessAddr: tc.signerAddr,
			TipHeight:   tipHeight,
		})
		if err != nil {
			return errors.Wrap(err, "failed to fetch stale heights from previous cashier")
		}
		response.Heights = append(response.Heights, previousResponse.Heights...)
	}
	for _, height := range response.Heights {
		if err := tc.PullTransfersByHeight(height); err != nil {
			return errors.Wrap(err, "failed to pull transfers by height")
		}
	}
	return nil
}

func (tc *tokenCashierBase) CheckTransfers() error {
	transfersToSettle, err := tc.recorder.TransfersToSettle(tc.cashierContractAddr.String())
	if err != nil {
		return errors.Wrap(err, "failed to fetch transfers to settle")
	}
	if tc.previousCashierAddr != nil {
		transfersFromPreviousCashier, err := tc.recorder.TransfersToSettle(tc.previousCashierAddr.String())
		if err != nil {
			return errors.Wrap(err, "failed to fetch transfers from previous cashier to settle")
		}
		transfersToSettle = append(transfersToSettle, transfersFromPreviousCashier...)
	}
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)

	for _, transfer := range transfersToSettle {
		response, err := relayer.Check(
			context.Background(),
			&services.CheckRequest{Id: transfer.ID()},
		)
		if err != nil {
			return errors.Wrap(err, "failed to check with relayer")
		}
		if response.Status == services.Status_SETTLED {
			if err := tc.recorder.SettleTransfer(transfer); err != nil {
				return errors.Wrap(err, "failed to settle transfer")
			}
		}
	}
	return nil
}

func (tc *tokenCashierBase) fetchTransfers(count uint16) error {
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	response, err := relayer.ListNewTX(context.Background(),
		&services.ListNewTXRequest{
			Count: uint32(count),
		})
	if err != nil {
		return errors.Wrap(err, "failed to list new tx")
	}

	tsfs, err := tc.recorder.UnsettledTransfers()
	if err != nil {
		return errors.Wrap(err, "failed to get unsettled transfers")
	}
	var (
		fetchHeights   = make(map[uint64]struct{})
		foundTransfers = make(map[string]struct{}, len(tsfs))
	)
	for _, t := range tsfs {
		foundTransfers[t] = struct{}{}
	}
	for _, tx := range response.Txs {
		if len(tx.TxHash) != 64 {
			log.Printf("invalid tx hash %x, skipping\n", tx.TxHash)
			continue
		}
		hash := hex.EncodeToString(tx.TxHash)
		if _, exist := foundTransfers[hash]; exist {
			continue
		}
		fetchHeights[tx.Height] = struct{}{}
	}

	for height := range fetchHeights {
		transfers, err := tc.pullTransfers(height, height)
		if err != nil {
			log.Printf("failed to pull transfers for height %d: %+v\n", height, err)
			continue
		}
		for _, transfer := range transfers {
			if err := tc.recorder.AddTransfer(transfer, TransferReady); err != nil {
				return errors.Wrap(err, "failed to add transfer")
			}
		}
	}
	return nil
}
