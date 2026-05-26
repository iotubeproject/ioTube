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
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
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
	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness"
	"github.com/iotexproject/iotex-address/address"
)

// ApprovalConfig configures the new pre-sign security guards. When `enabled`
// is false (or the block is absent) all guards are disabled and the witness
// behaves as before.
type ApprovalConfig struct {
	Enabled           bool          `json:"enabled" yaml:"enabled"`
	ServerListenAddr  string        `json:"serverListenAddr" yaml:"serverListenAddr"`
	LarkCardWebHook   string        `json:"larkCardWebHook" yaml:"larkCardWebHook"`
	LarkSigningSecret string        `json:"larkSigningSecret" yaml:"larkSigningSecret"`
	WindowDuration    time.Duration `json:"windowDuration" yaml:"windowDuration"`
	// ExplorerTxURL is the base URL for the source-chain block explorer's
	// transaction page. The transfer hash is appended directly, so it must
	// end with a slash or "?tx=" as appropriate.
	// Example: "https://bscscan.com/tx/" or "https://iotexscan.io/tx/"
	ExplorerTxURL string `json:"explorerTxURL" yaml:"explorerTxURL"`
}

// PriceFeedConfig configures the CoinGecko price source used by ApprovalGuard
// to convert per-token amounts into USD. When `enabled` is false the witness
// skips the refresh loop entirely; in that case any cashier with a USD limit
// will fail closed (Block until a price is available).
type PriceFeedConfig struct {
	Enabled         bool          `json:"enabled" yaml:"enabled"`
	BaseURL         string        `json:"baseURL" yaml:"baseURL"`
	APIKey          string        `json:"apiKey" yaml:"apiKey"`
	RefreshInterval time.Duration `json:"refreshInterval" yaml:"refreshInterval"`
	MaxPriceAge     time.Duration `json:"maxPriceAge" yaml:"maxPriceAge"`
	RequestTimeout  time.Duration `json:"requestTimeout" yaml:"requestTimeout"`
}

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain                 string         `json:"chain" yaml:"chain"`
	ClientURL             string         `json:"clientURL" yaml:"clientURL"`
	RelayerURL            string         `json:"relayerURL" yaml:"relayerURL"`
	Database              db.Config      `json:"database" yaml:"database"`
	PrivateKey            string         `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string         `json:"slackWebHook" yaml:"slackWebHook"`
	LarkWebHook           string         `json:"larkWebHook" yaml:"larkWebHook"`
	Approval              ApprovalConfig `json:"approval" yaml:"approval"`
	PriceFeed             PriceFeedConfig `json:"priceFeed" yaml:"priceFeed"`
	ConfirmBlockNumber    int            `json:"confirmBlockNumber" yaml:"confirmBlockNumber"`
	BatchSize             int            `json:"batchSize" yaml:"batchSize"`
	Interval              time.Duration  `json:"interval" yaml:"interval"`
	GrpcPort              int            `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort         int            `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	DisableTransferSubmit bool           `json:"disableTransferSubmit" yaml:"disableTransferSubmit"`
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
		QPSLimit           uint32 `json:"qpsLimit" yaml:"qpsLimit"`
		DisablePull        bool   `json:"disablePull" yaml:"disablePull"`
		WindowValueLimit   string `json:"windowValueLimit" yaml:"windowValueLimit"`
		SingleTxValueLimit string `json:"singleTxValueLimit" yaml:"singleTxValueLimit"`
	} `json:"cashiers" yaml:"cashiers"`
}

// TokenPair defines a token pair
type TokenPair struct {
	Token1         string   `json:"token1" yaml:"token1"`
	Token2         string   `json:"token2" yaml:"token2"`
	TokenMint      string   `json:"tokenMint" yaml:"tokenMint"`
	TokenProgramID string   `json:"tokenProgramID,omitempty" yaml:"tokenProgramID,omitempty"`
	Whitelist      []string `json:"whitelist" yaml:"whitelist"`
	// CoinGeckoID is the CoinGecko coin id (e.g. "weth", "wrapped-bitcoin")
	// used to fetch this token's USD price. Optional: when omitted, the
	// witness resolves it at startup via CoinGecko's contract-address lookup
	// (/coins/{platform}/contract/{token1}). Set it explicitly only to
	// override the auto-resolved id — useful for bridged tokens CoinGecko
	// doesn't index under the bridged contract.
	CoinGeckoID string `json:"coingeckoID" yaml:"coingeckoID"`
	// Decimals is the number of decimal places of the value returned by
	// AbstractTransfer.Amount() — i.e. after any DecimalRound adjustment.
	// Required when the parent cashier has a USD limit configured.
	Decimals int `json:"decimals" yaml:"decimals"`
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

	approvalGuards := make(map[string]*witness.ApprovalGuard)
	cgClient, priceCache := newPriceFeed(cfg)
	resolver := newCoingeckoResolver(cgClient, cfg.Chain)
	var allTokenMetas []witness.TokenMeta

	storeFactory := db.NewSQLStoreFactory()
	cashiers := make([]witness.TokenCashier, 0, len(cfg.Cashiers))
	switch cfg.Chain {
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

			cashierPubKey := solcommon.PublicKeyFromString(cc.CashierContractAddress)
			solRecorder := witness.NewSOLRecorder(
				storeFactory.NewStore(cfg.Database),
				cc.TransferTableName,
				pairs,
				decimalRound,
				destAddrDecoder,
			)
			cashier, err := witness.NewTokenCashierOnSolana(
				cc.ID,
				cc.RelayerURL,
				solClient,
				cashierPubKey,
				common.BytesToAddress(addr.Bytes()),
				solRecorder,
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
	default: // "heco", "bsc", "matic", "polis", "iotex-e", "iotex", "sepolia", "iotex-testnet", "ethereum":
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

			ethRecorder := witness.NewRecorder(
				storeFactory.NewStore(cfg.Database),
				cc.TransferTableName,
				pairs,
				whitelists,
				tokenMintPairs,
				decimalRound,
				destAddrDecoder,
			)
			cashierKey := cashierAddr.Hex()
			var ethTokenMetas []witness.TokenMeta
			if cc.ToSolana {
				// For ETH→SOL cashiers token1 is still an EVM address; same path.
				ethTokenMetas = tokenMetasForEthereum(cc.TokenPairs, resolver)
			} else if cfg.Chain == "iotex" || cfg.Chain == "iotex-e" || cfg.Chain == "iotex-testnet" {
				ethTokenMetas = tokenMetasForIotex(cc.TokenPairs, resolver)
			} else {
				ethTokenMetas = tokenMetasForEthereum(cc.TokenPairs, resolver)
			}
			allTokenMetas = append(allTokenMetas, ethTokenMetas...)
			guard := buildApprovalGuard(cfg.Approval, cashierKey, cc.WindowValueLimit, cc.SingleTxValueLimit, ethTokenMetas, priceCache, ethRecorder)
			if guard != nil {
				approvalGuards[cashierKey] = guard
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
				ethRecorder,
				uint64(cc.StartBlockHeight),
				uint8(cfg.ConfirmBlockNumber),
				signHandler,
				reverseRecorder,
				common.HexToAddress(cc.Reverse.CashierContractAddress),
				guard,
			)
			if err != nil {
				log.Fatalf("failed to create cashier %v\n", err)
			}
			cashiers = append(cashiers, cashier)
		}
	}

	if cfg.Approval.Enabled {
		if cfg.Approval.ServerListenAddr == "" {
			log.Fatal("approval is enabled but approval.serverListenAddr is not configured")
		}
		if len(approvalGuards) == 0 {
			log.Fatal("approval is enabled but no cashier has windowValueLimit or singleTxValueLimit configured")
		}
		if err := witness.NewApprovalServer(cfg.Approval.ServerListenAddr, cfg.Approval.LarkSigningSecret, approvalGuards).Start(); err != nil {
			log.Fatalf("failed to start approval server: %v\n", err)
		}
	}

	if priceRunner := startPriceFeedRunner(cfg, cgClient, priceCache, collectCoinGeckoIDsFromMetas(allTokenMetas)); priceRunner != nil {
		defer priceRunner.Close()
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

func buildApprovalGuard(
	ac ApprovalConfig,
	cashierKey string,
	windowLimitStr, singleTxLimitStr string,
	tokens []witness.TokenMeta,
	prices witness.PriceSource,
	recorder witness.AbstractRecorder,
) *witness.ApprovalGuard {
	if !ac.Enabled {
		return nil
	}
	windowLimit := parseUsdLimit(windowLimitStr)
	singleTxLimit := parseUsdLimit(singleTxLimitStr)
	if windowLimit == nil && singleTxLimit == nil {
		return nil
	}
	if prices == nil {
		log.Fatalf("approval guard for %s configured with USD limit but no price source", cashierKey)
	}
	for _, tk := range tokens {
		if tk.CoinGeckoID == "" {
			log.Fatalf(
				"approval guard for %s: could not resolve CoinGecko id for token %s (set coingeckoID in config to override)",
				cashierKey, tk.Token,
			)
		}
		if tk.Decimals < 0 {
			log.Fatalf(
				"approval guard for %s: token %s has negative decimals (%d)",
				cashierKey, tk.Token, tk.Decimals,
			)
		}
	}
	return witness.NewApprovalGuard(
		cashierKey,
		ac.WindowDuration,
		windowLimit,
		singleTxLimit,
		tokens,
		prices,
		recorder,
		ac.LarkCardWebHook,
		ac.ExplorerTxURL,
	)
}

// parseUsdLimit parses a decimal USD amount (e.g. "100000", "5000.50") into a
// *big.Float. Empty / "0" / non-positive values disable that limit dimension.
func parseUsdLimit(s string) *big.Float {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return nil
	}
	v, _, err := big.ParseFloat(s, 10, 80, big.ToNearestEven)
	if err != nil {
		log.Fatalf("invalid USD limit %q in approval config: %v", s, err)
	}
	if v.Sign() <= 0 {
		return nil
	}
	return v
}

// newPriceFeed creates the CoinGecko client and the in-process PriceCache.
// Returns (nil, nil) when priceFeed.enabled is false — the witness behaves
// exactly as it did before USD limits existed. Does NOT start the periodic
// refresh runner; that happens later via startPriceFeedRunner after all
// CoinGecko ids (explicit + auto-resolved) are known.
func newPriceFeed(cfg Configuration) (*util.CoinGeckoClient, *util.PriceCache) {
	if !cfg.PriceFeed.Enabled {
		return nil, nil
	}
	maxAge := cfg.PriceFeed.MaxPriceAge
	if maxAge == 0 {
		maxAge = 10 * time.Minute
	}
	client := util.NewCoinGeckoClient(cfg.PriceFeed.BaseURL, cfg.PriceFeed.APIKey, cfg.PriceFeed.RequestTimeout)
	return client, util.NewPriceCache(maxAge)
}

// startPriceFeedRunner starts the periodic price refresh against the given
// CoinGecko ids. Caller must defer runner.Close(). Returns nil when the
// price feed is disabled or no ids need fetching.
func startPriceFeedRunner(cfg Configuration, client *util.CoinGeckoClient, cache *util.PriceCache, ids []string) dispatcher.Runner {
	if client == nil || cache == nil {
		return nil
	}
	if len(ids) == 0 {
		log.Println("priceFeed.enabled=true but no coingecko ids resolved; skipping refresh loop")
		return nil
	}
	interval := cfg.PriceFeed.RefreshInterval
	if interval == 0 {
		interval = 2 * time.Minute
	}

	runner, err := dispatcher.NewRunner(interval, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), interval)
		defer cancel()
		prices, err := client.FetchUSDPrices(ctx, ids)
		if err != nil {
			util.Alert(fmt.Sprintf("price feed refresh failed: %v", err))
			return err
		}
		cache.Replace(prices)
		return nil
	})
	if err != nil {
		log.Fatalf("failed to create price feed runner: %v", err)
	}
	// Run once synchronously so guards have prices before signing starts.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if prices, err := client.FetchUSDPrices(ctx, ids); err != nil {
		log.Printf("initial price fetch failed (will retry on schedule): %v", err)
	} else {
		cache.Replace(prices)
	}
	if err := runner.Start(); err != nil {
		log.Fatalf("failed to start price feed runner: %v", err)
	}
	return runner
}

// collectCoinGeckoIDsFromMetas returns the deduped CoinGecko id list for the
// price feed. After auto-resolution every TokenMeta with a USD-limit cashier
// has its id populated; this is the canonical source.
func collectCoinGeckoIDsFromMetas(metas []witness.TokenMeta) []string {
	seen := make(map[string]struct{}, len(metas))
	out := make([]string, 0, len(metas))
	for _, m := range metas {
		id := strings.ToLower(strings.TrimSpace(m.CoinGeckoID))
		if id == "" {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// coingeckoPlatformForChain maps witness chain identifiers to CoinGecko's
// asset-platform string (used in /coins/{platform}/contract/{address}).
// Returns "" when no mapping is known — caller treats that as
// "auto-resolution unavailable on this chain".
func coingeckoPlatformForChain(chain string) string {
	switch chain {
	case "ethereum", "sepolia":
		return "ethereum"
	case "bsc":
		return "binance-smart-chain"
	case "matic":
		return "polygon-pos"
	case "iotex", "iotex-e", "iotex-testnet":
		return "iotex"
	case "solana":
		return "solana"
	case "heco":
		return "huobi-token"
	case "polis":
		return "polis-chain"
	}
	return ""
}

// coingeckoNativeIDForChain returns the CoinGecko id for the chain's native
// gas token, used when a TokenPair has the zero address as its token1
// (currently no in-tree config does, but the future-proofing is cheap).
func coingeckoNativeIDForChain(chain string) string {
	switch chain {
	case "ethereum", "sepolia":
		return "ethereum"
	case "bsc":
		return "binancecoin"
	case "matic":
		return "matic-network"
	case "iotex", "iotex-e", "iotex-testnet":
		return "iotex"
	case "solana":
		return "solana"
	}
	return ""
}

// coingeckoResolver auto-resolves CoinGecko ids from token contract
// addresses for the witness's bound chain. A nil resolver means
// auto-resolution is disabled (price feed off or chain not mapped); callers
// fall back to whatever id was set explicitly in config. Results are cached
// in-process so repeated lookups (same token across multiple cashiers) make
// at most one HTTP call.
type coingeckoResolver struct {
	client   *util.CoinGeckoClient
	chain    string
	platform string
	cache    map[string]string // key: lowercase hex addr with 0x → CG id ("" = lookup failed)
}

func newCoingeckoResolver(client *util.CoinGeckoClient, chain string) *coingeckoResolver {
	if client == nil {
		return nil
	}
	platform := coingeckoPlatformForChain(chain)
	if platform == "" {
		log.Printf("coingecko: no platform mapping for chain %q; auto-resolution disabled (config must set coingeckoID)", chain)
		return nil
	}
	return &coingeckoResolver{
		client:   client,
		chain:    chain,
		platform: platform,
		cache:    make(map[string]string),
	}
}

// resolveByEthAddress looks up the CoinGecko id for an EVM contract address.
// Returns "" if the resolver is nil, the chain has no platform mapping, or
// the API has no record — in any of those cases the caller's existing
// fallback (config override or buildApprovalGuard's hard check) takes over.
func (r *coingeckoResolver) resolveByEthAddress(addr common.Address) string {
	if r == nil {
		return ""
	}
	key := strings.ToLower(addr.Hex())
	if id, ok := r.cache[key]; ok {
		return id
	}
	if isZeroEthAddress(key) {
		id := coingeckoNativeIDForChain(r.chain)
		r.cache[key] = id
		return id
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	id, err := r.client.ResolveIDByContract(ctx, r.platform, key)
	if err != nil {
		if errors.Is(err, util.ErrCoinGeckoIDNotFound) {
			log.Printf("coingecko: no coin indexed for %s/%s", r.platform, key)
		} else {
			log.Printf("coingecko: auto-resolve failed for %s/%s: %v", r.platform, key, err)
		}
		r.cache[key] = ""
		return ""
	}
	r.cache[key] = id
	log.Printf("coingecko: %s/%s → %s", r.platform, key, id)
	return id
}

func isZeroEthAddress(hexAddr string) bool {
	h := strings.TrimPrefix(strings.ToLower(hexAddr), "0x")
	if len(h) == 0 {
		return false
	}
	for _, c := range h {
		if c != '0' {
			return false
		}
	}
	return true
}

// tokenMetasForEthereum builds TokenMeta for an EVM cashier — token1 is an
// 0x... 20-byte address. When resolver is non-nil and a pair omits
// CoinGeckoID, the id is auto-resolved via CoinGecko's contract lookup.
func tokenMetasForEthereum(pairs []TokenPair, resolver *coingeckoResolver) []witness.TokenMeta {
	out := make([]witness.TokenMeta, 0, len(pairs))
	for _, p := range pairs {
		if resolver == nil && p.CoinGeckoID == "" && p.Decimals == 0 {
			continue
		}
		addr, err := util.ParseEthAddress(p.Token1)
		if err != nil {
			log.Fatalf("tokenMetasForEthereum: invalid token1 %s: %v", p.Token1, err)
		}
		cgID := p.CoinGeckoID
		if cgID == "" {
			cgID = resolver.resolveByEthAddress(common.BytesToAddress(addr.Bytes()))
		}
		out = append(out, witness.TokenMeta{
			Token:       strings.ToLower(strings.TrimPrefix(addr.String(), "0x")),
			CoinGeckoID: cgID,
			Decimals:    p.Decimals,
		})
	}
	return out
}

// tokenMetasForIotex handles `io1...` bech32 addresses (IoTeX source-chain
// cashiers). util.ParseEthAddress also accepts io1 prefixes — verified via
// existing call sites — so we reuse it.
func tokenMetasForIotex(pairs []TokenPair, resolver *coingeckoResolver) []witness.TokenMeta {
	out := make([]witness.TokenMeta, 0, len(pairs))
	for _, p := range pairs {
		if resolver == nil && p.CoinGeckoID == "" && p.Decimals == 0 {
			continue
		}
		addr, err := util.ParseEthAddress(p.Token1)
		if err != nil {
			log.Fatalf("tokenMetasForIotex: invalid token1 %s: %v", p.Token1, err)
		}
		cgID := p.CoinGeckoID
		if cgID == "" {
			cgID = resolver.resolveByEthAddress(common.BytesToAddress(addr.Bytes()))
		}
		out = append(out, witness.TokenMeta{
			Token:       strings.ToLower(strings.TrimPrefix(addr.String(), "0x")),
			CoinGeckoID: cgID,
			Decimals:    p.Decimals,
		})
	}
	return out
}
