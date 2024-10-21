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
	"strconv"
	"strings"
	"time"

	soltypes "github.com/blocto/solana-go-sdk/types"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/relayer"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// ValidatorConfig defines the configuration of a validator
	ValidatorConfig struct {
		Address     string   `json:"address" yaml:"address"`
		Cashiers    []string `json:"cashiers" yaml:"cashiers"`
		WithPayload bool     `json:"withPayload" yaml:"withPayload"`
		FromSolana  bool     `json:"fromSolana" yaml:"fromSolana"`
	}
	// Configuration defines the configuration of the witness service
	Configuration struct {
		Chain                 string            `json:"chain" yaml:"chain"`
		ClientURL             string            `json:"clientURL" yaml:"clientURL"`
		EthConfirmBlockNumber uint16            `json:"ethConfirmBlockNumber" yaml:"ethConfirmBlockNumber"`
		EthDefaultGasPrice    uint64            `json:"ethDefaultGasPrice" yaml:"ethDefaultGasPrice"`
		EthGasPriceLimit      uint64            `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
		EthGasPriceHardLimit  uint64            `json:"ethGasPriceHardLimit" yaml:"ethGasPriceHardLimit"`
		EthGasPriceDeviation  int64             `json:"ethGasPriceDeviation" yaml:"ethGasPriceDeviation"`
		EthGasPriceGap        uint64            `json:"ethGasPriceGap" yaml:"ethGasPriceGap"`
		PrivateKey            string            `json:"privateKey" yaml:"privateKey"`
		Interval              time.Duration     `json:"interval" yaml:"interval"`
		Validators            []ValidatorConfig `json:"validators" yaml:"validators"`

		BonusTokens map[string]*big.Int `json:"bonusTokens" yaml:"bonusTokens"`
		Bonus       *big.Int            `json:"bonus" yaml:"bonus"`

		AlwaysReset             bool      `json:"alwaysReset" yaml:"alwaysReset"`
		SlackWebHook            string    `json:"slackWebHook" yaml:"slackWebHook"`
		LarkWebHook             string    `json:"larkWebHook" yaml:"larkWebHook"`
		GrpcPort                int       `json:"grpcPort" yaml:"grpcPort"`
		GrpcProxyPort           int       `json:"grpcProxyPort" yaml:"grpcProxyPort"`
		Database                db.Config `json:"database" yaml:"database"`
		ExplorerDatabase        db.Config `json:"explorerDatabase" yaml:"explorerDatabase"`
		TransferTableName       string    `json:"transferTableName" yaml:"transferTableName"`
		NewTransactionTableName string    `json:"newTransactionTableName" yaml:"newTransactionTableName"`
		WitnessTableName        string    `json:"witnessTableName" yaml:"witnessTableName"`
		ExplorerTableName       string    `json:"explorerTableName" yaml:"explorerTableName"`

		SolanaConfig struct {
			ValidatorAddress        string  `json:"validatorAddress" yaml:"validatorAddress"`
			RealmAddr               string  `json:"realmAddr" yaml:"realmAddr"`
			GoverningTokenMintAddr  string  `json:"governingTokenMintAddr" yaml:"governingTokenMintAddr"`
			GovernanceAddr          string  `json:"governanceAddr" yaml:"governanceAddr"`
			ProposalAddr            string  `json:"proposalAddr" yaml:"proposalAddr"`
			ProposalTransactionAddr string  `json:"proposalTransactionAddr" yaml:"proposalTransactionAddr"`
			Threshold               float64 `json:"threshold" yaml:"threshold"`
			QPSLimit                uint32  `json:"qpsLimit" yaml:"qpsLimit"`
		} `json:"solanaConfig" yaml:"solanaConfig"`
		SourceChain string `json:"sourceChain" yaml:"sourceChain"`
	}
)

var defaultConfig = Configuration{
	Chain:                 "iotex-e",
	Interval:              time.Hour,
	ClientURL:             "",
	EthConfirmBlockNumber: 20,
	EthGasPriceLimit:      1200000000000,
	EthGasPriceDeviation:  0,
	EthGasPriceGap:        0,
	GrpcPort:              8080,
	GrpcProxyPort:         8081,
	PrivateKey:            "",
	SlackWebHook:          "",
	LarkWebHook:           "",
	TransferTableName:     "relayer.transfers",
	WitnessTableName:      "relayer.witnesses",
}

var configFile = flag.String("config", "", "path of config file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename>")
		flag.PrintDefaults()
	}
}

// main performs the main routine of the application:
//  1. parses the args;
//  2. analyzes the declaration of the API
//  3. sets the implementation of the handlers
//  4. listens on the port we want
func main() {
	flag.Parse()
	opts := []config.YAMLOption{config.Static(defaultConfig), config.Expand(os.LookupEnv)}
	if *configFile != "" {
		opts = append(opts, config.File(*configFile))
	}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to create yaml config"))
	}
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	if port, ok := os.LookupEnv("RELAYER_GRPC_PORT"); ok {
		cfg.GrpcPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if port, ok := os.LookupEnv("RELAYER_GRPC_PROXY_PORT"); ok {
		cfg.GrpcProxyPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if client, ok := os.LookupEnv("RELAYER_CLIENT_URL"); ok {
		cfg.ClientURL = client
	}
	if pk, ok := os.LookupEnv("RELAYER_PRIVATE_KEY"); ok {
		cfg.PrivateKey = pk
	}
	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	if cfg.LarkWebHook != "" {
		util.SetLarkURL(cfg.LarkWebHook)
	}
	util.SetPrefix("relayer-" + cfg.Chain)

	storeFactory := db.NewSQLStoreFactory()
	log.Println("Creating service")
	var service services.RelayServiceServer
	if chain, ok := os.LookupEnv("RELAYER_CHAIN"); ok {
		cfg.Chain = chain
	}
	switch cfg.Chain {
	case "ethereum", "heco", "bsc", "matic", "polis", "iotex-e", "iotex", "iotex-testnet", "sepolia":
		if cfg.ClientURL == "" {
			break
		}
		privateKeys := []*ecdsa.PrivateKey{}
		for _, pk := range strings.Split(cfg.PrivateKey, ",") {
			privateKey, err := crypto.HexToECDSA(pk)
			if err != nil {
				log.Fatalf("failed to decode private key %v", err)
			}
			privateKeys = append(privateKeys, privateKey)
		}
		ethClient, err := ethclient.Dial(cfg.ClientURL)
		if err != nil {
			log.Fatalf("failed to create eth client %v\n", err)
		}
		validators := map[string]relayer.TransferValidator{}
		for _, vc := range cfg.Validators {
			validatorAddr, err := util.ParseEthAddress(vc.Address)
			if err != nil {
				log.Fatalf("failed to parse validator address %s: %+v", vc.Address, err)
			}
			version := relayer.NoPayload
			if vc.WithPayload {
				version = relayer.Payload
			}
			if vc.FromSolana {
				version = relayer.FromSolana
			}
			validator, err := relayer.NewTransferValidatorOnEthereum(
				ethClient,
				privateKeys,
				cfg.EthConfirmBlockNumber,
				new(big.Int).SetUint64(cfg.EthDefaultGasPrice),
				new(big.Int).SetUint64(cfg.EthGasPriceLimit),
				new(big.Int).SetUint64(cfg.EthGasPriceHardLimit),
				new(big.Int).SetInt64(cfg.EthGasPriceDeviation),
				new(big.Int).SetUint64(cfg.EthGasPriceGap),
				version,
				validatorAddr,
			)
			if err != nil {
				log.Fatalf("failed to create validator: %+v\n", err)
			}
			for _, cashier := range vc.Cashiers {
				cashierAddr, err := util.ParseAddress(cashier)
				if err != nil {
					log.Fatalf("failed to parse cashier address %s: %+v", cashier, err)
				}
				validators[cashierAddr.String()] = validator
			}
		}
		bonusSender, err := relayer.NewBonusSender(ethClient, privateKeys, cfg.BonusTokens, cfg.Bonus)
		if err != nil {
			log.Fatalf("failed to create bonus sender: %+v\n", err)
		}

		ethService, err := relayer.NewServiceOnEthereum(
			validators,
			bonusSender,
			relayer.NewRecorder(
				storeFactory.NewStore(cfg.Database),
				storeFactory.NewStore(cfg.ExplorerDatabase),
				cfg.TransferTableName,
				cfg.WitnessTableName,
				"",
				cfg.ExplorerTableName,
			),
			cfg.Interval,
		)
		if err != nil {
			log.Fatalf("failed to create relay service: %v\n", err)
		}
		if cfg.AlwaysReset {
			ethService.SetAlwaysRetry()
		}
		if err := ethService.Start(context.Background()); err != nil {
			log.Fatalf("failed to start solana relay service: %v\n", err)
		}
		defer ethService.Stop(context.Background())
		service = ethService
	case "solana":
		transferValidatorAddr, err := util.NewSOLAddressDecoder().DecodeString(cfg.SolanaConfig.ProposalAddr)
		if err != nil {
			log.Fatalf("failed to decode validator address %v", err)
		}

		solRecorder := relayer.NewSolRecorder(
			storeFactory.NewStore(cfg.Database),
			cfg.TransferTableName,
			cfg.WitnessTableName,
			util.NewETHAddressDecoder(),
			util.NewSOLAddressDecoder(),
		)

		privateKey, err := soltypes.AccountFromHex(cfg.PrivateKey)
		if err != nil {
			log.Fatalf("failed to decode private key %v", err)
		}
		solProcessor := relayer.NewSolProcessor(
			client.NewClient(cfg.ClientURL),
			cfg.Interval,
			&privateKey,
			relayer.VoteConfig{
				ProgramID:               cfg.SolanaConfig.ValidatorAddress,
				RealmAddr:               cfg.SolanaConfig.RealmAddr,
				GoverningTokenMintAddr:  cfg.SolanaConfig.GoverningTokenMintAddr,
				GovernanceAddr:          cfg.SolanaConfig.GovernanceAddr,
				ProposalAddr:            cfg.SolanaConfig.ProposalAddr,
				ProposalTransactionAddr: cfg.SolanaConfig.ProposalTransactionAddr,
				Threshold:               cfg.SolanaConfig.Threshold,
			},
			solRecorder,
			cfg.SolanaConfig.QPSLimit,
		)
		solanaService, err := relayer.NewServiceOnSolana(solRecorder, transferValidatorAddr)
		if err != nil {
			log.Fatalf("failed to create relay service: %v\n", err)
		}
		solanaService.SetProcessor(solProcessor)
		log.Fatalf("iotex chain is not supported anymore, please switch to iotex-e\n")
		if cfg.AlwaysReset {
			solanaService.SetAlwaysRetry()
		}
		if err := solanaService.Start(context.Background()); err != nil {
			log.Fatalf("failed to start solana relay service: %v\n", err)
		}
		defer solanaService.Stop(context.Background())
		service = solanaService
	default:
		log.Fatalf("unknown chain name '%s'\n", cfg.Chain)
	}

	relayer.StartServer(service, cfg.GrpcPort, cfg.GrpcProxyPort)

	select {}
}
