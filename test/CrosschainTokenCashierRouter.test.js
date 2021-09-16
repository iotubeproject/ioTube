const CrosschainTokenCashierRouter = artifacts.require('CrosschainTokenCashierRouter');
const TokenCashier = artifacts.require('TokenCashier');
const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const CC = artifacts.require('CC');
const WETH = artifacts.require('WETH9');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('CrosschainTokenCashierRouter', function([owner, minter, sender, receiver, fakeTokenAddress1, fakeTokenAddress2]) {
    beforeEach(async function() {
        // use shadow token as burnable erc20 token
        this.weth = await WETH.new();
        this.cc = await CC.new(this.weth.address, minter, "name", "symbol", 18);
        this.mintableTokenList = await TokenList.new();
        this.standardTokenList = await TokenList.new();
        this.tokenSafe = await TokenSafe.new();
        this.cashier = await TokenCashier.new(
            this.weth.address,
            [this.mintableTokenList.address, this.standardTokenList.address],
            ['0x0000000000000000000000000000000000000000', this.tokenSafe.address],
        );
        this.router = await CrosschainTokenCashierRouter.new(this.cashier.address);
        await this.weth.deposit({from: sender, value: 10000000000})
        assert.equal(await this.weth.balanceOf(sender), 10000000000);
    });
    describe("deposit", function() {
        describe("burn", function() {
            beforeEach(async function() {
                await this.mintableTokenList.addToken(this.cc.address, 10, 1000);
                assert.equal(await this.mintableTokenList.isAllowed(this.cc.address), true);
            });
            describe("enough quota", function() {
                beforeEach(async function() {
                    await this.weth.approve(this.router.address, 10000, {from: sender});
                    await this.router.approveCrosschainToken(this.cc.address);
                });
                describe('success', function() {
                    beforeEach(async function() {
                        assert.equal(await this.cashier.count(this.cc.address), 0);
                        assert.equal(await web3.eth.getBalance(this.cashier.address), 0);
                    });
                    it('no fee', async function() {
                        const depositResponse = await this.router.depositTo(this.cc.address, sender, 1000, {from: sender});
                        assert.equal(await this.weth.balanceOf(sender), 9999999000);
                        assert.equal(await this.cashier.count(this.cc.address), 1);
                        assert.equal(depositResponse.receipt.rawLogs[5].address, this.cashier.address);
                        assert.equal(depositResponse.receipt.rawLogs[5].topics[1].substring(26, 66), this.cc.address.substring(2).toLowerCase());
                        assert.equal(depositResponse.receipt.rawLogs[5].topics[2], 1);
                    });
                });
            });
        });
    });
});