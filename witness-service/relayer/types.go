// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"math/big"
	"time"

	soltypes "github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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
		cashier      util.Address
		token        util.Address
		index        uint64
		sender       util.Address
		txSender     util.Address
		recipient    util.Address
		amount       *big.Int
		payload      []byte
		fee          *big.Int
		id           common.Hash
		sourceTxHash []byte
		txHash       []byte
		timestamp    time.Time
		gas          uint64
		gasPrice     *big.Int
		relayer      common.Address
		nonce        uint64
		updateTime   time.Time
		status       ValidationStatusType
	}
	// Witness defines a witness structure
	Witness struct {
		addr      []byte
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

	AbstractRecorder interface {
		// Start starts the recorder
		Start(ctx context.Context) error
		// Stop stops the recorder
		Stop(ctx context.Context) error
		// AddWitness adds a witness and a transfer
		AddWitness(validator util.Address, transfer *Transfer, witness *Witness) (common.Hash, error)
		// ResetFailedTransfer resets a failed transfer
		ResetFailedTransfer(id common.Hash) error
		// Transfers returns a list of transfers
		Transfers(offset uint32, limit uint8, byUpdateTime bool, desc bool,
			queryOpts ...TransferQueryOption) ([]*Transfer, error)
		// Transfer returns a transfer by id
		Transfer(id common.Hash) (*Transfer, error)
		// Witnesses returns a list of witnesses by transfer id
		Witnesses(ids ...common.Hash) (map[common.Hash][]*Witness, error)
		// Count returns the number of transfers
		Count(opts ...TransferQueryOption) (int, error)
		// AddNewTX adds a new tx to the new tx table
		AddNewTX(height uint64, txHash []byte) error
		// NewTXs returns the new txs to be witnessed
		NewTXs(count uint32) ([]uint64, [][]byte, error)
	}

	SOLRawTransaction struct {
		id                   common.Hash
		cashier              util.Address
		token                util.Address
		index                uint64
		sender               util.Address
		txSender             util.Address
		recipient            util.Address
		amount               *big.Int
		payload              []byte
		fee                  *big.Int
		relayer              common.Address
		signature            soltypes.Signature
		lastValidBlockHeight uint64
		status               ValidationStatusType
	}
)

const (
	// WaitingForWitnesses stands for a transfer which needs more valid witnesses
	WaitingForWitnesses ValidationStatusType = "new"
	// ValidationInProcess stands for a transfer in process
	ValidationInProcess = "processing"
	// ValidationSubmitted stands for a transfer with validation submitted
	ValidationSubmitted = "validated"
	// ValidationValidationSettled stands for witness validation which has been settled
	ValidationValidationSettled = "valsettled"
	// ValidationExecuted stands for a transfer with execution submitted
	ValidationExecuted = "executed"
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
func UnmarshalTransferProto(transfer *types.Transfer, destAddrDecoder util.AddressDecoder,
) (*Transfer, error) {
	cashier, err := DecodeSourceAddrBytes(transfer.Cashier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode cashier")
	}
	token, err := destAddrDecoder.DecodeBytes(transfer.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode token")
	}
	index := uint64(transfer.Index)
	sender, err := DecodeSourceAddrBytes(transfer.Sender)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode sender")
	}
	var txSender util.Address
	if len(transfer.TxSender) > 0 {
		txSender, err = DecodeSourceAddrBytes(transfer.TxSender)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode tx sender")
		}
	}
	recipient, err := destAddrDecoder.DecodeBytes(transfer.Recipient)
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
	return &Transfer{
		cashier:      cashier,
		token:        token,
		index:        index,
		sender:       sender,
		recipient:    recipient,
		amount:       amount,
		payload:      transfer.Payload,
		fee:          fee,
		gas:          transfer.Gas,
		gasPrice:     gasPrice,
		timestamp:    transfer.Timestamp.AsTime(),
		txSender:     txSender,
		sourceTxHash: transfer.SourceTxHash,
	}, nil
}

// DecodeSourceAddrBytes decode bytes into util.Address based on the length of bytes
func DecodeSourceAddrBytes(bytes []byte) (util.Address, error) {
	switch len(bytes) {
	// ETH address
	case 20:
		return util.NewETHAddressDecoder().DecodeBytes(bytes)
	// Solana address
	case 32:
		return util.NewSOLAddressDecoder().DecodeBytes(bytes)
	default:
		return nil, errors.Errorf("invalid address length %d", len(bytes))
	}
}

// NewWitness creates a new witness struct
func NewWitness(witnessBytes []byte, signature []byte) (*Witness, error) {
	clone := make([]byte, len(signature))
	copy(clone, signature)

	return &Witness{
		addr:      witnessBytes,
		signature: signature,
	}, nil
}

func (transfer *Transfer) ID() common.Hash {
	return transfer.id
}

func (transfer *Transfer) TxHash() []byte {
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
		Recipient: transfer.recipient.Bytes(),
		Amount:    transfer.amount.String(),
		Payload:   transfer.payload,
		Fee:       transfer.fee.String(),
		Gas:       transfer.gas,
		GasPrice:  gasPrice,
		Timestamp: timestamppb.New(transfer.timestamp),
	}
}

func (w *Witness) Address() common.Address {
	return common.BytesToAddress(w.addr)
}
