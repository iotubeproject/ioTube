// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

interface ISwapper {
    function target_token() external view returns (IERC20);
    function source_token() external view returns (IERC20);
    function deposit(uint256 _amount) external;
    function withdraw(uint256 _amount) external;
    function balance() external view returns (uint256);
}

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

interface ICashier {
    function depositTo(address _token, address _to, uint256 _amount, bytes memory _payload) external payable;
}

contract SwapperCashierRouter {
    ISwapper public swapper;
    address public target_token;
    address public source_token;
    constructor(ISwapper _swapper) {
        swapper = _swapper;
        target_token = address(_swapper.target_token());
        source_token = address(_swapper.source_token());
        _approve();
    }

    function _approve() private {
        require(safeApprove(target_token, address(swapper), type(uint256).max), "failed to approve allowance to swapper");
    }

    function depositTo(address _cashier, address _to, uint256 _amount, bytes memory _payload) public payable {
        require(swapper.balance() >= _amount, "swapper insufficient balance");
        require(safeTransferFrom(target_token, msg.sender, address(this), _amount), "failed to transfer target token");
        swapper.withdraw(_amount);
        require(safeApprove(source_token, _cashier, _amount), "failed to approve allowance to cashier");
        ICashier(_cashier).depositTo{value: msg.value}(source_token, _to, _amount, _payload);
    }

    function safeTransferFrom(address _token, address _from, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, _from, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

    function safeApprove(address _token, address _spender, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('approve(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x095ea7b3, _spender, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}