// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract TokenPairs is Ownable {
    struct TokenPair {
        uint256 chainID1;
        address token1;
        uint256 chainID2;
        address token2;
        bool active;
    }

    uint256 public lastUpdatedHeight;
    uint256 immutable public chainID1;
    uint256 immutable public chainID2;
    TokenPair[] private _tokenPairs;
    mapping(address => bool) private _tokenExists;
    mapping(address => uint256) private _tokenToPairIndex;

    constructor(uint256 _chainID1, uint256 _chainID2) Ownable(msg.sender) {
        chainID1 = _chainID1;
        chainID2 = _chainID2;
    }

    function addTokenPair(
        uint256 _chainID1,
        address _token1,
        uint256 _chainID2,
        address _token2
    ) external onlyOwner {
        require(_chainID1 == chainID1, "Invalid chainID1");
        require(_chainID2 == chainID2, "Invalid chainID2");
        require(!_tokenExists[_token1], "Token1 already exists");
        require(!_tokenExists[_token2], "Token2 already exists");
        require(_token1 != address(0), "Invalid token1 address");
        require(_token2 != address(0), "Invalid token2 address");

        _tokenPairs.push(TokenPair({
            chainID1: _chainID1,
            token1: _token1,
            chainID2: _chainID2,
            token2: _token2,
            active: true
        }));
        uint256 pairIndex = _tokenPairs.length - 1;
        _tokenExists[_token1] = true;
        _tokenExists[_token2] = true;
        _tokenToPairIndex[_token1] = pairIndex;
        _tokenToPairIndex[_token2] = pairIndex;
        lastUpdatedHeight = block.number;
    }

    function updateTokenPair(
        uint256 _chainID1,
        address _token1,
        uint256 _chainID2,
        address _token2
    ) external onlyOwner {
        require(_chainID1 == chainID1, "Invalid chainID1");
        require(_chainID2 == chainID2, "Invalid chainID2");
        require(_token1 != address(0), "Invalid token1 address");
        require(_token2 != address(0), "Invalid token2 address");

        bool token1Exists = _tokenExists[_token1];
        bool token2Exists = _tokenExists[_token2];
        require(token1Exists != token2Exists, "One token must exist and the other must be new");

        address existingToken = token1Exists ? _token1 : _token2;
        address newToken = token1Exists ? _token2 : _token1;

        uint256 pairIndex = _tokenToPairIndex[existingToken];
        TokenPair storage pair = _tokenPairs[pairIndex];

        address oldToken = (pair.token1 == existingToken) ? pair.token2 : pair.token1;

        delete _tokenExists[oldToken];
        delete _tokenToPairIndex[oldToken];

        pair.chainID1 = _chainID1;
        pair.token1 = _token1;
        pair.chainID2 = _chainID2;
        pair.token2 = _token2;

        _tokenExists[newToken] = true;
        _tokenToPairIndex[newToken] = pairIndex;

        lastUpdatedHeight = block.number;
    }

    function activateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) external onlyOwner {
        require(_chainID1 == chainID1, "Invalid chainID1");
        require(_chainID2 == chainID2, "Invalid chainID2");
        require(_tokenExists[_token1], "Token not found");

        uint256 pairIndex = _tokenToPairIndex[_token1];
        TokenPair storage pair = _tokenPairs[pairIndex];

        require(pair.token2 == _token2, "Token2 mismatch");
        require(!pair.active, "Pair is already active");

        pair.active = true;
        lastUpdatedHeight = block.number;
    }

    function deactivateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) external onlyOwner {
        require(_chainID1 == chainID1, "Invalid chainID1");
        require(_chainID2 == chainID2, "Invalid chainID2");
        require(_tokenExists[_token1], "Token not found");

        uint256 pairIndex = _tokenToPairIndex[_token1];
        TokenPair storage pair = _tokenPairs[pairIndex];

        require(pair.token2 == _token2, "Token2 mismatch");
        require(pair.active, "Pair is already inactive");

        pair.active = false;
        lastUpdatedHeight = block.number;
    }

    function getTokenPairs(uint256 chainID) external view returns (address[] memory, address[] memory) {
        if (chainID != chainID1 && chainID != chainID2) {
            return (new address[](0), new address[](0));
        }

        uint256 activePairsCount = 0;
        for (uint256 i = 0; i < _tokenPairs.length; i++) {
            if (_tokenPairs[i].active) {
                activePairsCount++;
            }
        }

        address[] memory localTokens = new address[](activePairsCount);
        address[] memory remoteTokens = new address[](activePairsCount);
        uint256 counter = 0;

        if (chainID == chainID1) {
            for (uint256 i = 0; i < _tokenPairs.length; i++) {
                if (_tokenPairs[i].active) {
                    localTokens[counter] = _tokenPairs[i].token1;
                    remoteTokens[counter] = _tokenPairs[i].token2;
                    counter++;
                }
            }
        } else { // chainID == chainID2
            for (uint256 i = 0; i < _tokenPairs.length; i++) {
                if (_tokenPairs[i].active) {
                    localTokens[counter] = _tokenPairs[i].token2;
                    remoteTokens[counter] = _tokenPairs[i].token1;
                    counter++;
                }
            }
        }

        return (localTokens, remoteTokens);
    }
}