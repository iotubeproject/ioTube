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
| `deploy/` | The release bundle: `docker-compose.yml` and `.env.example`. With the per-route configs it is published as `iotube-witness-<version>.tar.gz` on each release (see [Install](#install)). |
| `docker-compose-witness.yml` | Setup from a clone (`start_witness.sh`): one container per route (per source chain the witness watches). All use the image `ghcr.io/iotubeproject/iotube-witness:${WITNESS_TAG}`. |
| `configs/` | The [`iotube-configs`](https://github.com/iotubeproject/iotube-configs) submodule. `configs/witness/<route>.yaml` are the bundle's per-route configs; `configs/witness-config-*.yaml` are the templates `start_witness.sh` copies (see below). |
| `.env.template` | Operator-supplied environment (data dir, DB password, image version). |
| `start_witness.sh` | First-time setup: copies config into the data dir, seeds `.env`, pulls the image and starts the containers. |
| `upgrade.sh` | In-place version upgrade of all or selected routes (pull + recreate, no config re-copy). |
| `Dockerfile.witness` | Builds the witness/relayer binary image. |

### config vs. secret (setup from a clone)

With `start_witness.sh`, each route reads **two** YAML files, merged at runtime:

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
- A witness key: a private key, or an AWS KMS key (see
  [Signing with AWS KMS](#signing-with-aws-kms)). Its address must be
  registered in the WitnessList contracts before its signatures count (see
  [Joining as a new witness](#joining-as-a-new-witness)).
- RPC endpoints for the source chains you watch. They must serve `eth_getLogs`
  and receipts for at least the last few days. Some public endpoints, for
  example BSC's `bsc-dataseed`, do not do this reliably.

### Local database

The witness keeps its state in a local MySQL-compatible database (`witness-db`).
Choose the engine at first setup:

| Image | Memory, idle | When |
|---|---|---|
| `mysql:8` | about 420 MB | existing installs; hosts with memory to spare |
| `mariadb:11.4` | about 85 MB | new installs (the release bundle's default) |

The witness runs the same on both. The release bundle takes it from
`WITNESS_DB_IMAGE` in `.env`. `start_witness.sh` asks on a new install (or
takes `WITNESS_DB_IMAGE` from the environment) and records the choice as
`WITNESS_DB_IMAGE` in `<data-dir>/etc/.env`, the `.env` that compose uses; later
runs keep it, and an existing data directory keeps the engine that created it.
MariaDB cannot open a MySQL 8 data directory (or the other way round), so the
script refuses an image that does not match the existing data directory;
switching an existing install means `mysqldump` from the old database and
restoring into a new data directory.

## Install

New installs use the release bundle. Every operator value is in one `.env`
file; the route configs and the compose file are not edited.

```bash
VER=X.Y.Z   # the release to run
curl -fLO https://github.com/iotubeproject/ioTube/releases/download/v$VER/iotube-witness-$VER.tar.gz
curl -fLO https://github.com/iotubeproject/ioTube/releases/download/v$VER/iotube-witness-$VER.tar.gz.sha256
sha256sum -c iotube-witness-$VER.tar.gz.sha256
tar xzf iotube-witness-$VER.tar.gz          # creates ./iotube-witness/
cd iotube-witness
cp .env.example .env && chmod 600 .env      # fill in the REQUIRED values
docker compose up -d
```

The bundle holds:

| Path | What it is |
|---|---|
| `docker-compose.yml` | The local database and one service per route. `COMPOSE_PROFILES` in `.env` selects the routes. |
| `.env.example` | Every setting, with comments: image version, routes, database, relayer host, signing key or KMS key, RPC endpoints, start heights, sign policy. |
| `config/<route>.yaml` | Route configs. They take the RPC endpoint and start height from `.env`. |
| `signpolicy/` | Mounted into the containers for sign-policy certificate files. |
| `relayer-tls/` | Mounted into the containers for relayer TLS files (see [Relayer TLS](#relayer-tls)). |

The database data goes to `./data/mysql`; the database `witness` is created on
the first start.

Check that each route started:

```bash
docker compose ps
docker logs --since 5m bsc-to-iotex      # container name = route name
```

The log starts with `witness vX.Y.Z (commit ...)` and the witness keeps
scanning. The witness reads its configuration only at startup: after changing
`.env`, run `docker compose up -d --force-recreate`.

Compose leaves alone the services of profiles that are not selected. To stop
a route, remove it from `COMPOSE_PROFILES` and also remove its container,
otherwise it keeps running and signing:

```bash
docker compose rm -sf bsc-to-iotex
```

To upgrade, unpack the new release over the same directory. That replaces
`docker-compose.yml`, `.env.example` and `config/`, and leaves `.env` and
`data/` alone. Then:

```bash
# names of new settings to copy from .env.example; prints names only, so
# no key or password from .env reaches the terminal
comm -13 <(grep -oE '^[A-Z0-9_]+=' .env | sort) <(grep -oE '^[A-Z0-9_]+=' .env.example | sort)
sed -i 's/^WITNESS_TAG=.*/WITNESS_TAG=X.Y.Z/' .env
docker compose up -d --force-recreate
```

To roll back, unpack the previous release and set its `WITNESS_TAG`.

### Moving an existing install to the bundle

The bundle replaces the `start_witness.sh` stack on a host; the two must not
run side by side (same `witness-db` container, database port and route ports,
and the same key signing twice). The database name and tables are the same, so
the existing data directory carries over and each route resumes from its
stored height, as long as its start height is not raised.

The bundle uses one signer, one relayer and one sign policy for all routes
(the policy gates only the cashiers it lists, so a policy that gated one route
before still gates only that route). If the old secret files use different
keys, relayers or policy services per chain, keep the setup from a clone.

1. Unpack the bundle next to the old deploy dir (`$IOTEX_WITNESS`) and create
   `.env` from `.env.example`. Move every setting over. In the old install a
   value in the YAML file wins over the same setting in the environment (the
   old `.env` or the compose `environment:`), so take the YAML value where both
   are set.

   | Old setting (`$IOTEX_WITNESS/etc/`) | Bundle |
   |---|---|
   | `.env`: `WITNESS_TAG`, `WITNESS_DB_IMAGE` (unset means `mysql:8`), `DB_ROOT_PASSWORD` | the same names in `.env`; keep the database engine |
   | `*.secret.yaml` `privateKey` or `kms.keyID`/`kms.region`, or `WITNESS_PRIVATE_KEY`, `WITNESS_KMS_*` in the environment | `WITNESS_PRIVATE_KEY`, or `WITNESS_KMS_KEY_ID` and `WITNESS_KMS_REGION` |
   | `*.secret.yaml` `clientURL` of each chain | `IOTEX_RPC_URL`, `BSC_RPC_URL`, `POLYGON_RPC_URL` |
   | `relayerURL` in the YAML, or `RELAYER_URL` | `RELAYER_URL` |
   | `startBlockHeight` of each cashier in `*-payload.yaml` | `IOTEX_TO_BSC_START_HEIGHT`, `BSC_TO_IOTEX_START_HEIGHT`, `IOTEX_TO_POLYGON_START_HEIGHT`, `POLYGON_TO_IOTEX_START_HEIGHT`; the same value or lower, never higher, or transfers below it are skipped |
   | `signPolicy` in the YAML, or `WITNESS_SIGN_POLICY_*` in the compose `environment:` | `WITNESS_SIGN_POLICY_URL`, `_CASHIERS` (comma-separated), `_CLIENT_CERT`, `_CLIENT_KEY`, `_CA`; files to `signpolicy/`, paths `/etc/iotube-witness/signpolicy/<file>`. Without it the gated cashiers are signed without a check. |
   | `relayerTLS` in the YAML, or `RELAYER_TLS_*` | `RELAYER_TLS_CA_FILE`, `_CERT_FILE`, `_KEY_FILE`, `_SERVER_NAME`; files to `relayer-tls/`, paths `/etc/iotube-witness/relayer-tls/<file>` |
   | enabled cashiers of the running containers | `COMPOSE_PROFILES`, one route per enabled cashier: `iotex-to-bsc-payload` = `iotex-to-bsc`, `iotex-to-matic-payload` = `iotex-to-polygon` (both in `iotex-payload-witness`), `bsc-to-iotex-payload` = `bsc-to-iotex`, `matic-to-iotex-payload` = `polygon-to-iotex` |

   Copy secrets with an editor; do not print the files. Check that nothing is
   left: every non-empty setting in the old `*.yaml` files, the old `.env` and
   the compose `environment:` has a place in the new `.env`.

2. Stop the old stack and move the data, from the directory that holds
   `iotube-witness/`. Each step runs only if the one before succeeded: the old
   containers must be gone and the target must not exist yet. Keep both
   directories on the same filesystem, so the move is a rename.

   ```bash
   (cd "$IOTEX_WITNESS/etc" && docker compose down) \
     && names=$(docker ps -a --format '{{.Names}}') \
     && ! grep -xE 'witness-db|iotex-payload-witness|bsc-payload-witness|matic-payload-witness' <<<"$names" \
     && test ! -e iotube-witness/data/mysql \
     && mkdir -p iotube-witness/data \
     && sudo mv "$IOTEX_WITNESS/data/mysql" iotube-witness/data/mysql \
     && echo "data moved"
   ```

   Without `data moved`, fix what the output shows and run it again.

3. Start: `(cd iotube-witness && docker compose up -d)`, then check the logs as
   in [Install](#install); with a sign policy, each route logs
   `sign policy: <n> cashier(s) via <url>`.

To go back, stop every bundle container, including routes no longer in
`COMPOSE_PROFILES` (`--profile "*"`), then move the data back and start the
old stack. Run it, like step 2, from the directory that holds
`iotube-witness/`:

```bash
(cd iotube-witness && docker compose --profile "*" down) \
  && names=$(docker ps -a --format '{{.Names}}') \
  && ! grep -xE 'witness-db|iotex-to-bsc|bsc-to-iotex|iotex-to-polygon|polygon-to-iotex' <<<"$names" \
  && test ! -e "$IOTEX_WITNESS/data/mysql" \
  && sudo mv iotube-witness/data/mysql "$IOTEX_WITNESS/data/mysql" \
  && echo "data moved back"
```

## Setup from a clone (start_witness.sh)

Existing installs use this layout. New installs should use the release bundle
(see [Install](#install)).


```bash
# clone with the configs submodule
git clone --recurse-submodules <this-repo-url>
cd ioTube/witness-service

# guided setup: prompts for the data dir (default ~/iotex-witness) and the
# local database, copies config, seeds .env, pulls and starts the containers
./start_witness.sh
```

The setup asks which local database to use: MySQL 8 (the default) or MariaDB,
which needs about 85 MB of memory instead of about 420 MB and suits small hosts
(see [Local database](#local-database)). To set it without the prompt, for
example in automation:

```bash
WITNESS_DB_IMAGE=mariadb:11.4 ./start_witness.sh   # or mysql:8
```

`start_witness.sh` seeds `.env`; check the REQUIRED vars — `IOTEX_WITNESS`
(absolute deploy path), `DB_ROOT_PASSWORD` and `WITNESS_TAG` (the release to run,
for example `7.2.0`). If one is left blank, `docker compose` fails fast and tells
you which one to set (rather than silently mounting empty dirs over the config
files).

Then fill in your `privateKey` and relayer `relayerURL` host in the relevant
`<data-dir>/etc/witness-config-<chain>.secret.yaml` and restart.

## Upgrade in place (setup from a clone)

Run a release version, never `latest` or `X.Y`: those tags move with every
release, so the next time the containers are recreated after a pull, the running
version changes although the configuration did not. Releases are listed at
<https://github.com/iotubeproject/ioTube/releases>.

1. Read the release notes for config, database or compose changes.
2. Update the scripts and config templates to the release:
   `git fetch --tags && git checkout vX.Y.Z && git submodule update --init`
   (on a branch checkout, `git pull --recurse-submodules` also works).
3. Set the version in `<data-dir>/etc/.env`, for example `WITNESS_TAG=7.2.0`.
4. Upgrade one route first, check it, then the rest:

   ```bash
   ./upgrade.sh bsc-payload-witness   # one route
   ./upgrade.sh                       # all witness routes
   ```

5. Check that the container logs `witness vX.Y.Z (commit ...)` at startup and
   keeps scanning (`docker logs --since 5m <container>`). The relayer also
   records the version from the witness heartbeat (see [Versions](#versions)).

To roll back, set the previous `WITNESS_TAG` and run `./upgrade.sh` again.

`upgrade.sh` refuses `latest` or an unset `WITNESS_TAG`, pulls and recreates only
the selected witness services (never the database), removes only dangling images,
prints the version each upgraded image reports, and migrates a legacy compose file
that still names the old local `witness:latest` tag. It does not re-copy the
compose file: if you edited the `image:` lines by hand, set the version there.

## Joining as a new witness

The bridge operators give you the deployment values: the relayer host, which
routes to run, and a start height for each route. Then:

1. **Install** the release bundle (see [Install](#install)). In `.env`, set
   `COMPOSE_PROFILES` to the routes you run, `RELAYER_URL` to the relayer
   host, the start heights you were given, your RPC endpoints, and
   `WITNESS_PRIVATE_KEY` or `WITNESS_KMS_KEY_ID`.

2. **Start and check**: `docker compose up -d`, then `docker compose ps` and
   `docker logs --since 5m <route>`. The log starts with
   `witness vX.Y.Z (commit ...)`, and the witness keeps scanning without
   restarts.

3. **Shadow check.** Until your address is registered, the relayer refuses your
   submissions with
   `failed to submit transfer (...): ... PermissionDenied ... not an active witness`.
   The relayer verifies the signature before it checks the witness list, so this
   message means your key signs correctly. When it shows for each route, send
   the bridge operators your address and the logs. They also compare the
   transfers your witness recorded (its local database) with the relayer's
   records, to confirm that your witness reads the same transfers.

4. **Registration.** The operators add your address to the WitnessList contract
   on each target chain. A transfer settles with more than 2/3 of the active
   witnesses, `floor(2n/3) + 1` signatures out of `n`. After registration, the
   `not an active witness` message stops and your signatures count.

Operating:
- Upgrade to each release the operators announce (see [Install](#install)).
- With `n` witnesses, at most `n - floor(2n/3) - 1` may be down at the same
  time: one of six, one of four. Tell the operators before planned downtime.
- If your key may be compromised, tell the operators at once so they can remove
  the address. With KMS, also disable the key.
- Keep the host closed to inbound traffic (see
  [Network exposure](#network-exposure)).

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

With the release bundle, set the witness values in `.env` and put the files in
`relayer-tls/` (mounted at `/etc/iotube-witness/relayer-tls/`).

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

## Sign policy (optional)

A witness can ask an external service before it signs transfers of selected
cashiers, for example for a manual review on a newly opened route. For a listed
cashier, the witness POSTs the transfer to `<url>/v1/sign-check` and signs only
if the answer is `{"decision":"approve"}`. Any other answer, an error or a
timeout withholds the signature, and the transfer is retried on the next cycle.
Other cashiers are not affected.

With the release bundle, set the `WITNESS_SIGN_POLICY_*` values in `.env` and
put certificate files in `signpolicy/` (mounted at
`/etc/iotube-witness/signpolicy/`). With a setup from a clone, configure it in
the route's secret config (`witness-config-<chain>.secret.yaml`):

```yaml
signPolicy:
  url: "https://<policy service>"
  cashiers:                       # source cashier addresses to gate
    - "0x<cashier>"
  clientCert: "/etc/iotube-witness/signpolicy/client.crt"   # optional, mTLS
  clientKey: "/etc/iotube-witness/signpolicy/client.key"
  # ca: "/etc/iotube-witness/signpolicy/ca.pem"             # optional
```

| Config | Env | Meaning |
|---|---|---|
| `signPolicy.url` | `WITNESS_SIGN_POLICY_URL` | service base url |
| `signPolicy.cashiers` (YAML list) | `WITNESS_SIGN_POLICY_CASHIERS` (comma-separated) | source cashier addresses to gate |
| `signPolicy.clientCert`, `signPolicy.clientKey` | `WITNESS_SIGN_POLICY_CLIENT_CERT`, `WITNESS_SIGN_POLICY_CLIENT_KEY` | PEM client certificate for a service that requires mTLS |
| `signPolicy.ca` | `WITNESS_SIGN_POLICY_CA` | PEM roots that replace the system roots (for a self-hosted CA) |

- The compose file of a setup from a clone does not pass these environment
  variables into the containers. To use them there, add them to the service's
  `environment:`, and mount the certificate files, for example
  `- $IOTEX_WITNESS/etc/signpolicy:/etc/iotube-witness/signpolicy:ro`.
- At startup the witness logs `sign policy: <n> cashier(s) via <url>` (and the
  client certificate's name and expiry). Without that line, no policy is active.

- The url must be https, unless the host is loopback, a private address, or a
  single-label name such as a docker service (`http://witness-admin:8090`).
- Redirects are not followed. Plain http does not use a proxy from the
  environment.
- A url without cashiers, or cashiers without a url, is a startup error.
- Cashiers that send to Solana (`toSolana`) cannot be gated: listing one, or
  its previous address, stops the witness at startup.
- Recreate the containers after changing the policy
  (`docker compose up -d --force-recreate`).

## Versions

Releases are tagged `vX.Y.Z` on commits of `master`, and CI builds the images
from the tag. It refuses to publish a tag that is not on `master`.

- **major**: a major feature launch, the generations in the
  [version history](../README.md#version-history) (V7 = payload support).
  Decided by hand.
- **minor**: smaller features and hardening.
- **patch**: small bug fixes.

Run `7.2.0` or later. Images are tagged `X.Y.Z` (fixed), and `X.Y` and
`latest` (both move with new releases); deployments should use `X.Y.Z`.

`witness -version` and `relayer -version` print the build, and both log it at
startup:

```
v7.2.0 (commit <full sha>, clean, go1.27.3, built 2026-11-09T00:00:00Z)
```

The version is `git describe --tags --dirty`, so a build after a tag reads like
`v7.1.0-3-gabc1234`. `make` sets it from git; Docker builds take it from the
`VERSION`, `COMMIT`, `GIT_STATUS` and `BUILD_TIME` build args (the build context
has no `.git`), otherwise it reads `NoBuildInfo`.

Each witness sends its version and commit with its heartbeat (`StaleHeights`),
and the relayer keeps the latest per witness and route:

```sql
SELECT witness, cashier, version, commitID, tipHeight, updateTime
FROM witness_heartbeats ORDER BY witness, cashier;
```

Heartbeats are advisory, self-reported telemetry: they show which operators
have upgraded, and must not be used for authorization or build attestation.
Relayers add the `version` and
`commitID` columns at startup; a relayer without `ALTER` permission logs a
warning and records heartbeats without them.

## Build the image

```bash
docker build -f witness-service/Dockerfile.witness -t iotube-witness .
```

Release CI (`.github/workflows/docker-publish.yml`) publishes
`ghcr.io/iotubeproject/iotube-witness:X.Y.Z` for each release tag; the compose
file runs the version set by `WITNESS_TAG`.

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
