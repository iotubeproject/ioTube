// SPDX-License-Identifier: MIT
pragma solidity = 0.8.20;

interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
}

interface ISwapper {
    function source_token() external view returns (address);
    function target_token() external view returns (address);

    function deposit(uint256 _amount) external;
}

// Same contract as USDCUnwrapper.sol but with different naming
contract TokenSwapperUnwrapper {
    IERC20 public source_token;
    IERC20 public target_token;
    ISwapper public swapper;

    constructor(ISwapper _swapper) {
        swapper = _swapper;
        source_token = IERC20(_swapper.source_token());
        target_token = IERC20(_swapper.target_token());
    }

    function onReceive(address _sender, address _token, uint256 _amount, bytes calldata _payload) external {
        require(_token == address(source_token), "TokenSwapperUnwrapper: invalid token");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        require(source_token.approve(address(swapper), _amount), "TokenSwapperUnwrapper: approve failed");
        swapper.deposit(_amount);
        require(target_token.transfer(recipient, _amount), "TokenSwapperUnwrapper: transfer failed");
    }
}
