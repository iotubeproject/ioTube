pragma solidity <6.0 >=0.4.24;

import "./UniqueAppendOnlyAddressList.sol";
import "../ownership/Ownable.sol";

contract WitnessList is Ownable, UniqueAppendOnlyAddressList {
    event WitnessAdded(address indexed witness);
    event WitnessRemoved(address indexed witness);

    function isAllowed(address _witness) public view returns (bool) {
        return isActive(_witness);
    }

    function addWitness(address _witness) public onlyOwner returns (bool success_) {
        if (activateItem(_witness)) {
            emit WitnessAdded(_witness);
            success_ = true;
        }
    }

    function removeWitness(address _witness) public onlyOwner returns (bool success_) {
        if (deactivateItem(_witness)) {
            emit WitnessRemoved(_witness);
            success_ = true;
        }
    }

}