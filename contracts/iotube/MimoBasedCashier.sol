pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

interface IMimoFactory {
    function getAmountsOut(uint256 amountIn, address[] memory path) external view returns (uint256[] memory amounts);
}

contract MimoBasedCashier is Ownable {
    event Exchanged(address indexed recipient, address indexed srcToken, uint256 amountIn, uint256 amountOut);
    uint256 public maxAmount;
    uint256 public transactionFee;
    uint8 public taxRate;
    IMimoFactory immutable public mimoFactory;
    mapping(address => address[]) private paths;
    mapping(bytes32 => bool) private transfers;
    mapping(address => bool) private operators;
    modifier onlyOperator {
        require(operators[msg.sender], "invalid operator");
        _;
    }

    constructor(address _factory, uint256 _maxAmount, uint256 _fee, uint8 _taxRate) Ownable(msg.sender) {
        maxAmount = _maxAmount;
        transactionFee = _fee;
        taxRate = _taxRate;
        mimoFactory = IMimoFactory(_factory);
    }

    function setTokenPath(address _token, address[] memory _path) public onlyOperator {
        if (_path.length == 0) {
            delete paths[_token];
            return;
        }
        paths[_token] = _path;
    }

    function pathOf(address _token) public view returns (address[] memory) {
        return paths[_token];
    }

    function genId(uint256 _chainId, bytes32 _txHash, address _recipient, address _token, uint256 _amount) public pure returns (bytes32) {
        return keccak256(abi.encode(_chainId, _txHash, _recipient, _token, _amount));
    }

    function deposit(uint256 _chainId, bytes32 _txHash, address payable _recipient, address _token, uint256 _amountIn) public onlyOperator {
        bytes32 id = genId(_chainId, _txHash, _recipient, _token, _amountIn);
        require(!transfers[id], "duplicate transfer");
        address[] memory path = paths[_token];
        require(path.length > 0, "token not supported");
        uint256[] memory outs = mimoFactory.getAmountsOut(_amountIn, path);
        uint256 amountOut = outs[path.length - 1];
        if (amountOut > transactionFee) {
            amountOut -= transactionFee;
            if (taxRate > 0) {
                amountOut = amountOut * (10000 - taxRate) / 10000;
            }
            if (amountOut > 0) {
                _recipient.transfer(amountOut);
            }
        } else {
            amountOut = 0;
        }
        emit Exchanged(_recipient, _token, _amountIn, amountOut);
    }

    function withdraw() public onlyOwner {
        payable(msg.sender).transfer(address(this).balance);
    }
}