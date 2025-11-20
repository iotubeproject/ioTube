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
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
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

		client                     *ethclient.Client
		validator                  validatorContract
		witnessListContractMapping map[common.Address]witnessListInfo
	}

	witnessListInfo struct {
		witnessListContract *contract.AddressListCaller
		witnesses           map[string]bool
	}

	validatorContract interface {
		WitnessListsAddr(*bind.CallOpts, common.Address) ([]common.Address, error)
		SettledHeight(common.Hash) (*big.Int, error)
		SettledTransaction(uint64, common.Hash) (common.Hash, error)
		SubmitTransfer(*bind.TransactOpts, *Transfer, ...[]byte) (*types.Transaction, error)
	}

	validatorWithPayload      contract.TransferValidatorWithPayload
	validatorWithoutPayload   contract.TransferValidator
	validatorForSolana        contract.TransferValidatorForSolana
	validatorWitnessCommittee contract.TransferValidatorV3
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
	case WitnessCommittee:
		validator, err := contract.NewTransferValidatorV3(addr, client)
		if err != nil {
			return nil, err
		}
		return (*validatorWitnessCommittee)(validator), nil
	default:
		return nil, errors.New("")
	}
}

func (v *validatorWithoutPayload) WitnessListsAddr(callOpts *bind.CallOpts, _ common.Address) ([]common.Address, error) {
	addr, err := v.WitnessList(callOpts)
	if err != nil {
		return nil, err
	}
	return []common.Address{addr}, nil
}

func (v *validatorWithPayload) WitnessListsAddr(callOpts *bind.CallOpts, _ common.Address) ([]common.Address, error) {
	addr, err := v.WitnessList(callOpts)
	if err != nil {
		return nil, err
	}
	return []common.Address{addr}, nil
}

func (v *validatorForSolana) WitnessListsAddr(callOpts *bind.CallOpts, _ common.Address) ([]common.Address, error) {
	addr, err := v.WitnessList(callOpts)
	if err != nil {
		return nil, err
	}
	return []common.Address{addr}, nil
}

func (v *validatorWitnessCommittee) WitnessListsAddr(callOpts *bind.CallOpts, tokenAddr common.Address) ([]common.Address, error) {
	addrs, err := v.GetWitnessLists(callOpts, tokenAddr)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.New("no witness list found")
	}
	for _, addr := range addrs {
		if addr == zeroAddress {
			return nil, errors.New("witness list address is zero address")
		}
	}
	return addrs, nil
}

func (v *validatorWithPayload) SettledHeight(id common.Hash) (*big.Int, error) {
	return v.Settles(&bind.CallOpts{}, id)
}

