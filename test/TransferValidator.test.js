const TokenSafe = artifacts.require('TokenSafe');
const MinterPool = artifacts.require('MinterPool');
const TransferValidator = artifacts.require('TransferValidator');
const TokenList = artifacts.require('TokenList');
const ShadowToken = artifacts.require('ShadowToken');
const WitnessList = artifacts.require('WitnessList');
const ethjs = require('eth-lib');
const {assertAsyncThrows} = require('./assert-async-throws');

const witnessPrivateKeys = [
    '0x388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418',
    '0x659cbb0e2411a44db63778987b1e22153c086a95eb6b18bdf89de078917abc63',
    '0x82d052c865f5763aad42add438569276c00d3d88a2d062d36b2bae914d58b8c8',
    '0xaa3680d5d48a8283413f7a108367c7299ca73f553735860a87b08f39395618b7',
    '0x0f62d96d6675f32685bbdb8ac13cda7c23436f63efbb9d07700d8669ff12b7c4',
];

contract('TransferValidator', function([owner, minter, sender, relayer, witness1, witness2, witness3, witness4, cashier, receiver]) {
    beforeEach(async function() {
        this.tokenSafe = await TokenSafe.new();
        this.minterPool = await MinterPool.new();
        // use shadow token as standard erc20 token
        this.mintableToken = await ShadowToken.new(this.minterPool.address, "0x0000000000000000000000000000000000000000", "token to mint", "mt", 18);
        this.shadowToken = await ShadowToken.new(minter, "0x0000000000000000000000000000000000000000", "token to withdraw", "wt", 18);
        await this.shadowToken.mint(this.tokenSafe.address, 100000000, {from: minter});
        this.mintableTokenList = await TokenList.new();
        this.standardTokenList = await TokenList.new();
        this.witnessList = await WitnessList.new();
        this.validator = await TransferValidator.new(
            [this.tokenSafe.address, this.minterPool.address],
            [this.standardTokenList.address, this.mintableTokenList.address],
            this.witnessList.address,
        );
        await this.tokenSafe.transferOwnership(this.validator.address);
        await this.minterPool.transferOwnership(this.validator.address);
    });
    it('token not in list', async function() {
        await assertAsyncThrows(this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, "", {from: relayer}));
    });
    it('upgrade', async function() {
        assert.equal(await this.tokenSafe.owner(), this.validator.address);
        assert.equal(await this.minterPool.owner(), this.validator.address);
        await assertAsyncThrows(this.validator.upgrade(owner, {from: sender}));
        assert.equal(await this.tokenSafe.owner(), this.validator.address);
        assert.equal(await this.minterPool.owner(), this.validator.address);
        await this.validator.upgrade(owner);
        assert.equal(await this.tokenSafe.owner(), owner);
        assert.equal(await this.minterPool.owner(), owner);
    });
    describe("withdraw standard token", function() {
        beforeEach(async function() {
            await this.standardTokenList.addToken(this.shadowToken.address, 1, 100000);
        });
        it('invalid signature length', async function() {
            await assertAsyncThrows(this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, "", {from: relayer}));
        });
        describe("different numbers of witnesses", function() {
            let key;
            beforeEach(async function() {
                await this.standardTokenList.addToken(this.shadowToken.address, 1, 100000);
                key = await this.validator.generateKey(cashier, this.shadowToken.address, 0, sender, receiver, 12345);
                assert.equal(await this.validator.settles(key), 0);    
            });
            it('one witness', async function() {
                await this.witnessList.addWitness(witness1);
                const signature = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                await this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature, {from: relayer});
                assert.notEqual(await this.validator.settles(key), 0);
                assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
                assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
            });
            it('two witnesses', async function() {
                await this.witnessList.addWitness(witness1);
                await this.witnessList.addWitness(witness2);
                const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                await this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2), {from: relayer});
                assert.notEqual(await this.validator.settles(key), 0);
                assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
                assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);
            });
            describe('three witnesses', function() {
                beforeEach(async function() {
                    await this.witnessList.addWitness(witness1);
                    await this.witnessList.addWitness(witness2);
                    await this.witnessList.addWitness(witness3);
                });
                it("three valid signatures", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    await this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
                    assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);    
                });
                it("insufficient signatures", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    await assertAsyncThrows(this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2), {from: relayer}));
                });
                it("signature from invalid witness", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[3]);
                    await assertAsyncThrows(this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer}));
                });
            });
            describe('four witnesses', function() {
                beforeEach(async function() {
                    await this.witnessList.addWitness(witness1);
                    await this.witnessList.addWitness(witness2);
                    await this.witnessList.addWitness(witness3);
                    await this.witnessList.addWitness(witness4);
                });
                it("three submissions", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    await this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
                    assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);    
                });
                it("four submissions", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    const signature4 = await ethjs.Account.sign(key, witnessPrivateKeys[3]);
                    await this.validator.submit(cashier, this.shadowToken.address, 0, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2) + signature4.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.shadowToken.balanceOf(receiver), 12345);
                    assert.equal(await this.shadowToken.balanceOf(this.tokenSafe.address), 99987655);    
                });            
            });
        });
    });
    describe("withdraw mintable token", function() {
        beforeEach(async function() {
            await this.mintableTokenList.addToken(this.mintableToken.address, 1, 100000);
        });
        it('invalid signature length', async function() {
            await assertAsyncThrows(this.validator.submit(cashier, this.mintableToken.address, 0, sender, receiver, 12345, "", {from: relayer}));
        });
        describe("different numbers of witnesses", function() {
            let key;
            beforeEach(async function() {
                await this.mintableTokenList.addToken(this.mintableToken.address, 1, 100000);
                key = await this.validator.generateKey(cashier, this.mintableToken.address, 321, sender, receiver, 12345);
                assert.equal(await this.validator.settles(key), 0);    
                assert.equal(await this.mintableToken.balanceOf(receiver), 0);
            });
            it('one witness', async function() {
                await this.witnessList.addWitness(witness1);
                const signature = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature, {from: relayer});
                assert.notEqual(await this.validator.settles(key), 0);
                assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
            });
            it('two witnesses', async function() {
                await this.witnessList.addWitness(witness1);
                await this.witnessList.addWitness(witness2);
                const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2), {from: relayer});
                assert.notEqual(await this.validator.settles(key), 0);
                assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
            });
            describe('three witnesses', function() {
                beforeEach(async function() {
                    await this.witnessList.addWitness(witness1);
                    await this.witnessList.addWitness(witness2);
                    await this.witnessList.addWitness(witness3);
                });
                it("three valid signatures", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
                });
                it("insufficient signatures", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    await assertAsyncThrows(this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2), {from: relayer}));
                });
                it("signature from invalid witness", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[3]);
                    await assertAsyncThrows(this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer}));
                });
            });
            describe('four witnesses', function() {
                beforeEach(async function() {
                    await this.witnessList.addWitness(witness1);
                    await this.witnessList.addWitness(witness2);
                    await this.witnessList.addWitness(witness3);
                    await this.witnessList.addWitness(witness4);
                });
                it("three submissions", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
                });
                it("four submissions", async function() {
                    const signature1 = await ethjs.Account.sign(key, witnessPrivateKeys[0]);
                    const signature2 = await ethjs.Account.sign(key, witnessPrivateKeys[1]);
                    const signature3 = await ethjs.Account.sign(key, witnessPrivateKeys[2]);
                    const signature4 = await ethjs.Account.sign(key, witnessPrivateKeys[3]);
                    await this.validator.submit(cashier, this.mintableToken.address, 321, sender, receiver, 12345, signature1 + signature2.substr(2) + signature3.substr(2) + signature4.substr(2), {from: relayer});
                    assert.notEqual(await this.validator.settles(key), 0);
                    assert.equal(await this.mintableToken.balanceOf(receiver), 12345);
                });            
            });
        });
    });
});