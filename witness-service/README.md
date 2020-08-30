## Running the service with docker-compose

### Prepare Config

Copy `service.yaml` to current directory, and fill in the following fields:
* iotex.privateKey
* ethereum.client
* ethereum.privateKey

### Run Service using Dockers

```
./start.sh
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
