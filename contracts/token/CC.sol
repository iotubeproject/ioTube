pragma solidity <0.6 >=0.4.24;

import "./MintableToken.sol";
import "./StandardToken.sol";
import "../lifecycle/Pausable.sol";

contract CC is StandardToken, MintableToken, Pausable {
    event MinterAdded(address indexed minter);
    event MinterRemoved(address indexed minter);

    modifier onlyMinter() {
        require(minter == msg.sender, "not the minter");
        _;
    }

    ERC20 public coToken;
    address public minter;
    string public name;
    string public symbol;
    uint8 public decimals;

    constructor(ERC20 _coToken, address _minter, string memory _name, string memory _symbol, uint8 _decimals) public {
        coToken = _coToken;
        minter = _minter;
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        emit MinterAdded(_minter);
    }

    function deposit(uint256 _amount) public whenNotPaused {
        require(coToken.transferFrom(msg.sender, address(this), _amount), "replace with safeTransferFrom");
        require(mintInternal(msg.sender, _amount), "failed to mint");
    }

    function withdraw(uint256 _amount) public whenNotPaused {
        require(burnInternal(msg.sender, _amount), "failed to burn");
        require(coToken.transfer(msg.sender, _amount), "replace with safeTransfer");
    }

    function mint(address _to, uint256 _amount) public onlyMinter whenNotPaused returns (bool) {
        return mintInternal(_to, _amount);
    }

    function mintInternal(address _to, uint256 _amount) internal returns (bool) {
        totalSupply_ = totalSupply_.add(_amount);
        balances[_to] = balances[_to].add(_amount);
        emit Transfer(address(0), _to, _amount);
        return true;
    }

    // user can also burn by sending token to address(0), but this function will emit Burned event
    function burn(uint256 _amount) public returns (bool) {
        return burnInternal(msg.sender, _amount);
    }

    function burnInternal(address _from, uint256 _amount) internal returns (bool) {
        require(balances[_from] >= _amount);
        totalSupply_ = totalSupply_.sub(_amount);
        balances[_from] = balances[_from].sub(_amount);
        emit Transfer(_from, address(0), _amount);
        return true;
    }
}