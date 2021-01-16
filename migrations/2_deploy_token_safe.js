const TokenSafe = artifacts.require('TokenSafe')

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const tokenSafe = await deployer.deploy(TokenSafe);
        console.log('Token Safe Address: ', tokenSafe.address);
    });
}
