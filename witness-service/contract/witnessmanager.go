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

// WitnessManagerMetaData contains all meta data concerning the WitnessManager contract.
var WitnessManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_witnessList\",\"type\":\"address\",\"internalType\":\"contractIWitnessList\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addExcludedWitnesses\",\"inputs\":[{\"name\":\"_witnesses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addWitnesses\",\"inputs\":[{\"name\":\"_witnesses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochNum\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExcludedWitnesses\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeWitnesses\",\"inputs\":[{\"name\":\"nextEpochNum\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"witnessesToAdd\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"witnessesToRemove\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeExcludedWitnesses\",\"inputs\":[{\"name\":\"_witnesses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeWitnesses\",\"inputs\":[{\"name\":\"_witnesses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEpochInterval\",\"inputs\":[{\"name\":\"_epochInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEpochNum\",\"inputs\":[{\"name\":\"_epochNum\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferWitnessListOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"witnessList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIWitnessList\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WitnessesProposed\",\"inputs\":[{\"name\":\"epochNum\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"witnessesHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// WitnessManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use WitnessManagerMetaData.ABI instead.
var WitnessManagerABI = WitnessManagerMetaData.ABI

// WitnessManager is an auto generated Go binding around an Ethereum contract.
type WitnessManager struct {
	WitnessManagerCaller     // Read-only binding to the contract
	WitnessManagerTransactor // Write-only binding to the contract
	WitnessManagerFilterer   // Log filterer for contract events
}

// WitnessManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type WitnessManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WitnessManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WitnessManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WitnessManagerSession struct {
	Contract     *WitnessManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WitnessManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WitnessManagerCallerSession struct {
	Contract *WitnessManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// WitnessManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WitnessManagerTransactorSession struct {
	Contract     *WitnessManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// WitnessManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type WitnessManagerRaw struct {
	Contract *WitnessManager // Generic contract binding to access the raw methods on
}

// WitnessManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WitnessManagerCallerRaw struct {
	Contract *WitnessManagerCaller // Generic read-only contract binding to access the raw methods on
}

// WitnessManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WitnessManagerTransactorRaw struct {
	Contract *WitnessManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWitnessManager creates a new instance of WitnessManager, bound to a specific deployed contract.
func NewWitnessManager(address common.Address, backend bind.ContractBackend) (*WitnessManager, error) {
	contract, err := bindWitnessManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WitnessManager{WitnessManagerCaller: WitnessManagerCaller{contract: contract}, WitnessManagerTransactor: WitnessManagerTransactor{contract: contract}, WitnessManagerFilterer: WitnessManagerFilterer{contract: contract}}, nil
}

// NewWitnessManagerCaller creates a new read-only instance of WitnessManager, bound to a specific deployed contract.
func NewWitnessManagerCaller(address common.Address, caller bind.ContractCaller) (*WitnessManagerCaller, error) {
	contract, err := bindWitnessManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WitnessManagerCaller{contract: contract}, nil
}

// NewWitnessManagerTransactor creates a new write-only instance of WitnessManager, bound to a specific deployed contract.
func NewWitnessManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*WitnessManagerTransactor, error) {
	contract, err := bindWitnessManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WitnessManagerTransactor{contract: contract}, nil
}

// NewWitnessManagerFilterer creates a new log filterer instance of WitnessManager, bound to a specific deployed contract.
func NewWitnessManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*WitnessManagerFilterer, error) {
	contract, err := bindWitnessManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WitnessManagerFilterer{contract: contract}, nil
}

// bindWitnessManager binds a generic wrapper to an already deployed contract.
func bindWitnessManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WitnessManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WitnessManager *WitnessManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WitnessManager.Contract.WitnessManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WitnessManager *WitnessManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WitnessManager.Contract.WitnessManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WitnessManager *WitnessManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WitnessManager.Contract.WitnessManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WitnessManager *WitnessManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WitnessManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WitnessManager *WitnessManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WitnessManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WitnessManager *WitnessManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WitnessManager.Contract.contract.Transact(opts, method, params...)
}

