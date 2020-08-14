pragma solidity ^0.4.24;

import "../ownership/Ownable.sol";

contract Cashier is Ownable {
    address public safe;
    uint256 public depositFee;
    uint256 public minAmount;
    uint256 public maxAmount;

    address[] customers;
    address[] receivers;
    uint256[] amounts;
    uint256[] fees;

    event Receipt(address indexed customer, address indexed receiver, uint256 amount, uint256 fee);

    constructor(address _safe, uint256 _fee, uint256 _minAmount, uint256 _maxAmount) public {
        require(_minAmount > 0);
        require(_maxAmount >= _minAmount);
        safe = _safe;
        depositFee = _fee;
        minAmount = _minAmount;
        maxAmount = _maxAmount;
    }

    function setDepositFee(uint256 _fee) external onlyOwner {
        depositFee = _fee;
    }

    function setSafe(address _safe) external onlyOwner {
        safe = _safe;
    }

    function setMinAmount(uint256 _minAmount) external onlyOwner {
        require(maxAmount >= _minAmount);
        require(_minAmount > 0);
        minAmount = _minAmount;
    }

    function setMaxAmount(uint256 _maxAmount) external onlyOwner {
        require(_maxAmount >= minAmount);
        maxAmount = _maxAmount;
    }

    function count() public view returns (uint256) {
        return customers.length;
    }

    function getRecords(uint256 _offset, uint256 _limit) public view returns(address[] _customers, address[] _receivers, uint256[] _amounts, uint256[] _fees) {
        require(_limit < 200);
        if (_offset >= customers.length) {
            return (_customers, _receivers, _amounts, _fees);
        }
        uint256 l = _limit;
        if (_limit > customers.length - _offset) {
            l = customers.length - _offset;
        }
        if (l == 0) {
            return (_customers, _receivers, _amounts, _fees);
        }
        _customers = new address[](l);
        _receivers = new address[](l);
        _amounts = new uint256[](l);
        _fees = new uint256[](l);
        for (uint256 i = 0; i < l; i++) {
            _customers[i] = customers[_offset + i];
            _receivers[i] = receivers[_offset + i];
            _amounts[i] = amounts[_offset + i];
            _fees[i] = fees[_offset + i];
        }

        return (_customers, _receivers, _amounts, _fees);
    }

    function withdraw() external onlyOwner {
        msg.sender.transfer(address(this).balance);
    }
}