pragma solidity >= 0.8.0;

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

interface IERC20Burnable is IERC20 {
    function mint(address to, uint256 amount) external returns (bool);
    function burn(uint256 amount) external;
}

contract USDCSwapper {
    IERC20Burnable public usdc_e;
    IERC20 public iousdc;

    constructor(IERC20Burnable _usdc_e, IERC20 _iousdc) {
        usdc_e = _usdc_e;
        iousdc = _iousdc;
    }

    function deposit(uint256 _amount) external {
        require(iousdc.transferFrom(msg.sender, address(this), _amount), "failed to deposit");
        require(usdc_e.mint(msg.sender, _amount), "failed to mint");
    }

    function withdraw(uint256 _amount) external {
        require(balance() >= _amount, "insufficient balance");
        require(usdc_e.transferFrom(msg.sender, address(this), _amount), "failed to withdraw");
        usdc_e.burn(_amount);
        require(iousdc.transfer(msg.sender, _amount), "failed to transfer");
    }

    function balance() public view returns (uint256) {
        return iousdc.balanceOf(address(this));
    }
}
