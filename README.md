<p align="center">
  <img src="https://github.com/iotexproject/ioTube/blob/master/ioTube.png" width="480px">
</p>

<p>
  <strong>A multi-assets, fully decentralized and bidirectional bridge for exchanging ERC20/XRC20 tokens between Ethereum and IoTeX.</strong>
  
ioTube is a decentralized cross-chain bridge that enables the bidirectional exchange of digital assets (e.g., tokens, stable coins) between IoTeX and other blockchain networks. The <a href="https://member.iotex.io/tools/iotex">first version</a> of ioTube was built by core-dev around 2019 June to facilitate the swap of IOTX-ERC20 and IOTX until now, and this released version generalizes the original ioTube to support multiple assets on Ethereum & IoTeX blockchains. In the future, we plan to add support more blockchain networks to increase the reach and impact of ioTube.
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

1. Edit `witness-service/service.yaml` to fill in the following fields:
* iotex.privateKey
* ethereum.privateKey
* ethereum.client

2. start containers
```
cd witness-service
./start.sh
```

3. Clean up everything by running
```
./clean-all.sh
```

### Transfer assets between IoTeX and Ethereum

Please use dApp ioTube https://tube.iotex.io. Please note that the service is still in beta mode.  

#### From ERC20 to XRC20

1. open https://tube.iotex.io/eth in a metamask installed browser. (eg. Firefox/Chrome/Brave + metamask) 
2. Choose supported ERC20 token from the list.
3. Enter the amount. 
4. Click `Approve` button to approve ERC20 token transfer and sign on metamask.
5. Click `Convert` button and sign on metamask.
6. After `12` confirmations of Ethereum network and `2/3 + 1` confirmations from witnesses, the XRC20 tokens will be minted and sent to your IoTeX address.
7. You can add the token to your <a href="http://iopay.iotex.io/">ioPay</a> to see and use them.

#### From XRC20 to ERC20

1. open https://tube.iotex.io/iotx in ioPay desktop supported broswers (eg. Chrome/Firefox/Brave with ioPay desktop installed) or ioPay Android/iOS.
2. Choose supported XRC20 token from the list.
3. Enter the amount.
4. Click `Approve` button to approve XRC20 token transfer and sign on ioPay.
5. Click `Convert` button and sign on ioPay.
6. After `1` confirmation of IoTeX network and `2/3 + 1` confirmations from witnesses, the XRC20 token will be burnt and ERC20 token will be sent to your ETH wallet.

#### Fees
Tube fees: `0`

Network fees: 

1. from ERC20 to XRC20: `0`
2. from XRC20 to ERC20: `4000 IOTX` (to cover the high ETH gas fee for witnesses)


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


## Current Supported Tokens
1. WETH (0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2) - min=0.1 max=1000 

Contract on ioETH (io1qfvgvmk6lpxkpqwlzanqx4atyzs86ryqjnfuad)

2. PAXG (0x45804880de22913dafe09f4980848ece6ecbaf78) - min=0.01 max=200

Contract on ioPAXG (TBA)

