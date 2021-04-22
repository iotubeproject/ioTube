pragma solidity <6.0 >=0.4.24;

import "../lifecycle/Pausable.sol";

interface IAllowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

interface IMinter {
    function mint(address, address, uint256) external returns(bool);
    function transferOwnership(address) external;
    function owner() external view returns(address);
}

contract TransferValidator is Pausable {
    event Settled(bytes32 indexed key, address[] witnesses);

    mapping(bytes32 => uint256) public settles;

    IMinter[] public minters;
    IAllowlist[] public tokenLists;
    IAllowlist public witnessList;

    constructor(IAllowlist _witnessList) public {
        witnessList = _witnessList;
    }

    function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) public view returns(bytes32) {
        return keccak256(abi.encodePacked(address(this), cashier, tokenAddr, index, from, to, amount));
    }

    function concatKeys(bytes32[] memory keys) public pure returns(bytes32) {
        return keccak256(abi.encodePacked(keys));
    }

    function submitMulti(address[] memory cashiers, address[] memory tokenAddrs, uint256[] memory indexes, address[] memory senders, address[] memory recipients, uint256[] memory amounts, bytes memory signatures) public whenNotPaused {
        require(cashiers.length == tokenAddrs.length && tokenAddrs.length == indexes.length && indexes.length == senders.length && senders.length == recipients.length && recipients.length == amounts.length, "invalid parameters");
        bytes32[] memory keys = new bytes32[](cashiers.length);
        for (uint256 i = 0; i < cashiers.length; i++) {
            keys[i] = generateKey(cashiers[i], tokenAddrs[i], indexes[i], senders[i], recipients[i], amounts[i]);
            require(settles[keys[i]] == 0, "transfer has been settled");
            for (uint256 j = 0; j < i; j++) {
                require(keys[i] != keys[j], "duplicate key");
            }
        }
        address[] memory witnesses = extractWitnesses(concatKeys(keys), signatures);
        require(witnesses.length > 0 && witnesses.length * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");
        for (uint256 i = 0; i < cashiers.length; i++) {
            settles[keys[i]] = block.number;
            require(minters[getTokenGroup(tokenAddrs[i])].mint(tokenAddrs[i], recipients[i], amounts[i]), "failed to mint token");
            emit Settled(keys[i], witnesses);
        }
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

    function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes memory signatures) public whenNotPaused {
        require(amount != 0, "amount cannot be zero");
        require(to != address(0), "recipient cannot be zero");
        require(signatures.length % 65 == 0, "invalid signature length");
        bytes32 key = generateKey(cashier, tokenAddr, index, from, to, amount);
        require(settles[key] == 0, "transfer has been settled");
        address[] memory witnesses = extractWitnesses(key, signatures);
        require(witnesses.length > 0 && witnesses.length * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");
        settles[key] = block.number;
        require(minters[getTokenGroup(tokenAddr)].mint(tokenAddr, to, amount), "failed to mint token");
        emit Settled(key, witnesses);
    }

    function numOfPairs() external view returns (uint256) {
        return tokenLists.length;
    }

    function addPair(IAllowlist _tokenList, IMinter _minter) external onlyOwner {
        tokenLists.push(_tokenList);
        minters.push(_minter);
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
    function recover(bytes32 hash, bytes memory signature, uint256 offset)
        internal
        pure
        returns (address)
    {
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
