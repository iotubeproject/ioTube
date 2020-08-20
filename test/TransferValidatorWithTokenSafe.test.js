const TokenSafe = artifacts.require('TokenSafe');
const TransferValidatorWithTokenSafe = artifacts.require('TransferValidatorWithTokenSafe');
const TokenList = artifacts.require('TokenList');
const ShadowToken = artifacts.require('ShadowToken');
const VoterList = artifacts.require('VoterList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TransferValidatorWithTokenSafe', function([owner, minter, sender, receiver, voter1, voter2, voter3, voter4, fakeTokenAddress]) {
    beforeEach(async function() {
        this.tokenSafe = await TokenSafe.new();
        // use shadow token as standard erc20 token
        this.shadowToken = await ShadowToken.new(minter, fakeTokenAddress);
        await this.shadowToken.mint(this.tokenSafe.address, 100000000, {from: minter});
        this.tokenList = await TokenList.new();
        this.voterList = await VoterList.new();
        this.validator = await TransferValidatorWithTokenSafe.new(10, this.tokenSafe.address, this.tokenList.address, this.voterList.address);
        await this.tokenSafe.transferOwnership(this.validator.address);
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
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
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
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
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
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
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
        await this.validator.vote(this.shadowToken.address, 0, sender, receiver, 12345, {from: voter4});
        assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
    });
    it('upgrade', async function() {
        assert.equal(await this.tokenSafe.owner(), this.validator.address);
        await assertAsyncThrows(this.validator.upgrade(owner, {from: sender}));
        assert.equal(await this.tokenSafe.owner(), this.validator.address);
        await this.validator.upgrade(owner);
        assert.equal(await this.tokenSafe.owner(), owner);
    });
});