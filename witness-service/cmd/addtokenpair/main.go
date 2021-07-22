// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"go.uber.org/config"
)

type Chain struct {
	URL                     string `json:"url" yaml:"url"`
	StandardTokenListAddr   string `json:"standardTokenListAddr" yaml:"standardTokenListAddr"`
	IoTeXProxyTokenListAddr string `json:"iotexProxyTokenListAddr" yaml:"iotexProxyTokenListAddr"`
	MinterAddr              string `json:"minterAddr" yaml:"minterAddr"`
	OperatorPrivateKey      string `json:"operatorPrivateKey" yaml:"operatorPrivateKey"`
}

type TokenConfig struct {
	MinAmount          string `json:"minAmount" yaml:"minAmount"`
	MaxAmount          string `json:"maxAmount" yaml:"maxAmount"`
	SourceTokenAddr    string `json:"sourceTokenAddr" yaml:"sourceTokenAddr"`
	SourceChain        string `json:"sourceChain" yaml:"sourceChain"`
	ShadowTokenAddr    string `json:"shadowTokenAddr" yaml:"shadowTokenAddr"`
	ShadowTokenName    string `json:"shadowTokenName" yaml:"shadowTokenName"`
	ShadowTokenSymbol  string `json:"shadowTokenSymbol" yaml:"shadowTokenSymbol"`
	ShadowTokenDecimal int    `json:"shadowTokenDecimal" yaml:"shadowTokenDecimal"`
}

type Config struct {
	IoTeX  Chain            `json:"iotex" yaml:"iotex"`
	Chains map[string]Chain `json:"chains" yaml:"chains"`
	Token  TokenConfig      `json:"token" yaml:"token"`
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
	sourceChain, ok := cfg.Chains[cfg.Token.SourceChain]
	if !ok {
		log.Fatalf("source chain %s is not defined in chain configs", cfg.Token.SourceChain)
	}
	minAmount, ok := new(big.Int).SetString(cfg.Token.MinAmount, 10)
	if !ok {
		log.Fatalln("failed to parse token min amount", cfg.Token.MinAmount)
	}
	maxAmount, ok := new(big.Int).SetString(cfg.Token.MaxAmount, 10)
	if !ok {
		log.Fatalln("failed to parse token max amount", cfg.Token.MaxAmount)
	}
	srcChainClient, err := ethclient.Dial(sourceChain.URL)
	if err != nil {
		log.Fatalln(err)
	}
	srcChainID, err := srcChainClient.ChainID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	srcPrivateKey, err := crypto.HexToECDSA(sourceChain.OperatorPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	srcChainAuth, err := bind.NewKeyedTransactorWithChainID(srcPrivateKey, srcChainID)
	if err != nil {
		log.Fatal(err)
	}
	srcTokenAddr := common.HexToAddress(cfg.Token.SourceTokenAddr)
	if err := addTokenToList(
		srcTokenAddr,
		common.HexToAddress(sourceChain.StandardTokenListAddr),
		minAmount,
		maxAmount,
		srcChainAuth,
		srcChainClient,
	); err != nil {
		log.Fatalln(err)
	}
	iotexChain := cfg.IoTeX
	targetChainClient, err := ethclient.Dial(iotexChain.URL)
	if err != nil {
		log.Fatalln(err)
	}
	targetChainID, err := targetChainClient.ChainID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	targetChainPrivateKey, err := crypto.HexToECDSA(iotexChain.OperatorPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	targetChainAuth, err := bind.NewKeyedTransactorWithChainID(targetChainPrivateKey, targetChainID)
	if err != nil {
		log.Fatal(err)
	}
	var shadowTokenAddr common.Address
	if cfg.Token.ShadowTokenAddr == "" {
		shadowTokenAddr, _, _, err = contract.DeployShadowToken(
			targetChainAuth,
			targetChainClient,
			common.HexToAddress(iotexChain.MinterAddr),
			srcTokenAddr,
			cfg.Token.ShadowTokenName,
			cfg.Token.ShadowTokenSymbol,
			uint8(cfg.Token.ShadowTokenDecimal),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		shadowTokenAddr = common.HexToAddress(cfg.Token.ShadowTokenAddr)
	}
	if err := addTokenToList(
		shadowTokenAddr,
		common.HexToAddress(sourceChain.IoTeXProxyTokenListAddr),
		minAmount,
		maxAmount,
		targetChainAuth,
		targetChainClient,
	); err != nil {
		log.Fatal(err)
	}
}

func addTokenToList(
	tokenAddr, tokenListAddr common.Address,
	minAmount, maxAmount *big.Int,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
) error {
	tokenList, err := contract.NewTokenList(tokenListAddr, backend)
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
	}
	return nil
}
