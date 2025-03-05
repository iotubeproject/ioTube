// SPDX-License-Identifier: MIT
pragma solidity = 0.8.20;

import "./interfaces/IUniswapV2Router02.sol";

interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
}

contract UniswapUnwrapper {
    event Swap(address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, address to);

    struct SwapData {
        address tokenOut;
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

    function transfer(address _token, address _to, uint256 _amount) external {
        // selector = bytes4(keccak256(bytes('transfer(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0xa9059cbb, _to, _amount));
        require(success && (data.length == 0 || abi.decode(data, (bool))), "UniswapUnwrapper: transfer failed");
        emit Swap(_token, address(0), _amount, 0, _to);
    }

    function onReceive(address _sender, address _token, uint256 _amount, bytes calldata _payload) external {
        address recipient = _sender;
        SwapData memory swapData = abi.decode(_payload, (SwapData));
        if (swapData.to != address(0)) {
            recipient = swapData.to;
        }
        if (swapData.deadline < block.timestamp) {
            this.transfer(_token, recipient, _amount);
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
        path[0] = _token;
        path[1] = weth;

        uint256[] memory amounts;
        try IUniswapV2Router02(router).getAmountsOut(_amount, path) returns (uint256[] memory _amounts) {
            amounts = _amounts;
        } catch {
            this.transfer(_token, recipient, _amount);
            return;
        }
        if (amounts[amounts.length - 1] < swapData.amountOutMin) {
            this.transfer(_token, recipient, _amount);
            return;
        }
        // The approve amount will be used right away
        // selector = bytes4(keccak256(bytes('approve(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x095ea7b3, router, _amount));
        require(success && (data.length == 0 || abi.decode(data, (bool))), "UniswapUnwrapper: approve failed");
        if (swapData.tokenOut == address(0)) {
            amounts = IUniswapV2Router02(router).swapExactTokensForETH(_amount, swapData.amountOutMin, path, recipient, swapData.deadline);
        } else if (swapData.tokenOut == weth) {
            amounts = IUniswapV2Router02(router).swapExactTokensForETH(_amount, swapData.amountOutMin, path, recipient, swapData.deadline);
        } else {
            amounts = IUniswapV2Router02(router).swapExactTokensForTokens(_amount, swapData.amountOutMin, path, recipient, swapData.deadline);
        }
        emit Swap(_token, swapData.tokenOut, _amount, amounts[amounts.length - 1], recipient);
    }
}
