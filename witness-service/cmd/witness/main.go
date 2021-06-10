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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain            string        `json:"chain" yaml:"chain"`
	ClientURL        string        `json:"clientURL" yaml:"clientURL"`
	Database         db.Config     `json:"database" yaml:"database"`
	PrivateKey       string        `json:"privateKey" yaml:"privateKey"`
	SlackWebHook     string        `json:"slackWebHook" yaml:"slackWebHook"`
	BatchSize        int           `json:"batchSize" yaml:"batchSize"`
	Interval         time.Duration `json:"interval" yaml:"interval"`
	GrpcPort         int           `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort    int           `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	Cashiers         []struct {
		ID                       string `json:"id" yaml:"id"`
		RelayerURL               string `json:"relayerURL" yaml:"relayerURL"`
		CashierContractAddress   string `json:"cashierContractAddress" yaml:"cashierContractAddress"`
		ValidatorContractAddress string `json:"vialidatorContractAddress" yaml:"validatorContractAddress"`
		TransferTableName        string `json:"transferTableName" yaml:"transferTableName"`
		TokenPairs               []struct {
			Token1 string `json:"token1" yaml:"token1"`
			Token2 string `json:"token2" yaml:"token2"`
		} `json:"tokenPairs" yaml:"tokenPairs"`
		StartBlockHeight int `json:"startBlockHeight" yaml:"startBlockHeight"`
	} `json:"cashiers" yaml:"cashiers"`
}

var (
	defaultConfig = Configuration{
		Chain:        	"ethereum",
		Interval:     	time.Minute,
		BatchSize:    	100,
		PrivateKey:   	"",
		SlackWebHook: 	"",
		ClientURL:    	"",
		GrpcPort:     	9080,
		GrpcProxyPort:  9081,
	}
)

var configFile = flag.String("config", "", "path of config file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename>")
		flag.PrintDefaults()
	}
}

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
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok {
		cfg.PrivateKey = pk
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PORT"); ok {
		cfg.GrpcPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PROXY_PORT"); ok {
		cfg.GrpcProxyPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}

	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to decode private key %v\n", err)
	}
	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}

	cashiers := make([]witness.TokenCashier, 0, len(cfg.Cashiers))
	switch cfg.Chain {
	case "iotex":
		conn, err := iotex.NewDefaultGRPCConn(cfg.ClientURL)
		if err != nil {
			log.Fatalf("failed ot create grpc connection %v\n", err)
		}
		iotexClient := iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
		// defer conn.Close()
		for _, cc := range cfg.Cashiers {
			cashierContractAddr, err := address.FromString(cc.CashierContractAddress)
			if err != nil {
				log.Fatalf("failed to parse cashier contract address %s, %v\n", cc.CashierContractAddress, err)
			}
			pairs := make(map[common.Address]common.Address)
			for _, pair := range cc.TokenPairs {
				ioAddr, err := address.FromString(pair.Token1)
				if err != nil {
					log.Fatalf("failed to parse iotex address %s, %v\n", pair.Token1, err)
				}
				if _, ok := pairs[common.BytesToAddress(ioAddr.Bytes())]; ok {
					log.Fatalf("duplicate token key %s\n", pair.Token1)
				}
				pairs[common.BytesToAddress(ioAddr.Bytes())] = common.HexToAddress(pair.Token2)
			}
			cashier, err := witness.NewTokenCashier(
				cc.ID,
				cc.RelayerURL,
				iotexClient,
				cashierContractAddr,
				common.HexToAddress(cc.ValidatorContractAddress),
				witness.NewRecorder(
					db.NewStore(cfg.Database),
					cc.TransferTableName,
					pairs,
				),
				uint64(cc.StartBlockHeight),
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	case "heco", "bsc", "matic":
		// heco and bsc are identical to ethereum
		fallthrough
	case "ethereum":
		ethClient, err := ethclient.DialContext(context.Background(), cfg.ClientURL)
		if err != nil {
			log.Fatal(err)
		}
		for _, cc := range cfg.Cashiers {
			addr, err := address.FromString(cc.ValidatorContractAddress)
			if err != nil {
				log.Fatalf("failed to parse validator contract address %v\n", err)
			}
			pairs := make(map[common.Address]common.Address)
			for _, pair := range cc.TokenPairs {
				if _, ok := pairs[common.HexToAddress(pair.Token1)]; ok {
					log.Fatalf("duplicate token key %s\n", pair.Token1)
				}
				ioAddr, err := address.FromString(pair.Token2)
				if err != nil {
					log.Fatalf("failed to parse iotex address %s, %v\n", pair.Token2, err)
				}
				pairs[common.HexToAddress(pair.Token1)] = common.BytesToAddress(ioAddr.Bytes())
			}
			cashier, err := witness.NewTokenCashierOnEthereum(
				cc.ID,
				cc.RelayerURL,
				ethClient,
				common.HexToAddress(cc.CashierContractAddress),
				common.BytesToAddress(addr.Bytes()),
				witness.NewRecorder(
					db.NewStore(cfg.Database),
					cc.TransferTableName,
					pairs,
				),
				uint64(cc.StartBlockHeight),
				12,
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	default:
		log.Fatalf("unknown chain name %s", cfg.Chain)
	}

	service, err := witness.NewService(
		privateKey,
		cashiers,
		uint16(cfg.BatchSize),
		cfg.Interval,
	)
	if err != nil {
		log.Fatalf("failed to create witness service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start witness service: %v\n", err)
	}
	defer service.Stop(context.Background())

	witness.StartServer(service, cfg.GrpcPort, cfg.GrpcProxyPort)

	log.Println("Serving...")
	select {}
}
