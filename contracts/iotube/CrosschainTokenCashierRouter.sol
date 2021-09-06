pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";
import "../token/ERC20.sol";

interface ICashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface ICrosschainToken {
    function deposit(uint256 _amount) external;
}

contract CrosschainTokenCashierRouter is Ownable {

    mapping(address => address) private tokenList;
    ICashier public cashier;

    constructor(ICashier _cashier) {
        cashier = _cashier;
    }

    function crosschainTokenOf(address _token) public view returns (address) {
        return tokenList[_token];
    }

    function addCrosschainToken(address _token, address _crosschainToken) public onlyOwner {
        require(ERC20(_token).approve(_crosschainToken, uint256(-1)), "failed to approve allowance to crosschain token");
        require(ERC20(_crosschainToken).approve(address(cashier), uint256(-1)), "failed to approve allowance to cashier");
        tokenList[_token] = _crosschainToken;
    }

    function depositTo(address _token, address _to, uint256 _amount) public payable {
        address ct = tokenList[_token];
        revert(ct != address(0), "invalid token");
        require(safeTransferFrom(_token, msg.sender, address(this), _amount), "failed to transfer token");
        ICrosschainToken(ct).deposit(_amount);
        cashier.depositTo(ct, _to, _amount);
    }

    function safeTransferFrom(address _token, address _from, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, _from, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

}