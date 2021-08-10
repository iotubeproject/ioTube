// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type (
	tokenCashierBase struct {
		id                     string
		recorder               *Recorder
		relayerURL             string
		validatorContractAddr  common.Address
		lastProcessBlockHeight uint64
		lastPullTimestamp      time.Time
		calcEndHeight          calcEndHeightFunc
		pullTransfers          pullTransfersFunc
	}
	calcEndHeightFunc func(startHeight uint64, count uint16) (uint64, error)
	pullTransfersFunc func(startHeight uint64, endHeight uint64) ([]*Transfer, error)
)

func newTokenCashierBase(
	id string,
	recorder *Recorder,
	relayerURL string,
	validatorContractAddr common.Address,
	startBlockHeight uint64,
	calcEndHeight calcEndHeightFunc,
	pullTransfers pullTransfersFunc,
) TokenCashier {
	return &tokenCashierBase{
		id:                     id,
		recorder:               recorder,
		relayerURL:             relayerURL,
		lastProcessBlockHeight: startBlockHeight,
		validatorContractAddr:  validatorContractAddr,
		calcEndHeight:          calcEndHeight,
		pullTransfers:          pullTransfers,
		lastPullTimestamp:      time.Now(),
	}
}

func (tc *tokenCashierBase) Start(ctx context.Context) error {
	return tc.recorder.Start(ctx)
}

func (tc *tokenCashierBase) Stop(ctx context.Context) error {
	return tc.recorder.Stop(ctx)
}

func (tc *tokenCashierBase) GetRecorder() *Recorder {
	return tc.recorder
}

func (tc *tokenCashierBase) PullTransfers(count uint16) error {
	startHeight, err := tc.recorder.TipHeight(tc.id)
	if err != nil {
		return err
	}
	if startHeight < tc.lastProcessBlockHeight {
		startHeight = tc.lastProcessBlockHeight
	}
	startHeight = startHeight + 1
	endHeight, err := tc.calcEndHeight(startHeight, count)
	if err != nil {
		if tc.lastPullTimestamp.Add(3 * time.Minute).After(time.Now()) {
			log.Printf("failed to get end height with start height %d, count %d: %+v\n", startHeight, endHeight, err)
			return nil
		}
		return errors.Wrapf(err, "failed to get end height with start height %d, count %d", startHeight, count)
	}
	tc.lastPullTimestamp = time.Now()
	log.Printf("fetching events from block %d to %d for %s\n", startHeight, endHeight, tc.id)
	transfers, err := tc.pullTransfers(startHeight, endHeight)
	if err != nil {
		return errors.Wrapf(err, "failed to pull transfers from %d to %d", startHeight, endHeight)
	}
	for _, transfer := range transfers {
		if err := tc.recorder.AddTransfer(transfer); err != nil {
			return errors.Wrap(err, "failed to add transfer")
		}
	}
	tc.lastProcessBlockHeight = endHeight

	return tc.recorder.UpdateSyncHeight(tc.id, endHeight)
}

func (tc *tokenCashierBase) SubmitTransfers(sign func(*Transfer, common.Address) (common.Hash, common.Address, []byte, error)) error {
	transfersToSubmit, err := tc.recorder.TransfersToSubmit()
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	for _, transfer := range transfersToSubmit {
		id, witness, signature, err := sign(transfer, tc.validatorContractAddr)
		if err != nil {
			return err
		}
		transfer.id = id
		response, err := relayer.Submit(
			context.Background(),
			&types.Witness{
				Transfer: &types.Transfer{
					Cashier:   transfer.cashier.Bytes(),
					Token:     transfer.coToken.Bytes(),
					Index:     int64(transfer.index),
					Sender:    transfer.sender.Bytes(),
					Recipient: transfer.recipient.Bytes(),
					Amount:    transfer.amount.String(),
					Fee:       transfer.fee.String(),
				},
				Address:   witness.Bytes(),
				Signature: signature,
			},
		)
		if err != nil {
			return err
		}
		if response.Success {
			if err := tc.recorder.ConfirmTransfer(transfer); err != nil {
				return err
			}
		} else {
			log.Printf("something went wrong when submitting transfer (%s, %s, %d) for %s\n", transfer.cashier, transfer.token, transfer.index, tc.id)
		}
	}
	return nil
}

func (tc *tokenCashierBase) CheckTransfers() error {
	transfersToSettle, err := tc.recorder.TransfersToSettle()
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)

	for _, transfer := range transfersToSettle {
		response, err := relayer.Check(
			context.Background(),
			&services.CheckRequest{Id: transfer.id.Bytes()},
		)
		if err != nil {
			return err
		}
		if response.Status == services.CheckResponse_SETTLED {
			if err := tc.recorder.SettleTransfer(transfer); err != nil {
				return err
			}
		}
	}
	return nil
}
