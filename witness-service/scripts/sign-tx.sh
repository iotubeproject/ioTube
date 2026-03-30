#!/bin/bash
# Sign a witness signature for a transfer transaction
# Usage: ./sign-tx.sh -config <config.yaml> -secret <secret.yaml> -tx <txhash>

set -e

# Colour codes
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_IOTEX_WITNESS="$HOME/iotex-witness"

# Determine IOTEX_WITNESS directory
function determineIotexWitness() {
    if [[ -z "$IOTEX_WITNESS" ]]; then
        echo "Input IOTEX_WITNESS directory. This is where your .env file is stored."
        echo -e "${RED}Input your directory \$IOTEX_WITNESS!!!${NC}"
        read -p "Input your \$IOTEX_WITNESS [default: $DEFAULT_IOTEX_WITNESS]: " inputdir
        IOTEX_WITNESS="${inputdir:-$DEFAULT_IOTEX_WITNESS}"
    fi
    echo -e "${YELLOW}IOTEX_WITNESS: $IOTEX_WITNESS${NC}"
}

function loadEnvFile() {
    local envFile="$IOTEX_WITNESS/etc/.env"
    if [[ -f "$envFile" ]]; then
        echo -e "${YELLOW}Loading environment from $envFile${NC}"
        # Export all variables from .env file
        set -a
        source "$envFile"
        set +a
    else
        echo -e "${YELLOW}No .env file found at $envFile${NC}"
    fi
}

# Parse command line arguments
CONFIG_FILE=""
SECRET_FILE=""
TX_HASH=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -secret)
            SECRET_FILE="$2"
            shift 2
            ;;
        -tx)
            TX_HASH="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 -config <config.yaml> -secret <secret.yaml> -tx <txhash>"
            echo ""
            echo "Options:"
            echo "  -config    Path to witness config file"
            echo "  -secret    Path to secret config file with privateKey"
            echo "  -tx        Transaction hash to sign"
            echo ""
            echo "Environment variables (can be set in \$IOTEX_WITNESS/etc/.env):"
            echo "  IOTEX_WITNESS        Base directory for .env file"
            echo "  WITNESS_PRIVATE_KEY  Private key for signing (fallback)"
            echo "  WITNESS_CLIENT_URL   Ethereum RPC URL (fallback)"
            echo ""
            echo "This script automatically:"
            echo "  - Prompts for IOTEX_WITNESS if not set"
            echo "  - Loads environment from \$IOTEX_WITNESS/etc/.env if it exists"
            echo "  - Fetches the transaction receipt from the blockchain"
            echo "  - Parses the Receipt event to extract transfer details"
            echo "  - Determines if destination is Ethereum or Solana"
            echo "  - Calls the appropriate signing script"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Check required parameters
if [[ -z "$TX_HASH" ]]; then
    echo -e "${RED}Error: -tx is required${NC}"
    exit 1
fi
if [[ -z "$CONFIG_FILE" ]]; then
    echo -e "${RED}Error: -config is required${NC}"
    exit 1
fi
if [[ -z "$SECRET_FILE" ]]; then
    echo -e "${RED}Error: -secret is required${NC}"
    exit 1
fi

# Determine IOTEX_WITNESS and load .env
determineIotexWitness
loadEnvFile

# Execute sign-tx binary
exec "$SCRIPT_DIR/../bin/sign-tx" \
    -config "$CONFIG_FILE" \
    -secret "$SECRET_FILE" \
    -tx "$TX_HASH"
