pragma solidity <6.0 >=0.4.24;

import "../token/ERC20.sol";
import "./Safe.sol";

contract TokenSafe is Safe {
    ERC20 public token;
    address public validator;

    constructor(address tokenAddr, address validatorAddr) public {
        token = ERC20(tokenAddr);
        validator = validatorAddr;
    }

    function withdraw(address _to, uint256 _amount) public returns (bool) {
        require(msg.sender == validator, "only validator could withdraw from safe");
        require(token.transfer(_to, _amount), "transfer failed");
        emit Withdrew(_to, _amount);
        return true;
    }
}