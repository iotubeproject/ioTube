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

// TransferValidatorV3MetaData contains all meta data concerning the TransferValidatorV3 contract.
var TransferValidatorV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_minterTokenList\",\"type\":\"address\"},{\"internalType\":\"contractIMinter\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"addMinterPair\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"addReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessTokenList\",\"type\":\"address\"},{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessList\",\"type\":\"address\"}],\"name\":\"addWitnessPair\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getMinter\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getWitnessList\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"minters\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numOfPairs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"receivers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"removeReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"witnessLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"witnessTokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TransferValidatorV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use TransferValidatorV3MetaData.ABI instead.
var TransferValidatorV3ABI = TransferValidatorV3MetaData.ABI

// TransferValidatorV3 is an auto generated Go binding around an Ethereum contract.
type TransferValidatorV3 struct {
	TransferValidatorV3Caller     // Read-only binding to the contract
	TransferValidatorV3Transactor // Write-only binding to the contract
	TransferValidatorV3Filterer   // Log filterer for contract events
}

// TransferValidatorV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type TransferValidatorV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferValidatorV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferValidatorV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferValidatorV3Session struct {
	Contract     *TransferValidatorV3 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TransferValidatorV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferValidatorV3CallerSession struct {
	Contract *TransferValidatorV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TransferValidatorV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferValidatorV3TransactorSession struct {
	Contract     *TransferValidatorV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TransferValidatorV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type TransferValidatorV3Raw struct {
	Contract *TransferValidatorV3 // Generic contract binding to access the raw methods on
}

// TransferValidatorV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferValidatorV3CallerRaw struct {
	Contract *TransferValidatorV3Caller // Generic read-only contract binding to access the raw methods on
}

// TransferValidatorV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferValidatorV3TransactorRaw struct {
	Contract *TransferValidatorV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferValidatorV3 creates a new instance of TransferValidatorV3, bound to a specific deployed contract.
func NewTransferValidatorV3(address common.Address, backend bind.ContractBackend) (*TransferValidatorV3, error) {
	contract, err := bindTransferValidatorV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3{TransferValidatorV3Caller: TransferValidatorV3Caller{contract: contract}, TransferValidatorV3Transactor: TransferValidatorV3Transactor{contract: contract}, TransferValidatorV3Filterer: TransferValidatorV3Filterer{contract: contract}}, nil
}

// NewTransferValidatorV3Caller creates a new read-only instance of TransferValidatorV3, bound to a specific deployed contract.
func NewTransferValidatorV3Caller(address common.Address, caller bind.ContractCaller) (*TransferValidatorV3Caller, error) {
	contract, err := bindTransferValidatorV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3Caller{contract: contract}, nil
}

// NewTransferValidatorV3Transactor creates a new write-only instance of TransferValidatorV3, bound to a specific deployed contract.
func NewTransferValidatorV3Transactor(address common.Address, transactor bind.ContractTransactor) (*TransferValidatorV3Transactor, error) {
	contract, err := bindTransferValidatorV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3Transactor{contract: contract}, nil
}

// NewTransferValidatorV3Filterer creates a new log filterer instance of TransferValidatorV3, bound to a specific deployed contract.
func NewTransferValidatorV3Filterer(address common.Address, filterer bind.ContractFilterer) (*TransferValidatorV3Filterer, error) {
	contract, err := bindTransferValidatorV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3Filterer{contract: contract}, nil
}

// bindTransferValidatorV3 binds a generic wrapper to an already deployed contract.
func bindTransferValidatorV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransferValidatorV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorV3 *TransferValidatorV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorV3.Contract.TransferValidatorV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorV3 *TransferValidatorV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.TransferValidatorV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorV3 *TransferValidatorV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.TransferValidatorV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorV3 *TransferValidatorV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorV3 *TransferValidatorV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorV3 *TransferValidatorV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.contract.Transact(opts, method, params...)
}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) GetMinter(opts *bind.CallOpts, tokenAddr common.Address) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "getMinter", tokenAddr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) GetMinter(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorV3.Contract.GetMinter(&_TransferValidatorV3.CallOpts, tokenAddr)
}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) GetMinter(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorV3.Contract.GetMinter(&_TransferValidatorV3.CallOpts, tokenAddr)
}

