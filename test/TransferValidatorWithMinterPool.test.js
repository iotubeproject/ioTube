const MinterPool = artifacts.require('MinterPool');
const TransferValidatorWithMinterPool = artifacts.require('TransferValidatorWithMinterPool');
const TokenList = artifacts.require('TokenList');
const ShadowToken = artifacts.require('ShadowToken');
const WitnessList = artifacts.require('WitnessList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TransferValidatorWithMinterPool', function([owner, minter, sender, receiver, witness1, witness2, witness3, witness4, fakeTokenAddress]) {
    beforeEach(async function() {
        this.minterPool = await MinterPool.new();
        // use shadow token as mintable erc20 token
        this.shadowToken = await ShadowToken.new(this.minterPool.address, fakeTokenAddress);
        this.tokenList = await TokenList.new();
        this.witnessList = await WitnessList.new();
        this.validator = await TransferValidatorWithMinterPool.new(10, this.minterPool.address, this.tokenList.address, this.witnessList.address);
        await this.minterPool.transferOwnership(this.validator.address);
    });
    it('witness not in list', async function() {
        await assertAsyncThrows(this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1}));
    });
    it('token not in list', async function() {
        await this.witnessList.addWitness(witness1);
        await assertAsyncThrows(this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1}));
    });
    it('one witness', async function() {
        await this.witnessList.addWitness(witness1);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('two witnesss', async function() {
        await this.witnessList.addWitness(witness1);
        await this.witnessList.addWitness(witness2);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness2});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('three witnesss', async function() {
        await this.witnessList.addWitness(witness1);
        await this.witnessList.addWitness(witness2);
        await this.witnessList.addWitness(witness3);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness2});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness3});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('four witnesss', async function() {
        await this.witnessList.addWitness(witness1);
        await this.witnessList.addWitness(witness2);
        await this.witnessList.addWitness(witness3);
        await this.witnessList.addWitness(witness4);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness2});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.submit(this.shadowToken.address, 0, sender, receiver, 12345, {from: witness3});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('upgrade', async function() {
        assert.equal(await this.minterPool.owner(), this.validator.address);
        await assertAsyncThrows(this.validator.upgrade(owner, {from: sender}));
        assert.equal(await this.minterPool.owner(), this.validator.address);
        await this.validator.upgrade(owner);
        assert.equal(await this.minterPool.owner(), owner);
    });
});