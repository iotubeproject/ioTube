// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/pkg/errors"
)

type iotexIterator struct {
	version             Version
	client              iotex.ReadOnlyClient
	cashierContractAddr address.Address
}

func (ii *iotexIterator) filterLogs(topic []byte, startHeight, endHeight uint64) (*iotexapi.GetLogsResponse, error) {
	return ii.client.API().GetLogs(context.Background(), &iotexapi.GetLogsRequest{
		Filter: &iotexapi.LogsFilter{
			Address: []string{ii.cashierContractAddr.String()},
			Topics: []*iotexapi.Topics{
				{
					Topic: [][]byte{topic},
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
}

func (ii *iotexIterator) extractTransfer(
	transferLog *iotextypes.Log,
	topic common.Hash,
) (*Transfer, error) {
	senderAddr := common.BytesToAddress(transferLog.Data[:32])
	amount := new(big.Int).SetBytes(transferLog.Data[64:96])

	receipt, err := ii.client.API().GetReceiptByAction(context.Background(), &iotexapi.GetReceiptByActionRequest{
		ActionHash: hex.EncodeToString(transferLog.ActHash),
	})
	if err != nil {
		return nil, err
	}
	tokenAddr := common.BytesToAddress(transferLog.Topics[1])
	tokenIoAddr, err := address.FromBytes(tokenAddr.Bytes())
	if err != nil {
		return nil, err
	}
	cashierAddr, err := address.FromString(transferLog.ContractAddress)
	if err != nil {
		return nil, err
	}
	cashier := common.BytesToAddress(cashierAddr.Bytes())
	var realAmount *big.Int
	for _, l := range receipt.ReceiptInfo.Receipt.Logs {
		if tokenIoAddr.String() == l.ContractAddress && common.BytesToHash(l.Topics[0]) == topic && (common.BytesToAddress(l.Topics[1]) == senderAddr || cashier == common.BytesToAddress(l.Topics[1])) {
			if realAmount != nil && common.BytesToHash(l.Topics[2]) != _ZeroHash {
				return nil, errors.Errorf("two transfers in one transaction %x", transferLog.ActHash)
			}
			realAmount = new(big.Int).SetBytes(l.Data)
		}
	}
	if realAmount == nil {
		return nil, errors.Errorf("failed to get the amount from transfer event for %x", transferLog.ActHash)
	}
	switch realAmount.Cmp(amount) {
	case 1:
		return nil, errors.Errorf("Invalid amount: %d < %d", amount, realAmount)
	case -1:
		log.Printf("\tAmount %d is reduced %d after tax\n", amount, realAmount)
	case 0:
		log.Printf("\tAmount %d is the same as real amount %d\n", amount, realAmount)
	}

	return &Transfer{
		cashier:     cashier,
		token:       tokenAddr,
		index:       new(big.Int).SetBytes(transferLog.Topics[2]).Uint64(),
		sender:      senderAddr,
		recipient:   common.BytesToAddress(transferLog.Data[32:64]),
		amount:      amount,
		fee:         new(big.Int).SetBytes(transferLog.Data[96:128]),
		blockHeight: transferLog.BlkHeight,
		txHash:      common.BytesToHash(transferLog.ActHash),
		payload:     transferLog.Data[128:],
	}, nil
}

func (ii *iotexIterator) Transfers(startHeight uint64, endHeight uint64) ([]*Transfer, error) {
	transfers := []*Transfer{}
	switch ii.version {
	case NoPayload:
		response, err := ii.filterLogs(_ReceiptEventTopic[:], startHeight, endHeight)
		if err != nil {
			return nil, err
		}
		for _, transferLog := range response.Logs {
			if len(transferLog.Data) != 128 {
				return nil, errors.Errorf("Invalid data length %d, 128 expected", len(transferLog.Data))
			}
			tsf, err := ii.extractTransfer(transferLog, _ReceiptEventTopic)
			if err != nil {
				return nil, err
			}
			transfers = append(transfers, tsf)
		}
	case Payload:
		response, err := ii.filterLogs(_ReceiptEventTopicWithPayload[:], startHeight, endHeight)
		if err != nil {
			return nil, err
		}
		for _, transferLog := range response.Logs {
			if len(transferLog.Data) >= 128 {
				return nil, errors.Errorf("Invalid data length %d < 128", len(transferLog.Data))
			}
			tsf, err := ii.extractTransfer(transferLog, _ReceiptEventTopicWithPayload)
			if err != nil {
				return nil, err
			}
			transfers = append(transfers, tsf)
		}
	}
	return transfers, nil
}

// NewTokenCashier creates a new TokenCashier
func NewTokenCashier(
	id string,
	version Version,
	relayerURL string,
	iotexClient iotex.ReadOnlyClient,
	cashierContractAddr address.Address,
	validatorContractAddr common.Address,
	recorder *Recorder,
	startBlockHeight uint64,
) (TokenCashier, error) {
	iter := &iotexIterator{
		version:             version,
		client:              iotexClient,
		cashierContractAddr: cashierContractAddr,
	}
	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr,
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, uint64, error) {
			chainMetaResponse, err := iotexClient.API().GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
			if err != nil {
				return 0, 0, err
			}
			tipHeight := chainMetaResponse.ChainMeta.Height
			if startHeight > tipHeight {
				return 0, 0, errors.Errorf("query height %d is larger than chain tip height %d", startHeight, tipHeight)
			}
			if count == 0 {
				count = 1
			}
			endHeight := startHeight + uint64(count) - 1
			if endHeight > tipHeight {
				endHeight = tipHeight
			}
			return endHeight, endHeight, nil
		},
		iter.Transfers,
		func(common.Address, *big.Int) bool {
			return true
		},
		func(context.Context) error {
			return nil
		},
		func(context.Context) error {
			return nil
		},
	), nil
}
