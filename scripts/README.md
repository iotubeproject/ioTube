# Scripts

# Add Crosschain Token

## Compile

```
yarn
yarn hardhat compile
```

## deploy Crosschain token on source chain

```
export C_TOKEN=0x...
export CO_TOKEN=0x...
export TOKEN_NAME=Crosschain ABC
export TOKEN_SYMBOL=CABC
export TOKEN_DECIMALS=18
export TOKEN_MAX=0.5
export TOKEN_MIN=0.06
export REF_CHAIN=bsc
yarn hardhat run scripts/add-crosschain-token-source.js --network mainnet
```

## deploy Crosschain token on destination chain

```
export C_TOKEN=0x...
export TOKEN_NAME=Ethereum ABC
export TOKEN_SYMBOL=ABC-E
export TOKEN_DECIMALS=18
export TOKEN_MAX=0.5
export TOKEN_MIN=0.06
yarn hardhat run scripts/add-crosschain-token-dest.js --network iotex
```
