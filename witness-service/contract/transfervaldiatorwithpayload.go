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

// TransferValidatorWithPayloadMetaData contains all meta data concerning the TransferValidatorWithPayload contract.
var TransferValidatorWithPayloadMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessList\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_tokenList\",\"type\":\"address\"},{\"internalType\":\"contractIMinter\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"addPair\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"addReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getMinter\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"minters\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numOfPairs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"receivers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"removeReceiver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"witnessList\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TransferValidatorWithPayloadABI is the input ABI used to generate the binding from.
// Deprecated: Use TransferValidatorWithPayloadMetaData.ABI instead.
var TransferValidatorWithPayloadABI = TransferValidatorWithPayloadMetaData.ABI

// TransferValidatorWithPayload is an auto generated Go binding around an Ethereum contract.
type TransferValidatorWithPayload struct {
	TransferValidatorWithPayloadCaller     // Read-only binding to the contract
	TransferValidatorWithPayloadTransactor // Write-only binding to the contract
	TransferValidatorWithPayloadFilterer   // Log filterer for contract events
}

// TransferValidatorWithPayloadCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransferValidatorWithPayloadCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorWithPayloadTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferValidatorWithPayloadTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorWithPayloadFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferValidatorWithPayloadFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorWithPayloadSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferValidatorWithPayloadSession struct {
	Contract     *TransferValidatorWithPayload // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// TransferValidatorWithPayloadCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferValidatorWithPayloadCallerSession struct {
	Contract *TransferValidatorWithPayloadCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// TransferValidatorWithPayloadTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferValidatorWithPayloadTransactorSession struct {
	Contract     *TransferValidatorWithPayloadTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// TransferValidatorWithPayloadRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransferValidatorWithPayloadRaw struct {
	Contract *TransferValidatorWithPayload // Generic contract binding to access the raw methods on
}

// TransferValidatorWithPayloadCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferValidatorWithPayloadCallerRaw struct {
	Contract *TransferValidatorWithPayloadCaller // Generic read-only contract binding to access the raw methods on
}

// TransferValidatorWithPayloadTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferValidatorWithPayloadTransactorRaw struct {
	Contract *TransferValidatorWithPayloadTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferValidatorWithPayload creates a new instance of TransferValidatorWithPayload, bound to a specific deployed contract.
