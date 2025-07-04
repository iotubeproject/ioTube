pragma solidity = 0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface ITreasury {
    function repay(uint256 _amount) external;
}

contract TetherSwapper {
    IERC20 public usdt;
    IERC20 public iousdt;
    ITreasury public treasury;

    constructor(IERC20 _usdt, IERC20 _iousdt, ITreasury _treasury) {
        usdt = _usdt;
        iousdt = _iousdt;
        treasury = _treasury;
    }

    function deposit(uint256 _amount) external {
        require(iousdt.transferFrom(msg.sender, address(this), _amount), "failed to deposit");
        require(usdt.transferFrom(address(treasury), msg.sender, _amount), "failed to transfer");
    }

    function withdraw(uint256 _amount) external {
        require(balance() >= _amount, "insufficient balance");
        require(usdt.transferFrom(msg.sender, address(this), _amount), "failed to withdraw");
        require(iousdt.transfer(msg.sender, _amount), "failed to transfer");
        require(usdt.approve(address(treasury), _amount), "failed to approve");
        treasury.repay(_amount);
    }

    function balance() public view returns (uint256) {
        return iousdt.balanceOf(address(this));
    }
}
