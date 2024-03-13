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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	btcchain "github.com/btcsuite/btcwallet/chain"
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
		ID                       string `json:"id" yaml:"id"`
		RelayerURL               string `json:"relayerURL" yaml:"relayerURL"`
		CashierContractAddress   string `json:"cashierContractAddress" yaml:"cashierContractAddress"`
		TokenSafeContractAddress string `json:"tokenSafeContractAddress" yaml:"tokenSafeContractAddress"`
		ValidatorContractAddress string `json:"validatorContractAddress" yaml:"validatorContractAddress"`
		TransferTableName        string `json:"transferTableName" yaml:"transferTableName"`
		TokenPairs               []struct {
			Token1 string `json:"token1" yaml:"token1"`
			Token2 string `json:"token2" yaml:"token2"`
		} `json:"tokenPairs" yaml:"tokenPairs"`
		StartBlockHeight int `json:"startBlockHeight" yaml:"startBlockHeight"`
		Reverse          struct {
			TransferTableName      string   `json:"transferTableName" yaml:"transferTableName"`
			CashierContractAddress string   `json:"cashierContractAddress" yaml:"cashierContractAddress"`
			Tokens                 []string `json:"tokens" yaml:"tokens"`
		}
		SupportBitcoinRecipient bool   `json:"supportBitcoinRecipient" yaml:"supportBitcoinRecipient"`
		MinTipFee               uint64 `json:"minTipFee" yaml:"minTipFee"`
	} `json:"cashiers" yaml:"cashiers"`
	BitcoinConfig struct {
		RPCHost   string `json:"rpcHost" yaml:"rpcHost"`
		RPCUser   string `json:"rpcUser" yaml:"rpcUser"`
		RPCPass   string `json:"rpcPass" yaml:"rpcPass"`
		EnableTLS bool   `json:"enableTLS" yaml:"enableTLS"`
	} `json:"bitcoinConfig" yaml:"bitcoinConfig"`
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

func main() {
	log.SetOutput(os.Stdout)

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
	if relayerURL, ok := os.LookupEnv("RELAYER_URL"); ok {
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
		for i, cc := range cfg.Cashiers {
			switch {
			case strings.HasPrefix(cc.RelayerURL, ":"):
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL + cc.RelayerURL
			case cc.RelayerURL == "":
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL
			}
		}
	}

	var bTCValidator *witness.BTCValidator
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
			var addrDecoder util.AddressDecoder
			if cc.SupportBitcoinRecipient {
				addrDecoder = util.NewBTCAddressDecoder(&chaincfg.TestNet3Params)
			} else {
				addrDecoder = util.NewETHAddressDecoder()
			}
			recorder := witness.NewRecorder(
				db.NewStore(cfg.Database),
				cc.TransferTableName,
				pairs,
				addrDecoder,
			)

			if cc.SupportBitcoinRecipient {
				if bTCValidator != nil {
					log.Fatalf("only one btc casiher is supported in source chain")
				}
				cli, err := rpcclient.New(&rpcclient.ConnConfig{
					Host:         cfg.BitcoinConfig.RPCHost,
					User:         cfg.BitcoinConfig.RPCUser,
					Pass:         cfg.BitcoinConfig.RPCPass,
					HTTPPostMode: true,
					DisableTLS:   !cfg.BitcoinConfig.EnableTLS,
				}, nil)
				if err != nil {
					log.Fatal(err)
				}
				defer cli.Shutdown()

				bTCValidator = witness.NewBTCTransferValidator(
					privateKey,
					cli,
					cfg.Interval,
					cc.RelayerURL,
					recorder,
				)
			}

			validatorContractAddress := common.HexToAddress(cc.ValidatorContractAddress).Bytes()
			if cc.SupportBitcoinRecipient {
				validatorContractAddress = []byte{}
			}

			cashier, err := witness.NewTokenCashier(
				cc.ID,
				cc.RelayerURL,
				iotexClient,
				cashierContractAddr,
				validatorContractAddress,
				recorder,
				uint64(cc.StartBlockHeight),
				addrDecoder,
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	case "heco", "bsc", "matic", "polis":
		// heco and bsc are identical to ethereum
		fallthrough
	case "ethereum":
		ethClient, err := ethclient.Dial(cfg.ClientURL)
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
			var reverseRecorder *witness.Recorder
			if cc.Reverse.CashierContractAddress != "" && cc.Reverse.TransferTableName != "" {
				pairs := make(map[common.Address]common.Address)
				for _, token := range cc.Reverse.Tokens {
					pairs[common.HexToAddress(token)] = common.HexToAddress(token)
				}
				reverseRecorder = witness.NewRecorder(
					db.NewStore(cfg.Database),
					cc.Reverse.TransferTableName,
					pairs,
					util.NewETHAddressDecoder(),
				)
			}
			cashier, err := witness.NewTokenCashierOnEthereum(
				cc.ID,
				cc.RelayerURL,
				ethClient,
				common.HexToAddress(cc.CashierContractAddress),
				common.HexToAddress(cc.TokenSafeContractAddress),
				common.BytesToAddress(addr.Bytes()),
				witness.NewRecorder(
					db.NewStore(cfg.Database),
					cc.TransferTableName,
					pairs,
					util.NewETHAddressDecoder(),
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
	case "bitcoin-testnet":
		if len(cfg.Cashiers) != 1 {
			log.Fatalln("only one cashier is supported for bitcoin")
		}

		// Testnet
		testnetParam := &chaincfg.TestNet3Params
		btcCfg := &btcchain.BitcoindConfig{
			ChainParams: testnetParam,
			Host:        cfg.BitcoinConfig.RPCHost,
			User:        cfg.BitcoinConfig.RPCUser,
			Pass:        cfg.BitcoinConfig.RPCPass,
			PollingConfig: &btcchain.PollingConfig{
				BlockPollingInterval: time.Millisecond * 100,
				TxPollingInterval:    time.Millisecond * 100,
			},
		}

		chainConn, err := btcchain.NewBitcoindConn(btcCfg)
		if err != nil {
			log.Fatalf("failed to create bitcoind connection %v\n", err)
		}
		chainConn.Start()
		defer chainConn.Stop()
		btcClient := chainConn.NewBitcoindClient()
		btcClient.Start()
		defer btcClient.Stop()

		// the string of CashierContractAddress is also used as the musig pubkey for bitcoin
		cashierPubkey, err := util.HexToPubkey(cfg.Cashiers[0].CashierContractAddress)
		if err != nil {
			log.Fatalf("failed to parse cashier contract address %s, error %v\n", cfg.Cashiers[0].CashierContractAddress, err)
		}

		cc := cfg.Cashiers[0]
		cashier, err := witness.NewTokenCashierOnBitcoin(
			cc.ID,
			cc.RelayerURL,
			btcClient,
			testnetParam,
			cashierPubkey,
			common.HexToAddress(cc.TokenSafeContractAddress),
			common.HexToAddress(cc.ValidatorContractAddress),
			witness.NewBTCRecorder(
				db.NewStore(cfg.Database),
				cc.TransferTableName,
			),
			uint64(cc.StartBlockHeight),
			uint8(cfg.ConfirmBlockNumber),
			cc.MinTipFee,
		)
		if err != nil {
			log.Fatalf("failed to create cashier %v\n", err)
		}
		cashiers = append(cashiers, cashier)
	default:
		log.Fatalf("unknown chain name %s", cfg.Chain)
	}

	service, err := witness.NewService(
		privateKey,
		cashiers,
		bTCValidator,
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
