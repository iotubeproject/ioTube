// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TokenCashierV2ABI is the input ABI used to generate the binding from.
const TokenCashierV2ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"counts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenList\",\"outputs\":[{\"internalType\":\"contractITokenList\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenCashierV2 is an auto generated Go binding around an Ethereum contract.
type TokenCashierV2 struct {
	TokenCashierV2Caller     // Read-only binding to the contract
	TokenCashierV2Transactor // Write-only binding to the contract
	TokenCashierV2Filterer   // Log filterer for contract events
}

// TokenCashierV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierV2Session struct {
	Contract     *TokenCashierV2   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCashierV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierV2CallerSession struct {
	Contract *TokenCashierV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TokenCashierV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierV2TransactorSession struct {
	Contract     *TokenCashierV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TokenCashierV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierV2Raw struct {
	Contract *TokenCashierV2 // Generic contract binding to access the raw methods on
}

// TokenCashierV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierV2CallerRaw struct {
	Contract *TokenCashierV2Caller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierV2TransactorRaw struct {
	Contract *TokenCashierV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashierV2 creates a new instance of TokenCashierV2, bound to a specific deployed contract.
func NewTokenCashierV2(address common.Address, backend bind.ContractBackend) (*TokenCashierV2, error) {
	contract, err := bindTokenCashierV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2{TokenCashierV2Caller: TokenCashierV2Caller{contract: contract}, TokenCashierV2Transactor: TokenCashierV2Transactor{contract: contract}, TokenCashierV2Filterer: TokenCashierV2Filterer{contract: contract}}, nil
}

// NewTokenCashierV2Caller creates a new read-only instance of TokenCashierV2, bound to a specific deployed contract.
func NewTokenCashierV2Caller(address common.Address, caller bind.ContractCaller) (*TokenCashierV2Caller, error) {
	contract, err := bindTokenCashierV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2Caller{contract: contract}, nil
}

// NewTokenCashierV2Transactor creates a new write-only instance of TokenCashierV2, bound to a specific deployed contract.
func NewTokenCashierV2Transactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierV2Transactor, error) {
	contract, err := bindTokenCashierV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2Transactor{contract: contract}, nil
}

// NewTokenCashierV2Filterer creates a new log filterer instance of TokenCashierV2, bound to a specific deployed contract.
func NewTokenCashierV2Filterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierV2Filterer, error) {
	contract, err := bindTokenCashierV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2Filterer{contract: contract}, nil
}

// bindTokenCashierV2 binds a generic wrapper to an already deployed contract.
func bindTokenCashierV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenCashierV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierV2 *TokenCashierV2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenCashierV2.Contract.TokenCashierV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierV2 *TokenCashierV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.TokenCashierV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierV2 *TokenCashierV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.TokenCashierV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierV2 *TokenCashierV2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenCashierV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierV2 *TokenCashierV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierV2 *TokenCashierV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Caller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "count", _token)
	return *ret0, err
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Session) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierV2.Contract.Count(&_TokenCashierV2.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2CallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierV2.Contract.Count(&_TokenCashierV2.CallOpts, _token)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Caller) Counts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "counts", arg0)
	return *ret0, err
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Session) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierV2.Contract.Counts(&_TokenCashierV2.CallOpts, arg0)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2CallerSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierV2.Contract.Counts(&_TokenCashierV2.CallOpts, arg0)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Caller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "depositFee")
	return *ret0, err
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2Session) DepositFee() (*big.Int, error) {
	return _TokenCashierV2.Contract.DepositFee(&_TokenCashierV2.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV2 *TokenCashierV2CallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashierV2.Contract.DepositFee(&_TokenCashierV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV2 *TokenCashierV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV2 *TokenCashierV2Session) Owner() (common.Address, error) {
	return _TokenCashierV2.Contract.Owner(&_TokenCashierV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV2 *TokenCashierV2CallerSession) Owner() (common.Address, error) {
	return _TokenCashierV2.Contract.Owner(&_TokenCashierV2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV2 *TokenCashierV2Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV2 *TokenCashierV2Session) Paused() (bool, error) {
	return _TokenCashierV2.Contract.Paused(&_TokenCashierV2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV2 *TokenCashierV2CallerSession) Paused() (bool, error) {
	return _TokenCashierV2.Contract.Paused(&_TokenCashierV2.CallOpts)
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashierV2 *TokenCashierV2Caller) TokenList(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenCashierV2.contract.Call(opts, out, "tokenList")
	return *ret0, err
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashierV2 *TokenCashierV2Session) TokenList() (common.Address, error) {
	return _TokenCashierV2.Contract.TokenList(&_TokenCashierV2.CallOpts)
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashierV2 *TokenCashierV2CallerSession) TokenList() (common.Address, error) {
	return _TokenCashierV2.Contract.TokenList(&_TokenCashierV2.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) Deposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "deposit", _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2Session) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Deposit(&_TokenCashierV2.TransactOpts, _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Deposit(&_TokenCashierV2.TransactOpts, _token, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "depositTo", _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2Session) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.DepositTo(&_TokenCashierV2.TransactOpts, _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.DepositTo(&_TokenCashierV2.TransactOpts, _token, _to, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV2 *TokenCashierV2Session) Pause() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Pause(&_TokenCashierV2.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Pause(&_TokenCashierV2.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV2 *TokenCashierV2Session) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.SetDepositFee(&_TokenCashierV2.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.SetDepositFee(&_TokenCashierV2.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV2 *TokenCashierV2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.TransferOwnership(&_TokenCashierV2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.TransferOwnership(&_TokenCashierV2.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV2 *TokenCashierV2Session) Unpause() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Unpause(&_TokenCashierV2.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Unpause(&_TokenCashierV2.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV2.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV2 *TokenCashierV2Session) Withdraw() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Withdraw(&_TokenCashierV2.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Withdraw(&_TokenCashierV2.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV2 *TokenCashierV2Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV2.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV2 *TokenCashierV2Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Fallback(&_TokenCashierV2.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV2 *TokenCashierV2TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV2.Contract.Fallback(&_TokenCashierV2.TransactOpts, calldata)
}

// TokenCashierV2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashierV2 contract.
type TokenCashierV2OwnershipTransferredIterator struct {
	Event *TokenCashierV2OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierV2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV2OwnershipTransferred)
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
		it.Event = new(TokenCashierV2OwnershipTransferred)
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
func (it *TokenCashierV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV2OwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashierV2 contract.
type TokenCashierV2OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierV2 *TokenCashierV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierV2OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierV2.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2OwnershipTransferredIterator{contract: _TokenCashierV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierV2 *TokenCashierV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierV2OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierV2.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV2OwnershipTransferred)
				if err := _TokenCashierV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashierV2 *TokenCashierV2Filterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierV2OwnershipTransferred, error) {
	event := new(TokenCashierV2OwnershipTransferred)
	if err := _TokenCashierV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierV2PauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashierV2 contract.
type TokenCashierV2PauseIterator struct {
	Event *TokenCashierV2Pause // Event containing the contract specifics and raw log

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
func (it *TokenCashierV2PauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV2Pause)
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
		it.Event = new(TokenCashierV2Pause)
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
func (it *TokenCashierV2PauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV2PauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV2Pause represents a Pause event raised by the TokenCashierV2 contract.
type TokenCashierV2Pause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierV2 *TokenCashierV2Filterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierV2PauseIterator, error) {

	logs, sub, err := _TokenCashierV2.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2PauseIterator{contract: _TokenCashierV2.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierV2 *TokenCashierV2Filterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierV2Pause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierV2.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV2Pause)
				if err := _TokenCashierV2.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashierV2 *TokenCashierV2Filterer) ParsePause(log types.Log) (*TokenCashierV2Pause, error) {
	event := new(TokenCashierV2Pause)
	if err := _TokenCashierV2.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierV2ReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashierV2 contract.
type TokenCashierV2ReceiptIterator struct {
	Event *TokenCashierV2Receipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierV2ReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV2Receipt)
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
		it.Event = new(TokenCashierV2Receipt)
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
func (it *TokenCashierV2ReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV2ReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV2Receipt represents a Receipt event raised by the TokenCashierV2 contract.
type TokenCashierV2Receipt struct {
	Token    common.Address
	Id       *big.Int
	Receiver common.Address
	Sender   common.Address
	Amount   *big.Int
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address receiver, address sender, uint256 amount, uint256 fee)
func (_TokenCashierV2 *TokenCashierV2Filterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierV2ReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierV2.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2ReceiptIterator{contract: _TokenCashierV2.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address receiver, address sender, uint256 amount, uint256 fee)
func (_TokenCashierV2 *TokenCashierV2Filterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierV2Receipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierV2.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV2Receipt)
				if err := _TokenCashierV2.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address receiver, address sender, uint256 amount, uint256 fee)
func (_TokenCashierV2 *TokenCashierV2Filterer) ParseReceipt(log types.Log) (*TokenCashierV2Receipt, error) {
	event := new(TokenCashierV2Receipt)
	if err := _TokenCashierV2.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierV2UnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashierV2 contract.
type TokenCashierV2UnpauseIterator struct {
	Event *TokenCashierV2Unpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierV2UnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV2Unpause)
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
		it.Event = new(TokenCashierV2Unpause)
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
func (it *TokenCashierV2UnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV2UnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV2Unpause represents a Unpause event raised by the TokenCashierV2 contract.
type TokenCashierV2Unpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierV2 *TokenCashierV2Filterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierV2UnpauseIterator, error) {

	logs, sub, err := _TokenCashierV2.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierV2UnpauseIterator{contract: _TokenCashierV2.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierV2 *TokenCashierV2Filterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierV2Unpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierV2.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV2Unpause)
				if err := _TokenCashierV2.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashierV2 *TokenCashierV2Filterer) ParseUnpause(log types.Log) (*TokenCashierV2Unpause, error) {
	event := new(TokenCashierV2Unpause)
	if err := _TokenCashierV2.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	return event, nil
}
