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
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/util"
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
	case *util.SOLAddressDecoder:
		tokenSOLCashierABI, err := abi.JSON(strings.NewReader(TokenCashierBitcoinABI))
		if err != nil {
			log.Panicf("failed to decode token cashier abi, %+v", err)
		}
		getTransferInfo = func(data []byte) (common.Address,
			util.Address, *big.Int, *big.Int, error) {
			return getSOLTransferInfo(tokenSOLCashierABI, data, util.NewSOLAddressDecoder())
		}
		logTopic = tokenSOLCashierABI.Events["Receipt"].ID.Bytes()
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
		func(util.Address, *big.Int) bool {
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

func getSOLTransferInfo(cashierABI abi.ABI, data []byte, addrDecoder util.AddressDecoder) (senderAddrr common.Address,
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
		return
	}
	amount = event.Amount
	fee = event.Fee

	return
}

// TokenCashierBitcoinMetaData contains all meta data concerning the TokenCashierBitcoin contract.
var TokenCashierBitcoinMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bitcoin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bitcoin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TokenCashierBitcoinABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenCashierBitcoinMetaData.ABI instead.
var TokenCashierBitcoinABI = TokenCashierBitcoinMetaData.ABI
