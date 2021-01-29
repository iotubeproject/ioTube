pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

interface MintableToken {
    function mint(address, uint256) external returns(bool);
}

contract MinterPool is Ownable {
    function mint(address _token, address _to, uint256 _amount) public onlyOwner returns (bool) {
        return MintableToken(_token).mint(_to, _amount);
    }
}