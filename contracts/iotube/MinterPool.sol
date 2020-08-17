pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../token/MintableToken.sol";

contract MinterPool is Ownable {
    constructor(address _owner) public {
        owner = _owner;
    }

    function mint(address _token, address _to, uint256 _amount) public onlyOwner returns (bool) {
        require(MintableToken(_token).mint(_to, _amount), "fail to mint");
        return true;
    }
}