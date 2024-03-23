// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// NewTokenCashier creates a new TokenCashier
func NewTokenCashier(
	id string,
	relayerURL string,
	iotexClient iotex.ReadOnlyClient,
	cashierContractAddr address.Address,
	validatorContractAddr []byte,
	recorder *Recorder,
	startBlockHeight uint64,
	addrDecoder util.AddressDecoder,
) (TokenCashier, error) {
	var (
		getTransferInfo func([]byte) (common.Address, util.Address, *big.Int, *big.Int, error)
		logTopic        []byte
	)
	switch addrDecoder.(type) {
	case *util.BTCAddressDecoder:
		tokenBTCCashierABI, err := abi.JSON(strings.NewReader(contract.TokenCashierBitcoinMetaData.ABI))
		if err != nil {
			log.Panicf("failed to decode token cashier abi, %+v", err)
		}
		getTransferInfo = func(data []byte) (common.Address,
			util.Address, *big.Int, *big.Int, error) {
			return getBTCTransferInfo(tokenBTCCashierABI, data, util.NewBTCAddressDecoder(&chaincfg.TestNet3Params))
		}
		logTopic = tokenBTCCashierABI.Events["Receipt"].ID.Bytes()
	case *util.ETHAddressDecoder:
		getTransferInfo = func(data []byte) (common.Address,
			util.Address, *big.Int, *big.Int, error) {
			return getETHTransferInfo(data, util.NewETHAddressDecoder())
		}
		logTopic = _ReceiptEventTopic.Bytes()
	default:
		return nil, errors.Errorf("unsupported address decoder %T", addrDecoder)
	}
	log.Printf("the address of recipient is monitored: %s\n", cashierContractAddr.String())
	log.Printf("the event topic %s is monitored\n", hex.EncodeToString(logTopic))

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
		func(startHeight uint64, endHeight uint64) ([]AbstractTransfer, error) {
			response, err := iotexClient.API().GetLogs(context.Background(), &iotexapi.GetLogsRequest{
				Filter: &iotexapi.LogsFilter{
					Address: []string{cashierContractAddr.String()},
					Topics: []*iotexapi.Topics{
						{
							Topic: [][]byte{
								logTopic,
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
			transfers := []AbstractTransfer{}
			if len(response.Logs) > 0 {
				log.Printf("\t%d transfers fetched from %d to %d\n", len(response.Logs), startHeight, endHeight)
				for _, transferLog := range response.Logs {
					if !bytes.Equal(logTopic, transferLog.Topics[0]) {
						return nil, errors.Errorf("Wrong event topic %s, %s expected", transferLog.Topics[0], _ReceiptEventTopic)
					}
					senderAddr, recipient, amount, fee, err := getTransferInfo(transferLog.Data)
					if err != nil {
						if errors.Cause(err) == errInvalidRecipient {
							log.Printf("Invalid recipient: %s Log data %x\n", err.Error(), transferLog.Data)
							continue
						}
						return nil, errors.Wrap(err, "failed to get transfer info")
					}
					receipt, err := iotexClient.API().GetReceiptByAction(context.Background(), &iotexapi.GetReceiptByActionRequest{
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
						if tokenIoAddr.String() == l.ContractAddress && common.BytesToHash(l.Topics[0]) == _TransferEventTopic && (common.BytesToAddress(l.Topics[1]) == senderAddr || cashier == common.BytesToAddress(l.Topics[1])) {
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

					transfers = append(transfers, &Transfer{
						cashier:     cashier,
						token:       tokenAddr,
						index:       new(big.Int).SetBytes(transferLog.Topics[2]).Uint64(),
						sender:      senderAddr,
						recipient:   recipient,
						amount:      amount,
						fee:         fee,
						blockHeight: transferLog.BlkHeight,
						txHash:      common.BytesToHash(transferLog.ActHash),
					})
				}
			}
			return transfers, nil
		},
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

var (
	errInvalidRecipient = errors.New("invalid recipient")
)

func getETHTransferInfo(data []byte, addrDecoder util.AddressDecoder) (senderAddrr common.Address,
	recipient util.Address, amount *big.Int, fee *big.Int, err error) {
	if len(data) != 128 {
		err = errors.Errorf("Invalid data length %d, 128 expected", len(data))
		return
	}
	senderAddrr = common.BytesToAddress(data[:32])
	recipient, err = addrDecoder.DecodeBytes(data[32:64])
	if err != nil {
		return
	}
	amount = new(big.Int).SetBytes(data[64:96])
	fee = new(big.Int).SetBytes(data[96:128])

	return
}

func getBTCTransferInfo(cashierABI abi.ABI, data []byte, addrDecoder util.AddressDecoder) (senderAddrr common.Address,
	recipient util.Address, amount *big.Int, fee *big.Int, err error) {
	var event struct {
		Sender    common.Address
		Recipient string
		Amount    *big.Int
		Fee       *big.Int
	}

	if err = cashierABI.UnpackIntoInterface(&event, "Receipt", data); err != nil {
		return
	}

	senderAddrr = event.Sender
	recipient, err = addrDecoder.DecodeString(event.Recipient)
	if err != nil {
		err = errors.Wrap(errInvalidRecipient,
			fmt.Sprintf("failed to decode recipient %s, %s", event.Recipient, err.Error()))
		return
	}
	amount = event.Amount
	fee = event.Fee

	return
}
