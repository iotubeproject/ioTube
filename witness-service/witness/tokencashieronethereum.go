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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

// NewTokenCashierOnEthereum creates a new TokenCashier on ethereum
func NewTokenCashierOnEthereum(
	id string,
	relayerURL string,
	ethereumClient *ethclient.Client,
	cashierContractAddr common.Address,
	tokenSafeContractAddr common.Address,
	validatorContractAddr common.Address,
	recorder *Recorder,
	startBlockHeight uint64,
	confirmBlockNumber uint8,
	reverseRecorder *Recorder,
	reverseCashierContractAddr common.Address,
) (TokenCashier, error) {
	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr.Bytes(),
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, uint64, error) {
			tipHeader, err := ethereumClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return 0, 0, errors.Wrap(err, "failed to query tip block header")
			}
			tipHeight := tipHeader.Number.Uint64()
			if count == 0 {
				count = 1
			}
			if tipHeight <= uint64(confirmBlockNumber) {
				return 0, 0, errors.Errorf("tip height %d is smaller than confirm block number %d", tipHeight, confirmBlockNumber)
			}
			endHeight := startHeight + uint64(count) - 1
			if tipHeight < endHeight {
				endHeight = tipHeight
			}
			return tipHeight - uint64(confirmBlockNumber), endHeight, nil
		},
		func(startHeight uint64, endHeight uint64) ([]AbstractTransfer, error) {
			logs, err := ethereumClient.FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(startHeight),
				ToBlock:   new(big.Int).SetUint64(endHeight),
				Addresses: []common.Address{cashierContractAddr},
				Topics: [][]common.Hash{
					{
						_ReceiptEventTopic,
					},
				},
			})
			if err != nil {
				return nil, err
			}
			transfers := []AbstractTransfer{}
			if len(logs) > 0 {
				log.Printf("\t%d transfers fetched from %d to %d\n", len(logs), startHeight, endHeight)
				for _, transferLog := range logs {
					if !bytes.Equal(_ReceiptEventTopic[:], transferLog.Topics[0][:]) {
						return nil, errors.Errorf("Wrong event topic %x, %x expected", transferLog.Topics[0], _ReceiptEventTopic)
					}
					tokenAddress := common.BytesToAddress(transferLog.Topics[1][:])
					senderAddress := common.BytesToAddress(transferLog.Data[:32])
					amount := new(big.Int).SetBytes(transferLog.Data[64:96])
					receipt, err := ethereumClient.TransactionReceipt(context.Background(), transferLog.TxHash)
					if err != nil {
						return nil, err
					}
					var realAmount *big.Int
					for _, l := range receipt.Logs {
						if l.Address == tokenAddress && l.Topics[0] == _TransferEventTopic && (l.Topics[1] == common.BytesToHash(senderAddress.Bytes()) || l.Topics[1] == common.BytesToHash(transferLog.Address.Bytes())) {
							if l.Topics[2] == common.BytesToHash(cashierContractAddr.Bytes()) || l.Topics[2] != _ZeroHash && l.Topics[2] == common.BytesToHash(tokenSafeContractAddr.Bytes()) {
								if realAmount != nil {
									return nil, errors.Errorf("two transfers in one transaction %x", transferLog.TxHash)
								}
								realAmount = new(big.Int).SetBytes(l.Data)
							}
						}
					}
					if realAmount == nil {
						return nil, errors.Errorf("failed to get the amount from transfer event for %x", transferLog.TxHash)
					}
					switch realAmount.Cmp(amount) {
					case 1:
						return nil, errors.Errorf("Invalid amount: %d < %d", amount, realAmount)
					case -1:
						log.Printf("\tAmount %d is reduced %d after tax\n", amount, realAmount)
					case 0:
						log.Printf("\tAmount %d is the same as real amount %d\n", amount, realAmount)
					}
					tx, err := ethereumClient.TransactionInBlock(context.Background(), transferLog.BlockHash, transferLog.TxIndex)
					if err != nil {
						return nil, errors.Wrap(err, "failed to fetch transaction")
					}
					from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
					if err != nil {
						return nil, errors.Wrap(err, "failed to extract sender")
					}
					recipient, err := util.NewETHAddressDecoder().DecodeBytes(transferLog.Data[32:64])
					if err != nil {
						return nil, errors.Wrap(err, "failed to decode recipient")
					}
					tsf := &Transfer{
						cashier:     transferLog.Address,
						token:       tokenAddress,
						index:       new(big.Int).SetBytes(transferLog.Topics[2][:]).Uint64(),
						sender:      senderAddress,
						recipient:   recipient,
						amount:      amount,
						fee:         new(big.Int).SetBytes(transferLog.Data[96:128]),
						blockHeight: transferLog.BlockNumber,
						txHash:      transferLog.TxHash,
					}
					if from != senderAddress {
						tsf.txSender = from
					}
					transfers = append(transfers, tsf)
				}
			}
			return transfers, nil
		},
		func(_token util.Address, amountToTransfer *big.Int) bool {
			if reverseRecorder == nil {
				return true
			}
			token := _token.Address().(common.Address)
			_coToken, ok := recorder.tokenPairs[token]
			if !ok {
				return false
			}
			coToken := _coToken.Address().(common.Address)
			if _, ok := reverseRecorder.tokenPairs[coToken]; !ok {
				return true
			}
			inAmount, err := reverseRecorder.AmountOfTransferred(reverseCashierContractAddr, coToken)
			if err != nil {
				return false
			}
			outAmount, err := recorder.AmountOfTransferred(cashierContractAddr, token)
			if err != nil {
				return false
			}
			return inAmount.Cmp(big.NewInt(0).Add(outAmount, amountToTransfer)) >= 0
		},
		func(ctx context.Context) error {
			if reverseRecorder != nil {
				return reverseRecorder.Start(ctx)
			}
			return nil
		},
		func(ctx context.Context) error {
			if reverseRecorder != nil {
				return reverseRecorder.Stop(ctx)
			}
			return nil
		},
		false,
	), nil
}
