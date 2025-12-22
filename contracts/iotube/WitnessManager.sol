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
    function threshold() external view returns (uint8);
    function setThreshold(uint8 _threshold) external;
}

contract WitnessManager is Ownable {
    using EnumerableSet for EnumerableSet.AddressSet;
    event WitnessesProposed(uint64 indexed epochNum, bytes32 witnessesHash);

    IWitnessList public witnessList;

    // protocol-related
    uint64 public epochNum;
    uint64 public epochInterval;
   
    // witness setting related
    uint64 public numNominees;
    uint64 private _nextNumNominees;
    EnumerableSet.AddressSet private _excludedList;
    address[] private _nextAddExcludedList;
    address[] private _nextRemoveExcludedList;

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

    function setNumNominees(uint64 _numNominees) public onlyOwner {
        numNominees = _numNominees;
    }

    /**
     * @notice Set the number of nominees for the next epoch.
     * @param nextNumNominees The number of nominees for the next epoch.
     */
    function setNextNumNominees(uint64 nextNumNominees) public onlyOwner {
        _nextNumNominees = nextNumNominees;
    }

    function setExcludedWitnesses(address[] calldata _witnesses) public onlyOwner {
        _excludedList.clear();
        for (uint256 i = 0; i < _witnesses.length; i++) {
            _excludedList.add(_witnesses[i]);
        }
    }

    /**
     * @notice Set the list of witnesses to be added to the excluded list in the next epoch.
     * @param _witnesses The list of witnesses to be added to the excluded list.
     */
    function setNextAddExcludedWitnesses(address[] calldata _witnesses) public onlyOwner {
        _nextAddExcludedList = _witnesses;
    }

    /**
     * @notice Set the list of witnesses to be removed from the excluded list in the next epoch.
     * @param _witnesses The list of witnesses to be removed from the excluded list.
     */
    function setNextRemoveExcludedWitnesses(address[] calldata _witnesses) public onlyOwner {
        _nextRemoveExcludedList = _witnesses;
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
        require(signatureCount * 100 > witnessList.numOfActive() * witnessList.threshold(), "insufficient witnesses");

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
        flushNextWitnessSettings();

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
        require(witnessList.numOfActive() == numNominees, "nominee count mismatch");
    }

    /**
    * @notice Flush the pending witness settings to the current settings.
    */
    function flushNextWitnessSettings() internal {
        if (_nextNumNominees != 0) {
           numNominees = _nextNumNominees;
           _nextNumNominees = 0;
        }
        if (_nextAddExcludedList.length > 0) {
            for (uint256 i = 0; i < _nextAddExcludedList.length; i++) {
                _excludedList.add(_nextAddExcludedList[i]);
            }
            delete _nextAddExcludedList;
        }
        if (_nextRemoveExcludedList.length > 0) {
            for (uint256 i = 0; i < _nextRemoveExcludedList.length; i++) {
                _excludedList.remove(_nextRemoveExcludedList[i]);
            }
            delete _nextRemoveExcludedList;
        }
    }
    
    function setWitnessListThreshold(uint8 _threshold) public onlyOwner {
        witnessList.setThreshold(_threshold);
    }

    function transferWitnessListOwnership(address newOwner) public onlyOwner {
        require(witnessList.owner() == address(this), "witness list is not owned by this contract");
        witnessList.transferOwnership(newOwner);
    }
}