const ShadowToken = artifacts.require('ShadowToken');
const VoterList = artifacts.require('VoterList');
const {assertAsyncThrows} = require('./assert-async-throws');

contract('VoterList', function([owner, stranger, voter1, voter2, voter3]) {
    beforeEach(async function() {
        this.voterList = await VoterList.new();
        assert.equal(await this.voterList.numOfAllowed(), 0);
    });
    it('voter not in list', async function() {
        assert.equal(await this.voterList.isAllowed(voter1), false);
    });
    it('add voter', async function() {
        await this.voterList.addVoter(voter1);
        assert.equal(await this.voterList.isAllowed(voter1), true);
        assert.equal(await this.voterList.numOfAllowed(), 1);
        assert.equal(await this.voterList.voters(0), voter1);
        await this.voterList.addVoter(voter2);
        assert.equal(await this.voterList.isAllowed(voter2), true);
        assert.equal(await this.voterList.numOfAllowed(), 2);
        assert.equal(await this.voterList.voters(1), voter2);
        await this.voterList.addVoter(voter3);
        assert.equal(await this.voterList.isAllowed(voter3), true);
        assert.equal(await this.voterList.numOfAllowed(), 3);
        assert.equal(await this.voterList.voters(2), voter3);
        await this.voterList.removeVoter(voter2);
        assert.equal(await this.voterList.isAllowed(voter2), false);
        assert.equal(await this.voterList.numOfAllowed(), 2);
    });
});