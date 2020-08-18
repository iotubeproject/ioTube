pragma solidity <6.0 >=0.4.24;

import "./TokenList.sol";
import "../lifecycle/Pausable.sol";

contract TokenCashierBase is Pausable {
    event Receipt(address indexed customer, address indexed receiver, address indexed token, uint256 amount, uint256 fee);

    struct Cashier {
        address[] customers;
        address[] receivers;
        uint256[] amounts;
        uint256[] fees;
    }

    TokenList public tokenList;
    mapping(address => Cashier) cashiers;
    uint256 public depositFee;

    function() external {
        revert();
    }

    function setDepositFee(uint256 _fee) public onlyOwner {
        depositFee = _fee;
    }

    function transferToSafe(address _token, uint256 _amount) internal returns (bool);

    function depositTo(address _token, address _to, uint256 _amount) public whenNotPaused payable {
        require(tokenList.isAllowed(_token), "token is not in list");
        require(_to != address(0), "invalid destination");
        require(msg.value >= depositFee, "insufficient balance");
        require(_amount >= tokenList.minAmount(_token), "amount too low");
        require(_amount <= tokenList.maxAmount(_token), "amount too high");
        require(transferToSafe(_token, _amount), "failed to put into safe");
        Cashier storage cashier = cashiers[_token];
        cashier.customers.push(msg.sender);
        cashier.receivers.push(_to);
        cashier.amounts.push(_amount);
        cashier.fees.push(msg.value);
        emit Receipt(msg.sender, _to, _token, _amount, msg.value);
    }

    function deposit(address _token, uint256 _amount) public payable {
        depositTo(_token, msg.sender, _amount);
    }

    function count(address _token) public view returns (uint256) {
        return cashiers[_token].customers.length;
    }

    function getRecords(address _token, uint256 _offset, uint256 _limit) public view returns(address[] memory _customers, address[] memory _receivers, uint256[] memory _amounts, uint256[] memory _fees) {
        require(_limit < 200);
        require(tokenList.isAllowed(_token), "token not in list");
        Cashier storage cashier = cashiers[_token];
        if (_offset >= cashier.customers.length) {
            return (_customers, _receivers, _amounts, _fees);
        }
        uint256 l = _limit;
        if (_limit > cashier.customers.length - _offset) {
            l = cashier.customers.length - _offset;
        }
        if (l == 0) {
            return (_customers, _receivers, _amounts, _fees);
        }
        _customers = new address[](l);
        _receivers = new address[](l);
        _amounts = new uint256[](l);
        _fees = new uint256[](l);
        for (uint256 i = 0; i < l; i++) {
            _customers[i] = cashier.customers[_offset + i];
            _receivers[i] = cashier.receivers[_offset + i];
            _amounts[i] = cashier.amounts[_offset + i];
            _fees[i] = cashier.fees[_offset + i];
        }

        return (_customers, _receivers, _amounts, _fees);
    }
}