// EpochInterval is a free data retrieval call binding the contract method 0x09b1ef26.
//
// Solidity: function epochInterval() view returns(uint64)
func (_WitnessManager *WitnessManagerCaller) EpochInterval(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _WitnessManager.contract.Call(opts, &out, "epochInterval")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EpochInterval is a free data retrieval call binding the contract method 0x09b1ef26.
//
// Solidity: function epochInterval() view returns(uint64)
func (_WitnessManager *WitnessManagerSession) EpochInterval() (uint64, error) {
	return _WitnessManager.Contract.EpochInterval(&_WitnessManager.CallOpts)
}

// EpochInterval is a free data retrieval call binding the contract method 0x09b1ef26.
//
// Solidity: function epochInterval() view returns(uint64)
func (_WitnessManager *WitnessManagerCallerSession) EpochInterval() (uint64, error) {
	return _WitnessManager.Contract.EpochInterval(&_WitnessManager.CallOpts)
}

// EpochNum is a free data retrieval call binding the contract method 0x05e3c05b.
//
// Solidity: function epochNum() view returns(uint64)
func (_WitnessManager *WitnessManagerCaller) EpochNum(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _WitnessManager.contract.Call(opts, &out, "epochNum")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EpochNum is a free data retrieval call binding the contract method 0x05e3c05b.
//
// Solidity: function epochNum() view returns(uint64)
func (_WitnessManager *WitnessManagerSession) EpochNum() (uint64, error) {
	return _WitnessManager.Contract.EpochNum(&_WitnessManager.CallOpts)
}

// EpochNum is a free data retrieval call binding the contract method 0x05e3c05b.
//
// Solidity: function epochNum() view returns(uint64)
func (_WitnessManager *WitnessManagerCallerSession) EpochNum() (uint64, error) {
	return _WitnessManager.Contract.EpochNum(&_WitnessManager.CallOpts)
}

// GetExcludedWitnesses is a free data retrieval call binding the contract method 0x7ab5095c.
//
// Solidity: function getExcludedWitnesses() view returns(address[])
func (_WitnessManager *WitnessManagerCaller) GetExcludedWitnesses(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _WitnessManager.contract.Call(opts, &out, "getExcludedWitnesses")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetExcludedWitnesses is a free data retrieval call binding the contract method 0x7ab5095c.
//
// Solidity: function getExcludedWitnesses() view returns(address[])
func (_WitnessManager *WitnessManagerSession) GetExcludedWitnesses() ([]common.Address, error) {
	return _WitnessManager.Contract.GetExcludedWitnesses(&_WitnessManager.CallOpts)
}

// GetExcludedWitnesses is a free data retrieval call binding the contract method 0x7ab5095c.
//
// Solidity: function getExcludedWitnesses() view returns(address[])
func (_WitnessManager *WitnessManagerCallerSession) GetExcludedWitnesses() ([]common.Address, error) {
	return _WitnessManager.Contract.GetExcludedWitnesses(&_WitnessManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessManager *WitnessManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WitnessManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessManager *WitnessManagerSession) Owner() (common.Address, error) {
	return _WitnessManager.Contract.Owner(&_WitnessManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessManager *WitnessManagerCallerSession) Owner() (common.Address, error) {
	return _WitnessManager.Contract.Owner(&_WitnessManager.CallOpts)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_WitnessManager *WitnessManagerCaller) WitnessList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WitnessManager.contract.Call(opts, &out, "witnessList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_WitnessManager *WitnessManagerSession) WitnessList() (common.Address, error) {
	return _WitnessManager.Contract.WitnessList(&_WitnessManager.CallOpts)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_WitnessManager *WitnessManagerCallerSession) WitnessList() (common.Address, error) {
	return _WitnessManager.Contract.WitnessList(&_WitnessManager.CallOpts)
}

// AddExcludedWitnesses is a paid mutator transaction binding the contract method 0xd583da9b.
//
// Solidity: function addExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactor) AddExcludedWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "addExcludedWitnesses", _witnesses)
}

// AddExcludedWitnesses is a paid mutator transaction binding the contract method 0xd583da9b.
//
// Solidity: function addExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerSession) AddExcludedWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.AddExcludedWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// AddExcludedWitnesses is a paid mutator transaction binding the contract method 0xd583da9b.
//
// Solidity: function addExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactorSession) AddExcludedWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.AddExcludedWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactor) AddWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "addWitnesses", _witnesses)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerSession) AddWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.AddWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactorSession) AddWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.AddWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// ProposeWitnesses is a paid mutator transaction binding the contract method 0x9c4343b8.
//
// Solidity: function proposeWitnesses(uint64 nextEpochNum, address[] witnessesToAdd, address[] witnessesToRemove, bytes[] signatures) returns()
func (_WitnessManager *WitnessManagerTransactor) ProposeWitnesses(opts *bind.TransactOpts, nextEpochNum uint64, witnessesToAdd []common.Address, witnessesToRemove []common.Address, signatures [][]byte) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "proposeWitnesses", nextEpochNum, witnessesToAdd, witnessesToRemove, signatures)
}

// ProposeWitnesses is a paid mutator transaction binding the contract method 0x9c4343b8.
//
// Solidity: function proposeWitnesses(uint64 nextEpochNum, address[] witnessesToAdd, address[] witnessesToRemove, bytes[] signatures) returns()
func (_WitnessManager *WitnessManagerSession) ProposeWitnesses(nextEpochNum uint64, witnessesToAdd []common.Address, witnessesToRemove []common.Address, signatures [][]byte) (*types.Transaction, error) {
	return _WitnessManager.Contract.ProposeWitnesses(&_WitnessManager.TransactOpts, nextEpochNum, witnessesToAdd, witnessesToRemove, signatures)
}

// ProposeWitnesses is a paid mutator transaction binding the contract method 0x9c4343b8.
//
// Solidity: function proposeWitnesses(uint64 nextEpochNum, address[] witnessesToAdd, address[] witnessesToRemove, bytes[] signatures) returns()
func (_WitnessManager *WitnessManagerTransactorSession) ProposeWitnesses(nextEpochNum uint64, witnessesToAdd []common.Address, witnessesToRemove []common.Address, signatures [][]byte) (*types.Transaction, error) {
	return _WitnessManager.Contract.ProposeWitnesses(&_WitnessManager.TransactOpts, nextEpochNum, witnessesToAdd, witnessesToRemove, signatures)
}

// RemoveExcludedWitnesses is a paid mutator transaction binding the contract method 0xb2dea92d.
//
// Solidity: function removeExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactor) RemoveExcludedWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "removeExcludedWitnesses", _witnesses)
}

// RemoveExcludedWitnesses is a paid mutator transaction binding the contract method 0xb2dea92d.
//
// Solidity: function removeExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerSession) RemoveExcludedWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.RemoveExcludedWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// RemoveExcludedWitnesses is a paid mutator transaction binding the contract method 0xb2dea92d.
//
// Solidity: function removeExcludedWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactorSession) RemoveExcludedWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.RemoveExcludedWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactor) RemoveWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "removeWitnesses", _witnesses)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerSession) RemoveWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.RemoveWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessManager *WitnessManagerTransactorSession) RemoveWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.RemoveWitnesses(&_WitnessManager.TransactOpts, _witnesses)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WitnessManager *WitnessManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WitnessManager *WitnessManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _WitnessManager.Contract.RenounceOwnership(&_WitnessManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WitnessManager *WitnessManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WitnessManager.Contract.RenounceOwnership(&_WitnessManager.TransactOpts)
}

