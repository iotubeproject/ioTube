// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

interface IAllowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

interface IMinter {
    function mint(address, address, uint256) external returns(bool);
    function transferOwnership(address) external;
    function owner() external view returns(address);
}

interface IReceiver {
    function onReceive(address sender, address token, uint256 amount, bytes calldata payload) external;
}

contract TransferValidatorV3 is Ownable, Pausable {
    event Settled(bytes32 indexed key, address[] witnesses);
    event ReceiverAdded(address receiver);
    event ReceiverRemoved(address receiver);

    mapping(bytes32 => uint256) public settles;
    mapping(address => bool) public receivers;

    IMinter[] public minters;
    IAllowlist[] public tokenLists;
    IAllowlist[] public witnessLists;
    IAllowlist[] public witnessTokenLists;


    constructor() Ownable(msg.sender) {}

    function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes memory payload) internal view returns(bytes32) {
        return keccak256(abi.encodePacked(address(this), cashier, tokenAddr, index, from, to, amount, payload));
    }

    function getWitnessLists(address tokenAddr) public view returns (IAllowlist[] memory) {
        uint256 count = 0;
        bool[] memory found = new bool[](witnessTokenLists.length);
        for (uint256 i = 0; i < witnessTokenLists.length; i++) {
            if (witnessTokenLists[i].isAllowed(tokenAddr)) {
                count++;
                found[i] = true;
            }
        }
        if (count == 0) {
            revert("Witness list not found");
        }

        IAllowlist[] memory lists = new IAllowlist[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < witnessTokenLists.length; i++) {
            if (found[i]) {
                lists[index] = witnessLists[i];
                index++;
            }
        }
        return lists;
    }

    function getMinter(address tokenAddr) public view returns (IMinter) {
        for (uint256 i = 0; i < tokenLists.length; i++) {
            if (tokenLists[i].isAllowed(tokenAddr)) {
                return minters[i];
            }
        }
        revert("invalid token address");
    }

    function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes[] memory signaturesArray, bytes memory payload) public whenNotPaused {
        require(amount != 0, "amount cannot be zero");
        require(to != address(0), "recipient cannot be zero");
        bytes32 key = generateKey(cashier, tokenAddr, index, from, to, amount, payload);
        require(settles[key] == 0, "transfer has been settled");
        IAllowlist[] memory witnessListArray = getWitnessLists(tokenAddr);
        require(witnessListArray.length == signaturesArray.length, "invalid signature length");
        uint256 totalSignatures = 0;
        for (uint256 i = 0; i < signaturesArray.length; i++) {
            require(signaturesArray[i].length % 65 == 0, "invalid signature length");
            totalSignatures += signaturesArray[i].length / 65;
        }
        address[] memory allWitnesses = new address[](totalSignatures);
        uint256 witnessIndex = 0;
        for (uint256 i = 0; i < witnessListArray.length; i++) {
            IAllowlist witnessList = witnessListArray[i];
            bytes memory signatures = signaturesArray[i];
            uint256 numOfSignatures = signatures.length / 65;
            address[] memory witnesses = new address[](numOfSignatures);
            for (uint256 j = 0; j < numOfSignatures; j++) {
                // TODO: check sig copy gas
                address witness = recover(key, signatures, j * 65);
                require(witnessList.isAllowed(witness), "invalid signature");
                for (uint256 k = 0; k < j; k++) {
                    require(witness != witnesses[k], "duplicate witness");
                }
                witnesses[j] = witness;
                allWitnesses[witnessIndex] = witness;
                witnessIndex++;
            }
            require(numOfSignatures * 3 > witnessList.numOfActive() * 2, "insufficient witnesses");
        }
        IMinter minter = getMinter(tokenAddr);
        settles[key] = block.number;
        require(minter.mint(tokenAddr, to, amount), "failed to mint token");
        if (receivers[to]) {
            IReceiver(to).onReceive(from, tokenAddr, amount, payload);
        }
        emit Settled(key, allWitnesses);
    }

    function numOfPairs() external view returns (uint256) {
        return tokenLists.length;
    }

    function addMinterPair(IAllowlist _minterTokenList, IMinter _minter) external onlyOwner {
        tokenLists.push(_minterTokenList);
        minters.push(_minter);
    }

    function addWitnessPair(IAllowlist _witnessTokenList, IAllowlist _witnessList) external onlyOwner {
        witnessTokenLists.push(_witnessTokenList);
        witnessLists.push(_witnessList);
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
