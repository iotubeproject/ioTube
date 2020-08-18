<p align="center">
  <img src="https://github.com/iotexproject/ioTube/blob/master/ioTube.png" width="480px">
</p>

<p>
  <strong>A fully decentralized and bidirectional bridge for exchanging ERC20/XRC20 tokens between Ethereum and IoTeX</strong>
</p>

<h3>
      <a href="https://github.com/iotexproject/ioTube#deployement">Deployement</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube#usage">Usage</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube/tree/master/docs">Documentation</a>
</h3>

&nbsp;

## Deployement
### Deploy Contracts on IoTeX
* Deploy a MinterPool `mp`
* Deploy a TokenList `tl`
* Deploy a VoterList `vl`
* Deploy a BurnableTokenCashier with `tl`
* Deploy a TransferValidatorWithMinterPool with `mp`, `tl`, `vl`
* Transfer ownership of `mp` to the TransferValidatorWithMinterPool
* Add voters to `vl`

### Deploy Contracts on Ethereum
* Deploy a TokenSafe `ts`
* Deploy a TokenList `tl`
* Deploy a VoterList `vl`
* Deploy a TokenCashierWithSafe with `tl` and `ts`
* Deploy a TransferValidatorWithTokenSafe with `ts`, `tl`, and `vl`
* Transfer ownership of `ts` to the TransferValidatorWithTokenSafe
* Add voters to `vl`

### Join as a Witness (Voter)

TBD

## Usage

### Transfer assets between IoTeX and Ethereum
TBD

### Add a ERC20 token to ioTube
* Add this token to `tl` on Ethereum
* Deploy a ShadowToken and add it to `tl` on IoTeX


## Security
- Each witnesses (voters) needs to be registered to `VoterList` contract on both chains based on PoA (Proof-of-Authority). Once enough number of witnesses (voters) joined, ioTube will be switched to use PoS (Proof-of-Stake) with slash policies. 
- Each transferring of assets from one chain to another needs the endorsement from more than 2/3 of the registered witnesses (voters); otherwise it will not be successful.
- IoTeX has instant finality, meaning, to a witness, one confirmed block indicates finalization of a `burn` event, while on the Ethereum side, a witness needs to wai 12 blocks before making any endorsement.
- To launch ioTube in a safe way, we have limit the min/max of the asset that can be move around. These limits can be lifted up once ioTube gets more stress validated.

## Gas Costs
Gas fees on IoTeX are negligible, both for bridge maintenance and for asset transfer. The estimated gas fees on Ethereum side are:
- To transfer token from Ethereum to IoTeX: ~XXX gas to set allowance and ~XXX to lock it;
- To transfer token back from IoTeX to ETH: ~XXX gas to unlock the token.

