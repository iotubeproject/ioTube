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
	"net"
	"os"
	"strconv"
	"time"

	uconfig "go.uber.org/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness-relayer/relayer"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain                 string        `json:"chain" yaml:"chain"`
	RefreshInterval       time.Duration `json:"refreshInterval" yaml:"refreshInterval"`
	ClientURL             string        `json:"clientURL" yaml:"clientURL"`
	EthConfirmBlockNumber uint8         `json:"ethConfirmBlockNumber" yaml:"ethConfirmBlockNumber"`
	EthGasPriceLimit      uint64        `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
	Port                  int           `json:"port" yaml:"port"`
	PrivateKey            string        `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string        `json:"slackWebHook" yaml:"slackWebHook"`
	ValidatorAddress      string        `json:"vialidatorAddress" yaml:"validatorAddress"`
	DatabaseURL           string        `json:"databaseURL" yaml:"databaseURL"`
	TransferTableName     string        `json:"transferTableName" yaml:"transferTableName"`
	WitnessTableName      string        `json:"witnessTableName" yaml:"witnessTableName"`
}

var defaultConfig = Configuration{
	Chain:                 "iotex",
	RefreshInterval:       time.Hour,
	ClientURL:             "",
	EthConfirmBlockNumber: 20,
	EthGasPriceLimit:      120000000000,
	Port:                  8080,
	PrivateKey:            "",
	SlackWebHook:          "",
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
	if port, ok := os.LookupEnv("RELAYER_PORT"); ok {
		cfg.Port, err = strconv.Atoi(port)
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
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to decode private key %v", err)
	}
	if validatorAddr, ok := os.LookupEnv("RELAYER_VALIDATOR_ADDRESS"); ok {
		cfg.ValidatorAddress = validatorAddr
	}
	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	log.Printf("Listening to port %d\n", cfg.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen to port: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	log.Println("Creating service")
	var transferValidator relayer.TransferValidator
	if chain, ok := os.LookupEnv("WITNESS_CHAIN"); ok {
		cfg.Chain = chain
	}
	switch cfg.Chain {
	case "ethereum":
		ethClient, err := ethclient.Dial(cfg.ClientURL)
		if err != nil {
			log.Fatalf("failed to create eth client %v\n", err)
		}
		if transferValidator, err = relayer.NewTransferValidatorOnEthereum(
			ethClient,
			privateKey,
			cfg.EthConfirmBlockNumber,
			new(big.Int).SetUint64(cfg.EthGasPriceLimit),
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
			privateKey,
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
			db.NewStore("mysql", cfg.DatabaseURL),
			cfg.TransferTableName,
			cfg.WitnessTableName,
		),
	)
	if err != nil {
		log.Fatalf("failed to create relay service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start relay service: %v\n", err)
	}
	defer service.Stop(context.Background())
	services.RegisterRelayServiceServer(grpcServer, service)
	log.Println("Registering...")
	reflection.Register(grpcServer)
	log.Println("Serving...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
