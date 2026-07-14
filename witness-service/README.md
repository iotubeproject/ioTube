# Witness Service for iotube

## How to deploy a new service?

0. install docker

1. go to the working directory of this repo, and then go into `witness-service` directory
2. create a path at `~/iotex-witness/etc`
3. copy `.env.template` to `.env` and set values — at minimum the two REQUIRED
   vars `IOTEX_WITNESS` (absolute deploy path) and `DB_ROOT_PASSWORD`. If either is
   left blank, `docker compose` fails fast and tells you which one to set (rather
   than silently mounting empty dirs over the config files).
4. set `clientURL` in `witness-config-ethereum.secret.yaml`
5. run `./start_witness.sh` to build and start running services.
