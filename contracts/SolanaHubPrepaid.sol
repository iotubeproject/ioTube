// SPDX-License-Identifier: MIT
pragma solidity >= 0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

interface ICashier {
    function depositFee() external view returns (uint256);
    function depositTo(address, string calldata, uint256, bytes calldata) external payable;
}

interface ICrosschainToken {
    function approve(address, uint256) external returns (bool);
    function transfer(address, uint256) external returns (bool);
}

contract SolanaHubPrepaid is Ownable {
    event OperatorAdded(address indexed operator);
    event OperatorRemoved(address indexed operator);
    address immutable public cashier;
    mapping(address => bool) operators;

    constructor(address _cashier) payable {
        cashier = _cashier;
    }

    receive() external payable {
    }

    function withdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    function addOperator(address _operator) external onlyOwner {
        require(_operator != address(0), "SolanaHubPrepaid: operator is the zero address");
        require(!operators[_operator], "SolanaHubPrepaid: already an operator");
        operators[_operator] = true;
        emit OperatorAdded(_operator);
    }

    function removeOperator(address _operator) external onlyOwner {
        require(operators[_operator], "SolanaHubPrepaid: not an operator");
        delete operators[_operator];
        emit OperatorRemoved(_operator);
    }

    function onReceive(address, ICrosschainToken _token, uint256 _amount, bytes calldata _payload) external {
        require(operators[msg.sender], "SolanaHubPrepaid: not an operator");
        ICashier c = ICashier(cashier);
        uint256 fee = c.depositFee();
        require(fee <= address(this).balance, "SolanaHubPrepaid: insufficient balance");
        (string memory recipient, bytes memory payload) = abi.decode(_payload, (string, bytes));
        require(_token.approve(address(c), _amount), "SolanaHubPrepaid: approve cashier failed");
        c.depositTo{value: fee}(address(_token), recipient, _amount, payload);
    }

}
