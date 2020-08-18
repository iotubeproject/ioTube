pragma solidity <0.6 >=0.4.24;

interface MintableToken {
    function mint(address, uint) external returns (bool);
    function burn(uint) external returns (bool);
}
