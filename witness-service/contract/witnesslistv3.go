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

// WitnessListV3MetaData contains all meta data concerning the WitnessListV3 contract.
var WitnessListV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"witnessUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_witness\",\"type\":\"address\"}],\"name\":\"addWitness\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success_\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_witnesses\",\"type\":\"address[]\"}],\"name\":\"addWitnesses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_witnesses\",\"type\":\"address[]\"}],\"name\":\"areAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"limit\",\"type\":\"uint8\"}],\"name\":\"getActiveItems\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"items_\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_item\",\"type\":\"address\"}],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_witness\",\"type\":\"address\"}],\"name\":\"isAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_item\",\"type\":\"address\"}],\"name\":\"isExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numOfActive\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_witness\",\"type\":\"address\"}],\"name\":\"removeWitness\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success_\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_witnesses\",\"type\":\"address[]\"}],\"name\":\"removeWitnesses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"}],\"name\":\"setThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newWitness\",\"type\":\"address\"}],\"name\":\"switchWitness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// WitnessListV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use WitnessListV3MetaData.ABI instead.
var WitnessListV3ABI = WitnessListV3MetaData.ABI

// WitnessListV3 is an auto generated Go binding around an Ethereum contract.
type WitnessListV3 struct {
	WitnessListV3Caller     // Read-only binding to the contract
	WitnessListV3Transactor // Write-only binding to the contract
	WitnessListV3Filterer   // Log filterer for contract events
}

// WitnessListV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type WitnessListV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessListV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type WitnessListV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessListV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WitnessListV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WitnessListV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WitnessListV3Session struct {
	Contract     *WitnessListV3    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WitnessListV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WitnessListV3CallerSession struct {
	Contract *WitnessListV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// WitnessListV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WitnessListV3TransactorSession struct {
	Contract     *WitnessListV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// WitnessListV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type WitnessListV3Raw struct {
	Contract *WitnessListV3 // Generic contract binding to access the raw methods on
}

// WitnessListV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WitnessListV3CallerRaw struct {
	Contract *WitnessListV3Caller // Generic read-only contract binding to access the raw methods on
}

// WitnessListV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WitnessListV3TransactorRaw struct {
	Contract *WitnessListV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewWitnessListV3 creates a new instance of WitnessListV3, bound to a specific deployed contract.
func NewWitnessListV3(address common.Address, backend bind.ContractBackend) (*WitnessListV3, error) {
	contract, err := bindWitnessListV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WitnessListV3{WitnessListV3Caller: WitnessListV3Caller{contract: contract}, WitnessListV3Transactor: WitnessListV3Transactor{contract: contract}, WitnessListV3Filterer: WitnessListV3Filterer{contract: contract}}, nil
}

// NewWitnessListV3Caller creates a new read-only instance of WitnessListV3, bound to a specific deployed contract.
func NewWitnessListV3Caller(address common.Address, caller bind.ContractCaller) (*WitnessListV3Caller, error) {
	contract, err := bindWitnessListV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WitnessListV3Caller{contract: contract}, nil
}

// NewWitnessListV3Transactor creates a new write-only instance of WitnessListV3, bound to a specific deployed contract.
func NewWitnessListV3Transactor(address common.Address, transactor bind.ContractTransactor) (*WitnessListV3Transactor, error) {
	contract, err := bindWitnessListV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WitnessListV3Transactor{contract: contract}, nil
}

// NewWitnessListV3Filterer creates a new log filterer instance of WitnessListV3, bound to a specific deployed contract.
func NewWitnessListV3Filterer(address common.Address, filterer bind.ContractFilterer) (*WitnessListV3Filterer, error) {
	contract, err := bindWitnessListV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WitnessListV3Filterer{contract: contract}, nil
}

// bindWitnessListV3 binds a generic wrapper to an already deployed contract.
func bindWitnessListV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WitnessListV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WitnessListV3 *WitnessListV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WitnessListV3.Contract.WitnessListV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WitnessListV3 *WitnessListV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WitnessListV3.Contract.WitnessListV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WitnessListV3 *WitnessListV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WitnessListV3.Contract.WitnessListV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WitnessListV3 *WitnessListV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WitnessListV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WitnessListV3 *WitnessListV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WitnessListV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WitnessListV3 *WitnessListV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WitnessListV3.Contract.contract.Transact(opts, method, params...)
}