// GetWitnessList is a free data retrieval call binding the contract method 0x9d123614.
//
// Solidity: function getWitnessList(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) GetWitnessList(opts *bind.CallOpts, tokenAddr common.Address) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "getWitnessList", tokenAddr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWitnessList is a free data retrieval call binding the contract method 0x9d123614.
//
// Solidity: function getWitnessList(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) GetWitnessList(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorV3.Contract.GetWitnessList(&_TransferValidatorV3.CallOpts, tokenAddr)
}

// GetWitnessList is a free data retrieval call binding the contract method 0x9d123614.
//
// Solidity: function getWitnessList(address tokenAddr) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) GetWitnessList(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorV3.Contract.GetWitnessList(&_TransferValidatorV3.CallOpts, tokenAddr)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) Minters(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "minters", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.Minters(&_TransferValidatorV3.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.Minters(&_TransferValidatorV3.CallOpts, arg0)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3Caller) NumOfPairs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "numOfPairs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3Session) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorV3.Contract.NumOfPairs(&_TransferValidatorV3.CallOpts)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorV3.Contract.NumOfPairs(&_TransferValidatorV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) Owner() (common.Address, error) {
	return _TransferValidatorV3.Contract.Owner(&_TransferValidatorV3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) Owner() (common.Address, error) {
	return _TransferValidatorV3.Contract.Owner(&_TransferValidatorV3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3Session) Paused() (bool, error) {
	return _TransferValidatorV3.Contract.Paused(&_TransferValidatorV3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) Paused() (bool, error) {
	return _TransferValidatorV3.Contract.Paused(&_TransferValidatorV3.CallOpts)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3Caller) Receivers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "receivers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3Session) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorV3.Contract.Receivers(&_TransferValidatorV3.CallOpts, arg0)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorV3.Contract.Receivers(&_TransferValidatorV3.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3Caller) Settles(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "settles", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3Session) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorV3.Contract.Settles(&_TransferValidatorV3.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorV3.Contract.Settles(&_TransferValidatorV3.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.TokenLists(&_TransferValidatorV3.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.TokenLists(&_TransferValidatorV3.CallOpts, arg0)
}

// WitnessLists is a free data retrieval call binding the contract method 0xe3871b8d.
//
// Solidity: function witnessLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) WitnessLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "witnessLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessLists is a free data retrieval call binding the contract method 0xe3871b8d.
//
// Solidity: function witnessLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) WitnessLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.WitnessLists(&_TransferValidatorV3.CallOpts, arg0)
}

// WitnessLists is a free data retrieval call binding the contract method 0xe3871b8d.
//
// Solidity: function witnessLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) WitnessLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.WitnessLists(&_TransferValidatorV3.CallOpts, arg0)
}

// WitnessTokenLists is a free data retrieval call binding the contract method 0x46da8afc.
//
// Solidity: function witnessTokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Caller) WitnessTokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorV3.contract.Call(opts, &out, "witnessTokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessTokenLists is a free data retrieval call binding the contract method 0x46da8afc.
//
// Solidity: function witnessTokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3Session) WitnessTokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.WitnessTokenLists(&_TransferValidatorV3.CallOpts, arg0)
}

// WitnessTokenLists is a free data retrieval call binding the contract method 0x46da8afc.
//
// Solidity: function witnessTokenLists(uint256 ) view returns(address)
func (_TransferValidatorV3 *TransferValidatorV3CallerSession) WitnessTokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorV3.Contract.WitnessTokenLists(&_TransferValidatorV3.CallOpts, arg0)
}

// AddMinterPair is a paid mutator transaction binding the contract method 0x071219be.
//
// Solidity: function addMinterPair(address _minterTokenList, address _minter) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) AddMinterPair(opts *bind.TransactOpts, _minterTokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "addMinterPair", _minterTokenList, _minter)
}

// AddMinterPair is a paid mutator transaction binding the contract method 0x071219be.
//
// Solidity: function addMinterPair(address _minterTokenList, address _minter) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) AddMinterPair(_minterTokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddMinterPair(&_TransferValidatorV3.TransactOpts, _minterTokenList, _minter)
}

