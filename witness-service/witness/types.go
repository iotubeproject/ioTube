// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// TransferStatus is the status of a transfer
	TransferStatus string

	// Transfer defines a record
	Transfer struct {
		cashier     common.Address
		token       common.Address
		coToken     common.Address
		index       uint64
		sender      common.Address
		recipient   common.Address
		amount      *big.Int
		fee         *big.Int
		id          common.Hash
		status      TransferStatus
		blockHeight uint64
		txHash      common.Hash
		timestamp   time.Time
		gas         uint64
		gasPrice    *big.Int
	}

	// Service manages to exchange iotex coin to ERC20 token on ethereum
	Service interface {
		// Start starts the service
		Start(context.Context) error
		// Stop stops the service
		Stop(context.Context) error
	}

	// TokenCashier defines the interface to pull transfers from chain in a block range
	TokenCashier interface {
		Start(context.Context) error
		Stop(context.Context) error
		GetRecorder() *Recorder
		PullTransfersByHeight(blockHeight uint64) error
		PullTransfers(blockCount uint16) error
		SubmitTransfers(func(*Transfer, common.Address) (common.Hash, common.Address, []byte, error)) error
		CheckTransfers() error
	}
)

const (
	// TransferNew stands for a new transfer
	TransferNew TransferStatus = "new"
	// TransferReady stands for a new transfer ready to sign
	TransferReady = "ready"
	// WitnessSubmitted stands for a witnessed transfer
	WitnessSubmitted = "submitted"
	// SubmissionConfirmed stands for a confirmed witness
	SubmissionConfirmed = "confirmed"
	// TransferSettled stands for a settled transfer
	TransferSettled = "settled"
)
