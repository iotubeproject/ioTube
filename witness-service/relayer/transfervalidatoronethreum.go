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
)

// transferValidatorOnEthreum defines the transfer validator
type transferValidatorOnEthreum struct {
	mu                 sync.RWMutex
	confirmBlockNumber uint8
	gasPriceLimit      *big.Int

	privateKey            *ecdsa.PrivateKey
	relayerAddr           common.Address
	validatorContractAddr common.Address

	client              *ethclient.Client
	validatorContract   *contract.TransferValidator
	witnessListContract *contract.AddressListCaller
	witnesses           map[string]bool
}

// NewTransferValidatorOnEthereum creates a new TransferValidator
func NewTransferValidatorOnEthereum(
	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	confirmBlockNumber uint8,
	gasPriceLimit *big.Int,
	validatorContractAddr common.Address,
) (TransferValidator, error) {
	validatorContract, err := contract.NewTransferValidator(validatorContractAddr, client)
	if err != nil {
		return nil, err
	}
	tv := &transferValidatorOnEthreum{
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
	witnessContractAddr, err := tv.validatorContract.WitnessList(callOpts)
	if err != nil {
		return nil, err
	}
	tv.witnessListContract, err = contract.NewAddressListCaller(witnessContractAddr, client)
	if err != nil {
		return nil, err
	}

	return tv, nil
}

func (tv *transferValidatorOnEthreum) Address() common.Address {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return tv.validatorContractAddr
}

func (tv *transferValidatorOnEthreum) refresh() error {
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

	log.Println("refresh Witnesses on Ethereum")
	activeWitnesses := make(map[string]bool)
	for _, w := range witnesses {
		log.Println("\t" + w.Hex())
		activeWitnesses[w.Hex()] = true
	}

	tv.witnesses = activeWitnesses
	return nil
}

func (tv *transferValidatorOnEthreum) isActiveWitness(witness common.Address) bool {
	val, ok := tv.witnesses[witness.Hex()]

	return ok && val
}

// Check returns true if a transfer has been settled
func (tv *transferValidatorOnEthreum) Check(transfer *Transfer) (StatusOnChainType, error) {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	pendingNonce, err := tv.pendingNonce()
	if err != nil {
		return StatusOnChainUnknown, err
	}
	header, err := tv.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return StatusOnChainUnknown, err
	}
	settleHeight, err := tv.validatorContract.Settles(&bind.CallOpts{}, transfer.id)
	if err != nil {
		return StatusOnChainUnknown, err
	}
	if settleHeight.Cmp(big.NewInt(0)) > 0 {
		if new(big.Int).Add(settleHeight, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return StatusOnChainNotConfirmed, nil
		}
		return StatusOnChainSettled, nil
	}
	r, err := tv.client.TransactionReceipt(context.Background(), transfer.txHash)
	if err != nil {
		return StatusOnChainUnknown, err
	}
	if r != nil {
		if new(big.Int).Add(r.BlockNumber, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return StatusOnChainNotConfirmed, nil
		}
		// no matter what the receipt status is, mark the validation as failure
		return StatusOnChainRejected, nil
	}
	if transfer.nonce < pendingNonce {
		return StatusOnChainNonceOverwritten, nil
	}
	return StatusOnChainNotConfirmed, nil
}

// Submit submits validation for a transfer
func (tv *transferValidatorOnEthreum) Submit(transfer *Transfer, witnesses []*Witness) (common.Hash, uint64, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	if err := tv.refresh(); err != nil {
		return common.Hash{}, 0, err
	}
	signatures := []byte{}
	numOfValidSignatures := 0
	for _, witness := range witnesses {
		if !tv.isActiveWitness(witness.addr) {
			log.Printf("witness %s is inactive\n", witness.addr.Hex())
			continue
		}
		signatures = append(signatures, witness.signature...)
		numOfValidSignatures++
	}
	if numOfValidSignatures*3 <= len(tv.witnesses)*2 {
		return common.Hash{}, 0, errInsufficientWitnesses
	}
	tOpts, err := tv.transactionOpts(300000)
	if err != nil {
		return common.Hash{}, 0, err
	}
	transaction, err := tv.validatorContract.Submit(tOpts, transfer.cashier, transfer.token, new(big.Int).SetUint64(transfer.index), transfer.sender, transfer.recipient, transfer.amount, signatures)
	if err != nil {
		return common.Hash{}, 0, err
	}
	return transaction.Hash(), transaction.Nonce(), nil
}

func (tv *transferValidatorOnEthreum) transactionOpts(gasLimit uint64) (*bind.TransactOpts, error) {
	opts := bind.NewKeyedTransactor(tv.privateKey)
	opts.Value = big.NewInt(0)
	opts.GasLimit = gasLimit
	gasPrice, err := tv.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get suggested gas price")
	}
	if gasPrice.Cmp(tv.gasPriceLimit) >= 0 {
		return nil, errors.Wrapf(errGasPriceTooHigh, "suggested gas price %d > limit %d", gasPrice, tv.gasPriceLimit)
	}
	opts.GasPrice = gasPrice
	balance, err := tv.client.BalanceAt(context.Background(), tv.relayerAddr, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get balance of operator account")
	}
	gasFee := new(big.Int).Mul(new(big.Int).SetUint64(opts.GasLimit), opts.GasPrice)
	if gasFee.Cmp(balance) > 0 {
		return nil, errors.Errorf("insuffient balance for gas fee")
	}
	nonce, err := tv.pendingNonce()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pending nonce for %s", tv.relayerAddr)
	}
	opts.Nonce = new(big.Int).SetUint64(nonce)

	return opts, nil
}

func (tv *transferValidatorOnEthreum) pendingNonce() (uint64, error) {
	return tv.client.PendingNonceAt(context.Background(), tv.relayerAddr)
}

func (tv *transferValidatorOnEthreum) callOpts() (*bind.CallOpts, error) {
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
