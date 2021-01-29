pragma solidity <6.0 >=0.4.24;

import "./TokenCashierBaseV2.sol";

contract BurnableTokenCashierV2 is TokenCashierBaseV2 {
    constructor(address _tokenList) public {
        tokenList = ITokenList(_tokenList);
    }

    function transferToSafe(address _token, uint256 _amount) internal returns (bool) {
        require(safeTransferFrom(_token, msg.sender, address(this), _amount), "fail to transfer token to cashier");
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x42966c68, _amount));
        require(success && (data.length == 0 || abi.decode(data, (bool))), "fail to burn token");
        return true;
    }
}