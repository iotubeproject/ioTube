pragma solidity >=0.8.0;

interface ICashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface IERC20 {
    function balanceOf(address _owner) external view returns(uint256);
}

interface ICrosschainToken {
    function balanceOf(address _account) external returns (uint256);
    function deposit(uint256 _amount) external;
    function coToken() external view returns (address);
}

interface WToken {
    function deposit() external payable;
}

contract CrosschainTokenCashierRouterV2 {

    address public ciotx;
    WToken public wiotx;

    constructor(ICrosschainToken _ciotx) {
        wiotx = WToken(_ciotx.coToken());
        ciotx = address(_ciotx);
        _approve(ciotx, address(wiotx));
    }

    function _approve(address _ctoken, address _coToken) private {
        require(safeApprove(_coToken, _ctoken, type(uint256).max), "failed to approve allowance to crosschain token");
    }

    function approveCrosschainToken(address _crosschainToken) public {
        _approve(_crosschainToken, ICrosschainToken(_crosschainToken).coToken());
    }

    function _depositToCashier(address _token, address _cashier, address _to, uint256 _amount, uint256 _value) private {
        ICrosschainToken ctoken = ICrosschainToken(_token);
        uint256 originBalance = ctoken.balanceOf(address(this));
        ctoken.deposit(_amount);
        uint256 newBalance = ctoken.balanceOf(address(this));
        require(newBalance > originBalance, "invalid balance");
        require(safeApprove(_token, _cashier, newBalance - originBalance), "failed to approve allowance to cashier");
        ICashier(_cashier).depositTo{value: _value}(_token, _to, newBalance - originBalance);
    }

    function depositCoinTo(address _cashier, address _to, uint256 _amount) public payable {
        require(_amount <= msg.value, "invalid amount");
        wiotx.deposit{value: _amount}();
        _depositToCashier(ciotx, _cashier, _to, _amount, msg.value - _amount);
    }

    function depositTo(address _cashier, address _crosschainToken, address _to, uint256 _amount) public payable {
        require(_crosschainToken != address(0), "invalid token");
        address coToken = ICrosschainToken(_crosschainToken).coToken();
        uint256 originBalance = IERC20(coToken).balanceOf(address(this));
        require(safeTransferFrom(coToken, msg.sender, address(this), _amount), "failed to transfer token");
        uint256 newBalance = IERC20(coToken).balanceOf(address(this));
        require(newBalance > originBalance, "invalid balance");
        _depositToCashier(_crosschainToken, _cashier, _to, newBalance - originBalance, msg.value);
    }

    function safeTransferFrom(address _token, address _from, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, _from, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

    function safeApprove(address _token, address _spender, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('approve(address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x095ea7b3, _spender, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}
