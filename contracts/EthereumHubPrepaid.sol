// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

interface ICashier {
    function depositFee() external view returns (uint256);
    function depositTo(address, address, uint256, bytes calldata) external payable;
}

interface ICrosschainToken {
    function approve(address, uint256) external returns (bool);
    function transfer(address, uint256) external returns (bool);
}

contract EthereumHubPrepaid is Ownable {
    event OperatorAdded(address indexed operator);
    event OperatorRemoved(address indexed operator);
    address immutable public cashier;
    mapping(address => bool) operators;

    constructor(address _cashier) Ownable(msg.sender) payable {
        cashier = _cashier;
    }

    receive() external payable {
    }

    function withdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    function addOperator(address _operator) external onlyOwner {
        require(_operator != address(0), "EthereumHubPrepaid: operator is the zero address");
        require(!operators[_operator], "EthereumHubPrepaid: already an operator");
        operators[_operator] = true;
        emit OperatorAdded(_operator);
    }

    function removeOperator(address _operator) external onlyOwner {
        require(operators[_operator], "EthereumHubPrepaid: not an operator");
        delete operators[_operator];
        emit OperatorRemoved(_operator);
    }

    function onReceive(address, ICrosschainToken _token, uint256 _amount, bytes calldata _payload) external {
        require(operators[tx.origin], "EthereumHubPrepaid: not an operator");
        ICashier c = ICashier(cashier);
        uint256 fee = c.depositFee();
        require(fee <= address(this).balance, "EthereumHubPrepaid: insufficient balance");
        (address recipient, bytes memory payload) = abi.decode(_payload, (address, bytes));
        require(_token.approve(address(c), _amount), "EthereumHubPrepaid: approve cashier failed");
        c.depositTo{value: fee}(address(_token), recipient, _amount, payload);
    }

}
