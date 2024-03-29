// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

import "../ownership/Ownable.sol";

interface ICrosschainToken {
    function withdrawTo(address, uint256) external;
}

contract Unwrapper is Ownable {
    mapping(address => bool) private whitelist;

    constructor() Ownable() {

    }

    function onReceive(address _sender, ICrosschainToken _token, uint256 _amount, bytes calldata _payload) external {
        require(whitelist[msg.sender], "invalid caller");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        _token.withdrawTo(recipient, _amount);
    }

    function addWhitelist(address _addr) external onlyOwner {
        whitelist[_addr] = true;
    }

    function removeWhitelist(address _addr) external onlyOwner {
        whitelist[_addr] = false;
    }
}

