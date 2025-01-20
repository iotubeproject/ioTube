pragma solidity = 0.8.20;

import "../interfaces/IUniswapV2Router02.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MockUniswapV2Router02 is IUniswapV2Router02 {
    address public immutable override WETH;
    uint256 public amountOut;
    
    constructor(address _weth) {
        WETH = _weth;
    }

    receive() external payable {
    }

    function setAmountOut(uint256 _amountOut) external {
        amountOut = _amountOut;
    }

    function swapExactTokensForTokens(
        uint256 _amountIn,
        uint256 _amountOutMin,
        address[] calldata _path,
        address _to,
        uint256 _deadline
    ) external override returns (uint256[] memory amounts) {
        require(_deadline >= block.timestamp, "UniswapV2Router: EXPIRED");
        require(_amountOutMin <= amountOut, "UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT");
        require(IERC20(_path[0]).transferFrom(msg.sender, address(this), _amountIn), "UniswapV2Router: transferFrom failed");
        amounts = new uint256[](_path.length);
        amounts[0] = _amountIn;
        for(uint i = 1; i < _path.length; i++) {
            amounts[i] = _amountIn; // Mock: return same amount
        }
        amounts[_path.length - 1] = amountOut;
        require(IERC20(_path[_path.length - 1]).transfer(_to, amountOut), "UniswapV2Router: transfer failed");
        
        return amounts;
    }

    function swapExactTokensForETH(
        uint256 _amountIn,
        uint256 _amountOutMin,
        address[] calldata _path,
        address _to,
        uint256 _deadline
    ) external override returns (uint256[] memory amounts) {
        require(_deadline >= block.timestamp, "UniswapV2Router: EXPIRED");
        require(_amountOutMin <= amountOut, "UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT");
        require(IERC20(_path[0]).transferFrom(msg.sender, address(this), _amountIn), "UniswapV2Router: transferFrom failed");
        amounts = new uint256[](_path.length);
        amounts[0] = _amountIn;
        for(uint i = 1; i < _path.length; i++) {
            amounts[i] = _amountIn; // Mock: return same amount
        }
        amounts[_path.length - 1] = amountOut;
        payable(_to).transfer(amountOut);
        return amounts;
    }

    function swapTokensForExactETH(
        uint256 _amountOut,
        uint256 _amountInMax,
        address[] calldata _path,
        address _to,
        uint256 _deadline
    ) external override returns (uint256[] memory amounts) {
        require(_deadline >= block.timestamp, "UniswapV2Router: EXPIRED");
        require(_amountInMax >= amountOut, "UniswapV2Router: EXCESSIVE_INPUT_AMOUNT");
        require(IERC20(_path[0]).transferFrom(msg.sender, address(this), amountOut), "UniswapV2Router: transferFrom failed");
        amounts = new uint256[](_path.length);
        amounts[0] = amountOut;
        for(uint i = 1; i < _path.length; i++) {
            amounts[i] = amountOut; // Mock: return same amount
        }
        amounts[_path.length - 1] = _amountOut;
        payable(_to).transfer(_amountOut);
        return amounts;
    }

    function getAmountsOut(uint256 _amountIn, address[] calldata _path)
        external
        view
        override
        returns (uint256[] memory amounts)
    {
        amounts = new uint256[](_path.length);
        amounts[0] = _amountIn;
        for(uint i = 1; i < _path.length; i++) {
            amounts[i] = _amountIn; // Mock: return same amount
        }
        amounts[_path.length - 1] = amountOut;
        return amounts;
    }

}
