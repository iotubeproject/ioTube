const hre = require("hardhat");
const addresses = require("./addresses");

async function main() {
  let tubeAddress = addresses[hre.network.name];
  if (hre.network.name == "iotex") {
    const ref = process.env.REF_CHAIN;
    if (ref === undefined || ref === "") {
      console.log("Must use env variable to provide ref chain: export REF_CHAIN=bsc");
      return;
    }
    tubeAddress = tubeAddress[ref.toLowerCase()];
  }
  const owner = new hre.ethers.Wallet(
    process.env[`PRIVATE_KEY_${hre.network.name.toUpperCase()}`],
    hre.ethers.provider
  )

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
  let c_token = process.env.C_TOKEN;
  if (c_token === undefined || c_token === "") {
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

    console.log(`Deploy cToken[${name}, ${symbol}, ${decimals}]...`);
    const cToken = await hre.ethers.deployContract("CrosschainERC20", [
      "0x0000000000000000000000000000000000000000",
      tubeAddress.minter_pool,
      name,
      symbol,
      decimals
    ], owner);
    await cToken.waitForDeployment();
    console.log(
      `cToken[${name}, ${symbol}, ${decimals}] deployed to ${cToken.target}`
    );
    c_token = cToken.target;
  } else {
    console.log("cToken:", c_token);
  }


  if (tubeAddress.proxy_token_list !== "") {
    console.log(`Add cToken ${c_token} to proxy token list...`);
    const proxyTokenList = await hre.ethers.getContractAt("TokenList", tubeAddress.proxy_token_list);
    const decimals = process.env.TOKEN_DECIMALS;
    if (decimals === undefined || decimals === "") {
      console.log("Must use env variable to provide token decimals: export TOKEN_DECIMALS=18");
      return;
    }
    let tx = await proxyTokenList.connect(owner).addToken(
      c_token,
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

  console.log(`Add cToken ${c_token} successful`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
