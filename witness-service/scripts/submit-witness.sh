#!/bin/bash
# Assemble witness signatures and submit to validator contract
# Usage: ./submit-witness.sh -config <relayer-config.yaml> -token <addr> -index <num> -sender <addr> -recipient <addr> -amount <num> -signatures <sig1,sig2,...>

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -eq 0 ]]; then
    cat << 'EOF'
Usage: ./submit-witness.sh -config <relayer-config.yaml> [options]

Assembles witness signatures and submits them to the validator contract.

Required Options:
  -config       Path to relayer config file (contains validator address, clientURL, privateKey)
  -token        Co-token address on the destination chain (CoToken from sign-witness output)
  -index        Transfer index
  -sender       Sender address
  -recipient    Recipient address
  -amount       Transfer amount (in wei)
  -signatures   Comma-separated hex signatures (from sign-witness output)

Optional Options:
  -secret       Path to secret config file (if privateKey is not in -config)
  -cashier      Cashier contract address (default: first cashier in config)
  -validator    Validator contract address (overrides config)
  -payload      Payload in hex (for transfers with payload)
  -gas-price    Gas price in wei
  -gas-limit    Gas limit (default: 500000)
  -dry-run      Print transaction without sending

Example:
  # After collecting signatures from witnesses:
  ./submit-witness.sh \
    -config configs/relayer-config-iotex-testnet.yaml \
    -cashier 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -token 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -index 123 \
    -sender 0x1234567890123456789012345678901234567890 \
    -recipient 0xabcdefabcdefabcdefabcdefabcdefabcdef \
    -amount 1000000000000000000 \
    -signatures "0xabc123...,0xdef456..."

  # Dry run to preview:
  ./submit-witness.sh ... -dry-run

EOF
    exit 1
fi

exec "$SCRIPT_DIR/../bin/submit-witness" "$@"
