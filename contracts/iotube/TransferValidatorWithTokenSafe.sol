pragma solidity <6.0 >=0.4.24;

import "./TokenSafe.sol";
import "./TransferValidatorBase.sol";

contract TransferValidatorWithTokenSafe is TransferValidatorBase {
    TokenSafe public safe;
    constructor(uint256 _expireHeight) public {
        setExpireHeight(_expireHeight);
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool) {
        return safe.withdrawToken(_token, _to, _amount);
    }

    function upgrade(address _newValidator) public onlyOwner {
        safe.transferOwnership(_newValidator);
    }
}