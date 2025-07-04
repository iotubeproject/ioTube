// SPDX-License-Identifier: MIT
pragma solidity = 0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface ITetherSwapper {
    function deposit(uint256 _amount) external;
    function usdt() external view returns (address);
    function iousdt() external view returns (address);
}

contract USDTUnwrapper {
    IERC20 public iousdt;
    IERC20 public usdt;
    ITetherSwapper public swapper;

    constructor(ITetherSwapper _swapper) {
        swapper = _swapper;
        iousdt = IERC20(_swapper.iousdt());
        usdt = IERC20(_swapper.usdt());
    }

    function onReceive(address _sender, address _token, uint256 _amount, bytes calldata _payload) external {
        require(_token == address(iousdt), "USDTUnwrapper: invalid token");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        require(iousdt.approve(address(swapper), _amount), "USDTUnwrapper: approve failed");
        swapper.deposit(_amount);
        require(usdt.transfer(recipient, _amount), "USDTUnwrapper: transfer failed");
    }
}
