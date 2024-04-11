pragma solidity <6.0 >=0.4.24;

import "../lifecycle/Pausable.sol";

interface ITokenList {
    function isAllowed(address) external returns (bool);
    function maxAmount(address) external returns (uint256);
    function minAmount(address) external returns (uint256);
}

interface IWrappedCoin {
    function deposit() external payable;
}

contract TokenCashierV3 is Pausable {
    event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload);
    // event ContractDestinationAdded(address indexed destination);
    // event ContractDestinationRemoved(address indexed destination);

    ITokenList[] public tokenLists;
    address[] public tokenSafes;
    mapping(address => uint256) public counts;
    uint256 public depositFee;
    IWrappedCoin public wrappedCoin;
    // mapping(address => bool) public contractDestinations;

    constructor(IWrappedCoin _wrappedCoin, ITokenList[] memory _tokenLists, address[] memory _tokenSafes) public {
        require(_tokenLists.length == _tokenSafes.length, "# of token lists is not equal to # of safes");
        wrappedCoin = _wrappedCoin;
        tokenLists = _tokenLists;
        tokenSafes = _tokenSafes;
    }

    function() external {
        revert();
    }

    function count(address _token) public view returns (uint256) {
        return counts[_token];
    }
/*
    function addContractDestination(address _dest) public onlyOwner {
        require(!contractDestinations[_dest], "already added");
        contractDestinations[_dest] = true;
        emit ContractDestinationAdded(_dest);
    }

    function removeContractDestination(address _dest) public onlyOwner {
        require(contractDestinations[_dest], "invalid destination");
        contractDestinations[_dest] = false;
        emit ContractDestinationRemoved(_dest);
    }
*/
    function setDepositFee(uint256 _fee) public onlyOwner {
        depositFee = _fee;
    }

    function depositTo(address _token, address _to, uint256 _amount, bytes memory _payload) public whenNotPaused payable {
        require(_to != address(0), "invalid destination");
        // require(_payload.length == 0 || contractDestinations[_to], "invalid destination with payload");
        bool isCoin = false;
        uint256 fee = msg.value;
        if (_token == address(0)) {
            require(msg.value >= _amount, "insufficient msg.value");
            fee = msg.value - _amount;
            wrappedCoin.deposit.value(_amount)();
            _token = address(wrappedCoin);
            isCoin = true;
        }
        require(fee >= depositFee, "insufficient fee");
        for (uint256 i = 0; i < tokenLists.length; i++) {
            if (tokenLists[i].isAllowed(_token)) {
                require(_amount >= tokenLists[i].minAmount(_token), "amount too low");
                require(_amount <= tokenLists[i].maxAmount(_token), "amount too high");
                if (tokenSafes[i] == address(0)) {
                    require(!isCoin && safeTransferFrom(_token, msg.sender, address(this), _amount), "fail to transfer token to cashier");
                    // selector = bytes4(keccak256(bytes('burn(uint256)')))
                    (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x42966c68, _amount));
                    require(success && (data.length == 0 || abi.decode(data, (bool))), "fail to burn token");
                } else {
                    if (isCoin) {
                        require(safeTransfer(_token, tokenSafes[i], _amount), "failed to put into safe");
                    } else {
                        require(safeTransferFrom(_token, msg.sender, tokenSafes[i], _amount), "failed to put into safe");
                    }
                }
                counts[_token] += 1;
                emit Receipt(_token, counts[_token], msg.sender, _to, _amount, fee, _payload);
                return;
            }
        }
        revert("not a whitelisted token");
    }

    function deposit(address _token, uint256 _amount, bytes memory _payload) public payable {
        depositTo(_token, msg.sender, _amount, _payload);
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