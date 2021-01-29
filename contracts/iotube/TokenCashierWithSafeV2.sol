pragma solidity <6.0 >=0.4.24;

import "./TokenCashierBaseV2.sol";

contract TokenCashierWithSafeV2 is TokenCashierBaseV2 {
    address public safe;
    constructor(address _tokenList, address _safe) public {
        tokenList = ITokenList(_tokenList);
        safe = _safe;
    }

    function transferToSafe(address _token, uint256 _amount) internal returns (bool) {
        return safeTransferFrom(_token, msg.sender, safe, _amount);
    }
}