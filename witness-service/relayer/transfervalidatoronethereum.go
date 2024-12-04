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
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
)

var zeroAddress = common.Address{}

type (
	// transferValidatorOnEthereum defines the transfer validator
	transferValidatorOnEthereum struct {
		mu                 sync.RWMutex
		confirmBlockNumber uint16
		defaultGasPrice    *big.Int
		gasPriceLimit      *big.Int
		gasPriceHardLimit  *big.Int
		gasPriceDeviation  *big.Int
		gasPriceGap        *big.Int
		support1559        bool

		chainID               *big.Int
		privateKeys           []*ecdsa.PrivateKey
		validatorContractAddr common.Address

		client              *ethclient.Client
		validator           validatorContract
		witnessListContract *contract.AddressListCaller
		witnesses           map[string]bool
	}

	validatorContract interface {
		WitnessListAddr(*bind.CallOpts) (common.Address, error)
		SettledHeight(common.Hash) (*big.Int, error)
		SettledTransaction(uint64, common.Hash) (common.Hash, error)
		SubmitTransfer(*bind.TransactOpts, *Transfer, []byte) (*types.Transaction, error)
	}

	validatorWithPayload    contract.TransferValidatorWithPayload
	validatorWithoutPayload contract.TransferValidator
	validatorForSolana      contract.TransferValidatorForSolana
)

func newValidatorContract(
	version Version,
	addr common.Address,
	client *ethclient.Client,
) (validatorContract, error) {
	switch version {
	case FromSolana:
		validator, err := contract.NewTransferValidatorForSolana(addr, client)
		if err != nil {
			return nil, err
		}
		return (*validatorForSolana)(validator), nil
	case Payload:
		validator, err := contract.NewTransferValidatorWithPayload(addr, client)
		if err != nil {
			return nil, err
		}
		return (*validatorWithPayload)(validator), nil
	case NoPayload:
		validator, err := contract.NewTransferValidator(addr, client)
		if err != nil {
			return nil, err
		}
		return (*validatorWithoutPayload)(validator), nil
	default:
		return nil, errors.New("")
	}
}

func (v *validatorWithoutPayload) WitnessListAddr(callOpts *bind.CallOpts) (common.Address, error) {
	return v.WitnessList(callOpts)
}

func (v *validatorWithPayload) WitnessListAddr(callOpts *bind.CallOpts) (common.Address, error) {
	return v.WitnessList(callOpts)
}

func (v *validatorForSolana) WitnessListAddr(callOpts *bind.CallOpts) (common.Address, error) {
	return v.WitnessList(callOpts)
}

func (v *validatorWithPayload) SettledHeight(id common.Hash) (*big.Int, error) {
	return v.Settles(&bind.CallOpts{}, id)
}

func (v *validatorWithoutPayload) SettledHeight(id common.Hash) (*big.Int, error) {
	return v.Settles(&bind.CallOpts{}, id)
}

func (v *validatorForSolana) SettledHeight(id common.Hash) (*big.Int, error) {
	return v.Settles(&bind.CallOpts{}, id)
}

func (v *validatorWithPayload) SettledTransaction(height uint64, id common.Hash) (common.Hash, error) {
	end := height + 1
	iter, err := v.FilterSettled(
		&bind.FilterOpts{Start: height, End: &end},
		[][32]byte{id},
	)
	if err != nil {
		return common.Hash{}, err
	}
	if !iter.Next() {
		return common.Hash{}, ethereum.NotFound

	}
	return iter.Event.Raw.TxHash, nil
}

func (v *validatorWithoutPayload) SettledTransaction(height uint64, id common.Hash) (common.Hash, error) {
	iter, err := v.FilterSettled(
		&bind.FilterOpts{Start: height},
		[][32]byte{id},
	)
	if err != nil {
		return common.Hash{}, err
	}
	if !iter.Next() {
		return common.Hash{}, ethereum.NotFound
	}
	return iter.Event.Raw.TxHash, nil
}

func (v *validatorForSolana) SettledTransaction(height uint64, id common.Hash) (common.Hash, error) {
	iter, err := v.FilterSettled(
		&bind.FilterOpts{Start: height},
		[][32]byte{id},
	)
	if err != nil {
		return common.Hash{}, err
	}
	if !iter.Next() {
		return common.Hash{}, ethereum.NotFound
	}
	return iter.Event.Raw.TxHash, nil
}

