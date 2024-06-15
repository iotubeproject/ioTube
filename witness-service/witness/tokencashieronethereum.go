// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/pkg/errors"
)

type iterator struct {
	version               Version
	client                *ethclient.Client
	cashierContractAddr   common.Address
	tokenSafeContractAddr common.Address
}

func newIterator(version Version, cashierContractAddr, tokenSafeContractAddr common.Address, client *ethclient.Client) (*iterator, error) {
	iter := &iterator{
		version:               version,
		cashierContractAddr:   cashierContractAddr,
		tokenSafeContractAddr: tokenSafeContractAddr,
		client:                client,
	}
	var err error
	switch version {
	case NoPayload, Payload:
	default:
		return nil, errors.Errorf("invalid version %s", version)
	}
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (iter *iterator) extract(
	tokenAddress, senderAddress, recipient common.Address,
	index uint64,
	amount, fee *big.Int,
	payload []byte,
	raw types.Log,
) (*Transfer, error) {
	receipt, err := iter.client.TransactionReceipt(context.Background(), raw.TxHash)
	if err != nil {
		return nil, err
	}
	var realAmount *big.Int
	for _, l := range receipt.Logs {
		if l.Address == tokenAddress && l.Topics[0] == _TransferEventTopic && (l.Topics[1] == senderAddress.Hash() || l.Topics[1] == raw.Address.Hash()) {
			if l.Topics[2] == iter.cashierContractAddr.Hash() || l.Topics[2] != _ZeroHash && l.Topics[2] == iter.tokenSafeContractAddr.Hash() {
				if realAmount != nil {
					return nil, errors.Errorf("two transfers in one transaction %x", raw.TxHash)
				}
				realAmount = new(big.Int).SetBytes(l.Data)
			}
		}
	}
	if realAmount == nil {
		return nil, errors.Errorf("failed to get the amount from transfer event for %x", raw.TxHash)
	}
	switch realAmount.Cmp(amount) {
	case 1:
		return nil, errors.Errorf("Invalid amount: %d < %d", amount, realAmount)
	case -1:
		log.Printf("\tAmount %d is reduced %d after tax\n", amount, realAmount)
	case 0:
		log.Printf("\tAmount %d is the same as real amount %d\n", amount, realAmount)
	}
	tx, err := iter.client.TransactionInBlock(context.Background(), raw.BlockHash, raw.TxIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch transaction")
	}
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract sender")
	}
	tsf := &Transfer{
		cashier:     raw.Address,
		token:       tokenAddress,
		index:       index,
		sender:      senderAddress,
		recipient:   recipient,
		amount:      amount,
		fee:         fee,
		blockHeight: raw.BlockNumber,
		txHash:      raw.TxHash,
		payload:     payload,
	}
	if from != senderAddress {
		tsf.txSender = from
	}

	return tsf, nil
}

func (iter *iterator) Transfers(start, end uint64) ([]*Transfer, error) {
	transfers := []*Transfer{}
	switch iter.version {
	case NoPayload:
		filter, err := contract.NewTokenCashierFilterer(iter.cashierContractAddr, iter.client)
		if err != nil {
			return nil, err
		}
		iterator, err := filter.FilterReceipt(
			&bind.FilterOpts{
				Start: start,
				End:   &end,
			},
			nil,
			nil,
		)
		if err != nil {
			return nil, err
		}
		for iterator.Next() {
			tsf, err := iter.extract(
				iterator.Event.Token,
				iterator.Event.Sender,
				iterator.Event.Recipient,
				iterator.Event.Id.Uint64(),
				iterator.Event.Amount,
				iterator.Event.Fee,
				nil,
				iterator.Event.Raw,
			)
			if err != nil {
				return nil, err
			}
			transfers = append(transfers, tsf)
		}
	case Payload:
		filter, err := contract.NewTokenCashierWithPayloadFilterer(iter.cashierContractAddr, iter.client)
		if err != nil {
			return nil, err
		}
		iterator, err := filter.FilterReceipt(
			&bind.FilterOpts{
				Start: start,
				End:   &end,
			},
			nil,
			nil,
		)
		if err != nil {
			return nil, err
		}
		for iterator.Next() {
			tsf, err := iter.extract(
				iterator.Event.Token,
				iterator.Event.Sender,
				iterator.Event.Recipient,
				iterator.Event.Id.Uint64(),
				iterator.Event.Amount,
				iterator.Event.Fee,
				iterator.Event.Payload,
				iterator.Event.Raw,
			)
			if err != nil {
				return nil, err
			}
			transfers = append(transfers, tsf)

		}
	default:
		return nil, errors.New("invalid version")
	}

	return transfers, nil
}

// NewTokenCashierOnEthereum creates a new TokenCashier on ethereum
func NewTokenCashierOnEthereum(
	id string,
	version Version,
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
	iter, err := newIterator(version, cashierContractAddr, tokenSafeContractAddr, ethereumClient)
	if err != nil {
		return nil, err
	}
	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr,
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
		iter.Transfers,
		func(token common.Address, amountToTransfer *big.Int) bool {
			if reverseRecorder == nil {
				return true
			}
			coToken, ok := recorder.tokenPairs[token]
			if !ok {
				return false
			}
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
	), nil
}
