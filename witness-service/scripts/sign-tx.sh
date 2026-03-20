#!/bin/bash
# Sign a witness signature for a transfer transaction
# Usage: ./sign-tx.sh -config <config.yaml> -secret <secret.yaml> -tx <txhash>

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 -config <config.yaml> -secret <secret.yaml> -tx <txhash>"
    echo ""
    echo "Options:"
    echo "  -config    Path to witness config file"
    echo "  -secret    Path to secret config file with privateKey"
    echo "  -tx        Transaction hash to sign"
    echo ""
    echo "This script automatically:"
    echo "  - Fetches the transaction receipt from the blockchain"
    echo "  - Parses the Receipt event to extract transfer details"
    echo "  - Determines if destination is Ethereum or Solana"
    echo "  - Calls the appropriate signing script"
    exit 1
fi

exec "$SCRIPT_DIR/../bin/sign-tx" "$@"
