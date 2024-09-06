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

// TokenCashierForSolanaMetaData contains all meta data concerning the TokenCashierForSolana contract.
var TokenCashierForSolanaMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"_wrappedCoin\",\"type\":\"address\"},{\"internalType\":\"contractITokenList[]\",\"name\":\"_tokenLists\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_tokenSafes\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"counts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractITokenList\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenSafes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wrappedCoin\",\"outputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TokenCashierForSolanaABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenCashierForSolanaMetaData.ABI instead.
var TokenCashierForSolanaABI = TokenCashierForSolanaMetaData.ABI

// TokenCashierForSolana is an auto generated Go binding around an Ethereum contract.
type TokenCashierForSolana struct {
	TokenCashierForSolanaCaller     // Read-only binding to the contract
	TokenCashierForSolanaTransactor // Write-only binding to the contract
	TokenCashierForSolanaFilterer   // Log filterer for contract events
}

// TokenCashierForSolanaCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierForSolanaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierForSolanaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierForSolanaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierForSolanaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierForSolanaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierForSolanaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierForSolanaSession struct {
	Contract     *TokenCashierForSolana // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TokenCashierForSolanaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierForSolanaCallerSession struct {
	Contract *TokenCashierForSolanaCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// TokenCashierForSolanaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierForSolanaTransactorSession struct {
	Contract     *TokenCashierForSolanaTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// TokenCashierForSolanaRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierForSolanaRaw struct {
	Contract *TokenCashierForSolana // Generic contract binding to access the raw methods on
}

// TokenCashierForSolanaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierForSolanaCallerRaw struct {
	Contract *TokenCashierForSolanaCaller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierForSolanaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierForSolanaTransactorRaw struct {
	Contract *TokenCashierForSolanaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashierForSolana creates a new instance of TokenCashierForSolana, bound to a specific deployed contract.
func NewTokenCashierForSolana(address common.Address, backend bind.ContractBackend) (*TokenCashierForSolana, error) {
	contract, err := bindTokenCashierForSolana(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolana{TokenCashierForSolanaCaller: TokenCashierForSolanaCaller{contract: contract}, TokenCashierForSolanaTransactor: TokenCashierForSolanaTransactor{contract: contract}, TokenCashierForSolanaFilterer: TokenCashierForSolanaFilterer{contract: contract}}, nil
}

// NewTokenCashierForSolanaCaller creates a new read-only instance of TokenCashierForSolana, bound to a specific deployed contract.
func NewTokenCashierForSolanaCaller(address common.Address, caller bind.ContractCaller) (*TokenCashierForSolanaCaller, error) {
	contract, err := bindTokenCashierForSolana(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaCaller{contract: contract}, nil
}

// NewTokenCashierForSolanaTransactor creates a new write-only instance of TokenCashierForSolana, bound to a specific deployed contract.
func NewTokenCashierForSolanaTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierForSolanaTransactor, error) {
	contract, err := bindTokenCashierForSolana(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaTransactor{contract: contract}, nil
}

// NewTokenCashierForSolanaFilterer creates a new log filterer instance of TokenCashierForSolana, bound to a specific deployed contract.
func NewTokenCashierForSolanaFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierForSolanaFilterer, error) {
	contract, err := bindTokenCashierForSolana(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaFilterer{contract: contract}, nil
}

// bindTokenCashierForSolana binds a generic wrapper to an already deployed contract.
func bindTokenCashierForSolana(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenCashierForSolanaMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierForSolana *TokenCashierForSolanaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierForSolana.Contract.TokenCashierForSolanaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierForSolana *TokenCashierForSolanaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.TokenCashierForSolanaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierForSolana *TokenCashierForSolanaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.TokenCashierForSolanaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierForSolana *TokenCashierForSolanaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierForSolana.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "count", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierForSolana.Contract.Count(&_TokenCashierForSolana.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierForSolana.Contract.Count(&_TokenCashierForSolana.CallOpts, _token)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) Counts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "counts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierForSolana.Contract.Counts(&_TokenCashierForSolana.CallOpts, arg0)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierForSolana.Contract.Counts(&_TokenCashierForSolana.CallOpts, arg0)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "depositFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) DepositFee() (*big.Int, error) {
	return _TokenCashierForSolana.Contract.DepositFee(&_TokenCashierForSolana.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashierForSolana.Contract.DepositFee(&_TokenCashierForSolana.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Owner() (common.Address, error) {
	return _TokenCashierForSolana.Contract.Owner(&_TokenCashierForSolana.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) Owner() (common.Address, error) {
	return _TokenCashierForSolana.Contract.Owner(&_TokenCashierForSolana.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Paused() (bool, error) {
	return _TokenCashierForSolana.Contract.Paused(&_TokenCashierForSolana.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) Paused() (bool, error) {
	return _TokenCashierForSolana.Contract.Paused(&_TokenCashierForSolana.CallOpts)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierForSolana.Contract.TokenLists(&_TokenCashierForSolana.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierForSolana.Contract.TokenLists(&_TokenCashierForSolana.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) TokenSafes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "tokenSafes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierForSolana.Contract.TokenSafes(&_TokenCashierForSolana.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierForSolana.Contract.TokenSafes(&_TokenCashierForSolana.CallOpts, arg0)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCaller) WrappedCoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierForSolana.contract.Call(opts, &out, "wrappedCoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaSession) WrappedCoin() (common.Address, error) {
	return _TokenCashierForSolana.Contract.WrappedCoin(&_TokenCashierForSolana.CallOpts)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierForSolana *TokenCashierForSolanaCallerSession) WrappedCoin() (common.Address, error) {
	return _TokenCashierForSolana.Contract.WrappedCoin(&_TokenCashierForSolana.CallOpts)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd13b5612.
//
// Solidity: function depositTo(address _token, string _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to string, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "depositTo", _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd13b5612.
//
// Solidity: function depositTo(address _token, string _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) DepositTo(_token common.Address, _to string, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.DepositTo(&_TokenCashierForSolana.TransactOpts, _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd13b5612.
//
// Solidity: function depositTo(address _token, string _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) DepositTo(_token common.Address, _to string, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.DepositTo(&_TokenCashierForSolana.TransactOpts, _token, _to, _amount, _payload)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Pause() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Pause(&_TokenCashierForSolana.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Pause(&_TokenCashierForSolana.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.SetDepositFee(&_TokenCashierForSolana.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.SetDepositFee(&_TokenCashierForSolana.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.TransferOwnership(&_TokenCashierForSolana.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.TransferOwnership(&_TokenCashierForSolana.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Unpause(&_TokenCashierForSolana.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Unpause(&_TokenCashierForSolana.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Withdraw(&_TokenCashierForSolana.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Withdraw(&_TokenCashierForSolana.TransactOpts)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.Transact(opts, "withdrawToken", _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.WithdrawToken(&_TokenCashierForSolana.TransactOpts, _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.WithdrawToken(&_TokenCashierForSolana.TransactOpts, _token)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Fallback(&_TokenCashierForSolana.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierForSolana *TokenCashierForSolanaTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierForSolana.Contract.Fallback(&_TokenCashierForSolana.TransactOpts, calldata)
}

// TokenCashierForSolanaOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaOwnershipTransferredIterator struct {
	Event *TokenCashierForSolanaOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierForSolanaOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierForSolanaOwnershipTransferred)
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
		it.Event = new(TokenCashierForSolanaOwnershipTransferred)
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
func (it *TokenCashierForSolanaOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierForSolanaOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierForSolanaOwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierForSolanaOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierForSolana.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaOwnershipTransferredIterator{contract: _TokenCashierForSolana.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierForSolanaOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierForSolana.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierForSolanaOwnershipTransferred)
				if err := _TokenCashierForSolana.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierForSolanaOwnershipTransferred, error) {
	event := new(TokenCashierForSolanaOwnershipTransferred)
	if err := _TokenCashierForSolana.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierForSolanaPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaPauseIterator struct {
	Event *TokenCashierForSolanaPause // Event containing the contract specifics and raw log

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
func (it *TokenCashierForSolanaPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierForSolanaPause)
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
		it.Event = new(TokenCashierForSolanaPause)
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
func (it *TokenCashierForSolanaPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierForSolanaPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierForSolanaPause represents a Pause event raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierForSolanaPauseIterator, error) {

	logs, sub, err := _TokenCashierForSolana.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaPauseIterator{contract: _TokenCashierForSolana.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierForSolanaPause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierForSolana.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierForSolanaPause)
				if err := _TokenCashierForSolana.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) ParsePause(log types.Log) (*TokenCashierForSolanaPause, error) {
	event := new(TokenCashierForSolanaPause)
	if err := _TokenCashierForSolana.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierForSolanaReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaReceiptIterator struct {
	Event *TokenCashierForSolanaReceipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierForSolanaReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierForSolanaReceipt)
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
		it.Event = new(TokenCashierForSolanaReceipt)
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
func (it *TokenCashierForSolanaReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierForSolanaReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierForSolanaReceipt represents a Receipt event raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaReceipt struct {
	Token     common.Address
	Id        *big.Int
	Sender    common.Address
	Recipient string
	Amount    *big.Int
	Fee       *big.Int
	Payload   []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0x5c19714fce15effc6b70576855990fdd03e4f73ed5ac526a275b8648c07c89ce.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierForSolanaReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierForSolana.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaReceiptIterator{contract: _TokenCashierForSolana.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0x5c19714fce15effc6b70576855990fdd03e4f73ed5ac526a275b8648c07c89ce.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierForSolanaReceipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierForSolana.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierForSolanaReceipt)
				if err := _TokenCashierForSolana.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0x5c19714fce15effc6b70576855990fdd03e4f73ed5ac526a275b8648c07c89ce.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, string recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) ParseReceipt(log types.Log) (*TokenCashierForSolanaReceipt, error) {
	event := new(TokenCashierForSolanaReceipt)
	if err := _TokenCashierForSolana.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierForSolanaUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaUnpauseIterator struct {
	Event *TokenCashierForSolanaUnpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierForSolanaUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierForSolanaUnpause)
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
		it.Event = new(TokenCashierForSolanaUnpause)
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
func (it *TokenCashierForSolanaUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierForSolanaUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierForSolanaUnpause represents a Unpause event raised by the TokenCashierForSolana contract.
type TokenCashierForSolanaUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierForSolanaUnpauseIterator, error) {

	logs, sub, err := _TokenCashierForSolana.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierForSolanaUnpauseIterator{contract: _TokenCashierForSolana.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierForSolanaUnpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierForSolana.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierForSolanaUnpause)
				if err := _TokenCashierForSolana.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashierForSolana *TokenCashierForSolanaFilterer) ParseUnpause(log types.Log) (*TokenCashierForSolanaUnpause, error) {
	event := new(TokenCashierForSolanaUnpause)
	if err := _TokenCashierForSolana.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
