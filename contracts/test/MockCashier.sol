pragma solidity = 0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MockCashier {
    event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload);

    uint256 public depositFee;
    
    constructor() {
    }

    receive() external payable {
    }

    function setDepositFee(uint256 _depositFee) external {
        depositFee = _depositFee;
    }

    function depositTo(address _token, address _recipient, uint256 _amount, bytes calldata _payload) external payable {
        require(msg.value >= depositFee, "MockCashier: invalid deposit fee");
        require(IERC20(_token).transferFrom(msg.sender, address(this), _amount), "MockCashier: transferFrom failed");
        emit Receipt(_token, 0, msg.sender, _recipient, _amount, msg.value, _payload);
    }

}
