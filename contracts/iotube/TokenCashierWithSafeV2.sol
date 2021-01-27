pragma solidity <6.0 >=0.4.24;

import "./TokenCashierBaseV2.sol";

contract TokenCashierWithSafeV2 is TokenCashierBaseV2 {
    address public safe;
    constructor(address _tokenList, address _safe) public {
        tokenList = ITokenList(_tokenList);
        safe = _safe;
    }

    function transferToSafe(address _token, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, msg.sender, safe, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

    function withdrawToken(address _token) public onlyOwner {
        // selector = bytes4(keccak256(bytes('balanceOf(address)')))
        (bool success, bytes memory balance) = _token.call(abi.encodeWithSelector(0x70a08231, address(this)));
        require(success, "failed to call balanceOf");
        uint256 bal = abi.decode(balance, (uint256));
        if (bal > 0) {
            bytes memory data;
            // selector = bytes4(keccak256(bytes('transfer(address,uint256)')))
            (success, data) = _token.call(abi.encodeWithSelector(0xa9059cbb, msg.sender, bal));
            require(success && (data.length == 0 || abi.decode(data, (bool))), "failed to withdraw token");
        }
    }
}