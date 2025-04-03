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

// Same contract as USDCSwapper.sol but with different naming
contract TokenSwapper {
    IERC20Burnable public target_token;
    IERC20 public source_token;

    constructor(IERC20Burnable _target_token, IERC20 _source_token) {
        target_token = _target_token;
        source_token = _source_token;
    }

    function deposit(uint256 _amount) external {
        require(source_token.transferFrom(msg.sender, address(this), _amount), "failed to deposit");
        require(target_token.mint(msg.sender, _amount), "failed to mint");
    }

    function withdraw(uint256 _amount) external {
        require(balance() >= _amount, "insufficient balance");
        require(target_token.transferFrom(msg.sender, address(this), _amount), "failed to withdraw");
        target_token.burn(_amount);
        require(source_token.transfer(msg.sender, _amount), "failed to transfer");
    }

    function balance() public view returns (uint256) {
        return source_token.balanceOf(address(this));
    }
}
