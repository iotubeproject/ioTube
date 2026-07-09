#!/bin/bash
#
# Upgrade running witnesses in place.
#
# Pulls the witness image (pinned by WITNESS_TAG in .env, default: latest) and
# recreates the witness containers. It does NOT re-copy config files and does
# NOT run a destructive image prune, so it is safe to run repeatedly and will
# not clobber local edits.
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
COMPOSE_FILE="$IOTEX_WITNESS/etc/docker-compose.yml"

if [[ ! -f "$COMPOSE_FILE" ]]; then
    echo "No compose file at $IOTEX_WITNESS/etc — run ./start_witness.sh first." >&2
    exit 1
fi

cd "$IOTEX_WITNESS/etc"

# Migrate legacy installs: witnesses set up before the ghcr-image change have a
# compose file that still names the old locally-retagged `witness:latest`, which
# `docker compose pull` cannot fetch. Rewrite just the image lines to the pinned
# ghcr reference (idempotent; preserves any route comment/uncomment edits).
if grep -q '^[[:space:]]*image: witness:latest[[:space:]]*$' "$COMPOSE_FILE"; then
    echo "Migrating legacy 'witness:latest' image references to ghcr..."
    sed -i 's#image: witness:latest#image: ghcr.io/iotubeproject/iotube-witness:${WITNESS_TAG:-latest}#' "$COMPOSE_FILE"
fi

# Only pull/recreate the witness services, never the database: an unqualified
# `docker compose pull` would also pull the database service's floating `mysql:8`
# tag, and the subsequent `up -d` would roll the production DB to a new digest.
mapfile -t WITNESS_SERVICES < <(docker compose config --services | grep -v -x 'database')
if [[ ${#WITNESS_SERVICES[@]} -eq 0 ]]; then
    echo "No witness services found in $COMPOSE_FILE." >&2
    exit 1
fi

echo "Pulling witness image (WITNESS_TAG from .env, default: latest)..."
docker compose pull "${WITNESS_SERVICES[@]}"

echo "Recreating witness containers..."
docker compose up -d "${WITNESS_SERVICES[@]}"

echo "Removing dangling images..."
docker image prune -f

echo "Done. Current status:"
docker compose ps
