const MinterPool = artifacts.require('MinterPool');
const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const TokenCashier = artifacts.require('TokenCashier');
const TransferValidator = artifacts.require('TransferValidator');
const WitnessList = artifacts.require('WitnessList');

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        // deploy a TokenSafe `ts`
        const tokenSafe = await deployer.deploy(TokenSafe);
        console.log('Token Safe Address: ', tokenSafe.address);

        // Deploy a TokenList `tl1`, which stores tokens using `ts`
        const standardTokenList = await deployer.deploy(TokenList);
        console.log('Standard Token List Address: ', standardTokenList.address);

        // Deploy a MinterPool `mp`
        const minterPool = await deployer.deploy(MinterPool);
        console.log('Minter Pool Address: ', minterPool.address);

        // Deploy a TokenList `tl2`, which stores tokens using `mp`
        const mintableTokenList = await deployer.deploy(TokenList);
        console.log('Mintable Token List Address: ', mintableTokenList.address);

        // Deploy a WitnessList `wl`
        const witnessList = await deployer.deploy(WitnessList);
        console.log('Witness List Address: ', witnessList.address);

        const wrappedCoin = process.env.WRAPPED_COIN;
        if (wrappedCoin == undefined) {
            console.log("Wrapped Coin is not specified, early quit.");
            return;
        }

        // Deploy a TokenCashier `tc` with `ts`, `tl1`, and `tl2`
        const tokenCashier = await deployer.deploy(
            TokenCashier,
            wrappedCoin,
            [standardTokenList.address, mintableTokenList.address],
            [tokenSafe.address, minterPool.address],
        );
        console.log('Token Cashier Address: ', tokenCashier.address);

        // Deploy a TransferValidator `tv` with `ts`, `mp`, `tl1`, `tl2`, and `wl`
        const validator = await deployer.deploy(
            TransferValidator,
            witnessList.address,
        );
        console.log("Add new pairs")
        await validator.addPair(standardTokenList.address, tokenSafe.address);
        await validator.addPair(mintableTokenList.address, minterPool.address),
    
        console.log('Transfer Validator Address: ', validator.address);
        // Transfer ownership of `mp` to `tv`
        await minterPool.transferOwnership(validator.address);
        console.log('Ownership of minter pool is transferred to validator');
    
        // Transfer ownership of `ts` to `tv`
        await tokenSafe.transferOwnership(validator.address);
        console.log('Ownership of token safe is transferred to validator');
    });
}
