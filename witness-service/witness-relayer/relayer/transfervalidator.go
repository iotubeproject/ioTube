// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
)

// TransferValidator defines the transfer validator
type TransferValidator struct {
	mu                 sync.RWMutex
	confirmBlockNumber uint8
	gasPriceLimit      *big.Int

	privateKey            *ecdsa.PrivateKey
	relayerAddr           common.Address
	validatorContractAddr common.Address

	client              *ethclient.Client
	validatorContract   *contract.TransferValidatorV2
	witnessListContract *contract.AddressListCaller
	witnesses           map[string]bool
}

// NewTransferValidator creates a new TransferValidator
func NewTransferValidator(
	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	confirmBlockNumber uint8,
	gasPriceLimit *big.Int,
	validatorContractAddr common.Address,
) (*TransferValidator, error) {
	validatorContract, err := contract.NewTransferValidatorV2(validatorContractAddr, client)
	if err != nil {
		return nil, err
	}
	tv := &TransferValidator{
		confirmBlockNumber: confirmBlockNumber,
		gasPriceLimit:      gasPriceLimit,

		privateKey:            privateKey,
		relayerAddr:           crypto.PubkeyToAddress(privateKey.PublicKey),
		validatorContractAddr: validatorContractAddr,

		client:            client,
		validatorContract: validatorContract,
	}
	callOpts, err := tv.callOpts()
	if err != nil {
		return nil, err
	}
	witnessContractAddr, err := tv.validatorContract.WhitelistedWitnesses(callOpts)
	if err != nil {
		return nil, err
	}
	tv.witnessListContract, err = contract.NewAddressListCaller(witnessContractAddr, client)
	if err != nil {
		return nil, err
	}
	if err := tv.Refresh(); err != nil {
		return nil, err
	}

	return tv, nil
}

// UnmarshalTransferProto unmalshals a witness proto
func (tv *TransferValidator) UnmarshalTransferProto(transfer *types.Transfer) (*Transfer, error) {
	return UnmarshalTransferProto(tv.validatorContractAddr, transfer)
}

// Refresh refreshes the data stored
func (tv *TransferValidator) Refresh() error {
	callOpts, err := tv.callOpts()
	if err != nil {
		return err
	}
	count, err := tv.witnessListContract.Count(callOpts)
	if err != nil {
		return errors.Wrap(err, "failed to call witness list contract")
	}
	offset := big.NewInt(0)
	limit := uint8(10)
	witnesses := []common.Address{}
	for offset.Cmp(count) < 0 {
		result, err := tv.witnessListContract.GetActiveItems(callOpts, offset, limit)
		if err != nil {
			return errors.Wrap(err, "failed to query list")
		}
		witnesses = append(witnesses, result.Items[0:result.Count.Int64()]...)
		offset.Add(offset, big.NewInt(int64(limit)))
	}

	log.Println("Refresh Witnesses on Ethereum")
	activeWitnesses := make(map[string]bool)
	for _, w := range witnesses {
		log.Println("\t" + w.String())
		activeWitnesses[w.String()] = true
	}

	tv.mu.Lock()
	defer tv.mu.Unlock()
	tv.witnesses = activeWitnesses
	return nil
}

// IsActiveWitness returns true if the input relayerAddress is an active witness on Ethereum
func (tv *TransferValidator) IsActiveWitness(witness common.Address) bool {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return tv.witnesses[witness.Hex()]
}

// NumOfActiveWitnesses returns the number of active witnesses on Ethereum
func (tv *TransferValidator) NumOfActiveWitnesses() int {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return len(tv.witnesses)
}

// Check returns true if a transfer has been settled
func (tv *TransferValidator) Check(transfer *Transfer) (confirmed bool, success bool, reset bool, err error) {
	pendingNonce, err := tv.pendingNonce()
	if err != nil {
		return false, false, false, err
	}
	header, err := tv.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, false, false, err
	}
	settleHeight, err := tv.validatorContract.Settles(&bind.CallOpts{}, transfer.id)
	if err != nil {
		return false, false, false, err
	}
	if settleHeight.Cmp(big.NewInt(0)) > 0 {
		if new(big.Int).Add(settleHeight, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return false, false, false, nil
		}
		return true, true, false, err
	}
	r, err := tv.client.TransactionReceipt(context.Background(), transfer.txHash)
	if err != nil {
		return false, false, false, err
	}
	if r != nil {
		if new(big.Int).Add(r.BlockNumber, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return false, false, false, nil
		}
		// no matter what the receipt status is, mark the validation as failure
		return true, false, false, nil
	}
	if transfer.nonce < pendingNonce {
		return true, true, true, nil
	}
	return false, false, false, nil
}

// Submit submits validation for a transfer
func (tv *TransferValidator) Submit(transfer *Transfer, signatures []byte) (common.Hash, uint64, error) {
	tOpts, err := tv.transactionOpts()
	if err != nil {
		return common.Hash{}, 0, err
	}
	transaction, err := tv.validatorContract.Submit(tOpts, transfer.cashier, transfer.token, new(big.Int).SetUint64(transfer.index), transfer.sender, transfer.recipient, transfer.amount, signatures)
	if err != nil {
		return common.Hash{}, 0, err
	}
	return transaction.Hash(), transaction.Nonce(), nil
}

func (tv *TransferValidator) transactionOpts() (*bind.TransactOpts, error) {
	opts := bind.NewKeyedTransactor(tv.privateKey)
	opts.Value = big.NewInt(0)
	opts.GasLimit = tv.gasPriceLimit.Uint64()
	gasPrice, err := tv.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get suggested gas price")
	}
	// Slightly higher than suggested gas price
	opts.GasPrice = gasPrice.Add(gasPrice, big.NewInt(1000000000))
	if opts.GasPrice.Cmp(tv.gasPriceLimit) >= 0 {
		return nil, errors.Errorf("suggested gas price is higher than limit %d", tv.gasPriceLimit)
	}
	balance, err := tv.client.BalanceAt(context.Background(), tv.relayerAddr, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get balance of operator account")
	}
	gasFee := new(big.Int).Mul(new(big.Int).SetUint64(opts.GasLimit), opts.GasPrice)
	if gasFee.Cmp(balance) > 0 {
		return nil, errors.Errorf("insuffient balance for gas fee on Ethereum")
	}
	nonce, err := tv.pendingNonce()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pending nonce for %s", tv.relayerAddr)
	}
	opts.Nonce = new(big.Int).SetUint64(nonce)

	return opts, nil
}

func (tv *TransferValidator) pendingNonce() (uint64, error) {
	return tv.client.PendingNonceAt(context.Background(), tv.relayerAddr)
}

func (tv *TransferValidator) callOpts() (*bind.CallOpts, error) {
	tipBlockHeader, err := tv.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	blockNumber := new(big.Int).Sub(tipBlockHeader.Number, big.NewInt(int64(tv.confirmBlockNumber)))
	if blockNumber.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.Errorf("Ethereum height %d is smaller than confirm height %d", tipBlockHeader.Number, tv.confirmBlockNumber)
	}

	return &bind.CallOpts{BlockNumber: blockNumber}, nil
}
