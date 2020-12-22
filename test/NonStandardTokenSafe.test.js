const TetherToken = artifacts.require('TetherToken');
const NonStandardTokenSafe = artifacts.require('NonStandardTokenSafe');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('NonStandardTokenSafe', function([owner, alice]) {
    beforeEach(async function() {
        this.nonStandardToken = await TetherToken.new(300, "name", "symbol", 18, {from: owner})
        this.nonStandardTokenSafe = await NonStandardTokenSafe.new({from: owner})
    });
    it('withdrawToken for non starndard token ', async function() {
        assert.equal(await this.nonStandardToken.balanceOf(owner), 300)
        await this.nonStandardToken.transfer(this.nonStandardTokenSafe.address, 100, {from: owner})
        assert.equal(await this.nonStandardToken.balanceOf(this.nonStandardTokenSafe.address), 100)
        await this.nonStandardTokenSafe.withdrawToken(this.nonStandardToken.address, alice, 100)
        assert.equal(await this.nonStandardToken.balanceOf(this.nonStandardTokenSafe.address), 0)
        assert.equal(await this.nonStandardToken.balanceOf(alice), 100)
    });
});