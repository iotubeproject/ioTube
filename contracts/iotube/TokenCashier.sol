pragma solidity <6.0 >=0.4.24;

import "./Cashier.sol";
import "../token/ERC20.sol";
import "../lifecycle/Pausable.sol";

interface BurnableToken {
    function burn(uint) external returns (bool);
}

contract TokenCashier is Pausable, Cashier {
    ERC20 public token;
    address public safe;
    BurnableToken private burner;

    constructor(
        address _tokenAddr,
        address _safe,
        uint256 _fee,
        uint256 _minAmount,
        uint256 _maxAmount
    ) public {
        token = ERC20(_tokenAddr);
        if (_safe == address(0)) {
            burner = BurnableToken(_tokenAddr);
        }
        safe = _safe;
        setDepositFee(_fee);
        setMaxAmount(_maxAmount);
        setMinAmount(_minAmount);
    }

    function() external {
        revert();
    }

    function depositTo(address _to, uint256 _amount) public whenNotPaused payable {
        require(_to != address(0), "invalid destination");
        require(msg.value >= depositFee, "insufficient balance");
        require(_amount >= minAmount, "insufficient amount");
        require(_amount <= maxAmount, "amount too high");
        if (safe != address(0)) {
            require(token.transferFrom(msg.sender, safe, _amount), "put into safe required");
        } else {
            require(token.transferFrom(msg.sender, address(this), _amount), "transfer token to cashier");
            require(burner.burn(_amount), "burn tokens");
        }
        customers.push(msg.sender);
        receivers.push(_to);
        amounts.push(_amount);
        fees.push(msg.value);
        emit Receipt(msg.sender, _to, _amount, msg.value);
    }

    function withdrawToken() public onlyOwner {
        uint256 bal = token.balanceOf(address(this));
        if (bal > 0) {
            token.transfer(msg.sender, bal);
        }
    }
}