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

type (
	// tokenCashierOnEthereum maintains the list of witnesses and tokens
	tokenCashierOnEthereum struct {
		cashierContractAddr common.Address
		ethereumClient      *ethclient.Client
		eventTopic          common.Hash
		confirmBlockNumber  uint64
	}
)

// NewTokenCashierOnEthereum creates a new TokenCashier on ethereum
func NewTokenCashierOnEthereum(cashierContractAddr common.Address, ethereumClient *ethclient.Client, confirmBlockNumber uint8) (TokenCashier, error) {
	tokenCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	return &tokenCashierOnEthereum{
		cashierContractAddr: cashierContractAddr,
		ethereumClient:      ethereumClient,
		eventTopic:          tokenCashierABI.Events[eventName].Id(),
		confirmBlockNumber:  uint64(confirmBlockNumber),
	}, nil
}

// PullTransfers pulls transfers by query token cashier receipts
func (tc *tokenCashierOnEthereum) PullTransfers(offset uint64, count uint16) (uint64, []*Transfer, error) {
	tipHeader, err := tc.ethereumClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, nil, err
	}
	tipHeight := tipHeader.Number.Uint64() - tc.confirmBlockNumber
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
	logs, err := tc.ethereumClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(offset),
		ToBlock:   new(big.Int).SetUint64(endHeight),
		Addresses: []common.Address{tc.cashierContractAddr},
		Topics: [][]common.Hash{
			{
				tc.eventTopic,
			},
		},
	})
	if err != nil {
		return 0, nil, err
	}
	log.Printf("\t%d transfers fetched", len(logs))
	transfers := []*Transfer{}
	for _, log := range logs {
		if !bytes.Equal(tc.eventTopic[:], log.Topics[0][:]) {
			return 0, nil, errors.Errorf("Wrong event topic %x, %x expected", log.Topics[0], tc.eventTopic)
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
	return endHeight, transfers, nil
}
