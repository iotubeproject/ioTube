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

## Supported tokens

Fee-on-transfer and rebasing tokens are not supported.

## Network exposure

The witness database and management interfaces are host-local services.
Restrict access to them to loopback with the deployment firewall or security
group.

- MySQL: the compose file publishes it on `127.0.0.1:13306` only. `upgrade.sh`
  does not recreate the database container; after pulling this change, apply it
  once with `docker compose -f docker-compose-witness.yml up -d database`.
- Witness gRPC (`grpcPort`, default 9080) and HTTP gateway (`grpcProxyPort`,
  default 9081) listen on all interfaces with `network_mode: host`. Block them
  in the host firewall (or the cloud security group) for everything except
  loopback.
- Check with `ss -ltnp`: 13306 should listen on `127.0.0.1` only, and 9080/9081
  must not be reachable from another host.

Transfers that fail source validation (`ALERT not signing transfer ...` in the
log) remain pending for operator reconciliation.

## Relayer TLS

TLS and mTLS are available for witness-to-relayer gRPC. Configure peer
verification for connections that cross host boundaries. Existing transport
settings stay unchanged until configured; a plaintext connection logs
`WARNING connecting to the relayer without TLS` at startup.

| Side | Config | Env | Effect |
|---|---|---|---|
| relayer | `tls.certFile`, `tls.keyFile` | `RELAYER_TLS_CERT_FILE`, `RELAYER_TLS_KEY_FILE` | gRPC accepts TLS only |
| relayer | `tls.caFile` | `RELAYER_TLS_CLIENT_CA_FILE` | witnesses must present a client certificate signed by this CA (mTLS) |
| witness | `relayerTLS.caFile` | `RELAYER_TLS_CA_FILE` | the relayer's certificate must be signed by this CA |
| witness | `relayerTLS.certFile`, `relayerTLS.keyFile` | `RELAYER_TLS_CERT_FILE`, `RELAYER_TLS_KEY_FILE` | client certificate for mTLS |
| witness | `relayerTLS.serverName` | `RELAYER_TLS_SERVER_NAME` | name checked against the relayer's certificate (default: the host in `relayerURL`) |

There is no mixed mode: a witness with `relayerTLS` cannot reach a plaintext
relayer, and a plaintext witness cannot reach a TLS relayer. Distribute the CA
(and client certificates) to the operators first, then switch the relayer and
all witnesses of that relayer in one maintenance window. Enabling mTLS
(`tls.caFile`) at the same time or later is the same kind of switch. Admin gRPC
calls to a TLS relayer must use TLS as well, for example `grpcurl -cacert`.

## Signing with AWS KMS

Instead of `privateKey`, an EVM witness can sign with an AWS KMS key, so the
private key never exists outside KMS:

| Config | Env | Meaning |
|---|---|---|
| `kms.keyID` | `WITNESS_KMS_KEY_ID` | key ID, key ARN, or alias (`alias/...`) |
| `kms.region` | `WITNESS_KMS_REGION` | AWS region (default: the AWS SDK default chain) |

- The key must be `KeySpec=ECC_SECG_P256K1`, `KeyUsage=SIGN_VERIFY`. The witness
  address is derived from its public key at startup and logged; register that
  address in the WitnessList contracts.
- `privateKey` and `kms.keyID` are mutually exclusive. KMS signs EVM transfers
  only: a config with a `toSolana` cashier (Ed25519) refuses to start.
- Credentials come from the AWS SDK default chain. On EC2 use an instance role;
  no access keys on disk. With IMDSv2 and the witness in a Docker bridge
  network, set the instance metadata hop limit to 2
  (`aws ec2 modify-instance-metadata-options --http-tokens required --http-put-response-hop-limit 2`).
- The witness calls `GetPublicKey` once at startup and `Sign` once per transfer
  (signatures are kept in memory, so retries of the same transfer do not call
  KMS again). If KMS is unreachable the transfer is not signed and is retried on
  the next cycle.

Create the key in KMS itself (`Origin=AWS_KMS`), not from imported key material.
Suggested hardening:

- Key policy: only the witness instance role may `kms:Sign` and
  `kms:GetPublicKey`, optionally restricted with `aws:SourceVpce` to a KMS VPC
  endpoint. Nobody gets `kms:ScheduleKeyDeletion` or `kms:PutKeyPolicy` in
  normal operation; keep a break-glass admin role that requires MFA.
- To stop a compromised witness at once, `kms:DisableKey`; then remove the
  address from the WitnessList contracts.
- Alert on CloudTrail `Sign` calls from any principal other than the instance
  role, and on `DisableKey`, `ScheduleKeyDeletion`, `PutKeyPolicy`,
  `CreateGrant`.
- Deleting the key loses the witness address for good (it cannot be re-created);
  rotate by registering a new key's address first.

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
