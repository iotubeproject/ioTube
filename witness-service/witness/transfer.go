package witness

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
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

func (t *Transfer) Cashier() common.Address {
	return t.cashier
}

func (t *Transfer) Token() common.Address {
	return t.token
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

// bTCTransfer defines a BTC record
type bTCTransfer struct {
	version     uint32
	cashier     *btcec.PublicKey // musig combined pubkey
	blockHeight uint64
	txHash      chainhash.Hash
	vout        uint64
	sender      *btcec.PublicKey
	recipient   common.Address
	amount      btcutil.Amount
	fee         btcutil.Amount
	pkScript    []byte
	status      TransferStatus
	metadata    []byte

	id      common.Hash
	coToken common.Address
}

func (t *bTCTransfer) Cashier() common.Address {
	return crypto.PubkeyToAddress(*t.sender.ToECDSA())
}

func (t *bTCTransfer) Token() common.Address {
	return t.coToken
}

func (t *bTCTransfer) Index() *big.Int {
	return btcIndexFromTxHashAndVout(t.txHash, t.vout)
}

func (t *bTCTransfer) ID() ([]byte, error) {
	for i := range t.id[:] {
		if t.id[i] != 0 {
			return t.id[:], nil
		}
	}
	return nil, errors.New("t id is empty")
}

func (t *bTCTransfer) SetID(id common.Hash) {
	t.id = id
}

func (t *bTCTransfer) BlockHeight() uint64 {
	return t.blockHeight
}

func (t *bTCTransfer) DataToSign() []byte {
	return bytes.Join([][]byte{
		t.Cashier().Bytes(),
		t.coToken.Bytes(),
		math.U256Bytes(btcIndexFromTxHashAndVout(t.txHash, t.vout)),
		crypto.PubkeyToAddress(*t.sender.ToECDSA()).Bytes(),
		t.recipient.Bytes(),
		math.U256Bytes(t.Amount()),
	}, []byte{})
}

func btcIndexFromTxHashAndVout(txHash chainhash.Hash, vout uint64) *big.Int {
	voutByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(voutByte, vout)
	index := crypto.Keccak256Hash(txHash[:], voutByte)
	return new(big.Int).SetBytes(index[:])
}

func (t *bTCTransfer) ToTypesTransfer() *types.Transfer {
	return &types.Transfer{
		Cashier:   t.Cashier().Bytes(),
		Token:     t.coToken.Bytes(),
		Index:     0,
		BtcIndex:  btcIndexFromTxHashAndVout(t.txHash, t.vout).String(),
		Sender:    crypto.PubkeyToAddress(*t.sender.ToECDSA()).Bytes(),
		Recipient: t.recipient.Bytes(),
		Amount:    t.Amount().String(),
		Fee:       t.fee.String(),
		TxSender:  crypto.PubkeyToAddress(*t.sender.ToECDSA()).Bytes(),
	}
}

func (t *bTCTransfer) Recipient() util.Address {
	addr, err := util.NewETHAddressDecoder().DecodeString(t.recipient.String())
	if err != nil {
		log.Panicf("failed to decode recipient: %v", err)
	}
	return addr
}

func (t *bTCTransfer) Amount() *big.Int {
	return big.NewInt(int64(t.amount))
}

func (t *bTCTransfer) Status() TransferStatus {
	return t.status
}
