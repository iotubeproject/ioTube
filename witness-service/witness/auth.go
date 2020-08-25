// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

// Auth maintains the list of witnesses and tokens
type Auth struct {
	mu                            sync.RWMutex
	iotexClient                   iotex.AuthedClient
	ethClientPool                 *EthClientPool
	erc20TokenListContract        common.Address
	xrc20TokenListContract        iotex.Contract
	witnessListContractOnEthereum common.Address
	witnessListContractOnIoTeX    iotex.Contract
	lastUpdateTime                time.Time
	witnessesOnEthereum           map[common.Address]bool
	witnessesOnIoTeX              map[address.Address]bool
	erc20ToXrc20                  map[common.Address]address.Address
	xrc20ToErc20                  map[address.Address]common.Address
}

// NewAuth creates a new auth
func NewAuth(
	ethClientPool *EthClientPool,
	iotexClient iotex.AuthedClient,
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
		ethClientPool:                 ethClientPool,
		iotexClient:                   iotexClient,
		erc20TokenListContract:        erc20TokenListContract,
		witnessListContractOnEthereum: witnessListContractOnEthereum,
		xrc20TokenListContract:        iotexClient.Contract(xrc20TokenListContract, addressListABI),
		witnessListContractOnIoTeX:    iotexClient.Contract(witnessListContractOnIoTeX, addressListABI),
	}, nil
}

// IoTeXClient returns the iotex client
func (auth *Auth) IoTeXClient() iotex.AuthedClient {
	return auth.iotexClient
}

// EthereumClientPool returns the ethereum client pool
func (auth *Auth) EthereumClientPool() *EthClientPool {
	return auth.ethClientPool
}

func (auth *Auth) loadAddressListOnEthereum(contractAddr common.Address) ([]common.Address, error) {
	var retval []common.Address
	if err := auth.ethClientPool.Execute(func(client *ethclient.Client) error {
		list, err := contract.NewAddressListCaller(contractAddr, client)
		if err != nil {
			return errors.Wrapf(err, "failed to create caller")
		}
		count, err := list.Count(&bind.CallOpts{})
		offset := big.NewInt(0)
		limit := uint8(10)
		retval = []common.Address{}
		for offset.Cmp(count) < 0 {
			result, err := list.GetActiveItems(&bind.CallOpts{}, offset, limit)
			if err != nil {
				return errors.Wrap(err, "failed to query list")
			}
			retval = append(retval, result.Items[0:result.Count.Int64()]...)
			offset.Add(offset, big.NewInt(int64(limit)))
		}
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to load address list")
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
	newXrc20ToErc20 := map[address.Address]common.Address{}
	newErc20ToXrc20 := map[common.Address]address.Address{}
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
		if ethAddr, ok := auth.xrc20ToErc20[ioAddr]; ok {
			newXrc20ToErc20[ioAddr] = ethAddr
			newErc20ToXrc20[ethAddr] = ioAddr
		} else {
			response, err := auth.iotexClient.ReadOnlyContract(ioAddr, shadowTokenABI).Read("coToken").Call(context.Background())
			if err != nil {
				return errors.Wrapf(err, "failed to get corresponding token of %s", ioAddr)
			}
			tokenOnEth := common.Address{}
			if err := response.Unmarshal(&tokenOnEth); err != nil {
				return errors.Wrapf(err, "failed to extract corresponding token of %s from %s", ioAddr, response)
			}
			newXrc20ToErc20[ioAddr] = tokenOnEth
			if erc20Addr, ok := newErc20ToXrc20[tokenOnEth]; ok {
				return errors.Wrapf(err, "two Xrc20 tokens %s and %s map to the same Erc20 %s", ioAddr, erc20Addr, tokenOnEth)
			}
			newErc20ToXrc20[tokenOnEth] = ioAddr
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
		if _, ok := newErc20ToXrc20[token]; !ok {
			return errors.Errorf("erc20 token %s doesn't have a match token", token)
		}
	}
	auth.mu.Lock()
	defer auth.mu.Unlock()
	auth.lastUpdateTime = time.Now()
	for _, w := range witnessesOnIoTeX {
		auth.witnessesOnIoTeX[w] = true
	}
	for _, w := range witnessesOnEth {
		auth.witnessesOnEthereum[w] = true
	}
	auth.erc20ToXrc20 = newErc20ToXrc20
	auth.xrc20ToErc20 = newXrc20ToErc20

	return nil
}

// IsActiveWitnessOnIoTeX returns true if the input address is an active witness on IoTeX
func (auth *Auth) IsActiveWitnessOnIoTeX(witness address.Address) bool {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.witnessesOnIoTeX[witness]
}

// IsActiveWitnessOnEthereum returns true if the input address is an active witness on Ethereum
func (auth *Auth) IsActiveWitnessOnEthereum(witness common.Address) bool {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	return auth.witnessesOnEthereum[witness]
}

// CorrespondingXrc20Token returns the corresponding Xrc20 token address on IoTeX
func (auth *Auth) CorrespondingXrc20Token(erc20 common.Address) (address.Address, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	xrc20, ok := auth.erc20ToXrc20[erc20]
	if !ok {
		return nil, errors.Errorf("cannot find corresponding XRC20 token address for %s", erc20)
	}

	return address.FromBytes(xrc20.Bytes())
}

// CorrespondingErc20Token returns the corresponding Erc20 token address on Ethereum
func (auth *Auth) CorrespondingErc20Token(xrc20 address.Address) (common.Address, error) {
	auth.mu.RLock()
	defer auth.mu.RUnlock()

	erc20, ok := auth.xrc20ToErc20[xrc20]
	if !ok {
		return common.Address{}, errors.Errorf("cannot find corresponding ERC20 token address for %s", erc20)
	}

	return common.BytesToAddress(erc20.Bytes()), nil
}
