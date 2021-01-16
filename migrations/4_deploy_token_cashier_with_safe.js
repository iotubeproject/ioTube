const TokenCashierWithSafeV2 = artifacts.require('TokenCashierWithSafeV2')
const TokenList = artifacts.require('TokenList')
const TokenSafe = artifacts.require('TokenSafe')

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const tokenSafe = await TokenSafe.deployed();
        const tokenList = await TokenList.deployed();
        const tokenCashier = await deployer.deploy(TokenCashierWithSafeV2, tokenList.address, tokenSafe.address);
        console.log('Token Cashier with Safe V2 Address: ', tokenCashier.address);
    });
}
