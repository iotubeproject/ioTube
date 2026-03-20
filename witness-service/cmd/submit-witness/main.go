// Copyright (c) 2024 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

// Configuration defines the relayer configuration
type Configuration struct {
	Chain              string            `json:"chain" yaml:"chain"`
	ClientURL          string            `json:"clientURL" yaml:"clientURL"`
	PrivateKey         string            `json:"privateKey" yaml:"privateKey"`
	ConfirmBlockNumber int               `json:"confirmBlockNumber" yaml:"confirmBlockNumber"`
	GrpcPort           int               `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort      int               `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	Interval           string            `json:"interval" yaml:"interval"`
	Database           DatabaseConfig    `json:"database" yaml:"database"`
	TransferTableName  string            `json:"transferTableName" yaml:"transferTableName"`
	WitnessTableName   string            `json:"witnessTableName" yaml:"witnessTableName"`
	Bonus              string            `json:"bonus" yaml:"bonus"`
	SlackWebHook       string            `json:"slackWebHook" yaml:"slackWebHook"`
	Validators         []ValidatorConfig `json:"validators" yaml:"validators"`
}

// DatabaseConfig defines database configuration
type DatabaseConfig struct {
	URI    string `json:"uri" yaml:"uri"`
	Driver string `json:"driver" yaml:"driver"`
}

// ValidatorConfig defines the configuration for a validator
type ValidatorConfig struct {
	Address   string   `json:"address" yaml:"address"`
	WithPayload bool   `json:"withPayload" yaml:"withPayload"`
	Cashiers  []CashierRef `json:"cashiers" yaml:"cashiers"`
}

// CashierRef references a cashier
type CashierRef struct {
	Address string `json:"address" yaml:"address"`
}

