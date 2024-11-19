// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

interface IWIOTX {
    function withdraw(uint256) external;
}

contract WIOTXUnwrapper {
    address public wiotx;

    constructor(address _wiotx) {
        wiotx = _wiotx;
    }

    receive() external payable {
    }

    function onReceive(address _sender, IWIOTX _token, uint256 _amount, bytes calldata _payload) external {
        require(address(_token) == wiotx, "WIOTXUnwrapper: invalid token");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        _token.withdraw(_amount);
        payable(recipient).transfer(_amount);
    }
}
