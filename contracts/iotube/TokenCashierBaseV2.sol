pragma solidity <6.0 >=0.4.24;

import "./TokenList.sol";
import "../lifecycle/Pausable.sol";

contract TokenCashierBaseV2 is Pausable {
    event Receipt(uint256 indexed id, address indexed token, address indexed receiver, address sender, uint256 amount, uint256 fee);

    TokenList public tokenList;
    mapping(address => uint256) public counts;
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
        counts[_token] += 1;
        emit Receipt(counts[_token], _token, _to, msg.sender, _amount, msg.value);
    }

    function deposit(address _token, uint256 _amount) public payable {
        depositTo(_token, msg.sender, _amount);
    }

    function count(address _token) public view returns (uint256) {
        return counts[_token];
    }

    function withdraw() external onlyOwner {
        msg.sender.transfer(address(this).balance);
    }
}