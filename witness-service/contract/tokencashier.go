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

// TokenCashierABI is the input ABI used to generate the binding from.
const TokenCashierABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenList\",\"outputs\":[{\"internalType\":\"contractTokenList\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"getRecords\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"customers_\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"receivers_\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"fees_\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenCashier is an auto generated Go binding around an Ethereum contract.
type TokenCashier struct {
	TokenCashierCaller     // Read-only binding to the contract
	TokenCashierTransactor // Write-only binding to the contract
	TokenCashierFilterer   // Log filterer for contract events
}

// TokenCashierCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierSession struct {
	Contract     *TokenCashier     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCashierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierCallerSession struct {
	Contract *TokenCashierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenCashierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierTransactorSession struct {
	Contract     *TokenCashierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenCashierRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierRaw struct {
	Contract *TokenCashier // Generic contract binding to access the raw methods on
}

// TokenCashierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierCallerRaw struct {
	Contract *TokenCashierCaller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierTransactorRaw struct {
	Contract *TokenCashierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashier creates a new instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashier(address common.Address, backend bind.ContractBackend) (*TokenCashier, error) {
	contract, err := bindTokenCashier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashier{TokenCashierCaller: TokenCashierCaller{contract: contract}, TokenCashierTransactor: TokenCashierTransactor{contract: contract}, TokenCashierFilterer: TokenCashierFilterer{contract: contract}}, nil
}

// NewTokenCashierCaller creates a new read-only instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierCaller(address common.Address, caller bind.ContractCaller) (*TokenCashierCaller, error) {
	contract, err := bindTokenCashier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierCaller{contract: contract}, nil
}

// NewTokenCashierTransactor creates a new write-only instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierTransactor, error) {
	contract, err := bindTokenCashier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierTransactor{contract: contract}, nil
}

// NewTokenCashierFilterer creates a new log filterer instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierFilterer, error) {
	contract, err := bindTokenCashier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierFilterer{contract: contract}, nil
}

