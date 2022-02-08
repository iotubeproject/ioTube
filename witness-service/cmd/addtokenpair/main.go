// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/pkg/errors"
	"go.uber.org/config"
)

var zeroAddr = common.Address{}

type ChainPair struct {
	ChainA Chain `json:"chainA" yaml:"chainA"`
	ChainB Chain `json:"chainB" yaml:"chainB"`
}

type Chain struct {
	URL                   string `json:"url" yaml:"url"`
	MinterPoolAddr        string `json:"minterPoolAddr" yaml:"minterPoolAddr"`
	OperatorPrivateKey    string `json:"operatorPrivateKey" yaml:"operatorPrivateKey"`
	CreatorPrivateKey     string `json:"creatorPrivateKey" yaml:"creatorPrivateKey"`
	StandardTokenListAddr string `json:"standardTokenListAddr" yaml:"standardTokenListAddr"`
	ProxyTokenListAddr    string `json:"proxyTokenListAddr" yaml:"proxyTokenListAddr"`
	RouterAddr            string `json:"routerAddr" yaml:"routerAddr"`
}

type TokenConfig struct {
	IsProxy  bool   `json:"isProxy" yaml:"isProxy"`
	Address  string `json:"address" yaml:"address"`
	Name     string `json:"name" yaml:"name"`
	Symbol   string `json:"symbol" yaml:"symbol"`
	Decimals uint8  `json:"decimals" yaml:"decimals"`
}

type Config struct {
	ChainPairs            map[string]ChainPair `json:"chainPairs" yaml:"chainPairs"`
	PairName              string               `json:"pairName" yaml:"pairName"`
	OriginTokenAddr       string               `json:"originTokenAddr" yaml:"originTokenAddr"`
	IsOriginTokenOnChainB bool                 `json:"isOriginTokenOnChainB" yaml:"isOriginTokenOnChainB"`
	MinAmount             string               `json:"minAmount" yaml:"minAmount"`
	MaxAmount             string               `json:"maxAmount" yaml:"maxAmount"`
	TokenOnChainA         TokenConfig          `json:"tokenOnChainA" yaml:"tokenOnChainA"`
	TokenOnChainB         TokenConfig          `json:"tokenOnChainB" yaml:"tokenOnChainB"`
}

