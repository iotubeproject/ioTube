pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../lifecycle/Pausable.sol";

interface Allowlist {
    function isAllowed(address) external view returns (bool);
    function numOfAllowed() external view returns (uint256);
}

contract TransferValidatorBase is Ownable, Pausable {
    event Settled(
        address indexed token,
        uint256 indexed index,
        address indexed from,
        address to,
        uint256 amount,
        uint256 blockNumber,
        address[] voters);

    struct Vote {
        address voter;
        uint256 blockNumber;
    }

    struct Transfer {
        address tokenAddr;
        uint256 index;
        address from;
        address to;
        uint256 amount;
        uint256 settleHeight;
        bool flag;
    }

    mapping(bytes32 => Transfer) public transfers;
    mapping(bytes32 => Vote[]) public votes;

    Allowlist public whitelistedTokens;
    Allowlist public whitelistedVoters;
    uint256 public expireHeight;
    
    function setExpireHeight(uint256 _expireHeight) public onlyOwner {
        expireHeight = _expireHeight;
    }

    function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) public pure returns(bytes32) {
        return keccak256(abi.encodePacked(tokenAddr, index, from, to, amount));
    }

    function settled(bytes32 key) public view returns(bool) {
        return transfers[key].settleHeight > 0;
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool);

    function vote(address tokenAddr, uint256 index, address from, address to, uint256 amount) public whenNotPaused {
        require(whitelistedVoters.isAllowed(msg.sender), "not whitelisted voters");
        require(whitelistedTokens.isAllowed(tokenAddr), "not whitelisted tokens");
        bytes32 key = generateKey(tokenAddr, index, from, to, amount);
        if (settled(key)) {
            return;
        }
        if (!transfers[key].flag) {
            transfers[key] = Transfer(tokenAddr, index, from, to, amount, 0, true);
        }
        uint256 l = votes[key].length;
        uint256 numOfValidVoters = 0;
        uint256 i;
        bool isUpdate = false;
        address[] memory voters = new address[](whitelistedVoters.numOfAllowed());
        for (i = 0; i < l; i++) {
            if (votes[key][i].voter == msg.sender) {
                votes[key][i].blockNumber = block.number;
                isUpdate = true;
            }
            // vote's block number is always less or equal to block.number
            if ((expireHeight >= block.number - votes[key][i].blockNumber) && whitelistedVoters.isAllowed(votes[key][i].voter)) {
                voters[numOfValidVoters] = votes[key][i].voter;
                numOfValidVoters++;
            }
        }
        if (!isUpdate) {
            votes[key].push(Vote(msg.sender, block.number));
            voters[numOfValidVoters] = msg.sender;
            numOfValidVoters++;
        }
        if (numOfValidVoters * 3 > whitelistedVoters.numOfAllowed() * 2) {
            transfers[key].settleHeight = block.number;
            address[] memory trimmedVoters = new address[](numOfValidVoters);
            for (i = 0; i < numOfValidVoters; i++) {
                trimmedVoters[i] = voters[i];
            }
            require(withdrawToken(tokenAddr, to, amount), "withdraw success");
            emit Settled(tokenAddr, index, from, to, amount, block.number, trimmedVoters);
        }
    }
}