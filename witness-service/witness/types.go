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

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// Version
	Version string
	// TransferStatus is the status of a transfer
	TransferStatus string

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
		ID() string
		GetRecorder() AbstractRecorder
		PullTransfersByHeight(blockHeight uint64) error
		PullTransfers(blockCount uint16) error
		SubmitTransfers() error
		CheckTransfers() error
		ProcessStales() error
	}

	AbstractRecorder interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		AddTransfer(tx AbstractTransfer, status TransferStatus) error
		UpsertTransfer(tx AbstractTransfer) error
		TipHeight(id string) (uint64, error)
		UpdateSyncHeight(id string, height uint64) error
		Transfer(id common.Hash) (AbstractTransfer, error)
		UnsettledTransfers() ([]string, error)
		TransfersToSubmit(string) ([]AbstractTransfer, error)
		TransfersToSettle(string) ([]AbstractTransfer, error)
		SettleTransfer(tx AbstractTransfer) error
		ConfirmTransfer(tx AbstractTransfer) error
		MarkTransferAsPending(tx AbstractTransfer) error
	}

	AbstractTransfer interface {
		Cashier() util.Address
		Token() util.Address
		CoToken() util.Address
		Index() *big.Int
		Sender() util.Address
		Recipient() util.Address
		Payload() []byte
		Amount() *big.Int
		ID() []byte
		SetID(common.Hash)
		Status() TransferStatus
		BlockHeight() uint64
		ToTypesTransfer() *types.Transfer
	}

	SignHandler func(AbstractTransfer, []byte) (common.Hash, []byte, []byte, error)
)

const (
	// TransferNew stands for a new transfer
	TransferNew TransferStatus = "new"
	// TransferPending stands for a pending transfer
	TransferPending = "pending"
	// TransferReady stands for a new transfer ready to sign
	TransferReady = "ready"
	// SubmissionConfirmed stands for a confirmed witness
	SubmissionConfirmed = "confirmed"
	// TransferSettled stands for a settled transfer
	TransferSettled = "settled"
	// TransferInvalid stands for an invalid transfer
	TransferInvalid = "invalid"

	// NoPayload is version without payload
	NoPayload Version = "no-payload"
	// Payload is version with payload
	Payload Version = "payload"
	// ToSolana is version with payload to solana
	ToSolana Version = "to-solana"
)
