// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-core/pkg/unit"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// Auth maintains the list of witnesses and tokens
type Auth struct {
	mu                      sync.RWMutex
	lastUpdateTime          time.Time
	ethConfirmBlockNumber   uint8
	gasPriceLimitOnEthereum *big.Int

	witnessAddressOnIoTeX    address.Address
	witnessPKOnEthereum      *ecdsa.PrivateKey
	witnessAddressOnEthereum common.Address

	iotexClient iotex.AuthedClient
	ethClient   *ethclient.Client

	erc20TokenListContract        common.Address
	xrc20TokenListContract        iotex.Contract
	witnessListContractOnEthereum common.Address
	witnessListContractOnIoTeX    iotex.Contract
	witnessesOnEthereum           map[string]bool
	witnessesOnIoTeX              map[string]bool
	erc20ToXrc20                  map[string]address.Address
	xrc20ToErc20                  map[string]common.Address
}

// NewAuth creates a new auth
func NewAuth(
	ethClient *ethclient.Client,
	iotexClient iotex.AuthedClient,
	witnessPKOnEthereum *ecdsa.PrivateKey,
	witnessAddressOnIoTeX address.Address,
	ethConfirmBlockNumber uint8,
	gasPriceLimitOnEthereum *big.Int,
	witnessListContractOnEthereum common.Address,
	witnessListContractOnIoTeX address.Address,
	erc20TokenListContract common.Address,
	xrc20TokenListContract address.Address,
) (*Auth, error) {
	addressListABI, err := abi.JSON(strings.NewReader(contract.AddressListABI))
	if err != nil {
		return nil, err
	}

	return &Auth{
		ethConfirmBlockNumber:   ethConfirmBlockNumber,
		gasPriceLimitOnEthereum: gasPriceLimitOnEthereum,

		witnessAddressOnIoTeX:    witnessAddressOnIoTeX,
		witnessPKOnEthereum:      witnessPKOnEthereum,
		witnessAddressOnEthereum: crypto.PubkeyToAddress(witnessPKOnEthereum.PublicKey),

		ethClient:   ethClient,
		iotexClient: iotexClient,

		erc20TokenListContract:        erc20TokenListContract,
		xrc20TokenListContract:        iotexClient.Contract(xrc20TokenListContract, addressListABI),
		witnessListContractOnEthereum: witnessListContractOnEthereum,
		witnessListContractOnIoTeX:    iotexClient.Contract(witnessListContractOnIoTeX, addressListABI),
		witnessesOnEthereum:           make(map[string]bool),
		witnessesOnIoTeX:              make(map[string]bool),
	}, nil
}

// EthConfirmBlockNumber defines the number of blocks which is considerred as confirmed
func (auth *Auth) EthConfirmBlockNumber() *big.Int {
	return big.NewInt(int64(auth.ethConfirmBlockNumber))
}

// IoTeXClient returns the iotex client
func (auth *Auth) IoTeXClient() iotex.AuthedClient {
	return auth.iotexClient
}

// EthereumClient returns the ethereum client
func (auth *Auth) EthereumClient() *ethclient.Client {
	return auth.ethClient
}

// Erc20Tokens reutrns the erc20 tokens in whitelist
func (auth *Auth) Erc20Tokens() []common.Address {
	auth.mu.RLock()
	defer auth.mu.RUnlock()
	tokens := []common.Address{}
	for _, token := range auth.xrc20ToErc20 {
		tokens = append(tokens, token)
	}
	return tokens
}

// Xrc20Tokens returns the xrc20 tokens in whitelist
func (auth *Auth) Xrc20Tokens() []address.Address {
	auth.mu.RLock()
	defer auth.mu.RUnlock()
	tokens := []address.Address{}
	for _, token := range auth.erc20ToXrc20 {
		tokens = append(tokens, token)
	}
	return tokens
}

// NewTransactionOpts prepares the transaction opts
func (auth *Auth) NewTransactionOpts(value *big.Int, gasLimit uint64) (*bind.TransactOpts, error) {
	opts := bind.NewKeyedTransactor(auth.witnessPKOnEthereum)
	opts.Value = value
	opts.GasLimit = gasLimit
	gasPrice, err := auth.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get suggested gas price")
	}
	// Slightly higher than suggested gas price
	opts.GasPrice = gasPrice.Add(gasPrice, big.NewInt(1000000000))
	if opts.GasPrice.Cmp(auth.gasPriceLimitOnEthereum) >= 0 {
		return nil, errors.Errorf("suggested gas price is higher than limit %d", auth.gasPriceLimitOnEthereum)
	}
	balance, err := auth.ethClient.BalanceAt(context.Background(), auth.witnessAddressOnEthereum, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get balance of operator account")
	}
	gasFee := new(big.Int).Mul(new(big.Int).SetUint64(opts.GasLimit), opts.GasPrice)
	if gasFee.Cmp(balance) > 0 {
		return nil, errors.Errorf("insuffient balance for gas fee on Ethereum")
	}
	nonce, err := auth.ethClient.PendingNonceAt(context.Background(), auth.witnessAddressOnEthereum)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pending nonce for %s", auth.witnessAddressOnEthereum)
	}
	opts.Nonce = new(big.Int).SetUint64(nonce)

	return opts, nil
}

