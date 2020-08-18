module.exports.assertAsyncThrows = async function(promise) {
    try {
        await promise;
    } catch (_) {
        return;
    }
    assert.fail();
}