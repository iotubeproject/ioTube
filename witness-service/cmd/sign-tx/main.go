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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/iotex-address/address"
	"go.uber.org/config"
)

// ReceiptEventData holds parsed event data
type ReceiptEventData struct {
	Token     common.Address
	ID        *big.Int
	Sender    common.Address
	Recipient string // Can be address or Solana pubkey
	Amount    *big.Int
	Fee       *big.Int
	Payload   []byte
	IsSolana  bool
}

// Configuration defines the configuration
type Configuration struct {
	Chain                 string          `json:"chain" yaml:"chain"`
	ClientURL             string          `json:"clientURL" yaml:"clientURL"`
	RelayerURL            string          `json:"relayerURL" yaml:"relayerURL"`
	PrivateKey            string          `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string          `json:"slackWebHook" yaml:"slackWebHook"`
	LarkWebHook           string          `json:"larkWebHook" yaml:"larkWebHook"`
	ConfirmBlockNumber    int             `json:"confirmBlockNumber" yaml:"confirmBlockNumber"`
	BatchSize             int             `json:"batchSize" yaml:"batchSize"`
	Interval              string          `json:"interval" yaml:"interval"`
	GrpcPort              int             `json:"grpcPort" yaml:"grpcPort"`
	GrpcProxyPort         int             `json:"grpcProxyPort" yaml:"grpcProxyPort"`
	DisableTransferSubmit bool            `json:"disableTransferSubmit" yaml:"disableTransferSubmit"`
	Database              DatabaseConfig  `json:"database" yaml:"database"`
	Cashiers              []CashierConfig `json:"cashiers" yaml:"cashiers"`
}

// DatabaseConfig defines database configuration
type DatabaseConfig struct {
	URI    string `json:"uri" yaml:"uri"`
	Driver string `json:"driver" yaml:"driver"`
}

// CashierConfig defines the configuration for a cashier
type CashierConfig struct {
	ID                         string      `json:"id" yaml:"id"`
	WithPayload                bool        `json:"withPayload" yaml:"withPayload"`
	RelayerURL                 string      `json:"relayerURL" yaml:"relayerURL"`
	CashierContractAddress     string      `json:"cashierContractAddress" yaml:"cashierContractAddress"`
	TokenSafeContractAddress   string      `json:"tokenSafeContractAddress" yaml:"tokenSafeContractAddress"`
	ValidatorContractAddress   string      `json:"validatorContractAddress" yaml:"validatorContractAddress"`
	PreviousCashierContractAddress string `json:"previousCashierContractAddress" yaml:"previousCashierContractAddress"`
	StartBlockHeight           int64       `json:"startBlockHeight" yaml:"startBlockHeight"`
	TransferTableName          string      `json:"transferTableName" yaml:"transferTableName"`
	TokenPairs                 []TokenPair `json:"tokenPairs" yaml:"tokenPairs"`
	ToSolana                   bool        `json:"toSolana" yaml:"toSolana"`
}

// TokenPair defines a token pair
type TokenPair struct {
	Token1 string `json:"token1" yaml:"token1"`
	Token2 string `json:"token2" yaml:"token2"`
}

var (
	configFile = flag.String("config", "", "path to witness config file")
	secretFile = flag.String("secret", "", "path to secret config file")
	txHash     = flag.String("tx", "", "transaction hash")
)

// Event topic hashes
var (
	// Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee)
	receiptEventTopic = crypto.Keccak256Hash([]byte("Receipt(address,uint256,address,address,uint256,uint256)"))
	// Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload)
	receiptWithPayloadEventTopic = crypto.Keccak256Hash([]byte("Receipt(address,uint256,address,address,uint256,uint256,bytes)"))
	// Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee, bytes payload)
	receiptSolanaEventTopic = crypto.Keccak256Hash([]byte("Receipt(address,uint256,address,string,uint256,uint256,bytes)"))
)

// parseAddressToHex converts both io1... and 0x... address formats to lowercase hex
func parseAddressToHex(addr string) string {
	addr = strings.TrimSpace(addr)
	if strings.HasPrefix(addr, "io1") {
		// IoTeX native address format - convert to hex
		ioAddr, err := address.FromString(addr)
		if err != nil {
			return strings.ToLower(addr)
		}
		return strings.ToLower(common.BytesToAddress(ioAddr.Bytes()).Hex())
	}
	// Already hex format
	return strings.ToLower(addr)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -config <config.yaml> -secret <secret.yaml> -tx <txhash>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nSigns a witness signature for a transfer transaction.")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *txHash == "" {
		log.Fatal("-tx is required")
	}
	if *configFile == "" {
		log.Fatal("-config is required")
	}
	if *secretFile == "" {
		log.Fatal("-secret is required")
	}

	// Load config
	opts := []config.YAMLOption{config.Expand(os.LookupEnv), config.File(*configFile), config.File(*secretFile)}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	// Override with environment variables if set
	if pk, ok := os.LookupEnv("WITNESS_PRIVATE_KEY"); ok && cfg.PrivateKey == "" {
		cfg.PrivateKey = pk
	}
	if url, ok := os.LookupEnv("WITNESS_CLIENT_URL"); ok && cfg.ClientURL == "" {
		cfg.ClientURL = url
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

	// Get transaction receipt
	txHash := common.HexToHash(*txHash)
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v\n", err)
	}

	// Find the Receipt event in logs
	var eventData *ReceiptEventData
	var cashierCfg *CashierConfig

	for _, logEntry := range receipt.Logs {
		// Check if this log is from one of our cashiers
		logAddrHex := strings.ToLower(logEntry.Address.Hex())
		for i := range cfg.Cashiers {
			cc := &cfg.Cashiers[i]
			cashierAddrHex := parseAddressToHex(cc.CashierContractAddress)
			if logAddrHex == cashierAddrHex {
				// Try to parse as Receipt event
				data, err := parseReceiptEvent(logEntry, cc.ToSolana)
				if err == nil && data != nil {
					eventData = data
					cashierCfg = cc
					break
				}
			}
		}
		if eventData != nil {
			break
		}
	}

	if eventData == nil {
		log.Fatal("Could not find Receipt event in transaction logs")
	}

	// Find co-token from token pairs
	var coToken string
	tokenHex := strings.ToLower(eventData.Token.Hex())
	for _, pair := range cashierCfg.TokenPairs {
		// Parse token1 address to normalized hex format
		token1Hex := parseAddressToHex(pair.Token1)
		if token1Hex == tokenHex {
			coToken = parseAddressToHex(pair.Token2)
			break
		}
	}
	if coToken == "" {
		log.Fatalf("Could not find token pair for token %s in cashier %s", eventData.Token.Hex(), cashierCfg.ID)
	}

	// Determine which script to call
	scriptDir := getScriptDir()
	var cmd *exec.Cmd
	args := []string{
		"-config", *configFile,
		"-secret", *secretFile,
		"-cashier", cashierCfg.ID,
		"-token", eventData.Token.Hex(),
		"-cotoken", coToken,
		"-index", eventData.ID.String(),
		"-sender", eventData.Sender.Hex(),
		"-recipient", eventData.Recipient,
		"-amount", eventData.Amount.String(),
	}

	if len(eventData.Payload) > 0 {
		args = append(args, "-payload", hex.EncodeToString(eventData.Payload))
	}

	if cashierCfg.ToSolana {
		fmt.Println("Calling sign-to-solana.sh...")
		cmd = exec.Command(filepath.Join(scriptDir, "sign-to-solana.sh"), args...)
	} else {
		fmt.Println("Calling sign-to-eth.sh...")
		cmd = exec.Command(filepath.Join(scriptDir, "sign-to-eth.sh"), args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run signing script: %v\n", err)
	}
}

func parseReceiptEvent(logEntry *types.Log, isSolana bool) (*ReceiptEventData, error) {
	// Receipt event has topics: [eventSig, token (indexed), id (indexed)]
	// and data: [sender, recipient, amount, fee, payload?]

	if len(logEntry.Topics) < 3 {
		return nil, fmt.Errorf("not enough topics")
	}

	// Determine event type based on topic
	topic := logEntry.Topics[0]
	var isWithPayload bool

	if isSolana {
		if topic != receiptSolanaEventTopic {
			return nil, fmt.Errorf("event signature mismatch: got %x, expected %x (solana)", topic, receiptSolanaEventTopic)
		}
	} else {
		if topic == receiptEventTopic {
			isWithPayload = false
		} else if topic == receiptWithPayloadEventTopic {
			isWithPayload = true
		} else {
			return nil, fmt.Errorf("event signature mismatch: got %x, expected %x or %x", topic, receiptEventTopic, receiptWithPayloadEventTopic)
		}
	}

	token := common.BytesToAddress(logEntry.Topics[1].Bytes())
	id := new(big.Int).SetBytes(logEntry.Topics[2].Bytes())

	// Parse non-indexed data from the event
	// For Ethereum (no payload): sender (32), recipient (32), amount (32), fee (32)
	// For Ethereum (with payload): sender (32), recipient (32), amount (32), fee (32), payload offset (32), payload len (32), payload
	// For Solana: sender (32), recipient offset (32), amount (32), fee (32), payload offset (32)

	data := logEntry.Data
	if len(data) < 128 {
		return nil, fmt.Errorf("data too short: %d bytes", len(data))
	}

	sender := common.BytesToAddress(data[0:32])
	amount := new(big.Int).SetBytes(data[64:96])
	fee := new(big.Int).SetBytes(data[96:128])

	var recipient string
	var payload []byte

	if isSolana {
		// Recipient is a string (dynamic)
		recipientOffset := new(big.Int).SetBytes(data[32:64]).Int64()
		if int(recipientOffset)+32 > len(data) {
			return nil, fmt.Errorf("invalid recipient offset")
		}
		recipientLen := new(big.Int).SetBytes(data[recipientOffset : recipientOffset+32]).Int64()
		if int(recipientOffset+32+recipientLen) > len(data) {
			return nil, fmt.Errorf("invalid recipient length")
		}
		recipient = string(data[recipientOffset+32 : recipientOffset+32+recipientLen])

		// Payload (if present)
		if len(data) > 160 {
			payloadOffset := new(big.Int).SetBytes(data[128:160]).Int64()
			if int(payloadOffset)+32 <= len(data) {
				payloadLen := new(big.Int).SetBytes(data[payloadOffset : payloadOffset+32]).Int64()
				if int(payloadOffset+32+payloadLen) <= len(data) {
					payload = data[payloadOffset+32 : payloadOffset+32+payloadLen]
				}
			}
		}
	} else if isWithPayload {
		// Recipient is an address (fixed)
		recipient = common.BytesToAddress(data[32:64]).Hex()

		// Payload is at offset in data[128:160]
		if len(data) >= 160 {
			payloadOffset := new(big.Int).SetBytes(data[128:160]).Int64()
			if int(payloadOffset)+32 <= len(data) {
				payloadLen := new(big.Int).SetBytes(data[payloadOffset : payloadOffset+32]).Int64()
				if int(payloadOffset+32+payloadLen) <= len(data) {
					payload = data[payloadOffset+32 : payloadOffset+32+payloadLen]
				}
			}
		}
	} else {
		// Recipient is an address (no payload)
		recipient = common.BytesToAddress(data[32:64]).Hex()
	}

	return &ReceiptEventData{
		Token:     token,
		ID:        id,
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Payload:   payload,
		IsSolana:  isSolana,
	}, nil
}

func getScriptDir() string {
	execPath, err := os.Executable()
	if err != nil {
		return "."
	}
	// Binary is in bin/, scripts are in scripts/
	binDir := filepath.Dir(execPath)
	scriptDir := filepath.Join(binDir, "..", "scripts")
	// Check if scripts directory exists
	if _, err := os.Stat(scriptDir); err == nil {
		return scriptDir
	}
	// Fallback to same directory as binary
	return binDir
}
