package util

import (
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