// AreAllowed is a free data retrieval call binding the contract method 0x4b6fa7c5.
//
// Solidity: function areAllowed(address[] _witnesses) view returns(bool)
func (_WitnessListV3 *WitnessListV3Caller) AreAllowed(opts *bind.CallOpts, _witnesses []common.Address) (bool, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "areAllowed", _witnesses)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AreAllowed is a free data retrieval call binding the contract method 0x4b6fa7c5.
//
// Solidity: function areAllowed(address[] _witnesses) view returns(bool)
func (_WitnessListV3 *WitnessListV3Session) AreAllowed(_witnesses []common.Address) (bool, error) {
	return _WitnessListV3.Contract.AreAllowed(&_WitnessListV3.CallOpts, _witnesses)
}

// AreAllowed is a free data retrieval call binding the contract method 0x4b6fa7c5.
//
// Solidity: function areAllowed(address[] _witnesses) view returns(bool)
func (_WitnessListV3 *WitnessListV3CallerSession) AreAllowed(_witnesses []common.Address) (bool, error) {
	return _WitnessListV3.Contract.AreAllowed(&_WitnessListV3.CallOpts, _witnesses)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_WitnessListV3 *WitnessListV3Caller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_WitnessListV3 *WitnessListV3Session) Count() (*big.Int, error) {
	return _WitnessListV3.Contract.Count(&_WitnessListV3.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_WitnessListV3 *WitnessListV3CallerSession) Count() (*big.Int, error) {
	return _WitnessListV3.Contract.Count(&_WitnessListV3.CallOpts)
}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_WitnessListV3 *WitnessListV3Caller) GetActiveItems(opts *bind.CallOpts, offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "getActiveItems", offset, limit)

	outstruct := new(struct {
		Count *big.Int
		Items []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Count = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Items = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_WitnessListV3 *WitnessListV3Session) GetActiveItems(offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	return _WitnessListV3.Contract.GetActiveItems(&_WitnessListV3.CallOpts, offset, limit)
}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_WitnessListV3 *WitnessListV3CallerSession) GetActiveItems(offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	return _WitnessListV3.Contract.GetActiveItems(&_WitnessListV3.CallOpts, offset, limit)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3Caller) IsActive(opts *bind.CallOpts, _item common.Address) (bool, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "isActive", _item)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3Session) IsActive(_item common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsActive(&_WitnessListV3.CallOpts, _item)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3CallerSession) IsActive(_item common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsActive(&_WitnessListV3.CallOpts, _item)
}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _witness) view returns(bool)
func (_WitnessListV3 *WitnessListV3Caller) IsAllowed(opts *bind.CallOpts, _witness common.Address) (bool, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "isAllowed", _witness)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _witness) view returns(bool)
func (_WitnessListV3 *WitnessListV3Session) IsAllowed(_witness common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsAllowed(&_WitnessListV3.CallOpts, _witness)
}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _witness) view returns(bool)
func (_WitnessListV3 *WitnessListV3CallerSession) IsAllowed(_witness common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsAllowed(&_WitnessListV3.CallOpts, _witness)
}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3Caller) IsExist(opts *bind.CallOpts, _item common.Address) (bool, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "isExist", _item)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3Session) IsExist(_item common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsExist(&_WitnessListV3.CallOpts, _item)
}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_WitnessListV3 *WitnessListV3CallerSession) IsExist(_item common.Address) (bool, error) {
	return _WitnessListV3.Contract.IsExist(&_WitnessListV3.CallOpts, _item)
}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_WitnessListV3 *WitnessListV3Caller) NumOfActive(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "numOfActive")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_WitnessListV3 *WitnessListV3Session) NumOfActive() (*big.Int, error) {
	return _WitnessListV3.Contract.NumOfActive(&_WitnessListV3.CallOpts)
}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_WitnessListV3 *WitnessListV3CallerSession) NumOfActive() (*big.Int, error) {
	return _WitnessListV3.Contract.NumOfActive(&_WitnessListV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessListV3 *WitnessListV3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessListV3 *WitnessListV3Session) Owner() (common.Address, error) {
	return _WitnessListV3.Contract.Owner(&_WitnessListV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WitnessListV3 *WitnessListV3CallerSession) Owner() (common.Address, error) {
	return _WitnessListV3.Contract.Owner(&_WitnessListV3.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_WitnessListV3 *WitnessListV3Caller) Threshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WitnessListV3.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_WitnessListV3 *WitnessListV3Session) Threshold() (uint8, error) {
	return _WitnessListV3.Contract.Threshold(&_WitnessListV3.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_WitnessListV3 *WitnessListV3CallerSession) Threshold() (uint8, error) {
	return _WitnessListV3.Contract.Threshold(&_WitnessListV3.CallOpts)
}

// AddWitness is a paid mutator transaction binding the contract method 0x59e26be1.
//
// Solidity: function addWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3Transactor) AddWitness(opts *bind.TransactOpts, _witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "addWitness", _witness)
}

// AddWitness is a paid mutator transaction binding the contract method 0x59e26be1.
//
// Solidity: function addWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3Session) AddWitness(_witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.AddWitness(&_WitnessListV3.TransactOpts, _witness)
}

// AddWitness is a paid mutator transaction binding the contract method 0x59e26be1.
//
// Solidity: function addWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3TransactorSession) AddWitness(_witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.AddWitness(&_WitnessListV3.TransactOpts, _witness)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3Transactor) AddWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "addWitnesses", _witnesses)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3Session) AddWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.AddWitnesses(&_WitnessListV3.TransactOpts, _witnesses)
}

