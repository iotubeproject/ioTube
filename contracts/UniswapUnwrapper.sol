// SPDX-License-Identifier: MIT
pragma solidity = 0.8.20;

import "./interfaces/IUniswapV2Router02.sol";

interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
}

contract UniswapUnwrapper {
    event Swap (address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, address to);

    struct SwapData {
        address tokenIn;
        address tokenOut;
        uint256 amountIn;
        uint256 amountOutMin;
        address to;
        uint256 deadline;
    }
    address immutable public router;

    constructor(address _router) {
        router = _router;
    }

    receive() external payable {
    }

    function onReceive(address _sender, IERC20 _token, uint256 _amount, bytes calldata _payload) external {
        address recipient = _sender;
        SwapData memory swapData = abi.decode(_payload, (SwapData));
        require(swapData.tokenIn == address(_token), "UniswapUnwrapper: invalid input token");
        require(swapData.amountIn == _amount, "UniswapUnwrapper: invalid amount");
        if (swapData.to != address(0)) {
            recipient = swapData.to;
        }
        if (swapData.deadline < block.timestamp) {
            require(_token.transfer(recipient, _amount), "UniswapUnwrapper: transfer failed");
            emit Swap(swapData.tokenIn, address(0), swapData.amountIn, 0, recipient);
            return;
        }
        address weth = IUniswapV2Router02(router).WETH();
        address[] memory path;
        if (swapData.tokenOut == address(0)) {
            path = new address[](2);
        } else if (swapData.tokenOut == weth) {
            path = new address[](2);
        } else {
            path = new address[](3);
            path[2] = swapData.tokenOut;
        }
        path[0] = swapData.tokenIn;
        path[1] = weth;

        uint256[] memory amounts = IUniswapV2Router02(router).getAmountsOut(swapData.amountIn, path);
        if (amounts[amounts.length - 1] < swapData.amountOutMin) {
            require(_token.transfer(recipient, _amount), "UniswapUnwrapper: transfer failed");
            emit Swap(swapData.tokenIn, address(0), swapData.amountIn, 0, recipient);
            return;
        }
        _token.approve(router, swapData.amountIn);
        if (swapData.tokenOut == address(0)) {
            amounts = IUniswapV2Router02(router).swapExactTokensForETH(swapData.amountIn, swapData.amountOutMin, path, recipient, swapData.deadline);
        } else if (swapData.tokenOut == weth) {
            amounts = IUniswapV2Router02(router).swapExactTokensForETH(swapData.amountIn, swapData.amountOutMin, path, recipient, swapData.deadline);
        } else {
            amounts = IUniswapV2Router02(router).swapExactTokensForTokens(swapData.amountIn, swapData.amountOutMin, path, recipient, swapData.deadline);
        }
        emit Swap(swapData.tokenIn, swapData.tokenOut, swapData.amountIn, amounts[amounts.length - 1], recipient);
    }
}
