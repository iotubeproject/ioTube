const TokenCashierWithSafe = artifacts.require('TokenCashierWithSafe');
const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const ShadowToken = artifacts.require('ShadowToken');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TokenCashierWithSafe', function([owner, minter, sender, receiver, fakeTokenAddress1, fakeTokenAddress2]) {
    beforeEach(async function() {
        this.tokenSafe = await TokenSafe.new();
        // use shadow token as standard erc20 token
        this.shadowToken = await ShadowToken.new(minter, fakeTokenAddress1);
        await this.shadowToken.mint(sender, 10000000000, {from: minter});
        assert.equal(await this.shadowToken.balanceOf(sender), 10000000000);
        this.tokenList = await TokenList.new();
        this.cashier = await TokenCashierWithSafe.new(this.tokenList.address, this.tokenSafe.address);
    });
    it('deposit not-in-list token', async function() {
        await this.shadowToken.approve(this.cashier.address, 10000);
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 10000, {from: sender}));
    });
    it('deposit with no quota', async function() {
        await this.tokenList.addToken(this.shadowToken.address, 10, 1000);
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 10000, {from: sender}));
    });
    it('deposit amount is invalid', async function() {
        await this.tokenList.addToken(this.shadowToken.address, 10, 1000);
        await this.shadowToken.approve(this.cashier.address, 10, {from: sender});
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 1000, {from: sender}));
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 1, {from: sender}));
        await this.shadowToken.approve(this.cashier.address, 10000, {from: sender});
        await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 10000, {from: sender}));
    });
    it('deposit', async function() {
        await this.tokenList.addToken(this.shadowToken.address, 10, 1000);
        assert.equal(await this.cashier.count(this.shadowToken.address), 0);
        await this.shadowToken.approve(this.cashier.address, 10000, {from: sender});
        await this.cashier.deposit(this.shadowToken.address, 1000, {from: sender});
        assert.equal(await this.shadowToken.balanceOf(sender), 9999999000);
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 1000);
        assert.equal(await this.cashier.count(this.shadowToken.address), 1);
        await this.cashier.depositTo(this.shadowToken.address, receiver, 500, {from: sender});
        assert.equal(await this.shadowToken.balanceOf(sender), 9999998500);
        assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 1500);
        assert.equal(await this.cashier.count(this.shadowToken.address), 2);
        const records = await this.cashier.getRecords(this.shadowToken.address, 0, 10);
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