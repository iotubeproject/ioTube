pragma solidity <6.0 >=0.4.24;

import "../lifecycle/Pausable.sol";

interface TokenList {
    function isAllowed(address) external returns (bool);
    function maxAmount(address) external returns (uint256);
    function minAmount(address) external returns (uint256);
}

contract TokenCashier is Pausable {
    event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee);

    TokenList[] public tokenLists;
    address[] public tokenSafes;
    mapping(address => uint256) public counts;
    uint256 public depositFee;

    constructor(TokenList[] memory _tokenLists, address[] memory _tokenSafes) public {
        require(_tokenLists.length == _tokenSafes.length, "# of token lists is not equal to # of safes");
        tokenLists = _tokenLists;
        tokenSafes = _tokenSafes;
    }

    function() external {
        revert();
    }

    function count(address _token) public view returns (uint256) {
        return counts[_token];
    }

    function setDepositFee(uint256 _fee) public onlyOwner {
        depositFee = _fee;
    }

    function depositTo(address _token, address _to, uint256 _amount) public whenNotPaused payable {
        require(_to != address(0), "invalid destination");
        require(msg.value >= depositFee, "insufficient balance");
        for (uint256 i = 0; i < tokenLists.length; i++) {
            if (tokenLists[i].isAllowed(_token)) {
                require(_amount >= tokenLists[i].minAmount(_token), "amount too low");
                require(_amount <= tokenLists[i].maxAmount(_token), "amount too high");
                if (tokenSafes[i] == address(0)) {
                    require(safeTransferFrom(_token, msg.sender, address(this), _amount), "fail to transfer token to cashier");
                    (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x42966c68, _amount));
                    require(success && (data.length == 0 || abi.decode(data, (bool))), "fail to burn token");
                } else {
                    require(safeTransferFrom(_token, msg.sender, tokenSafes[i], _amount), "failed to put into safe");
                }
                counts[_token] += 1;
                emit Receipt(_token, counts[_token], msg.sender, _to, _amount, msg.value);
                return;
            }
        }
        revert("not a whitelisted token");
    }

    function deposit(address _token, uint256 _amount) public payable {
        depositTo(_token, msg.sender, _amount);
    }

    function withdraw() external onlyOwner {
        msg.sender.transfer(address(this).balance);
    }

    function withdrawToken(address _token) public onlyOwner {
        // selector = bytes4(keccak256(bytes('balanceOf(address)')))
        (bool success, bytes memory balance) = _token.call(abi.encodeWithSelector(0x70a08231, address(this)));
        require(success, "failed to call balanceOf");
        uint256 bal = abi.decode(balance, (uint256));
        if (bal > 0) {
            require(safeTransfer(_token, msg.sender, bal), "failed to withdraw token");
        }
    }

    function safeTransferFrom(address _token, address _from, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, _from, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

    function safeTransfer(address _token, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transfer(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0xa9059cbb, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}