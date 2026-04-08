#!/bin/bash
# Sign a witness signature for Solana chain destinations (Ed25519)
# Usage: ./sign-to-solana.sh -config <config.yaml> -secret <secret.yaml> -cashier <id> -cotoken <address> -index <num> -sender <address> -recipient <pubkey> -amount <amount>

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SIGN_WITNESS="$SCRIPT_DIR/../bin/sign-witness"

# Load .env if IOTEX_WITNESS is set
if [[ -n "$IOTEX_WITNESS" && -f "$IOTEX_WITNESS/etc/.env" ]]; then
    set -a
    source "$IOTEX_WITNESS/etc/.env"
    set +a
fi

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 -config <config.yaml> -secret <secret.yaml> -cashier <id> -cotoken <address> -index <num> -sender <address> -recipient <pubkey> -amount <amount>"
    echo ""
    echo "Options:"
    echo "  -config    Path to witness config file"
    echo "  -secret    Path to secret config file with privateKey"
    echo "  -cashier   Cashier ID from config"
    echo "  -cotoken   Co-token address (Solana token mint)"
    echo "  -index     Transfer index"
    echo "  -sender    Sender address"
    echo "  -recipient Recipient Solana public key"
    echo "  -amount    Transfer amount"
    echo "  -payload   Optional payload in hex"
    exit 1
fi

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -config)    CONFIG="$2"; shift 2 ;;
        -secret)    SECRET="$2"; shift 2 ;;
        -cashier)   CASHIER="$2"; shift 2 ;;
        -token)     TOKEN="$2"; shift 2 ;;
        -cotoken)   COTOKEN="$2"; shift 2 ;;
        -index)     INDEX="$2"; shift 2 ;;
        -sender)    SENDER="$2"; shift 2 ;;
        -recipient) RECIPIENT="$2"; shift 2 ;;
        -amount)    AMOUNT="$2"; shift 2 ;;
        -payload)   PAYLOAD="$2"; shift 2 ;;
        *)          shift ;;
    esac
done

# Validate required args
[[ -z "$CONFIG" ]] && { echo "Error: -config is required"; exit 1; }
[[ -z "$SECRET" ]] && { echo "Error: -secret is required"; exit 1; }
[[ -z "$CASHIER" ]] && { echo "Error: -cashier is required"; exit 1; }
[[ -z "$COTOKEN" ]] && { echo "Error: -cotoken is required"; exit 1; }
[[ -z "$INDEX" ]] && { echo "Error: -index is required"; exit 1; }
[[ -z "$SENDER" ]] && { echo "Error: -sender is required"; exit 1; }
[[ -z "$RECIPIENT" ]] && { echo "Error: -recipient is required"; exit 1; }
[[ -z "$AMOUNT" ]] && { echo "Error: -amount is required"; exit 1; }

# Extract cashier and validator addresses from config
# First try to match by cashier ID (case-insensitive)
CASHIER_ADDR=$(grep -i -A 20 "id: \"$CASHIER\"" "$CONFIG" | grep -m1 "cashierContractAddress:" | awk '{print $2}' | tr -d '"')
VALIDATOR_ADDR=$(grep -i -A 20 "id: \"$CASHIER\"" "$CONFIG" | grep -m1 "validatorContractAddress:" | awk '{print $2}' | tr -d '"')
CASHIER_ID_FROM_CONFIG=""

# If not found by ID, try to match by cashier contract address (case-insensitive hex comparison)
if [[ -z "$CASHIER_ADDR" ]]; then
    # Normalize the input address (remove 0x prefix, convert to lowercase)
    CASHIER_NORM=$(echo "$CASHIER" | sed 's/^0x//' | tr '[:upper:]' '[:lower:]')

    # Get context around the matching cashierContractAddress (10 lines before, 10 lines after)
    CONTEXT=$(grep -i -B 10 -A 10 "cashierContractAddress.*$CASHIER_NORM" "$CONFIG")

    if [[ -n "$CONTEXT" ]]; then
        CASHIER_ADDR=$(echo "$CONTEXT" | grep -i "cashierContractAddress:" | head -1 | awk '{print $2}' | tr -d '"')
        VALIDATOR_ADDR=$(echo "$CONTEXT" | grep -i "validatorContractAddress:" | head -1 | awk '{print $2}' | tr -d '"')
        # Extract the cashier ID from context (look for "- id:" pattern)
        CASHIER_ID_FROM_CONFIG=$(echo "$CONTEXT" | grep -E "id:.*\"" | head -1 | sed 's/.*id:[[:space:]]*//' | tr -d '"')
    fi
fi

[[ -z "$CASHIER_ADDR" ]] && { echo "Error: Could not find cashier with ID or address: $CASHIER"; exit 1; }
[[ -z "$VALIDATOR_ADDR" ]] && { echo "Error: Could not find validatorContractAddress for $CASHIER"; exit 1; }

# Use extracted cashier ID if we found it by address
if [[ -n "$CASHIER_ID_FROM_CONFIG" ]]; then
    CASHIER="$CASHIER_ID_FROM_CONFIG"
fi

# Build and execute command
CMD="$SIGN_WITNESS -to-solana -config $CONFIG -secret $SECRET -cashier $CASHIER -cashier-address $CASHIER_ADDR -validator-address $VALIDATOR_ADDR -cotoken $COTOKEN -index $INDEX -sender $SENDER -recipient $RECIPIENT -amount $AMOUNT"
[[ -n "$TOKEN" ]] && CMD="$CMD -token $TOKEN"
[[ -n "$PAYLOAD" ]] && CMD="$CMD -payload $PAYLOAD"

exec $CMD
