pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../lifecycle/Pausable.sol";

interface Allowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

contract TransferValidatorBase is Ownable, Pausable {
    event Settled(
        address indexed token,
        uint256 indexed index,
        address indexed from,
        address to,
        uint256 amount,
        uint256 blockNumber,
        address[] witnesses);

    struct Submission {
        address witness;
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
    mapping(bytes32 => Submission[]) public submissions;

    Allowlist public whitelistedTokens;
    Allowlist public whitelistedWitnesss;
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

    function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) public whenNotPaused {
        require(whitelistedWitnesss.isAllowed(msg.sender), "not whitelisted witnesses");
        require(whitelistedTokens.isAllowed(tokenAddr), "not whitelisted tokens");
        bytes32 key = generateKey(tokenAddr, index, from, to, amount);
        if (settled(key)) {
            return;
        }
        if (!transfers[key].flag) {
            transfers[key] = Transfer(tokenAddr, index, from, to, amount, 0, true);
        }
        uint256 l = submissions[key].length;
        uint256 numOfValidWitnesss = 0;
        uint256 i;
        bool isUpdate = false;
        address[] memory witnesses = new address[](whitelistedWitnesss.numOfActive());
        for (i = 0; i < l; i++) {
            if (submissions[key][i].witness == msg.sender) {
                submissions[key][i].blockNumber = block.number;
                isUpdate = true;
            }
            // block number is always less or equal to block.number
            if ((expireHeight >= block.number - submissions[key][i].blockNumber) && whitelistedWitnesss.isAllowed(submissions[key][i].witness)) {
                witnesses[numOfValidWitnesss] = submissions[key][i].witness;
                numOfValidWitnesss++;
            }
        }
        if (!isUpdate) {
            submissions[key].push(Submission(msg.sender, block.number));
            witnesses[numOfValidWitnesss] = msg.sender;
            numOfValidWitnesss++;
        }
        if (numOfValidWitnesss * 3 > whitelistedWitnesss.numOfActive() * 2) {
            transfers[key].settleHeight = block.number;
            address[] memory trimmedWitnesss = new address[](numOfValidWitnesss);
            for (i = 0; i < numOfValidWitnesss; i++) {
                trimmedWitnesss[i] = witnesses[i];
            }
            require(withdrawToken(tokenAddr, to, amount), "withdraw success");
            emit Settled(tokenAddr, index, from, to, amount, block.number, trimmedWitnesss);
        }
    }
}