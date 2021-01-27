pragma solidity <6.0 >=0.4.24;

import "./TransferValidatorBaseV2.sol";

interface ITokenSafe {
    function transferOwnership(address _newOwner) external;
    function withdrawToken(address _token, address _to, uint256 _amount) external returns (bool);
}

contract TransferValidatorWithTokenSafeV2 is TransferValidatorBaseV2 {
    ITokenSafe public safe;
    constructor(address _safe, address _tokenList, address _witnessList) public {
        safe = ITokenSafe(_safe);
        whitelistedTokens = Allowlist(_tokenList);
        whitelistedWitnesses = Allowlist(_witnessList);
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool) {
        return safe.withdrawToken(_token, _to, _amount);
    }

    function upgrade(address _newValidator) external onlyOwner {
        safe.transferOwnership(_newValidator);
    }
}
