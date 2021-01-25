pragma solidity <6.0 >=0.4.24;

import "./TransferValidatorBaseV2.sol";

interface IMinterPool {
    function mint(address, address, uint256) external returns(bool);
    function transferOwnership(address) external;
}

contract TransferValidatorWithMinterPoolV2 is TransferValidatorBaseV2 {
    IMinterPool public pool;

    constructor(address _minterPool, address _tokenList, address _witnessList) public {
        pool = IMinterPool(_minterPool);
        whitelistedTokens = Allowlist(_tokenList);
        whitelistedWitnesses = Allowlist(_witnessList);
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool) {
        return pool.mint(_token, _to, _amount);
    }

    function upgrade(address _newValidator) public onlyOwner {
        pool.transferOwnership(_newValidator);
    }
}