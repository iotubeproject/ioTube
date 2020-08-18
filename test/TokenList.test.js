const ShadowToken = artifacts.require('ShadowToken');
const TokenList = artifacts.require('TokenList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TokenList', function([owner, minter, stranger, fakeTokenAddress1, fakeTokenAddress2]) {
    beforeEach(async function() {
        this.tokenList = await TokenList.new();
        this.shadowToken1 = await ShadowToken.new(minter, fakeTokenAddress1);
        this.shadowToken2 = await ShadowToken.new(minter, fakeTokenAddress2);
        assert.equal(await this.tokenList.numOfAllowed(), 0);
    });
    it('token not in list', async function() {
        assert.equal(await this.tokenList.isAllowed(this.shadowToken1.address), false);
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 0);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 0);
        await assertAsyncThrows(this.tokenList.setMinAmount(this.shadowToken1.address, 1));
        await assertAsyncThrows(this.tokenList.setMaxAmount(this.shadowToken1.address, 10));
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 0);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 0);
    });
    it('add token then delete one', async function() {
        await assertAsyncThrows(this.tokenList.addToken(this.shadowToken1.address, 0, 9));
        assert.equal(await this.tokenList.isAllowed(this.shadowToken1.address), false);
        await assertAsyncThrows(this.tokenList.addToken(this.shadowToken1.address, 10, 9));
        assert.equal(await this.tokenList.isAllowed(this.shadowToken1.address), false);
        await this.tokenList.addToken(this.shadowToken1.address, 9, 100);
        assert.equal(await this.tokenList.isAllowed(this.shadowToken1.address), true);
        assert.equal(await this.tokenList.numOfAllowed(), 1);
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 9);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 100);
        await this.tokenList.addToken(this.shadowToken2.address, 99, 101);
        assert.equal(await this.tokenList.isAllowed(this.shadowToken2.address), true);
        assert.equal(await this.tokenList.numOfAllowed(), 2);
        assert.equal(await this.tokenList.minAmount(this.shadowToken2.address), 99);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken2.address), 101);
        await this.tokenList.removeToken(this.shadowToken2.address);
        assert.equal(await this.tokenList.isAllowed(this.shadowToken2.address), false);
        assert.equal(await this.tokenList.numOfAllowed(), 1);
        // update settings for active token
        await this.tokenList.setMinAmount(this.shadowToken1.address, 100);
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 100);
        await assertAsyncThrows(this.tokenList.setMinAmount(this.shadowToken1.address, 0));
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 100);
        await assertAsyncThrows(this.tokenList.setMinAmount(this.shadowToken1.address, 101));
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 100);
        await this.tokenList.setMaxAmount(this.shadowToken1.address, 200);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 200);
        await assertAsyncThrows(this.tokenList.setMaxAmount(this.shadowToken1.address, 99));
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 200);
        // update settings for inactive token
        await this.tokenList.setMinAmount(this.shadowToken2.address, 100);
        assert.equal(await this.tokenList.minAmount(this.shadowToken2.address), 100);
        await assertAsyncThrows(this.tokenList.setMinAmount(this.shadowToken2.address, 0));
        assert.equal(await this.tokenList.minAmount(this.shadowToken2.address), 100);
        await assertAsyncThrows(this.tokenList.setMinAmount(this.shadowToken2.address, 102));
        assert.equal(await this.tokenList.minAmount(this.shadowToken2.address), 100);
        await this.tokenList.setMaxAmount(this.shadowToken2.address, 200);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken2.address), 200);
        await assertAsyncThrows(this.tokenList.setMaxAmount(this.shadowToken2.address, 99));
        assert.equal(await this.tokenList.maxAmount(this.shadowToken2.address), 200);
    });
    it('add tokens', async function() {
        await assertAsyncThrows(this.tokenList.addTokens([this.shadowToken1.address, this.shadowToken2.address], [9, 100], [100, 99]));
        await assertAsyncThrows(this.tokenList.addTokens([this.shadowToken1.address, this.shadowToken2.address], [0, 10], [100, 99]));
        await this.tokenList.addTokens([this.shadowToken1.address, this.shadowToken2.address], [9, 10], [100, 99]);
        assert.equal(await this.tokenList.minAmount(this.shadowToken1.address), 9);
        assert.equal(await this.tokenList.minAmount(this.shadowToken2.address), 10);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken1.address), 100);
        assert.equal(await this.tokenList.maxAmount(this.shadowToken2.address), 99);
    });
});