// bindTokenCashier binds a generic wrapper to an already deployed contract.
func bindTokenCashier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenCashierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashier *TokenCashierRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenCashier.Contract.TokenCashierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashier *TokenCashierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.Contract.TokenCashierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashier *TokenCashierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashier.Contract.TokenCashierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashier *TokenCashierCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenCashier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashier *TokenCashierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashier *TokenCashierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashier.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierCaller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenCashier.contract.Call(opts, out, "count", _token)
	return *ret0, err
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Count(&_TokenCashier.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierCallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Count(&_TokenCashier.CallOpts, _token)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenCashier.contract.Call(opts, out, "depositFee")
	return *ret0, err
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierSession) DepositFee() (*big.Int, error) {
	return _TokenCashier.Contract.DepositFee(&_TokenCashier.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierCallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashier.Contract.DepositFee(&_TokenCashier.CallOpts)
}

// GetRecords is a free data retrieval call binding the contract method 0x1bdb94a4.
//
// Solidity: function getRecords(address _token, uint256 _offset, uint256 _limit) view returns(address[] customers_, address[] receivers_, uint256[] amounts_, uint256[] fees_)
func (_TokenCashier *TokenCashierCaller) GetRecords(opts *bind.CallOpts, _token common.Address, _offset *big.Int, _limit *big.Int) (struct {
	Customers []common.Address
	Receivers []common.Address
	Amounts   []*big.Int
	Fees      []*big.Int
}, error) {
	ret := new(struct {
		Customers []common.Address
		Receivers []common.Address
		Amounts   []*big.Int
		Fees      []*big.Int
	})
	out := ret
	err := _TokenCashier.contract.Call(opts, out, "getRecords", _token, _offset, _limit)
	return *ret, err
}

// GetRecords is a free data retrieval call binding the contract method 0x1bdb94a4.
//
// Solidity: function getRecords(address _token, uint256 _offset, uint256 _limit) view returns(address[] customers_, address[] receivers_, uint256[] amounts_, uint256[] fees_)
func (_TokenCashier *TokenCashierSession) GetRecords(_token common.Address, _offset *big.Int, _limit *big.Int) (struct {
	Customers []common.Address
	Receivers []common.Address
	Amounts   []*big.Int
	Fees      []*big.Int
}, error) {
	return _TokenCashier.Contract.GetRecords(&_TokenCashier.CallOpts, _token, _offset, _limit)
}

// GetRecords is a free data retrieval call binding the contract method 0x1bdb94a4.
//
// Solidity: function getRecords(address _token, uint256 _offset, uint256 _limit) view returns(address[] customers_, address[] receivers_, uint256[] amounts_, uint256[] fees_)
func (_TokenCashier *TokenCashierCallerSession) GetRecords(_token common.Address, _offset *big.Int, _limit *big.Int) (struct {
	Customers []common.Address
	Receivers []common.Address
	Amounts   []*big.Int
	Fees      []*big.Int
}, error) {
	return _TokenCashier.Contract.GetRecords(&_TokenCashier.CallOpts, _token, _offset, _limit)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenCashier.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierSession) Owner() (common.Address, error) {
	return _TokenCashier.Contract.Owner(&_TokenCashier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierCallerSession) Owner() (common.Address, error) {
	return _TokenCashier.Contract.Owner(&_TokenCashier.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TokenCashier.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierSession) Paused() (bool, error) {
	return _TokenCashier.Contract.Paused(&_TokenCashier.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierCallerSession) Paused() (bool, error) {
	return _TokenCashier.Contract.Paused(&_TokenCashier.CallOpts)
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashier *TokenCashierCaller) TokenList(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenCashier.contract.Call(opts, out, "tokenList")
	return *ret0, err
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashier *TokenCashierSession) TokenList() (common.Address, error) {
	return _TokenCashier.Contract.TokenList(&_TokenCashier.CallOpts)
}

// TokenList is a free data retrieval call binding the contract method 0x9e2c58ca.
//
// Solidity: function tokenList() view returns(address)
func (_TokenCashier *TokenCashierCallerSession) TokenList() (common.Address, error) {
	return _TokenCashier.Contract.TokenList(&_TokenCashier.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactor) Deposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "deposit", _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierSession) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.Deposit(&_TokenCashier.TransactOpts, _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactorSession) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.Deposit(&_TokenCashier.TransactOpts, _token, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "depositTo", _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.DepositTo(&_TokenCashier.TransactOpts, _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactorSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.DepositTo(&_TokenCashier.TransactOpts, _token, _to, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierSession) Pause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Pause(&_TokenCashier.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierTransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Pause(&_TokenCashier.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.SetDepositFee(&_TokenCashier.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.SetDepositFee(&_TokenCashier.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.TransferOwnership(&_TokenCashier.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.TransferOwnership(&_TokenCashier.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierSession) Unpause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Unpause(&_TokenCashier.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierTransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Unpause(&_TokenCashier.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.Contract.Fallback(&_TokenCashier.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.Contract.Fallback(&_TokenCashier.TransactOpts, calldata)
}

// TokenCashierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashier contract.
type TokenCashierOwnershipTransferredIterator struct {
	Event *TokenCashierOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierOwnershipTransferred)
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
		it.Event = new(TokenCashierOwnershipTransferred)
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
func (it *TokenCashierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierOwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashier contract.
type TokenCashierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashier *TokenCashierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierOwnershipTransferredIterator{contract: _TokenCashier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashier *TokenCashierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierOwnershipTransferred)
				if err := _TokenCashier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierOwnershipTransferred, error) {
	event := new(TokenCashierOwnershipTransferred)
	if err := _TokenCashier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashier contract.
type TokenCashierPauseIterator struct {
	Event *TokenCashierPause // Event containing the contract specifics and raw log

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
func (it *TokenCashierPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierPause)
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
		it.Event = new(TokenCashierPause)
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
func (it *TokenCashierPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierPause represents a Pause event raised by the TokenCashier contract.
type TokenCashierPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashier *TokenCashierFilterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierPauseIterator, error) {

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierPauseIterator{contract: _TokenCashier.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashier *TokenCashierFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierPause) (event.Subscription, error) {

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierPause)
				if err := _TokenCashier.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParsePause(log types.Log) (*TokenCashierPause, error) {
	event := new(TokenCashierPause)
	if err := _TokenCashier.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashier contract.
type TokenCashierReceiptIterator struct {
	Event *TokenCashierReceipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierReceipt)
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
		it.Event = new(TokenCashierReceipt)
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
func (it *TokenCashierReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierReceipt represents a Receipt event raised by the TokenCashier contract.
type TokenCashierReceipt struct {
	Customer common.Address
	Receiver common.Address
	Token    common.Address
	Amount   *big.Int
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0x5a724f2629bc1fea6af562d66f78fa16618c939eacef15aeea7b962c0044779f.
//
// Solidity: event Receipt(address indexed customer, address indexed receiver, address indexed token, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) FilterReceipt(opts *bind.FilterOpts, customer []common.Address, receiver []common.Address, token []common.Address) (*TokenCashierReceiptIterator, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Receipt", customerRule, receiverRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierReceiptIterator{contract: _TokenCashier.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0x5a724f2629bc1fea6af562d66f78fa16618c939eacef15aeea7b962c0044779f.
//
// Solidity: event Receipt(address indexed customer, address indexed receiver, address indexed token, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierReceipt, customer []common.Address, receiver []common.Address, token []common.Address) (event.Subscription, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Receipt", customerRule, receiverRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierReceipt)
				if err := _TokenCashier.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0x5a724f2629bc1fea6af562d66f78fa16618c939eacef15aeea7b962c0044779f.
//
// Solidity: event Receipt(address indexed customer, address indexed receiver, address indexed token, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) ParseReceipt(log types.Log) (*TokenCashierReceipt, error) {
	event := new(TokenCashierReceipt)
	if err := _TokenCashier.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenCashierUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashier contract.
type TokenCashierUnpauseIterator struct {
	Event *TokenCashierUnpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierUnpause)
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
		it.Event = new(TokenCashierUnpause)
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
func (it *TokenCashierUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierUnpause represents a Unpause event raised by the TokenCashier contract.
type TokenCashierUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashier *TokenCashierFilterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierUnpauseIterator, error) {

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierUnpauseIterator{contract: _TokenCashier.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashier *TokenCashierFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierUnpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierUnpause)
				if err := _TokenCashier.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParseUnpause(log types.Log) (*TokenCashierUnpause, error) {
	event := new(TokenCashierUnpause)
	if err := _TokenCashier.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	return event, nil
}
