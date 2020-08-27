const MinterPool = artifacts.require('MinterPool');
const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('MinterPool', function([owner, stranger, fakeTokenAddress]) {
    beforeEach(async function() {
        this.minterPool = await MinterPool.new();
        this.shadowToken = await ShadowToken.new(this.minterPool.address, fakeTokenAddress, "name", "symbol");
    });
    it('check minter', async function() {
        assert.equal(await this.shadowToken.minter(), this.minterPool.address);
    });
    it('is owner', async function() {
        assert.equal(await this.shadowToken.balanceOf(owner), 0);
        await this.minterPool.mint(this.shadowToken.address, owner, 10000000000);
        assert.equal(await this.shadowToken.balanceOf(owner), 10000000000);
    });
    it('not owner', async function() {
        assert.equal(await this.shadowToken.balanceOf(stranger), 0);
        await assertAsyncThrows(this.minterPool.mint(this.shadowToken.address, stranger, 1234567, {from: stranger}));
        assert.equal(await this.shadowToken.balanceOf(stranger), 0);
    });
});