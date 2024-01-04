package util

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

type btcAddress struct {
	address btcutil.Address
}

func (a *btcAddress) Bytes() []byte {
	return []byte(a.String())
}

func (a *btcAddress) String() string {
	return a.address.EncodeAddress()
}

func (a *btcAddress) Address() interface{} {
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

type BTCAddressDecoder struct {
	chain *chaincfg.Params
}

func NewBTCAddressDecoder(chain *chaincfg.Params) AddressDecoder {
	return &BTCAddressDecoder{chain: chain}
}

func (d *BTCAddressDecoder) DecodeBytes(b []byte) (Address, error) {
	addr, err := btcutil.DecodeAddress(string(b), d.chain)
	if err != nil {
		return nil, err
	}
	return &btcAddress{address: addr}, nil
}

func (d *BTCAddressDecoder) DecodeString(s string) (Address, error) {
	addr, err := btcutil.DecodeAddress(s, d.chain)
	if err != nil {
		return nil, err
	}
	return &btcAddress{address: addr}, nil
}
