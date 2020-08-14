pragma solidity ^0.4.24;

import "./Cashier.sol";
import "../lifecycle/Pausable.sol";

contract CoinCashier is Pausable, Cashier {
    uint256 public gasLimit;

    constructor(
        address _safe,
        uint256 _fee,
        uint256 _minAmount,
        uint256 _maxAmount
    ) Cashier(_safe, _fee, _minAmount, _maxAmount) public {
        require(_safe != address(0), "coin cashier address cannot be 0x0");
        gasLimit = 4000000;
    }

    function () external payable {
        deposit();
    }

    function deposit() public payable {
        depositTo(msg.sender);
    }

    function depositTo(address _to) public whenNotPaused payable {
        require(_to != address(0));
        require(msg.value >= minAmount + depositFee);
        uint256 amount = msg.value - depositFee;
        require(amount <= maxAmount);

        if (safe.call.value(amount).gas(gasLimit)()) {
            customers.push(msg.sender);
            receivers.push(_to);
            amounts.push(amount);
            fees.push(depositFee);

            emit Receipt(msg.sender, _to, amount, depositFee);
        }
    }

    function setGasLimit(uint256 _gasLimit) external onlyOwner {
        gasLimit = _gasLimit;
    }

}
