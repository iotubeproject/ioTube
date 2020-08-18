pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

interface ERC20Token {
    function transfer(address to, uint256 value) external returns (bool);
}

contract TokenSafe is Ownable {
    function withdrawToken(address _token, address _to, uint256 _amount) public onlyOwner returns (bool) {
        require(ERC20Token(_token).transfer(_to, _amount), "failed to transfer token");
        return true;
    }
}