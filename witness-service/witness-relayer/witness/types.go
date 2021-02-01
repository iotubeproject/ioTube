// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// TransferStatus is the status of a transfer
	TransferStatus string

	// Transfer defines a record
	Transfer struct {
		cashier     common.Address
		token       common.Address
		index       uint64
		sender      common.Address
		recipient   common.Address
		amount      *big.Int
		id          common.Hash
		status      TransferStatus
		signature   []byte
		blockHeight uint64
		txHash      common.Hash
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
		PullTransfers(blockOffset uint64, blockCount uint16) (uint64, []*Transfer, error)
	}

	eventReceipt struct {
		token     common.Address
		id        *big.Int
		sender    common.Address
		recipient common.Address
		amount    *big.Int
		fee       *big.Int
	}
)

const (
	eventName = "Receipt"

	// TransferNew stands for a new transfer
	TransferNew TransferStatus = "new"
	// WitnessSubmitted stands for a witnessed transfer
	WitnessSubmitted = "submitted"
	// SubmissionConfirmed stands for a confirmed witness
	SubmissionConfirmed = "confirmed"
	// TransferSettled stands for a settled transfer
	TransferSettled = "settled"
)