func (v *validatorWithPayload) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signatures []byte) (*types.Transaction, error) {
	cashier, err := util.ParseEthAddress(transfer.cashier.String())
	if err != nil {
		return nil, err
	}
	token, err := util.ParseEthAddress(transfer.token.String())
	if err != nil {
		return nil, err
	}
	sender, err := util.ParseEthAddress(transfer.sender.String())
	if err != nil {
		return nil, err
	}
	recipient, err := util.ParseEthAddress(transfer.recipient.String())
	if err != nil {
		return nil, err
	}
	// opts.GasLimit = 0
	opts.NoSend = true
	nonceSet := opts.Nonce != nil
	tx, err := v.Submit(opts, cashier, token, new(big.Int).SetUint64(transfer.index), sender, recipient, transfer.amount, signatures, transfer.payload)
	if err != nil {
		return nil, err
	}
	if !nonceSet {
		opts.Nonce = nil
	}
	opts.NoSend = false
	opts.GasLimit = tx.Gas() * 11 / 10

	return v.Submit(opts, cashier, token, new(big.Int).SetUint64(transfer.index), sender, recipient, transfer.amount, signatures, transfer.payload)
}

func (v *validatorWithoutPayload) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signatures []byte) (*types.Transaction, error) {
	cashier, err := util.ParseEthAddress(transfer.cashier.String())
	if err != nil {
		return nil, err
	}
	token, err := util.ParseEthAddress(transfer.token.String())
	if err != nil {
		return nil, err
	}
	sender, err := util.ParseEthAddress(transfer.sender.String())
	if err != nil {
		return nil, err
	}
	recipient, err := util.ParseEthAddress(transfer.recipient.String())
	if err != nil {
		return nil, err
	}
	return v.Submit(opts, cashier, token, new(big.Int).SetUint64(transfer.index), sender, recipient, transfer.amount, signatures)
}

func (v *validatorForSolana) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signatures []byte) (*types.Transaction, error) {
	token, err := util.ParseEthAddress(transfer.token.String())
	if err != nil {
		return nil, err
	}
	recipient, err := util.ParseEthAddress(transfer.recipient.String())
	if err != nil {
		return nil, err
	}
	return v.Submit(opts, transfer.cashier.Bytes(), token, new(big.Int).SetUint64(transfer.index), transfer.sender.Bytes(), recipient, transfer.amount, signatures, transfer.payload)
}

