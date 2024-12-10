const MockERC20 = artifacts.require("MockERC20");
const MockUniswapV2Router02 = artifacts.require("MockUniswapV2Router02");
const UniswapUnwrapper = artifacts.require("UniswapUnwrapper");
const {ethers, AbiCoder} = require("ethers");

const expectSwap = function(log, tokenIn, tokenOut, amountIn, amountOut, to) {
    assert.equal(log.event, "Swap");
    assert.equal(log.args.tokenIn, tokenIn);
    assert.equal(log.args.tokenOut, tokenOut);
    assert.equal(log.args.amountIn, amountIn);
    assert.equal(log.args.amountOut, amountOut);
    assert.equal(log.args.to, to);
}

contract("UniswapUnwrapper", function([owner, recipient]) {
  let unwrapper;
  let router;
  let token;
  let weth;

  beforeEach(async function() {
    // Deploy mock WETH
    weth = await MockERC20.new("Wrapped ETH", "WETH");

    // Deploy mock router
    router = await MockUniswapV2Router02.new(weth.address);

    // Deploy mock token
    token = await MockERC20.new("Test Token", "TEST");

    // Deploy unwrapper
    unwrapper = await UniswapUnwrapper.new(router.address);

    // Approve tokens
    await token.approve(unwrapper.address, ethers.MaxUint256);
    await token.approve(router.address, ethers.MaxUint256);
  });

  it("should swap tokens for ETH", async function() {
    const amountIn = ethers.parseEther("100");
    const amountOutMin = ethers.parseEther("90");
    const deadline = Math.floor(Date.now() / 1000) + 3600;
    await router.sendTransaction({value: Number(amountOutMin)});
  
    const swapData = {
      tokenIn: token.address,
      tokenOut: ethers.ZeroAddress, // ETH
      amountIn: amountIn,
      amountOutMin: amountOutMin,
      to: recipient,
      deadline: deadline
    };
  
    const encodedData = AbiCoder.defaultAbiCoder().encode(
        ["tuple(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address to, uint256 deadline)"],
        [swapData]
    );

    await router.setAmountOut(amountOutMin);
    await token.transfer(unwrapper.address, amountIn);
    const tx = await unwrapper.onReceive(owner, token.address, amountIn, encodedData);
    expectSwap(tx.logs[0], token.address, ethers.ZeroAddress, amountIn, amountOutMin, recipient);
  });

  it("should swap tokens for tokens", async function() {
    const amountIn = ethers.parseEther("100");
    const amountOutMin = ethers.parseEther("90");
    const deadline = Math.floor(Date.now() / 1000) + 3600;

    const token2 = await MockERC20.new("Test Token 2", "TEST2");
    await token2.transfer(router.address, amountOutMin);

    const swapData = {
      tokenIn: token.address,
      tokenOut: token2.address,
      amountIn: amountIn,
      amountOutMin: amountOutMin,
      to: recipient,
      deadline: deadline
    };

    const encodedData = AbiCoder.defaultAbiCoder().encode(
      ["tuple(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address to, uint256 deadline)"],
      [swapData]
    );

    await router.setAmountOut(amountOutMin);
    await token.transfer(unwrapper.address, amountIn);
    const tx = await unwrapper.onReceive(owner, token.address, amountIn, encodedData);
    expectSwap(tx.logs[0], token.address, token2.address, amountIn, amountOutMin, recipient);
  });

  it("should transfer token in if deadline passed", async function() {
    const amountIn = ethers.parseEther("100");
    const amountOutMin = ethers.parseEther("90");
    const deadline = Math.floor(Date.now() / 1000) - 3600; // Passed deadline

    const swapData = {
      tokenIn: token.address,
      tokenOut: ethers.ZeroAddress,
      amountIn: amountIn,
      amountOutMin: amountOutMin,
      to: recipient,
      deadline: deadline
    };

    const encodedData = AbiCoder.defaultAbiCoder().encode(
      ["tuple(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address to, uint256 deadline)"],
      [swapData]
    );

    await token.transfer(unwrapper.address, amountIn);
    const tx = await unwrapper.onReceive(owner, token.address, amountIn, encodedData);
    expectSwap(tx.logs[0], token.address, ethers.ZeroAddress, amountIn, 0, recipient);
  });

  it("should transfer token in if slippage too high", async function() {
    const amountIn = ethers.parseEther("100");
    const amountOutMin = ethers.parseEther("90");
    const deadline = Math.floor(Date.now() / 1000) + 3600;

    const swapData = {
      tokenIn: token.address,
      tokenOut: ethers.ZeroAddress,
      amountIn: amountIn,
      amountOutMin: amountOutMin,
      to: recipient,
      deadline: deadline
    };

    const encodedData = AbiCoder.defaultAbiCoder().encode(
      ["tuple(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address to, uint256 deadline)"],
      [swapData]
    );

    await router.setAmountOut(amountOutMin - ethers.toBigInt(1)); // Set amount less than minimum
    await token.transfer(unwrapper.address, amountIn);
    const tx = await unwrapper.onReceive(owner, token.address, amountIn, encodedData);
    expectSwap(tx.logs[0], token.address, ethers.ZeroAddress, amountIn, 0, recipient);
  });
});
