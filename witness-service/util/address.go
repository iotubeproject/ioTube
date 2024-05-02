package util

import (
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
)

type Address interface {
	Bytes() []byte
	String() string
	Address() interface{}
}

type ethAddress struct {
	address common.Address
}

func (a *ethAddress) Bytes() []byte {
	return a.address.Bytes()
}

func (a *ethAddress) String() string {
	return a.address.Hex()
}

func (a *ethAddress) Address() interface{} {
	return a.address
}

func ETHAddressToAddress(addr common.Address) Address {
	return &ethAddress{
		address: addr,
	}
}

type solAddress struct {
	address solcommon.PublicKey
}

func (s *solAddress) Bytes() []byte {
	return s.address.Bytes()
}

func (s *solAddress) String() string {
	return s.address.String()
}

func (s *solAddress) Address() interface{} {
	return s.address
}

func SOLAddressToAddress(addr solcommon.PublicKey) Address {
	return &solAddress{
		address: addr,
	}
}

type AddressDecoder interface {
	DecodeBytes([]byte) (Address, error)
	DecodeString(string) (Address, error)
}

type ETHAddressDecoder struct{}

func NewETHAddressDecoder() AddressDecoder {
	return &ETHAddressDecoder{}
}

func (d *ETHAddressDecoder) DecodeBytes(b []byte) (Address, error) {
	return &ethAddress{
		address: common.BytesToAddress(b),
	}, nil
}

func (d *ETHAddressDecoder) DecodeString(s string) (Address, error) {
	return &ethAddress{
		address: common.HexToAddress(s),
	}, nil
}

type SOLAddressDecoder struct{}

func NewSOLAddressDecoder() AddressDecoder {
	return &SOLAddressDecoder{}
}

func (s *SOLAddressDecoder) DecodeBytes(b []byte) (Address, error) {
	return &solAddress{
		address: solcommon.PublicKeyFromBytes(b),
	}, nil
}

func (s *SOLAddressDecoder) DecodeString(str string) (Address, error) {
	return &solAddress{
		address: solcommon.PublicKeyFromString(str),
	}, nil
}
