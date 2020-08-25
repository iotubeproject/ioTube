pragma solidity <6.0 >=0.4.24;

import "./UniqueAppendOnlyAddressList.sol";
import "../ownership/Ownable.sol";

contract VoterList is Ownable, UniqueAppendOnlyAddressList {
    event VoterAdded(address indexed voter);
    event VoterRemoved(address indexed voter);

    function isAllowed(address _voter) public view returns (bool) {
        return isActive(_voter);
    }

    function addVoter(address _voter) public onlyOwner returns (bool success_) {
        if (activateItem(_voter)) {
            emit VoterAdded(_voter);
            success_ = true;
        }
    }

    function removeVoter(address _voter) public onlyOwner returns (bool success_) {
        if (deactivateItem(_voter)) {
            emit VoterRemoved(_voter);
            success_ = true;
        }
    }

}