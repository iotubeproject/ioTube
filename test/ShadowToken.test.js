const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('ShadowToken', function([owner, minter, stranger, fakeTokenAddress]) {
    beforeEach(async function() {
        this.shadowToken = await ShadowToken.new(minter, fakeTokenAddress, "name", "symbol", 12);
    });
    it('check values', async function() {
        assert.equal(await this.shadowToken.coToken(), fakeTokenAddress);
        assert.equal(await this.shadowToken.name(), "name");
        assert.equal(await this.shadowToken.symbol(), "symbol");
        assert.equal(await this.shadowToken.decimals(), 12);
    });
    it('is a minter', async function() {
        assert.equal(await this.shadowToken.minter(), minter);
        assert.equal(await this.shadowToken.balanceOf(owner), 0);
        await this.shadowToken.mint(owner, 10000000000, {from: minter});
        assert.equal(await this.shadowToken.balanceOf(owner), 10000000000);
        await this.shadowToken.burn(1);
        assert.equal(await this.shadowToken.balanceOf(owner), 9999999999);
    });
    it('not a minter', async function() {
        assert.equal(await this.shadowToken.balanceOf(owner), 0);
        await assertAsyncThrows(this.shadowToken.mint(owner, 1234567, {from: stranger}));
        assert.equal(await this.shadowToken.balanceOf(owner), 0);
    });
});