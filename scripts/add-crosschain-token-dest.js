const hre = require("hardhat");
const addresses = require("./addresses");

async function main() {
  const tubeAddress = addresses[hre.network.name];
  const co_token = process.env.CO_TOKEN;
  if (co_token === undefined || co_token === "") {
    console.log("Must use env variable to provide co-token address: export CO_TOKEN=0x...");
    return;
  }
  const name = process.env.TOKEN_NAME;
  if (name === undefined || name === "") {
    console.log("Must use env variable to provide token name: export TOKEN_NAME=Crosschain ABC");
    return;
  }
  const symbol = process.env.TOKEN_SYMBOL;
  if (symbol === undefined || symbol === "") {
    console.log("Must use env variable to provide token symbol: export TOKEN_SYMBOL=CABC");
    return;
  }
  const decimals = process.env.TOKEN_DECIMALS;
  if (decimals === undefined || decimals === "") {
    console.log("Must use env variable to provide token decimals: export TOKEN_DECIMALS=18");
    return;
  }
  const max = process.env.TOKEN_MAX;
  if (max === undefined || max === "") {
    console.log("Must use env variable to provide token max: export TOKEN_MAX=1.5");
    return;
  }
  const min = process.env.TOKEN_MIN;
  if (min === undefined || min === "") {
    console.log("Must use env variable to provide token min: export TOKEN_MIN=1");
    return;
  }

  console.log(`Deploy cToken[${name}, ${symbol}, ${decimals}]...`);
  const cToken = await hre.ethers.deployContract("CrosschainERC20", [
    "0x0000000000000000000000000000000000000000",
    tubeAddress.minter_pool,
    name,
    symbol,
    decimals
  ]);
  await cToken.waitForDeployment();
  console.log(
    `cToken[${name}, ${symbol}, ${decimals}] deployed to ${cToken.target}`
  );
  
  if (tubeAddress.proxy_token_list !== "") {
    console.log(`Add cToken ${cToken.target} to proxy token list...`);
    const proxyTokenList = await hre.ethers.getContractAt("TokenList", tubeAddress.proxy_token_list);
    let tx = await proxyTokenList.addToken(
      cToken.target,
      hre.ethers.FixedNumber.fromString(min, {decimals: Number(decimals)}).value,
      hre.ethers.FixedNumber.fromString(max, {decimals: Number(decimals)}).value
    );
    let receipt = await tx.wait();
    if (receipt.status !== 0) {
      console.log(`Add cToken to proxy token tx fail, txHash: ${tx.hash}`);
      return;
    }
  } else {
    console.log(`Skip add proxy token list`);
  }

  console.log(`Add cToken ${cToken.target} successful`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
