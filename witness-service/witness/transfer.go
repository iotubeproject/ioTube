package witness

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// Transfer defines a record
type Transfer struct {
	cashier     common.Address
	token       common.Address
	coToken     common.Address
	index       uint64
	sender      common.Address
	recipient   util.Address
	amount      *big.Int
	fee         *big.Int
	id          common.Hash
	status      TransferStatus
	blockHeight uint64
	txHash      common.Hash
	timestamp   time.Time
	gas         uint64
	gasPrice    *big.Int
	txSender    common.Address
}

func (t *Transfer) Cashier() util.Address {
	return util.ETHAddressToAddress(t.cashier)
}

func (t *Transfer) Token() util.Address {
	return util.ETHAddressToAddress(t.token)
}

func (t *Transfer) Index() *big.Int {
	return new(big.Int).SetUint64(t.index)
}

func (t *Transfer) ID() ([]byte, error) {
	for i := range t.id[:] {
		if t.id[i] != 0 {
			return t.id[:], nil
		}
	}
	return nil, errors.New("t id is empty")
}

func (t *Transfer) SetID(id common.Hash) {
	t.id = id
}

func (t *Transfer) BlockHeight() uint64 {
	return t.blockHeight
}

func (t *Transfer) DataToSign() []byte {
	return bytes.Join([][]byte{
		t.cashier.Bytes(),
		t.coToken.Bytes(),
		math.U256Bytes(new(big.Int).SetUint64(t.index)),
		t.sender.Bytes(),
		t.recipient.Bytes(),
		math.U256Bytes(t.amount),
	}, []byte{})
}

func (t *Transfer) ToTypesTransfer() *types.Transfer {
	gasPrice := "0"
	if t.gasPrice != nil {
		gasPrice = t.gasPrice.String()
	}
	return &types.Transfer{
		Cashier:   t.cashier.Bytes(),
		Token:     t.coToken.Bytes(),
		Index:     int64(t.index),
		Sender:    t.sender.Bytes(),
		Recipient: t.recipient.Bytes(),
		Amount:    t.amount.String(),
		Timestamp: timestamppb.New(t.timestamp),
		Fee:       t.fee.String(),
		TxSender:  t.txSender.Bytes(),
		Gas:       t.gas,
		GasPrice:  gasPrice,
	}
}

func (t *Transfer) Recipient() util.Address {
	return t.recipient
}

func (t *Transfer) Amount() *big.Int {
	return t.amount
}

func (t *Transfer) Status() TransferStatus {
	return t.status
}