// CallOnIoTeX fills the parameters and call
func (auth *Auth) CallOnIoTeX(
	caller iotex.ExecuteContractCaller,
	gasLimit uint64,
) ([]byte, error) {
	gasPrice := big.NewInt(unit.Qev)
	res, err := auth.iotexClient.API().GetAccount(
		context.Background(),
		&iotexapi.GetAccountRequest{Address: auth.witnessAddressOnIoTeX.String()},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get account of %s", auth.witnessAddressOnIoTeX)
	}
	balance, ok := big.NewInt(0).SetString(res.AccountMeta.Balance, 10)
	if !ok {
		return nil, errors.Wrapf(err, "failed to convert balance %s of account %s", res.AccountMeta.Balance, auth.witnessAddressOnIoTeX)
	}
	if balance.Cmp(new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))) < 0 {
		util.Alert("IOTX native balance has dropped to " + balance.String() + ", please refill account for gas " + auth.witnessAddressOnIoTeX.String())
	}

	actionHash, err := caller.
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		Call(context.Background())
	if err != nil {
		return nil, err
	}
	return actionHash[:], nil
}

// CallOptsOnEthereum returns call opts of confirmed height on ethereum
func (auth *Auth) CallOptsOnEthereum() (*bind.CallOpts, error) {
	tipBlockHeader, err := auth.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	blockNumber := new(big.Int).Sub(tipBlockHeader.Number, big.NewInt(int64(auth.ethConfirmBlockNumber)))
	if blockNumber.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.Errorf("Ethereum height %d is smaller than confirm height %d", tipBlockHeader.Number, auth.ethConfirmBlockNumber)
	}

	return &bind.CallOpts{BlockNumber: blockNumber}, nil
}

func (auth *Auth) loadAddressListOnEthereum(contractAddr common.Address) ([]common.Address, error) {
	var retval []common.Address
	callOpts, err := auth.CallOptsOnEthereum()
	if err != nil {
		return nil, err
	}
	list, err := contract.NewAddressListCaller(contractAddr, auth.ethClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create caller for %s", contractAddr.String())
	}
	count, err := list.Count(callOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call address list on %s", contractAddr.String())
	}
	offset := big.NewInt(0)
	limit := uint8(10)
	retval = []common.Address{}
	for offset.Cmp(count) < 0 {
		result, err := list.GetActiveItems(callOpts, offset, limit)
		if err != nil {
			return nil, errors.Wrap(err, "failed to query list")
		}
		retval = append(retval, result.Items[0:result.Count.Int64()]...)
		offset.Add(offset, big.NewInt(int64(limit)))
	}

	return retval, nil
}

func (auth *Auth) loadAddressListOnIoTeX(c iotex.Contract) ([]address.Address, error) {
	response, err := c.Read("count").Call(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to call witness list contract")
	}
	var count *big.Int
	if err := response.Unmarshal(&count); err != nil {
		return nil, errors.Wrap(err, "failed to parse list count")
	}
	fullList := []address.Address{}
	offset := big.NewInt(0)
	limit := uint8(10)
	for offset.Cmp(count) < 0 {
		response, err := c.Read("getActiveItems", offset, limit).Call(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "failed to call witness list contract")
		}
		result := struct {
			Count *big.Int
			Items []common.Address
		}{}
		if err := response.Unmarshal(&result); err != nil {
			return nil, errors.Wrap(err, "failed to parse addresses")
		}
		for i := int64(0); i < result.Count.Int64(); i++ {
			addr, err := address.FromBytes(result.Items[i].Bytes())
			if err != nil {
				return nil, err
			}
			fullList = append(fullList, addr)
		}
		offset.Add(offset, big.NewInt(int64(limit)))
	}
	return fullList, nil
}

// LastUpdateTime returns the last update time of the component
func (auth *Auth) LastUpdateTime() time.Time {
	auth.mu.RLock()
	defer auth.mu.RUnlock()
	return auth.lastUpdateTime
}