func (v *validatorWitnessCommittee) SettledHeight(id common.Hash) (*big.Int, error) {
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

func (v *validatorWitnessCommittee) SettledTransaction(height uint64, id common.Hash) (common.Hash, error) {
	iter, err := v.FilterSettled(
		&bind.FilterOpts{Start: height},
		[][32]byte{id},
	)
	if err != nil {
		return common.Hash{}, err
	}
	if !iter.Next() {
		if err := iter.Close(); err != nil {
			// ignore close error
		}
		return common.Hash{}, ethereum.NotFound
	}
	txHash := iter.Event.Raw.TxHash
	if err := iter.Close(); err != nil {
		// ignore close error
	}
	return txHash, nil
}

func (v *validatorWithPayload) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signaturesArr ...[]byte) (*types.Transaction, error) {
	if len(signaturesArr) != 1 {
		return nil, errors.New("invalid signature length")
	}
	signatures := signaturesArr[0]
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

func (v *validatorWithoutPayload) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signaturesArr ...[]byte) (*types.Transaction, error) {
	if len(signaturesArr) != 1 {
		return nil, errors.New("invalid signature length")
	}
	signatures := signaturesArr[0]
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

func (v *validatorForSolana) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signaturesArr ...[]byte) (*types.Transaction, error) {
	if len(signaturesArr) != 1 {
		return nil, errors.New("invalid signature length")
	}
	signatures := signaturesArr[0]
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

func (v *validatorWitnessCommittee) SubmitTransfer(opts *bind.TransactOpts, transfer *Transfer, signaturesArr ...[]byte) (*types.Transaction, error) {
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
	// Estimate gas first
	opts.NoSend = true
	nonceSet := opts.Nonce != nil
	tx, err := v.Submit(opts, cashier, token, new(big.Int).SetUint64(transfer.index), sender, recipient, transfer.amount, signaturesArr, transfer.payload)
	if err != nil {
		return nil, err
	}
	if !nonceSet {
		opts.Nonce = nil
	}
	opts.NoSend = false
	opts.GasLimit = tx.Gas() * 11 / 10

	return v.Submit(opts, cashier, token, new(big.Int).SetUint64(transfer.index), sender, recipient, transfer.amount, signaturesArr, transfer.payload)
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

		client:                     client,
		validator:                  validator,
		witnessListContractMapping: make(map[common.Address]witnessListInfo),
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

func (tv *transferValidatorOnEthereum) refresh(witnessLists []common.Address, callOpts *bind.CallOpts) error {
	for _, witnessList := range witnessLists {
		info, ok := tv.witnessListContractMapping[witnessList]
		if !ok {
			caller, err := contract.NewAddressListCaller(witnessList, tv.client)
			if err != nil {
				return err
			}
			info = witnessListInfo{
				witnessListContract: caller,
				witnesses:           make(map[string]bool),
			}
		}
		witnesses, err := fetchWitnessesFromContract(info.witnessListContract, callOpts)
		if err != nil {
			return err
		}
		info.witnesses = witnesses
		tv.witnessListContractMapping[witnessList] = info
	}
	return nil
}

type witnessListContract interface {
	NumOfActive(opts *bind.CallOpts) (*big.Int, error)
	Count(opts *bind.CallOpts) (*big.Int, error)
	GetActiveItems(*bind.CallOpts, *big.Int, uint8) (struct {
		Count *big.Int
		Items []common.Address
	}, error)
}

func fetchWitnessesFromContract(contract witnessListContract, callOpts *bind.CallOpts) (map[string]bool, error) {
	numOfActive, err := contract.NumOfActive(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get number of active witnesses")
	}
	if numOfActive.Cmp(big.NewInt(0)) == 0 {
		return make(map[string]bool), nil
	}
	count, err := contract.Count(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of witnesses")
	}
	offset := big.NewInt(0)
	limit := uint8(100)
	witnesses := make([]common.Address, 0, int(numOfActive.Int64()))
	for offset.Cmp(count) < 0 && big.NewInt(int64(len(witnesses))).Cmp(numOfActive) < 0 {
		result, err := contract.GetActiveItems(callOpts, offset, limit)
		if err != nil {
			return nil, errors.Wrap(err, "failed to query list")
		}
		witnesses = append(witnesses, result.Items[0:result.Count.Int64()]...)
		offset.Add(offset, big.NewInt(int64(limit)))
	}

	activeWitnesses := make(map[string]bool)
	for _, w := range witnesses {
		activeWitnesses[w.Hex()] = true
	}

	return activeWitnesses, nil
}

func (tv *transferValidatorOnEthereum) filterValidWitnesses(witnesses []*Witness, witnessLists []common.Address) ([][]*Witness, [][]byte, error) {
	validWitnesses := [][]*Witness{}
	validSignatures := [][]byte{}

	for _, witnessListAddr := range witnessLists {
		witnessListInfo, exist := tv.witnessListContractMapping[witnessListAddr]
		if !exist {
			return nil, nil, errors.Errorf("witness list contract not found for witness list %s", witnessListAddr.Hex())
		}

		currentListWitnesses := []*Witness{}
		signatures := []byte{}
		numOfValidSignatures := 0

		for _, witness := range witnesses {
			witnessAddr := witness.Address()
			if !witnessListInfo.witnesses[witnessAddr.Hex()] {
				continue
			}
			signatures = append(signatures, witness.signature...)
			currentListWitnesses = append(currentListWitnesses, witness)
			numOfValidSignatures++
		}
		if numOfValidSignatures*3 <= len(witnessListInfo.witnesses)*2 {
			return nil, nil, errInsufficientWitnesses
		}
		validWitnesses = append(validWitnesses, currentListWitnesses)
		validSignatures = append(validSignatures, signatures)
	}

	return validWitnesses, validSignatures, nil
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
	settledHeight, err := tv.validator.SettledHeight(transfer.id)
	if err != nil {
		return StatusOnChainUnknown, err
	}
	if settledHeight.Cmp(big.NewInt(0)) > 0 {
		if new(big.Int).Add(settledHeight, big.NewInt(int64(tv.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return StatusOnChainNotConfirmed, nil
		}
		tx, _, err := tv.client.TransactionByHash(context.Background(), common.BytesToHash(transfer.txHash))
		if errors.Cause(err) == ethereum.NotFound {
			var txHash common.Hash
			txHash, err = tv.validator.SettledTransaction(
				settledHeight.Uint64(),
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
		settleBlockHeader, err := tv.client.HeaderByNumber(context.Background(), settledHeight)
		if err != nil {
			return StatusOnChainNotConfirmed, err
		}
		transfer.timestamp = time.Unix(int64(settleBlockHeader.Time), 0)
		return StatusOnChainSettled, nil
	}
	if transfer.nonce < nonce {
		return StatusOnChainNonceOverwritten, nil
	}
	if transfer.nonce != nonce {
		return StatusOnChainNotConfirmed, nil
	}
	if transfer.updateTime.Before(time.Now().Add(-10 * time.Minute)) {
		log.Printf("transfer %s with nonce %d needs speed up, %s %s %d\n", transfer.id, transfer.nonce, transfer.updateTime.String(), time.Now(), nonce)
		return StatusOnChainNeedSpeedUp, nil
	}
	if transfer.updateTime.After(time.Now().Add(-10 * time.Minute)) {
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
	callOpts, err := tv.callOpts()
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	witnessLists, err := tv.validator.WitnessListsAddr(callOpts, transfer.token.Address().(common.Address))
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	if err := tv.refresh(witnessLists, callOpts); err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	_, signatures, err := tv.filterValidWitnesses(witnesses, witnessLists)
	if err != nil {
		if errors.Cause(err) == errInsufficientWitnesses {
			return common.Hash{}, common.Address{}, 0, nil, errInsufficientWitnesses
		}
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	var privateKey *ecdsa.PrivateKey
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
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "current gas price %s is not significantly larger than old gas price %s", gasPrice, transfer.gasPrice)
		}
		// TODO: increase price in tOpts when speeding up?
		tOpts.Nonce = big.NewInt(0).SetUint64(transfer.nonce)
	}
	transaction, err := tv.validator.SubmitTransfer(tOpts, transfer, signatures...)
	switch errors.Cause(err) {
	case nil:
		return transaction.Hash(), tOpts.From, transaction.Nonce(), transaction.GasPrice(), nil
	case txpool.ErrUnderpriced:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	case ethereum.NotFound:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	default:
		if strings.Contains(err.Error(), "could not replace existing tx") {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
		}
		if strings.Contains(err.Error(), "transfer has been settled") {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(vm.ErrExecutionReverted, err.Error())
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
