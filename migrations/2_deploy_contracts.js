const MinterPool = artifacts.require('MinterPool');
const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const TokenCashier = artifacts.require('TokenCashier');
const TransferValidator = artifacts.require('TransferValidator');
const WitnessList = artifacts.require('WitnessList');

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const witnessList = await deployer.deploy(WitnessList);
        console.log('Witness List Address: ', witnessList.address);
        const minterPool = await deployer.deploy(MinterPool);
        console.log('Minter Pool Address: ', minterPool.address);
        const tokenSafe = await deployer.deploy(TokenSafe);
        console.log('Token Safe Address: ', tokenSafe.address);
        const mintableTokenList = await deployer.deploy(TokenList);
        console.log('Mintable Token List Address: ', mintableTokenList.address);
        const standardTokenList = await deployer.deploy(TokenList);
        console.log('Standard Token List Address: ', standardTokenList.address);
        const tokenCashier = await deployer.deploy(
            TokenCashier,
            "0x0000000000000000000000000000000000000000",
            [standardTokenList.address, mintableTokenList.address],
            [tokenSafe.address, minterPool.address],
        );
        console.log('Token Cashier Address: ', tokenCashier.address);
        const validator = await deployer.deploy(
            TransferValidator,
            witnessList.address,
        );
        console.log("Add new pairs")
        await validator.addPair(standardTokenList.address, tokenSafe.address);
        await validator.addPair(mintableTokenList.address, minterPool.address),
        console.log('Transfer Validator Address: ', validator.address);
        await minterPool.transferOwnership(validator.address);
        console.log('Ownership of minter pool is transferred to validator');
        await tokenSafe.transferOwnership(validator.address);
        console.log('Ownership of token safe is transferred to validator');
    });
}
