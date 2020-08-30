<p align="center">
  <img src="https://github.com/iotexproject/ioTube/blob/master/ioTube.png" width="480px">
</p>

<p>
  <strong>A multi-assets, fully decentralized and bidirectional bridge for exchanging ERC20/XRC20 tokens between Ethereum and IoTeX.</strong>
  
ioTube is a decentralized cross-chain bridge that enables the bidirectional exchange of digital assets (e.g., tokens, stable coins) between IoTeX and other blockchain networks. The first version of ioTube was built by core-dev around 2019 June to facilitate the swap of IOTX-ERC20 and IOTX until now, and this released version generalizes the original ioTube by making it more decentralized and supports multiple assets on Ethereum & IoTeX blockchains. In the future, we plan to add support for more blockchain networks to increase the reach and impact of ioTube.
</p>

<h3>
      <a href="https://github.com/iotexproject/ioTube#deployement">Deployement</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube#usage">Usage</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube/tree/master/docs">Documentation</a>
</h3>

&nbsp;

## Deployment
Different from traditional bridges, ioTube comes with two components:
- **a golang service** witnessing what has happened on both changes and relay the finalized information
- **a set of smart contracts** pre-deployed on both chains letting the legitimate witnesses relay information back and forth to facilitate cross-chain transferring of assets

### Deploy Contracts on IoTeX
* Deploy a MinterPool `mp`
* Deploy a TokenList `tl`
* Deploy a WitnessList `wl`
* Deploy a BurnableTokenCashier with `tl`
* Deploy a TransferValidatorWithMinterPool with `mp`, `tl`, `wl`
* Transfer ownership of `mp` to the TransferValidatorWithMinterPool
* Add witnesses to `wl`

### Deploy Contracts on Ethereum
* Deploy a TokenSafe `ts`
* Deploy a TokenList `tl`
* Deploy a WitnessList `wl`
* Deploy a TokenCashierWithSafe with `tl` and `ts`
* Deploy a TransferValidatorWithTokenSafe with `ts`, `tl`, and `wl`
* Transfer ownership of `ts` to the TransferValidatorWithTokenSafe
* Add witnesses to `wl`

### Join as a Witness

#### Option 1: Run Witness Service on Dockers

1. Prepare config by copying `service.yaml` to current directory, and fill in the following fields:
* iotex.privateKey
* ethereum.privateKey
* ethereum.client

2. start containers
```
./start.sh
```

#### Option 2: Run Witness Service on Heroku (not ready yet)
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


### Transfer assets between IoTeX and Ethereum
TBD

### Add an ERC20 token to ioTube
* Add this token to `tl` on Ethereum
* Deploy a ShadowToken and add it to `tl` on IoTeX


## Security
- Each witness needs to be registered to `WitnessList` contract on both chains based on PoA (Proof-of-Authority).
- Each transferring of assets from one chain to another needs the endorsement from more than 2/3 of the registered witnesses; otherwise, it will not be successful.
- IoTeX has instant finality, meaning, to a witness, one confirmed block indicates finalization of a `burn` event, while on the Ethereum side, a witness needs to wai 12 blocks before making any endorsement.
- To launch ioTube reliably, we have limited the min/max of the asset that can be moved around. These limits can be lifted once ioTube gets more stress validated.

## Gas Costs
Gas fees on IoTeX are negligible, both for bridge maintenance and for asset transfer. The estimated gas fees on Ethereum side are:
- To transfer token from Ethereum to IoTeX: ~XXX gas to set allowance and ~XXX to lock it;
- To transfer token back from IoTeX to ETH: ~XXX gas to unlock the token.

