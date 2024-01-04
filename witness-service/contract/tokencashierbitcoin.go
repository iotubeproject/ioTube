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

// TokenCashierBitcoinMetaData contains all meta data concerning the TokenCashierBitcoin contract.
var TokenCashierBitcoinMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bitcoin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bitcoin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TokenCashierBitcoinABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenCashierBitcoinMetaData.ABI instead.
var TokenCashierBitcoinABI = TokenCashierBitcoinMetaData.ABI

// TokenCashierBitcoin is an auto generated Go binding around an Ethereum contract.
type TokenCashierBitcoin struct {
	TokenCashierBitcoinCaller     // Read-only binding to the contract
	TokenCashierBitcoinTransactor // Write-only binding to the contract
	TokenCashierBitcoinFilterer   // Log filterer for contract events
}

// TokenCashierBitcoinCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierBitcoinCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierBitcoinTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierBitcoinTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierBitcoinFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierBitcoinFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierBitcoinSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierBitcoinSession struct {
	Contract     *TokenCashierBitcoin // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TokenCashierBitcoinCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierBitcoinCallerSession struct {
	Contract *TokenCashierBitcoinCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TokenCashierBitcoinTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierBitcoinTransactorSession struct {
	Contract     *TokenCashierBitcoinTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TokenCashierBitcoinRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierBitcoinRaw struct {
	Contract *TokenCashierBitcoin // Generic contract binding to access the raw methods on
}

// TokenCashierBitcoinCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierBitcoinCallerRaw struct {
	Contract *TokenCashierBitcoinCaller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierBitcoinTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierBitcoinTransactorRaw struct {
	Contract *TokenCashierBitcoinTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashierBitcoin creates a new instance of TokenCashierBitcoin, bound to a specific deployed contract.
func NewTokenCashierBitcoin(address common.Address, backend bind.ContractBackend) (*TokenCashierBitcoin, error) {
	contract, err := bindTokenCashierBitcoin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoin{TokenCashierBitcoinCaller: TokenCashierBitcoinCaller{contract: contract}, TokenCashierBitcoinTransactor: TokenCashierBitcoinTransactor{contract: contract}, TokenCashierBitcoinFilterer: TokenCashierBitcoinFilterer{contract: contract}}, nil
}

// NewTokenCashierBitcoinCaller creates a new read-only instance of TokenCashierBitcoin, bound to a specific deployed contract.
func NewTokenCashierBitcoinCaller(address common.Address, caller bind.ContractCaller) (*TokenCashierBitcoinCaller, error) {
	contract, err := bindTokenCashierBitcoin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinCaller{contract: contract}, nil
}

// NewTokenCashierBitcoinTransactor creates a new write-only instance of TokenCashierBitcoin, bound to a specific deployed contract.
func NewTokenCashierBitcoinTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierBitcoinTransactor, error) {
	contract, err := bindTokenCashierBitcoin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinTransactor{contract: contract}, nil
}

// NewTokenCashierBitcoinFilterer creates a new log filterer instance of TokenCashierBitcoin, bound to a specific deployed contract.
func NewTokenCashierBitcoinFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierBitcoinFilterer, error) {
	contract, err := bindTokenCashierBitcoin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinFilterer{contract: contract}, nil
}

// bindTokenCashierBitcoin binds a generic wrapper to an already deployed contract.
func bindTokenCashierBitcoin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenCashierBitcoinMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierBitcoin *TokenCashierBitcoinRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierBitcoin.Contract.TokenCashierBitcoinCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierBitcoin *TokenCashierBitcoinRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.TokenCashierBitcoinTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierBitcoin *TokenCashierBitcoinRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.TokenCashierBitcoinTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierBitcoin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.contract.Transact(opts, method, params...)
}

// Bitcoin is a free data retrieval call binding the contract method 0xced35070.
//
// Solidity: function bitcoin() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinCaller) Bitcoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierBitcoin.contract.Call(opts, &out, "bitcoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bitcoin is a free data retrieval call binding the contract method 0xced35070.
//
// Solidity: function bitcoin() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Bitcoin() (common.Address, error) {
	return _TokenCashierBitcoin.Contract.Bitcoin(&_TokenCashierBitcoin.CallOpts)
}

