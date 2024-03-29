const TransferValidator = artifacts.require('TransferValidatorV3');
const Unwrapper = artifacts.require('Unwrapper');
const MinterPool = artifacts.require('MinterPool');
const ShadowToken = artifacts.require('ShadowToken');
const TokenList = artifacts.require('TokenList');
const WitnessList = artifacts.require('WitnessList');
const CrosschainERC20 = artifacts.require('CrosschainERC20');
const StandardToken = artifacts.require('StandardToken');
const Account = require('eth-lib/lib/account');
const {AbiCoder} = require('ethers');

const witnessPrivateKeys = [
    '0x388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418',
    '0x659cbb0e2411a44db63778987b1e22153c086a95eb6b18bdf89de078917abc63',
    '0x82d052c865f5763aad42add438569276c00d3d88a2d062d36b2bae914d58b8c8',
    '0xaa3680d5d48a8283413f7a108367c7299ca73f553735860a87b08f39395618b7',
    '0x0f62d96d6675f32685bbdb8ac13cda7c23436f63efbb9d07700d8669ff12b7c4',
];

contract('TransferValidatorV3', function([owner, minter, sender, relayer, witness1, witness2, witness3, witness4, cashier, receiver]) {
    beforeEach(async function() {
        this.witnessList = await WitnessList.new();
        this.validator = await TransferValidator.new(this.witnessList.address);
        this.unwrapper = await Unwrapper.new();
        this.minterPool = await MinterPool.new();
        this.mintableToken = await ShadowToken.new(this.minterPool.address, "0x0000000000000000000000000000000000000000", "token to mint", "mt", 18, {from: minter});
        this.mintableTokenList = await TokenList.new();
        await this.mintableTokenList.addToken(this.mintableToken.address, 1, 100000);
        await this.validator.addPair(this.mintableTokenList.address, this.minterPool.address);
        await this.unwrapper.addWhitelist(this.validator.address);
        await this.validator.addReceiver(this.unwrapper.address);
        await this.witnessList.addWitness(witness1);
    });
    it('is not reciever', async function() {
        await this.minterPool.transferOwnership(this.validator.address);
        const payload = "0x";
        const key = await this.validator.generateKey(cashier, this.mintableToken.address, 321, sender, receiver, 12345, payload);
        const signature = await Account.sign(key, witnessPrivateKeys[0]);
        assert.equal(await this.validator.settles(key), 0);
        const tx = await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature, payload);
        assert.notEqual(await this.validator.settles(key), 0);
        assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
        assert.equal(tx.logs.length, 1);
        assert.equal(tx.logs[0].event, "Settled");
        assert.equal(tx.logs[0].args.key, key);
        assert.equal(tx.logs[0].args.witnesses.length, 1);
        assert.equal(tx.logs[0].args.witnesses[0], witness1);
    })
    it('is reciever', async function() {
        this.crosschainERC20 = await CrosschainERC20.new(this.mintableToken.address, this.minterPool.address, "token to mint", "cerc20", 18);
        await this.minterPool.mint(this.mintableToken.address, this.crosschainERC20.address, 30000);
        await this.minterPool.transferOwnership(this.validator.address);
        await this.mintableTokenList.addToken(this.crosschainERC20.address, 1, 100000);
        const payload = AbiCoder.defaultAbiCoder().encode(["address"], [receiver]);
        const key = await this.validator.generateKey(cashier, this.crosschainERC20.address, 321, sender, this.unwrapper.address, 12345, payload);
        const signature = await Account.sign(key, witnessPrivateKeys[0]);
        assert.equal(await this.validator.settles(key), 0);
        const tx = await this.validator.submit(cashier, this.crosschainERC20.address, 321, sender, this.unwrapper.address, 12345, signature, payload);
        assert.notEqual(await this.validator.settles(key), 0);
        assert.equal(await this.mintableToken.balanceOf(this.unwrapper.address), 0);
        assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
        assert.equal(tx.logs.length, 1);
        assert.equal(tx.logs[0].event, "Settled");
        assert.equal(tx.logs[0].args.key, key);
        assert.equal(tx.logs[0].args.witnesses.length, 1);
        assert.equal(tx.logs[0].args.witnesses[0], witness1);
    })
})