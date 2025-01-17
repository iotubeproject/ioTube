const MockERC20 = artifacts.require("MockERC20");
const MockCashier = artifacts.require("MockCashier");
const EthereumHubPrepaid = artifacts.require("EthereumHubPrepaid");
const {ethers, AbiCoder} = require("ethers");
const {assertAsyncThrows} = require('./assert-async-throws');

const expectSwap = function(log, tokenIn, tokenOut, amountIn, amountOut, to) {
    assert.equal(log.event, "Swap");
    assert.equal(log.args.tokenIn, tokenIn);
    assert.equal(log.args.tokenOut, tokenOut);
    assert.equal(log.args.amountIn, amountIn);
    assert.equal(log.args.amountOut, amountOut);
    assert.equal(log.args.to, to);
}

contract("EthereumHubPrepaid", function([owner, operator, hacker, recipient]) {
  let unwrapper;
  let cashier;
  let token;
  let weth;

  beforeEach(async function() {
    // Deploy mock WETH
    weth = await MockERC20.new("Wrapped ETH", "WETH");

    // Deploy mock router
    cashier = await MockCashier.new();

    // Deploy mock token
    token = await MockERC20.new("Test Token", "TEST");

    // Deploy unwrapper
    unwrapper = await EthereumHubPrepaid.new(cashier.address, {value: 1000});
    await unwrapper.addOperator(operator);
    // Approve tokens
    await token.approve(unwrapper.address, ethers.MaxUint256);
    await token.approve(cashier.address, ethers.MaxUint256);
  });

  it('insufficient funds', async function() {
    await cashier.setDepositFee(10000);
    await assertAsyncThrows(unwrapper.onReceive(owner, token.address, ethers.parseEther("1"), "0x", {from: operator}));
  });

  it('not operator', async function() {
    await assertAsyncThrows(unwrapper.onReceive(hacker, token.address, ethers.parseEther("1"), "0x", {from: hacker}));
  });

  describe('deposit', function() {
    beforeEach(async function() {
      await cashier.setDepositFee(10);
    });
    it('should deposit tokens', async function() {
      const amountIn = ethers.parseEther("100");
      await token.transfer(unwrapper.address, amountIn);
      const payload = AbiCoder.defaultAbiCoder().encode(["address", "bytes"], [recipient, "0x"]);
      await unwrapper.onReceive(owner, token.address, amountIn, payload, {from: operator});
      assert.equal(await token.balanceOf(unwrapper.address), 0);
      assert.equal(await token.balanceOf(cashier.address), amountIn);
    });
  });
});
