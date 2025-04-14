// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

interface IWrapper {
    function target_token() external view returns (IERC20Burnable);
    function source_token() external view returns (IERC20);
    function deposit(uint256 _amount) external;
    function withdraw(uint256 _amount) external;
    function balance() external view returns (uint256);
}

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

interface ICashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface IERC20Burnable is IERC20 {
    function mint(address to, uint256 amount) external returns (bool);
    function burn(uint256 amount) external;
}

contract WrapperCashierRouter {
    IWrapper public wrapper;
    IERC20Burnable public target_token;
    address public source_token;
    constructor(IWrapper _wrapper) {
        wrapper = _wrapper;
        target_token = _wrapper.target_token();
        source_token = address(_wrapper.source_token());
        approveWrapper();
    }

    function approveWrapper() public {
        require(safeApprove(address(target_token), address(wrapper), type(uint256).max), "failed to approve allowance to wrapper");
    }

    function depositTo(address _cashier, address _to, uint256 _amount) public payable {
        require(wrapper.balance() >= _amount, "insufficient balance");
        require(target_token.transferFrom(msg.sender, address(this), _amount), "failed to transfer target token");
        wrapper.withdraw(_amount);
        require(safeApprove(source_token, _cashier, _amount), "failed to approve allowance to cashier");
        ICashier(_cashier).depositTo{value: msg.value}(source_token, _to, _amount);
    }

    function safeApprove(address _token, address _spender, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('approve(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x095ea7b3, _spender, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}