// AddMinterPair is a paid mutator transaction binding the contract method 0x071219be.
//
// Solidity: function addMinterPair(address _minterTokenList, address _minter) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) AddMinterPair(_minterTokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddMinterPair(&_TransferValidatorV3.TransactOpts, _minterTokenList, _minter)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) AddReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "addReceiver", _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddReceiver(&_TransferValidatorV3.TransactOpts, _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddReceiver(&_TransferValidatorV3.TransactOpts, _receiver)
}

// AddWitnessPair is a paid mutator transaction binding the contract method 0x8e26c22c.
//
// Solidity: function addWitnessPair(address _witnessTokenList, address _witnessList) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) AddWitnessPair(opts *bind.TransactOpts, _witnessTokenList common.Address, _witnessList common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "addWitnessPair", _witnessTokenList, _witnessList)
}

// AddWitnessPair is a paid mutator transaction binding the contract method 0x8e26c22c.
//
// Solidity: function addWitnessPair(address _witnessTokenList, address _witnessList) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) AddWitnessPair(_witnessTokenList common.Address, _witnessList common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddWitnessPair(&_TransferValidatorV3.TransactOpts, _witnessTokenList, _witnessList)
}

// AddWitnessPair is a paid mutator transaction binding the contract method 0x8e26c22c.
//
// Solidity: function addWitnessPair(address _witnessTokenList, address _witnessList) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) AddWitnessPair(_witnessTokenList common.Address, _witnessList common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.AddWitnessPair(&_TransferValidatorV3.TransactOpts, _witnessTokenList, _witnessList)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) RemoveReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "removeReceiver", _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.RemoveReceiver(&_TransferValidatorV3.TransactOpts, _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.RemoveReceiver(&_TransferValidatorV3.TransactOpts, _receiver)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) RenounceOwnership() (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.RenounceOwnership(&_TransferValidatorV3.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.RenounceOwnership(&_TransferValidatorV3.TransactOpts)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) Submit(opts *bind.TransactOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.Submit(&_TransferValidatorV3.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.Submit(&_TransferValidatorV3.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.TransferOwnership(&_TransferValidatorV3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.TransferOwnership(&_TransferValidatorV3.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV3 *TransferValidatorV3Transactor) Upgrade(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.contract.Transact(opts, "upgrade", _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV3 *TransferValidatorV3Session) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.Upgrade(&_TransferValidatorV3.TransactOpts, _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV3 *TransferValidatorV3TransactorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV3.Contract.Upgrade(&_TransferValidatorV3.TransactOpts, _newValidator)
}

