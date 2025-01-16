// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

import "./interfaces/IUniswapV2Router02.sol";

interface ICashier {
    function depositFee() external view returns (uint256);
    function depositTo(address, address, uint256, bytes calldata) external payable;
}

interface ICrosschainToken {
    function approve(address, uint256) external returns (bool);
    function transfer(address, uint256) external returns (bool);
}

contract EthereumHub {
    address immutable public cashier;
    address immutable public router;

    constructor(address _cashier, address _router) {
        cashier = _cashier;
        router = _router;
    }

    receive() external payable {
    }

    function onReceive(address _sender, ICrosschainToken _token, uint256 _amount, bytes calldata _payload) external {
        ICashier c = ICashier(cashier);
        uint256 fee = c.depositFee();
        if (fee > 0) {
            IUniswapV2Router02 r = IUniswapV2Router02(router);
            address[] memory path = new address[](2);
            path[0] = address(_token);
            path[1] = r.WETH();
            uint256[] memory amounts;
            try r.getAmountsOut(_amount, path) returns (uint256[] memory _amounts) {
                amounts = _amounts;
            } catch {
                require(_token.transfer(_sender, _amount), "EthereumHub: estimate amounts out failed");
                return;
            }
            if (amounts[1] <= fee) {
                require(_token.transfer(_sender, _amount), "EthereumHub: transfer failed");
                return;
            }
            require(_token.approve(address(r), _amount), "EthereumHub: approve router failed");
            amounts = r.swapTokensForExactETH(fee, _amount, path, address(this), block.timestamp + 10);
            _amount = _amount - amounts[0];
        }
        (address recipient, bytes memory payload) = abi.decode(_payload, (address, bytes));
        require(_token.approve(address(c), _amount), "EthereumHub: approve cashier failed");
        c.depositTo{value: fee}(address(_token), recipient, _amount, payload);
    }

}
