pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract TokenConfigList is Ownable {
    event Inactivated(address indexed token);
    event Activated(address indexed token);
    event Upserted(address indexed token, address indexed targetToken, uint256 minAmount, uint256 maxAmount, bool whitelistOnly);
    event Whitelisted(address indexed token, address indexed account);
    event Unwhitelisted(address indexed token, address indexed account);
    struct Config {
        address sourceToken;
        address targetToken;
        uint256 minAmount;
        uint256 maxAmount;
        bool activated;
        bool whitelistOnly;
    }
    mapping(address => mapping(address => bool)) whitelists;

    mapping(address => uint256) private tokenIndex;
    Config[] private configs;

    function mustIndexByToken(address _token) internal view returns (uint256) {
        uint256 idx = tokenIndex[_token];
        require(idx != 0, "TokenList: token not exists");
        return idx - 1;
    }

    function count() public view returns (uint256) {
        return configs.length;
    }

    function getConfig(uint256 _index) public view returns (Config memory) {
        require(_index < configs.length, "TokenList: index out of bounds");
        return configs[_index];
    }

    function isActivated(address _token) public view returns (bool) {
        return configs[mustIndexByToken(_token)].activated;
    }

    function isWhitelisted(address _token, address _user) public view returns (bool) {
        if (!configs[mustIndexByToken(_token)].whitelistOnly) {
            return true;
        }
        return whitelists[_token][_user];
    }

    function pause(address _token) public onlyOwner {
        configs[mustIndexByToken(_token)].activated = false;
        emit Inactivated(_token);
    }

    function activate(address _token) public onlyOwner {
        configs[mustIndexByToken(_token)].activated = true;
        emit Activated(_token);
    }

    function upsert(address _token, address _targetToken, uint256 _minAmount, uint256 _maxAmount, bool _whitelistOnly) public onlyOwner {
        require(_token != address(0), "TokenList: token is zero address");
        require(_targetToken != address(0), "TokenList: targetToken is zero address");
        require(_minAmount <= _maxAmount, "TokenList: minAmount > maxAmount");
        uint256 idx = tokenIndex[_token];
        if (idx == 0) {
            configs.push(Config(_token, _targetToken, _minAmount, _maxAmount, false, _whitelistOnly));
            tokenIndex[_token] = configs.length;
        } else {
            Config storage cfg = configs[idx - 1];
            cfg.targetToken = _targetToken;
            cfg.minAmount = _minAmount;
            cfg.maxAmount = _maxAmount;
            cfg.whitelistOnly = _whitelistOnly;
        }
        emit Upserted(_token, _targetToken, _minAmount, _maxAmount, _whitelistOnly);
    }

    function whitelist(address _token, address _user) public onlyOwner {
        require(_user != address(0), "TokenList: user is zero address");
        require(configs[mustIndexByToken(_token)].whitelistOnly, "TokenList: not whitelist only");
        require(!whitelists[_token][_user], "TokenList: user already whitelisted");
        whitelists[_token][_user] = true;
        emit Whitelisted(_token, _user);
    }

    function unwhitelist(address _token, address _user) public onlyOwner {
        require(_user != address(0), "TokenList: user is zero address");
        require(configs[mustIndexByToken(_token)].whitelistOnly, "TokenList: not whitelist only");
        require(whitelists[_token][_user], "TokenList: user not whitelisted");
        whitelists[_token][_user] = false;
        emit Unwhitelisted(_token, _user);
    }

}

contract CashierConfig is Ownable {

    string public id;
    bool immutable public withPayload;
    uint256 immutable public startBlockHeight;
    address immutable public cashier;
    address immutable public validator;
    address immutable public tokenSafe;
    address immutable public tokenList;

    constructor(
        string memory _id,
        bool _withPayload,
        uint256 _startBlockHeight,
        address _cashier,
        address _validator,
        address _tokenSafe,
        address _tokenList
    ) {
        id = _id;
        withPayload = _withPayload;
        startBlockHeight = _startBlockHeight;
        cashier = _cashier;
        validator = _validator;
        tokenSafe = _tokenSafe;
        tokenList = _tokenList;
    }
}