var configFile = flag.String("config", "", "path of config file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename>")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *configFile == "" {
		log.Fatalln("config file is not specified")
	}
	opts := []config.YAMLOption{config.Expand(os.LookupEnv), config.File(*configFile)}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Config
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	pair, ok := cfg.ChainPairs[cfg.PairName]
	if !ok {
		err = errors.Errorf("chain pair %s is not defined in config", cfg.PairName)
		return
	}
	minAmount, ok := new(big.Int).SetString(cfg.MinAmount, 10)
	if !ok {
		err = errors.Errorf("failed to parse token min amount %s", cfg.MinAmount)
		return
	}
	maxAmount, ok := new(big.Int).SetString(cfg.MaxAmount, 10)
	if !ok {
		err = errors.Errorf("failed to parse token max amount %s", cfg.MaxAmount)
		return
	}
	if cfg.IsOriginTokenOnChainB {
		err = addTokens(
			common.HexToAddress(cfg.OriginTokenAddr),
			pair.ChainB,
			pair.ChainA,
			cfg.TokenOnChainB,
			cfg.TokenOnChainA,
			minAmount,
			maxAmount,
		)
	} else {
		err = addTokens(
			common.HexToAddress(cfg.OriginTokenAddr),
			pair.ChainA,
			pair.ChainB,
			cfg.TokenOnChainA,
			cfg.TokenOnChainB,
			minAmount,
			maxAmount,
		)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func addTokens(
	originTokenAddr common.Address,
	nativeChain, foreignChain Chain,
	tokenOnNativeChain, tokenOnForeignChain TokenConfig,
	minAmount, maxAmount *big.Int,
) error {
	operatorPrivateKey, err := crypto.HexToECDSA(nativeChain.OperatorPrivateKey)
	if err != nil {
		return err
	}
	creatorPrivateKey, err := crypto.HexToECDSA(nativeChain.CreatorPrivateKey)
	if err != nil {
		return err
	}
	name, symbol, decimals, err := addToken(
		nativeChain.URL,
		operatorPrivateKey,
		creatorPrivateKey,
		common.HexToAddress(nativeChain.MinterPoolAddr),
		common.HexToAddress(nativeChain.StandardTokenListAddr),
		common.HexToAddress(nativeChain.ProxyTokenListAddr),
		common.HexToAddress(nativeChain.RouterAddr),
		originTokenAddr,
		minAmount,
		maxAmount,
		tokenOnNativeChain,
	)

	if err != nil {
		return errors.Wrap(err, "failed to add token to source chain")
	}
	if tokenOnForeignChain.Name == "" {
		tokenOnForeignChain.Name = name
	}
	if tokenOnForeignChain.Symbol == "" {
		tokenOnForeignChain.Symbol = symbol
	}
	if tokenOnForeignChain.Decimals == 0 {
		tokenOnForeignChain.Decimals = decimals
	}
	foreignChainOperatorPrivateKey, err := crypto.HexToECDSA(foreignChain.OperatorPrivateKey)
	if err != nil {
		return err
	}
	foreignChainCreatorPrivateKey, err := crypto.HexToECDSA(foreignChain.OperatorPrivateKey)
	if err != nil {
		return err
	}
	_, _, _, err = addToken(
		foreignChain.URL,
		foreignChainOperatorPrivateKey,
		foreignChainCreatorPrivateKey,
		common.HexToAddress(foreignChain.MinterPoolAddr),
		common.HexToAddress(foreignChain.StandardTokenListAddr),
		common.HexToAddress(foreignChain.ProxyTokenListAddr),
		common.HexToAddress(foreignChain.RouterAddr),
		zeroAddr,
		minAmount,
		maxAmount,
		tokenOnForeignChain,
	)
	return err
}

func addToken(
	url string,
	operatorPrivateKey, creatorPrivateKey *ecdsa.PrivateKey,
	minterPool, standardTokenList, proxyTokenList, router, originTokenAddr common.Address,
	minAmount, maxAmount *big.Int,
	tokenConfig TokenConfig,
) (
	name, symbol string, decimals uint8, err error,
) {
	chainClient, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	chainID, err := chainClient.ChainID(context.Background())
	if err != nil {
		return
	}
	operatorAuth, err := bind.NewKeyedTransactorWithChainID(operatorPrivateKey, chainID)
	if err != nil {
		return
	}
	if originTokenAddr != zeroAddr {
		// using CrosschainERC20 as ERC20 token
		var originToken *contract.CrosschainERC20
		originToken, err = contract.NewCrosschainERC20(originTokenAddr, chainClient)
		if err != nil {
			return
		}
		name, err = originToken.Name(nil)
		if err != nil {
			return
		}
		symbol, err = originToken.Symbol(nil)
		if err != nil {
			return
		}
		decimals, err = originToken.Decimals(nil)
		if err != nil {
			return
		}
	}
	if tokenConfig.Name == "" {
		if name == "" {
			err = errors.New("invalid token name")
			return
		}
		tokenConfig.Name = "Crosschain " + name
	}
	if tokenConfig.Symbol == "" {
		if symbol == "" {
			err = errors.New("invalid token symbol")
			return
		}
		tokenConfig.Symbol = "C" + symbol
	}
	if tokenConfig.Decimals == 0 {
		tokenConfig.Decimals = decimals
	}
	var tokenListAddr common.Address
	var tokenAddr common.Address
	if tokenConfig.IsProxy {
		if tokenConfig.Address != "" {
			tokenAddr = common.HexToAddress(tokenConfig.Address)
		} else {
			var creatorAuth *bind.TransactOpts
			creatorAuth, err = bind.NewKeyedTransactorWithChainID(creatorPrivateKey, chainID)
			if err != nil {
				return
			}
			log.Printf("Deploying crosschain token (%s, %s, %s, %s, %d)\n", originTokenAddr.String(), minterPool.String(), tokenConfig.Name, tokenConfig.Symbol, tokenConfig.Decimals)
			var tx *types.Transaction
			tokenAddr, tx, _, err = contract.DeployCrosschainERC20(
				creatorAuth,
				chainClient,
				originTokenAddr,
				minterPool,
				tokenConfig.Name,
				tokenConfig.Symbol,
				tokenConfig.Decimals,
			)
			if err != nil {
				return
			}
			log.Printf("Waiting token %s deployment for %s with tx %s\n", tokenAddr.String(), originTokenAddr, tx.Hash().Hex())
			waitUntilConfirm(chainClient, tx)
		}
		if router != zeroAddr {
			var ct *contract.CrosschainERC20
			ct, err = contract.NewCrosschainERC20(tokenAddr, chainClient)
			if err != nil {
				return
			}
			var coAddr common.Address
			coAddr, err = ct.CoToken(nil)
			if err != nil {
				return
			}
			if coAddr != zeroAddr {
				var ctcr *contract.CrosschainTokenCashierRouter
				ctcr, err = contract.NewCrosschainTokenCashierRouter(router, chainClient)
				if err != nil {
					return
				}
				var cashier common.Address
				cashier, err = ctcr.Cashier(nil)
				if err != nil {
					return
				}
				var allowance *big.Int
				allowance, err = ct.Allowance(nil, router, cashier)
				if err != nil {
					return
				}
				if allowance.Sign() == 0 {
					var tx *types.Transaction
					tx, err = ctcr.ApproveCrosschainToken(operatorAuth, tokenAddr)
					if err != nil {
						return
					}
					log.Printf("Adding %s to router via tx %s\n", tokenAddr, tx.Hash().Hex())
					waitUntilConfirm(chainClient, tx)
				}
			}
		}
		tokenListAddr = proxyTokenList
	} else {
		tokenAddr = originTokenAddr
		tokenListAddr = standardTokenList
	}
	log.Printf(">> C-Token address on %d: %s\n", chainID, tokenAddr)
	err = addTokenToList(
		tokenAddr,
		tokenListAddr,
		minAmount,
		maxAmount,
		operatorAuth,
		chainClient,
	)

	return
}

func addTokenToList(
	tokenAddr, tokenListAddr common.Address,
	minAmount, maxAmount *big.Int,
	auth *bind.TransactOpts,
	client *ethclient.Client,
) error {
	if tokenAddr == zeroAddr {
		return errors.New("invalid token address")
	}
	tokenList, err := contract.NewTokenList(tokenListAddr, client)
	if err != nil {
		return err
	}
	active, err := tokenList.IsActive(nil, tokenAddr)
	if err != nil {
		return err
	}
	if !active {
		tx, err := tokenList.AddToken(auth, tokenAddr, minAmount, maxAmount)
		if err != nil {
			return err
		}
		log.Printf("Adding token %x to token list %x via tx %s", tokenAddr, tokenListAddr, tx.Hash())
		waitUntilConfirm(client, tx)
	}
	return nil
}

func waitUntilConfirm(client *ethclient.Client, tx *types.Transaction) error {
	h := tx.Hash()
	for {
		log.Println("Wait for 5s...")
		time.Sleep(5 * time.Second)
		receipt, err := client.TransactionReceipt(context.Background(), h)
		switch errors.Cause(err) {
		case ethereum.NotFound:
			// do nothing
		case nil:
			if receipt.Status != types.ReceiptStatusSuccessful {
				return errors.Errorf("transaction %x is rejected", h)
			}
			for {
				tip, err := client.BlockNumber(context.Background())
				if err != nil {
					return err
				}
				if tip > receipt.BlockNumber.Uint64()+20 {
					return nil
				}
			}
		default:
			return err
		}
	}
}
