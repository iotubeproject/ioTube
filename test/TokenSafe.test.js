const TetherToken = artifacts.require('TetherToken');
const TokenSafe = artifacts.require('TokenSafe');
const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TokenSafe', function([owner, miner, stranger, alice, fakeTokenAddress]) {
    beforeEach(async function() {
        this.tokenSafe = await TokenSafe.new();
        assert.equal(await this.tokenSafe.owner(), owner);
        this.standardToken = await ShadowToken.new(miner, fakeTokenAddress, "name", "symbol", 18);
        this.nonStandardToken = await TetherToken.new(300, "name", "symbol", 18, {from: owner})
        await this.standardToken.mint(this.tokenSafe.address, 10000000000, {from: miner});
        assert.equal(await this.standardToken.balanceOf(this.tokenSafe.address), 10000000000);
        assert.equal(await this.nonStandardToken.balanceOf(owner), 300);
        await this.nonStandardToken.transfer(this.tokenSafe.address, 100, {from: owner})
        assert.equal(await this.nonStandardToken.balanceOf(this.tokenSafe.address), 100);
    });
    describe('not owner', function() {
        it('stardard token', async function() {
            assert.equal(await this.standardToken.balanceOf(stranger), 0);
            await assertAsyncThrows(this.tokenSafe.mint(this.standardToken.address, stranger, 1234567, {from: stranger}));
            assert.equal(await this.standardToken.balanceOf(this.tokenSafe.address), 10000000000);
            assert.equal(await this.standardToken.balanceOf(stranger), 0);
        });
        it('non-stardard token', async function() {
            assert.equal(await this.nonStandardToken.balanceOf(stranger), 0);
            await assertAsyncThrows(this.tokenSafe.mint(this.nonStandardToken.address, stranger, 12, {from: stranger}));
            assert.equal(await this.nonStandardToken.balanceOf(this.tokenSafe.address), 100);
            assert.equal(await this.nonStandardToken.balanceOf(stranger), 0);
        });
    });
    describe('owner', function() {
        it('standard token', async function() {
            assert.equal(await this.standardToken.balanceOf(alice), 0);
            await this.tokenSafe.mint(this.standardToken.address, alice, 1000000000);
            assert.equal(await this.standardToken.balanceOf(this.tokenSafe.address), 9000000000);
            assert.equal(await this.standardToken.balanceOf(alice), 1000000000);
        });
        it('non-starndard token ', async function() {
            await this.tokenSafe.mint(this.nonStandardToken.address, alice, 100)
            assert.equal(await this.nonStandardToken.balanceOf(this.tokenSafe.address), 0)
            assert.equal(await this.nonStandardToken.balanceOf(alice), 100)
        });
    });
});