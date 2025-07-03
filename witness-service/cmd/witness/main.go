// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	solclient "github.com/blocto/solana-go-sdk/client"
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness"
	"github.com/iotexproject/iotex-address/address"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain                 string        `json:"chain" yaml:"chain"`
	ClientURL             string        `json:"clientURL" yaml:"clientURL"`
	RelayerURL            string        `json:"relayerURL" yaml:"relayerURL"`
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
		ID                             string      `json:"id" yaml:"id"`
		RelayerURL                     string      `json:"relayerURL" yaml:"relayerURL"`
		WithPayload                    bool        `json:"withPayload" yaml:"withPayload"`
		CashierContractAddress         string      `json:"cashierContractAddress" yaml:"cashierContractAddress"`
		PreviousCashierContractAddress string      `json:"previousCashierContractAddress" yaml:"previousCashierContractAddress"`
		TokenSafeContractAddress       string      `json:"tokenSafeContractAddress" yaml:"tokenSafeContractAddress"`
		ValidatorContractAddress       string      `json:"vialidatorContractAddress" yaml:"validatorContractAddress"`
		TransferTableName              string      `json:"transferTableName" yaml:"transferTableName"`
		TokenPairs                     []TokenPair `json:"tokenPairs" yaml:"tokenPairs"`
		StartBlockHeight               int         `json:"startBlockHeight" yaml:"startBlockHeight"`
		ToSolana                       bool        `json:"toSolana" yaml:"toSolana"`
		DecimalRound                   []struct {
			Token1 string `json:"token1" yaml:"token1"`
			Amount int    `json:"amount" yaml:"amount"`
		} `json:"decimalRound" yaml:"decimalRound"`
		Reverse struct {
			TransferTableName      string   `json:"transferTableName" yaml:"transferTableName"`
			CashierContractAddress string   `json:"cashierContractAddress" yaml:"cashierContractAddress"`
			Tokens                 []string `json:"tokens" yaml:"tokens"`
		}
		QPSLimit    uint32 `json:"qpsLimit" yaml:"qpsLimit"`
		DisablePull bool   `json:"disablePull" yaml:"disablePull"`
	} `json:"cashiers" yaml:"cashiers"`
}

