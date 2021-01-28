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
	"os"
	"strconv"
	"time"

	uconfig "go.uber.org/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness-relayer/witness"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	IoTeXApiURL              string        `json:"iotexApiURL" yaml:"iotexApiURL"`
	RelayerURL               string        `json:"relayerURL" yaml:"relayerURL"`
	PrivateKey               string        `json:"privateKey" yaml:"privateKey"`
	SlackWebHook             string        `json:"slackWebHook" yaml:"slackWebHook"`
	ValidatorContractAddress string        `json:"vialidatorContractAddress" yaml:"validatorContractAddress"`
	CashierContractAddress   string        `json:"cashierContractAddress" yaml:"cashierContractAddress"`
	StartBlockHeight         int           `json:"startBlockHeight" yaml:"startBlockHeight"`
	BatchSize                int           `json:"batchSize" yaml:"batchSize"`
	ProcessInterval          time.Duration `json:"processInterval" yaml:"processInterval"`
	DatabaseURL              string        `json:"databaseURL" yaml:"databaseURL"`
	TransferTableName        string        `json:"transferTableName" yaml:"transferTableName"`
}

var defaultConfig = Configuration{
	ProcessInterval:          time.Minute,
	RelayerURL:               "",
	StartBlockHeight:         9305000,
	BatchSize:                100,
	PrivateKey:               "",
	SlackWebHook:             "",
	IoTeXApiURL:              "api.iotex.one:443",
	TransferTableName:        "transfers",
	ValidatorContractAddress: "",
	CashierContractAddress:   "",
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
	opts := make([]uconfig.YAMLOption, 0)
	opts = append(opts, uconfig.Static(defaultConfig))
	opts = append(opts, uconfig.Expand(os.LookupEnv))
	if *configFile != "" {
		opts = append(opts, uconfig.File(*configFile))
	}
	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Configuration
	if err := yaml.Get(uconfig.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	if height, ok := os.LookupEnv("WITNESS_START_BLOCK_HEIGHT"); ok {
		cfg.StartBlockHeight, err = strconv.Atoi(height)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if relayerURL, ok := os.LookupEnv("WITNESS_RELAYER_URL"); ok {
		cfg.RelayerURL = relayerURL
	}
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok {
		cfg.PrivateKey = pk
	}
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to decode private key %v", err)
	}
	if validatorAddr, ok := os.LookupEnv("WITNESS_VALIDATOR_CONTRACT_ADDRESS"); ok {
		cfg.ValidatorContractAddress = validatorAddr
	}
	if cashierAddr, ok := os.LookupEnv("WITNESS_CASHIER_CONTRACT_ADDRESS"); ok {
		cfg.CashierContractAddress = cashierAddr
	}
	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	if url, ok := os.LookupEnv("WITNESS_IOTEX_API_URL"); ok {
		cfg.IoTeXApiURL = url
	}
	conn, err := iotex.NewDefaultGRPCConn(cfg.IoTeXApiURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	service, err := witness.NewService(
		cfg.RelayerURL,
		common.HexToAddress(cfg.ValidatorContractAddress),
		witness.NewTokenCashier(client),
		witness.NewRecorder(
			db.NewStore("mysql", cfg.DatabaseURL),
			cfg.TransferTableName,
		),
		privateKey,
		uint64(cfg.StartBlockHeight),
		uint16(cfg.BatchSize),
		cfg.ProcessInterval,
	)
	if err != nil {
		log.Fatalf("failed to create relay service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start relay service: %v\n", err)
	}
	defer service.Stop(context.Background())
	log.Println("Serving...")
	select {}
}
