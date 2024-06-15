// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

// transferValidatorOnIoTeX defines the transfer validator
type transferValidatorOnIoTeX struct {
	mu            sync.RWMutex
	gasLimit      uint64
	gasPrice      *big.Int
	bonus         *big.Int
	bonusTokens   map[string]*big.Int
	bonusRecorder map[string]time.Time

	relayerAddr           address.Address
	validatorContractAddr address.Address
	version               Version

	client                 iotex.AuthedClient
	validatorContract      iotex.Contract
	unpack                 func(name string, data []byte) ([]interface{}, error)
	witnessListContract    iotex.Contract
	witnessListContractABI abi.ABI
	witnesses              map[string]bool
}

// newTransferValidatorOnIoTeX creates a new TransferValidator on IoTeX
func newTransferValidatorOnIoTeX(
	client iotex.AuthedClient,
	version Version,
	validatorContractAddr address.Address,
	bonusTokens map[string]*big.Int,
	bonus *big.Int,
) (*transferValidatorOnIoTeX, error) {
	validatorContractIoAddr, err := address.FromBytes(validatorContractAddr.Bytes())
	if err != nil {
		return nil, err
	}
	var validatorABI abi.ABI
	switch version {
	case NoPayload:
		validatorABI, err = abi.JSON(strings.NewReader(contract.TransferValidatorABI))
	case Payload:
		validatorABI, err = abi.JSON(strings.NewReader(contract.TransferValidatorWithPayloadABI))
	}
	if err != nil {
		return nil, err
	}
	validatorContract := client.Contract(validatorContractAddr, validatorABI)

	data, err := validatorContract.Read("witnessList").Call(context.Background())
	if err != nil {
		return nil, err
	}

	ret, err := validatorABI.Unpack("witnessList", data.Raw)
	if err != nil {
		return nil, err
	}
	witnessContractAddr, ok := ret[0].(common.Address)
	if !ok {
		return nil, errors.Errorf("invalid type %s, common.Address expected", reflect.TypeOf(ret[0]))
	}
	witnessContractIoAddr, err := address.FromBytes(witnessContractAddr.Bytes())
	if err != nil {
		return nil, err
	}
	witnessContractABI, err := abi.JSON(strings.NewReader(contract.AddressListABI))
	if err != nil {
		return nil, err
	}

	return &transferValidatorOnIoTeX{
		gasLimit:      2000000,
		gasPrice:      big.NewInt(1000000000000),
		bonus:         bonus,
		bonusTokens:   bonusTokens,
		bonusRecorder: map[string]time.Time{},

		relayerAddr:           client.Account().Address(),
		validatorContractAddr: validatorContractIoAddr,
		version:               version,

		client:                 client,
		validatorContract:      validatorContract,
		unpack:                 validatorABI.Unpack,
		witnessListContract:    client.Contract(witnessContractIoAddr, witnessContractABI),
		witnessListContractABI: witnessContractABI,
	}, nil
}

func (tv *transferValidatorOnIoTeX) Address() common.Address {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return common.BytesToAddress(tv.validatorContractAddr.Bytes())
}

func (tv *transferValidatorOnIoTeX) refresh() error {
	witnesses := []common.Address{}
	countData, err := tv.witnessListContract.Read("count").Call(context.Background())
	if err != nil {
		return err
	}
	ret, err := countData.Unmarshal()
	if err != nil {
		return err
	}
	count, ok := ret[0].(*big.Int)
	if !ok {
		return errors.Errorf("invalid type %s, *big.Int expected", reflect.TypeOf(ret[0]))
	}
	offset := big.NewInt(0)
	limit := uint8(10)
	for offset.Cmp(count) < 0 {
		data, err := tv.witnessListContract.Read("getActiveItems", offset, limit).Call(context.Background())
		if err != nil {
			return err
		}
		ret, err := tv.witnessListContractABI.Unpack("getActiveItems", data.Raw)
		if err != nil {
			return err
		}
		count, ok := ret[0].(*big.Int)
		if !ok {
			return errors.Errorf("invalid type %s, *big.Int expected", reflect.TypeOf(ret[0]))
		}
		items, ok := ret[1].([]common.Address)
		if !ok {
			return errors.Errorf("invalid type %s, []common.Address expected", reflect.TypeOf(ret[0]))
		}

		witnesses = append(witnesses, items[:int(count.Int64())]...)
		offset.Add(offset, big.NewInt(int64(limit)))
	}
	// log.Println("refresh Witnesses on IoTeX")
	activeWitnesses := make(map[string]bool)
	for _, w := range witnesses {
		_, err := address.FromBytes(w.Bytes())
		if err != nil {
			return err
		}
		// log.Println("\t" + addr.String())
		activeWitnesses[w.Hex()] = true
	}
	tv.witnesses = activeWitnesses
	gasPrice, err := tv.client.API().SuggestGasPrice(context.Background(), &iotexapi.SuggestGasPriceRequest{})
	if err != nil {
		return err
	}
	tv.gasPrice = new(big.Int).SetUint64(gasPrice.GasPrice)

	return nil
}

