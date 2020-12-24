pragma solidity <6.0 >=0.4.24;

import "./TokenList.sol";

interface TokenCashier {
    function tokenList() external view returns(TokenList);

    function depositTo(address _token, address _to, uint256 _amount) external payable;

    function count(address _token) external view returns (uint256);

    function getRecords(address _token, uint256 _offset, uint256 _limit) external view returns(address[] memory customers_, address[] memory receivers_, uint256[] memory amounts_, uint256[] memory fees_);
}

contract ComboTokenCashier {
    TokenCashier public standardTokenCashier;
    TokenCashier public nonStandardTokenCashier;
    
    constructor(address _standardTokenCashier, address _nonStandardTokenCashier) public {
        standardTokenCashier = TokenCashier(_standardTokenCashier);
        nonStandardTokenCashier = TokenCashier(_nonStandardTokenCashier);
    }

    function tokenList() external pure returns(TokenList) {
        require(false, "not implemented");
    }

    function depositTo(address _token, address _to, uint256 _amount) external payable {
        if (standardTokenCashier.tokenList().isAllowed(_token)) {
            return standardTokenCashier.depositTo(_token, _to, _amount);
        }
        return nonStandardTokenCashier.depositTo(_token, _to, _amount);
    }

    function count(address _token) public view returns (uint256) {
        if (standardTokenCashier.tokenList().isAllowed(_token)) {
            return standardTokenCashier.count(_token);
        }
        return nonStandardTokenCashier.count(_token);
    }

    function getRecords(address _token, uint256 _offset, uint256 _limit) public view returns(address[] memory customers_, address[] memory receivers_, uint256[] memory amounts_, uint256[] memory fees_) {
        if (standardTokenCashier.tokenList().isAllowed(_token)) {
            return standardTokenCashier.getRecords(_token, _offset, _limit);
        }
        return nonStandardTokenCashier.getRecords(_token, _offset, _limit);
    }

}