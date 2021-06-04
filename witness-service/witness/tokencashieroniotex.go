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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
)

// NewTokenCashier creates a new TokenCashier
func NewTokenCashier(
	id string,
	relayerURL string,
	iotexClient iotex.ReadOnlyClient,
	cashierContractAddr address.Address,
	validatorContractAddr common.Address,
	recorder *Recorder,
	startBlockHeight uint64,
) (TokenCashier, error) {
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	eventTopic := tokenCashierABI.Events[eventName].ID.Bytes()
	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr,
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, error) {
			chainMetaResponse, err := iotexClient.API().GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
			if err != nil {
				return 0, err
			}
			tipHeight := chainMetaResponse.ChainMeta.Height
			if startHeight >= tipHeight {
				return 0, errors.Errorf("query height %d is larger than chain tip height %d", startHeight, tipHeight)
			}
			if count == 0 {
				count = 1
			}
			endHeight := startHeight + uint64(count) - 1
			if endHeight > tipHeight {
				endHeight = tipHeight
			}
			return endHeight, nil
		},
		func(startHeight uint64, endHeight uint64) ([]*Transfer, error) {
			response, err := iotexClient.API().GetLogs(context.Background(), &iotexapi.GetLogsRequest{
				Filter: &iotexapi.LogsFilter{
					Address: []string{cashierContractAddr.String()},
					Topics: []*iotexapi.Topics{
						{
							Topic: [][]byte{
								eventTopic,
							},
						},
					},
				},
				Lookup: &iotexapi.GetLogsRequest_ByRange{
					ByRange: &iotexapi.GetLogsByRange{
						FromBlock: startHeight,
						// TODO: this is a bug, which should be fixed in iotex-core
						ToBlock: endHeight,
					},
				},
			})
			if err != nil {
				return nil, err
			}
			transfers := []*Transfer{}
			if len(response.Logs) > 0 {
				log.Printf("\t%d transfers fetched", len(response.Logs))
				for _, log := range response.Logs {
					if bytes.Compare(eventTopic, log.Topics[0]) != 0 {
						return nil, errors.Errorf("Wrong event topic %s, %s expected", log.Topics[0], eventTopic)
					}
					if len(log.Data) != 128 {
						return nil, errors.Errorf("Invalid data length %d, 128 expected", len(log.Data))
					}
					cashier, err := address.FromString(log.ContractAddress)
					if err != nil {
						return nil, err
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
			}
			return transfers, nil
		},
	), nil
}
