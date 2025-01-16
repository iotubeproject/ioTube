// SPDX-License-Identifier: MIT
pragma solidity <6.0 >=0.4.24;
pragma experimental ABIEncoderV2;

import "../lifecycle/Pausable.sol";

interface IAllowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

interface IMinter {
    function mint(address, address, uint256) external returns (bool);
    function transferOwnership(address) external;
    function owner() external view returns (address);
}

interface IReceiver {
    function onReceive(address sender, address token, uint256 amount, bytes calldata payload) external;
}

contract TransferValidatorForSolana is Pausable {
    event Settled(bytes32 indexed key, address[] witnesses);
    event ReceiverAdded(address receiver);
    event ReceiverRemoved(address receiver);
    mapping(bytes32 => uint256) public settles;
    mapping(address => bool) public receivers;

    IMinter[] public minters;
    IAllowlist[] public tokenLists;
    IAllowlist public witnessList;

    constructor(IAllowlist _witnessList) public {
        witnessList = _witnessList;
    }

    function generateKey(
        bytes memory cashier,
        address tokenAddr,
        uint256 index,
        bytes memory from,
        address to,
        uint256 amount,
        bytes memory payload
    ) public view returns (bytes32) {
        return keccak256(abi.encodePacked(address(this), cashier, tokenAddr, index, from, to, amount, payload));
    }

    function getTokenGroup(address tokenAddr) public view returns (uint256) {
        for (uint256 i = 0; i < tokenLists.length; i++) {
            if (tokenLists[i].isAllowed(tokenAddr)) {
                return i;
            }
        }
        require(false, "invalid token address");
    }

    function extractWitnesses(bytes32 key, bytes memory signatures) public view returns (address[] memory witnesses) {
        uint256 numOfSignatures = signatures.length / 65;
        witnesses = new address[](numOfSignatures);
        for (uint256 i = 0; i < numOfSignatures; i++) {
            address witness = recover(key, signatures, i * 65);
            require(witnessList.isAllowed(witness), "invalid signatures");
            for (uint256 j = 0; j < i; j++) {
                require(witness != witnesses[j], "duplicate witness");
            }
            witnesses[i] = witness;
        }
    }

    function submit(bytes memory cashier, address tokenAddr, uint256 index, bytes memory from,
        address to, uint256 amount, bytes memory signatures, bytes memory payload) public whenNotPaused {
        require(amount != 0, "amount cannot be zero");
        require(to != address(0), "recipient cannot be zero");
        require(signatures.length % 65 == 0, "invalid signature length");
        require((from.length == 20 && cashier.length == 20) || (from.length == 32 && cashier.length == 32), "illegal address length");
        bytes32 key = generateKey(cashier, tokenAddr, index, from, to, amount, payload);
        require(settles[key] == 0, "transfer has been settled");
        address[] memory witnesses = extractWitnesses(key, signatures);
        require(witnesses.length > 0 && witnesses.length * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");
        settles[key] = block.number;
        require(minters[getTokenGroup(tokenAddr)].mint(tokenAddr, to, amount), "failed to mint token");
        if (receivers[to]) {
            IReceiver(to).onReceive(this.owner(), tokenAddr, amount, payload);
        }
        emit Settled(key, witnesses);
    }

    function numOfPairs() external view returns (uint256) {
        return tokenLists.length;
    }

    function addPair(IAllowlist _tokenList, IMinter _minter) external onlyOwner {
        tokenLists.push(_tokenList);
        minters.push(_minter);
    }

    function addReceiver(address _receiver) external onlyOwner {
        require(!receivers[_receiver], "already a receiver");
        receivers[_receiver] = true;
        emit ReceiverAdded(_receiver);
    }

    function removeReceiver(address _receiver) external onlyOwner {
        require(receivers[_receiver], "invalid receiver");
        receivers[_receiver] = false;
        emit ReceiverRemoved(_receiver);
    }

    function upgrade(address _newValidator) external onlyOwner {
        address contractAddr = address(this);
        for (uint256 i = 0; i < minters.length; i++) {
            IMinter minter = minters[i];
            if (minter.owner() == contractAddr) {
                minter.transferOwnership(_newValidator);
            }
        }
    }

    /**
     * @dev Recover signer address from a message by using their signature
     * @param hash bytes32 message, the hash is the signed message. What is recovered is the signer address.
     * @param signature bytes signature, the signature is generated using web3.eth.sign()
     */
    function recover(bytes32 hash, bytes memory signature, uint256 offset) internal pure returns (address) {
        bytes32 r;
        bytes32 s;
        uint8 v;

        // Divide the signature in r, s and v variables with inline assembly.
        assembly {
            r := mload(add(signature, add(offset, 0x20)))
            s := mload(add(signature, add(offset, 0x40)))
            v := byte(0, mload(add(signature, add(offset, 0x60))))
        }

        // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
        if (v < 27) {
            v += 27;
        }

        // If the version is correct return the signer address
        if (v != 27 && v != 28) {
            return (address(0));
        }
        // solium-disable-next-line arg-overflow
        return ecrecover(hash, v, r, s);
    }
}
