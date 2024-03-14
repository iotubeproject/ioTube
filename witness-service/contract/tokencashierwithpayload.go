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

// TokenCashierWithPayloadMetaData contains all meta data concerning the TokenCashierWithPayload contract.
var TokenCashierWithPayloadMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"_wrappedCoin\",\"type\":\"address\"},{\"internalType\":\"contractITokenList[]\",\"name\":\"_tokenLists\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_tokenSafes\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"counts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"depositTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getSafeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractITokenList\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenSafes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wrappedCoin\",\"outputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TokenCashierWithPayloadABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenCashierWithPayloadMetaData.ABI instead.
var TokenCashierWithPayloadABI = TokenCashierWithPayloadMetaData.ABI

// TokenCashierWithPayload is an auto generated Go binding around an Ethereum contract.
type TokenCashierWithPayload struct {
	TokenCashierWithPayloadCaller     // Read-only binding to the contract
	TokenCashierWithPayloadTransactor // Write-only binding to the contract
	TokenCashierWithPayloadFilterer   // Log filterer for contract events
}

// TokenCashierWithPayloadCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierWithPayloadCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierWithPayloadTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierWithPayloadTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierWithPayloadFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierWithPayloadFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierWithPayloadSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierWithPayloadSession struct {
	Contract     *TokenCashierWithPayload // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TokenCashierWithPayloadCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierWithPayloadCallerSession struct {
	Contract *TokenCashierWithPayloadCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TokenCashierWithPayloadTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierWithPayloadTransactorSession struct {
	Contract     *TokenCashierWithPayloadTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TokenCashierWithPayloadRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierWithPayloadRaw struct {
	Contract *TokenCashierWithPayload // Generic contract binding to access the raw methods on
}

// TokenCashierWithPayloadCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierWithPayloadCallerRaw struct {
	Contract *TokenCashierWithPayloadCaller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierWithPayloadTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierWithPayloadTransactorRaw struct {
	Contract *TokenCashierWithPayloadTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashierWithPayload creates a new instance of TokenCashierWithPayload, bound to a specific deployed contract.
func NewTokenCashierWithPayload(address common.Address, backend bind.ContractBackend) (*TokenCashierWithPayload, error) {
	contract, err := bindTokenCashierWithPayload(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayload{TokenCashierWithPayloadCaller: TokenCashierWithPayloadCaller{contract: contract}, TokenCashierWithPayloadTransactor: TokenCashierWithPayloadTransactor{contract: contract}, TokenCashierWithPayloadFilterer: TokenCashierWithPayloadFilterer{contract: contract}}, nil
}

// NewTokenCashierWithPayloadCaller creates a new read-only instance of TokenCashierWithPayload, bound to a specific deployed contract.
func NewTokenCashierWithPayloadCaller(address common.Address, caller bind.ContractCaller) (*TokenCashierWithPayloadCaller, error) {
	contract, err := bindTokenCashierWithPayload(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadCaller{contract: contract}, nil
}

// NewTokenCashierWithPayloadTransactor creates a new write-only instance of TokenCashierWithPayload, bound to a specific deployed contract.
func NewTokenCashierWithPayloadTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierWithPayloadTransactor, error) {
	contract, err := bindTokenCashierWithPayload(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadTransactor{contract: contract}, nil
}

// NewTokenCashierWithPayloadFilterer creates a new log filterer instance of TokenCashierWithPayload, bound to a specific deployed contract.
func NewTokenCashierWithPayloadFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierWithPayloadFilterer, error) {
	contract, err := bindTokenCashierWithPayload(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadFilterer{contract: contract}, nil
}

// bindTokenCashierWithPayload binds a generic wrapper to an already deployed contract.
func bindTokenCashierWithPayload(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenCashierWithPayloadMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierWithPayload *TokenCashierWithPayloadRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierWithPayload.Contract.TokenCashierWithPayloadCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierWithPayload *TokenCashierWithPayloadRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.TokenCashierWithPayloadTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierWithPayload *TokenCashierWithPayloadRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.TokenCashierWithPayloadTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashierWithPayload.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "count", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.Count(&_TokenCashierWithPayload.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.Count(&_TokenCashierWithPayload.CallOpts, _token)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) Counts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "counts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.Counts(&_TokenCashierWithPayload.CallOpts, arg0)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.Counts(&_TokenCashierWithPayload.CallOpts, arg0)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "depositFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) DepositFee() (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.DepositFee(&_TokenCashierWithPayload.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashierWithPayload.Contract.DepositFee(&_TokenCashierWithPayload.CallOpts)
}

// GetSafeAddress is a free data retrieval call binding the contract method 0xa287bdf1.
//
// Solidity: function getSafeAddress(address _token) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) GetSafeAddress(opts *bind.CallOpts, _token common.Address) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "getSafeAddress", _token)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSafeAddress is a free data retrieval call binding the contract method 0xa287bdf1.
//
// Solidity: function getSafeAddress(address _token) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) GetSafeAddress(_token common.Address) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.GetSafeAddress(&_TokenCashierWithPayload.CallOpts, _token)
}

// GetSafeAddress is a free data retrieval call binding the contract method 0xa287bdf1.
//
// Solidity: function getSafeAddress(address _token) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) GetSafeAddress(_token common.Address) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.GetSafeAddress(&_TokenCashierWithPayload.CallOpts, _token)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Owner() (common.Address, error) {
	return _TokenCashierWithPayload.Contract.Owner(&_TokenCashierWithPayload.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) Owner() (common.Address, error) {
	return _TokenCashierWithPayload.Contract.Owner(&_TokenCashierWithPayload.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Paused() (bool, error) {
	return _TokenCashierWithPayload.Contract.Paused(&_TokenCashierWithPayload.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) Paused() (bool, error) {
	return _TokenCashierWithPayload.Contract.Paused(&_TokenCashierWithPayload.CallOpts)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.TokenLists(&_TokenCashierWithPayload.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.TokenLists(&_TokenCashierWithPayload.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) TokenSafes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "tokenSafes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.TokenSafes(&_TokenCashierWithPayload.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashierWithPayload.Contract.TokenSafes(&_TokenCashierWithPayload.CallOpts, arg0)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCaller) WrappedCoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashierWithPayload.contract.Call(opts, &out, "wrappedCoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) WrappedCoin() (common.Address, error) {
	return _TokenCashierWithPayload.Contract.WrappedCoin(&_TokenCashierWithPayload.CallOpts)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashierWithPayload *TokenCashierWithPayloadCallerSession) WrappedCoin() (common.Address, error) {
	return _TokenCashierWithPayload.Contract.WrappedCoin(&_TokenCashierWithPayload.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) Deposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "deposit", _token, _amount, _payload)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Deposit(_token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Deposit(&_TokenCashierWithPayload.TransactOpts, _token, _amount, _payload)
}

// Deposit is a paid mutator transaction binding the contract method 0x49bdc2b8.
//
// Solidity: function deposit(address _token, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) Deposit(_token common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Deposit(&_TokenCashierWithPayload.TransactOpts, _token, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "depositTo", _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.DepositTo(&_TokenCashierWithPayload.TransactOpts, _token, _to, _amount, _payload)
}

// DepositTo is a paid mutator transaction binding the contract method 0xd33b5bb9.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount, bytes _payload) payable returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int, _payload []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.DepositTo(&_TokenCashierWithPayload.TransactOpts, _token, _to, _amount, _payload)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Pause() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Pause(&_TokenCashierWithPayload.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Pause(&_TokenCashierWithPayload.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.RenounceOwnership(&_TokenCashierWithPayload.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.RenounceOwnership(&_TokenCashierWithPayload.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.SetDepositFee(&_TokenCashierWithPayload.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.SetDepositFee(&_TokenCashierWithPayload.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.TransferOwnership(&_TokenCashierWithPayload.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.TransferOwnership(&_TokenCashierWithPayload.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Unpause(&_TokenCashierWithPayload.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Unpause(&_TokenCashierWithPayload.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Withdraw(&_TokenCashierWithPayload.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Withdraw(&_TokenCashierWithPayload.TransactOpts)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.Transact(opts, "withdrawToken", _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.WithdrawToken(&_TokenCashierWithPayload.TransactOpts, _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.WithdrawToken(&_TokenCashierWithPayload.TransactOpts, _token)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Fallback(&_TokenCashierWithPayload.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashierWithPayload *TokenCashierWithPayloadTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashierWithPayload.Contract.Fallback(&_TokenCashierWithPayload.TransactOpts, calldata)
}

// TokenCashierWithPayloadOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadOwnershipTransferredIterator struct {
	Event *TokenCashierWithPayloadOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierWithPayloadOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierWithPayloadOwnershipTransferred)
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
		it.Event = new(TokenCashierWithPayloadOwnershipTransferred)
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
func (it *TokenCashierWithPayloadOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierWithPayloadOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierWithPayloadOwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierWithPayloadOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierWithPayload.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadOwnershipTransferredIterator{contract: _TokenCashierWithPayload.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierWithPayloadOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashierWithPayload.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierWithPayloadOwnershipTransferred)
				if err := _TokenCashierWithPayload.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierWithPayloadOwnershipTransferred, error) {
	event := new(TokenCashierWithPayloadOwnershipTransferred)
	if err := _TokenCashierWithPayload.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierWithPayloadPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadPauseIterator struct {
	Event *TokenCashierWithPayloadPause // Event containing the contract specifics and raw log

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
func (it *TokenCashierWithPayloadPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierWithPayloadPause)
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
		it.Event = new(TokenCashierWithPayloadPause)
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
func (it *TokenCashierWithPayloadPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierWithPayloadPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierWithPayloadPause represents a Pause event raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierWithPayloadPauseIterator, error) {

	logs, sub, err := _TokenCashierWithPayload.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadPauseIterator{contract: _TokenCashierWithPayload.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierWithPayloadPause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierWithPayload.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierWithPayloadPause)
				if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) ParsePause(log types.Log) (*TokenCashierWithPayloadPause, error) {
	event := new(TokenCashierWithPayloadPause)
	if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierWithPayloadReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadReceiptIterator struct {
	Event *TokenCashierWithPayloadReceipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierWithPayloadReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierWithPayloadReceipt)
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
		it.Event = new(TokenCashierWithPayloadReceipt)
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
func (it *TokenCashierWithPayloadReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierWithPayloadReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierWithPayloadReceipt represents a Receipt event raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadReceipt struct {
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
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierWithPayloadReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierWithPayload.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadReceiptIterator{contract: _TokenCashierWithPayload.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0xd2be25887579d6d0dc43743403c85c398b3873c57506ad20610cef12f2a3c9d2.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee, bytes payload)
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierWithPayloadReceipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashierWithPayload.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierWithPayloadReceipt)
				if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Receipt", log); err != nil {
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
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) ParseReceipt(log types.Log) (*TokenCashierWithPayloadReceipt, error) {
	event := new(TokenCashierWithPayloadReceipt)
	if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierWithPayloadUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadUnpauseIterator struct {
	Event *TokenCashierWithPayloadUnpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierWithPayloadUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierWithPayloadUnpause)
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
		it.Event = new(TokenCashierWithPayloadUnpause)
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
func (it *TokenCashierWithPayloadUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierWithPayloadUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierWithPayloadUnpause represents a Unpause event raised by the TokenCashierWithPayload contract.
type TokenCashierWithPayloadUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierWithPayloadUnpauseIterator, error) {

	logs, sub, err := _TokenCashierWithPayload.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierWithPayloadUnpauseIterator{contract: _TokenCashierWithPayload.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierWithPayloadUnpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashierWithPayload.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierWithPayloadUnpause)
				if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashierWithPayload *TokenCashierWithPayloadFilterer) ParseUnpause(log types.Log) (*TokenCashierWithPayloadUnpause, error) {
	event := new(TokenCashierWithPayloadUnpause)
	if err := _TokenCashierWithPayload.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
