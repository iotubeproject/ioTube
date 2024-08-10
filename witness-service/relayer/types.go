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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	// ValidationStatusType type of transfer validation status
	ValidationStatusType string
	// StatusOnChainType type of transfer status on chain
	StatusOnChainType int
	// Transfer defines a transfer structure
	Transfer struct {
		cashier     common.Address
		token       common.Address
		index       uint64
		sender      common.Address
		txSender    common.Address
		recipient   common.Address
		amount      *big.Int
		payload     []byte
		fee         *big.Int
		blockHeight uint64
		id          common.Hash
		txHash      common.Hash
		timestamp   time.Time
		gas         uint64
		gasPrice    *big.Int
		relayer     common.Address
		nonce       uint64
		updateTime  time.Time
		status      ValidationStatusType
	}
	// Witness defines a witness structure
	Witness struct {
		addr      common.Address
		signature []byte
	}

	// BonusSender defines the interface of a bonus sender
	BonusSender interface {
		// SendBonus sends bonus to a transfer
		SendBonus(transfer *Transfer) error
		// Size returns the number of senders
		Size() int
	}

	// TransferValidator defines the interface of a transfer validator
	TransferValidator interface {
		// Size returns the number of relayers
		Size() int
		// Address returns the transfer validator contract address
		Address() common.Address
		// Check returns transfer status on chain
		Check(transfer *Transfer) (StatusOnChainType, error)
		// Submit submits validation for a transfer
		Submit(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error)
		// SpeedUp resubmits validation with higher gas price
		SpeedUp(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error)
	}
)

const (
	// WaitingForWitnesses stands for a transfer which needs more valid witnesses
	WaitingForWitnesses ValidationStatusType = "new"
	// ValidationInProcess stands for a transfer in process
	ValidationInProcess = "processing"
	// ValidationSubmitted stands for a transfer with validation submitted
	ValidationSubmitted = "validated"
	// BonusPending stands for a transfer with pending bonus
	BonusPending = "bonus"
	// TransferSettled stands for a transfer which has been settled
	TransferSettled = "settled"
	// ValidationFailed stands for the validation of a transfer failed
	ValidationFailed = "failed"
	// ValidationRejected stands for the validation of a transfer is rejected
	ValidationRejected = "rejected"
)

const (
	StatusOnChainUnknown StatusOnChainType = iota
	StatusOnChainNotConfirmed
	StatusOnChainNeedSpeedUp
	StatusOnChainRejected
	StatusOnChainNonceOverwritten
	StatusOnChainSettled
)

var errInsufficientWitnesses = errors.New("insufficient witnesses")
var errGasPriceTooHigh = errors.New("gas price is too high")
var errNoncritical = errors.New("error before submission")

// UnmarshalTransferProto unmarshals a transfer proto
func UnmarshalTransferProto(transfer *types.Transfer) (*Transfer, error) {
	cashier := common.BytesToAddress(transfer.Cashier)
	token := common.BytesToAddress(transfer.Token)
	index := uint64(transfer.Index)
	sender := common.BytesToAddress(transfer.Sender)
	var txSender common.Address
	if len(transfer.TxSender) > 0 {
		txSender = common.BytesToAddress(transfer.TxSender)
	}
	recipient := common.BytesToAddress(transfer.Recipient)
	amount, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok || amount.Sign() == -1 {
		return nil, errors.Errorf("invalid amount %s", transfer.Amount)
	}
	fee, ok := new(big.Int).SetString(transfer.Fee, 10)
	if !ok || fee.Sign() == -1 {
		fee = big.NewInt(0)
	}
	var gasPrice *big.Int
	if transfer.GasPrice != "" {
		gasPrice, ok = new(big.Int).SetString(transfer.GasPrice, 10)
		if !ok || gasPrice.Sign() == -1 {
			return nil, errors.Errorf("invalid gas price %s", transfer.GasPrice)
		}
	}

	return &Transfer{
		cashier:     cashier,
		token:       token,
		index:       index,
		sender:      sender,
		recipient:   recipient,
		amount:      amount,
		fee:         fee,
		gas:         transfer.Gas,
		gasPrice:    gasPrice,
		timestamp:   transfer.Timestamp.AsTime(),
		txSender:    txSender,
		blockHeight: transfer.BlockHeight,
		payload:     transfer.Payload,
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

func (transfer *Transfer) GenID(validatorAddr common.Address) {
	transfer.id = crypto.Keccak256Hash(
		validatorAddr.Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		math.U256Bytes(new(big.Int).SetUint64(transfer.index)),
		transfer.sender.Bytes(),
		transfer.recipient.Bytes(),
		math.U256Bytes(transfer.amount),
		transfer.payload,
	)
}

func (transfer *Transfer) ID() common.Hash {
	return transfer.id
}

func (transfer *Transfer) TxHash() common.Hash {
	return transfer.txHash
}

func (transfer *Transfer) Status() ValidationStatusType {
	return transfer.status
}

func (transfer *Transfer) ToTypesTransfer() *types.Transfer {
	gasPrice := "0"
	if transfer.gasPrice != nil {
		gasPrice = transfer.gasPrice.String()
	}

	return &types.Transfer{
		Cashier:   transfer.cashier.Bytes(),
		Token:     transfer.token.Bytes(),
		Index:     int64(transfer.index),
		Sender:    transfer.sender.Bytes(),
		TxSender:  transfer.txSender.Bytes(),
		Recipient: transfer.recipient.Bytes(),
		Amount:    transfer.amount.String(),
		Fee:       transfer.fee.String(),
		Gas:       transfer.gas,
		GasPrice:  gasPrice,
		Timestamp: timestamppb.New(transfer.timestamp),
		Payload:   transfer.payload,
	}
}

func (w *Witness) Address() common.Address {
	return w.addr
}
