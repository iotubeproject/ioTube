// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

interface IWrapper {
    function usdc_e() external view returns (IERC20Burnable);
    function iousdc() external view returns (IERC20);
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

contract USDCRouter {
    IWrapper public wrapper;
    IERC20Burnable public usdc_e;
    address public iousdc;
    constructor(IWrapper _wrapper) {
        wrapper = _wrapper;
        usdc_e = _wrapper.usdc_e();
        iousdc = address(_wrapper.iousdc());
    }

    function depositTo(address _cashier, address _to, uint256 _amount) public payable {
        require(wrapper.balance() >= _amount, "insufficient balance");
        require(usdc_e.transferFrom(msg.sender, address(this), _amount), "failed to transfer usdc.e");
        wrapper.deposit(_amount);
        require(safeApprove(iousdc, _cashier, _amount), "failed to approve allowance to cashier");
        ICashier(_cashier).depositTo{value: msg.value}(iousdc, _to, _amount);
    }

    function safeApprove(address _token, address _spender, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('approve(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x095ea7b3, _spender, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}
