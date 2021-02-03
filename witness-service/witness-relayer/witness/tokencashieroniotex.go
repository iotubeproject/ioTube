// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-ethereum/accounts/abi"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
)

type (
	// tokenCashierOnIoTeX maintains the list of witnesses and tokens
	tokenCashierOnIoTeX struct {
		cashierContractAddr address.Address
		iotexClient         iotex.ReadOnlyClient
		tokenCashierABI     abi.ABI
	}
)

// NewTokenCashier creates a new TokenCashier
func NewTokenCashier(cashierContractAddr address.Address, iotexClient iotex.ReadOnlyClient) (TokenCashier, error) {
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	return &tokenCashierOnIoTeX{
		cashierContractAddr: cashierContractAddr,
		iotexClient:         iotexClient,
		tokenCashierABI:     tokenCashierABI,
	}, nil
}

// PullTransfers pulls transfers by query token cashier receipts
func (tc *tokenCashierOnIoTeX) PullTransfers(offset uint64, count uint16) (uint64, []*Transfer, error) {
	topicToFilter := tc.tokenCashierABI.Events[eventName].Id().Bytes()
	chainMetaResponse, err := tc.iotexClient.API().GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	if err != nil {
		return 0, nil, err
	}
	tipHeight := chainMetaResponse.ChainMeta.Height
	if offset >= tipHeight {
		return 0, nil, errors.Errorf("query height %d is larger than chain tip height %d", offset, tipHeight)
	}
	if count == 0 {
		count = 1
	}
	endHeight := offset + uint64(count) - 1
	if endHeight > tipHeight {
		endHeight = tipHeight
	}
	log.Printf("fetching events from block %d\n", offset)
	response, err := tc.iotexClient.API().GetLogs(context.Background(), &iotexapi.GetLogsRequest{
		Filter: &iotexapi.LogsFilter{
			Address: []string{tc.cashierContractAddr.String()},
			Topics: []*iotexapi.Topics{
				{
					Topic: [][]byte{
						topicToFilter,
					},
				},
			},
		},
		Lookup: &iotexapi.GetLogsRequest_ByRange{
			ByRange: &iotexapi.GetLogsByRange{
				FromBlock: offset,
				// TODO: this is a bug, which should be fixed in iotex-core
				Count: endHeight,
			},
		},
	})
	if err != nil {
		return 0, nil, err
	}
	log.Printf("\t%d transfers fetched", len(response.Logs))
	transfers := []*Transfer{}
	for _, log := range response.Logs {
		if bytes.Compare(topicToFilter, log.Topics[0]) != 0 {
			return 0, nil, errors.Errorf("Wrong event topic %s, %s expected", log.Topics[0], topicToFilter)
		}
		if len(log.Data) != 128 {
			return 0, nil, errors.Errorf("Invalid data length %d, 128 expected", len(log.Data))
		}
		cashier, err := address.FromString(log.ContractAddress)
		if err != nil {
			return 0, nil, err
		}
		transfers = append(transfers, &Transfer{
			cashier:     common.BytesToAddress(cashier.Bytes()),
			token:       common.BytesToAddress(log.Topics[1]),
			index:       new(big.Int).SetBytes(log.Topics[2]).Uint64(),
			sender:      common.BytesToAddress(log.Data[:32]),
			recipient:   common.BytesToAddress(log.Data[32:64]),
			amount:      new(big.Int).SetBytes(log.Data[64:96]),
			blockHeight: log.BlkHeight,
			txHash:      common.BytesToHash(log.ActHash),
		})
	}
	return endHeight, transfers, nil
}