// NewTransferValidatorOnEthereum creates a new TransferValidator
func NewTransferValidatorOnEthereum(
	client *ethclient.Client,
	privateKeys []*ecdsa.PrivateKey,
	confirmBlockNumber uint16,
	defaultGasPrice *big.Int,
	gasPriceLimit *big.Int,
	gasPriceHardLimit *big.Int,
	gasPriceDeviation *big.Int,
	gasPriceGap *big.Int,
	version Version,
	validatorContractAddr common.Address,
	support1559 bool,
) (*transferValidatorOnEthereum, error) {
	validator, err := newValidatorContract(version, validatorContractAddr, client)
	if err != nil {
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	if gasPriceHardLimit == nil || gasPriceHardLimit.Cmp(big.NewInt(0)) == 0 {
		gasPriceHardLimit = gasPriceLimit
	}
	log.Printf("Create transfer validator for chain %d\n", chainID)
	tv := &transferValidatorOnEthereum{
		confirmBlockNumber: confirmBlockNumber,
		defaultGasPrice:    defaultGasPrice,
		gasPriceLimit:      gasPriceLimit,
		gasPriceHardLimit:  gasPriceHardLimit,
		gasPriceDeviation:  gasPriceDeviation,
		gasPriceGap:        gasPriceGap,
		support1559:        support1559,

		chainID:               chainID,
		privateKeys:           privateKeys,
		validatorContractAddr: validatorContractAddr,

		client:    client,
		validator: validator,
	}

	callOpts, err := tv.callOpts()
	if err != nil {
		return nil, err
	}
	witnessContractAddr, err := tv.validator.WitnessListAddr(callOpts)
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
	SettledHeight, err := tv.validator.SettledHeight(transfer.id)
	if err == nil {
		if SettledHeight.Cmp(big.NewInt(0)) > 0 {
			if new(big.Int).Add(SettledHeight, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
				return StatusOnChainNotConfirmed, nil
			}
			tx, _, err := tv.client.TransactionByHash(context.Background(), common.BytesToHash(transfer.txHash))
			if errors.Cause(err) == ethereum.NotFound {
				var txHash common.Hash
				txHash, err = tv.validator.SettledTransaction(
					SettledHeight.Uint64(),
					transfer.id,
				)
				if err != nil {
					return StatusOnChainNotConfirmed, err
				}
				transfer.txHash = txHash.Bytes()
				tx, _, err = tv.client.TransactionByHash(context.Background(), common.BytesToHash(transfer.txHash))
			}
			if err != nil {
				return StatusOnChainNotConfirmed, err
			}
			transfer.gas = tx.Gas()
			settleBlockHeader, err := tv.client.HeaderByNumber(context.Background(), SettledHeight)
			if err != nil {
				return StatusOnChainUnknown, err
			}
			transfer.timestamp = time.Unix(int64(settleBlockHeader.Time), 0)
			return StatusOnChainSettled, nil
		}
		var r *types.Receipt
		r, err = tv.client.TransactionReceipt(context.Background(), common.BytesToHash(transfer.txHash))
		if err == nil {
			if r == nil {
				return StatusOnChainNotConfirmed, nil
			}
			if new(big.Int).Add(r.BlockNumber, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
				return StatusOnChainNotConfirmed, nil
			}
		}
	}
	log.Println(err)
	switch errors.Cause(err) {
	case ethereum.NotFound:
		if transfer.nonce <= nonce {
			if transfer.updateTime.Add(10 * time.Minute).Before(time.Now()) {
				return StatusOnChainNonceOverwritten, nil
			}
			return StatusOnChainNotConfirmed, nil
		}
		if transfer.updateTime.Before(time.Now().Add(-10*time.Minute)) && transfer.nonce == nonce+1 {
			log.Printf("transfer %s with nonce %d needs speed up, %s %s %d\n", transfer.id, transfer.nonce, transfer.updateTime.String(), time.Now(), nonce)
			return StatusOnChainNeedSpeedUp, nil
		}
		return StatusOnChainNotConfirmed, nil
	case nil:
		break
	default:
		return StatusOnChainUnknown, err
	}
	if transfer.updateTime.After(time.Now().Add(-20 * time.Minute)) {
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
		if !tv.isActiveWitness(witness.Address()) {
			log.Printf("witness %s is inactive\n", witness.Address().Hex())
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
	tOpts, err := tv.transactionOpts(privateKey, transfer.timestamp)
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	if isSpeedUp {
		var gasPrice *big.Int
		if tv.support1559 {
			gasPrice = tOpts.GasFeeCap
		} else {
			gasPrice = tOpts.GasPrice
		}
		if new(big.Int).Sub(gasPrice, transfer.gasPrice).Cmp(tv.gasPriceGap) < 0 {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "current gas price %s is not significantly larger than old gas price %s", tOpts.GasPrice, transfer.gasPrice)
		}
		tOpts.Nonce = big.NewInt(0).SetUint64(transfer.nonce)
	}
	transaction, err := tv.validator.SubmitTransfer(tOpts, transfer, signatures)
	switch errors.Cause(err) {
	case nil:
		return transaction.Hash(), tOpts.From, transaction.Nonce(), transaction.GasPrice(), nil
	case core.ErrUnderpriced:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	case ethereum.NotFound:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	default:
		if strings.Contains(err.Error(), "could not replace existing tx") {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
		}
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

func (tv *transferValidatorOnEthereum) transactionOpts(privateKey *ecdsa.PrivateKey, ts time.Time) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, tv.chainID)
	if err != nil {
		return nil, err
	}
	opts.Value = big.NewInt(0)
	gasPrice, err := tv.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get suggested gas price")
	}
	if gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice = tv.defaultGasPrice
	}
	if tv.gasPriceDeviation != nil && new(big.Int).Add(tv.gasPriceDeviation, gasPrice).Sign() > 0 {
		gasPrice = new(big.Int).Add(tv.gasPriceDeviation, gasPrice)
	}
	gasPriceLimit := tv.gasPriceLimit
	if time.Now().Before(ts.Add(30 * time.Minute)) {
		gasPriceLimit = tv.gasPriceHardLimit
	}
	if gasPrice.Cmp(gasPriceLimit) >= 0 {
		return nil, errors.Wrapf(errGasPriceTooHigh, "suggested gas price %d > limit %d", gasPrice, gasPriceLimit)
	}
	if tv.support1559 {
		opts.GasFeeCap = gasPrice
		gasTipCap, err := tv.client.SuggestGasTipCap(context.Background())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get suggested gas tip cap")
		}
		if gasTipCap.Cmp(big.NewInt(0)) == 0 {
			gasTipCap = big.NewInt(1)
		}
		if gasTipCap.Cmp(gasPrice) > 0 {
			return nil, errors.Errorf("suggested gas tip cap %d > gas price %d", gasTipCap, gasPrice)
		}
		opts.GasTipCap = gasTipCap
	} else {
		opts.GasPrice = gasPrice
	}

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
