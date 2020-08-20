const MinterPool = artifacts.require('MinterPool');
const TransferValidatorWithMinterPool = artifacts.require('TransferValidatorWithMinterPool');
const TokenList = artifacts.require('TokenList');
const ShadowToken = artifacts.require('ShadowToken');
const VoterList = artifacts.require('VoterList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TransferValidatorWithMinterPool', function([owner, minter, sender, receiver, voter1, voter2, voter3, voter4, fakeTokenAddress]) {
    beforeEach(async function() {
        this.minterPool = await MinterPool.new();
        // use shadow token as standard erc20 token
        this.shadowToken = await ShadowToken.new(this.minterPool.address, fakeTokenAddress);
        this.tokenList = await TokenList.new();
        this.voterList = await VoterList.new();
        this.validator = await TransferValidatorWithMinterPool.new(10, this.minterPool.address, this.tokenList.address, this.voterList.address);
        await this.minterPool.transferOwnership(this.validator.address);
    });
    it('voter not in list', async function() {
        await assertAsyncThrows(this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1}));
    });
    it('token not in list', async function() {
        await this.voterList.addVoter(voter1);
        await assertAsyncThrows(this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1}));
    });
    it('one voter', async function() {
        await this.voterList.addVoter(voter1);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('two voters', async function() {
        await this.voterList.addVoter(voter1);
        await this.voterList.addVoter(voter2);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter2});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('three voters', async function() {
        await this.voterList.addVoter(voter1);
        await this.voterList.addVoter(voter2);
        await this.voterList.addVoter(voter3);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter2});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter3});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
    it('four voters', async function() {
        await this.voterList.addVoter(voter1);
        await this.voterList.addVoter(voter2);
        await this.voterList.addVoter(voter3);
        await this.voterList.addVoter(voter4);
        await this.tokenList.addToken(this.shadowToken.address, 1, 100000);
        const key = await this.validator.generateKey(this.shadowToken.address, 0, sender, receiver, 12345);
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter1});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter2});
        assert.equal(await this.validator.settled(key), false);
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter3});
        assert.equal(await this.validator.settled(key), true);
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
    });
});