func (tv *transferValidatorOnIoTeX) isActiveWitness(witness common.Address) bool {
	val, ok := tv.witnesses[witness.Hex()]

	return ok && val
}

func (tv *transferValidatorOnIoTeX) Size() int {
	return 1
}

func (tv *transferValidatorOnIoTeX) SendBonus(transfer *Transfer) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	threshold, ok := tv.bonusTokens[transfer.token.Hex()]
	if !ok || transfer.amount.Cmp(threshold) < 0 {
		return nil
	}
	return tv.sendBonus(transfer.recipient)
}

// Check returns true if a transfer has been settled
func (tv *transferValidatorOnIoTeX) Check(transfer *Transfer) (StatusOnChainType, error) {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	settleHeightData, err := tv.validatorContract.Read("settles", transfer.id).Call(context.Background())
	if err != nil {
		return StatusOnChainUnknown, err
	}
	ret, err := tv.unpack("settles", settleHeightData.Raw)
	if err != nil {
		return StatusOnChainUnknown, err
	}
	settleHeight, ok := ret[0].(*big.Int)
	if !ok {
		return StatusOnChainUnknown, errors.Errorf("invalid type %s", reflect.TypeOf(ret[0]))
	}

	if settleHeight.Cmp(big.NewInt(0)) > 0 {
		response, err := tv.client.API().GetReceiptByAction(
			context.Background(),
			&iotexapi.GetReceiptByActionRequest{ActionHash: transfer.txHash.String()[2:]},
		)
		if err != nil {
			return StatusOnChainUnknown, err
		}
		transfer.gas = response.ReceiptInfo.Receipt.GasConsumed
		metaResponse, err := tv.client.API().GetBlockMetas(context.Background(), &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByHash{
				ByHash: &iotexapi.GetBlockMetaByHashRequest{
					BlkHash: response.ReceiptInfo.BlkHash,
				},
			},
		})
		if err != nil {
			return StatusOnChainUnknown, err
		}
		transfer.timestamp = metaResponse.BlkMetas[0].Timestamp.AsTime()

		return StatusOnChainSettled, nil
	}
	response, err := tv.client.API().GetReceiptByAction(context.Background(), &iotexapi.GetReceiptByActionRequest{
		ActionHash: hex.EncodeToString(transfer.txHash.Bytes()),
	})
	switch status.Code(err) {
	case codes.NotFound:
		return StatusOnChainNeedSpeedUp, nil
	case codes.OK:
		break
	default:
		return StatusOnChainUnknown, err
	}
	if response != nil {
		// no matter what the receipt status is, mark the validation as failure
		return StatusOnChainRejected, nil
	}

	return StatusOnChainNotConfirmed, nil
}

func (tv *transferValidatorOnIoTeX) sendBonus(recipient common.Address) error {
	addr, err := address.FromBytes(recipient.Bytes())
	if err != nil {
		log.Panic("failed to convert address", recipient)
	}
	accountResponse, err := tv.client.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: addr.String(),
	})
	if err != nil {
		return err
	}
	if accountResponse.AccountMeta.IsContract || accountResponse.AccountMeta.PendingNonce >= 1 {
		return nil
	}
	switch balance, ok := big.NewInt(0).SetString(accountResponse.AccountMeta.Balance, 10); {
	case !ok:
		return errors.Errorf("failed to get balance of %s", addr.String())
	case balance.Cmp(big.NewInt(0)) > 0:
		return nil
	}
	if t, ok := tv.bonusRecorder[addr.String()]; !ok || time.Now().After(t.Add(24*time.Hour)) {
		_, err = tv.client.Transfer(addr, tv.bonus).SetGasPrice(tv.gasPrice).SetGasLimit(10000).Call(context.Background())
		if err != nil {
			return err
		}
		tv.bonusRecorder[addr.String()] = time.Now()
	}
	return nil
}