// Bitcoin is a free data retrieval call binding the contract method 0xced35070.
//
// Solidity: function bitcoin() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerSession) Bitcoin() (common.Address, error) {
	return _TokenCashierBitcoin.Contract.Bitcoin(&_TokenCashierBitcoin.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierBitcoin.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Count() (*big.Int, error) {
	return _TokenCashierBitcoin.Contract.Count(&_TokenCashierBitcoin.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerSession) Count() (*big.Int, error) {
	return _TokenCashierBitcoin.Contract.Count(&_TokenCashierBitcoin.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierBitcoin.contract.Call(opts, &out, "depositFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) DepositFee() (*big.Int, error) {
	return _TokenCashierBitcoin.Contract.DepositFee(&_TokenCashierBitcoin.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashierBitcoin.Contract.DepositFee(&_TokenCashierBitcoin.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierBitcoin.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Owner() (common.Address, error) {
	return _TokenCashierBitcoin.Contract.Owner(&_TokenCashierBitcoin.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerSession) Owner() (common.Address, error) {
	return _TokenCashierBitcoin.Contract.Owner(&_TokenCashierBitcoin.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierBitcoin *TokenCashierBitcoinCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenCashierBitcoin.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Paused() (bool, error) {
	return _TokenCashierBitcoin.Contract.Paused(&_TokenCashierBitcoin.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierBitcoin *TokenCashierBitcoinCallerSession) Paused() (bool, error) {
	return _TokenCashierBitcoin.Contract.Paused(&_TokenCashierBitcoin.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bitcoin) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) Initialize(opts *bind.TransactOpts, _bitcoin common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "initialize", _bitcoin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bitcoin) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Initialize(_bitcoin common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Initialize(&_TokenCashierBitcoin.TransactOpts, _bitcoin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bitcoin) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) Initialize(_bitcoin common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Initialize(&_TokenCashierBitcoin.TransactOpts, _bitcoin)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Pause() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Pause(&_TokenCashierBitcoin.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Pause(&_TokenCashierBitcoin.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.RenounceOwnership(&_TokenCashierBitcoin.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.RenounceOwnership(&_TokenCashierBitcoin.TransactOpts)
}

// Report is a paid mutator transaction binding the contract method 0x24a74b9a.
//
// Solidity: function report(address _sender, string _to, uint256 _amount) payable returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) Report(opts *bind.TransactOpts, _sender common.Address, _to string, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "report", _sender, _to, _amount)
}

// Report is a paid mutator transaction binding the contract method 0x24a74b9a.
//
// Solidity: function report(address _sender, string _to, uint256 _amount) payable returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Report(_sender common.Address, _to string, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Report(&_TokenCashierBitcoin.TransactOpts, _sender, _to, _amount)
}

// Report is a paid mutator transaction binding the contract method 0x24a74b9a.
//
// Solidity: function report(address _sender, string _to, uint256 _amount) payable returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) Report(_sender common.Address, _to string, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Report(&_TokenCashierBitcoin.TransactOpts, _sender, _to, _amount)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.SetDepositFee(&_TokenCashierBitcoin.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.SetDepositFee(&_TokenCashierBitcoin.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.TransferOwnership(&_TokenCashierBitcoin.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.TransferOwnership(&_TokenCashierBitcoin.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Unpause(&_TokenCashierBitcoin.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Unpause(&_TokenCashierBitcoin.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Withdraw(&_TokenCashierBitcoin.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.Withdraw(&_TokenCashierBitcoin.TransactOpts)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.contract.Transact(opts, "withdrawToken", _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.WithdrawToken(&_TokenCashierBitcoin.TransactOpts, _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierBitcoin *TokenCashierBitcoinTransactorSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierBitcoin.Contract.WithdrawToken(&_TokenCashierBitcoin.TransactOpts, _token)
}

// TokenCashierBitcoinInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinInitializedIterator struct {
	Event *TokenCashierBitcoinInitialized // Event containing the contract specifics and raw log

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
func (it *TokenCashierBitcoinInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierBitcoinInitialized)
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
		it.Event = new(TokenCashierBitcoinInitialized)
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
func (it *TokenCashierBitcoinInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierBitcoinInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierBitcoinInitialized represents a Initialized event raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) FilterInitialized(opts *bind.FilterOpts) (*TokenCashierBitcoinInitializedIterator, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinInitializedIterator{contract: _TokenCashierBitcoin.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TokenCashierBitcoinInitialized) (event.Subscription, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierBitcoinInitialized)
				if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) ParseInitialized(log types.Log) (*TokenCashierBitcoinInitialized, error) {
	event := new(TokenCashierBitcoinInitialized)
	if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierBitcoinOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinOwnershipTransferredIterator struct {
	Event *TokenCashierBitcoinOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierBitcoinOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierBitcoinOwnershipTransferred)
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
		it.Event = new(TokenCashierBitcoinOwnershipTransferred)
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
func (it *TokenCashierBitcoinOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierBitcoinOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierBitcoinOwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierBitcoinOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierBitcoin.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinOwnershipTransferredIterator{contract: _TokenCashierBitcoin.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierBitcoinOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierBitcoin.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierBitcoinOwnershipTransferred)
				if err := _TokenCashierBitcoin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierBitcoinOwnershipTransferred, error) {
	event := new(TokenCashierBitcoinOwnershipTransferred)
	if err := _TokenCashierBitcoin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierBitcoinPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinPauseIterator struct {
	Event *TokenCashierBitcoinPause // Event containing the contract specifics and raw log

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
func (it *TokenCashierBitcoinPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierBitcoinPause)
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
		it.Event = new(TokenCashierBitcoinPause)
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
func (it *TokenCashierBitcoinPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierBitcoinPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierBitcoinPause represents a Pause event raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierBitcoinPauseIterator, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinPauseIterator{contract: _TokenCashierBitcoin.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierBitcoinPause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierBitcoinPause)
				if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Pause", log); err != nil {
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

// ParsePause is a log parse operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) ParsePause(log types.Log) (*TokenCashierBitcoinPause, error) {
	event := new(TokenCashierBitcoinPause)
	if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierBitcoinReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinReceiptIterator struct {
	Event *TokenCashierBitcoinReceipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierBitcoinReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierBitcoinReceipt)
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
		it.Event = new(TokenCashierBitcoinReceipt)
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
func (it *TokenCashierBitcoinReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierBitcoinReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierBitcoinReceipt represents a Receipt event raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinReceipt struct {
	Token     common.Address
	Id        *big.Int
	Sender    common.Address
	Recipient string
	Amount    *big.Int
	Fee       *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0xc59b3b73be4e0b3a1cd9e30aa4a4c13883a774a91c35fc607669f741dc9296e6.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierBitcoinReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierBitcoin.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinReceiptIterator{contract: _TokenCashierBitcoin.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0xc59b3b73be4e0b3a1cd9e30aa4a4c13883a774a91c35fc607669f741dc9296e6.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierBitcoinReceipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierBitcoin.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierBitcoinReceipt)
				if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0xc59b3b73be4e0b3a1cd9e30aa4a4c13883a774a91c35fc607669f741dc9296e6.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee)
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) ParseReceipt(log types.Log) (*TokenCashierBitcoinReceipt, error) {
	event := new(TokenCashierBitcoinReceipt)
	if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierBitcoinUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinUnpauseIterator struct {
	Event *TokenCashierBitcoinUnpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierBitcoinUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierBitcoinUnpause)
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
		it.Event = new(TokenCashierBitcoinUnpause)
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
func (it *TokenCashierBitcoinUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierBitcoinUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierBitcoinUnpause represents a Unpause event raised by the TokenCashierBitcoin contract.
type TokenCashierBitcoinUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierBitcoinUnpauseIterator, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierBitcoinUnpauseIterator{contract: _TokenCashierBitcoin.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierBitcoinUnpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierBitcoin.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierBitcoinUnpause)
				if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Unpause", log); err != nil {
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

// ParseUnpause is a log parse operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierBitcoin *TokenCashierBitcoinFilterer) ParseUnpause(log types.Log) (*TokenCashierBitcoinUnpause, error) {
	event := new(TokenCashierBitcoinUnpause)
	if err := _TokenCashierBitcoin.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
