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

// TransferValidatorV2ABI is the input ABI used to generate the binding from.
const TransferValidatorV2ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"whitelistedTokens\",\"outputs\":[{\"internalType\":\"contractAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"whitelistedWitnesses\",\"outputs\":[{\"internalType\":\"contractAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TransferValidatorV2 is an auto generated Go binding around an Ethereum contract.
type TransferValidatorV2 struct {
	TransferValidatorV2Caller     // Read-only binding to the contract
	TransferValidatorV2Transactor // Write-only binding to the contract
	TransferValidatorV2Filterer   // Log filterer for contract events
}

// TransferValidatorV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type TransferValidatorV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferValidatorV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferValidatorV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferValidatorV2Session struct {
	Contract     *TransferValidatorV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TransferValidatorV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferValidatorV2CallerSession struct {
	Contract *TransferValidatorV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TransferValidatorV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferValidatorV2TransactorSession struct {
	Contract     *TransferValidatorV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TransferValidatorV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type TransferValidatorV2Raw struct {
	Contract *TransferValidatorV2 // Generic contract binding to access the raw methods on
}

// TransferValidatorV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferValidatorV2CallerRaw struct {
	Contract *TransferValidatorV2Caller // Generic read-only contract binding to access the raw methods on
}

// TransferValidatorV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferValidatorV2TransactorRaw struct {
	Contract *TransferValidatorV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferValidatorV2 creates a new instance of TransferValidatorV2, bound to a specific deployed contract.
func NewTransferValidatorV2(address common.Address, backend bind.ContractBackend) (*TransferValidatorV2, error) {
	contract, err := bindTransferValidatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2{TransferValidatorV2Caller: TransferValidatorV2Caller{contract: contract}, TransferValidatorV2Transactor: TransferValidatorV2Transactor{contract: contract}, TransferValidatorV2Filterer: TransferValidatorV2Filterer{contract: contract}}, nil
}

// NewTransferValidatorV2Caller creates a new read-only instance of TransferValidatorV2, bound to a specific deployed contract.
func NewTransferValidatorV2Caller(address common.Address, caller bind.ContractCaller) (*TransferValidatorV2Caller, error) {
	contract, err := bindTransferValidatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2Caller{contract: contract}, nil
}

// NewTransferValidatorV2Transactor creates a new write-only instance of TransferValidatorV2, bound to a specific deployed contract.
func NewTransferValidatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*TransferValidatorV2Transactor, error) {
	contract, err := bindTransferValidatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2Transactor{contract: contract}, nil
}

// NewTransferValidatorV2Filterer creates a new log filterer instance of TransferValidatorV2, bound to a specific deployed contract.
func NewTransferValidatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*TransferValidatorV2Filterer, error) {
	contract, err := bindTransferValidatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2Filterer{contract: contract}, nil
}

