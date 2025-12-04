// SPDX-License-Identifier: MIT
pragma solidity <6.0 >=0.4.24;

import "./UniqueAppendOnlyAddressList.sol";

contract WitnessTokenList is UniqueAppendOnlyAddressList {
    event TokenAdded(address indexed token);
    event TokenRemoved(address indexed token);

    function isAllowed(address _token) public view returns (bool) {
        return isActive(_token);
    }

    function addToken(address _token) public onlyOwner returns (bool success_) {
        if (activateItem(_token)) {
            emit TokenAdded(_token);
            success_ = true;
        }
    }

    function removeToken(address _token) public onlyOwner returns (bool success_) {
        if (deactivateItem(_token)) {
            emit TokenRemoved(_token);
            success_ = true;
        }
    }
}