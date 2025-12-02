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

// RollDPoSProtocolContractMetaData contains all meta data concerning the RollDPoSProtocolContract contract.
var RollDPoSProtocolContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNumber\",\"type\":\"uint256\"}],\"name\":\"EpochHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNumber\",\"type\":\"uint256\"}],\"name\":\"EpochLastHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"epochLastHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"EpochNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NumCandidateDelegates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"numCandidateDelegates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NumDelegates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"numDelegates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"NumSubEpochs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"numSubEpochs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"SubEpochNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subEpochNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RollDPoSProtocolContractABI is the input ABI used to generate the binding from.
// Deprecated: Use RollDPoSProtocolContractMetaData.ABI instead.
var RollDPoSProtocolContractABI = RollDPoSProtocolContractMetaData.ABI

// RollDPoSProtocolContract is an auto generated Go binding around an Ethereum contract.
type RollDPoSProtocolContract struct {
	RollDPoSProtocolContractCaller     // Read-only binding to the contract
	RollDPoSProtocolContractTransactor // Write-only binding to the contract
	RollDPoSProtocolContractFilterer   // Log filterer for contract events
}

// RollDPoSProtocolContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollDPoSProtocolContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollDPoSProtocolContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollDPoSProtocolContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollDPoSProtocolContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollDPoSProtocolContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollDPoSProtocolContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollDPoSProtocolContractSession struct {
	Contract     *RollDPoSProtocolContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RollDPoSProtocolContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollDPoSProtocolContractCallerSession struct {
	Contract *RollDPoSProtocolContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// RollDPoSProtocolContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollDPoSProtocolContractTransactorSession struct {
	Contract     *RollDPoSProtocolContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// RollDPoSProtocolContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollDPoSProtocolContractRaw struct {
	Contract *RollDPoSProtocolContract // Generic contract binding to access the raw methods on
}

// RollDPoSProtocolContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollDPoSProtocolContractCallerRaw struct {
	Contract *RollDPoSProtocolContractCaller // Generic read-only contract binding to access the raw methods on
}

// RollDPoSProtocolContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollDPoSProtocolContractTransactorRaw struct {
	Contract *RollDPoSProtocolContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollDPoSProtocolContract creates a new instance of RollDPoSProtocolContract, bound to a specific deployed contract.
func NewRollDPoSProtocolContract(address common.Address, backend bind.ContractBackend) (*RollDPoSProtocolContract, error) {
	contract, err := bindRollDPoSProtocolContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollDPoSProtocolContract{RollDPoSProtocolContractCaller: RollDPoSProtocolContractCaller{contract: contract}, RollDPoSProtocolContractTransactor: RollDPoSProtocolContractTransactor{contract: contract}, RollDPoSProtocolContractFilterer: RollDPoSProtocolContractFilterer{contract: contract}}, nil
}

// NewRollDPoSProtocolContractCaller creates a new read-only instance of RollDPoSProtocolContract, bound to a specific deployed contract.
func NewRollDPoSProtocolContractCaller(address common.Address, caller bind.ContractCaller) (*RollDPoSProtocolContractCaller, error) {
	contract, err := bindRollDPoSProtocolContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollDPoSProtocolContractCaller{contract: contract}, nil
}

// NewRollDPoSProtocolContractTransactor creates a new write-only instance of RollDPoSProtocolContract, bound to a specific deployed contract.
func NewRollDPoSProtocolContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RollDPoSProtocolContractTransactor, error) {
	contract, err := bindRollDPoSProtocolContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollDPoSProtocolContractTransactor{contract: contract}, nil
}

// NewRollDPoSProtocolContractFilterer creates a new log filterer instance of RollDPoSProtocolContract, bound to a specific deployed contract.
func NewRollDPoSProtocolContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RollDPoSProtocolContractFilterer, error) {
	contract, err := bindRollDPoSProtocolContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollDPoSProtocolContractFilterer{contract: contract}, nil
}

// bindRollDPoSProtocolContract binds a generic wrapper to an already deployed contract.
func bindRollDPoSProtocolContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollDPoSProtocolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollDPoSProtocolContract.Contract.RollDPoSProtocolContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollDPoSProtocolContract.Contract.RollDPoSProtocolContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollDPoSProtocolContract.Contract.RollDPoSProtocolContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollDPoSProtocolContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollDPoSProtocolContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollDPoSProtocolContract *RollDPoSProtocolContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollDPoSProtocolContract.Contract.contract.Transact(opts, method, params...)
}

