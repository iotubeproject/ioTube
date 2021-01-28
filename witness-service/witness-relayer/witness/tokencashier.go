// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-ethereum/accounts/abi"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

// TokenCashier maintains the list of witnesses and tokens
type TokenCashier struct {
	mu                      sync.RWMutex
	lastUpdateTime          time.Time
	gasPriceLimitOnEthereum *big.Int

	iotexClient iotex.ReadOnlyClient
}

// NewTokenCashier creates a new TokenCashier
func NewTokenCashier(iotexClient iotex.ReadOnlyClient) *TokenCashier {
	return &TokenCashier{iotexClient: iotexClient}
}

type receipt struct {
	token     common.Address
	id        *big.Int
	sender    common.Address
	recipient common.Address
	amount    *big.Int
	fee       *big.Int
}

// FetchTransfers fetches transfers by query token cashier receipts
func (tc *TokenCashier) FetchTransfers(start uint64, end uint64) ([]*Transfer, error) {
	response, err := tc.iotexClient.API().GetLogs(context.Background(), &iotexapi.GetLogsRequest{})
	if err != nil {
		return nil, err
	}
	transfers := []*Transfer{}
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierV2ABI))
	if err != nil {
		return nil, err
	}

	for _, log := range response.Logs {
		cashier, err := address.FromString(log.ContractAddress)
		if err != nil {
			return nil, err
		}
		var r receipt
		// TODO: verify topics[0]
		if err := tokenCashierABI.Unpack(&r, "Receipt", log.Data); err != nil {
			return nil, err
		}
		transfers = append(transfers, &Transfer{
			cashier:     common.BytesToAddress(cashier.Bytes()),
			token:       common.BytesToAddress(log.Topics[1]),
			index:       new(big.Int).SetBytes(log.Topics[2]).Uint64(),
			sender:      r.sender,
			recipient:   r.recipient,
			amount:      r.amount,
			blockHeight: log.BlkHeight,
			txHash:      common.BytesToHash(log.ActHash),
		})
	}
	return transfers, nil
}
