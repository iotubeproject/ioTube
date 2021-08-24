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
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

var zeroAddress = common.Address{}

// transferValidatorOnEthereum defines the transfer validator
type transferValidatorOnEthereum struct {
	mu                 sync.RWMutex
	confirmBlockNumber uint8
	gasPriceLimit      *big.Int
	gasPriceDeviation  *big.Int
	gasPriceGap        *big.Int

	chainID               *big.Int
	privateKeys           []*ecdsa.PrivateKey
	validatorContractAddr common.Address

	client              *ethclient.Client
	validatorContract   *contract.TransferValidator
	witnessListContract *contract.AddressListCaller
	witnesses           map[string]bool
}

// NewTransferValidatorOnEthereum creates a new TransferValidator
func NewTransferValidatorOnEthereum(
	client *ethclient.Client,
	privateKeys []*ecdsa.PrivateKey,
	confirmBlockNumber uint8,
	gasPriceLimit *big.Int,
	gasPriceDeviation *big.Int,
	gasPriceGap *big.Int,
	validatorContractAddr common.Address,
) (TransferValidator, error) {
	validatorContract, err := contract.NewTransferValidator(validatorContractAddr, client)
	if err != nil {
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	log.Printf("Create transfer validator for chain %d\n", chainID)
	tv := &transferValidatorOnEthereum{
		confirmBlockNumber: confirmBlockNumber,
		gasPriceLimit:      gasPriceLimit,
		gasPriceDeviation:  gasPriceDeviation,
		gasPriceGap:        gasPriceGap,

		chainID:               chainID,
		privateKeys:           privateKeys,
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

func (tv *transferValidatorOnEthereum) Size() int {
	return len(tv.privateKeys)
}

func (tv *transferValidatorOnEthereum) Address() common.Address {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	return tv.validatorContractAddr
}

func (tv *transferValidatorOnEthereum) refresh() error {
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

	// log.Println("refresh Witnesses")
	activeWitnesses := make(map[string]bool)
	for _, w := range witnesses {
		// log.Println("\t" + w.Hex())
		activeWitnesses[w.Hex()] = true
	}

	tv.witnesses = activeWitnesses
	return nil
}

func (tv *transferValidatorOnEthereum) isActiveWitness(witness common.Address) bool {
	val, ok := tv.witnesses[witness.Hex()]

	return ok && val
}

// Check returns true if a transfer has been settled
func (tv *transferValidatorOnEthereum) Check(transfer *Transfer) (StatusOnChainType, error) {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	if transfer.relayer == zeroAddress {
		return StatusOnChainUnknown, errors.New("relayer is null")
	}
	// Fetch confirmed nonce before all the other checks
	nonce, err := tv.client.NonceAt(context.Background(), transfer.relayer, nil)
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
		// contract status: settled
		if new(big.Int).Add(settleHeight, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return StatusOnChainNotConfirmed, nil
		}
		tx, _, err := tv.client.TransactionByHash(context.Background(), transfer.txHash)
		if err != nil {
			return StatusOnChainUnknown, err
		}
		transfer.gas = tx.Gas()
		settleBlockHeader, err := tv.client.HeaderByNumber(context.Background(), settleHeight)
		if err != nil {
			return StatusOnChainUnknown, err
		}
		transfer.timestamp = time.Unix(int64(settleBlockHeader.Time), 0)
		return StatusOnChainSettled, nil
	}
	r, err := tv.client.TransactionReceipt(context.Background(), transfer.txHash)
	switch errors.Cause(err) {
	case ethereum.NotFound:
		if transfer.nonce < nonce {
			if transfer.updateTime.Add(5 * time.Minute).Before(time.Now()) {
				return StatusOnChainNonceOverwritten, nil
			}
			return StatusOnChainNotConfirmed, nil
		}
		if transfer.updateTime.Before(time.Now().Add(-10*time.Minute)) && transfer.nonce == nonce+1 {
			return StatusOnChainNeedSpeedUp, nil
		}
		return StatusOnChainNotConfirmed, nil
	case nil:
		break
	default:
		return StatusOnChainUnknown, err
	}
	if r == nil {
		return StatusOnChainNotConfirmed, nil
	}
	if new(big.Int).Add(r.BlockNumber, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
		return StatusOnChainNotConfirmed, nil
	}
	// no matter what the receipt status is, mark the validation as failure
	return StatusOnChainRejected, nil
}

func (tv *transferValidatorOnEthereum) privateKeyOfRelayer(relayer common.Address) (*ecdsa.PrivateKey, error) {
	for _, pk := range tv.privateKeys {
		if relayer == crypto.PubkeyToAddress(pk.PublicKey) {
			return pk, nil
		}
	}
	return nil, errors.Errorf("no private key for relayer %s", relayer.Hex())
}

func (tv *transferValidatorOnEthereum) submit(transfer *Transfer, witnesses []*Witness, isSpeedUp bool) (common.Hash, common.Address, uint64, *big.Int, error) {
	if err := tv.refresh(); err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
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
		return common.Hash{}, common.Address{}, 0, nil, errInsufficientWitnesses
	}
	var privateKey *ecdsa.PrivateKey
	var err error
	if isSpeedUp {
		privateKey, err = tv.privateKeyOfRelayer(transfer.relayer)
		if err != nil {
			return common.Hash{}, common.Address{}, 0, nil, err
		}
	} else {
		privateKey = tv.privateKeys[transfer.index%uint64(len(tv.privateKeys))]
	}
	tOpts, err := tv.transactionOpts(300000, privateKey)
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	if tv.gasPriceDeviation != nil && new(big.Int).Add(tv.gasPriceDeviation, tOpts.GasPrice).Sign() > 0 {
		tOpts.GasPrice = new(big.Int).Add(tv.gasPriceDeviation, tOpts.GasPrice)
	}
	if isSpeedUp {
		if new(big.Int).Sub(tOpts.GasPrice, transfer.gasPrice).Cmp(tv.gasPriceGap) < 0 {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "current gas price %s is not significantly larger than old gas price %s", tOpts.GasPrice, transfer.gasPrice)
		}
		tOpts.Nonce = tOpts.Nonce.SetUint64(transfer.nonce)
	}
	transaction, err := tv.validatorContract.Submit(tOpts, transfer.cashier, transfer.token, new(big.Int).SetUint64(transfer.index), transfer.sender, transfer.recipient, transfer.amount, signatures)
	switch errors.Cause(err) {
	case nil:
		return transaction.Hash(), tOpts.From, transaction.Nonce(), transaction.GasPrice(), nil
	case core.ErrUnderpriced:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	default:
		return common.Hash{}, common.Address{}, 0, nil, err
	}
}