var (
	configFile    = flag.String("config", "", "path to relayer config file")
	secretFile    = flag.String("secret", "", "path to secret config file")
	validatorAddr = flag.String("validator", "", "validator contract address (overrides config)")
	cashierAddr   = flag.String("cashier", "", "cashier contract address")
	tokenAddr     = flag.String("token", "", "source token address")
	index         = flag.Uint("index", 0, "transfer index")
	senderAddr    = flag.String("sender", "", "sender address")
	recipientAddr = flag.String("recipient", "", "recipient address")
	amount        = flag.String("amount", "", "transfer amount")
	payloadHex    = flag.String("payload", "", "payload in hex (optional)")
	signatures    = flag.String("signatures", "", "comma-separated hex signatures")
	gasPrice      = flag.String("gas-price", "", "gas price in wei (optional)")
	gasLimit      = flag.Uint("gas-limit", 500000, "gas limit")
	dryRun        = flag.Bool("dry-run", false, "just print the transaction without sending")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -config <relayer-config.yaml> -secret <secret.yaml> [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nSubmits assembled witness signatures to the validator contract.")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintln(os.Stderr, "  # With signatures from multiple witnesses:")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-config relayer-config.yaml -secret secret.yaml \\")
		fmt.Fprintln(os.Stderr, "    -cashier 0x... -token 0x... -index 123 -sender 0x... -recipient 0x... -amount 1000000000000000000 \\")
		fmt.Fprintln(os.Stderr, `    -signatures "0xabc123...,0xdef456..."`)
	}
	flag.Parse()

	if *configFile == "" {
		log.Fatal("-config is required")
	}
	if *cashierAddr == "" {
		log.Fatal("-cashier is required")
	}
	if *tokenAddr == "" {
		log.Fatal("-token is required")
	}
	if *senderAddr == "" {
		log.Fatal("-sender is required")
	}
	if *recipientAddr == "" {
		log.Fatal("-recipient is required")
	}
	if *amount == "" {
		log.Fatal("-amount is required")
	}
	if *signatures == "" {
		log.Fatal("-signatures is required")
	}

	// Load config
	opts := []config.YAMLOption{config.Expand(os.LookupEnv), config.File(*configFile)}
	if *secretFile != "" {
		opts = append(opts, config.File(*secretFile))
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
	if pk, ok := os.LookupEnv("RELAYER_PRIVATE_KEY"); ok && cfg.PrivateKey == "" {
		cfg.PrivateKey = pk
	}

	if cfg.PrivateKey == "" {
		log.Fatal("Private key is required. Set it in the secret config file or RELAYER_PRIVATE_KEY env var.")
	}

	// Get validator address
	valAddr := *validatorAddr
	if valAddr == "" && len(cfg.Validators) > 0 {
		valAddr = cfg.Validators[0].Address
	}
	if valAddr == "" {
		log.Fatal("-validator or validator in config is required")
	}

	// Connect to client
	if cfg.ClientURL == "" {
		log.Fatal("clientURL is required in config")
	}
	client, err := ethclient.Dial(cfg.ClientURL)
	if err != nil {
		log.Fatalf("Failed to connect to client: %v\n", err)
	}
	defer client.Close()

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v\n", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v\n", err)
	}

	// Parse addresses
	cashier := common.HexToAddress(*cashierAddr)
	token := common.HexToAddress(*tokenAddr)
	sender := common.HexToAddress(*senderAddr)
	recipient := common.HexToAddress(*recipientAddr)
	validator := common.HexToAddress(valAddr)

	// Parse amount
	amountVal := new(big.Int)
	if _, ok := amountVal.SetString(*amount, 10); !ok {
		log.Fatalf("Invalid amount: %s\n", *amount)
	}

	// Parse payload
	var payload []byte
	if *payloadHex != "" {
		payload, err = hex.DecodeString(strings.TrimPrefix(*payloadHex, "0x"))
		if err != nil {
			log.Fatalf("Invalid payload hex: %v\n", err)
		}
	}

	// Parse and concatenate signatures
	sigList := strings.Split(*signatures, ",")
	concatSigs := []byte{}
	for i, sig := range sigList {
		sig = strings.TrimSpace(sig)
		sig = strings.TrimPrefix(sig, "0x")
		sigBytes, err := hex.DecodeString(sig)
		if err != nil {
			log.Fatalf("Invalid signature %d: %v\n", i, err)
		}
		// Ethereum signatures are 65 bytes (r, s, v)
		if len(sigBytes) != 65 {
			log.Fatalf("Signature %d has invalid length %d (expected 65)\n", i, len(sigBytes))
		}
		concatSigs = append(concatSigs, sigBytes...)
	}

	fmt.Println("=== Submitting to Validator ===")
	fmt.Printf("Validator: %s\n", validator.Hex())
	fmt.Printf("Cashier: %s\n", cashier.Hex())
	fmt.Printf("Token: %s\n", token.Hex())
	fmt.Printf("Index: %d\n", *index)
	fmt.Printf("Sender: %s\n", sender.Hex())
	fmt.Printf("Recipient: %s\n", recipient.Hex())
	fmt.Printf("Amount: %s\n", amountVal.String())
	if len(payload) > 0 {
		fmt.Printf("Payload: 0x%x\n", payload)
	}
	fmt.Printf("Signatures: %d (%d bytes total)\n", len(sigList), len(concatSigs))
	fmt.Printf("Relayer: %s\n", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

	if *dryRun {
		fmt.Println("\n[DRY RUN] Transaction not sent.")
		return
	}

	// Create transact opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v\n", err)
	}

	// Set gas price if specified
	if *gasPrice != "" {
		gp := new(big.Int)
		if _, ok := gp.SetString(*gasPrice, 10); !ok {
			log.Fatalf("Invalid gas price: %s\n", *gasPrice)
		}
		auth.GasPrice = gp
	} else {
		// Get suggested gas price
		gp, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("Failed to get gas price: %v\n", err)
		}
		auth.GasPrice = gp
	}
	auth.GasLimit = uint64(*gasLimit)

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v\n", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Create validator contract binding
	validatorContract, err := contract.NewTransferValidator(validator, client)
	if err != nil {
		log.Fatalf("Failed to create validator binding: %v\n", err)
	}

	// Submit transaction
	fmt.Println("\nSubmitting transaction...")
	var tx *types.Transaction

	// Determine if with payload based on config or flag
	withPayload := len(payload) > 0
	for _, v := range cfg.Validators {
		if v.Address == valAddr {
			withPayload = v.WithPayload
			break
		}
	}

	if withPayload {
		// Use TransferValidatorWithPayload
		validatorWithPayload, err := contract.NewTransferValidatorWithPayload(validator, client)
		if err != nil {
			log.Fatalf("Failed to create validator with payload binding: %v\n", err)
		}
		tx, err = validatorWithPayload.Submit(auth, cashier, token, big.NewInt(int64(*index)), sender, recipient, amountVal, concatSigs, payload)
	} else {
		tx, err = validatorContract.Submit(auth, cashier, token, big.NewInt(int64(*index)), sender, recipient, amountVal, concatSigs)
	}

	if err != nil {
		log.Fatalf("Failed to submit transaction: %v\n", err)
	}

	fmt.Println("=== Transaction Submitted ===")
	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Printf("Gas Price: %s wei\n", tx.GasPrice().String())
	fmt.Printf("Gas Limit: %d\n", tx.Gas())

	// Wait for receipt
	fmt.Println("\nWaiting for confirmation...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to get receipt: %v\n", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Println("=== Transaction Confirmed ===")
		fmt.Printf("Status: Success\n")
		fmt.Printf("Block: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

		// Parse Settled event
		for _, vLog := range receipt.Logs {
			if vLog.Address == validator {
				fmt.Printf("Log: %x\n", vLog.Topics)
			}
		}
	} else {
		fmt.Println("=== Transaction Failed ===")
		fmt.Printf("Status: Failed\n")
		os.Exit(1)
	}
}

func mustParseUint64(s string) uint64 {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalf("Invalid uint64: %s\n", s)
	}
	return val
}
