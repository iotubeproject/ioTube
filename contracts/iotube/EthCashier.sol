pragma solidity <6.0 >=0.4.24;

import "../ownership/Ownable.sol";

interface TokenCashier {
    function depositTo(address _token, address _to, uint256 _amount) external payable;
}

interface WrappedEther {
    function approve(address _user, uint _amount) external returns (bool);
    function deposit() external payable;
}

contract EthCashier is Ownable {
    TokenCashier public tokenCashier;
    WrappedEther public wrappedEther;

    constructor(address _tokenCashier, address _wrappedEther) public {
        tokenCashier = TokenCashier(_tokenCashier);
        wrappedEther = WrappedEther(_wrappedEther);
    }

    function deposit() public payable {
        depositTo(msg.sender);
    }

    function depositTo(address _to) public payable {
        wrappedEther.deposit.value(msg.value)();
        require(wrappedEther.approve(address(tokenCashier), msg.value), "approve failure");
        tokenCashier.depositTo(address(wrappedEther), _to, msg.value);
    }

    function upgrade(address _newTokenCashier) public onlyOwner {
        tokenCashier = TokenCashier(_newTokenCashier);
    }
}
