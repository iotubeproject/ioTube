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
	"strings"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	Chain                    string            `json:"chain" yaml:"chain"`
	ClientURL                string            `json:"clientURL" yaml:"clientURL"`
	PrivateKey               string            `json:"privateKey" yaml:"privateKey"`
	EthConfirmBlockNumber    int               `json:"ethConfirmBlockNumber" yaml:"ethConfirmBlockNumber"`
	EthGasPriceLimit         int64             `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
	GrpcPort                 int               `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort            int               `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	Interval                 string            `json:"interval" yaml:"interval"`
	Database                 DatabaseConfig    `json:"database" yaml:"database"`
	TransferTableName        string            `json:"transferTableName" yaml:"transferTableName"`
	NewTransactionTableName  string            `json:"newTransactionTableName" yaml:"newTransactionTableName"`
	WitnessTableName         string            `json:"witnessTableName" yaml:"witnessTableName"`
	Bonus                    string            `json:"bonus" yaml:"bonus"`
	SlackWebHook             string            `json:"slackWebHook" yaml:"slackWebHook"`
	SolanaConfig             interface{}       `json:"solanaConfig" yaml:"solanaConfig"`
	Unwrappers               interface{}       `json:"unwrappers" yaml:"unwrappers"`
	Validators               []ValidatorConfig `json:"validators" yaml:"validators"`
}

// DatabaseConfig defines database configuration
type DatabaseConfig struct {
	URI    string `json:"uri" yaml:"uri"`
	Driver string `json:"driver" yaml:"driver"`
}

// ValidatorConfig defines the configuration for a validator
type ValidatorConfig struct {
	Address     string       `json:"address" yaml:"address"`
	WithPayload bool         `json:"withPayload" yaml:"withPayload"`
	FromSolana  bool         `json:"fromSolana" yaml:"fromSolana"`
	Cashiers    []CashierRef `json:"cashiers" yaml:"cashiers"`
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
	tokenAddr     = flag.String("token", "", "co-token address on the destination chain")
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

	// Get cashier address from flag or config
	resolvedCashier, err := resolveCashier(*cashierAddr, cfg.Validators, valAddr)
	if err != nil {
		log.Fatal(err)
	}
	*cashierAddr = resolvedCashier

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
	concatSigs, err := parseSignatures(*signatures)
	if err != nil {
		log.Fatalf("Invalid signatures: %v\n", err)
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

	// Determine if with payload based on config
	withPayload := resolveWithPayload(cfg.Validators, valAddr)

	if *dryRun {
		callData, err := buildSubmitCallData(withPayload, cashier, token, big.NewInt(int64(*index)), sender, recipient, amountVal, concatSigs, payload)
		if err != nil {
			log.Fatalf("Failed to pack calldata: %v\n", err)
		}
		relayer := crypto.PubkeyToAddress(privateKey.PublicKey)
		estimatedGas, err := client.EstimateGas(context.Background(), goethereum.CallMsg{
			From: relayer,
			To:   &validator,
			Data: callData,
		})
		if err != nil {
			log.Fatalf("Gas estimation failed (params may be invalid): %v\n", err)
		}
		gp, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("Failed to get gas price: %v\n", err)
		}
		gasCost := new(big.Int).Mul(big.NewInt(int64(estimatedGas)), gp)
		fmt.Printf("\n[DRY RUN] Estimated gas: %d\n", estimatedGas)
		fmt.Printf("[DRY RUN] Gas price:     %s wei\n", gp.String())
		fmt.Printf("[DRY RUN] Estimated cost: %s wei\n", gasCost.String())
		fmt.Println("[DRY RUN] Transaction not sent.")
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

// buildSubmitCallData packs the ABI calldata for the submit function.
func buildSubmitCallData(withPayload bool, cashier, token common.Address, index *big.Int, sender, recipient common.Address, amount *big.Int, signatures []byte, payload []byte) ([]byte, error) {
	var (
		parsedABI abi.ABI
		err       error
	)
	if withPayload {
		parsedABI, err = abi.JSON(strings.NewReader(contract.TransferValidatorWithPayloadMetaData.ABI))
		if err != nil {
			return nil, err
		}
		return parsedABI.Pack("submit", cashier, token, index, sender, recipient, amount, signatures, payload)
	}
	parsedABI, err = abi.JSON(strings.NewReader(contract.TransferValidatorABI))
	if err != nil {
		return nil, err
	}
	return parsedABI.Pack("submit", cashier, token, index, sender, recipient, amount, signatures)
}

// parseSignatures parses a comma-separated list of hex-encoded 65-byte signatures
// and returns the concatenated bytes.
func parseSignatures(sigStr string) ([]byte, error) {
	sigList := strings.Split(sigStr, ",")
	concatSigs := make([]byte, 0, len(sigList)*65)
	for i, sig := range sigList {
		sig = strings.TrimSpace(sig)
		sig = strings.TrimPrefix(sig, "0x")
		sigBytes, err := hex.DecodeString(sig)
		if err != nil {
			return nil, fmt.Errorf("invalid signature %d: %v", i, err)
		}
		if len(sigBytes) != 65 {
			return nil, fmt.Errorf("signature %d has invalid length %d (expected 65)", i, len(sigBytes))
		}
		concatSigs = append(concatSigs, sigBytes...)
	}
	return concatSigs, nil
}

// resolveCashier returns the cashier address from the flag value if set, or falls back
// to the first cashier listed under the matching validator in config.
func resolveCashier(flagValue string, validators []ValidatorConfig, valAddr string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	for _, v := range validators {
		if strings.EqualFold(v.Address, valAddr) && len(v.Cashiers) > 0 {
			return v.Cashiers[0].Address, nil
		}
	}
	return "", fmt.Errorf("-cashier is required (or set cashiers in config validators)")
}

// resolveWithPayload returns true if the validator matching valAddr has withPayload set.
func resolveWithPayload(validators []ValidatorConfig, valAddr string) bool {
	for _, v := range validators {
		if strings.EqualFold(v.Address, valAddr) {
			return v.WithPayload
		}
	}
	return false
}
