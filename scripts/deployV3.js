const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  const WitnessList = await ethers.getContractFactory("WitnessListV3");
  const witnessList = await WitnessList.deploy();
  await witnessList.waitForDeployment();
  console.log("Witness List Address:", witnessList.target);

  const WitnessManager = await ethers.getContractFactory("WitnessManager");
  const witnessManager = await WitnessManager.deploy(deployer.address, witnessList.target);
  await witnessManager.waitForDeployment();
  console.log("WitnessManager Address:", witnessManager.target);

  console.log("Transferring ownership of WitnessListV3 to WitnessManager...");
  let tx = await witnessList.transferOwnership(witnessManager.target);
  await tx.wait();
  console.log("Ownership of WitnessListV3 transferred.");

  const MinterPool = await ethers.getContractFactory("MinterPool");
  const minterPool = await MinterPool.deploy();
  await minterPool.waitForDeployment();
  console.log("Minter Pool Address:", minterPool.target);

  const TokenSafe = await ethers.getContractFactory("TokenSafe");
  const tokenSafe = await TokenSafe.deploy();
  await tokenSafe.waitForDeployment();
  console.log("Token Safe Address:", tokenSafe.target);

  const TokenList = await ethers.getContractFactory("TokenList");

  const mintableTokenList = await TokenList.deploy();
  await mintableTokenList.waitForDeployment();
  console.log("Mintable Token List Address:", mintableTokenList.target);

  const standardTokenList = await TokenList.deploy();
  await standardTokenList.waitForDeployment();
  console.log("Standard Token List Address:", standardTokenList.target);

  const WitnessTokenList = await ethers.getContractFactory("WitnessTokenList");
  const witnessTokenList = await WitnessTokenList.deploy();
  await witnessTokenList.waitForDeployment();
  console.log("WitnessTokenList Address:", witnessTokenList.target);

  const wrappedCoin = process.env.WRAPPED_COIN_ADDRESS;
  if (!wrappedCoin) {
        console.log("WRAPPED_COIN_ADDRESS is not specified in .env file, early quit.");
    return;
  }
  console.log("Wrapped Coin Address:", wrappedCoin);

  const TokenCashierWithPayload = await ethers.getContractFactory("TokenCashierWithPayload");
  const tokenCashier = await TokenCashierWithPayload.deploy(wrappedCoin);
  await tokenCashier.waitForDeployment();
  console.log("TokenCashierWithPayload Address:", tokenCashier.target);

  console.log("Adding token lists to TokenCashierWithPayload...");
  tx = await tokenCashier.addTokenList(standardTokenList.target, tokenSafe.target);
  await tx.wait();
  tx = await tokenCashier.addTokenList(mintableTokenList.target, "0x0000000000000000000000000000000000000000");
  await tx.wait();
  console.log("Token lists added.");

  const TransferValidatorV3 = await ethers.getContractFactory("TransferValidatorV3");
  const validator = await TransferValidatorV3.deploy();
  await validator.waitForDeployment();
  console.log("TransferValidatorV3 Address:", validator.target);

  console.log("Adding new pairs to TransferValidatorV3...");
  tx = await validator.addWitnessPair(witnessTokenList.target, witnessList.target);
  await tx.wait();
  tx = await validator.addMinterPair(standardTokenList.target, tokenSafe.target);
  await tx.wait();
  tx = await validator.addMinterPair(mintableTokenList.target, minterPool.target);
  await tx.wait();
  console.log("Pairs added.");

  console.log("Transferring ownership of MinterPool to TransferValidator...");
  tx = await minterPool.transferOwnership(validator.target);
  await tx.wait();
  console.log("Ownership of MinterPool transferred.");

  console.log("Transferring ownership of TokenSafe to TransferValidator...");
  tx = await tokenSafe.transferOwnership(validator.target);
  await tx.wait();
  console.log("Ownership of TokenSafe transferred.");

  // save addresses to file
  const addresses = {
    witnessList: witnessList.target,
    witnessManager: witnessManager.target,
    minterPool: minterPool.target,
    tokenSafe: tokenSafe.target,
    mintableTokenList: mintableTokenList.target,
    standardTokenList: standardTokenList.target,
    witnessTokenList: witnessTokenList.target,
    tokenCashier: tokenCashier.target,
    transferValidator: validator.target,
  };

  const addressesPath = path.join(__dirname, "deployment-addressesV3.json");
  fs.writeFileSync(addressesPath, JSON.stringify(addresses, null, 2));

  console.log("Deployment addresses saved to:", addressesPath);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
