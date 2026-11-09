#!/bin/bash
#
# Upgrade running witnesses in place.
#
# Runs the witness image version set by WITNESS_TAG in <data-dir>/etc/.env (a
# release such as 7.2.0; moving tags such as `latest` or `7.2` are refused), pulls it and recreates the witness containers. It does NOT re-copy
# config files, never touches the database, and only removes dangling images,
# so it is safe to run repeatedly and will not clobber local edits.
#
# Usage:
#   ./upgrade.sh                         # all witness routes
#   ./upgrade.sh bsc-payload-witness     # only the named routes (upgrade one first)
#
# Roll back by setting the previous WITNESS_TAG and running it again.
#
# For a config/route change (e.g. adding a new chain), refresh the repo first:
#   git pull --recurse-submodules
#   ./start_witness.sh          # re-copies config templates, then (re)starts
#
set -euo pipefail

defaultdatadir="$HOME/iotex-witness"
IOTEX_WITNESS="${IOTEX_WITNESS:-$defaultdatadir}"
COMPOSE_FILE="$IOTEX_WITNESS/etc/docker-compose.yml"
WITNESS_IMAGE="ghcr.io/iotubeproject/iotube-witness"

if [[ ! -f "$COMPOSE_FILE" ]]; then
    echo "No compose file at $IOTEX_WITNESS/etc — run ./start_witness.sh first." >&2
    exit 1
fi

cd "$IOTEX_WITNESS/etc"

# Migrate legacy installs: witnesses set up before the ghcr-image change have a
# compose file that still names the old locally-retagged `witness:latest`, which
# `docker compose pull` cannot fetch. Rewrite just the image lines to the ghcr
# reference (idempotent; preserves any route comment/uncomment edits).
if grep -q '^[[:space:]]*image: witness:latest[[:space:]]*$' "$COMPOSE_FILE"; then
    echo "Migrating legacy 'witness:latest' image references to ghcr..."
    sed -i "s#image: witness:latest#image: ${WITNESS_IMAGE}:\${WITNESS_TAG:?WITNESS_TAG is unset - set it in .env to a release such as 7.2.0}#" "$COMPOSE_FILE"
fi

# Fail with compose's own message (e.g. WITNESS_TAG unset) before selecting services.
docker compose config -q

# Only pull/recreate the services that run the witness image. Selecting by image
# (rather than "every service except database") keeps a witness upgrade from also
# pulling/recreating the optional cron/backup helpers' floating Docker Hub images
# if an operator has uncommented them, and never touches the database.
mapfile -t ALL_SERVICES < <(docker compose config --services)
WITNESS_SERVICES=()
for svc in "${ALL_SERVICES[@]}"; do
    if docker compose config --images "$svc" 2>/dev/null | grep -q 'iotube-witness'; then
        WITNESS_SERVICES+=("$svc")
    fi
done
if [[ ${#WITNESS_SERVICES[@]} -eq 0 ]]; then
    echo "No witness services ($WITNESS_IMAGE) found in $COMPOSE_FILE." >&2
    exit 1
fi

# Optional arguments narrow the upgrade to the named witness services.
if [[ $# -gt 0 ]]; then
    for svc in "$@"; do
        if [[ ! " ${WITNESS_SERVICES[*]} " == *" $svc "* ]]; then
            echo "'$svc' is not a witness service in $COMPOSE_FILE (have: ${WITNESS_SERVICES[*]})." >&2
            exit 1
        fi
    done
    WITNESS_SERVICES=("$@")
fi

# A fixed image version: a digest, a release (7.2.0) or a commit build
# (sha-<commit>). Moving tags (latest, 7.2, master) and untagged images are refused.
isFixedImage() {
    local ref=$1 name tag
    [[ $ref == *@sha256:* ]] && return 0
    name=${ref##*/}                 # last path component; a registry port is before it
    [[ $name == *:* ]] || return 1  # untagged = latest
    tag=${name##*:}
    [[ ${#tag} -le 128 ]] && [[ $tag =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ || $tag =~ ^sha-[0-9a-f]{7,40}$ ]]
}

# Every selected service must run a fixed version.
for svc in "${WITNESS_SERVICES[@]}"; do
    image=$(docker compose config --images "$svc")
    if ! isFixedImage "$image"; then
        echo "$svc would run '$image', which is not a fixed version. Set WITNESS_TAG in $IOTEX_WITNESS/etc/.env to a release such as 7.2.0 (https://github.com/iotubeproject/ioTube/releases)." >&2
        exit 1
    fi
    echo "$svc -> $image"
done

echo "Pulling the witness image..."
docker compose pull "${WITNESS_SERVICES[@]}"

echo "Recreating witness containers..."
docker compose up -d --no-deps "${WITNESS_SERVICES[@]}"

echo "Removing dangling images..."
docker image prune -f

echo "Done. Running versions (also logged at startup):"
for svc in "${WITNESS_SERVICES[@]}"; do
    image=$(docker compose config --images "$svc")
    echo "  $svc: $(docker run --rm "$image" witness -version 2>/dev/null || echo 'unknown (image predates -version)')"
done
docker compose ps "${WITNESS_SERVICES[@]}"
echo "Check the logs, e.g. 'docker logs --since 5m ${WITNESS_SERVICES[0]}', before upgrading the next route."
