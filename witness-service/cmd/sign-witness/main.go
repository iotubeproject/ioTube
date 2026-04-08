// Copyright (c) 2024 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/util/instruction"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	Chain                 string          `json:"chain" yaml:"chain"`
	ClientURL             string          `json:"clientURL" yaml:"clientURL"`
	RelayerURL            string          `json:"relayerURL" yaml:"relayerURL"`
	Database              db.Config       `json:"database" yaml:"database"`
	PrivateKey            string          `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string          `json:"slackWebHook" yaml:"slackWebHook"`
	LarkWebHook           string          `json:"larkWebHook" yaml:"larkWebHook"`
	ConfirmBlockNumber    int             `json:"confirmBlockNumber" yaml:"confirmBlockNumber"`
	BatchSize             int             `json:"batchSize" yaml:"batchSize"`
	Interval              time.Duration   `json:"interval" yaml:"interval"`
	GrpcPort              int             `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort         int             `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	DisableTransferSubmit bool            `json:"disableTransferSubmit" yaml:"disableTransferSubmit"`
	Cashiers              []CashierConfig `json:"cashiers" yaml:"cashiers"`
}

// TokenPair defines a token pair
type TokenPair struct {
	Token1         string   `json:"token1" yaml:"token1"`
	Token2         string   `json:"token2" yaml:"token2"`
	TokenMint      string   `json:"tokenMint" yaml:"tokenMint"`
	TokenProgramID string   `json:"tokenProgramID,omitempty" yaml:"tokenProgramID,omitempty"`
	Whitelist      []string `json:"whitelist" yaml:"whitelist"`
}

// DecimalRound defines decimal rounding config
type DecimalRound struct {
	Token1 string `json:"token1" yaml:"token1"`
	Amount int    `json:"amount" yaml:"amount"`
}

// CashierConfig defines the configuration for a cashier
type CashierConfig struct {
	ID                             string         `json:"id" yaml:"id"`
	RelayerURL                     string         `json:"relayerURL" yaml:"relayerURL"`
	WithPayload                    bool           `json:"withPayload" yaml:"withPayload"`
	CashierContractAddress         string         `json:"cashierContractAddress" yaml:"cashierContractAddress"`
	PreviousCashierContractAddress string         `json:"previousCashierContractAddress" yaml:"previousCashierContractAddress"`
	TokenSafeContractAddress       string         `json:"tokenSafeContractAddress" yaml:"tokenSafeContractAddress"`
	ValidatorContractAddress       string         `json:"validatorContractAddress" yaml:"validatorContractAddress"`
	TransferTableName              string         `json:"transferTableName" yaml:"transferTableName"`
	TokenPairs                     []TokenPair    `json:"tokenPairs" yaml:"tokenPairs"`
	DecimalRound                   []DecimalRound `json:"decimalRound" yaml:"decimalRound"`
	StartBlockHeight               int            `json:"startBlockHeight" yaml:"startBlockHeight"`
	ToSolana                       bool           `json:"toSolana" yaml:"toSolana"`
	QPSLimit                       uint32         `json:"qpsLimit" yaml:"qpsLimit"`
	DisablePull                    bool           `json:"disablePull" yaml:"disablePull"`
}

// TransferData holds the transfer information
type TransferData struct {
	Cashier   common.Address
	Token     common.Address
	CoToken   common.Address
	Index     *big.Int
	Sender    common.Address
	Recipient []byte
	Payload   []byte
	Amount    *big.Int
}

// Address implements util.Address interface
type transferAddress struct {
	bytes []byte
	str   string
}

func (a *transferAddress) Bytes() []byte   { return a.bytes }
func (a *transferAddress) String() string  { return a.str }
func (a *transferAddress) Address() any { return a.bytes }

var (
	configFile       = flag.String("config", "", "path of config file")
	secretConfigFile = flag.String("secret", "", "path of secret config file")
	cashierID        = flag.String("cashier", "", "cashier ID (from config)")
	cashierAddr      = flag.String("cashier-address", "", "cashier contract address")
	validatorAddr    = flag.String("validator-address", "", "validator contract address")
	tokenAddr        = flag.String("token", "", "source token address (for reference in output)")
	coTokenAddr      = flag.String("cotoken", "", "co-token address")
	index            = flag.String("index", "", "transfer index")
	senderAddr       = flag.String("sender", "", "sender address")
	recipientAddr    = flag.String("recipient", "", "recipient address")
	amountStr        = flag.String("amount", "", "transfer amount")
	payloadHex       = flag.String("payload", "", "payload in hex (optional)")
	toSolana         = flag.Bool("to-solana", false, "destination is solana (uses Ed25519)")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[options]")
		fmt.Fprintln(os.Stderr, "\nGenerates a witness signature for a transfer transaction.")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  # Using config file:")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-config witness-config.yaml -secret secret.yaml -cashier iotex-testnet-to-bsc-testnet -cotoken 0x... -index 123 -sender 0x... -recipient 0x... -amount 1000000000000000000")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  # Manual mode (private key from secret config):")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-secret secret.yaml -cashier-address 0x... -validator-address 0x... -cotoken 0x... -index 123 -sender 0x... -recipient 0x... -amount 1000000000000000000")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  # For Solana destination:")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-config witness-config.yaml -secret secret.yaml -to-solana -cashier iotex-solana -cotoken <solana-token> -index 123 -sender 0x... -recipient <solana-pubkey> -amount 1000000")
	}
}

func main() {
	flag.Parse()

	// Load configuration
	opts := []config.YAMLOption{config.Expand(os.LookupEnv)}
	if *configFile != "" {
		opts = append(opts, config.File(*configFile))
	}
	if *secretConfigFile != "" {
		opts = append(opts, config.File(*secretConfigFile))
	}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	// Allow override from env
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok && cfg.PrivateKey == "" {
		cfg.PrivateKey = pk
	}

	if cfg.PrivateKey == "" {
		log.Fatal("Private key is required. Set it in the secret config file or WITNESS_PRIVATE_KEY env var.")
	}

	// Get cashier config if specified (case-insensitive matching)
	var cashierCfg *CashierConfig
	if *cashierID != "" {
		cashierIDLower := strings.ToLower(*cashierID)
		for _, cc := range cfg.Cashiers {
			if strings.ToLower(cc.ID) == cashierIDLower {
				cashierCfg = &cc
				break
			}
		}
		if cashierCfg == nil {
			log.Fatalf("Cashier ID '%s' not found in config\n", *cashierID)
		}
	}

	// Resolve addresses
	cashier := resolveAddress(*cashierAddr, cashierCfg, func(cc CashierConfig) string { return cc.CashierContractAddress }, "cashier")
	validator := resolveAddress(*validatorAddr, cashierCfg, func(cc CashierConfig) string { return cc.ValidatorContractAddress }, "validator")

	// Determine if destination is Solana
	isToSolana := *toSolana
	if cashierCfg != nil && cashierCfg.ToSolana {
		isToSolana = true
	}

	// Parse required arguments
	if *coTokenAddr == "" {
		log.Fatal("-cotoken is required")
	}
	if *index == "" {
		log.Fatal("-index is required")
	}
	if *senderAddr == "" {
		log.Fatal("-sender is required")
	}
	if *recipientAddr == "" {
		log.Fatal("-recipient is required")
	}
	if *amountStr == "" {
		log.Fatal("-amount is required")
	}

	// Parse co-token address
	var coToken []byte
	if isToSolana {
		coToken = solcommon.PublicKeyFromString(*coTokenAddr).Bytes()
	} else {
		coToken = common.HexToAddress(*coTokenAddr).Bytes()
	}

	// Parse index
	indexVal := new(big.Int)
	if _, ok := indexVal.SetString(*index, 10); !ok {
		log.Fatalf("Invalid index: %s\n", *index)
	}

	// Parse sender
	sender := common.HexToAddress(*senderAddr)

	// Parse recipient
	var recipient []byte
	if isToSolana {
		recipient = solcommon.PublicKeyFromString(*recipientAddr).Bytes()
	} else {
		recipient = common.HexToAddress(*recipientAddr).Bytes()
	}

	// Parse amount
	amount := new(big.Int)
	if _, ok := amount.SetString(*amountStr, 10); !ok {
		log.Fatalf("Invalid amount: %s\n", *amountStr)
	}

	// Parse payload (optional)
	var payload []byte
	if *payloadHex != "" {
		payload, err = hex.DecodeString(strings.TrimPrefix(*payloadHex, "0x"))
		if err != nil {
			log.Fatalf("Invalid payload hex: %v\n", err)
		}
	}

	// Generate signature
	var id common.Hash
	var witnessPubKey []byte
	var signature []byte

	if isToSolana {
		// Ed25519 signature for Solana
		privateKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
		if err != nil {
			log.Fatalf("Failed to decode private key: %v\n", err)
		}
		var edPrivateKey ed25519.PrivateKey
		switch len(privateKeyBytes) {
		case ed25519.PrivateKeySize:
			edPrivateKey = ed25519.PrivateKey(privateKeyBytes)
		case 32: // Seed from 32 bytes
			edPrivateKey = ed25519.NewKeyFromSeed(privateKeyBytes)
		default:
			log.Fatalf("Invalid private key length %d\n", len(privateKeyBytes))
		}

		// Serialize payload using Borsh
		data, err := instruction.SerializePayload(
			validator.Bytes(),
			cashier.Bytes(),
			coToken,
			indexVal.Uint64(),
			sender.Hex(),
			recipient,
			amount.Uint64(),
			payload,
		)
		if err != nil {
			log.Fatalf("Failed to serialize payload: %v\n", err)
		}
		id = crypto.Keccak256Hash(data)
		signature = ed25519.Sign(edPrivateKey, id[:])
		witnessPubKey = edPrivateKey.Public().(ed25519.PublicKey)
	} else {
		// Secp256k1 signature for Ethereum chains
		privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
		if err != nil {
			log.Fatalf("Failed to decode private key: %v\n", err)
		}

		// Create the transfer ID hash
		id = crypto.Keccak256Hash(
			validator.Bytes(),
			cashier.Bytes(),
			coToken,
			math.U256Bytes(indexVal),
			sender.Bytes(),
			recipient,
			math.U256Bytes(amount),
			payload,
		)
		signature, err = crypto.Sign(id.Bytes(), privateKey)
		if err != nil {
			log.Fatalf("Failed to sign: %v\n", err)
		}
		witnessPubKey = crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()
	}

	// Output results
	fmt.Println("=== Witness Signature ===")
	fmt.Printf("Transfer ID: 0x%x\n", id.Bytes())
	fmt.Printf("Witness Public Key: 0x%x\n", witnessPubKey)
	fmt.Printf("Signature: 0x%x\n", signature)
	fmt.Println()
	fmt.Println("=== Details ===")
	fmt.Printf("Validator: %s\n", validator.Hex())
	fmt.Printf("Cashier: %s\n", cashier.Hex())
	if isToSolana {
		fmt.Printf("CoToken: 0x%x\n", coToken)
		fmt.Printf("Recipient: 0x%x\n", recipient)
	} else {
		fmt.Printf("CoToken: %s\n", common.BytesToAddress(coToken).Hex())
		fmt.Printf("Recipient: %s\n", common.BytesToAddress(recipient).Hex())
	}
	fmt.Printf("Index: %s\n", indexVal.String())
	fmt.Printf("Sender: %s\n", sender.Hex())
	fmt.Printf("Amount: %s\n", amount.String())
	if len(payload) > 0 {
		fmt.Printf("Payload: 0x%x\n", payload)
	}
	fmt.Printf("To Solana: %v\n", isToSolana)
}

func resolveAddress(flagValue string, cashierCfg *CashierConfig, getFromConfig func(CashierConfig) string, name string) common.Address {
	if flagValue != "" {
		parsed, err := util.ParseEthAddress(flagValue)
		if err != nil {
			log.Fatalf("Failed to parse -%s-address: %v\n", name, err)
		}
		return parsed
	}
	if cashierCfg != nil {
		addrStr := getFromConfig(*cashierCfg)
		if addrStr != "" {
			// Handle both Ethereum and IoTeX addresses
			parsed, err := util.ParseEthAddress(addrStr)
			if err != nil {
				log.Fatalf("Failed to parse %s address from config: %v\n", name, err)
			}
			return parsed
		}
	}
	log.Fatalf("-%s or -cashier with %s in config is required\n", name, name)
	return common.Address{}
}
