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

Follow instructions [here](https://github.com/iotexproject/ioTube/blob/master/witness-service/README.md)

## Usage

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