// Submit submits validation for a transfer
func (tv *transferValidatorOnEthereum) Submit(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	return tv.submit(transfer, witnesses, false)
}

// SpeedUp creases the transaction gas price
func (tv *transferValidatorOnEthereum) SpeedUp(transfer *Transfer, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()
	if tv.gasPriceGap == nil || tv.gasPriceGap.Cmp(big.NewInt(0)) < 0 {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "gas price gas is not set")
	}

	return tv.submit(transfer, witnesses, true)
}

func (tv *transferValidatorOnEthereum) transactionOpts(gasLimit uint64, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	relayerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, tv.chainID)
	if err != nil {
		return nil, err
	}
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
	balance, err := tv.client.BalanceAt(context.Background(), relayerAddr, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get balance of operator account")
	}
	gasFee := new(big.Int).Mul(new(big.Int).SetUint64(opts.GasLimit), opts.GasPrice)
	if gasFee.Cmp(balance) > 0 {
		return nil, errors.Errorf("insuffient balance for gas fee")
	}
	nonce, err := tv.client.PendingNonceAt(context.Background(), relayerAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pending nonce for %s", relayerAddr)
	}
	opts.Nonce = new(big.Int).SetUint64(nonce)

	return opts, nil
}

func (tv *transferValidatorOnEthereum) callOpts() (*bind.CallOpts, error) {
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
