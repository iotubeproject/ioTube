<p align="center">
  <img src="https://github.com/iotexproject/ioTube/blob/master/ioTube.png" width="480px">
</p>

&nbsp;

# ioTube
A decentralized bridge between Ethereum and IoTeX for exchanging native coins and ERC20/XRC20 tokens

# Deploy on IoTeX
## Deploy a MinterPool `mp`
## Deploy a TokenList `tl`
## Deploy a VoterList `vl`
## Deploy a BurnableTokenCashier with `tl`
## Deploy a TokenValidatorWithMinterPool with `mp`, `tl`, `vl`
## Transfer ownership of `mp` to the TokenValidatorWithMinterPool
## Add voters to `vl`
## Deploy a ShadowToken and add it to `tl`

# Deploy on Ethereum
## Deploy a TokenSafe `ts`
## Deploy a TokenList `tl`
## Deploy a VoterList `vl`
## Deploy a TokenCashierWithSafe with `tl` and `ts`
## Deploy a TokenValidatorWithTokenSafe with `ts`, `tl`, and `vl`
## Transfer ownership of `ts` to the TokenValidatorWithTokenSafe
## Add voters to `vl`
## Add existing ERC20 token to `tl`