func NewTransferValidatorWithPayload(address common.Address, backend bind.ContractBackend) (*TransferValidatorWithPayload, error) {
	contract, err := bindTransferValidatorWithPayload(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayload{TransferValidatorWithPayloadCaller: TransferValidatorWithPayloadCaller{contract: contract}, TransferValidatorWithPayloadTransactor: TransferValidatorWithPayloadTransactor{contract: contract}, TransferValidatorWithPayloadFilterer: TransferValidatorWithPayloadFilterer{contract: contract}}, nil
}

// NewTransferValidatorWithPayloadCaller creates a new read-only instance of TransferValidatorWithPayload, bound to a specific deployed contract.
func NewTransferValidatorWithPayloadCaller(address common.Address, caller bind.ContractCaller) (*TransferValidatorWithPayloadCaller, error) {
	contract, err := bindTransferValidatorWithPayload(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadCaller{contract: contract}, nil
}

// NewTransferValidatorWithPayloadTransactor creates a new write-only instance of TransferValidatorWithPayload, bound to a specific deployed contract.
func NewTransferValidatorWithPayloadTransactor(address common.Address, transactor bind.ContractTransactor) (*TransferValidatorWithPayloadTransactor, error) {
	contract, err := bindTransferValidatorWithPayload(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadTransactor{contract: contract}, nil
}

// NewTransferValidatorWithPayloadFilterer creates a new log filterer instance of TransferValidatorWithPayload, bound to a specific deployed contract.
func NewTransferValidatorWithPayloadFilterer(address common.Address, filterer bind.ContractFilterer) (*TransferValidatorWithPayloadFilterer, error) {
	contract, err := bindTransferValidatorWithPayload(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadFilterer{contract: contract}, nil
}

// bindTransferValidatorWithPayload binds a generic wrapper to an already deployed contract.
func bindTransferValidatorWithPayload(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransferValidatorWithPayloadMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorWithPayload.Contract.TransferValidatorWithPayloadCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.TransferValidatorWithPayloadTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.TransferValidatorWithPayloadTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorWithPayload.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.contract.Transact(opts, method, params...)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3d0dc063.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) GenerateKey(opts *bind.CallOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "generateKey", cashier, tokenAddr, index, from, to, amount, payload)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateKey is a free data retrieval call binding the contract method 0x3d0dc063.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidatorWithPayload.Contract.GenerateKey(&_TransferValidatorWithPayload.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3d0dc063.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidatorWithPayload.Contract.GenerateKey(&_TransferValidatorWithPayload.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) GetMinter(opts *bind.CallOpts, tokenAddr common.Address) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "getMinter", tokenAddr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) GetMinter(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.GetMinter(&_TransferValidatorWithPayload.CallOpts, tokenAddr)
}

// GetMinter is a free data retrieval call binding the contract method 0xbc73b641.
//
// Solidity: function getMinter(address tokenAddr) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) GetMinter(tokenAddr common.Address) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.GetMinter(&_TransferValidatorWithPayload.CallOpts, tokenAddr)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) Minters(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "minters", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.Minters(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.Minters(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) NumOfPairs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "numOfPairs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorWithPayload.Contract.NumOfPairs(&_TransferValidatorWithPayload.CallOpts)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorWithPayload.Contract.NumOfPairs(&_TransferValidatorWithPayload.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Owner() (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.Owner(&_TransferValidatorWithPayload.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) Owner() (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.Owner(&_TransferValidatorWithPayload.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Paused() (bool, error) {
	return _TransferValidatorWithPayload.Contract.Paused(&_TransferValidatorWithPayload.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) Paused() (bool, error) {
	return _TransferValidatorWithPayload.Contract.Paused(&_TransferValidatorWithPayload.CallOpts)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) Receivers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "receivers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorWithPayload.Contract.Receivers(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorWithPayload.Contract.Receivers(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) Settles(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "settles", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorWithPayload.Contract.Settles(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorWithPayload.Contract.Settles(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.TokenLists(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.TokenLists(&_TransferValidatorWithPayload.CallOpts, arg0)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCaller) WitnessList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorWithPayload.contract.Call(opts, &out, "witnessList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) WitnessList() (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.WitnessList(&_TransferValidatorWithPayload.CallOpts)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadCallerSession) WitnessList() (common.Address, error) {
	return _TransferValidatorWithPayload.Contract.WitnessList(&_TransferValidatorWithPayload.CallOpts)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) AddPair(opts *bind.TransactOpts, _tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "addPair", _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.AddPair(&_TransferValidatorWithPayload.TransactOpts, _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.AddPair(&_TransferValidatorWithPayload.TransactOpts, _tokenList, _minter)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) AddReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "addReceiver", _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.AddReceiver(&_TransferValidatorWithPayload.TransactOpts, _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.AddReceiver(&_TransferValidatorWithPayload.TransactOpts, _receiver)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Pause() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Pause(&_TransferValidatorWithPayload.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) Pause() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Pause(&_TransferValidatorWithPayload.TransactOpts)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) RemoveReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "removeReceiver", _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.RemoveReceiver(&_TransferValidatorWithPayload.TransactOpts, _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.RemoveReceiver(&_TransferValidatorWithPayload.TransactOpts, _receiver)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) RenounceOwnership() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.RenounceOwnership(&_TransferValidatorWithPayload.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.RenounceOwnership(&_TransferValidatorWithPayload.TransactOpts)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) Submit(opts *bind.TransactOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Submit(&_TransferValidatorWithPayload.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x73c6d87b.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Submit(&_TransferValidatorWithPayload.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.TransferOwnership(&_TransferValidatorWithPayload.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.TransferOwnership(&_TransferValidatorWithPayload.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Unpause() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Unpause(&_TransferValidatorWithPayload.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) Unpause() (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Unpause(&_TransferValidatorWithPayload.TransactOpts)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactor) Upgrade(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.contract.Transact(opts, "upgrade", _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Upgrade(&_TransferValidatorWithPayload.TransactOpts, _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadTransactorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorWithPayload.Contract.Upgrade(&_TransferValidatorWithPayload.TransactOpts, _newValidator)
}

// TransferValidatorWithPayloadOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadOwnershipTransferredIterator struct {
	Event *TransferValidatorWithPayloadOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadOwnershipTransferred)
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
		it.Event = new(TransferValidatorWithPayloadOwnershipTransferred)
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
func (it *TransferValidatorWithPayloadOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadOwnershipTransferred represents a OwnershipTransferred event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferValidatorWithPayloadOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadOwnershipTransferredIterator{contract: _TransferValidatorWithPayload.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadOwnershipTransferred)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParseOwnershipTransferred(log types.Log) (*TransferValidatorWithPayloadOwnershipTransferred, error) {
	event := new(TransferValidatorWithPayloadOwnershipTransferred)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorWithPayloadPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadPauseIterator struct {
	Event *TransferValidatorWithPayloadPause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadPause)
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
		it.Event = new(TransferValidatorWithPayloadPause)
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
func (it *TransferValidatorWithPayloadPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadPause represents a Pause event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterPause(opts *bind.FilterOpts) (*TransferValidatorWithPayloadPauseIterator, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadPauseIterator{contract: _TransferValidatorWithPayload.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadPause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadPause)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParsePause(log types.Log) (*TransferValidatorWithPayloadPause, error) {
	event := new(TransferValidatorWithPayloadPause)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorWithPayloadReceiverAddedIterator is returned from FilterReceiverAdded and is used to iterate over the raw logs and unpacked data for ReceiverAdded events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadReceiverAddedIterator struct {
	Event *TransferValidatorWithPayloadReceiverAdded // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadReceiverAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadReceiverAdded)
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
		it.Event = new(TransferValidatorWithPayloadReceiverAdded)
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
func (it *TransferValidatorWithPayloadReceiverAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadReceiverAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadReceiverAdded represents a ReceiverAdded event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadReceiverAdded struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverAdded is a free log retrieval operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterReceiverAdded(opts *bind.FilterOpts) (*TransferValidatorWithPayloadReceiverAddedIterator, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadReceiverAddedIterator{contract: _TransferValidatorWithPayload.contract, event: "ReceiverAdded", logs: logs, sub: sub}, nil
}

// WatchReceiverAdded is a free log subscription operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchReceiverAdded(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadReceiverAdded) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadReceiverAdded)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParseReceiverAdded(log types.Log) (*TransferValidatorWithPayloadReceiverAdded, error) {
	event := new(TransferValidatorWithPayloadReceiverAdded)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorWithPayloadReceiverRemovedIterator is returned from FilterReceiverRemoved and is used to iterate over the raw logs and unpacked data for ReceiverRemoved events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadReceiverRemovedIterator struct {
	Event *TransferValidatorWithPayloadReceiverRemoved // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadReceiverRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadReceiverRemoved)
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
		it.Event = new(TransferValidatorWithPayloadReceiverRemoved)
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
func (it *TransferValidatorWithPayloadReceiverRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadReceiverRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadReceiverRemoved represents a ReceiverRemoved event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadReceiverRemoved struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverRemoved is a free log retrieval operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterReceiverRemoved(opts *bind.FilterOpts) (*TransferValidatorWithPayloadReceiverRemovedIterator, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadReceiverRemovedIterator{contract: _TransferValidatorWithPayload.contract, event: "ReceiverRemoved", logs: logs, sub: sub}, nil
}

// WatchReceiverRemoved is a free log subscription operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchReceiverRemoved(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadReceiverRemoved) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadReceiverRemoved)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParseReceiverRemoved(log types.Log) (*TransferValidatorWithPayloadReceiverRemoved, error) {
	event := new(TransferValidatorWithPayloadReceiverRemoved)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorWithPayloadSettledIterator is returned from FilterSettled and is used to iterate over the raw logs and unpacked data for Settled events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadSettledIterator struct {
	Event *TransferValidatorWithPayloadSettled // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadSettled)
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
		it.Event = new(TransferValidatorWithPayloadSettled)
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
func (it *TransferValidatorWithPayloadSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadSettled represents a Settled event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadSettled struct {
	Key       [32]byte
	Witnesses []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterSettled(opts *bind.FilterOpts, key [][32]byte) (*TransferValidatorWithPayloadSettledIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadSettledIterator{contract: _TransferValidatorWithPayload.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadSettled, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadSettled)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Settled", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParseSettled(log types.Log) (*TransferValidatorWithPayloadSettled, error) {
	event := new(TransferValidatorWithPayloadSettled)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorWithPayloadUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadUnpauseIterator struct {
	Event *TransferValidatorWithPayloadUnpause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorWithPayloadUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorWithPayloadUnpause)
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
		it.Event = new(TransferValidatorWithPayloadUnpause)
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
func (it *TransferValidatorWithPayloadUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorWithPayloadUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorWithPayloadUnpause represents a Unpause event raised by the TransferValidatorWithPayload contract.
type TransferValidatorWithPayloadUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) FilterUnpause(opts *bind.FilterOpts) (*TransferValidatorWithPayloadUnpauseIterator, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorWithPayloadUnpauseIterator{contract: _TransferValidatorWithPayload.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TransferValidatorWithPayloadUnpause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorWithPayload.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorWithPayloadUnpause)
				if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TransferValidatorWithPayload *TransferValidatorWithPayloadFilterer) ParseUnpause(log types.Log) (*TransferValidatorWithPayloadUnpause, error) {
	event := new(TransferValidatorWithPayloadUnpause)
	if err := _TransferValidatorWithPayload.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
