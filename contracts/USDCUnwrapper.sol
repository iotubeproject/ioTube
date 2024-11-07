// SPDX-License-Identifier: MIT
pragma solidity = 0.8.20;

interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
}

interface ISwapper {
    function usdc_e() external view returns (address);
    function iousdc() external view returns (address);

    function deposit(uint256 _amount) external;
}

contract USDCUnwrapper {
    IERC20 public iousdc;
    IERC20 public usdc_e;
    ISwapper public swapper;

    constructor(ISwapper _swapper) {
        swapper = _swapper;
        iousdc = IERC20(_swapper.iousdc());
        usdc_e = IERC20(_swapper.usdc_e());
    }

    function onReceive(address _sender, address _token, uint256 _amount, bytes calldata _payload) external {
        require(_token == address(iousdc), "USDCUnwrapper: invalid token");
        address recipient = _sender;
        if (_payload.length == 32) {
            (recipient) = abi.decode(_payload, (address));
        }
        require(iousdc.approve(address(swapper), _amount), "USDCUnwrapper: approve failed");
        swapper.deposit(_amount);
        require(usdc_e.transfer(recipient, _amount), "USDCUnwrapper: transfer failed");
    }
}
