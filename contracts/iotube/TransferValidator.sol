pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../lifecycle/Pausable.sol";
import "../math/SafeMath.sol";
import "./Safe.sol";

contract TransferValidator is Ownable {
    using SafeMath for uint256;

    event Settled(
        address token,
        uint indexed index,
        address indexed from,
        address indexed to,
        uint256 amount,
        uint256 blockNumber,
        address[] voters);
    event TokenAddedToWhitelist(address token);
    event VoterAddedToWhitelist(address voter);
    event TokenRemovedFromWhitelist(address token);
    event VoterRemovedFromWhitelist(address voter);

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

    struct WhitelistedToken {
        Safe safe;
        bool flag;
    }

    mapping(bytes32 => Transfer) public transfers;
    mapping(bytes32 => Vote[]) public votes;
    mapping(address => WhitelistedToken) public whitelistedTokens;
    uint256 public numOfWhitelistedTokens;
    mapping(address => bool) public whitelistedVoters;
    uint256 public numOfWhitelistedVoters;
    uint256 public expireHeight;

    constructor(uint256 _expireHeight) public {
        setExpireHeight(_expireHeight);
    }

    function setExpireHeight(uint256 _expireHeight) public onlyOwner {
        expireHeight = _expireHeight;
    }

    function addWhitelistedToken(address _token, address _safe) public onlyOwner returns (bool success_) {
        if (!whitelistedTokens[_token].flag) {
            whitelistedTokens[_token] = WhitelistedToken(Safe(_safe), true);
            numOfWhitelistedTokens++;
            emit TokenAddedToWhitelist(_token);
            success_ = true;
        }
    }

    function addWhitelistedTokens(address[] memory _tokens, address[] memory _safes) public onlyOwner returns (bool success_) {
        require(_tokens.length == _safes.length, "tokens and safes do not match");
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (addWhitelistedToken(_tokens[i], _safes[i])) {
                success_ = true;
            }
        }
    }

    function removeWhitelistedToken(address _token) public onlyOwner returns (bool success_) {
        if (whitelistedTokens[_token].flag) {
            whitelistedTokens[_token].flag = false;
            numOfWhitelistedTokens--;
            emit TokenRemovedFromWhitelist(_token);
            success_ = true;
        }
    }

    function addWhitelistedVoter(address _voter) public onlyOwner returns (bool success_) {
        if (!whitelistedVoters[_voter]) {
            whitelistedVoters[_voter] = true;
            numOfWhitelistedVoters++;
            emit VoterAddedToWhitelist(_voter);
            success_ = true;
        }
    }

    function addWhitelistedVoters(address[] memory voters) public onlyOwner returns (bool success_) {
        for (uint256 i = 0; i < voters.length; i++) {
            if (addWhitelistedVoter(voters[i])) {
                success_ = true;
            }
        }
    }

    function removeWhitelistedVoter(address _voter) public onlyOwner returns (bool success_) {
        if (whitelistedVoters[_voter]) {
            whitelistedVoters[_voter] = false;
            numOfWhitelistedVoters--;
            emit VoterRemovedFromWhitelist(_voter);
            success_ = true;
        }
    }

    function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) public pure returns(bytes32) {
        return keccak256(abi.encodePacked(tokenAddr, index, from, to, amount));
    }

    function settled(bytes32 key) public view returns(bool) {
        return transfers[key].settleHeight > 0;
    }

    function vote(address tokenAddr, uint256 index, address from, address to, uint256 amount) public {
        require(whitelistedVoters[msg.sender], "not whitelisted voters");
        require(whitelistedTokens[tokenAddr].flag, "not whitelisted tokens");
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
        address[] memory voters = new address[](numOfWhitelistedVoters);
        for (i = 0; i < l; i++) {
            if (votes[key][i].voter == msg.sender) {
                votes[key][i].blockNumber = block.number;
                isUpdate = true;
            }
            // vote's block number is always less or equal to block.number
            if ((expireHeight >= block.number - votes[key][i].blockNumber) && whitelistedVoters[votes[key][i].voter]) {
                voters[numOfValidVoters] = votes[key][i].voter;
                numOfValidVoters++;
            }
        }
        if (!isUpdate) {
            votes[key].push(Vote(msg.sender, block.number));
            voters[numOfValidVoters] = msg.sender;
            numOfValidVoters++;
        }
        if (numOfValidVoters * 3 > numOfWhitelistedVoters * 2) {
            transfers[key].settleHeight = block.number;
            address[] memory trimmedVoters = new address[](numOfValidVoters);
            for (i = 0; i < numOfValidVoters; i++) {
                trimmedVoters[i] = voters[i];
            }
            require(whitelistedTokens[tokenAddr].safe.withdraw(to, amount), "withdraw success");
            emit Settled(tokenAddr, index, from, to, amount, block.number, trimmedVoters);
        }
    }
}