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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/pkg/errors"
)

type witnessManagerOnEthereum struct {
	mu                 sync.RWMutex
	witnessManagerAddr common.Address
	witnessManager     *contract.WitnessManager
	confirmBlockNumber uint16
	privateKeys        []*ecdsa.PrivateKey

	client            *ethclient.Client
	defaultGasPrice   *big.Int
	gasPriceLimit     *big.Int
	gasPriceHardLimit *big.Int
	gasPriceDeviation *big.Int
	gasPriceGap       *big.Int
	support1559       bool

	chainID             *big.Int
	witnessListContract *contract.AddressListCaller
	witnesses           map[string]bool
}

func NewWitnessManagerOnEthereum(
	client *ethclient.Client,
	privateKeys []*ecdsa.PrivateKey,
	confirmBlockNumber uint16,
	defaultGasPrice *big.Int,
	gasPriceLimit *big.Int,
	gasPriceHardLimit *big.Int,
	gasPriceDeviation *big.Int,
	gasPriceGap *big.Int,
	support1559 bool,
	witnessManagerAddr common.Address,
) (WitnessManager, error) {
	witnessManager, err := contract.NewWitnessManager(witnessManagerAddr, client)
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
	log.Printf("Create witness manager for chain %d\n", chainID)

	witnessContractAddr, err := witnessManager.WitnessList(nil)
	if err != nil {
		return nil, err
	}
	witnessListContract, err := contract.NewAddressListCaller(witnessContractAddr, client)
	if err != nil {
		return nil, err
	}

	return &witnessManagerOnEthereum{
		witnessManagerAddr:  witnessManagerAddr,
		witnessManager:      witnessManager,
		confirmBlockNumber:  confirmBlockNumber,
		privateKeys:         privateKeys,
		client:              client,
		defaultGasPrice:     defaultGasPrice,
		gasPriceLimit:       gasPriceLimit,
		gasPriceHardLimit:   gasPriceHardLimit,
		gasPriceDeviation:   gasPriceDeviation,
		gasPriceGap:         gasPriceGap,
		support1559:         support1559,
		chainID:             chainID,
		witnessListContract: witnessListContract,
	}, nil
}

func (w *witnessManagerOnEthereum) Size() int {
	return len(w.privateKeys)
}

func (w *witnessManagerOnEthereum) Address() common.Address {
	return w.witnessManagerAddr
}

// TODO: refactoring with transferValidatorOnEthereum.Check()
func (w *witnessManagerOnEthereum) Check(candidates *WitnessCandidates) (StatusOnChainType, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if candidates.relayer == zeroAddress {
		return StatusOnChainUnknown, errors.New("relayer is null")
	}
	if candidates.witnessManager != w.witnessManagerAddr {
		return StatusOnChainUnknown, errors.New("witness manager is not the same")
	}

	// Fetch confirmed nonce before all the other checks
	nonce, err := w.client.NonceAt(context.Background(), candidates.relayer, nil)
	if err != nil {
		return StatusOnChainUnknown, errors.Wrap(err, "failed to get nonce")
	}
	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return StatusOnChainUnknown, errors.Wrap(err, "failed to get epoch on contract")
	}
	switch {
	case epochOnContract >= candidates.epoch:
		header, err := w.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return StatusOnChainUnknown, errors.Wrap(err, "failed to get header")
		}
		if new(big.Int).Add(big.NewInt(int64(candidates.blockHeight)), big.NewInt(int64(w.confirmBlockNumber))).Cmp(header.Number) > 0 {
			return StatusOnChainNotConfirmed, nil
		}
		return StatusOnChainSettled, nil
	case candidates.nonce < nonce:
		return StatusOnChainNonceOverwritten, nil
	case candidates.nonce > nonce:
		return StatusOnChainNotConfirmed, nil
	case candidates.updateTime.Before(time.Now().Add(-10 * time.Minute)):
		log.Printf("witness candidates %s with nonce %d needs speed up, %s %s %d\n", candidates.id, candidates.nonce, candidates.updateTime.String(), time.Now(), nonce)
		return StatusOnChainNeedSpeedUp, nil
	default:
		return StatusOnChainNotConfirmed, nil
	}
}

func (w *witnessManagerOnEthereum) Submit(candidates *WitnessCandidates, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.submit(candidates, witnesses, false)
}

func (w *witnessManagerOnEthereum) SpeedUp(candidates *WitnessCandidates, witnesses []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.gasPriceGap == nil || w.gasPriceGap.Cmp(big.NewInt(0)) < 0 {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "gas price gas is not set")
	}
	return w.submit(candidates, witnesses, true)
}

