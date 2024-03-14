const TokenCashier = artifacts.require('TokenCashier');
const TokenList = artifacts.require('TokenList');
const TokenSafe = artifacts.require('TokenSafe');
const ShadowToken = artifacts.require('ShadowToken');
const WETH = artifacts.require('WETH9');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('TokenCashier', function([owner, minter, sender, receiver, fakeTokenAddress1, fakeTokenAddress2]) {
    beforeEach(async function() {
        // use shadow token as burnable erc20 token
        this.shadowToken = await ShadowToken.new(minter, fakeTokenAddress1, "name", "symbol", 18);
        await this.shadowToken.mint(sender, 10000000000, {from: minter});
        assert.equal(await this.shadowToken.balanceOf(sender), 10000000000);
        this.mintableTokenList = await TokenList.new();
        this.standardTokenList = await TokenList.new();
        this.weth = await WETH.new();
        this.tokenSafe = await TokenSafe.new();
        this.cashier = await TokenCashier.new(
            this.weth.address,
            [this.mintableTokenList.address, this.standardTokenList.address],
            ['0x0000000000000000000000000000000000000000', this.tokenSafe.address],
        );
    });
    it('deposit fee', async function() {
        assert.equal(await this.cashier.depositFee(), 0);
        await this.cashier.setDepositFee(12345678);
        assert.equal(await this.cashier.depositFee(), 12345678);        
    });
    describe("deposit", function() {
        it('not-in-list token', async function() {
            await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 10000, {from: sender}));
        });
        describe("burn", function() {
            beforeEach(async function() {
                await this.mintableTokenList.addToken(this.shadowToken.address, 10, 1000);
                assert.equal(await this.mintableTokenList.isAllowed(this.shadowToken.address), true);
            });
            it('no quota', async function() {
                await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 100, {from: sender}));
            });
            describe("enough quota", function() {
                beforeEach(async function() {
                    await this.shadowToken.approve(this.cashier.address, 10000, {from: sender});
                });
                it('invalid amount', async function() {
                    await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 9, {from: sender}));
                    await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 1001, {from: sender}));
                });
                it('insufficient deposit fee', async function() {
                    await this.cashier.setDepositFee(1234);
                    await assertAsyncThrows(this.cashier.deposit(this.shadowToken.address, 1000, {from: sender, value: 12}));
                });
                describe('success', function() {
                    beforeEach(async function() {
                        assert.equal(await this.cashier.count(this.shadowToken.address), 0);
                        assert.equal(await web3.eth.getBalance(this.cashier.address), 0);
                    });
                    it('no fee', async function() {
                        const depositResponse = await this.cashier.deposit(this.shadowToken.address, 1000, {from: sender});
                        assert.equal(await this.shadowToken.balanceOf(sender), 9999999000);
                        assert.equal(await this.cashier.count(this.shadowToken.address), 1);
                        assert.equal(depositResponse.logs.length, 1);
                        assert.equal(depositResponse.logs[0].event, "Receipt");
                        assert.equal(depositResponse.logs[0].address, this.cashier.address);
                        assert.equal(depositResponse.logs[0].args.token, this.shadowToken.address);
                        assert.equal(depositResponse.logs[0].args.id, 1);
                        assert.equal(depositResponse.logs[0].args.sender, sender);
                        assert.equal(depositResponse.logs[0].args.recipient, sender);
                        assert.equal(depositResponse.logs[0].args.amount, 1000);
                        assert.equal(depositResponse.logs[0].args.fee, 0);
                    });
                    it('depositTo with fee', async function() {
                        await this.cashier.setDepositFee(1234);
                        const depositResponse = await this.cashier.depositTo(this.shadowToken.address, receiver, 500, {from: sender, value: 1234});
                        assert.equal(await this.shadowToken.balanceOf(sender), 9999999500);
                        assert.equal(await this.cashier.count(this.shadowToken.address), 1);
                        assert.equal(await web3.eth.getBalance(this.cashier.address), 1234);
                        assert.equal(depositResponse.logs.length, 1);
                        assert.equal(depositResponse.logs[0].event, "Receipt");
                        assert.equal(depositResponse.logs[0].address, this.cashier.address);
                        assert.equal(depositResponse.logs[0].args.token, this.shadowToken.address);
                        assert.equal(depositResponse.logs[0].args.id, 1);
                        assert.equal(depositResponse.logs[0].args.sender, sender);
                        assert.equal(depositResponse.logs[0].args.recipient, receiver);
                        assert.equal(depositResponse.logs[0].args.amount, 500);
                        assert.equal(depositResponse.logs[0].args.fee, 1234);
                        const initBalance = await web3.eth.getBalance(owner);
                        const withdrawResponse = await this.cashier.withdraw();
                        const newBalance = await web3.eth.getBalance(owner);
                        assert.equal(
                            web3.utils.toBN(initBalance).sub(web3.utils.toBN(withdrawResponse.receipt.gasUsed * withdrawResponse.receipt.effectiveGasPrice)).toString(),
                            web3.utils.toBN(newBalance).sub(web3.utils.toBN(1234)).toString(),
                        );
                    });
                });
            });
        });
        it("deposit into safe", async function() {
            await this.standardTokenList.addToken(this.shadowToken.address, 10, 1000);
            await this.standardTokenList.addToken(this.weth.address, 1, 100000);
            assert.equal(await this.standardTokenList.isAllowed(this.shadowToken.address), true);
            assert.equal(await this.cashier.count(this.shadowToken.address), 0);
            assert.equal(await web3.eth.getBalance(this.cashier.address), 0);
            await this.shadowToken.approve(this.cashier.address, 10000, {from: sender});
            let depositResponse = await this.cashier.deposit(this.shadowToken.address, 1000, {from: sender});
            assert.equal(await this.shadowToken.balanceOf(sender), 9999999000);
            assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 1000);
            assert.equal(await this.cashier.count(this.shadowToken.address), 1);
            assert.equal(await web3.eth.getBalance(this.cashier.address), 0);
            assert.equal(depositResponse.logs.length, 1);
            assert.equal(depositResponse.logs[0].event, "Receipt");
            assert.equal(depositResponse.logs[0].address, this.cashier.address);
            assert.equal(depositResponse.logs[0].args.token, this.shadowToken.address);
            assert.equal(depositResponse.logs[0].args.id, 1);
            assert.equal(depositResponse.logs[0].args.sender, sender);
            assert.equal(depositResponse.logs[0].args.recipient, sender);
            assert.equal(depositResponse.logs[0].args.amount, 1000);
            assert.equal(depositResponse.logs[0].args.fee, 0);
            await this.cashier.setDepositFee(123);
            depositResponse = await this.cashier.depositTo(this.shadowToken.address, receiver, 500, {from: sender, value: 1234});
            assert.equal(await this.shadowToken.balanceOf(sender), 9999998500);
            assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 1500);
            assert.equal(await this.cashier.count(this.shadowToken.address), 2);
            assert.equal(await web3.eth.getBalance(this.cashier.address), 1234);
            assert.equal(depositResponse.logs.length, 1);
            assert.equal(depositResponse.logs[0].event, "Receipt");
            assert.equal(depositResponse.logs[0].address, this.cashier.address);
            assert.equal(depositResponse.logs[0].args.token, this.shadowToken.address);
            assert.equal(depositResponse.logs[0].args.id, 2);
            assert.equal(depositResponse.logs[0].args.sender, sender);
            assert.equal(depositResponse.logs[0].args.recipient, receiver);
            assert.equal(depositResponse.logs[0].args.amount, 500);
            assert.equal(depositResponse.logs[0].args.fee, 1234);
            const balanceBeforeDeposit = await web3.eth.getBalance(sender);
            // const gasPrice = await web3.eth.getGasPrice();
            depositResponse = await this.cashier.depositTo("0x0000000000000000000000000000000000000000", receiver, 50000, {from: sender, value: 54321});
            assert.equal(await web3.eth.getBalance(sender), web3.utils.toBN(balanceBeforeDeposit).sub(web3.utils.toBN(depositResponse.receipt.gasUsed * depositResponse.receipt.effectiveGasPrice + 54321)));
            assert.equal(await this.weth.balanceOf(this.tokenSafe.address), 50000);
            assert.equal(await this.cashier.count(this.weth.address), 1);
            assert.equal(await web3.eth.getBalance(this.cashier.address), 5555);
            assert.equal(depositResponse.logs.length, 1);
            assert.equal(depositResponse.logs[0].event, "Receipt");
            assert.equal(depositResponse.logs[0].address, this.cashier.address);
            assert.equal(depositResponse.logs[0].args.token, this.weth.address);
            assert.equal(depositResponse.logs[0].args.id, 1);
            assert.equal(depositResponse.logs[0].args.sender, sender);
            assert.equal(depositResponse.logs[0].args.recipient, receiver);
            assert.equal(depositResponse.logs[0].args.amount, 50000);
            assert.equal(depositResponse.logs[0].args.fee, 4321);
        });
    });
});