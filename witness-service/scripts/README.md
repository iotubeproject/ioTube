# Witness Signing Tools

A set of tools for generating and submitting witness signatures for cross-chain transfers on ioTube.

## Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         Witness Signing Workflow                          │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   Transaction ──▶ sign-tx.sh ──▶ Signature 1                            │
│       │                           Signature 2                            │
│       │                           Signature N                            │
│       │                                  │                               │
│       │                                  ▼                               │
│       │                         submit-witness.sh                        │
│       │                                  │                               │
│       ▼                                  ▼                               │
│   Blockchain ◀──────────────────▶ Validator Contract                    │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

## Installation

```bash
# Build all binaries
GOTOOLCHAIN=go1.23.7 go build -o sign-witness ./cmd/sign-witness/
GOTOOLCHAIN=go1.23.7 go build -o sign-tx ./cmd/sign-tx/
GOTOOLCHAIN=go1.23.7 go build -o submit-witness ./cmd/submit-witness/
```

## Quick Start

### 1. Sign by Transaction Hash (Recommended)

The easiest way - just provide a transaction hash:

```bash
./scripts/sign-tx.sh \
  -config configs/witness-config-iotex-testnet.yaml \
  -secret secret.yaml \
  -tx 0x1234567890abcdef...
```

This automatically:
- Fetches the transaction receipt from the blockchain
- Parses the `Receipt` event to extract transfer details
- Determines if destination is Ethereum or Solana
- Generates the appropriate signature

### 2. Sign Manually

#### For Ethereum Destinations

```bash
./scripts/sign-to-eth.sh \
  -config configs/witness-config-iotex-testnet.yaml \
  -secret secret.yaml \
  -cashier iotex-testnet-to-bsc-testnet \
  -token 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
  -index 123 \
  -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
  -recipient 0x1234567890123456789012345678901234567890 \
  -amount 1000000000000000000
```

#### For Solana Destinations (Ed25519)

```bash
./scripts/sign-to-solana.sh \
  -config configs/witness-config-iotex-solana.yaml \
  -secret secret.yaml \
  -cashier iotex-to-solana \
  -token io158elyywekvljpp4gqzzzsdkk0ufehxn6g3aa0u \
  -index 456 \
  -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
  -recipient 9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM \
  -amount 1000000
```

### 3. Submit Signatures to Validator

After collecting signatures from multiple witnesses:

```bash
./scripts/submit-witness.sh \
  -config configs/relayer-config-iotex-testnet.yaml \
  -secret relayer-secret.yaml \
  -cashier 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
  -token 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
  -index 123 \
  -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
  -recipient 0x1234567890123456789012345678901234567890 \
  -amount 1000000000000000000 \
  -signatures "0xabc123...,0xdef456...,0x789..."
```

## Configuration Files

### Witness Config (sign-*.sh)

```yaml
chain: "iotex-testnet"
clientURL: "https://babel-api.testnet.iotex.io"
cashiers:
  - id: "iotex-testnet-to-bsc-testnet"
    cashierContractAddress: "0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3"
    validatorContractAddress: "0xEf503971Aec1BC3cF3D896742Fa82975dCcB3162"
    toSolana: false  # Set true for Solana destinations
    tokenPairs:
      - token1: "0xToken1Address"
        token2: "0xToken2Address"
```

### Secret Config

```yaml
privateKey: "your-private-key-hex-without-0x-prefix"
```

### Relayer Config (submit-witness.sh)

```yaml
chain: "iotex-testnet"
clientURL: "https://babel-api.testnet.iotex.io"
privateKey: "relayer-private-key"
validators:
  - address: "0xValidatorContractAddress"
    withPayload: true
    cashiers:
      - address: "0xCashierContractAddress"
```

## Script Reference

| Script | Purpose |
|--------|---------|
| `sign-tx.sh` | Sign by transaction hash (auto-detects all details) |
| `sign-to-eth.sh` | Sign for Ethereum destinations (secp256k1) |
| `sign-to-solana.sh` | Sign for Solana destinations (Ed25519) |
| `submit-witness.sh` | Submit assembled signatures to validator |

## Output Format

All signing scripts output:

```
=== Witness Signature ===
Transfer ID: 0x...
Witness Public Key: 0x...
Signature: 0x...

=== Details ===
Validator: 0x...
Cashier: 0x...
CoToken: 0x...
Index: 123
Sender: 0x...
Recipient: 0x...
Amount: 1000000000000000000
To Solana: false
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WITNESS_PRIVATE_KEY` | Private key (alternative to secret file) |
| `RELAYER_PRIVATE_KEY` | Relayer private key for submit-witness |

## Signature Types

- **Ethereum (secp256k1)**: 65 bytes (r: 32, s: 32, v: 1)
- **Solana (Ed25519)**: 64 bytes

## Testing

```bash
# Run test suite
./scripts/test-scripts.sh
```

## Troubleshooting

### "Could not find cashierContractAddress"
- Check that the cashier ID matches exactly (case-sensitive)
- Verify the config file path is correct

### "Could not find token2 for token X"
- Ensure the token address is in the `tokenPairs` list
- Check that addresses match exactly (case-insensitive for hex)

### "Invalid signature length"
- Ethereum signatures must be 65 bytes
- Solana signatures must be 64 bytes
- Ensure you're using the correct script for the destination chain
