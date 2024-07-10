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
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	Chain                 string        `json:"chain" yaml:"chain"`
	ClientURL             string        `json:"clientURL" yaml:"clientURL"`
	RelayerURL            string        `jsong:"relayerURL" yaml:"relayerURL"`
	Database              db.Config     `json:"database" yaml:"database"`
	PrivateKey            string        `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string        `json:"slackWebHook" yaml:"slackWebHook"`
	LarkWebHook           string        `json:"larkWebHook" yaml:"larkWebHook"`
	ConfirmBlockNumber    int           `json:"confirmBlockNumber" yaml:"confirmBlockNumber"`
	BatchSize             int           `json:"batchSize" yaml:"batchSize"`
	Interval              time.Duration `json:"interval" yaml:"interval"`
	GrpcPort              int           `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort         int           `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	DisableTransferSubmit bool          `json:"disableTransferSubmit" yaml:"disableTransferSubmit"`
	Cashiers              []struct {
		ID                       string      `json:"id" yaml:"id"`
		RelayerURL               string      `json:"relayerURL" yaml:"relayerURL"`
		WithPayload              bool        `json:"withPayload" yaml:"withPayload"`
		CashierContractAddress   string      `json:"cashierContractAddress" yaml:"cashierContractAddress"`
		TokenSafeContractAddress string      `json:"tokenSafeContractAddress" yaml:"tokenSafeContractAddress"`
		ValidatorContractAddress string      `json:"vialidatorContractAddress" yaml:"validatorContractAddress"`
		TransferTableName        string      `json:"transferTableName" yaml:"transferTableName"`
		TokenPairs               []TokenPair `json:"tokenPairs" yaml:"tokenPairs"`
		StartBlockHeight         int         `json:"startBlockHeight" yaml:"startBlockHeight"`
		Reverse                  struct {
			TransferTableName      string   `json:"transferTableName" yaml:"transferTableName"`
			CashierContractAddress string   `json:"cashierContractAddress" yaml:"cashierContractAddress"`
			Tokens                 []string `json:"tokens" yaml:"tokens"`
		}
	} `json:"cashiers" yaml:"cashiers"`
}

// TokenPair defines a token pair
type TokenPair struct {
	Token1 string `json:"token1" yaml:"token1"`
	Token2 string `json:"token2" yaml:"token2"`
}

var (
	defaultConfig = Configuration{
		Chain:              "ethereum",
		Interval:           time.Minute,
		BatchSize:          100,
		ConfirmBlockNumber: 20,
		PrivateKey:         "",
		SlackWebHook:       "",
		LarkWebHook:        "",
		ClientURL:          "",
		GrpcPort:           9080,
		GrpcProxyPort:      9081,
	}

	configFile       = flag.String("config", "", "path of config file")
	secretConfigFile = flag.String("secret", "", "path of secret config file")

	continuously = "continuously"

	blocksFlag = flag.String("blocks", continuously, "block heights")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename> -secret <filename> -blocks <height,height...>")
		flag.PrintDefaults()
	}
}

func parseTokenPairs(tokenPairs []TokenPair) map[common.Address]common.Address {
	pairs := make(map[common.Address]common.Address)
	for _, pair := range tokenPairs {
		token1, err := util.ParseAddress(pair.Token1)
		if err != nil {
			log.Fatalf("failed to parse token1 address %s, %v\n", pair.Token1, err)
		}
		if _, ok := pairs[token1]; ok {
			log.Fatalf("duplicate token key %s\n", pair.Token1)
		}
		token2, err := util.ParseAddress(pair.Token2)
		if err != nil {
			log.Fatalf("failed to parse token2 address %s, %v\n", pair.Token1, err)
		}
		pairs[token1] = token2
	}

	return pairs
}

func main() {
	flag.Parse()
	opts := []config.YAMLOption{config.Static(defaultConfig), config.Expand(os.LookupEnv)}
	if *configFile != "" {
		opts = append(opts, config.File(*configFile))
	}
	if *secretConfigFile != "" {
		opts = append(opts, config.File(*secretConfigFile))
	}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok && pk != "" {
		cfg.PrivateKey = pk
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PORT"); ok && port != "" {
		cfg.GrpcPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PROXY_PORT"); ok && port != "" {
		cfg.GrpcProxyPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if relayerURL, ok := os.LookupEnv("RELAYER_URL"); ok && relayerURL != "" {
		cfg.RelayerURL = relayerURL
	}

	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	if cfg.LarkWebHook != "" {
		util.SetLarkURL(cfg.LarkWebHook)
	}
	var privateKey *ecdsa.PrivateKey
	if cfg.PrivateKey != "" {
		privateKey, err = crypto.HexToECDSA(cfg.PrivateKey)
		if err != nil {
			log.Fatalf("failed to decode private key %v\n", err)
		}
		util.SetPrefix("witness-" + cfg.Chain + ":" + crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		log.Println("Witness Service for " + crypto.PubkeyToAddress(privateKey.PublicKey).Hex() + " on chain " + cfg.Chain)
	} else {
		log.Println("No Private Key")
	}
	if cfg.RelayerURL != "" {
		_, port, err := net.SplitHostPort(cfg.RelayerURL)
		if err != nil {
			log.Fatalf("failed to split relayer url %s: %v\n", cfg.RelayerURL, err)
		}
		for i, cc := range cfg.Cashiers {
			switch {
			case strings.HasPrefix(cc.RelayerURL, ":") && port == "":
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL + cc.RelayerURL
			case cc.RelayerURL == "":
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL
			}
		}
	}

	storeFactory := db.NewSQLStoreFactory()
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
			version := witness.NoPayload
			if cc.WithPayload {
				version = witness.Payload
			}
			cashier, err := witness.NewTokenCashier(
				cc.ID,
				version,
				cc.RelayerURL,
				iotexClient,
				cashierContractAddr,
				common.HexToAddress(cc.ValidatorContractAddress),
				witness.NewRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.TransferTableName,
					parseTokenPairs(cc.TokenPairs),
				),
				uint64(cc.StartBlockHeight),
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	case "heco", "bsc", "matic", "polis", "iotex-e", "sepolia", "iotex-testnet", "ethereum":
		ethClient, err := ethclient.Dial(cfg.ClientURL)
		if err != nil {
			log.Fatal(err)
		}
		for _, cc := range cfg.Cashiers {
			validatorAddr, err := util.ParseAddress(cc.ValidatorContractAddress)
			if err != nil {
				log.Fatalf("failed to parse validator contract address %s: %v\n", cc.ValidatorContractAddress, err)
			}
			var reverseRecorder *witness.Recorder
			if cc.Reverse.CashierContractAddress != "" && cc.Reverse.TransferTableName != "" {
				pairs := make(map[common.Address]common.Address)
				for _, token := range cc.Reverse.Tokens {
					pairs[common.HexToAddress(token)] = common.HexToAddress(token)
				}
				reverseRecorder = witness.NewRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.Reverse.TransferTableName,
					pairs,
				)
			}
			cashierAddr, err := util.ParseAddress(cc.CashierContractAddress)
			if err != nil {
				log.Fatalf("invalid cashier address %s: %+v\n", cc.CashierContractAddress, err)
			}
			tokenSafeAddr, err := util.ParseAddress(cc.TokenSafeContractAddress)
			if err != nil {
				log.Fatalf("invalid token safe address %s: %+v\n", cc.TokenSafeContractAddress, err)
			}
			version := witness.NoPayload
			if cc.WithPayload {
				version = witness.Payload
			}
			cashier, err := witness.NewTokenCashierOnEthereum(
				cc.ID,
				version,
				cc.RelayerURL,
				ethClient,
				cashierAddr,
				tokenSafeAddr,
				validatorAddr,
				witness.NewRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.TransferTableName,
					parseTokenPairs(cc.TokenPairs),
				),
				uint64(cc.StartBlockHeight),
				uint8(cfg.ConfirmBlockNumber),
				reverseRecorder,
				common.HexToAddress(cc.Reverse.CashierContractAddress),
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
		cfg.DisableTransferSubmit,
	)
	if err != nil {
		log.Fatalf("failed to create witness service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start witness service: %v\n", err)
	}
	defer service.Stop(context.Background())

	if *blocksFlag != continuously {
		re := regexp.MustCompile(`^([0-9]*)-([0-9]*)$`)
		for _, hstr := range strings.Split(*blocksFlag, ",") {
			log.Printf("Processing %s\n", hstr)
			var start, end uint64
			if re.MatchString(hstr) {
				matches := re.FindStringSubmatch(hstr)
				start, err = strconv.ParseUint(matches[1], 10, 64)
				if err != nil {
					log.Fatalf("invalid start in %s: %v\n", hstr, err)
				}
				end, err = strconv.ParseUint(matches[2], 10, 64)
				if err != nil {
					log.Fatalf("invalid end in %s: %v\n", hstr, err)
				}
			} else {
				start, err = strconv.ParseUint(hstr, 10, 64)
				if err != nil {
					log.Fatalf("invalid height %s: %v\n", hstr, err)
				}
				end = start
			}
			for height := start; height <= end; height++ {
				if err := service.ProcessOneBlock(height); err != nil {
					log.Fatalf("failed to process block %d: %v\n", height, err)
				}
			}
		}
		log.Println("Done")
		return
	}
	witness.StartServer(service, cfg.GrpcPort, cfg.GrpcProxyPort)

	log.Println("Serving...")
	select {}
}