// EpochHeight is a free data retrieval call binding the contract method 0x1c659e71.
//
// Solidity: function EpochHeight(uint256 epochNumber) view returns(uint256 epochHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) EpochHeight(opts *bind.CallOpts, epochNumber *big.Int) (struct {
	EpochHeight *big.Int
	Height      *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "EpochHeight", epochNumber)

	outstruct := new(struct {
		EpochHeight *big.Int
		Height      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EpochHeight = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// EpochHeight is a free data retrieval call binding the contract method 0x1c659e71.
//
// Solidity: function EpochHeight(uint256 epochNumber) view returns(uint256 epochHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) EpochHeight(epochNumber *big.Int) (struct {
	EpochHeight *big.Int
	Height      *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochHeight(&_RollDPoSProtocolContract.CallOpts, epochNumber)
}

// EpochHeight is a free data retrieval call binding the contract method 0x1c659e71.
//
// Solidity: function EpochHeight(uint256 epochNumber) view returns(uint256 epochHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) EpochHeight(epochNumber *big.Int) (struct {
	EpochHeight *big.Int
	Height      *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochHeight(&_RollDPoSProtocolContract.CallOpts, epochNumber)
}

// EpochLastHeight is a free data retrieval call binding the contract method 0xbcabebbd.
//
// Solidity: function EpochLastHeight(uint256 epochNumber) view returns(uint256 epochLastHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) EpochLastHeight(opts *bind.CallOpts, epochNumber *big.Int) (struct {
	EpochLastHeight *big.Int
	Height          *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "EpochLastHeight", epochNumber)

	outstruct := new(struct {
		EpochLastHeight *big.Int
		Height          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EpochLastHeight = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// EpochLastHeight is a free data retrieval call binding the contract method 0xbcabebbd.
//
// Solidity: function EpochLastHeight(uint256 epochNumber) view returns(uint256 epochLastHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) EpochLastHeight(epochNumber *big.Int) (struct {
	EpochLastHeight *big.Int
	Height          *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochLastHeight(&_RollDPoSProtocolContract.CallOpts, epochNumber)
}

// EpochLastHeight is a free data retrieval call binding the contract method 0xbcabebbd.
//
// Solidity: function EpochLastHeight(uint256 epochNumber) view returns(uint256 epochLastHeight, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) EpochLastHeight(epochNumber *big.Int) (struct {
	EpochLastHeight *big.Int
	Height          *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochLastHeight(&_RollDPoSProtocolContract.CallOpts, epochNumber)
}

