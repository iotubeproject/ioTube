pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

interface TokenCashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface WrappedCoin {
    function approve(address _user, uint _amount) external returns (bool);
    function deposit() external payable;
}

contract EthCashier is Ownable {
    TokenCashier public tokenCashier;
    WrappedCoin public wrappedCoin;

    constructor(address _tokenCashier, address _wrappedCoin) public {
        tokenCashier = TokenCashier(_tokenCashier);
        wrappedCoin = WrappedCoin(_wrappedCoin);
    }

    function deposit() public payable {
        depositTo(msg.sender);
    }

    function depositTo(address _to) public payable {
        wrappedCoin.deposit.value(msg.value)();
        require(wrappedCoin.approve(address(tokenCashier), msg.value), "approve failure");
        tokenCashier.depositTo(address(wrappedCoin), _to, msg.value);
    }

    function upgrade(address _newTokenCashier) public onlyOwner {
        tokenCashier = TokenCashier(_newTokenCashier);
    }
}