func (w *witnessManagerOnEthereum) submit(candidates *WitnessCandidates, witnesses []*Witness, isSpeedUp bool) (common.Hash, common.Address, uint64, *big.Int, error) {
	if err := w.validateEpoch(candidates.epoch); err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errInvalidData, err.Error())
	}

	if err := w.refresh(); err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	signatures := [][]byte{}
	numOfValidSignatures := 0
	for _, witness := range witnesses {
		if !w.isActiveWitness(witness.Address()) {
			log.Printf("witness %s is inactive\n", witness.Address().Hex())
			continue
		}
		signatures = append(signatures, witness.signature)
		numOfValidSignatures++
	}
	if numOfValidSignatures*3 <= len(w.witnesses)*2 {
		return common.Hash{}, common.Address{}, 0, nil, errInsufficientWitnesses
	}

	var privateKey *ecdsa.PrivateKey
	var err error
	if isSpeedUp {
		privateKey, err = w.privateKeyOfRelayer(candidates.relayer)
		if err != nil {
			return common.Hash{}, common.Address{}, 0, nil, err
		}
	} else {
		privateKey = w.privateKeys[candidates.epoch%uint64(len(w.privateKeys))]
	}
	tOpts, err := w.transactionOpts(privateKey, candidates.updateTime)
	if err != nil {
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	}
	if isSpeedUp {
		var gasPrice *big.Int
		if w.support1559 {
			gasPrice = tOpts.GasFeeCap
		} else {
			gasPrice = tOpts.GasPrice
		}
		if new(big.Int).Sub(gasPrice, candidates.gasPrice).Cmp(w.gasPriceGap) < 0 {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrapf(errNoncritical, "current gas price %s is not significantly larger than old gas price %s", gasPrice, candidates.gasPrice)
		}
		tOpts.Nonce = new(big.Int).SetUint64(candidates.nonce)
	}

	witnessesToAdd, witnessesToRemove := candidates.Witnesses()
	transaction, err := w.witnessManager.ProposeWitnesses(tOpts, candidates.epoch, witnessesToAdd, witnessesToRemove, signatures)

	switch errors.Cause(err) {
	case nil:
		return transaction.Hash(), tOpts.From, transaction.Nonce(), transaction.GasPrice(), nil
	case txpool.ErrUnderpriced, ethereum.NotFound:
		return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
	default:
		if strings.Contains(err.Error(), "could not replace existing tx") {
			return common.Hash{}, common.Address{}, 0, nil, errors.Wrap(errNoncritical, err.Error())
		}
		return common.Hash{}, common.Address{}, 0, nil, err
	}
}

func (w *witnessManagerOnEthereum) validateEpoch(epoch uint64) error {
	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on contract")
	}
	epochInterval, err := w.witnessManager.EpochInterval(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch interval")
	}
	if epoch != epochOnContract+epochInterval {
		log.Printf("epoch number is invalid, candidates.epoch: %d, epochOnContract: %d, epochInterval: %d\n", epoch, epochOnContract, epochInterval)
		return errors.New("epoch number is invalid")
	}

	return nil
}

// TODO: refacotring with transferValidatorOnEthereum.transactionOpts()
func (w *witnessManagerOnEthereum) transactionOpts(privateKey *ecdsa.PrivateKey, ts time.Time) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, w.chainID)
	if err != nil {
		return nil, err
	}
	opts.Value = big.NewInt(0)
	gasPrice, err := w.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get suggested gas price")
	}
	if gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice = w.defaultGasPrice
	}
	if w.gasPriceDeviation != nil && new(big.Int).Add(w.gasPriceDeviation, gasPrice).Sign() > 0 {
		gasPrice = new(big.Int).Add(w.gasPriceDeviation, gasPrice)
	}
	gasPriceLimit := w.gasPriceLimit
	if time.Now().Before(ts.Add(30 * time.Minute)) {
		gasPriceLimit = w.gasPriceHardLimit
	}
	if gasPrice.Cmp(gasPriceLimit) >= 0 {
		return nil, errors.Wrapf(errGasPriceTooHigh, "suggested gas price %d > limit %d", gasPrice, gasPriceLimit)
	}
	if w.support1559 {
		opts.GasFeeCap = gasPrice
		gasTipCap, err := w.client.SuggestGasTipCap(context.Background())
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

// TODO: refactoring with transferValidatorOnEthereum.refresh()
func (w *witnessManagerOnEthereum) refresh() error {
	witnesses, err := fetchWitnessesFromContract(w.witnessListContract, nil)
	if err != nil {
		return err
	}
	w.witnesses = witnesses
	return nil
}

func (w *witnessManagerOnEthereum) isActiveWitness(witness common.Address) bool {
	val, ok := w.witnesses[witness.Hex()]

	return ok && val
}

func (w *witnessManagerOnEthereum) privateKeyOfRelayer(relayer common.Address) (*ecdsa.PrivateKey, error) {
	for _, pk := range w.privateKeys {
		if relayer == crypto.PubkeyToAddress(pk.PublicKey) {
			return pk, nil
		}
	}
	return nil, errors.Errorf("no private key for relayer %s", relayer.Hex())
}
