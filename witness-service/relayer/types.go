// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/pkg/errors"
)

type (
	// ValidationStatusType type of transfer validation status
	ValidationStatusType string
	// StatusOnChainType type of transfer status on chain
	StatusOnChainType int
	// Transfer defines a transfer structure
	Transfer struct {
		cashier    common.Address
		token      common.Address
		index      uint64
		sender     common.Address
		recipient  common.Address
		amount     *big.Int
		id         common.Hash
		txHash     common.Hash
		nonce      uint64
		updateTime time.Time
		status     ValidationStatusType
	}
	// Witness defines a witness structure
	Witness struct {
		addr      common.Address
		signature []byte
	}

	// TransferValidator defines the interface of a transfer validator
	TransferValidator interface {
		// Address returns the transfer validator contract address
		Address() common.Address
		// Check returns transfer status on chain
		Check(transfer *Transfer) (StatusOnChainType, error)
		// Submit submits validation for a transfer
		Submit(transfer *Transfer, witnesses []*Witness) (common.Hash, uint64, error)
	}
)

const (
	// waitingForWitnesses stands for a transfer which needs more valid witnesses
	waitingForWitnesses ValidationStatusType = "new"
	// validationInProcess stands for a transfer in process
	validationInProcess = "processing"
	// validationSubmitted stands for a transfer with validation submitted
	validationSubmitted = "validated"
	// transferSettled stands for a transfer which has been settled
	transferSettled = "settled"
	// validationFailed stands for the validation of a transfer failed
	validationFailed = "failed"
)

const (
	StatusOnChainUnknown StatusOnChainType = iota
	StatusOnChainNotConfirmed
	StatusOnChainRejected
	StatusOnChainNonceOverwritten
	StatusOnChainSettled
)

var errInsufficientWitnesses = errors.New("insufficient witnesses")
var errGasPriceTooHigh = errors.New("gas price is too high")

// UnmarshalTransferProto unmarshals a transfer proto
func UnmarshalTransferProto(validatorAddr common.Address, transfer *types.Transfer) (*Transfer, error) {
	cashier := common.BytesToAddress(transfer.Cashier)
	token := common.BytesToAddress(transfer.Token)
	index := uint64(transfer.Index)
	sender := common.BytesToAddress(transfer.Sender)
	recipient := common.BytesToAddress(transfer.Recipient)
	amount, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok || amount.Sign() == -1 {
		return nil, errors.Errorf("invalid amount %s", transfer.Amount)
	}
	id := crypto.Keccak256Hash(
		validatorAddr.Bytes(),
		cashier.Bytes(),
		token.Bytes(),
		math.U256Bytes(new(big.Int).SetUint64(index)),
		sender.Bytes(),
		recipient.Bytes(),
		math.U256Bytes(amount),
	)

	return &Transfer{
		cashier:   cashier,
		token:     token,
		index:     index,
		sender:    sender,
		recipient: recipient,
		amount:    amount,
		id:        id,
	}, nil
}

// NewWitness creates a new witness struct
func NewWitness(witnessAddr common.Address, signature []byte) (*Witness, error) {
	clone := make([]byte, len(signature))
	copy(clone, signature)

	return &Witness{
		addr:      witnessAddr,
		signature: signature,
	}, nil
}
