// SPDX-License-Identifier: MIT
pragma solidity >=0.8.20;

import "./UniqueAppendOnlyAddressList.sol";

contract WitnessListV3 is UniqueAppendOnlyAddressList {
    event WitnessUpdated();

    function isAllowed(address _witness) public view returns (bool) {
        return isActive(_witness);
    }

    function areAllowed(address[] calldata _witnesses) public view returns (bool) {
        for (uint256 i = 0; i < _witnesses.length; i++) {
            if (!isActive(_witnesses[i])) {
                return false;
            }
        }
        return true;
    }

    function addWitness(address _witness) public onlyOwner returns (bool success_) {
        if (activateItem(_witness)) {
            success_ = true;
            emit WitnessUpdated();
        }
    }

    function addWitnesses(address[] calldata _witnesses) public onlyOwner {
        for (uint256 i = 0; i < _witnesses.length; i++) {
            require(activateItem(_witnesses[i]), "witness already active");
        }
        emit WitnessUpdated();
    }


    function removeWitness(address _witness) public onlyOwner returns (bool success_) {
        if (deactivateItem(_witness)) {
            success_ = true;
            emit WitnessUpdated();
        }
    }

    function removeWitnesses(address[] calldata _witnesses) public onlyOwner {
        for (uint256 i = 0; i < _witnesses.length; i++) {
            require(deactivateItem(_witnesses[i]), "witness not active");
        }
        emit WitnessUpdated();
    }

    function switchWitness(address _newWitness) public {
        address witness = msg.sender;
        require(deactivateItem(witness), "WitnessList: deactivate witness failed");
        require(activateItem(_newWitness), "WitnessList: activate witness failed");
        emit WitnessUpdated();
    }

}