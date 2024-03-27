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

// TokenCashierV3ABI is the input ABI used to generate the binding from.
const TokenCashierV3ABI = "[{\"inputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"_wrappedCoin\",\"type\":\"address\"},{\"internalType\":\"contractITokenList[]\",\"name\":\"_tokenLists\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_tokenSafes\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"counts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractITokenList\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenSafes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wrappedCoin\",\"outputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenCashierV3 is an auto generated Go binding around an Ethereum contract.
type TokenCashierV3 struct {
	TokenCashierV3Caller     // Read-only binding to the contract
	TokenCashierV3Transactor // Write-only binding to the contract
	TokenCashierV3Filterer   // Log filterer for contract events
}

// TokenCashierV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierV3Session struct {
	Contract     *TokenCashierV3   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCashierV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierV3CallerSession struct {
	Contract *TokenCashierV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TokenCashierV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierV3TransactorSession struct {
	Contract     *TokenCashierV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TokenCashierV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierV3Raw struct {
	Contract *TokenCashierV3 // Generic contract binding to access the raw methods on
}

// TokenCashierV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierV3CallerRaw struct {
	Contract *TokenCashierV3Caller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierV3TransactorRaw struct {
	Contract *TokenCashierV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashierV3 creates a new instance of TokenCashierV3, bound to a specific deployed contract.
func NewTokenCashierV3(address common.Address, backend bind.ContractBackend) (*TokenCashierV3, error) {
	contract, err := bindTokenCashierV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3{TokenCashierV3Caller: TokenCashierV3Caller{contract: contract}, TokenCashierV3Transactor: TokenCashierV3Transactor{contract: contract}, TokenCashierV3Filterer: TokenCashierV3Filterer{contract: contract}}, nil
}

// NewTokenCashierV3Caller creates a new read-only instance of TokenCashierV3, bound to a specific deployed contract.
func NewTokenCashierV3Caller(address common.Address, caller bind.ContractCaller) (*TokenCashierV3Caller, error) {
	contract, err := bindTokenCashierV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3Caller{contract: contract}, nil
}

// NewTokenCashierV3Transactor creates a new write-only instance of TokenCashierV3, bound to a specific deployed contract.
func NewTokenCashierV3Transactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierV3Transactor, error) {
	contract, err := bindTokenCashierV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3Transactor{contract: contract}, nil
}

// NewTokenCashierV3Filterer creates a new log filterer instance of TokenCashierV3, bound to a specific deployed contract.
func NewTokenCashierV3Filterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierV3Filterer, error) {
	contract, err := bindTokenCashierV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3Filterer{contract: contract}, nil
}

// bindTokenCashierV3 binds a generic wrapper to an already deployed contract.
func bindTokenCashierV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenCashierV3ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierV3 *TokenCashierV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierV3.Contract.TokenCashierV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierV3 *TokenCashierV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.TokenCashierV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierV3 *TokenCashierV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.TokenCashierV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierV3 *TokenCashierV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierV3 *TokenCashierV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierV3 *TokenCashierV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Caller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "count", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Session) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierV3.Contract.Count(&_TokenCashierV3.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3CallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierV3.Contract.Count(&_TokenCashierV3.CallOpts, _token)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Caller) Counts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "counts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Session) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierV3.Contract.Counts(&_TokenCashierV3.CallOpts, arg0)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3CallerSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierV3.Contract.Counts(&_TokenCashierV3.CallOpts, arg0)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Caller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "depositFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3Session) DepositFee() (*big.Int, error) {
	return _TokenCashierV3.Contract.DepositFee(&_TokenCashierV3.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierV3 *TokenCashierV3CallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashierV3.Contract.DepositFee(&_TokenCashierV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV3 *TokenCashierV3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV3 *TokenCashierV3Session) Owner() (common.Address, error) {
	return _TokenCashierV3.Contract.Owner(&_TokenCashierV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierV3 *TokenCashierV3CallerSession) Owner() (common.Address, error) {
	return _TokenCashierV3.Contract.Owner(&_TokenCashierV3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV3 *TokenCashierV3Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV3 *TokenCashierV3Session) Paused() (bool, error) {
	return _TokenCashierV3.Contract.Paused(&_TokenCashierV3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierV3 *TokenCashierV3CallerSession) Paused() (bool, error) {
	return _TokenCashierV3.Contract.Paused(&_TokenCashierV3.CallOpts)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3Caller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3Session) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierV3.Contract.TokenLists(&_TokenCashierV3.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3CallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierV3.Contract.TokenLists(&_TokenCashierV3.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3Caller) TokenSafes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "tokenSafes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3Session) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierV3.Contract.TokenSafes(&_TokenCashierV3.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierV3 *TokenCashierV3CallerSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierV3.Contract.TokenSafes(&_TokenCashierV3.CallOpts, arg0)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierV3 *TokenCashierV3Caller) WrappedCoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierV3.contract.Call(opts, &out, "wrappedCoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierV3 *TokenCashierV3Session) WrappedCoin() (common.Address, error) {
	return _TokenCashierV3.Contract.WrappedCoin(&_TokenCashierV3.CallOpts)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierV3 *TokenCashierV3CallerSession) WrappedCoin() (common.Address, error) {
	return _TokenCashierV3.Contract.WrappedCoin(&_TokenCashierV3.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) Deposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "deposit", _token, _amount, _payload)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3Session) Deposit(_token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Deposit(&_TokenCashierV3.TransactOpts, _token, _amount, _payload)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) Deposit(_token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Deposit(&_TokenCashierV3.TransactOpts, _token, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "depositTo", _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3Session) DepositTo(_token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.DepositTo(&_TokenCashierV3.TransactOpts, _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.DepositTo(&_TokenCashierV3.TransactOpts, _token, _to, _amount, _payload)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV3 *TokenCashierV3Session) Pause() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Pause(&_TokenCashierV3.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Pause(&_TokenCashierV3.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV3 *TokenCashierV3Session) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.SetDepositFee(&_TokenCashierV3.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.SetDepositFee(&_TokenCashierV3.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV3 *TokenCashierV3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.TransferOwnership(&_TokenCashierV3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.TransferOwnership(&_TokenCashierV3.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV3 *TokenCashierV3Session) Unpause() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Unpause(&_TokenCashierV3.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Unpause(&_TokenCashierV3.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV3 *TokenCashierV3Session) Withdraw() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Withdraw(&_TokenCashierV3.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Withdraw(&_TokenCashierV3.TransactOpts)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.contract.Transact(opts, "withdrawToken", _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierV3 *TokenCashierV3Session) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.WithdrawToken(&_TokenCashierV3.TransactOpts, _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.WithdrawToken(&_TokenCashierV3.TransactOpts, _token)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV3 *TokenCashierV3Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV3.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV3 *TokenCashierV3Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Fallback(&_TokenCashierV3.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierV3 *TokenCashierV3TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierV3.Contract.Fallback(&_TokenCashierV3.TransactOpts, calldata)
}

// TokenCashierV3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashierV3 contract.
type TokenCashierV3OwnershipTransferredIterator struct {
	Event *TokenCashierV3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierV3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV3OwnershipTransferred)
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
		it.Event = new(TokenCashierV3OwnershipTransferred)
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
func (it *TokenCashierV3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV3OwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashierV3 contract.
type TokenCashierV3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierV3 *TokenCashierV3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierV3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierV3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3OwnershipTransferredIterator{contract: _TokenCashierV3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierV3 *TokenCashierV3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierV3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierV3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV3OwnershipTransferred)
				if err := _TokenCashierV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashierV3 *TokenCashierV3Filterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierV3OwnershipTransferred, error) {
	event := new(TokenCashierV3OwnershipTransferred)
	if err := _TokenCashierV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierV3PauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashierV3 contract.
type TokenCashierV3PauseIterator struct {
	Event *TokenCashierV3Pause // Event containing the contract specifics and raw log

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
func (it *TokenCashierV3PauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV3Pause)
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
		it.Event = new(TokenCashierV3Pause)
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
func (it *TokenCashierV3PauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV3PauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV3Pause represents a Pause event raised by the TokenCashierV3 contract.
type TokenCashierV3Pause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierV3 *TokenCashierV3Filterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierV3PauseIterator, error) {

	logs, sub, err := _TokenCashierV3.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3PauseIterator{contract: _TokenCashierV3.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierV3 *TokenCashierV3Filterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierV3Pause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierV3.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV3Pause)
				if err := _TokenCashierV3.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashierV3 *TokenCashierV3Filterer) ParsePause(log types.Log) (*TokenCashierV3Pause, error) {
	event := new(TokenCashierV3Pause)
	if err := _TokenCashierV3.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierV3ReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashierV3 contract.
type TokenCashierV3ReceiptIterator struct {
	Event *TokenCashierV3Receipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierV3ReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV3Receipt)
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
		it.Event = new(TokenCashierV3Receipt)
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
func (it *TokenCashierV3ReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV3ReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV3Receipt represents a Receipt event raised by the TokenCashierV3 contract.
type TokenCashierV3Receipt struct {
	Token     common.Address
	Id        *big.Int
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Fee       *big.Int
	Payload   []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0xd2be25887579d6d0dc43743403c85c398b3873c57506ad20610cef12f2a3c9d2.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierV3 *TokenCashierV3Filterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierV3ReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierV3.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3ReceiptIterator{contract: _TokenCashierV3.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0xd2be25887579d6d0dc43743403c85c398b3873c57506ad20610cef12f2a3c9d2.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierV3 *TokenCashierV3Filterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierV3Receipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierV3.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV3Receipt)
				if err := _TokenCashierV3.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0xd2be25887579d6d0dc43743403c85c398b3873c57506ad20610cef12f2a3c9d2.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierV3 *TokenCashierV3Filterer) ParseReceipt(log types.Log) (*TokenCashierV3Receipt, error) {
	event := new(TokenCashierV3Receipt)
	if err := _TokenCashierV3.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierV3UnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashierV3 contract.
type TokenCashierV3UnpauseIterator struct {
	Event *TokenCashierV3Unpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierV3UnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierV3Unpause)
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
		it.Event = new(TokenCashierV3Unpause)
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
func (it *TokenCashierV3UnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierV3UnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierV3Unpause represents a Unpause event raised by the TokenCashierV3 contract.
type TokenCashierV3Unpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierV3 *TokenCashierV3Filterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierV3UnpauseIterator, error) {

	logs, sub, err := _TokenCashierV3.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierV3UnpauseIterator{contract: _TokenCashierV3.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierV3 *TokenCashierV3Filterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierV3Unpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierV3.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierV3Unpause)
				if err := _TokenCashierV3.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashierV3 *TokenCashierV3Filterer) ParseUnpause(log types.Log) (*TokenCashierV3Unpause, error) {
	event := new(TokenCashierV3Unpause)
	if err := _TokenCashierV3.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