// AddWitnesses is a paid mutator transaction binding the contract method 0x14e61a54.
//
// Solidity: function addWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3TransactorSession) AddWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.AddWitnesses(&_WitnessListV3.TransactOpts, _witnesses)
}

// RemoveWitness is a paid mutator transaction binding the contract method 0xee2f13cd.
//
// Solidity: function removeWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3Transactor) RemoveWitness(opts *bind.TransactOpts, _witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "removeWitness", _witness)
}

// RemoveWitness is a paid mutator transaction binding the contract method 0xee2f13cd.
//
// Solidity: function removeWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3Session) RemoveWitness(_witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.RemoveWitness(&_WitnessListV3.TransactOpts, _witness)
}

// RemoveWitness is a paid mutator transaction binding the contract method 0xee2f13cd.
//
// Solidity: function removeWitness(address _witness) returns(bool success_)
func (_WitnessListV3 *WitnessListV3TransactorSession) RemoveWitness(_witness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.RemoveWitness(&_WitnessListV3.TransactOpts, _witness)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3Transactor) RemoveWitnesses(opts *bind.TransactOpts, _witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "removeWitnesses", _witnesses)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3Session) RemoveWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.RemoveWitnesses(&_WitnessListV3.TransactOpts, _witnesses)
}

// RemoveWitnesses is a paid mutator transaction binding the contract method 0xfca2ec37.
//
// Solidity: function removeWitnesses(address[] _witnesses) returns()
func (_WitnessListV3 *WitnessListV3TransactorSession) RemoveWitnesses(_witnesses []common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.RemoveWitnesses(&_WitnessListV3.TransactOpts, _witnesses)
}

// SetThreshold is a paid mutator transaction binding the contract method 0xe5a98603.
//
// Solidity: function setThreshold(uint8 _threshold) returns()
func (_WitnessListV3 *WitnessListV3Transactor) SetThreshold(opts *bind.TransactOpts, _threshold uint8) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "setThreshold", _threshold)
}

// SetThreshold is a paid mutator transaction binding the contract method 0xe5a98603.
//
// Solidity: function setThreshold(uint8 _threshold) returns()
func (_WitnessListV3 *WitnessListV3Session) SetThreshold(_threshold uint8) (*types.Transaction, error) {
	return _WitnessListV3.Contract.SetThreshold(&_WitnessListV3.TransactOpts, _threshold)
}

// SetThreshold is a paid mutator transaction binding the contract method 0xe5a98603.
//
// Solidity: function setThreshold(uint8 _threshold) returns()
func (_WitnessListV3 *WitnessListV3TransactorSession) SetThreshold(_threshold uint8) (*types.Transaction, error) {
	return _WitnessListV3.Contract.SetThreshold(&_WitnessListV3.TransactOpts, _threshold)
}

// SwitchWitness is a paid mutator transaction binding the contract method 0x95935ecf.
//
// Solidity: function switchWitness(address _newWitness) returns()
func (_WitnessListV3 *WitnessListV3Transactor) SwitchWitness(opts *bind.TransactOpts, _newWitness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "switchWitness", _newWitness)
}

// SwitchWitness is a paid mutator transaction binding the contract method 0x95935ecf.
//
// Solidity: function switchWitness(address _newWitness) returns()
func (_WitnessListV3 *WitnessListV3Session) SwitchWitness(_newWitness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.SwitchWitness(&_WitnessListV3.TransactOpts, _newWitness)
}

