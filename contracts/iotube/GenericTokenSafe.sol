pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

contract GenericTokenSafe is Ownable {
    function withdrawToken(address _token, address _to, uint256 _amount) public onlyOwner returns (bool) {
        // selector = bytes4(keccak256(bytes('transfer(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0xa9059cbb, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));	    
    }
}
