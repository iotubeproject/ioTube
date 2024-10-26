// Copyright (c) 2024 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"log"
	"strings"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
)

func ParseEthAddress(addr string) (common.Address, error) {
	if strings.HasPrefix(addr, "io") {
		ioAddr, err := address.FromString(addr)
		if err != nil {
			log.Fatalf("failed to parse iotex address %s, %v\n", addr, err)
		}
		return common.BytesToAddress(ioAddr.Bytes()), nil
	}
	return common.HexToAddress(addr), nil
}

func ParseAddress(addr string) (Address, error) {
	if len(addr) == 44 {
		return SOLAddressToAddress(solcommon.PublicKeyFromString(addr)), nil
	}
	ethaddr, err := ParseEthAddress(addr)
	if err != nil {
		return nil, err
	}
	return ETHAddressToAddress(ethaddr), nil
}

func ParseAddressBytes(b []byte) (Address, error) {
	if len(b) == 32 {
		return SOLAddressToAddress(solcommon.PublicKeyFromBytes(b)), nil
	}
	return ETHAddressToAddress(common.BytesToAddress(b)), nil
}

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

func (d *ETHAddressDecoder) DecodeString(addr string) (Address, error) {
	if strings.HasPrefix(addr, "io") {
		ioAddr, err := address.FromString(addr)
		if err != nil {
			log.Fatalf("failed to parse iotex address %s, %v\n", addr, err)
		}
		return &ethAddress{
			common.BytesToAddress(ioAddr.Bytes()),
		}, nil
	}
	return &ethAddress{
		address: common.HexToAddress(addr),
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

func NewAddressDecoder(chain string) AddressDecoder {
	switch chain {
	default:
		return NewETHAddressDecoder()
	case "sol", "solana":
		return NewSOLAddressDecoder()
	}
}
