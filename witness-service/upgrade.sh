#!/bin/bash
#
# Upgrade running witnesses in place.
#
# Pulls the witness image (pinned by WITNESS_TAG in .env, default: latest) and
# recreates the containers. It does NOT re-copy config files and does NOT run a
# destructive image prune, so it is safe to run repeatedly and will not clobber
# local edits.
#
# For a config/route change (e.g. adding a new chain), refresh the repo first:
#   git pull --recurse-submodules
#   ./start_witness.sh          # re-copies config templates, then (re)starts
#
# Usage:  ./upgrade.sh
#
set -euo pipefail

defaultdatadir="$HOME/iotex-witness"
IOTEX_WITNESS="${IOTEX_WITNESS:-$defaultdatadir}"

if [[ ! -f "$IOTEX_WITNESS/etc/docker-compose.yml" ]]; then
    echo "No compose file at $IOTEX_WITNESS/etc — run ./start_witness.sh first." >&2
    exit 1
fi

cd "$IOTEX_WITNESS/etc"

echo "Pulling witness image (WITNESS_TAG from .env, default: latest)..."
docker compose pull

echo "Recreating containers..."
docker compose up -d

echo "Removing dangling images..."
docker image prune -f

echo "Done. Current status:"
docker compose ps
