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

    function getRecords(address _token, uint256 _offset, uint256 _limit) public view returns(address[] memory customers_, address[] memory receivers_, uint256[] memory amounts_, uint256[] memory fees_) {
        require(_limit < 200);
        require(tokenList.isAllowed(_token), "token not in list");
        Cashier storage cashier = cashiers[_token];
        if (_offset >= cashier.customers.length) {
            return (customers_, receivers_, amounts_, fees_);
        }
        uint256 l = _limit;
        if (_limit > cashier.customers.length - _offset) {
            l = cashier.customers.length - _offset;
        }
        if (l == 0) {
            return (customers_, receivers_, amounts_, fees_);
        }
        customers_ = new address[](l);
        receivers_ = new address[](l);
        amounts_ = new uint256[](l);
        fees_ = new uint256[](l);
        for (uint256 i = 0; i < l; i++) {
            customers_[i] = cashier.customers[_offset + i];
            receivers_[i] = cashier.receivers[_offset + i];
            amounts_[i] = cashier.amounts[_offset + i];
            fees_[i] = cashier.fees[_offset + i];
        }

        return (customers_, receivers_, amounts_, fees_);
    }

    function withdraw() external onlyOwner {
        msg.sender.transfer(address(this).balance);
    }
}