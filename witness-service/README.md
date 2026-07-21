# Witness Service for ioTube

A **witness** watches a source chain's `TokenCashier`, signs each cross-chain
transfer it sees, and submits that signature to a **relayer** over gRPC. Once a
relayer has more than 2/3 of the registered witnesses' signatures for a transfer,
it submits the aggregated proof to the target chain's `TransferValidator`, which
mints/unlocks the funds. This directory builds and runs a witness.

This README covers **how the pieces fit together and how to build/run/upgrade**.
It intentionally contains no deployment-specific values (relayer endpoints, keys,
addresses) — those are supplied by the operator at deploy time.

## Layout

| Path | What it is |
|---|---|
| `docker-compose-witness.yml` | One container per route (per source chain the witness watches). All use the image `ghcr.io/iotubeproject/iotube-witness:${WITNESS_TAG:-latest}`. |
| `configs/` | The [`iotube-configs`](https://github.com/iotubeproject/iotube-configs) submodule — the config templates (see below). |
| `.env.template` | Operator-supplied environment (data dir, DB password, image tag). |
| `start_witness.sh` | First-time setup: copies config into the data dir, seeds `.env`, builds & starts the containers. |
| `upgrade.sh` | In-place image upgrade (pull + recreate, no config re-copy). |
| `Dockerfile.witness` | Builds the witness/relayer binary image. |

### config vs. secret

Each route reads **two** YAML files, merged at runtime:

- **`witness-config-<chain>-payload.yaml`** — the non-secret config: which cashier
  contracts to watch, and per cashier a `relayerURL: ":<port>"` where the **port**
  selects the target relayer. Shipped in the repo; the same for every operator.
- **`witness-config-<chain>.secret.yaml`** — the operator-filled part: your witness
  `privateKey` and the `relayerURL` **host** of the relayer you submit to. These
  ship as **empty placeholders** and are filled in at deploy.

So an operator only has to safeguard their **own private key** and point at a
relayer host; everything else comes from the shared config.

## Requirements

- Docker and Docker Compose v2 (`docker compose`).
- A witness private key — an EVM key that has been registered as a witness on-chain.

## First-time setup

```bash
# clone with the configs submodule
git clone --recurse-submodules <this-repo-url>
cd ioTube/witness-service

# guided setup: prompts for the data dir (default ~/iotex-witness),
# copies config, seeds .env, builds and starts the witness containers
./start_witness.sh
```

`start_witness.sh` seeds `.env` from `.env.template`; set at minimum the two
REQUIRED vars — `IOTEX_WITNESS` (absolute deploy path) and `DB_ROOT_PASSWORD`. If
either is left blank, `docker compose` fails fast and tells you which one to set
(rather than silently mounting empty dirs over the config files).

Then fill in your `privateKey` and relayer `relayerURL` host in the relevant
`<data-dir>/etc/witness-config-<chain>.secret.yaml` and restart.

## Upgrade in place

Pull the newest witness image (pinned by `WITNESS_TAG` in `.env`, default `latest`)
and recreate the witness containers — **without** re-copying config or touching the
database:

```bash
git pull --recurse-submodules && ./upgrade.sh
```

`upgrade.sh` pulls/recreates only the witness services (never the database), does a
non-destructive image prune, and migrates a legacy compose file if it still pins the
old local `witness:latest` tag. To pin or roll back a version, set
`WITNESS_TAG=sha-<commit>` in `<data-dir>/etc/.env` and re-run `./upgrade.sh`.

## Build the image

```bash
docker build -f witness-service/Dockerfile.witness -t iotube-witness .
```

Release CI publishes `ghcr.io/iotubeproject/iotube-witness:latest`
(`.github/workflows/docker-publish.yml`), which is what the compose file pulls by
default.

## Relayer API boundaries

The relayer HTTP gateway exposes only the read methods `check`, `list`,
`listnewtx`, and `lookup`. Witness submissions and heartbeats use gRPC directly.
Administrative methods (`Reset`, `Retry`, and `SubmitNewTX`) are local-only and
disabled unless `RELAYER_ADMIN_TOKEN` is set. An administrative gRPC caller must
run inside the relayer container, connect through loopback, and send
`authorization: Bearer <token>` metadata.

Persisted witness submissions must carry a valid 65-byte signature from a currently
active on-chain witness. For rolling-upgrade compatibility, legacy unsigned
pre-announcements receive a no-op acknowledgement but are never stored. If the
relayer cannot refresh the witness set, it rejects signed submissions until the
chain RPC recovers. Transfer proposals are stored by their signed transfer ID, and
only proposals with more than two-thirds of the current active witnesses enter the
settlement queue. On first startup after upgrading, legacy MySQL transfer tables are
migrated from the `(cashier, token, tidx)` primary key to the transfer ID; back up
the relayer database before deploying the upgrade.