func (tv *transferValidatorOnIoTeX) submit(transfer *Transfer, witnesses []*Witness, resubmit bool) (common.Hash, common.Address, uint64, *big.Int, error) {
	if err := tv.refresh(); err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	signatures := []byte{}
	numOfValidSignatures := 0
	for _, witness := range witnesses {
		if !tv.isActiveWitness(witness.addr) {
			addr, err := address.FromBytes(witness.addr.Bytes())
			if err != nil {
				return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
			}
			log.Printf("witness %s is inactive\n", addr.String())
			continue
		}
		signatures = append(signatures, witness.signature...)
		numOfValidSignatures++
	}
	if numOfValidSignatures*3 <= len(tv.witnesses)*2 {
		return common.Hash{}, common.Address{}, 0, nil, errInsufficientWitnesses
	}
	accountMeta, err := tv.relayerAccountMeta()
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "failed to get account of %s, %v", tv.relayerAddr.String(), err)
	}
	balance, ok := big.NewInt(0).SetString(accountMeta.Balance, 10)
	if !ok {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "failed to convert balance %s of account %s, %v", accountMeta.Balance, tv.relayerAddr.String(), err)
	}
	if balance.Cmp(new(big.Int).Mul(tv.gasPrice, new(big.Int).SetUint64(tv.gasLimit))) < 0 {
		util.Alert("IOTX native balance has dropped to " + balance.String() + ", please refill account for gas " + tv.relayerAddr.String())
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "insufficient balance %s of account %s", balance, tv.relayerAddr)
	}
	var nonce uint64
	if resubmit {
		nonce = transfer.nonce
	} else {
		nonce = accountMeta.PendingNonce
	}

	var actionHash hash.Hash256
	switch tv.version {
	case NoPayload:
		actionHash, err = tv.validatorContract.Execute(
			"submit",
			transfer.cashier,
			transfer.token,
			new(big.Int).SetUint64(transfer.index),
			transfer.sender,
			transfer.recipient,
			transfer.amount,
			signatures,
		).SetGasPrice(tv.gasPrice).
			SetGasLimit(tv.gasLimit).
			SetNonce(nonce).
			Call(context.Background())
	case Payload:
		actionHash, err = tv.validatorContract.Execute(
			"submit",
			transfer.cashier,
			transfer.token,
			new(big.Int).SetUint64(transfer.index),
			transfer.sender,
			transfer.recipient,
			transfer.amount,
			transfer.payload,
			signatures,
		).SetGasPrice(tv.gasPrice).
			SetGasLimit(tv.gasLimit).
			SetNonce(nonce).
			Call(context.Background())
	}
	if err != nil {
		if errors.Cause(err).Error() == "rpc error: code = Internal desc = exceeds block gas limit" {
			err = errors.Wrap(errNoncritical, err.Error())
		}
		return common.Hash{}, common.Address{}, 0, nil, err
	}

	return common.BytesToHash(actionHash[:]), common.BytesToAddress(tv.relayerAddr.Bytes()), nonce, tv.gasPrice, nil
}

// Submit submits validation for a transfer
func (tv *transferValidatorOnIoTeX) Submit(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	return tv.submit(transfer, witnesses, false)
}

func (tv *transferValidatorOnIoTeX) SpeedUp(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	return tv.submit(transfer, witnesses, true)
}

func (tv *transferValidatorOnIoTeX) relayerAccountMeta() (*iotextypes.AccountMeta, error) {
	response, err := tv.client.API().GetAccount(context.Background(), &iotexapi.GetAccountRequest{
		Address: tv.relayerAddr.String(),
	})
	if err != nil {
		return nil, err
	}
	return response.AccountMeta, nil
}
