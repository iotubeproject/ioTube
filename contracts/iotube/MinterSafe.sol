pragma solidity <6.0 >=0.4.24;

import "../token/MintableToken.sol";
import "./Safe.sol";

contract MinterSafe is Safe {
    MintableToken public token;
    address validator;

    constructor(address tokenAddr, address validatorAddr) public {
        token = MintableToken(tokenAddr);
        validator = validatorAddr;
    }

    function withdraw(address _to, uint256 _amount) public returns (bool) {
        require(msg.sender == validator, "only validator could withdraw from safe");
        require(token.mint(_to, _amount), "transfer failed");
        emit Withdrew(_to, _amount);
        return true;
    }
}