// SetEpochInterval is a paid mutator transaction binding the contract method 0x88275029.
//
// Solidity: function setEpochInterval(uint64 _epochInterval) returns()
func (_WitnessManager *WitnessManagerTransactor) SetEpochInterval(opts *bind.TransactOpts, _epochInterval uint64) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "setEpochInterval", _epochInterval)
}

// SetEpochInterval is a paid mutator transaction binding the contract method 0x88275029.
//
// Solidity: function setEpochInterval(uint64 _epochInterval) returns()
func (_WitnessManager *WitnessManagerSession) SetEpochInterval(_epochInterval uint64) (*types.Transaction, error) {
	return _WitnessManager.Contract.SetEpochInterval(&_WitnessManager.TransactOpts, _epochInterval)
}

// SetEpochInterval is a paid mutator transaction binding the contract method 0x88275029.
//
// Solidity: function setEpochInterval(uint64 _epochInterval) returns()
func (_WitnessManager *WitnessManagerTransactorSession) SetEpochInterval(_epochInterval uint64) (*types.Transaction, error) {
	return _WitnessManager.Contract.SetEpochInterval(&_WitnessManager.TransactOpts, _epochInterval)
}

// SetEpochNum is a paid mutator transaction binding the contract method 0xe31ed4ad.
//
// Solidity: function setEpochNum(uint64 _epochNum) returns()
func (_WitnessManager *WitnessManagerTransactor) SetEpochNum(opts *bind.TransactOpts, _epochNum uint64) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "setEpochNum", _epochNum)
}

// SetEpochNum is a paid mutator transaction binding the contract method 0xe31ed4ad.
//
// Solidity: function setEpochNum(uint64 _epochNum) returns()
func (_WitnessManager *WitnessManagerSession) SetEpochNum(_epochNum uint64) (*types.Transaction, error) {
	return _WitnessManager.Contract.SetEpochNum(&_WitnessManager.TransactOpts, _epochNum)
}

// SetEpochNum is a paid mutator transaction binding the contract method 0xe31ed4ad.
//
// Solidity: function setEpochNum(uint64 _epochNum) returns()
func (_WitnessManager *WitnessManagerTransactorSession) SetEpochNum(_epochNum uint64) (*types.Transaction, error) {
	return _WitnessManager.Contract.SetEpochNum(&_WitnessManager.TransactOpts, _epochNum)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.TransferOwnership(&_WitnessManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.TransferOwnership(&_WitnessManager.TransactOpts, newOwner)
}

// TransferWitnessListOwnership is a paid mutator transaction binding the contract method 0x1424f20f.
//
// Solidity: function transferWitnessListOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerTransactor) TransferWitnessListOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.contract.Transact(opts, "transferWitnessListOwnership", newOwner)
}

