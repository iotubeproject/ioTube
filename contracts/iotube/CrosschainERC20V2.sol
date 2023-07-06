// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract CrosschainERC20V2 is ERC20Burnable {
    using SafeERC20 for ERC20;

    event MinterSet(address indexed minter);

    modifier onlyMinter() {
        require(minter == msg.sender, "not the minter");
        _;
    }

    ERC20 public coToken;
    address public minter;
    uint8 private decimals_;

    constructor(
        ERC20 _coToken,
        address _minter,
        string memory _name,
        string memory _symbol,
        uint8 _decimals
    ) ERC20(_name, _symbol) {
        coToken = _coToken;
        minter = _minter;
        decimals_  = _decimals;
        emit MinterSet(_minter);
    }

    function decimals() public view virtual override returns (uint8) {
        return decimals_;
    }

    function transferMintership(address _newMinter) public onlyMinter {
        minter = _newMinter;
        emit MinterSet(_newMinter);
    }

    function deposit(uint256 _amount) public {
        depositTo(msg.sender, _amount);
    }

    function depositTo(address _to, uint256 _amount) public {
        require(address(coToken) != address(0), "no co-token");
        uint256 originBalance = coToken.balanceOf(address(this));
        coToken.safeTransferFrom(msg.sender, address(this), _amount);
        uint256 newBalance = coToken.balanceOf(address(this));
        require(newBalance > originBalance, "invalid balance");
        _mint(_to, newBalance - originBalance);
    }

    function withdraw(uint256 _amount) public {
        withdrawTo(msg.sender, _amount);
    }

    function withdrawTo(address _to, uint256 _amount) public {
        require(address(coToken) != address(0), "no co-token");
        require(_amount != 0, "amount is 0");
        _burn(msg.sender, _amount);
        coToken.safeTransfer(_to, _amount);
    }

    function mint(address _to, uint256 _amount) public onlyMinter returns (bool) {
        require(_amount != 0, "amount is 0");
        _mint(_to, _amount);
        return true;
    }
}
