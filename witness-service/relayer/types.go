// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"sort"
	"time"

	soltypes "github.com/blocto/solana-go-sdk/types"
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
		cashier      util.Address
		token        util.Address
		index        uint64
		sender       util.Address
		txSender     util.Address
		recipient    util.Address
		ataOwner     util.Address
		amount       *big.Int
		payload      []byte
		fee          *big.Int
		blockHeight  uint64
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

	// WitnessManager defines the interface for managing witness candidates
	WitnessManager interface {
		// Size returns the number of relayers
		Size() int
		// Address returns the witness manager contract address
		Address() common.Address
		// Check returns the status of witness candidates on chain
		Check(candidates *WitnessCandidates) (StatusOnChainType, error)
		// Submit submits witness candidate updates (additions and removals)
		Submit(candidates *WitnessCandidates, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error)
		// SpeedUp speeds up witness candidates
		SpeedUp(candidates *WitnessCandidates, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error)
	}

	SOLRawTransaction struct {
		id                   common.Hash
		cashier              util.Address
		token                util.Address
		index                uint64
		sender               util.Address
		txSender             util.Address
		recipient            util.Address
		ataOwner             util.Address
		amount               *big.Int
		payload              []byte
		fee                  *big.Int
		relayer              common.Address
		signature            soltypes.Signature
		lastValidBlockHeight uint64
		status               ValidationStatusType
		timestamp            time.Time
	}
)

const (
	// WaitingForWitnesses stands for a transfer which needs more valid witnesses
	WaitingForWitnesses ValidationStatusType = "new"
	// ValidationInProcess stands for a transfer in process
	ValidationInProcess = "processing"
	// InsufficientFeeRejected stands for a transfer with insufficient fee
	InsufficientFeeRejected = "insufficient"
	// ValidationSubmitted stands for a transfer with validation submitted
	ValidationSubmitted = "validated"
	// BonusPending stands for a transfer with pending bonus
	BonusPending = "bonus"
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
	// ValidationNeedSpeedUp stands for the validation of a transfer needs speed up
	ValidationNeedSpeedUp = "speedup"
)

const (
	StatusOnChainUnknown StatusOnChainType = iota
	StatusOnChainNotConfirmed
	StatusOnChainNeedSpeedUp
	StatusOnChainRejected
	StatusOnChainNonceOverwritten
	StatusOnChainSettled
)

var (
	errInsufficientWitnesses = errors.New("insufficient witnesses")
	errGasPriceTooHigh       = errors.New("gas price is too high")
	errNoncritical           = errors.New("error before submission")
	errInvalidData           = errors.New("invalid data")
)

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
	var ataOwner util.Address
	if len(transfer.AtaOwner) > 0 {
		ataOwner, err = destAddrDecoder.DecodeBytes(transfer.AtaOwner)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode ata owner")
		}
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
		ataOwner:     ataOwner,
		amount:       amount,
		fee:          fee,
		gas:          transfer.Gas,
		gasPrice:     gasPrice,
		timestamp:    transfer.Timestamp.AsTime(),
		txSender:     txSender,
		sourceTxHash: transfer.SourceTxHash,
		blockHeight:  transfer.BlockHeight,
		payload:      transfer.Payload,
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
		signature: clone,
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
		TxSender:  transfer.txSender.Bytes(),
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

// UnmarshalWitnessListProto unmarshals a witness list proto
func UnmarshalWitnessListProto(
	request *types.WitnessesList,
	destAddrDecoder util.AddressDecoder,
) (*WitnessCandidates, *Witness, error) {
	if request.Candidates == nil {
		return nil, nil, errors.New("candidates is nil")
	}
	witnessManager, err := destAddrDecoder.DecodeBytes(request.Candidates.WitnessManagerAddress)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode witness manager address")
	}
	witnessManagerAddr, ok := witnessManager.Address().(common.Address)
	if !ok {
		return nil, nil, errors.New("witness manager address is not common.Address")
	}
	witnessesToAdd := make([]util.Address, len(request.Candidates.WitnessesToAdd))
	for i, witnessBytes := range request.Candidates.WitnessesToAdd {
		witness, err := destAddrDecoder.DecodeBytes(witnessBytes)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to decode witness to add %x", witnessBytes)
		}
		witnessesToAdd[i] = witness
	}
	witnessesToRemove := make([]util.Address, len(request.Candidates.WitnessesToRemove))
	for i, witnessBytes := range request.Candidates.WitnessesToRemove {
		witness, err := destAddrDecoder.DecodeBytes(witnessBytes)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to decode witness to remove %x", witnessBytes)
		}
		witnessesToRemove[i] = witness
	}
	var witness *Witness
	if len(request.Address) > 0 && len(request.Signature) > 0 {
		witness, err = NewWitness(request.Address, request.Signature)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create witness")
		}
	}

	return &WitnessCandidates{
		witnessManager:  witnessManagerAddr,
		epoch:           request.Candidates.Epoch,
		witnessToAdd:    witnessesToAdd,
		witnessToRemove: witnessesToRemove,
	}, witness, nil
}

type (
	// WitnessCandidates defines a record of witness list update
	WitnessCandidates struct {
		id              common.Hash
		witnessManager  common.Address
		epoch           uint64
		witnessToAdd    []util.Address
		witnessToRemove []util.Address
		updateTime      time.Time
		status          ValidationStatusType
		txHash          common.Hash
		blockHeight     uint64
		gas             uint64
		nonce           uint64
		relayer         common.Address
		gasPrice        *big.Int
	}
)

func (cand *WitnessCandidates) GenID() error {
	witnessesToAdd, witnessesToRemove := cand.Witnesses()
	witnessesToAddBytes, err := packWitnesses(witnessesToAdd)
	if err != nil {
		return err
	}
	witnessesToRemoveBytes, err := packWitnesses(witnessesToRemove)
	if err != nil {
		return err
	}
	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, cand.epoch)
	data := bytes.Join([][]byte{
		cand.witnessManager.Bytes(),
		epochBytes,
		witnessesToAddBytes,
		witnessesToRemoveBytes,
	}, nil)
	cand.id = crypto.Keccak256Hash(data)
	return nil
}

func packWitnesses(witnesses []common.Address) ([]byte, error) {
	sort.Slice(witnesses, func(i, j int) bool {
		return bytes.Compare(witnesses[i][:], witnesses[j][:]) < 0
	})
	var packed []byte
	for _, witness := range witnesses {
		packed = append(packed, common.LeftPadBytes(witness.Bytes(), 32)...)
	}
	return packed, nil
}

func (cand *WitnessCandidates) Witnesses() ([]common.Address, []common.Address) {
	witnessesToAdd := []common.Address{}
	for _, witness := range cand.witnessToAdd {
		witnessesToAdd = append(witnessesToAdd, witness.Address().(common.Address))
	}
	witnessesToRemove := []common.Address{}
	for _, witness := range cand.witnessToRemove {
		witnessesToRemove = append(witnessesToRemove, witness.Address().(common.Address))
	}
	return witnessesToAdd, witnessesToRemove
}

func (cand *WitnessCandidates) ID() common.Hash {
	return cand.id
}
