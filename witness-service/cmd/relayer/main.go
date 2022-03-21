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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/relayer"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain                 string        `json:"chain" yaml:"chain"`
	ClientURL             string        `json:"clientURL" yaml:"clientURL"`
	EthConfirmBlockNumber uint8         `json:"ethConfirmBlockNumber" yaml:"ethConfirmBlockNumber"`
	EthGasPriceLimit      uint64        `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
	EthGasPriceDeviation  int64         `json:"ethGasPriceDeviation" yaml:"ethGasPriceDeviation"`
	EthGasPriceGap        uint64        `json:"ethGasPriceGap" yaml:"ethGasPriceGap"`
	PrivateKey            string        `json:"privateKey" yaml:"privateKey"`
	Interval              time.Duration `json:"interval" yaml:"interval"`
	ValidatorAddress      string        `json:"vialidatorAddress" yaml:"validatorAddress"`

	SlackWebHook      string    `json:"slackWebHook" yaml:"slackWebHook"`
	LarkWebHook       string    `json:"larkWebHook" yaml:"larkWebHook"`
	GrpcPort          int       `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort     int       `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	Database          db.Config `json:"database" yaml:"database"`
	TransferTableName string    `json:"transferTableName" yaml:"transferTableName"`
	WitnessTableName  string    `json:"witnessTableName" yaml:"witnessTableName"`
}

var defaultConfig = Configuration{
	Chain:                 "iotex",
	Interval:              time.Hour,
	ClientURL:             "",
	EthConfirmBlockNumber: 20,
	EthGasPriceLimit:      120000000000,
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
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	flag.Parse()
	opts := []config.YAMLOption{config.Static(defaultConfig), config.Expand(os.LookupEnv)}
	if *configFile != "" {
		opts = append(opts, config.File(*configFile))
	}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
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
	if validatorAddr, ok := os.LookupEnv("RELAYER_VALIDATOR_ADDRESS"); ok {
		cfg.ValidatorAddress = validatorAddr
	}
	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	if cfg.LarkWebHook != "" {
		util.SetLarkURL(cfg.LarkWebHook)
	}
	util.SetPrefix("relayer-" + cfg.Chain)

	log.Println("Creating service")
	var transferValidator relayer.TransferValidator
	if chain, ok := os.LookupEnv("RELAYER_CHAIN"); ok {
		cfg.Chain = chain
	}
	switch cfg.Chain {
	case "heco", "bsc", "matic", "polis":
		// heco and bsc are idential to ethereum
		fallthrough
	case "ethereum":
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
		if transferValidator, err = relayer.NewTransferValidatorOnEthereum(
			ethClient,
			privateKeys,
			cfg.EthConfirmBlockNumber,
			new(big.Int).SetUint64(cfg.EthGasPriceLimit),
			new(big.Int).SetInt64(cfg.EthGasPriceDeviation),
			new(big.Int).SetUint64(cfg.EthGasPriceGap),
			common.HexToAddress(cfg.ValidatorAddress),
		); err != nil {
			log.Fatalf("failed to create transfer validator: %v\n", err)
		}
	case "iotex":
		conn, err := iotex.NewDefaultGRPCConn(cfg.ClientURL)
		if err != nil {
			log.Fatal(err)
		}
		// defer conn.Close()
		acc, err := account.HexStringToAccount(cfg.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		validatorContractAddr, err := address.FromString(cfg.ValidatorAddress)
		if err != nil {
			log.Fatalf("failed to parse validator contract address %s\n", cfg.ValidatorAddress)
		}
		if transferValidator, err = relayer.NewTransferValidatorOnIoTeX(
			iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc),
			validatorContractAddr,
		); err != nil {
			log.Fatalf("failed to create transfer validator: %v\n", err)
		}
	default:
		log.Fatalf("unknown chain name '%s'\n", cfg.Chain)
	}
	service, err := relayer.NewService(
		transferValidator,
		relayer.NewRecorder(
			db.NewStore(cfg.Database),
			cfg.TransferTableName,
			cfg.WitnessTableName,
		),
		cfg.Interval,
	)
	if err != nil {
		log.Fatalf("failed to create relay service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start relay service: %v\n", err)
	}
	defer service.Stop(context.Background())

	relayer.StartServer(service, cfg.GrpcPort, cfg.GrpcProxyPort)

	select {}
}
