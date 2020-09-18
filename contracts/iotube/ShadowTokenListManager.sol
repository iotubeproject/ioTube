pragma solidity <6.0 >=0.4.24;

import "./TokenList.sol";
import "../token/ShadowToken.sol";

contract ShadowTokenListManager is Ownable {
    TokenList public tokenList;

    constructor(address _addr) public {
        tokenList = TokenList(_addr);
    }

    function addShadowToken(address _minter, address _coToken, string memory _name, string memory _symbol, uint256 _min, uint256 _max) public onlyOwner returns (bool success_) {
        ShadowToken st = new ShadowToken(_minter, _coToken, _name, _symbol);
        return tokenList.addToken(address(st), _min, _max);
    }

    function removeToken(address _token) public onlyOwner returns (bool success_) {
        return tokenList.removeToken(_token);
    }

    function setMinAmount(address _token, uint256 _minAmount) public onlyOwner {
        tokenList.setMinAmount(_token, _minAmount);
    }

    function setMaxAmount(address _token, uint256 _maxAmount) public onlyOwner {
        tokenList.setMaxAmount(_token, _maxAmount);
    }

    function upgrade(address _newManager) public onlyOwner {
        tokenList.transferOwnership(_newManager);
    }
}