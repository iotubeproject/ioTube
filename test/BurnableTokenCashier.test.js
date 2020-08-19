const BurnableTokenCashier = artifacts.require('BurnableTokenCashier');
const TokenList = artifacts.require('TokenList');
const MinterPool = artifacts.require('MinterPool');
const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('BurnableTokenCashier', function([owner, minter, sender, receiver, fakeTokenAddress1, fakeTokenAddress2]) {
    beforeEach(async function() {
        // use shadow token as burnable erc20 token
        this.shadowToken1 = await ShadowToken.new(minter, fakeTokenAddress1);
        this.shadowToken2 = await ShadowToken.new(minter, fakeTokenAddress2);
        await this.shadowToken1.mint(sender, 10000000000, {from: minter});
        assert.equal(await this.shadowToken1.balanceOf(sender), 10000000000);
        await this.shadowToken2.mint(sender, 10000000000, {from: minter});
        assert.equal(await this.shadowToken2.balanceOf(sender), 10000000000);
        this.tokenList = await TokenList.new();
        this.cashier = await BurnableTokenCashier.new(this.tokenList.address);
    });
    it('deposit not-in-list token', async function() {
        await this.shadowToken1.approve(this.cashier.address, 10000);
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken1.address, 10000, {from: sender}));
    });
    it('deposit with no quota', async function() {
        await this.tokenList.addToken(this.shadowToken1.address, 10, 1000);
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken1.address, 10000, {from: sender}));
    });
    it('deposit amount is invalid', async function() {
        await this.tokenList.addToken(this.shadowToken1.address, 10, 1000);
        await this.shadowToken1.approve(this.cashier.address, 10, {from: sender});
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken1.address, 1000, {from: sender}));
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken1.address, 1, {from: sender}));
        await this.shadowToken1.approve(this.cashier.address, 10000, {from: sender});
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken1.address, 10000, {from: sender}));
    });
    it('deposit', async function() {
        await this.tokenList.addToken(this.shadowToken1.address, 10, 1000);
        assert.equal(await this.cashier.count(this.shadowToken1.address), 0);
        await this.shadowToken1.approve(this.cashier.address, 10000, {from: sender});
        await this.cashier.deposit(this.shadowToken1.address, 1000, {from: sender});
        assert.equal(await this.shadowToken1.balanceOf(sender), 9999999000);
        assert.equal(await this.cashier.count(this.shadowToken1.address), 1);
        await this.cashier.depositTo(this.shadowToken1.address, receiver, 500, {from: sender});
        assert.equal(await this.shadowToken1.balanceOf(sender), 9999998500);
        assert.equal(await this.cashier.count(this.shadowToken1.address), 2);
        const records = await this.cashier.getRecords(this.shadowToken1.address, 0, 10);
        assert.equal(records.customers_.length, 2);
        assert.equal(records.customers_[0], sender);
        assert.equal(records.customers_[1], sender);
        assert.equal(records.receivers_.length, 2);
        assert.equal(records.receivers_[0], sender);
        assert.equal(records.receivers_[1], receiver);
        assert.equal(records.amounts_.length, 2);
        assert.equal(records.amounts_[0], 1000);
        assert.equal(records.amounts_[1], 500);
        assert.equal(records.fees_.length, 2);
        assert.equal(records.fees_[0], 0);
        assert.equal(records.fees_[1], 0);
    });
});