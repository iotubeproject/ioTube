# Witness Service for iotube

## How to deploy a new service?

0. install docker

1. go to the working directory of this repo, and then go into `witness-service` directory
2. create a path at `~/iotex-witness/etc`
3. copy .env.template `~/iotex-witness/etc/.env`, and set values
4. set `clientURL` in `witness-config-ethereum.secret.yaml`
5. run `./start_witness.sh` to build and start running services.

## How to setup a Bitcoin node for Bitcoin Bridge

1. Follow the instructions in https://bitcoin.org/en/full-node

2. connect ioTube service with the node by filling `rpcHost`, `rpcUser`, `rpcPass`, `enableTLS` fields in bitcoinConfig