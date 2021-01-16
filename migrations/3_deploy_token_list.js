const TokenList = artifacts.require('TokenList')

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const tokenList = await deployer.deploy(TokenList);
        console.log('Token List Address: ', tokenList.address);
    });
}
