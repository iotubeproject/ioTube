const ShadowToken = artifacts.require('ShadowToken');
const WitnessList = artifacts.require('WitnessList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('WitnessList', function([owner, stranger, witness1, witness2, witness3]) {
    beforeEach(async function() {
        this.witnessList = await WitnessList.new();
        assert.equal(await this.witnessList.numOfActive(), 0);
    });
    it('witness not in list', async function() {
        assert.equal(await this.witnessList.isAllowed(witness1), false);
    });
    it('add witness', async function() {
        await this.witnessList.addWitness(witness1);
        assert.equal(await this.witnessList.isAllowed(witness1), true);
        assert.equal(await this.witnessList.numOfActive(), 1);
        let {count_: count1, items_: items1} = await this.witnessList.getActiveItems(0, 1);
        assert.equal(count1, 1);
        assert.equal(items1[0], witness1);
        await this.witnessList.addWitness(witness2);
        assert.equal(await this.witnessList.isAllowed(witness2), true);
        assert.equal(await this.witnessList.numOfActive(), 2);
        let {count_: count2, items_: items2} = await this.witnessList.getActiveItems(1, 1);
        assert.equal(count2, 1);
        assert.equal(items2[0], witness2);
        await this.witnessList.addWitness(witness3);
        assert.equal(await this.witnessList.isAllowed(witness3), true);
        assert.equal(await this.witnessList.numOfActive(), 3);
        let {count_: count3, items_: items3} = await this.witnessList.getActiveItems(2, 1);
        assert.equal(count3, 1);
        assert.equal(items3[0], witness3);
        await this.witnessList.removeWitness(witness2);
        assert.equal(await this.witnessList.isAllowed(witness2), false);
        assert.equal(await this.witnessList.numOfActive(), 2);
    });
    it('switch witness', async function() {
        await this.witnessList.addWitness(witness1);
        await this.witnessList.addWitness(witness2);
        assert.equal(await this.witnessList.isAllowed(witness1), true);
        assert.equal(await this.witnessList.isAllowed(witness2), true);
        assert.equal(await this.witnessList.isAllowed(stranger), false);
        assert.equal(await this.witnessList.isAllowed(witness3), false);
        await assertAsyncThrows(this.witnessList.switchWitness(witness3, {from: stranger}));
        assert.equal(await this.witnessList.isAllowed(witness1), true);
        assert.equal(await this.witnessList.isAllowed(witness2), true);
        assert.equal(await this.witnessList.isAllowed(stranger), false);
        assert.equal(await this.witnessList.isAllowed(witness3), false);
        await assertAsyncThrows(this.witnessList.switchWitness(witness1, {from: witness2}));
        assert.equal(await this.witnessList.isAllowed(witness1), true);
        assert.equal(await this.witnessList.isAllowed(witness2), true);
        assert.equal(await this.witnessList.isAllowed(stranger), false);
        assert.equal(await this.witnessList.isAllowed(witness3), false);
        await this.witnessList.switchWitness(witness3, {from: witness2});
        assert.equal(await this.witnessList.isAllowed(witness1), true);
        assert.equal(await this.witnessList.isAllowed(witness2), false);
        assert.equal(await this.witnessList.isAllowed(stranger), false);
        assert.equal(await this.witnessList.isAllowed(witness3), true);
    });
});