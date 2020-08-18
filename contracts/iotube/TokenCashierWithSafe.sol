pragma solidity <6.0 >=0.4.24;

import "./TokenCashierBase.sol";

interface ERC20Token {
    function transferFrom(address from, address to, uint256 value) external returns (bool);
    function balanceOf(address who) external view returns (uint256);
    function transfer(address to, uint256 value) external returns (bool);
}

contract TokenCashierWithSafe is TokenCashierBase {
    address public safe;
    constructor(address _tokenList, address _safe) public {
        tokenList = TokenList(_tokenList);
        safe = _safe;
    }

    function transferToSafe(address _token, uint256 _amount) internal onlyOwner returns (bool) {
        ERC20Token token = ERC20Token(_token);
        require(token.transferFrom(msg.sender, safe, _amount), "fail to transfer token to safe");
        return true;
    }

    function withdrawToken(address _token) public onlyOwner {
        ERC20Token token = ERC20Token(_token);
        uint256 bal = token.balanceOf(address(this));
        if (bal > 0) {
            token.transfer(msg.sender, bal);
        }
    }
}