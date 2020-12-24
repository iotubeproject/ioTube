pragma solidity <6.0 >=0.4.24;

interface Allowlist {
    function isAllowed(address) external view returns (bool);
    function numOfActive() external view returns (uint256);
}

interface TransferValidator {
    function whitelistedTokens() external view returns (Allowlist);
    function getStatus(address tokenAddr, uint256 index, address from, address to, uint256 amount) external view returns (uint256 settleHeight_, uint256 numOfWhitelistedWitnesses_, uint256 numOfValidWitnesses_, address[] memory witnesses_, bool includingMsgSender_);
    function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) external;
}

contract ComboTransferValidatorBase {
    TransferValidator public standardTransferValidator;
    TransferValidator public nonStandardTransferValidator;
    constructor(address _standardTokenValidator, address _nonStandardTokenValidator) public {
        standardTransferValidator = TransferValidator(_standardTokenValidator);
        nonStandardTransferValidator = TransferValidator(_nonStandardTokenValidator);
    }

    function generateKey(address _tokenAddr, uint256 _index, address _from, address _to, uint256 _amount) public pure returns(bytes32) {
        return keccak256(abi.encodePacked(_tokenAddr, _index, _from, _to, _amount));
    }

    function getStatus(address _tokenAddr, uint256 _index, address _from, address _to, uint256 _amount) public view returns (uint256 settleHeight_, uint256 numOfWhitelistedWitnesses_, uint256 numOfValidWitnesses_, address[] memory witnesses_, bool includingMsgSender_) {
        if (standardTransferValidator.whitelistedTokens().isAllowed(_tokenAddr)) {
            return standardTransferValidator.getStatus(_tokenAddr, _index, _from, _to, _amount);
        }
        return nonStandardTransferValidator.getStatus(_tokenAddr, _index, _from, _to, _amount);
    }

    function submit(address _tokenAddr, uint256 _index, address _from, address _to, uint256 _amount) public {
        if (standardTransferValidator.whitelistedTokens().isAllowed(_tokenAddr)) {
            standardTransferValidator.submit(_tokenAddr, _index, _from, _to, _amount);
        } else {
            nonStandardTransferValidator.submit(_tokenAddr, _index, _from, _to, _amount);
        }
    }
}