// EpochNumber is a free data retrieval call binding the contract method 0x2f374006.
//
// Solidity: function EpochNumber(uint256 blockHeight) view returns(uint256 epochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) EpochNumber(opts *bind.CallOpts, blockHeight *big.Int) (struct {
	EpochNumber *big.Int
	Height      *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "EpochNumber", blockHeight)

	outstruct := new(struct {
		EpochNumber *big.Int
		Height      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EpochNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// EpochNumber is a free data retrieval call binding the contract method 0x2f374006.
//
// Solidity: function EpochNumber(uint256 blockHeight) view returns(uint256 epochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) EpochNumber(blockHeight *big.Int) (struct {
	EpochNumber *big.Int
	Height      *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochNumber(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}

// EpochNumber is a free data retrieval call binding the contract method 0x2f374006.
//
// Solidity: function EpochNumber(uint256 blockHeight) view returns(uint256 epochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) EpochNumber(blockHeight *big.Int) (struct {
	EpochNumber *big.Int
	Height      *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.EpochNumber(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}

// NumCandidateDelegates is a free data retrieval call binding the contract method 0xe303b7ef.
//
// Solidity: function NumCandidateDelegates() view returns(uint256 numCandidateDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) NumCandidateDelegates(opts *bind.CallOpts) (struct {
	NumCandidateDelegates *big.Int
	Height                *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "NumCandidateDelegates")

	outstruct := new(struct {
		NumCandidateDelegates *big.Int
		Height                *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NumCandidateDelegates = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// NumCandidateDelegates is a free data retrieval call binding the contract method 0xe303b7ef.
//
// Solidity: function NumCandidateDelegates() view returns(uint256 numCandidateDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) NumCandidateDelegates() (struct {
	NumCandidateDelegates *big.Int
	Height                *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumCandidateDelegates(&_RollDPoSProtocolContract.CallOpts)
}

// NumCandidateDelegates is a free data retrieval call binding the contract method 0xe303b7ef.
//
// Solidity: function NumCandidateDelegates() view returns(uint256 numCandidateDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) NumCandidateDelegates() (struct {
	NumCandidateDelegates *big.Int
	Height                *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumCandidateDelegates(&_RollDPoSProtocolContract.CallOpts)
}

// NumDelegates is a free data retrieval call binding the contract method 0x029ddb3e.
//
// Solidity: function NumDelegates() view returns(uint256 numDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) NumDelegates(opts *bind.CallOpts) (struct {
	NumDelegates *big.Int
	Height       *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "NumDelegates")

	outstruct := new(struct {
		NumDelegates *big.Int
		Height       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NumDelegates = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// NumDelegates is a free data retrieval call binding the contract method 0x029ddb3e.
//
// Solidity: function NumDelegates() view returns(uint256 numDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) NumDelegates() (struct {
	NumDelegates *big.Int
	Height       *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumDelegates(&_RollDPoSProtocolContract.CallOpts)
}

// NumDelegates is a free data retrieval call binding the contract method 0x029ddb3e.
//
// Solidity: function NumDelegates() view returns(uint256 numDelegates, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) NumDelegates() (struct {
	NumDelegates *big.Int
	Height       *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumDelegates(&_RollDPoSProtocolContract.CallOpts)
}

// NumSubEpochs is a free data retrieval call binding the contract method 0xbe0be990.
//
// Solidity: function NumSubEpochs(uint256 blockHeight) view returns(uint256 numSubEpochs, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) NumSubEpochs(opts *bind.CallOpts, blockHeight *big.Int) (struct {
	NumSubEpochs *big.Int
	Height       *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "NumSubEpochs", blockHeight)

	outstruct := new(struct {
		NumSubEpochs *big.Int
		Height       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NumSubEpochs = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// NumSubEpochs is a free data retrieval call binding the contract method 0xbe0be990.
//
// Solidity: function NumSubEpochs(uint256 blockHeight) view returns(uint256 numSubEpochs, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) NumSubEpochs(blockHeight *big.Int) (struct {
	NumSubEpochs *big.Int
	Height       *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumSubEpochs(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}

// NumSubEpochs is a free data retrieval call binding the contract method 0xbe0be990.
//
// Solidity: function NumSubEpochs(uint256 blockHeight) view returns(uint256 numSubEpochs, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) NumSubEpochs(blockHeight *big.Int) (struct {
	NumSubEpochs *big.Int
	Height       *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.NumSubEpochs(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}

// SubEpochNumber is a free data retrieval call binding the contract method 0x0f3a0358.
//
// Solidity: function SubEpochNumber(uint256 blockHeight) view returns(uint256 subEpochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCaller) SubEpochNumber(opts *bind.CallOpts, blockHeight *big.Int) (struct {
	SubEpochNumber *big.Int
	Height         *big.Int
}, error) {
	var out []interface{}
	err := _RollDPoSProtocolContract.contract.Call(opts, &out, "SubEpochNumber", blockHeight)

	outstruct := new(struct {
		SubEpochNumber *big.Int
		Height         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubEpochNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SubEpochNumber is a free data retrieval call binding the contract method 0x0f3a0358.
//
// Solidity: function SubEpochNumber(uint256 blockHeight) view returns(uint256 subEpochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractSession) SubEpochNumber(blockHeight *big.Int) (struct {
	SubEpochNumber *big.Int
	Height         *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.SubEpochNumber(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}

// SubEpochNumber is a free data retrieval call binding the contract method 0x0f3a0358.
//
// Solidity: function SubEpochNumber(uint256 blockHeight) view returns(uint256 subEpochNumber, uint256 height)
func (_RollDPoSProtocolContract *RollDPoSProtocolContractCallerSession) SubEpochNumber(blockHeight *big.Int) (struct {
	SubEpochNumber *big.Int
	Height         *big.Int
}, error) {
	return _RollDPoSProtocolContract.Contract.SubEpochNumber(&_RollDPoSProtocolContract.CallOpts, blockHeight)
}
