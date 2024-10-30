[![tests](https://github.com/iotubeproject/ioTube/actions/workflows/main.yml/badge.svg)](https://github.com/iotubeproject/ioTube/actions?query=workflow%3Atests)

<p align="center">
<img src="https://github.com/user-attachments/assets/9f249713-25d3-4de1-9458-69972b976699" width="720">
</p>

# ioTube - the decentralized bridge for Ethereum, Binance Smart Chain, Polygon(Matic), Solana and IoTeX
<a href="https://iotex.io/devdiscord" target="_blank">
  <img src="https://github.com/iotexproject/halogrants/blob/880eea4af074b082a75608c7376bd7a8eaa1ac21/img/btn-discord.svg" height="36px">
</a>

## Version history 

### V7.0 Payload support
TBD

### V6.0 Solana Bridge
Announcement: https://depinscan.io/news/2024-09-16/iotex-launches-solana-bridge

### V5.0 Web3 Support + Transactions + More Assets + Crosschain Tokens
Announcement: https://iotex.io/blog/cross-chain-polygon-web3/

### V4.0 Multi-Chain support
Launch date: Apr 13, 2021 ([Annoucement](https://medium.com/iotex/iotube-v4-cross-chain-bridge-for-iotex-ethereum-and-binance-smart-chain-9670c86723e2))

Because of Ethereum's high gas cost, many projects and users also adopted Binance Smart Chain (BSC) and Huobi Eco Chain (Heco). The need of supporting BSC and Heco is increasing. The demand of cross-chain support from Ethereum and BSC or other blockchains are increasing. We'd love to support them in ioTube.  

Matic support is launched on Jun 10, 2021 ([Annoucement](https://iotex.medium.com/iotube-v4-iotex-polygon-matic-cross-chain-token-swaps-are-live-bb2ae5bf41b4))

### V3.0 Rebuilt with largely fee reduction
Launch date: Feb 8, 2021 (<a href="https://community.iotex.io/t/iotube-v3-faster-cheaper-and-unified/2001">Announcement</a>)

Ethereum has been suffering high gas cost for a long period. The core-dev team rebuilt ioTube to largely reduced the gas cost on ETH side by introducing a relayer and putting signature offchain.

### V2.0 Multi Token Tube (IoTeX <-> Ethereum)
Launch date: Aug 31, 2020 (<a href="https://iotex.medium.com/iotube-cross-chain-bridge-to-connect-iotex-with-the-blockchain-universe-b0f5b08c1943">Announcement</a>)

V2 generalizes the ioTube V1 to support multiple assets on Ethereum & IoTeX blockchains.

### V1.0 IOTX tube 
Launch date: Apr 2019 (<a href="https://iotex.medium.com/everything-you-need-to-know-about-iotex-mainnet-alpha-b8d790e0bd55">Announcement</a>)

The first version of ioTube was built by core-dev in Apr. 2019 to facilitate the swap of IOTX-E (ERC20 version of IOTX on Ethereum) and IOTX mainnet token.


<h3>
      <a href="https://github.com/iotexproject/ioTube#deployement">Deployement</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube#usage">Usage</a>
      <span> | </span>
      <a href="https://github.com/iotexproject/ioTube/tree/master/docs">Documentation</a>
</h3>

Any bug report or feedback? Please submit an issue or discuss in https://discord.gg/jRqqSyGfUD.

## Submit a token to ioTube?

Please refer to token submission guide first: https://docs.iotube.org/introduction/token-submission

Feel free to reach out to https://github.com/guo for further discussion.

## Deployment
Different from traditional bridges, ioTube comes with two components:
- **a golang service** witnessing what has happened on both changes and relay the finalized information
- **a set of smart contracts** pre-deployed on both chains letting the legitimate witnesses relay information back and forth to facilitate cross-chain transferring of assets

### Deploy Contracts on IoTeX/Ethereum
* Deploy a TokenSafe `ts`
* Deploy a TokenList `tl1`, which stores tokens using `ts`
* Deploy a MinterPool `mp`
* Deploy a TokenList `tl2`, which stores tokens using `mp`
* Deploy a WitnessList `wl`
* Deploy a TokenCashier `tc` with `ts`, `tl1`, and `tl2`
* Deploy a TransferValidator `tv` with `ts`, `mp`, `tl1`, `tl2`, and `wl`
* Transfer ownership of `mp` to `tv`
* Transfer ownership of `ts` to `tv`
* Add witnesses to `wl`
* Add tokens to `tl1` or `tl2`

### Join as a Relayer

1. Edit `witness-service/relayer-config-iotex.yaml` and `witness-service/relayer-config-ethereum.yaml` to fill in the following fields:
* privateKey
* clientURL
* validatorContractAddress

2. start containers
```
cd witness-service
./start_relayer.sh
```

### Join as a Witness

1. Edit `witness-service/witness-config-iotex.yaml` and `witness-service/witness-config-ethereum.yaml` to fill in the following fields:
* privateKey
* clientURL
* validatorContractAddress
* cashierContractAddress

2. start containers
```
cd witness-service
./start_witness.sh
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
2. from XRC20 to ERC20: `2000 IOTX` (to cover the high ETH gas fee for witnesses)


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
- To transfer token from Ethereum to IoTeX: ~43,952 gas to set allowance and ~151,122 to lock it;
- To transfer token back from IoTeX to ETH: ~422,525 gas to unlock the token.

## Tube of IoTeX <-> Ethereum

Contracts on IoTeX 

- Wrapped IOTX: io15qr5fzpxsnp7garl4m7k355rafzqn8grrm0grz
- Token Safe: io1cj3f498390srqv765nnvaxuk0rpxyadzpfjz75
- Minter Pool: io1g7va274ltufv5nh4xawfmt0clel6tfz58p7n5r
- Standard Token List: io1t89whrwyfr0supctsqcx9n7ex5dd8yusfqhyfz 
- Proxy Token List: io1dn8nqk3pmmll990xz6a94fpradtrljxmmx5p8j 
- Witness List: io1hezp6d7y3246c5gklnnkh0z95qfld4zdsphhsw
- Token Cashier: io1gsr52ahqzklaf7flqar8r0269f2utkw9349qg8
- Transfer Validator: io10xr64as4krm5nufd5l2ddc43al6tl0smumkg7y <del>io1dwaxh2ml4fd2wg8cg35vhfsgdcyzrczffp3vus</del>

Contacts on Ethereum

- WETH: 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
- Standard Token List: 0x7c0bef36e1b1cbeb1f1a5541300786a7b608aede
- Proxy Token List: 0x73ffdfc98983ad59fb441fc5fe855c1589e35b3e
- Witness List: 0x8598dF1Ec0ac7dfBa802f4bDD93A6B93bd0AD83f 
- Token Safe: 0xc2e0f31d739cb3153ba5760a203b3bd7c27f0d7a 
- Minter Pool: 0x964f4f19bc823e72cc1f806021937cfc06f63b45 
- Token Cashier: 0xa0fd7430852361931b23a31f84374ba3314e1682
- Transfer Validator: 0xd8165188ccc135b3a3b2a5d2bc3af9d94753d955


## Tube of IoTeX <-> BSC (Binance Smart Chain) 

IoTeX Side
- Validator: io10xr64as4krm5nufd5l2ddc43al6tl0smumkg7y  (same for all)
- Standard Token list: io1h2d3r0d20t58sv6h707ppc959kvs8wjsurrtnk
- Proxy Token list: io17r9ehjstwj4gfqzwpm08fjnd606h04h2m6r92f
- Token Cashier: io1zjlng7je02kxyvjq4eavswp6uxvfvcnh2a0a3d

BSC Side

- Witness List Address:  0x8119411F5A78F73784A1B87dE43d452DA4A1EE3F
- Minter Pool Address:  0xf72CFb704d49aC7BB7FFa420AE5f084C671A29be
- Token Safe Address:  0xFBe9A4138AFDF1fA639a8c2818a0C4513fc4CE4B
- Mintable Token List Address:  0xa6ae9312D0AA3CC74d969Fcd4806d7729A321EE3
- Standard Token List Address:  0x0d793F4D4287265B9bdA86b7a4083193E8743b34
- Token Cashier Address:  0x797f1465796fd89ea7135e76dbc7cdb136bba1ca <del>0x082020Ae0B38fD1bef48895c6cFf4428e420F400</del>
- Transfer Validator Address:  0x116404F86e97846110EA08cd52fC2882d4AD3123

## Tube of IoTeX <-> Heco (Huobi Eco Chain)

Please note this is still being added. They are not in production use yet. Please contact us if you want to use them without interfaces.

IoTeX Side:
- Validator: io10xr64as4krm5nufd5l2ddc43al6tl0smumkg7y  (same for all)
- Standard Token list: io1kh0vgtxyamdkzvrlvga3r8l7r2plm7phd9ywv9
- Proxy Token list:  io18uqxuel6d93hluua4d9jxs8rjw3r2qe5g3adgk
- Token Cashier:  io1s6m6j3cdj0j7hlgupx0pww7wscvkepdjfwgkra

Heco Side:

- Witness List Address:  0x2f1a0BCa4005eBfD6A589850F436c8D8f9c2aEd2
- Minter Pool Address:  0xd2165D222B3dAF2528Fc1b1Aa2DB18B8821EE623
- Token Safe Address:  0x1E58cA53d90fe9B37F7f6AEB548b4BC7c6292C17
- Mintable Token List Address:  0x12af43ef94B05A0a3447A05eEE629C7D88A30a5f
- Standard Token List Address:  0xA239F03Cda98A7d2AaAA51e7bF408e5d73399e45
- Token Cashier Address:  0xC8DC8dCDFd94f9Cb953f379a7aD8Da5fdC303F3E
- Transfer Validator Address:  0xDe9395d2f4940aA501f9a27B98592589D14Bb0f7

## Tube of IoTeX <-> Polygon (formerly Matic) 

Launched on Jun 10, 2021.

We started adding 0x address to this doc because of IoTeX start supporting 0x address and web3 from IoTeX V1.2 (Babel API). Some IoTeX address are same as other tubes and we include 0x addresses here.

IoTeX side:

- Validator: io10xr64as4krm5nufd5l2ddc43al6tl0smumkg7y 
  - (0x7987aaf615b0f749f12da7d4d6e2b1eff4bfbe1b)
- Standard Token list: io197rk3nff9622pkncvuvhfwyms73esdtwph4rlq
  - (0x2F8768cD292E94A0Da78671974B89B87a398356E)
- Proxy Token list: io16at6mlcwcsrqutz2zhuhwam87h988r9fcdauk8
  - (0xD757adFF0eC4060e2c4A15f9777767f5Ca738Ca9)
- Token Cashier:   
  - io12s9f9hv4zsr7umy5hxt6g0k0xr4x6pxdp5w998
  - (0x540a92dd951407ee6c94b997a43ecf30ea6d04cd)
- Token safe: from old. 
  - (0xc4A29a94f12be03033daa4e6Ce9b9678c26275a2)
- Minter pool (old); io1g7va274ltufv5nh4xawfmt0clel6tfz58p7n5r
  - (0x4799d57abf5f12ca4ef5375c9dadf8fe7fa5a454)
  
Matic side:

- Witness List Address:  0x1E58cA53d90fe9B37F7f6AEB548b4BC7c6292C17
- Minter Pool Address:  0x12af43ef94B05A0a3447A05eEE629C7D88A30a5f
- Token Safe Address:  0xA239F03Cda98A7d2AaAA51e7bF408e5d73399e45
- Mintable Token List Address:  0xC8DC8dCDFd94f9Cb953f379a7aD8Da5fdC303F3E
- Standard Token List Address:  0xDe9395d2f4940aA501f9a27B98592589D14Bb0f7
- Token Cashier Address:  0xf72CFb704d49aC7BB7FFa420AE5f084C671A29be
- Transfer Validator Address: 0xFBe9A4138AFDF1fA639a8c2818a0C4513fc4CE4B


## Crosschain IOTX (WIP)

- CIOTX on IoTeX: `0x99B2B0eFb56E62E36960c20cD5ca8eC6ABD5557A`
- CIOTX on Polygon: `0x300211Def2a644b036A9bdd3e58159bb2074d388`
- CIOTX on BSC: `0x2aaF50869739e317AB80A57Bf87cAA35F5b60598`
