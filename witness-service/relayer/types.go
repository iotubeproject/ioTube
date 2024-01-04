// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// ValidationStatusType type of transfer validation status
	ValidationStatusType string
	// StatusOnChainType type of transfer status on chain
	StatusOnChainType int
	// Transfer defines a transfer structure
	Transfer struct {
		cashier common.Address
		token   common.Address
		// index type is changed to big.Int from u64 to support both u64 and [32]byte(U256)
		index      *big.Int
		indexType  IndexType
		sender     common.Address
		txSender   common.Address
		recipient  util.Address
		amount     *big.Int
		fee        *big.Int
		id         common.Hash
		txHash     common.Hash
		timestamp  time.Time
		gas        uint64
		gasPrice   *big.Int
		relayer    common.Address
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

	IndexType uint8

	BTCRawTransaction struct {
		txHash       chainhash.Hash
		txSerialized []byte
		status       ValidationStatusType
		transferID   map[uint64]common.Hash // transfer linked to
		retryTimes   uint8
	}

	BTCAddress struct {
		pubKey  *btcec.PublicKey
		btcAddr []byte
		ethAddr common.Address
	}
)

const (
	// WaitingForWitnesses stands for a transfer which needs more valid witnesses
	WaitingForWitnesses ValidationStatusType = "new"
	// ValidationInProcess stands for a transfer in process
	ValidationInProcess = "processing"
	// ValidationSubmitted stands for a transfer with validation submitted
	ValidationSubmitted = "validated"
	// TransferSigned stands for a transfer which has been signed
	TransferSigned = "signed"
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

const (
	LegacyIndex IndexType = iota
	BTCIndex
)

var errInsufficientWitnesses = errors.New("insufficient witnesses")
var errGasPriceTooHigh = errors.New("gas price is too high")
var errNoncritical = errors.New("error before submission")

// UnmarshalTransferProto unmarshals a transfer proto
func UnmarshalTransferProto(validatorAddr []byte, transfer *types.Transfer, addrDecoder util.AddressDecoder) (*Transfer, error) {
	cashier := common.BytesToAddress(transfer.Cashier)
	token := common.BytesToAddress(transfer.Token)
	var (
		index     *big.Int
		indexType IndexType
	)
	if transfer.Index != 0 && len(transfer.BtcIndex) > 0 {
		return nil, errors.Errorf("invalid index %d and btc index %s", transfer.Index, transfer.BtcIndex)
	} else if transfer.Index != 0 {
		index = new(big.Int).SetInt64(transfer.Index)
		indexType = LegacyIndex
	} else {
		var ok bool
		index, ok = new(big.Int).SetString(transfer.BtcIndex, 10)
		if !ok {
			return nil, errors.Errorf("invalid btc index %s", transfer.BtcIndex)
		}
		indexType = BTCIndex
	}
	sender := common.BytesToAddress(transfer.Sender)
	recipient, err := addrDecoder.DecodeBytes(transfer.Recipient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode recipient")
	}
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
	id := crypto.Keccak256Hash(
		validatorAddr,
		cashier.Bytes(),
		token.Bytes(),
		math.U256Bytes(index),
		sender.Bytes(),
		recipient.Bytes(),
		math.U256Bytes(amount),
	)

	return &Transfer{
		cashier:   cashier,
		token:     token,
		index:     new(big.Int).Set(index),
		indexType: indexType,
		sender:    sender,
		recipient: recipient,
		amount:    amount,
		fee:       fee,
		id:        id,
		gas:       transfer.Gas,
		gasPrice:  gasPrice,
		timestamp: transfer.Timestamp.AsTime(),
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
		Index:     transfer.index.Int64(),
		Sender:    transfer.sender.Bytes(),
		Recipient: transfer.recipient.Bytes(),
		Amount:    transfer.amount.String(),
		Fee:       transfer.fee.String(),
		Gas:       transfer.gas,
		GasPrice:  gasPrice,
		Timestamp: timestamppb.New(transfer.timestamp),
	}
}

func (w *Witness) Address() common.Address {
	return w.addr
}
