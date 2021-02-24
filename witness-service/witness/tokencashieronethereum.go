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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/go-ethereum/accounts/abi"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/pkg/errors"
)

// NewTokenCashierOnEthereum creates a new TokenCashier on ethereum
func NewTokenCashierOnEthereum(
	id string,
	relayerURL string,
	ethereumClient *ethclient.Client,
	cashierContractAddr common.Address,
	validatorContractAddr common.Address,
	recorder *Recorder,
	startBlockHeight uint64,
	confirmBlockNumber uint8,
) (TokenCashier, error) {
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	eventTopic := tokenCashierABI.Events[eventName].Id()
	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr,
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, error) {
			tipHeader, err := ethereumClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return 0, err
			}
			tipHeight := tipHeader.Number.Uint64() - uint64(confirmBlockNumber)
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
			logs, err := ethereumClient.FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(startHeight),
				ToBlock:   new(big.Int).SetUint64(endHeight),
				Addresses: []common.Address{cashierContractAddr},
				Topics: [][]common.Hash{
					{
						eventTopic,
					},
				},
			})
			if err != nil {
				return nil, err
			}
			transfers := []*Transfer{}
			if len(logs) > 0 {
				log.Printf("\t%d transfers fetched", len(logs))
				for _, log := range logs {
					if !bytes.Equal(eventTopic[:], log.Topics[0][:]) {
						return nil, errors.Errorf("Wrong event topic %x, %x expected", log.Topics[0], eventTopic)
					}
					transfers = append(transfers, &Transfer{
						cashier:     log.Address,
						token:       common.BytesToAddress(log.Topics[1][:]),
						index:       new(big.Int).SetBytes(log.Topics[2][:]).Uint64(),
						sender:      common.BytesToAddress(log.Data[:32]),
						recipient:   common.BytesToAddress(log.Data[32:64]),
						amount:      new(big.Int).SetBytes(log.Data[64:96]),
						blockHeight: log.BlockNumber,
						txHash:      log.TxHash,
					})
				}
			}
			return transfers, nil
		},
	), nil
}
