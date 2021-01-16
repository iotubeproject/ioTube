const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const TransferValidatorWithTokenSafe = artifacts.require('TransferValidatorWithTokenSafe');
const WitnessList = artifacts.require('WitnessList');

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const tokenList = await TokenList.deployed();
        const tokenSafe = await TokenSafe.deployed();
        const witnessList = await deployer.deploy(WitnessList);
        console.log('Witness List Address: ', witnessList.address);
        const validator = await deployer.deploy(TransferValidatorWithTokenSafe, tokenSafe.address, tokenList.address, witnessList.address);
        console.log('Transfer Validator Address: ', validator.address);
    });
}
