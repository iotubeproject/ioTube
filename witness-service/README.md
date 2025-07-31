# Witness Service for iotube

## How to deploy a new service?

0. install docker and git clone iotube repo.

1. create a working path at `~/iotex-witness/` with `~/iotex-witness/etc`. 
2. copy `ioTube/witness-service/.env.template` to `~/iotex-witness/etc/.env`, and set `RELAYER_URL` and `WITNESS_PRIVATE_KEY`
3. `cd ioTube/witness-service/ && ./start_witness.sh` to build and start service.
4. update `witness-config-[chainname].secret.yaml`  （currently eth, bsc, matic, iotex are in effective.):  `clientURL` to your own RPC
5. update `witness-config-[chainname].yaml`: `startBlockHeight` to the tip of the blockchain. (only fetch new data onwards.)