// SwitchWitness is a paid mutator transaction binding the contract method 0x95935ecf.
//
// Solidity: function switchWitness(address _newWitness) returns()
func (_WitnessListV3 *WitnessListV3TransactorSession) SwitchWitness(_newWitness common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.SwitchWitness(&_WitnessListV3.TransactOpts, _newWitness)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessListV3 *WitnessListV3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WitnessListV3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessListV3 *WitnessListV3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.TransferOwnership(&_WitnessListV3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WitnessListV3 *WitnessListV3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WitnessListV3.Contract.TransferOwnership(&_WitnessListV3.TransactOpts, newOwner)
}

// WitnessListV3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WitnessListV3 contract.
type WitnessListV3OwnershipTransferredIterator struct {
	Event *WitnessListV3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WitnessListV3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WitnessListV3OwnershipTransferred)
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
		it.Event = new(WitnessListV3OwnershipTransferred)
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
func (it *WitnessListV3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WitnessListV3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WitnessListV3OwnershipTransferred represents a OwnershipTransferred event raised by the WitnessListV3 contract.
type WitnessListV3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WitnessListV3 *WitnessListV3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WitnessListV3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WitnessListV3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WitnessListV3OwnershipTransferredIterator{contract: _WitnessListV3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WitnessListV3 *WitnessListV3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WitnessListV3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WitnessListV3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WitnessListV3OwnershipTransferred)
				if err := _WitnessListV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WitnessListV3 *WitnessListV3Filterer) ParseOwnershipTransferred(log types.Log) (*WitnessListV3OwnershipTransferred, error) {
	event := new(WitnessListV3OwnershipTransferred)
	if err := _WitnessListV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WitnessListV3WitnessUpdatedIterator is returned from FilterWitnessUpdated and is used to iterate over the raw logs and unpacked data for WitnessUpdated events raised by the WitnessListV3 contract.
type WitnessListV3WitnessUpdatedIterator struct {
	Event *WitnessListV3WitnessUpdated // Event containing the contract specifics and raw log

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
func (it *WitnessListV3WitnessUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WitnessListV3WitnessUpdated)
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
		it.Event = new(WitnessListV3WitnessUpdated)
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
func (it *WitnessListV3WitnessUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WitnessListV3WitnessUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WitnessListV3WitnessUpdated represents a WitnessUpdated event raised by the WitnessListV3 contract.
type WitnessListV3WitnessUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterWitnessUpdated is a free log retrieval operation binding the contract event 0xe7fedfcf58972c20a587bf272b0a8d43b03b107ace47d9692d662b0e66f62ff5.
//
// Solidity: event witnessUpdated()
func (_WitnessListV3 *WitnessListV3Filterer) FilterWitnessUpdated(opts *bind.FilterOpts) (*WitnessListV3WitnessUpdatedIterator, error) {

	logs, sub, err := _WitnessListV3.contract.FilterLogs(opts, "witnessUpdated")
	if err != nil {
		return nil, err
	}
	return &WitnessListV3WitnessUpdatedIterator{contract: _WitnessListV3.contract, event: "witnessUpdated", logs: logs, sub: sub}, nil
}

// WatchWitnessUpdated is a free log subscription operation binding the contract event 0xe7fedfcf58972c20a587bf272b0a8d43b03b107ace47d9692d662b0e66f62ff5.
//
// Solidity: event witnessUpdated()
func (_WitnessListV3 *WitnessListV3Filterer) WatchWitnessUpdated(opts *bind.WatchOpts, sink chan<- *WitnessListV3WitnessUpdated) (event.Subscription, error) {

	logs, sub, err := _WitnessListV3.contract.WatchLogs(opts, "witnessUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WitnessListV3WitnessUpdated)
				if err := _WitnessListV3.contract.UnpackLog(event, "witnessUpdated", log); err != nil {
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

// ParseWitnessUpdated is a log parse operation binding the contract event 0xe7fedfcf58972c20a587bf272b0a8d43b03b107ace47d9692d662b0e66f62ff5.
//
// Solidity: event witnessUpdated()
func (_WitnessListV3 *WitnessListV3Filterer) ParseWitnessUpdated(log types.Log) (*WitnessListV3WitnessUpdated, error) {
	event := new(WitnessListV3WitnessUpdated)
	if err := _WitnessListV3.contract.UnpackLog(event, "witnessUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
