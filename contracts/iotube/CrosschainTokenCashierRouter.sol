pragma solidity <6.0 >=0.4.24;

import "../token/ERC20.sol";

interface ICashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface ICrosschainToken {
    function deposit(uint256 _amount) external;
    function coToken() external view returns (ERC20);
}

contract CrosschainTokenCashierRouter {

    ICashier public cashier;

    constructor(ICashier _cashier) public {
        cashier = _cashier;
    }

    function approveCrosschainToken(address _crosschainToken) public {
        ERC20 token = ICrosschainToken(_crosschainToken).coToken();
        require(token.approve(_crosschainToken, uint256(-1)), "failed to approve allowance to crosschain token");
        require(ERC20(_crosschainToken).approve(address(cashier), uint256(-1)), "failed to approve allowance to cashier");
    }

    function depositTo(address _crosschainToken, address _to, uint256 _amount) public payable {
        require(_crosschainToken != address(0), "invalid token");
        ERC20 token = ICrosschainToken(_crosschainToken).coToken();
        require(safeTransferFrom(address(token), msg.sender, address(this), _amount), "failed to transfer token");
        ICrosschainToken(_crosschainToken).deposit(_amount);
        cashier.depositTo(_crosschainToken, _to, _amount);
    }

    function safeTransferFrom(address _token, address _from, address _to, uint256 _amount) internal returns (bool) {
        // selector = bytes4(keccak256(bytes('transferFrom(address,address,uint256)')))
        (bool success, bytes memory data) = _token.call(abi.encodeWithSelector(0x23b872dd, _from, _to, _amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }
}
