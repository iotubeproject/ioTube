pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

contract TokenList is Ownable {
    event TokenAdded(address indexed token, uint256 minAmount, uint256 maxAmount);
    event TokenUpdated(address indexed token, uint256 minAmount, uint256 maxAmount);
    event TokenRemoved(address indexed token);

    struct Setting {
        uint256 minAmount;
        uint256 maxAmount;
        bool active;
        bool flag;
    }

    mapping(address => Setting) public settings;
    address[] public tokens;
    uint256 private num;

    function isAllowed(address _token) public view returns (bool) {
        return settings[_token].active;
    }
    
    function addToken(address _token, uint256 _min, uint256 _max) public onlyOwner returns (bool success_) {
        if (!settings[_token].active) {
            if (!settings[_token].flag) {
                tokens.push(_token);
            }
            require(_min > 0 && _max > _min, "invalid parameters");
            settings[_token] = Setting(_min, _max, true, true);
            num++;
            emit TokenAdded(_token, _min, _max);
            success_ = true;
        }
    }

    function addTokens(address[] memory _tokens, uint256[] memory _mins, uint256[] memory _maxs) public onlyOwner returns (bool success_) {
        require(_tokens.length == _mins.length && _mins.length == _maxs.length, "invalid parameters");
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (addToken(_tokens[i], _mins[i], _maxs[i])) {
                success_ = true;
            }
        }
    }

    function removeToken(address _token) public onlyOwner returns (bool success_) {
        if (settings[_token].active) {
            settings[_token].active = false;
            num--;
            emit TokenRemoved(_token);
            success_ = true;
        }
    }

    function setMinAmount(address _token, uint256 _minAmount) public onlyOwner {
        require(settings[_token].flag, "token not added");
        require(settings[_token].maxAmount >= _minAmount);
        require(_minAmount > 0);
        settings[_token].minAmount = _minAmount;
    }

    function setMaxAmount(address _token, uint256 _maxAmount) public onlyOwner {
        require(settings[_token].flag, "token not added");
        require(_maxAmount >= settings[_token].minAmount);
        settings[_token].maxAmount = _maxAmount;
    }

    function minAmount(address _token) public view returns (uint256 minAmount_) {
        if (settings[_token].flag) {
            minAmount_ = settings[_token].minAmount;
        }
    }

    function maxAmount(address _token) public view returns (uint256 maxAmount_) {
        if (settings[_token].flag) {
            maxAmount_ = settings[_token].maxAmount;
        }
    }

    function numOfAllowed() public view returns (uint256) {
        return num;
    }
}