// TokenPair defines a token pair
type TokenPair struct {
	Token1         string   `json:"token1" yaml:"token1"`
	Token2         string   `json:"token2" yaml:"token2"`
	TokenMint      string   `json:"tokenMint" yaml:"tokenMint"`
	TokenProgramID string   `json:"tokenProgramID,omitempty" yaml:"tokenProgramID,omitempty"`
	Whitelist      []string `json:"whitelist" yaml:"whitelist"`
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

func parseTokenPairs(
	tokenPairs []TokenPair,
	destAddrDecoder util.AddressDecoder,
) (map[common.Address]util.Address, map[string][2]util.Address, map[common.Address]map[common.Address]struct{}) {
	pairs := make(map[common.Address]util.Address)
	tokenMintPairs := make(map[string][2]util.Address)
	whitelists := make(map[common.Address]map[common.Address]struct{})
	for _, pair := range tokenPairs {
		token1, err := util.ParseEthAddress(pair.Token1)
		if err != nil {
			log.Fatalf("failed to parse token1 address %s, %v\n", pair.Token1, err)
		}
		if _, ok := pairs[token1]; ok {
			log.Fatalf("duplicate token key %s\n", pair.Token1)
		}
		token2, err := destAddrDecoder.DecodeString(pair.Token2)
		if err != nil {
			log.Fatalf("failed to parse token2 address %s, %v\n", pair.Token2, err)
		}
		pairs[token1] = token2
		if _, ok := destAddrDecoder.(*util.SOLAddressDecoder); ok && len(pair.TokenMint) > 0 {
			mint, err := destAddrDecoder.DecodeString(pair.TokenMint)
			if err != nil {
				log.Fatalf("failed to decode mint address %s, %v\n", pair.TokenMint, err)
			}
			tokenProgram := util.SOLAddressToAddress(solcommon.TokenProgramID)
			if len(pair.TokenProgramID) > 0 && pair.TokenProgramID != tokenProgram.String() {
				tokenProgram, err = destAddrDecoder.DecodeString(pair.TokenProgramID)
				if err != nil {
					log.Fatalf("failed to decode token program id %s, %v\n", pair.TokenProgramID, err)
				}
			}
			// slot 0 is token mint, slot 1 is token program id
			tokenMintPairs[token2.String()] = [2]util.Address{mint, tokenProgram}
		}
		if len(pair.Whitelist) > 0 {
			whitelist := make(map[common.Address]struct{})
			for _, addr := range pair.Whitelist {
				a, err := util.ParseEthAddress(addr)
				if err != nil {
					log.Fatalf("failed to parse whitelist address %s, %v\n", addr, err)
				}
				whitelist[a] = struct{}{}
			}
			whitelists[token1] = whitelist
		}
	}
	return pairs, tokenMintPairs, whitelists
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
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok && cfg.PrivateKey == "" {
		cfg.PrivateKey = pk
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PORT"); ok && cfg.GrpcPort == 0 {
		cfg.GrpcPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if port, ok := os.LookupEnv("WITNESS_GRPC_PROXY_PORT"); ok && cfg.GrpcProxyPort == 0 {
		cfg.GrpcProxyPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if relayerURL, ok := os.LookupEnv("RELAYER_URL"); ok && cfg.RelayerURL == "" {
		cfg.RelayerURL = relayerURL
	}

	// TODO: load more parameters from env
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	if cfg.LarkWebHook != "" {
		util.SetLarkURL(cfg.LarkWebHook)
	}

	if cfg.RelayerURL != "" {
		hasPort := false
		if strings.Contains(cfg.RelayerURL, ":") {
			hasPort = true
			_, _, err := net.SplitHostPort(cfg.RelayerURL)
			if err != nil {
				log.Fatalf("failed to split relayer url %s: %v\n", cfg.RelayerURL, err)
			}
		}
		for i, cc := range cfg.Cashiers {
			switch {
			case strings.HasPrefix(cc.RelayerURL, ":") && !hasPort:
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL + cc.RelayerURL
			case cc.RelayerURL == "":
				cfg.Cashiers[i].RelayerURL = cfg.RelayerURL
			}
		}
	}

	storeFactory := db.NewSQLStoreFactory()
	cashiers := make([]witness.TokenCashier, 0, len(cfg.Cashiers))
	switch cfg.Chain {
	case "heco", "bsc", "matic", "polis", "iotex-e", "iotex", "sepolia", "iotex-testnet", "ethereum":
		ethClient, err := ethclient.Dial(cfg.ClientURL)
		if err != nil {
			log.Fatal(err)
		}
		for _, cc := range cfg.Cashiers {
			var (
				signHandler     witness.SignHandler
				destAddrDecoder util.AddressDecoder
			)
			if !cc.ToSolana {
				destAddrDecoder = util.NewETHAddressDecoder()

				if cfg.PrivateKey != "" {
					privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
					if err != nil {
						log.Fatalf("failed to decode private key %v\n", err)
					}
					util.SetPrefix("witness-" + cfg.Chain + ":" + crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
					log.Println("Witness Service for " + crypto.PubkeyToAddress(privateKey.PublicKey).Hex() + " on chain " + cfg.Chain)
					signHandler = witness.NewSecp256k1SignHandler(privateKey)
				} else {
					log.Println("No Private Key")
				}
			} else {
				destAddrDecoder = util.NewSOLAddressDecoder()

				if cfg.PrivateKey != "" {
					privateKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
					if err != nil {
						log.Fatalf("failed to decode private key %v\n", err)
					}
					var edPrivateKey ed25519.PrivateKey
					switch len(privateKeyBytes) {
					case ed25519.PrivateKeySize:
						edPrivateKey = ed25519.PrivateKey(privateKeyBytes)
					case 32: // Seed from 32 bytes
						edPrivateKey = ed25519.NewKeyFromSeed(privateKeyBytes)
					default:
						log.Fatalf("invalid private key length %d\n", len(privateKeyBytes))
					}
					pbk := solcommon.PublicKeyFromBytes(edPrivateKey.Public().(ed25519.PublicKey)).String()
					log.Println("Witness Service for " + pbk + " on chain " + cfg.Chain)
					signHandler = witness.NewEd25519SignHandler(&edPrivateKey)
				} else {
					log.Println("No Private Key")
				}
			}

			validatorAddr, err := destAddrDecoder.DecodeString(cc.ValidatorContractAddress)
			if err != nil {
				log.Fatalf("failed to parse validator contract address %s: %v\n", cc.ValidatorContractAddress, err)
			}

			var reverseRecorder *witness.Recorder
			if cc.Reverse.CashierContractAddress != "" && cc.Reverse.TransferTableName != "" {
				pairs := make(map[common.Address]util.Address)
				for _, token := range cc.Reverse.Tokens {
					pairs[common.HexToAddress(token)] = util.ETHAddressToAddress(common.HexToAddress(token))
				}
				reverseRecorder = witness.NewRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.Reverse.TransferTableName,
					pairs,
					map[common.Address]map[common.Address]struct{}{},
					map[string][2]util.Address{},
					map[common.Address]int{},
					destAddrDecoder,
				)
			}
			cashierAddr, err := util.ParseEthAddress(cc.CashierContractAddress)
			if err != nil {
				log.Fatalf("invalid cashier address %s: %+v\n", cc.CashierContractAddress, err)
			}
			previousCashierAddr, err := util.ParseEthAddress(cc.PreviousCashierContractAddress)
			if err != nil {
				log.Fatalf("invalid previous cashier address %s: %+v\n", cc.PreviousCashierContractAddress, err)
			}
			tokenSafeAddr, err := util.ParseEthAddress(cc.TokenSafeContractAddress)
			if err != nil {
				log.Fatalf("invalid token safe address %s: %+v\n", cc.TokenSafeContractAddress, err)
			}
			pairs, tokenMintPairs, whitelists := parseTokenPairs(cc.TokenPairs, destAddrDecoder)
			var version witness.Version
			switch {
			case cc.ToSolana:
				version = witness.ToSolana
			case cc.WithPayload:
				version = witness.Payload
			default:
				version = witness.NoPayload
			}
			decimalRound := make(map[common.Address]int)
			for _, r := range cc.DecimalRound {
				addr, err := address.FromString(r.Token1)
				if err != nil {
					log.Fatalf("failed to parse token address %s, %v\n", r.Token1, err)
				}
				decimalRound[common.BytesToAddress(addr.Bytes())] = r.Amount
			}

			cashier, err := witness.NewTokenCashierOnEthereum(
				cc.ID,
				version,
				cc.RelayerURL,
				ethClient,
				cashierAddr,
				previousCashierAddr,
				tokenSafeAddr,
				validatorAddr.Bytes(),
				witness.NewRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.TransferTableName,
					pairs,
					whitelists,
					tokenMintPairs,
					decimalRound,
					destAddrDecoder,
				),
				uint64(cc.StartBlockHeight),
				uint8(cfg.ConfirmBlockNumber),
				signHandler,
				reverseRecorder,
				common.HexToAddress(cc.Reverse.CashierContractAddress),
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	case "solana":
		solClient := solclient.NewClient(cfg.ClientURL)
		for _, cc := range cfg.Cashiers {
			addr, err := util.ParseEthAddress(cc.ValidatorContractAddress)
			if err != nil {
				log.Fatalf("failed to parse validator contract address %v\n", err)
			}
			destAddrDecoder := util.NewETHAddressDecoder()
			pairs := make(map[solcommon.PublicKey]util.Address)
			for _, pair := range cc.TokenPairs {
				token := solcommon.PublicKeyFromString(pair.Token1)
				if _, ok := pairs[token]; ok {
					log.Fatalf("duplicate token key %s\n", pair.Token1)
				}
				token2, err := destAddrDecoder.DecodeString(pair.Token2)
				if err != nil {
					log.Fatalf("failed to parse iotex address %s, %v\n", pair.Token2, err)
				}
				pairs[token] = token2
			}
			decimalRound := make(map[solcommon.PublicKey]int)
			for _, pair := range cc.DecimalRound {
				token := solcommon.PublicKeyFromString(pair.Token1)
				decimalRound[token] = pair.Amount
			}
			var signHandler witness.SignHandler
			if cfg.PrivateKey != "" {
				privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
				if err != nil {
					log.Fatalf("failed to decode private key %v\n", err)
				}
				util.SetPrefix("witness-" + cfg.Chain + ":" + crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
				log.Println("Witness Service for " + crypto.PubkeyToAddress(privateKey.PublicKey).Hex() + " on chain " + cfg.Chain)
				signHandler = witness.NewSecp256k1SignHandler(privateKey)
			} else {
				log.Println("No Private Key")
			}

			cashier, err := witness.NewTokenCashierOnSolana(
				cc.ID,
				cc.RelayerURL,
				solClient,
				solcommon.PublicKeyFromString(cc.CashierContractAddress),
				common.BytesToAddress(addr.Bytes()),
				witness.NewSOLRecorder(
					storeFactory.NewStore(cfg.Database),
					cc.TransferTableName,
					pairs,
					decimalRound,
					destAddrDecoder,
				),
				uint64(cc.StartBlockHeight),
				cc.QPSLimit,
				signHandler,
				cc.DisablePull,
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
