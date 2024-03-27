// Copyright (c) 2024 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
)

func ParseAddress(addr string) (common.Address, error) {
	if strings.HasPrefix(addr, "io") {
		ioAddr, err := address.FromString(addr)
		if err != nil {
			log.Fatalf("failed to parse iotex address %s, %v\n", addr, err)
		}
		return common.BytesToAddress(ioAddr.Bytes()), nil
	}
	return common.HexToAddress(addr), nil
}
