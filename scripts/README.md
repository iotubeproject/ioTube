# How to use?

Please download the scripts to your local path.  Please make sure ioctl has been installed with latest version.

```
cd path-to-scripts
chmod +x *.sh
```

To wrap IOTX to WIOTX,
```
./iotx2wiotx.sh [amount to wrap]
```

To unwrap WIOTX to IOTX,
```
./wiotx2iotx.sh [amount to unwrap]
```

# Add Crosschain Token

## deploy Crosschain token on source chain

```
export CO_TOKEN=0x...
export TOKEN_NAME=Crosschain ABC
export TOKEN_SYMBOL=CABC
export TOKEN_DECIMALS=18
export TOKEN_MAX=10000000000000000000
export TOKEN_MIN=1000000000000000000
yarn hardhat run scripts/add-crosschain-token-source.js --network mainnet
```

## deploy Crosschain token on destination chain

```
export TOKEN_NAME=Ethereum ABC
export TOKEN_SYMBOL=ABC-E
export TOKEN_DECIMALS=18
export TOKEN_MAX=10000000000000000000
export TOKEN_MIN=1000000000000000000
yarn hardhat run scripts/add-crosschain-token-dest.js --network iotex
```
