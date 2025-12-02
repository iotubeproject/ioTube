// SPDX-License-Identifier: MIT
pragma solidity >=0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

interface IWitnessList {
    function numOfActive() external view returns (uint256);
    function areAllowed(address[] calldata _witnesses) external view returns (bool);
    function addWitnesses(address[] memory _witnesses) external;
    function removeWitnesses(address[] memory _witnesses) external;
    function owner() external view returns (address);
    function transferOwnership(address newOwner) external;
}

contract WitnessManager is Ownable {
    using EnumerableSet for EnumerableSet.AddressSet;
    event WitnessesProposed(uint64 indexed epochNum, bytes32 witnessesHash);

    IWitnessList public witnessList;

    // protocol-related
    uint64 public epochNum;
    uint64 public epochInterval;
    EnumerableSet.AddressSet private _excludedList;
    // TODO: define x out of 36 number in the contract  

    constructor(address initialOwner, IWitnessList _witnessList) Ownable(initialOwner) {
        witnessList = _witnessList;
    }

    function addWitnesses(address[] memory _witnesses) public onlyOwner {
        witnessList.addWitnesses(_witnesses);
    }

    function removeWitnesses(address[] memory _witnesses) public onlyOwner {
        witnessList.removeWitnesses(_witnesses);
    }

    function setEpochNum(uint64 _epochNum) public onlyOwner {
        epochNum = _epochNum;
    }
    
    function setEpochInterval(uint64 _epochInterval) public onlyOwner {
        epochInterval = _epochInterval;
    }

    function addExcludedWitnesses(address[] calldata _witnesses) public onlyOwner {
        for (uint256 i = 0; i < _witnesses.length; i++) {
            _excludedList.add(_witnesses[i]);
        }
    }

    function removeExcludedWitnesses(address[] calldata _witnesses) public onlyOwner {
        for (uint256 i = 0; i < _witnesses.length; i++) {
            require(_excludedList.remove(_witnesses[i]), "witness not found");
        }
    }

    function getExcludedWitnesses() public view returns (address[] memory) {
        return _excludedList.values();
    }

    /**
     * @notice Propose a new witness list for the next epoch.
     * @param nextEpochNum The next epoch number.
     * @param witnessesToAdd The list of witnesses to add, sorted in lexicographical order.
     * @param witnessesToRemove The list of witnesses to remove, sorted in lexicographical order.
     * @param signatures The signatures of the current witnesses.
     */
    function proposeWitnesses(uint64 nextEpochNum, address[] calldata witnessesToAdd, address[] calldata witnessesToRemove, bytes[] calldata signatures) public {
        // validate epoch number
        require(nextEpochNum == epochNum + epochInterval, "invalid next epoch number");
        // require(nextEpochNum >  epochNum && (nextEpochNum - epochNum) % epochInterval == 0, "invalid next epoch number");

        // extract hash from nextEpochNum and witnesses
        bytes32 hash;
        bytes memory packed = abi.encode(address(this), nextEpochNum, witnessesToAdd, witnessesToRemove);
        assembly {
            hash := keccak256(add(packed, 32), mload(packed))
        }

        // validate current witness signature
        uint256 signatureCount = signatures.length;
        require(signatureCount * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");

        address[] memory signers = new address[](signatureCount);
        for (uint256 i = 0; i < signatureCount; i++) {
            address signer = ECDSA.recover(hash, signatures[i]);
            for (uint j = 0; j < i; j++) {
                require(signers[j] != signer, "duplicate signature");
            }
            signers[i] = signer;
        }
        require(witnessList.areAllowed(signers), "invalid signatures");

        // update witness list
        updateWitnessList(witnessesToAdd, witnessesToRemove);
        epochNum = nextEpochNum;

        emit WitnessesProposed(nextEpochNum, hash);
    }

    function updateWitnessList(address[] calldata witnessesToAdd, address[] calldata witnessesToRemove) internal {
        uint256 currentWitnessCount = witnessList.numOfActive();
        
        witnessList.removeWitnesses(witnessesToRemove);
        
        for (uint256 i = 0; i < witnessesToAdd.length; i++) {
            require(witnessesToAdd[i] != address(0), "invalid new witness");
        }
        witnessList.addWitnesses(witnessesToAdd);

        require(witnessList.numOfActive()  == currentWitnessCount + witnessesToAdd.length - witnessesToRemove.length, "witness count mismatch");
    }

    function transferWitnessListOwnership(address newOwner) public onlyOwner {
        require(witnessList.owner() == address(this), "witness list is not owned by this contract");
        witnessList.transferOwnership(newOwner);
    }
}