// TransferValidatorV3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferValidatorV3 contract.
type TransferValidatorV3OwnershipTransferredIterator struct {
	Event *TransferValidatorV3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3OwnershipTransferred)
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
		it.Event = new(TransferValidatorV3OwnershipTransferred)
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
func (it *TransferValidatorV3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3OwnershipTransferred represents a OwnershipTransferred event raised by the TransferValidatorV3 contract.
type TransferValidatorV3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferValidatorV3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3OwnershipTransferredIterator{contract: _TransferValidatorV3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3OwnershipTransferred)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParseOwnershipTransferred(log types.Log) (*TransferValidatorV3OwnershipTransferred, error) {
	event := new(TransferValidatorV3OwnershipTransferred)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorV3PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the TransferValidatorV3 contract.
type TransferValidatorV3PausedIterator struct {
	Event *TransferValidatorV3Paused // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3Paused)
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
		it.Event = new(TransferValidatorV3Paused)
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
func (it *TransferValidatorV3PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3Paused represents a Paused event raised by the TransferValidatorV3 contract.
type TransferValidatorV3Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterPaused(opts *bind.FilterOpts) (*TransferValidatorV3PausedIterator, error) {

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3PausedIterator{contract: _TransferValidatorV3.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3Paused) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3Paused)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParsePaused(log types.Log) (*TransferValidatorV3Paused, error) {
	event := new(TransferValidatorV3Paused)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorV3ReceiverAddedIterator is returned from FilterReceiverAdded and is used to iterate over the raw logs and unpacked data for ReceiverAdded events raised by the TransferValidatorV3 contract.
type TransferValidatorV3ReceiverAddedIterator struct {
	Event *TransferValidatorV3ReceiverAdded // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3ReceiverAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3ReceiverAdded)
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
		it.Event = new(TransferValidatorV3ReceiverAdded)
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
func (it *TransferValidatorV3ReceiverAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3ReceiverAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3ReceiverAdded represents a ReceiverAdded event raised by the TransferValidatorV3 contract.
type TransferValidatorV3ReceiverAdded struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverAdded is a free log retrieval operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterReceiverAdded(opts *bind.FilterOpts) (*TransferValidatorV3ReceiverAddedIterator, error) {

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3ReceiverAddedIterator{contract: _TransferValidatorV3.contract, event: "ReceiverAdded", logs: logs, sub: sub}, nil
}

// WatchReceiverAdded is a free log subscription operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchReceiverAdded(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3ReceiverAdded) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3ReceiverAdded)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
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

// ParseReceiverAdded is a log parse operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParseReceiverAdded(log types.Log) (*TransferValidatorV3ReceiverAdded, error) {
	event := new(TransferValidatorV3ReceiverAdded)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorV3ReceiverRemovedIterator is returned from FilterReceiverRemoved and is used to iterate over the raw logs and unpacked data for ReceiverRemoved events raised by the TransferValidatorV3 contract.
type TransferValidatorV3ReceiverRemovedIterator struct {
	Event *TransferValidatorV3ReceiverRemoved // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3ReceiverRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3ReceiverRemoved)
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
		it.Event = new(TransferValidatorV3ReceiverRemoved)
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
func (it *TransferValidatorV3ReceiverRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3ReceiverRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3ReceiverRemoved represents a ReceiverRemoved event raised by the TransferValidatorV3 contract.
type TransferValidatorV3ReceiverRemoved struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverRemoved is a free log retrieval operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterReceiverRemoved(opts *bind.FilterOpts) (*TransferValidatorV3ReceiverRemovedIterator, error) {

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3ReceiverRemovedIterator{contract: _TransferValidatorV3.contract, event: "ReceiverRemoved", logs: logs, sub: sub}, nil
}

// WatchReceiverRemoved is a free log subscription operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchReceiverRemoved(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3ReceiverRemoved) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3ReceiverRemoved)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
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

// ParseReceiverRemoved is a log parse operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParseReceiverRemoved(log types.Log) (*TransferValidatorV3ReceiverRemoved, error) {
	event := new(TransferValidatorV3ReceiverRemoved)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorV3SettledIterator is returned from FilterSettled and is used to iterate over the raw logs and unpacked data for Settled events raised by the TransferValidatorV3 contract.
type TransferValidatorV3SettledIterator struct {
	Event *TransferValidatorV3Settled // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3SettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3Settled)
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
		it.Event = new(TransferValidatorV3Settled)
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
func (it *TransferValidatorV3SettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3SettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3Settled represents a Settled event raised by the TransferValidatorV3 contract.
type TransferValidatorV3Settled struct {
	Key       [32]byte
	Witnesses []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterSettled(opts *bind.FilterOpts, key [][32]byte) (*TransferValidatorV3SettledIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3SettledIterator{contract: _TransferValidatorV3.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3Settled, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3Settled)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "Settled", log); err != nil {
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

// ParseSettled is a log parse operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParseSettled(log types.Log) (*TransferValidatorV3Settled, error) {
	event := new(TransferValidatorV3Settled)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorV3UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the TransferValidatorV3 contract.
type TransferValidatorV3UnpausedIterator struct {
	Event *TransferValidatorV3Unpaused // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV3UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV3Unpaused)
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
		it.Event = new(TransferValidatorV3Unpaused)
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
func (it *TransferValidatorV3UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV3UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV3Unpaused represents a Unpaused event raised by the TransferValidatorV3 contract.
type TransferValidatorV3Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) FilterUnpaused(opts *bind.FilterOpts) (*TransferValidatorV3UnpausedIterator, error) {

	logs, sub, err := _TransferValidatorV3.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV3UnpausedIterator{contract: _TransferValidatorV3.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *TransferValidatorV3Unpaused) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV3.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV3Unpaused)
				if err := _TransferValidatorV3.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_TransferValidatorV3 *TransferValidatorV3Filterer) ParseUnpaused(log types.Log) (*TransferValidatorV3Unpaused, error) {
	event := new(TransferValidatorV3Unpaused)
	if err := _TransferValidatorV3.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
