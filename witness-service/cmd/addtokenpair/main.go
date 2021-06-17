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

type SourceChain struct {
	URL                string `json:"url" yaml:"url"`
	TokenListAddr      string `json:"tokenListAddr" yaml:"tokenListAddr"`
	TokenAddr          string `json:"tokenAddr" yaml:"tokenAddr"`
	OperatorPrivateKey string `json:"operatorPrivateKey" yaml:"operatorPrivateKey"`
}

type TargetChain struct {
	URL                string `json:"url" yaml:"url"`
	TokenListAddr      string `json:"tokenListAddr" yaml:"tokenListAddr"`
	MinterAddr         string `json:"minterAddr" yaml:"minterAddr"`
	TokenAddr          string `json:"tokenAddr" yaml:"tokenAddr"`
	TokenName          string `json:"tokenName" yaml:"tokenName"`
	TokenSymbol        string `json:"tokenSymbol" yaml:"tokenSymbol"`
	TokenDecimal       int    `json:"tokenDecimal" yaml:"tokenDecimal"`
	OperatorPrivateKey string `json:"operatorPrivateKey" yaml:"operatorPrivateKey"`
}

// Configuration defines the configuration of the witness service
type Configuration struct {
	MinAmount string      `json:"minAmount" yaml:"minAmount"`
	MaxAmount string      `json:"maxAmount" yaml:"maxAmount"`
	Source    SourceChain `json:"source" yaml:"source"`
	Target    TargetChain `json:"target" yaml:"target"`
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
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	minAmount, ok := new(big.Int).SetString(cfg.MinAmount, 10)
	if !ok {
		log.Fatalln("failed to parse token min amount", cfg.MinAmount)
	}
	maxAmount, ok := new(big.Int).SetString(cfg.MinAmount, 10)
	if !ok {
		log.Fatalln("failed to parse token max amount", cfg.MaxAmount)
	}
	srcChainClient, err := ethclient.Dial(cfg.Source.URL)
	if err != nil {
		log.Fatalln(err)
	}
	srcChainID, err := srcChainClient.ChainID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	srcPrivateKey, err := crypto.HexToECDSA(cfg.Source.OperatorPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	srcChainAuth, err := bind.NewKeyedTransactorWithChainID(srcPrivateKey, srcChainID)
	if err != nil {
		log.Fatal(err)
	}
	srcTokenAddr := common.HexToAddress(cfg.Source.TokenAddr)
	if err := addTokenToList(
		srcTokenAddr,
		common.HexToAddress(cfg.Source.TokenListAddr),
		minAmount,
		maxAmount,
		srcChainAuth,
		srcChainClient,
	); err != nil {
		log.Fatalln(err)
	}
	targetChainClient, err := ethclient.Dial(cfg.Target.URL)
	if err != nil {
		log.Fatalln(err)
	}
	targetChainID, err := targetChainClient.ChainID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	targetChainPrivateKey, err := crypto.HexToECDSA(cfg.Target.OperatorPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	targetChainAuth, err := bind.NewKeyedTransactorWithChainID(targetChainPrivateKey, targetChainID)
	if err != nil {
		log.Fatal(err)
	}
	var shadowTokenAddr common.Address
	if cfg.Target.TokenAddr == "" {
		shadowTokenAddr, _, _, err = contract.DeployShadowToken(
			targetChainAuth,
			targetChainClient,
			common.HexToAddress(cfg.Target.MinterAddr),
			srcTokenAddr,
			cfg.Target.TokenName,
			cfg.Target.TokenSymbol,
			uint8(cfg.Target.TokenDecimal),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		shadowTokenAddr = common.HexToAddress(cfg.Target.TokenAddr)
	}
	if err := addTokenToList(
		shadowTokenAddr,
		common.HexToAddress(cfg.Target.TokenListAddr),
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