// Refresh refreshes the data stored
func (auth *Auth) Refresh() error {
	witnessesOnIoTeX, err := auth.loadAddressListOnIoTeX(auth.witnessListContractOnIoTeX)
	if err != nil {
		return err
	}
	witnessesOnEth, err := auth.loadAddressListOnEthereum(auth.witnessListContractOnEthereum)
	if err != nil {
		return err
	}
	newXrc20ToErc20 := map[string]common.Address{}
	newErc20ToXrc20 := map[string]address.Address{}
	shadowTokenABI, err := abi.JSON(strings.NewReader(contract.ShadowTokenABI))
	if err != nil {
		return err
	}
	tokensOnIoTeX, err := auth.loadAddressListOnIoTeX(auth.xrc20TokenListContract)
	if err != nil {
		return err
	}
	for _, token := range tokensOnIoTeX {
		ioAddr, err := address.FromBytes(token.Bytes())
		if err != nil {
			return err
		}
		if ethAddr, ok := auth.xrc20ToErc20[ioAddr.String()]; ok {
			newXrc20ToErc20[ioAddr.String()] = ethAddr
			newErc20ToXrc20[ethAddr.String()] = ioAddr
		} else {
			response, err := auth.iotexClient.ReadOnlyContract(ioAddr, shadowTokenABI).Read("coToken").Call(context.Background())
			if err != nil {
				return errors.Wrapf(err, "failed to get corresponding token of %s", ioAddr)
			}
			tokenOnEth := common.Address{}
			if err := response.Unmarshal(&tokenOnEth); err != nil {
				return errors.Wrapf(err, "failed to extract corresponding token of %s from %s", ioAddr, response)
			}
			newXrc20ToErc20[ioAddr.String()] = tokenOnEth
			if erc20Addr, ok := newErc20ToXrc20[tokenOnEth.String()]; ok {
				return errors.Wrapf(err, "two Xrc20 tokens %s and %s map to the same Erc20 %s", ioAddr, erc20Addr, tokenOnEth)
			}
			newErc20ToXrc20[tokenOnEth.String()] = ioAddr
		}
	}
	tokensOnEth, err := auth.loadAddressListOnEthereum(auth.erc20TokenListContract)
	if err != nil {
		return err
	}
	if len(tokensOnEth) != len(tokensOnIoTeX) {
		return errors.Errorf("num of tokens on eth %d is not equal to num of tokens on iotex %d", tokensOnEth, tokensOnIoTeX)
	}
	for _, token := range tokensOnEth {
		if _, ok := newErc20ToXrc20[token.String()]; !ok {
			return errors.Errorf("erc20 token %s doesn't have a match token", token)
		}
	}
	auth.mu.Lock()
	defer auth.mu.Unlock()
	auth.lastUpdateTime = time.Now()
	str := "auth data refreshed\n  Witnesses on IoTeX"
	for _, w := range witnessesOnIoTeX {
		str += "\n    " + w.String()
		auth.witnessesOnIoTeX[w.String()] = true
	}
	str += "\n  Witnesses on Ethereum"
	for _, w := range witnessesOnEth {
		str += "\n    " + w.String()
		auth.witnessesOnEthereum[w.String()] = true
	}
	auth.erc20ToXrc20 = newErc20ToXrc20
	auth.xrc20ToErc20 = newXrc20ToErc20
	str += "\n  Token pairs"
	for key, value := range auth.erc20ToXrc20 {
		str += "\n    " + key + "<=>" + value.String()
	}
	log.Println(str)
	return nil
}

// IsActiveWitnessOnIoTeX returns true if the input address is an active witness on IoTeX
func (auth *Auth) IsActiveWitnessOnIoTeX(witness address.Address) bool {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.witnessesOnIoTeX[witness.String()]
}

// IsActiveWitnessOnEthereum returns true if the input address is an active witness on Ethereum
func (auth *Auth) IsActiveWitnessOnEthereum() bool {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.witnessesOnEthereum[auth.witnessAddressOnEthereum.String()]
}

// CorrespondingXrc20Token returns the corresponding Xrc20 token address on IoTeX
func (auth *Auth) CorrespondingXrc20Token(erc20 common.Address) (address.Address, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	xrc20, ok := auth.erc20ToXrc20[erc20.String()]
	if !ok {
		return nil, errors.Errorf("cannot find corresponding XRC20 token address for %s", erc20)
	}

	return xrc20, nil
}

// CorrespondingErc20Token returns the corresponding Erc20 token address on Ethereum
func (auth *Auth) CorrespondingErc20Token(xrc20 address.Address) (common.Address, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	erc20, ok := auth.xrc20ToErc20[xrc20.String()]
	if !ok {
		return common.Address{}, errors.Errorf("cannot find corresponding ERC20 token address for %s", erc20)
	}

	return erc20, nil
}
