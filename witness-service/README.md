# iotex-tube

[![CircleCI](https://circleci.com/gh/iotexproject/iotex-antenna-go.svg?style=svg)](https://circleci.com/gh/iotexproject/iotex-tube)
[![Go version](https://img.shields.io/badge/go-1.11.5-blue.svg)](https://github.com/moovweb/gvm)
[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Tube is a bi-directional service that allows our users to swap the IoTeX mainnet native token to the [IoTeX Network ERC20 token](https://etherscan.io/token/0x6fb3e0a217407efff7ca062d46c26e5d60a14d69), or vice versa.

Click [IoTeX Tube docs](https://github.com/iotexproject/iotex-bootstrap/blob/master/tube/tube.md) for detailed documentation of the tube service.

## Getting started

### Minimum requirements

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |

## Running the service with docker

### Prepare docker image

```make docker```

### Prepare config file

Copy `service.yaml` to current directory, and fill in the following fields:
* db.url
* iotex.privateKey
* ethereum.client
* ethereum.privateKey

### Run service in docker

```
docker run -d --restart on-failure -name witness \
         -p 8080:8080 \
         -v=service.yaml:/etc/iotube-witness/service.yaml:ro \
         iotex/iotube-witness:v0.1 \
         iotube-witness \
         -config=/etc/iotube-witness/service.yaml
```

## Run service in Heroku (not ready yet)
The service is currently deployed on Heroku, with the following env variables:

| Variable | Description |
|----------|-------------|
| ERC20_CONTRACT_ADDRESS | IOTX network ERC20 contract: 0x6fB3e0A217407EFFf7Ca062D46c26E5d60a14d69 |
| ETH_CLIENT_URLS | Ethereum network endpoint: https://mainnet.infura.io/v3/b355cae6fafc4302b106b937ee6c15af |
| ETH_FAUCET_ADDRESS | swap contract deployed on Ethereum: 0x450cab2535d57ce9df625297d796aee266611728 |
| ETH_GAS_PRICE_LIMIT | the maximum Ethereum gas price we set: 20000000000 |
| ETH_MAINNET_PK | the private key of wallet sending out Ethereum |
| EXCHANGER_DIRECTION | "e2n" or "n2e" |
| FAUCET_BALANCE_THRESHOLD | 20000000000000000000000: once wallet balance drop below this amount, alert will be sent to slack channel |
| FAUCET_WEB_HOOK | the slack alert channel webhook URL: https://hooks.slack.com/services/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 |
| IO_ENDPOINT | IoTeX mainnet endpoint: api.iotex.one:443 |
| IO_FAUCET_ADDRESS | swap contract deployed on mainnet: io1p99pprm79rftj4r6kenfjcp8jkp6zc6mytuah5 |
| IO_MAINNET_PK | private key of wallet sending out IoTeX mainnet token |
| TABLE_NAME | MySql table name of swap record |
| TRANSFER_DATABASE_URL | MySql table URL |
