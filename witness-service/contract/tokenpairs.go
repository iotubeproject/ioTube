// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TokenPairsMetaData contains all meta data concerning the TokenPairs contract.
var TokenPairsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_chainID1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_chainID2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activateTokenPair\",\"inputs\":[{\"name\":\"_chainID1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token1\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainID2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token2\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addTokenPair\",\"inputs\":[{\"name\":\"_chainID1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token1\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainID2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token2\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"chainID1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chainID2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deactivateTokenPair\",\"inputs\":[{\"name\":\"_chainID1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token1\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainID2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token2\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getTokenPairs\",\"inputs\":[{\"name\":\"chainID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdatedHeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTokenPair\",\"inputs\":[{\"name\":\"_chainID1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token1\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_chainID2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_token2\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// TokenPairsABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenPairsMetaData.ABI instead.
var TokenPairsABI = TokenPairsMetaData.ABI

// TokenPairs is an auto generated Go binding around an Ethereum contract.
type TokenPairs struct {
	TokenPairsCaller     // Read-only binding to the contract
	TokenPairsTransactor // Write-only binding to the contract
	TokenPairsFilterer   // Log filterer for contract events
}

// TokenPairsCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenPairsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenPairsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenPairsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenPairsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenPairsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenPairsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenPairsSession struct {
	Contract     *TokenPairs       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenPairsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenPairsCallerSession struct {
	Contract *TokenPairsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TokenPairsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenPairsTransactorSession struct {
	Contract     *TokenPairsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TokenPairsRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenPairsRaw struct {
	Contract *TokenPairs // Generic contract binding to access the raw methods on
}

// TokenPairsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenPairsCallerRaw struct {
	Contract *TokenPairsCaller // Generic read-only contract binding to access the raw methods on
}

// TokenPairsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenPairsTransactorRaw struct {
	Contract *TokenPairsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenPairs creates a new instance of TokenPairs, bound to a specific deployed contract.
func NewTokenPairs(address common.Address, backend bind.ContractBackend) (*TokenPairs, error) {
	contract, err := bindTokenPairs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenPairs{TokenPairsCaller: TokenPairsCaller{contract: contract}, TokenPairsTransactor: TokenPairsTransactor{contract: contract}, TokenPairsFilterer: TokenPairsFilterer{contract: contract}}, nil
}

// NewTokenPairsCaller creates a new read-only instance of TokenPairs, bound to a specific deployed contract.
func NewTokenPairsCaller(address common.Address, caller bind.ContractCaller) (*TokenPairsCaller, error) {
	contract, err := bindTokenPairs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenPairsCaller{contract: contract}, nil
}

// NewTokenPairsTransactor creates a new write-only instance of TokenPairs, bound to a specific deployed contract.
func NewTokenPairsTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenPairsTransactor, error) {
	contract, err := bindTokenPairs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenPairsTransactor{contract: contract}, nil
}

// NewTokenPairsFilterer creates a new log filterer instance of TokenPairs, bound to a specific deployed contract.
func NewTokenPairsFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenPairsFilterer, error) {
	contract, err := bindTokenPairs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenPairsFilterer{contract: contract}, nil
}

// bindTokenPairs binds a generic wrapper to an already deployed contract.
func bindTokenPairs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenPairsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenPairs *TokenPairsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenPairs.Contract.TokenPairsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenPairs *TokenPairsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPairs.Contract.TokenPairsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenPairs *TokenPairsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenPairs.Contract.TokenPairsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenPairs *TokenPairsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenPairs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenPairs *TokenPairsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPairs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenPairs *TokenPairsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenPairs.Contract.contract.Transact(opts, method, params...)
}

