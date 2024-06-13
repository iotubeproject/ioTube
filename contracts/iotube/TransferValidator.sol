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

    function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes memory signatures) public whenNotPaused {
        require(amount != 0, "amount cannot be zero");
        require(to != address(0), "recipient cannot be zero");
        require(signatures.length % 65 == 0, "invalid signature length");
        bytes32 key = generateKey(cashier, tokenAddr, index, from, to, amount);
        require(settles[key] == 0, "transfer has been settled");
        for (uint256 it = 0; it < tokenLists.length; it++) {
            if (tokenLists[it].isAllowed(tokenAddr)) {
                uint256 numOfSignatures = signatures.length / 65;
                address[] memory witnesses = new address[](numOfSignatures);
                for (uint256 i = 0; i < numOfSignatures; i++) {
                    address witness = recover(key, signatures, i * 65);
                    require(witnessList.isAllowed(witness), "invalid signature");
                    for (uint256 j = 0; j < i; j++) {
                        require(witness != witnesses[j], "duplicate witness");
                    }
                    witnesses[i] = witness;
                }
                require(numOfSignatures * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");
                settles[key] = block.number;
                require(minters[it].mint(tokenAddr, to, amount), "failed to mint token");
                emit Settled(key, witnesses);

                return;
            }
        }
        revert("not a whitelisted token");
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
