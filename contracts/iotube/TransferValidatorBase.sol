pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../lifecycle/Pausable.sol";

interface Allowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

contract TransferValidatorBase is Ownable, Pausable {
    event WitnessSubmitted(
        address indexed token,
        uint256 indexed index,
        address indexed witness,
        address from,
        address to,
        uint256 amount);

    event Settled(
        address indexed token,
        uint256 indexed index,
        address from,
        address to,
        uint256 amount,
        address[] witnesses);

    struct Transfer {
        address tokenAddr;
        uint256 index;
        address from;
        address to;
        uint256 amount;
        uint256 settleHeight;
    }

    mapping(bytes32 => Transfer) public transfers;
    mapping(bytes32 => mapping(address => uint256)) public witnessMap;
    mapping(bytes32 => address[]) public witnessList;

    Allowlist public whitelistedTokens;
    Allowlist public whitelistedWitnesses;
    
    function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) public pure returns(bytes32) {
        return keccak256(abi.encodePacked(tokenAddr, index, from, to, amount));
    }

    function getStatusInternal(bytes32 key) internal view returns(uint256 settleHeight_, uint256 numOfWhitelistedWitnesses_, uint256 numOfValidWitnesses_, address[] memory witnesses_, bool includingMsgSender_) {
        settleHeight_ = transfers[key].settleHeight;
        numOfWhitelistedWitnesses_ = whitelistedWitnesses.numOfActive();
        witnesses_ = new address[](numOfWhitelistedWitnesses_);
        for (uint256 i = 0; i < witnessList[key].length; i++) {
            if (whitelistedWitnesses.isAllowed(witnessList[key][i])) {
                witnesses_[numOfValidWitnesses_++] = witnessList[key][i];
                if (witnessList[key][i] == msg.sender) {
                    includingMsgSender_ = true;
                }
            }
        }
    }

    function getStatus(address tokenAddr, uint256 index, address from, address to, uint256 amount) public view returns (uint256 settleHeight_, uint256 numOfWhitelistedWitnesses_, uint256 numOfValidWitnesses_, address[] memory witnesses_, bool includingMsgSender_) {
        return getStatusInternal(generateKey(tokenAddr, index, from, to, amount));
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool);

    function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) public whenNotPaused {
        require(whitelistedWitnesses.isAllowed(msg.sender), "not whitelisted witnesses");
        require(whitelistedTokens.isAllowed(tokenAddr), "not whitelisted tokens");
        require(amount != 0, "amount cannot be zero");
        bytes32 key = generateKey(tokenAddr, index, from, to, amount);
        if (transfers[key].amount == 0) {
            transfers[key] = Transfer(tokenAddr, index, from, to, amount, 0);
        }
        (uint256 settleHeight, uint256 numOfWhitelisted, uint256 numOfValid, address[] memory witnesses, bool submitted) = getStatusInternal(key);
        if (settleHeight != 0) {
            return;
        }
        if (!submitted) {
            witnessMap[key][msg.sender] = block.number;
            witnessList[key].push(msg.sender);
            witnesses[numOfValid++] = msg.sender;
            emit WitnessSubmitted(tokenAddr, index, msg.sender, from, to, amount);
        }
        if (numOfValid * 3 > numOfWhitelisted * 2) {
            transfers[key].settleHeight = block.number;
            require(withdrawToken(tokenAddr, to, amount), "withdraw success");
            emit Settled(tokenAddr, index, from, to, amount, witnesses);
        }
    }
}