// ChainID1 is a free data retrieval call binding the contract method 0xe7728e45.
//
// Solidity: function chainID1() view returns(uint256)
func (_TokenPairs *TokenPairsCaller) ChainID1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenPairs.contract.Call(opts, &out, "chainID1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChainID1 is a free data retrieval call binding the contract method 0xe7728e45.
//
// Solidity: function chainID1() view returns(uint256)
func (_TokenPairs *TokenPairsSession) ChainID1() (*big.Int, error) {
	return _TokenPairs.Contract.ChainID1(&_TokenPairs.CallOpts)
}

// ChainID1 is a free data retrieval call binding the contract method 0xe7728e45.
//
// Solidity: function chainID1() view returns(uint256)
func (_TokenPairs *TokenPairsCallerSession) ChainID1() (*big.Int, error) {
	return _TokenPairs.Contract.ChainID1(&_TokenPairs.CallOpts)
}

// ChainID2 is a free data retrieval call binding the contract method 0x512a0230.
//
// Solidity: function chainID2() view returns(uint256)
func (_TokenPairs *TokenPairsCaller) ChainID2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenPairs.contract.Call(opts, &out, "chainID2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChainID2 is a free data retrieval call binding the contract method 0x512a0230.
//
// Solidity: function chainID2() view returns(uint256)
func (_TokenPairs *TokenPairsSession) ChainID2() (*big.Int, error) {
	return _TokenPairs.Contract.ChainID2(&_TokenPairs.CallOpts)
}

// ChainID2 is a free data retrieval call binding the contract method 0x512a0230.
//
// Solidity: function chainID2() view returns(uint256)
func (_TokenPairs *TokenPairsCallerSession) ChainID2() (*big.Int, error) {
	return _TokenPairs.Contract.ChainID2(&_TokenPairs.CallOpts)
}

// GetTokenPairs is a free data retrieval call binding the contract method 0xcdbd5a04.
//
// Solidity: function getTokenPairs(uint256 chainID) view returns(address[], address[])
func (_TokenPairs *TokenPairsCaller) GetTokenPairs(opts *bind.CallOpts, chainID *big.Int) ([]common.Address, []common.Address, error) {
	var out []interface{}
	err := _TokenPairs.contract.Call(opts, &out, "getTokenPairs", chainID)

	if err != nil {
		return *new([]common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err

}

// GetTokenPairs is a free data retrieval call binding the contract method 0xcdbd5a04.
//
// Solidity: function getTokenPairs(uint256 chainID) view returns(address[], address[])
func (_TokenPairs *TokenPairsSession) GetTokenPairs(chainID *big.Int) ([]common.Address, []common.Address, error) {
	return _TokenPairs.Contract.GetTokenPairs(&_TokenPairs.CallOpts, chainID)
}

// GetTokenPairs is a free data retrieval call binding the contract method 0xcdbd5a04.
//
// Solidity: function getTokenPairs(uint256 chainID) view returns(address[], address[])
func (_TokenPairs *TokenPairsCallerSession) GetTokenPairs(chainID *big.Int) ([]common.Address, []common.Address, error) {
	return _TokenPairs.Contract.GetTokenPairs(&_TokenPairs.CallOpts, chainID)
}

// LastUpdatedHeight is a free data retrieval call binding the contract method 0xc29fab14.
//
// Solidity: function lastUpdatedHeight() view returns(uint256)
func (_TokenPairs *TokenPairsCaller) LastUpdatedHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenPairs.contract.Call(opts, &out, "lastUpdatedHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdatedHeight is a free data retrieval call binding the contract method 0xc29fab14.
//
// Solidity: function lastUpdatedHeight() view returns(uint256)
func (_TokenPairs *TokenPairsSession) LastUpdatedHeight() (*big.Int, error) {
	return _TokenPairs.Contract.LastUpdatedHeight(&_TokenPairs.CallOpts)
}

// LastUpdatedHeight is a free data retrieval call binding the contract method 0xc29fab14.
//
// Solidity: function lastUpdatedHeight() view returns(uint256)
func (_TokenPairs *TokenPairsCallerSession) LastUpdatedHeight() (*big.Int, error) {
	return _TokenPairs.Contract.LastUpdatedHeight(&_TokenPairs.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenPairs *TokenPairsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenPairs.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenPairs *TokenPairsSession) Owner() (common.Address, error) {
	return _TokenPairs.Contract.Owner(&_TokenPairs.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenPairs *TokenPairsCallerSession) Owner() (common.Address, error) {
	return _TokenPairs.Contract.Owner(&_TokenPairs.CallOpts)
}

// ActivateTokenPair is a paid mutator transaction binding the contract method 0x49183992.
//
// Solidity: function activateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactor) ActivateTokenPair(opts *bind.TransactOpts, _chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "activateTokenPair", _chainID1, _token1, _chainID2, _token2)
}

// ActivateTokenPair is a paid mutator transaction binding the contract method 0x49183992.
//
// Solidity: function activateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsSession) ActivateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.ActivateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// ActivateTokenPair is a paid mutator transaction binding the contract method 0x49183992.
//
// Solidity: function activateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactorSession) ActivateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.ActivateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// AddTokenPair is a paid mutator transaction binding the contract method 0xdb1f7f21.
//
// Solidity: function addTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactor) AddTokenPair(opts *bind.TransactOpts, _chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "addTokenPair", _chainID1, _token1, _chainID2, _token2)
}

// AddTokenPair is a paid mutator transaction binding the contract method 0xdb1f7f21.
//
// Solidity: function addTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsSession) AddTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.AddTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// AddTokenPair is a paid mutator transaction binding the contract method 0xdb1f7f21.
//
// Solidity: function addTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactorSession) AddTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.AddTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// DeactivateTokenPair is a paid mutator transaction binding the contract method 0x1cf937be.
//
// Solidity: function deactivateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactor) DeactivateTokenPair(opts *bind.TransactOpts, _chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "deactivateTokenPair", _chainID1, _token1, _chainID2, _token2)
}

// DeactivateTokenPair is a paid mutator transaction binding the contract method 0x1cf937be.
//
// Solidity: function deactivateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsSession) DeactivateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.DeactivateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// DeactivateTokenPair is a paid mutator transaction binding the contract method 0x1cf937be.
//
// Solidity: function deactivateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactorSession) DeactivateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.DeactivateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenPairs *TokenPairsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenPairs *TokenPairsSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenPairs.Contract.RenounceOwnership(&_TokenPairs.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenPairs *TokenPairsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenPairs.Contract.RenounceOwnership(&_TokenPairs.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenPairs *TokenPairsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenPairs *TokenPairsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.TransferOwnership(&_TokenPairs.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenPairs *TokenPairsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.TransferOwnership(&_TokenPairs.TransactOpts, newOwner)
}

// UpdateTokenPair is a paid mutator transaction binding the contract method 0x5eb30fe9.
//
// Solidity: function updateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactor) UpdateTokenPair(opts *bind.TransactOpts, _chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.contract.Transact(opts, "updateTokenPair", _chainID1, _token1, _chainID2, _token2)
}

// UpdateTokenPair is a paid mutator transaction binding the contract method 0x5eb30fe9.
//
// Solidity: function updateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsSession) UpdateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.UpdateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// UpdateTokenPair is a paid mutator transaction binding the contract method 0x5eb30fe9.
//
// Solidity: function updateTokenPair(uint256 _chainID1, address _token1, uint256 _chainID2, address _token2) returns()
func (_TokenPairs *TokenPairsTransactorSession) UpdateTokenPair(_chainID1 *big.Int, _token1 common.Address, _chainID2 *big.Int, _token2 common.Address) (*types.Transaction, error) {
	return _TokenPairs.Contract.UpdateTokenPair(&_TokenPairs.TransactOpts, _chainID1, _token1, _chainID2, _token2)
}

// TokenPairsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenPairs contract.
type TokenPairsOwnershipTransferredIterator struct {
	Event *TokenPairsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TokenPairsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPairsOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TokenPairsOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TokenPairsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenPairsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenPairsOwnershipTransferred represents a OwnershipTransferred event raised by the TokenPairs contract.
type TokenPairsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenPairs *TokenPairsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenPairsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenPairs.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenPairsOwnershipTransferredIterator{contract: _TokenPairs.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenPairs *TokenPairsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenPairsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenPairs.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenPairsOwnershipTransferred)
				if err := _TokenPairs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenPairs *TokenPairsFilterer) ParseOwnershipTransferred(log types.Log) (*TokenPairsOwnershipTransferred, error) {
	event := new(TokenPairsOwnershipTransferred)
	if err := _TokenPairs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
