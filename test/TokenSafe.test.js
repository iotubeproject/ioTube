const TokenSafe = artifacts.require('TokenSafe');
const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TokenSafe', function([owner, miner, stranger, fakeTokenAddress]) {
    beforeEach(async function() {
        this.tokenSafe = await TokenSafe.new();
        this.shadowToken = await ShadowToken.new(miner, fakeTokenAddress);
        await this.shadowToken.mint(this.tokenSafe.address, 10000000000, {from: miner});
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 10000000000);
    });
    it('is owner', async function() {
        assert.equal(await this.tokenSafe.owner(), owner);
    });
    it('is owner', async function() {
        assert.equal(await this.shadowToken.balanceOf(owner), 0);
        await this.tokenSafe.withdrawToken(this.shadowToken.address, owner, 1000000000);
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 9000000000);
        assert.equal(await this.shadowToken.balanceOf(owner), 1000000000);
    });
    it('not owner', async function() {
        assert.equal(await this.shadowToken.balanceOf(stranger), 0);
        await assertAsyncThrows(this.tokenSafe.withdrawToken(this.shadowToken.address, stranger, 1234567, {from: stranger}));
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 10000000000);
        assert.equal(await this.shadowToken.balanceOf(stranger), 0);
    });
});