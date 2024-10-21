// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

interface ICrosschainToken {
    function withdrawTo(address, uint256) external;
    function coToken() external view returns (address);
}

interface IWEth {
    function withdraw(uint256) external;
}

contract Unwrapper is Ownable {
    event Whitelisted(address indexed);
    event Unwhitelisted(address indexed);
    mapping(address => bool) public whitelist;
    IWEth public weth;

    constructor(address _weth) Ownable(msg.sender) {
        weth = IWEth(_weth);
    }

    receive() external payable {
    }

    function onReceive(address _sender, ICrosschainToken _token, uint256 _amount, bytes calldata _payload) external {
        require(whitelist[msg.sender], "invalid caller");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        address coToken = _token.coToken();
        if (coToken == address(weth)) {
            _token.withdrawTo(address(this), _amount);
            IWEth(coToken).withdraw(_amount);
            payable(recipient).transfer(_amount);
            return;
        }
        _token.withdrawTo(recipient, _amount);
    }

    function addWhitelist(address _addr) external onlyOwner {
        whitelist[_addr] = true;
        emit Whitelisted(_addr);
    }

    function removeWhitelist(address _addr) external onlyOwner {
        whitelist[_addr] = false;
        emit Unwhitelisted(_addr);
    }
}
