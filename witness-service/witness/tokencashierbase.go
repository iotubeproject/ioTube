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
		tokenPairs             TokenPairs
		relayerURL             string
		validatorContractAddr  []byte
		startBlockHeight       uint64
		lastProcessBlockHeight uint64
		lastPatrolBlockHeight  uint64
		lastPullTimestamp      time.Time
		calcConfirmHeight      calcConfirmHeightFunc
		pullTransfers          pullTransfersFunc
		idHasher               IDHasher
		signHandler            SignHandler
		hasEnoughBalance       hasEnoughBalanceFunc
		start                  startStopFunc
		stop                   startStopFunc
		disablePull            bool
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
	tokenPairs TokenPairs,
	relayerURL string,
	validatorContractAddr []byte,
	startBlockHeight uint64,
	calcConfirmHeight calcConfirmHeightFunc,
	pullTransfers pullTransfersFunc,
	idHasher IDHasher,
	signHandler SignHandler,
	hasEnoughBalance hasEnoughBalanceFunc,
	start startStopFunc,
	stop startStopFunc,
	disablePull bool,
) TokenCashier {
	return &tokenCashierBase{
		id:                     id,
		cashierContractAddr:    cashierContractAddr,
		previousCashierAddr:    previousCashierAddr,
		recorder:               recorder,
		tokenPairs:             tokenPairs,
		relayerURL:             relayerURL,
		startBlockHeight:       startBlockHeight,
		lastProcessBlockHeight: startBlockHeight,
		validatorContractAddr:  validatorContractAddr,
		calcConfirmHeight:      calcConfirmHeight,
		pullTransfers:          pullTransfers,
		idHasher:               idHasher,
		signHandler:            signHandler,
		hasEnoughBalance:       hasEnoughBalance,
		lastPullTimestamp:      time.Now(),
		start:                  start,
		stop:                   stop,
		disablePull:            disablePull,
	}
}

func (tc *tokenCashierBase) Start(ctx context.Context) error {
	if err := tc.recorder.Start(ctx); err != nil {
		return err
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

func (tc *tokenCashierBase) PullTransfers(count uint16) error {
	if tc.disablePull {
		return tc.fetchTransfers(count)
	}
	if count == 0 {
		count = 1
	}
	patrolSize := uint64(count) * 6
	if patrolSize > 1000 {
		patrolSize = 1000
	}
	startHeight, err := tc.recorder.TipHeight(tc.id)
	if err != nil {
		return err
	}
	startHeight = max(startHeight+1, tc.startBlockHeight)
	if tc.lastPatrolBlockHeight == 0 && startHeight > patrolSize {
		tc.lastPatrolBlockHeight = max(startHeight-patrolSize, tc.startBlockHeight)
	}
	confirmHeight, endHeight, err := tc.calcConfirmHeight(startHeight, count)
	if err != nil {
		if tc.lastPullTimestamp.Add(3 * time.Minute).After(time.Now()) {
			log.Printf("failed to get end height with start height %d, count %d: %+v\n", startHeight, confirmHeight, err)
			return nil
		}
		return errors.Wrapf(err, "failed to get end height and tip height with start height %d, count %d", startHeight, count)
	}
	var shouldPatrol bool
	switch {
	case endHeight < startHeight-1:
		return errors.Errorf("end height %d is less than start height %d - 1", endHeight, startHeight)
	case endHeight == startHeight-1:
		if endHeight > tc.lastPatrolBlockHeight+patrolSize {
			shouldPatrol = true
		} else {
			return nil
		}
	case endHeight > startHeight-1:
		if startHeight > tc.lastPatrolBlockHeight+patrolSize {
			shouldPatrol = true
			endHeight = startHeight
		}
	}
	var transfers []AbstractTransfer
	if shouldPatrol {
		log.Printf("fetching events from block %d to %d for %s with patrol\n", tc.lastPatrolBlockHeight, endHeight, tc.id)
		transfers, err = tc.pullTransfers(tc.lastPatrolBlockHeight, endHeight)
		if err != nil {
			return errors.Wrapf(err, "failed to pull transfers from %d to %d with patrol", tc.lastPatrolBlockHeight, endHeight)
		}
		tc.lastPatrolBlockHeight = endHeight
	} else {
		// log.Printf("fetching events from block %d to %d for %s\n", startHeight, endHeight, tc.id)
		transfers, err = tc.pullTransfers(startHeight, endHeight)
		if err != nil {
			return errors.Wrapf(err, "failed to pull transfers from %d to %d", startHeight, endHeight)
		}
	}
	tc.lastPullTimestamp = time.Now()
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

	endHeight = min(confirmHeight, endHeight)
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
		id, err := tc.idHasher(transfer, tc.validatorContractAddr)
		if err != nil {
			return err
		}
		transfer.SetID(id)
		pubkey, signature, err := tc.signHandler(id.Bytes())
		if err != nil {
			return err
		}
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
	response, err := relayer.StaleHeights(context.Background(), &services.StaleHeightsRequest{
		Cashier: tc.cashierContractAddr.Bytes(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch stale heights")
	}
	if tc.previousCashierAddr != nil {
		previousResponse, err := relayer.StaleHeights(context.Background(), &services.StaleHeightsRequest{
			Cashier: tc.previousCashierAddr.Bytes(),
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

func (tc *tokenCashierBase) RefreshTokenPairs() error {
	if tc.tokenPairs == nil {
		return nil
	}
	return tc.tokenPairs.Update()
}
