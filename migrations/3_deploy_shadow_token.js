const ShadowToken = artifacts.require('ShadowToken');

module.exports = function(deployer, network, accounts) {
    return deployer.then(async () => {
        const shadowToken = await deployer.deploy(ShadowToken, 
            "0xDdB5a53737fF9bDe253857f4b165EC12337366dd",
            "0xd0a1e359811322d97991e03f863a0c30c2cf029c",
            "Tube ETH",
            "TETH",
            18
        );
        console.log('Shadow Token Address: ', shadowToken.address);
    });
}
