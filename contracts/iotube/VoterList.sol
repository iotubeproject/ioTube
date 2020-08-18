pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

contract VoterList is Ownable {
    event VoterAdded(address indexed voter);
    event VoterRemoved(address indexed voter);
    struct Setting {
        bool active;
        bool flag;
    }

    mapping(address => Setting) public settings;
    address[] public voters;
    uint256 private num;
    
    function isAllowed(address _voter) public view returns (bool) {
        return settings[_voter].active;
    }

    function addVoter(address _voter) public onlyOwner returns (bool success_) {
        if (!settings[_voter].active) {
            if (!settings[_voter].flag) {
                voters.push(_voter);
            }
            settings[_voter] = Setting(true, true);
            num++;
            emit VoterAdded(_voter);
            success_ = true;
        }
    }

    function removeVoter(address _voter) public onlyOwner returns (bool success_) {
        if (settings[_voter].active) {
            settings[_voter].active = false;
            num--;
            emit VoterRemoved(_voter);
            success_ = true;
        }
    }

    function numOfAllowed() public view returns (uint256) {
        return num;
    }

}