// TransferWitnessListOwnership is a paid mutator transaction binding the contract method 0x1424f20f.
//
// Solidity: function transferWitnessListOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerSession) TransferWitnessListOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.TransferWitnessListOwnership(&_WitnessManager.TransactOpts, newOwner)
}

// TransferWitnessListOwnership is a paid mutator transaction binding the contract method 0x1424f20f.
//
// Solidity: function transferWitnessListOwnership(address newOwner) returns()
func (_WitnessManager *WitnessManagerTransactorSession) TransferWitnessListOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessManager.Contract.TransferWitnessListOwnership(&_WitnessManager.TransactOpts, newOwner)
}

// WitnessManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WitnessManager contract.
type WitnessManagerOwnershipTransferredIterator struct {
	Event *WitnessManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WitnessManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WitnessManagerOwnershipTransferred)
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
		it.Event = new(WitnessManagerOwnershipTransferred)
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
func (it *WitnessManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WitnessManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WitnessManagerOwnershipTransferred represents a OwnershipTransferred event raised by the WitnessManager contract.
type WitnessManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WitnessManager *WitnessManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WitnessManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WitnessManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WitnessManagerOwnershipTransferredIterator{contract: _WitnessManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WitnessManager *WitnessManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WitnessManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WitnessManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WitnessManagerOwnershipTransferred)
				if err := _WitnessManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WitnessManager *WitnessManagerFilterer) ParseOwnershipTransferred(log types.Log) (*WitnessManagerOwnershipTransferred, error) {
	event := new(WitnessManagerOwnershipTransferred)
	if err := _WitnessManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WitnessManagerWitnessesProposedIterator is returned from FilterWitnessesProposed and is used to iterate over the raw logs and unpacked data for WitnessesProposed events raised by the WitnessManager contract.
type WitnessManagerWitnessesProposedIterator struct {
	Event *WitnessManagerWitnessesProposed // Event containing the contract specifics and raw log

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
func (it *WitnessManagerWitnessesProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WitnessManagerWitnessesProposed)
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
		it.Event = new(WitnessManagerWitnessesProposed)
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
func (it *WitnessManagerWitnessesProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WitnessManagerWitnessesProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WitnessManagerWitnessesProposed represents a WitnessesProposed event raised by the WitnessManager contract.
type WitnessManagerWitnessesProposed struct {
	EpochNum      uint64
	WitnessesHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterWitnessesProposed is a free log retrieval operation binding the contract event 0xdf59af18f69003d67971a471c3358c0025506dc48b886c69ebf573d3c483338e.
//
// Solidity: event WitnessesProposed(uint64 indexed epochNum, bytes32 witnessesHash)
func (_WitnessManager *WitnessManagerFilterer) FilterWitnessesProposed(opts *bind.FilterOpts, epochNum []uint64) (*WitnessManagerWitnessesProposedIterator, error) {

	var epochNumRule []interface{}
	for _, epochNumItem := range epochNum {
		epochNumRule = append(epochNumRule, epochNumItem)
	}

	logs, sub, err := _WitnessManager.contract.FilterLogs(opts, "WitnessesProposed", epochNumRule)
	if err != nil {
		return nil, err
	}
	return &WitnessManagerWitnessesProposedIterator{contract: _WitnessManager.contract, event: "WitnessesProposed", logs: logs, sub: sub}, nil
}

// WatchWitnessesProposed is a free log subscription operation binding the contract event 0xdf59af18f69003d67971a471c3358c0025506dc48b886c69ebf573d3c483338e.
//
// Solidity: event WitnessesProposed(uint64 indexed epochNum, bytes32 witnessesHash)
func (_WitnessManager *WitnessManagerFilterer) WatchWitnessesProposed(opts *bind.WatchOpts, sink chan<- *WitnessManagerWitnessesProposed, epochNum []uint64) (event.Subscription, error) {

	var epochNumRule []interface{}
	for _, epochNumItem := range epochNum {
		epochNumRule = append(epochNumRule, epochNumItem)
	}

	logs, sub, err := _WitnessManager.contract.WatchLogs(opts, "WitnessesProposed", epochNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WitnessManagerWitnessesProposed)
				if err := _WitnessManager.contract.UnpackLog(event, "WitnessesProposed", log); err != nil {
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

// ParseWitnessesProposed is a log parse operation binding the contract event 0xdf59af18f69003d67971a471c3358c0025506dc48b886c69ebf573d3c483338e.
//
// Solidity: event WitnessesProposed(uint64 indexed epochNum, bytes32 witnessesHash)
func (_WitnessManager *WitnessManagerFilterer) ParseWitnessesProposed(log types.Log) (*WitnessManagerWitnessesProposed, error) {
	event := new(WitnessManagerWitnessesProposed)
	if err := _WitnessManager.contract.UnpackLog(event, "WitnessesProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
