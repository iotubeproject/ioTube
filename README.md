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
Launch date: Nov 6, 2024

This new version enables crosschain calls by adding a new field 'payload'. The payload of a crosschain transfer will be decoded by the whitelisted recipient contracts. For example, an unwrapper for the crosschain tokens is whitelisted such that users don't need to unwrap the crosschain tokens by themselves.

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

### Join as a Witness

1. Prepare the working directory $IOTEX_WITNESS, default ~/iotex-witness. Create file $IOTEX_WINTESS/etc/.env and add:
* WITNESS_PRIVATE_KEY
* RELAYER_URL

2. start containers
```
cd witness-service
./start_witness.sh
```

3. modify some configs in $IOTEX_WITNESS/etc/*.secret.yaml:
* clientURL

4. Clean up everything by running
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

1. from the other chains to IoTeX: `0`
2. from IoTeX to Ethereum: `500 IOTX`
3. from IoTeX to BSC: `20 IOTX`
4. from IoTeX to Polygon: `2 IOTX`
all the fees are used to cover the gas fee


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
- Token Cashier: 0x4Aad62e13cBd5cD6F5BcEc649cc24f6185011fe1
- Transfer Validator: 0x915241176a644A29F46F5F4F8Bd49309e461a9bD

Contacts on Ethereum

- WETH: 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
- Standard Token List: 0x7c0bef36e1b1cbeb1f1a5541300786a7b608aede
- Proxy Token List: 0x73ffdfc98983ad59fb441fc5fe855c1589e35b3e
- Witness List: 0x8598dF1Ec0ac7dfBa802f4bDD93A6B93bd0AD83f 
- Token Safe: 0xc2e0f31d739cb3153ba5760a203b3bd7c27f0d7a 
- Minter Pool: 0x964f4f19bc823e72cc1f806021937cfc06f63b45 
- Token Cashier: 0x1B9AA865d74b2B77fFdbCF507B56a7b3AB43Bac4
- Transfer Validator: 0xE7eBA1CEA51EC9B3AcCC16728e3B8786560c59d5

## Tube of IoTeX <-> BSC (Binance Smart Chain) 

IoTeX Side
- Validator: 0x915241176a644A29F46F5F4F8Bd49309e461a9bD  (same for all)
- Standard Token list: io1h2d3r0d20t58sv6h707ppc959kvs8wjsurrtnk
- Proxy Token list: io17r9ehjstwj4gfqzwpm08fjnd606h04h2m6r92f
- Token Cashier: 0xa016f866a606221E9C9E6ab3a942Bfc81F6074f4

BSC Side

- Witness List Address:  0x8119411F5A78F73784A1B87dE43d452DA4A1EE3F
- Minter Pool Address:  0xf72CFb704d49aC7BB7FFa420AE5f084C671A29be
- Token Safe Address:  0xFBe9A4138AFDF1fA639a8c2818a0C4513fc4CE4B
- Mintable Token List Address:  0xa6ae9312D0AA3CC74d969Fcd4806d7729A321EE3
- Standard Token List Address:  0x0d793F4D4287265B9bdA86b7a4083193E8743b34
- Token Cashier Address:  0x78de1E0b76523Ac6E190F89FFC46571346940204
- Transfer Validator Address:  0x95C6F6Af2c0Fa069768203FDa963d7626efC794a

## Tube of IoTeX <-> Polygon (formerly Matic) 

Launched on Jun 10, 2021.

We started adding 0x address to this doc because of IoTeX start supporting 0x address and web3 from IoTeX V1.2 (Babel API). Some IoTeX address are same as other tubes and we include 0x addresses here.

IoTeX side:

- Validator: 0x915241176a644A29F46F5F4F8Bd49309e461a9bD
- Standard Token list: io197rk3nff9622pkncvuvhfwyms73esdtwph4rlq
  - (0x2F8768cD292E94A0Da78671974B89B87a398356E)
- Proxy Token list: io16at6mlcwcsrqutz2zhuhwam87h988r9fcdauk8
  - (0xD757adFF0eC4060e2c4A15f9777767f5Ca738Ca9)
- Token Cashier: 0x8114746E4308a4d3Ff2a74B66414fF35657Fa0E2
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
- Token Cashier Address:  0x990B503f8C7353f1caB6f9D5bbF8f0Be2718D731
- Transfer Validator Address: 0x87E2D48De6CC2029fFc1a915462e4Aa597890cd6


## Crosschain IOTX (WIP)

- CIOTX on IoTeX: `0x99B2B0eFb56E62E36960c20cD5ca8eC6ABD5557A`
- CIOTX on Polygon: `0x300211Def2a644b036A9bdd3e58159bb2074d388`
- CIOTX on BSC: `0x2aaF50869739e317AB80A57Bf87cAA35F5b60598`
