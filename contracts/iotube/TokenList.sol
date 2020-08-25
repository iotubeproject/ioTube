pragma solidity <6.0 >=0.4.24;

import "./UniqueAppendOnlyAddressList.sol";
import "../ownership/Ownable.sol";

contract TokenList is Ownable, UniqueAppendOnlyAddressList {
    event TokenAdded(address indexed token, uint256 minAmount, uint256 maxAmount);
    event TokenUpdated(address indexed token, uint256 minAmount, uint256 maxAmount);
    event TokenRemoved(address indexed token);

    struct Setting {
        uint256 minAmount;
        uint256 maxAmount;
    }

    mapping(address => Setting) private settings;

    function isAllowed(address _token) public view returns (bool) {
        return isActive(_token);
    }

    function addToken(address _token, uint256 _min, uint256 _max) public onlyOwner returns (bool success_) {
        if (activateItem(_token)) {
            require(_min > 0 && _max > _min, "invalid parameters");
            settings[_token] = Setting(_min, _max);
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
        if (deactivateItem(_token)) {
            emit TokenRemoved(_token);
            success_ = true;
        }
    }

    function setMinAmount(address _token, uint256 _minAmount) public onlyOwner {
        require(isExist(_token), "token not added");
        require(settings[_token].maxAmount >= _minAmount);
        require(_minAmount > 0);
        settings[_token].minAmount = _minAmount;
    }

    function setMaxAmount(address _token, uint256 _maxAmount) public onlyOwner {
        require(isExist(_token), "token not added");
        require(_maxAmount >= settings[_token].minAmount);
        settings[_token].maxAmount = _maxAmount;
    }

    function minAmount(address _token) public view returns (uint256 minAmount_) {
        if (isExist(_token)) {
            minAmount_ = settings[_token].minAmount;
        }
    }

    function maxAmount(address _token) public view returns (uint256 maxAmount_) {
        if (isExist(_token)) {
            maxAmount_ = settings[_token].maxAmount;
        }
    }

}