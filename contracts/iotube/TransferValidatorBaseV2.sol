pragma solidity <6.0 >=0.4.24;

import "../lifecycle/Pausable.sol";

interface Allowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

contract TransferValidatorBaseV2 is Pausable {
    event Settled(
        address indexed token,
        uint256 indexed index,
        address from,
        address to,
        uint256 amount,
        address[] witnesses);

    mapping(bytes32 => uint256) public settles;

    Allowlist public whitelistedTokens;
    Allowlist public whitelistedWitnesses;
    
    function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) public view returns(bytes32) {
        return keccak256(abi.encodePacked(address(this), tokenAddr, index, from, to, amount));
    }

    function withdrawToken(address _token, address _to, uint256 _amount) internal returns(bool);

    function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes memory signatures) public whenNotPaused {
        require(whitelistedTokens.isAllowed(tokenAddr), "not whitelisted tokens");
        require(amount != 0, "amount cannot be zero");
        require(signatures.length % 65 == 0, "invalid signature length");
        bytes32 key = generateKey(tokenAddr, index, from, to, amount);
        require(settles[key] == 0, "transfer has been settled");
        uint256 numOfSignatures = signatures.length / 65;
        address[] memory witnesses = new address[](numOfSignatures);
        for (uint256 i = 0; i < numOfSignatures; i++) {
            address witness = recover(key, signatures, i * 65);
            require(whitelistedWitnesses.isAllowed(witness), "invalid signature");
            for (uint256 j = 0; j < i; j++) {
                require(witness != witnesses[j], "duplicate witness");
            }
            witnesses[i] = witness;
        }
        require(numOfSignatures * 3 > whitelistedWitnesses.numOfActive() * 2, "insufficient witnesses");
        settles[key] = block.number;
        require(withdrawToken(tokenAddr, to, amount), "failed to withdraw");
        emit Settled(tokenAddr, index, from, to, amount, witnesses);
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
