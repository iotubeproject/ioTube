pragma solidity <6.0 >=0.4.21;

interface Safe {
    event Withdrew(address to, uint256 amount);
    event Deposited(address from, uint256 amount);
    function withdraw(address _to, uint256 _amount) external returns (bool);
}