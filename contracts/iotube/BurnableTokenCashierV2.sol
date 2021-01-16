pragma solidity <6.0 >=0.4.24;

import "./TokenCashierBaseV2.sol";

interface BurnableToken {
    function burn(uint) external returns (bool);
    function transferFrom(address from, address to, uint256 value) external returns (bool);
    function balanceOf(address who) external view returns (uint256);
    function transfer(address to, uint256 value) external returns (bool);
}

contract BurnableTokenCashierV2 is TokenCashierBaseV2 {
    constructor(address _tokenList) public {
        tokenList = TokenList(_tokenList);
    }

    function transferToSafe(address _token, uint256 _amount) internal returns (bool) {
        BurnableToken token = BurnableToken(_token);
        require(token.transferFrom(msg.sender, address(this), _amount), "fail to transfer token to cashier");
        require(token.burn(_amount), "fail to burn token");
        return true;
    }

    function withdrawToken(address _token) public onlyOwner {
        BurnableToken token = BurnableToken(_token);
        uint256 bal = token.balanceOf(address(this));
        if (bal > 0) {
            token.transfer(msg.sender, bal);
        }
    }
}