// bindTransferValidatorV2 binds a generic wrapper to an already deployed contract.
func bindTransferValidatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TransferValidatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorV2 *TransferValidatorV2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TransferValidatorV2.Contract.TransferValidatorV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorV2 *TransferValidatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.TransferValidatorV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorV2 *TransferValidatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.TransferValidatorV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorV2 *TransferValidatorV2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TransferValidatorV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorV2 *TransferValidatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorV2 *TransferValidatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.contract.Transact(opts, method, params...)
}

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidatorV2 *TransferValidatorV2Caller) GenerateKey(opts *bind.CallOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "generateKey", cashier, tokenAddr, index, from, to, amount)
	return *ret0, err
}

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidatorV2 *TransferValidatorV2Session) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidatorV2.Contract.GenerateKey(&_TransferValidatorV2.CallOpts, cashier, tokenAddr, index, from, to, amount)
}

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidatorV2.Contract.GenerateKey(&_TransferValidatorV2.CallOpts, cashier, tokenAddr, index, from, to, amount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Session) Owner() (common.Address, error) {
	return _TransferValidatorV2.Contract.Owner(&_TransferValidatorV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) Owner() (common.Address, error) {
	return _TransferValidatorV2.Contract.Owner(&_TransferValidatorV2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV2 *TransferValidatorV2Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV2 *TransferValidatorV2Session) Paused() (bool, error) {
	return _TransferValidatorV2.Contract.Paused(&_TransferValidatorV2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) Paused() (bool, error) {
	return _TransferValidatorV2.Contract.Paused(&_TransferValidatorV2.CallOpts)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV2 *TransferValidatorV2Caller) Settles(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "settles", arg0)
	return *ret0, err
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV2 *TransferValidatorV2Session) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorV2.Contract.Settles(&_TransferValidatorV2.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorV2.Contract.Settles(&_TransferValidatorV2.CallOpts, arg0)
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Caller) WhitelistedTokens(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "whitelistedTokens")
	return *ret0, err
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Session) WhitelistedTokens() (common.Address, error) {
	return _TransferValidatorV2.Contract.WhitelistedTokens(&_TransferValidatorV2.CallOpts)
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) WhitelistedTokens() (common.Address, error) {
	return _TransferValidatorV2.Contract.WhitelistedTokens(&_TransferValidatorV2.CallOpts)
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Caller) WhitelistedWitnesses(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidatorV2.contract.Call(opts, out, "whitelistedWitnesses")
	return *ret0, err
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2Session) WhitelistedWitnesses() (common.Address, error) {
	return _TransferValidatorV2.Contract.WhitelistedWitnesses(&_TransferValidatorV2.CallOpts)
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidatorV2 *TransferValidatorV2CallerSession) WhitelistedWitnesses() (common.Address, error) {
	return _TransferValidatorV2.Contract.WhitelistedWitnesses(&_TransferValidatorV2.CallOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorV2 *TransferValidatorV2Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV2.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorV2 *TransferValidatorV2Session) Pause() (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Pause(&_TransferValidatorV2.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorV2 *TransferValidatorV2TransactorSession) Pause() (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Pause(&_TransferValidatorV2.TransactOpts)
}

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidatorV2 *TransferValidatorV2Transactor) Submit(opts *bind.TransactOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidatorV2.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures)
}

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidatorV2 *TransferValidatorV2Session) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Submit(&_TransferValidatorV2.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures)
}

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidatorV2 *TransferValidatorV2TransactorSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Submit(&_TransferValidatorV2.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV2 *TransferValidatorV2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV2 *TransferValidatorV2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.TransferOwnership(&_TransferValidatorV2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorV2 *TransferValidatorV2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.TransferOwnership(&_TransferValidatorV2.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorV2 *TransferValidatorV2Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorV2.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorV2 *TransferValidatorV2Session) Unpause() (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Unpause(&_TransferValidatorV2.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorV2 *TransferValidatorV2TransactorSession) Unpause() (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Unpause(&_TransferValidatorV2.TransactOpts)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV2 *TransferValidatorV2Transactor) Upgrade(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.contract.Transact(opts, "upgrade", _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV2 *TransferValidatorV2Session) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Upgrade(&_TransferValidatorV2.TransactOpts, _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorV2 *TransferValidatorV2TransactorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorV2.Contract.Upgrade(&_TransferValidatorV2.TransactOpts, _newValidator)
}

// TransferValidatorV2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferValidatorV2 contract.
type TransferValidatorV2OwnershipTransferredIterator struct {
	Event *TransferValidatorV2OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV2OwnershipTransferred)
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
		it.Event = new(TransferValidatorV2OwnershipTransferred)
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
func (it *TransferValidatorV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV2OwnershipTransferred represents a OwnershipTransferred event raised by the TransferValidatorV2 contract.
type TransferValidatorV2OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorV2 *TransferValidatorV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferValidatorV2OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorV2.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2OwnershipTransferredIterator{contract: _TransferValidatorV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorV2 *TransferValidatorV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferValidatorV2OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorV2.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV2OwnershipTransferred)
				if err := _TransferValidatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferValidatorV2 *TransferValidatorV2Filterer) ParseOwnershipTransferred(log types.Log) (*TransferValidatorV2OwnershipTransferred, error) {
	event := new(TransferValidatorV2OwnershipTransferred)
	if err := _TransferValidatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TransferValidatorV2PauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TransferValidatorV2 contract.
type TransferValidatorV2PauseIterator struct {
	Event *TransferValidatorV2Pause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV2PauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV2Pause)
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
		it.Event = new(TransferValidatorV2Pause)
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
func (it *TransferValidatorV2PauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV2PauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV2Pause represents a Pause event raised by the TransferValidatorV2 contract.
type TransferValidatorV2Pause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorV2 *TransferValidatorV2Filterer) FilterPause(opts *bind.FilterOpts) (*TransferValidatorV2PauseIterator, error) {

	logs, sub, err := _TransferValidatorV2.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2PauseIterator{contract: _TransferValidatorV2.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorV2 *TransferValidatorV2Filterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TransferValidatorV2Pause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV2.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV2Pause)
				if err := _TransferValidatorV2.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TransferValidatorV2 *TransferValidatorV2Filterer) ParsePause(log types.Log) (*TransferValidatorV2Pause, error) {
	event := new(TransferValidatorV2Pause)
	if err := _TransferValidatorV2.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TransferValidatorV2SettledIterator is returned from FilterSettled and is used to iterate over the raw logs and unpacked data for Settled events raised by the TransferValidatorV2 contract.
type TransferValidatorV2SettledIterator struct {
	Event *TransferValidatorV2Settled // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV2SettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV2Settled)
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
		it.Event = new(TransferValidatorV2Settled)
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
func (it *TransferValidatorV2SettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV2SettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV2Settled represents a Settled event raised by the TransferValidatorV2 contract.
type TransferValidatorV2Settled struct {
	Key       [32]byte
	Witnesses []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorV2 *TransferValidatorV2Filterer) FilterSettled(opts *bind.FilterOpts, key [][32]byte) (*TransferValidatorV2SettledIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorV2.contract.FilterLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2SettledIterator{contract: _TransferValidatorV2.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorV2 *TransferValidatorV2Filterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorV2Settled, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorV2.contract.WatchLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV2Settled)
				if err := _TransferValidatorV2.contract.UnpackLog(event, "Settled", log); err != nil {
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
func (_TransferValidatorV2 *TransferValidatorV2Filterer) ParseSettled(log types.Log) (*TransferValidatorV2Settled, error) {
	event := new(TransferValidatorV2Settled)
	if err := _TransferValidatorV2.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TransferValidatorV2UnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TransferValidatorV2 contract.
type TransferValidatorV2UnpauseIterator struct {
	Event *TransferValidatorV2Unpause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorV2UnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorV2Unpause)
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
		it.Event = new(TransferValidatorV2Unpause)
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
func (it *TransferValidatorV2UnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorV2UnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorV2Unpause represents a Unpause event raised by the TransferValidatorV2 contract.
type TransferValidatorV2Unpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorV2 *TransferValidatorV2Filterer) FilterUnpause(opts *bind.FilterOpts) (*TransferValidatorV2UnpauseIterator, error) {

	logs, sub, err := _TransferValidatorV2.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorV2UnpauseIterator{contract: _TransferValidatorV2.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorV2 *TransferValidatorV2Filterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TransferValidatorV2Unpause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorV2.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorV2Unpause)
				if err := _TransferValidatorV2.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TransferValidatorV2 *TransferValidatorV2Filterer) ParseUnpause(log types.Log) (*TransferValidatorV2Unpause, error) {
	event := new(TransferValidatorV2Unpause)
	if err := _TransferValidatorV2.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	return event, nil
}
