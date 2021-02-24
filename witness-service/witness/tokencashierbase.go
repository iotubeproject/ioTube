// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

type (
	tokenCashierBase struct {
		id                     string
		recorder               *Recorder
		validatorContractAddr  common.Address
		lastProcessBlockHeight uint64
		calcEndHeight          calcEndHeightFunc
		pullTransfers          pullTransfersFunc
	}
	calcEndHeightFunc func(startHeight uint64, count uint16) (uint64, error)
	pullTransfersFunc func(startHeight uint64, endHeight uint64) ([]*Transfer, error)
)

func newTokenCashierBase(
	id string,
	recorder *Recorder,
	validatorContractAddr common.Address,
	startBlockHeight uint64,
	calcEndHeight calcEndHeightFunc,
	pullTransfers pullTransfersFunc,
) TokenCashier {
	return &tokenCashierBase{
		id:                     id,
		recorder:               recorder,
		lastProcessBlockHeight: startBlockHeight,
		validatorContractAddr:  validatorContractAddr,
		calcEndHeight:          calcEndHeight,
		pullTransfers:          pullTransfers,
	}
}

func (tc *tokenCashierBase) Start(ctx context.Context) error {
	return tc.recorder.Start(ctx)
}

func (tc *tokenCashierBase) Stop(ctx context.Context) error {
	return tc.recorder.Stop(ctx)
}

func (tc *tokenCashierBase) PullTransfers(count uint16) error {
	startHeight, err := tc.recorder.TipHeight()
	if err != nil {
		return err
	}
	if startHeight < tc.lastProcessBlockHeight {
		startHeight = tc.lastProcessBlockHeight
	}
	endHeight, err := tc.calcEndHeight(startHeight, count)
	if err != nil {
		return err
	}
	log.Printf("fetching events from block %d for %s\n", startHeight, tc.id)
	transfers, err := tc.pullTransfers(startHeight, endHeight)
	if err != nil {
		return err
	}
	for _, transfer := range transfers {
		if err := tc.recorder.AddTransfer(transfer); err != nil {
			return err
		}
	}
	tc.lastProcessBlockHeight = endHeight

	return nil
}

func (tc *tokenCashierBase) SubmitTransfers(submit func(*Transfer, common.Address) (bool, error)) error {
	transfersToSubmit, err := tc.recorder.TransfersToSubmit()
	if err != nil {
		return err
	}
	for _, transfer := range transfersToSubmit {
		succeed, err := submit(transfer, tc.validatorContractAddr)
		if err != nil {
			return err
		}
		if succeed {
			if err := tc.recorder.ConfirmTransfer(transfer); err != nil {
				return err
			}
		} else {
			log.Printf("something went wrong when submitting transfer (%s, %s, %d) for %s\n", transfer.cashier, transfer.token, transfer.index, tc.id)
		}
	}
	return nil
}

func (tc *tokenCashierBase) CheckTransfers(check func(*Transfer) (bool, error)) error {
	transfersToSettle, err := tc.recorder.TransfersToSettle()
	if err != nil {
		return err
	}

	for _, transfer := range transfersToSettle {
		settled, err := check(transfer)
		if err != nil {
			return err
		}
		if settled {
			if err := tc.recorder.SettleTransfer(transfer); err != nil {
				return err
			}
		}
	}
	return nil
}
