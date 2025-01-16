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

// TransferValidatorForSolanaMetaData contains all meta data concerning the TransferValidatorForSolana contract.
var TransferValidatorForSolanaMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessList\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ReceiverRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_tokenList\",\"type\":\"address\"},{\"internalType\":\"contractIMinter\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"addPair\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"addReceiver\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"extractWitnesses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"cashier\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getTokenGroup\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"minters\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numOfPairs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"receivers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"}],\"name\":\"removeReceiver\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"cashier\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"witnessList\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TransferValidatorForSolanaABI is the input ABI used to generate the binding from.
// Deprecated: Use TransferValidatorForSolanaMetaData.ABI instead.
var TransferValidatorForSolanaABI = TransferValidatorForSolanaMetaData.ABI

// TransferValidatorForSolana is an auto generated Go binding around an Ethereum contract.
type TransferValidatorForSolana struct {
	TransferValidatorForSolanaCaller     // Read-only binding to the contract
	TransferValidatorForSolanaTransactor // Write-only binding to the contract
	TransferValidatorForSolanaFilterer   // Log filterer for contract events
}

// TransferValidatorForSolanaCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransferValidatorForSolanaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorForSolanaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferValidatorForSolanaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorForSolanaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferValidatorForSolanaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorForSolanaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferValidatorForSolanaSession struct {
	Contract     *TransferValidatorForSolana // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TransferValidatorForSolanaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferValidatorForSolanaCallerSession struct {
	Contract *TransferValidatorForSolanaCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// TransferValidatorForSolanaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferValidatorForSolanaTransactorSession struct {
	Contract     *TransferValidatorForSolanaTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// TransferValidatorForSolanaRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransferValidatorForSolanaRaw struct {
	Contract *TransferValidatorForSolana // Generic contract binding to access the raw methods on
}

// TransferValidatorForSolanaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferValidatorForSolanaCallerRaw struct {
	Contract *TransferValidatorForSolanaCaller // Generic read-only contract binding to access the raw methods on
}

// TransferValidatorForSolanaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferValidatorForSolanaTransactorRaw struct {
	Contract *TransferValidatorForSolanaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferValidatorForSolana creates a new instance of TransferValidatorForSolana, bound to a specific deployed contract.
func NewTransferValidatorForSolana(address common.Address, backend bind.ContractBackend) (*TransferValidatorForSolana, error) {
	contract, err := bindTransferValidatorForSolana(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolana{TransferValidatorForSolanaCaller: TransferValidatorForSolanaCaller{contract: contract}, TransferValidatorForSolanaTransactor: TransferValidatorForSolanaTransactor{contract: contract}, TransferValidatorForSolanaFilterer: TransferValidatorForSolanaFilterer{contract: contract}}, nil
}

// NewTransferValidatorForSolanaCaller creates a new read-only instance of TransferValidatorForSolana, bound to a specific deployed contract.
func NewTransferValidatorForSolanaCaller(address common.Address, caller bind.ContractCaller) (*TransferValidatorForSolanaCaller, error) {
	contract, err := bindTransferValidatorForSolana(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaCaller{contract: contract}, nil
}

// NewTransferValidatorForSolanaTransactor creates a new write-only instance of TransferValidatorForSolana, bound to a specific deployed contract.
func NewTransferValidatorForSolanaTransactor(address common.Address, transactor bind.ContractTransactor) (*TransferValidatorForSolanaTransactor, error) {
	contract, err := bindTransferValidatorForSolana(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaTransactor{contract: contract}, nil
}

// NewTransferValidatorForSolanaFilterer creates a new log filterer instance of TransferValidatorForSolana, bound to a specific deployed contract.
func NewTransferValidatorForSolanaFilterer(address common.Address, filterer bind.ContractFilterer) (*TransferValidatorForSolanaFilterer, error) {
	contract, err := bindTransferValidatorForSolana(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaFilterer{contract: contract}, nil
}

// bindTransferValidatorForSolana binds a generic wrapper to an already deployed contract.
func bindTransferValidatorForSolana(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransferValidatorForSolanaMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorForSolana *TransferValidatorForSolanaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorForSolana.Contract.TransferValidatorForSolanaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorForSolana *TransferValidatorForSolanaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.TransferValidatorForSolanaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorForSolana *TransferValidatorForSolanaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.TransferValidatorForSolanaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidatorForSolana.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.contract.Transact(opts, method, params...)
}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) ExtractWitnesses(opts *bind.CallOpts, key [32]byte, signatures []byte) ([]common.Address, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "extractWitnesses", key, signatures)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) ExtractWitnesses(key [32]byte, signatures []byte) ([]common.Address, error) {
	return _TransferValidatorForSolana.Contract.ExtractWitnesses(&_TransferValidatorForSolana.CallOpts, key, signatures)
}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) ExtractWitnesses(key [32]byte, signatures []byte) ([]common.Address, error) {
	return _TransferValidatorForSolana.Contract.ExtractWitnesses(&_TransferValidatorForSolana.CallOpts, key, signatures)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) GenerateKey(opts *bind.CallOpts, cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "generateKey", cashier, tokenAddr, index, from, to, amount, payload)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) GenerateKey(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidatorForSolana.Contract.GenerateKey(&_TransferValidatorForSolana.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) GenerateKey(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidatorForSolana.Contract.GenerateKey(&_TransferValidatorForSolana.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) GetTokenGroup(opts *bind.CallOpts, tokenAddr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "getTokenGroup", tokenAddr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) GetTokenGroup(tokenAddr common.Address) (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.GetTokenGroup(&_TransferValidatorForSolana.CallOpts, tokenAddr)
}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) GetTokenGroup(tokenAddr common.Address) (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.GetTokenGroup(&_TransferValidatorForSolana.CallOpts, tokenAddr)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) Minters(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "minters", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorForSolana.Contract.Minters(&_TransferValidatorForSolana.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorForSolana.Contract.Minters(&_TransferValidatorForSolana.CallOpts, arg0)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) NumOfPairs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "numOfPairs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.NumOfPairs(&_TransferValidatorForSolana.CallOpts)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.NumOfPairs(&_TransferValidatorForSolana.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Owner() (common.Address, error) {
	return _TransferValidatorForSolana.Contract.Owner(&_TransferValidatorForSolana.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) Owner() (common.Address, error) {
	return _TransferValidatorForSolana.Contract.Owner(&_TransferValidatorForSolana.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Paused() (bool, error) {
	return _TransferValidatorForSolana.Contract.Paused(&_TransferValidatorForSolana.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) Paused() (bool, error) {
	return _TransferValidatorForSolana.Contract.Paused(&_TransferValidatorForSolana.CallOpts)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) Receivers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "receivers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorForSolana.Contract.Receivers(&_TransferValidatorForSolana.CallOpts, arg0)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers(address ) view returns(bool)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) Receivers(arg0 common.Address) (bool, error) {
	return _TransferValidatorForSolana.Contract.Receivers(&_TransferValidatorForSolana.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) Settles(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "settles", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.Settles(&_TransferValidatorForSolana.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidatorForSolana.Contract.Settles(&_TransferValidatorForSolana.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorForSolana.Contract.TokenLists(&_TransferValidatorForSolana.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidatorForSolana.Contract.TokenLists(&_TransferValidatorForSolana.CallOpts, arg0)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCaller) WitnessList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidatorForSolana.contract.Call(opts, &out, "witnessList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) WitnessList() (common.Address, error) {
	return _TransferValidatorForSolana.Contract.WitnessList(&_TransferValidatorForSolana.CallOpts)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidatorForSolana *TransferValidatorForSolanaCallerSession) WitnessList() (common.Address, error) {
	return _TransferValidatorForSolana.Contract.WitnessList(&_TransferValidatorForSolana.CallOpts)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) AddPair(opts *bind.TransactOpts, _tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "addPair", _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.AddPair(&_TransferValidatorForSolana.TransactOpts, _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.AddPair(&_TransferValidatorForSolana.TransactOpts, _tokenList, _minter)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) AddReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "addReceiver", _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.AddReceiver(&_TransferValidatorForSolana.TransactOpts, _receiver)
}

// AddReceiver is a paid mutator transaction binding the contract method 0x69d83ed1.
//
// Solidity: function addReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) AddReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.AddReceiver(&_TransferValidatorForSolana.TransactOpts, _receiver)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Pause() (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Pause(&_TransferValidatorForSolana.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) Pause() (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Pause(&_TransferValidatorForSolana.TransactOpts)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) RemoveReceiver(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "removeReceiver", _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.RemoveReceiver(&_TransferValidatorForSolana.TransactOpts, _receiver)
}

// RemoveReceiver is a paid mutator transaction binding the contract method 0x6552d8b4.
//
// Solidity: function removeReceiver(address _receiver) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) RemoveReceiver(_receiver common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.RemoveReceiver(&_TransferValidatorForSolana.TransactOpts, _receiver)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) Submit(opts *bind.TransactOpts, cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Submit(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Submit(&_TransferValidatorForSolana.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) Submit(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Submit(&_TransferValidatorForSolana.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.TransferOwnership(&_TransferValidatorForSolana.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.TransferOwnership(&_TransferValidatorForSolana.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Unpause() (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Unpause(&_TransferValidatorForSolana.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) Unpause() (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Unpause(&_TransferValidatorForSolana.TransactOpts)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactor) Upgrade(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.contract.Transact(opts, "upgrade", _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Upgrade(&_TransferValidatorForSolana.TransactOpts, _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidatorForSolana *TransferValidatorForSolanaTransactorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidatorForSolana.Contract.Upgrade(&_TransferValidatorForSolana.TransactOpts, _newValidator)
}

// TransferValidatorForSolanaOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaOwnershipTransferredIterator struct {
	Event *TransferValidatorForSolanaOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaOwnershipTransferred)
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
		it.Event = new(TransferValidatorForSolanaOwnershipTransferred)
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
func (it *TransferValidatorForSolanaOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaOwnershipTransferred represents a OwnershipTransferred event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferValidatorForSolanaOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaOwnershipTransferredIterator{contract: _TransferValidatorForSolana.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaOwnershipTransferred)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParseOwnershipTransferred(log types.Log) (*TransferValidatorForSolanaOwnershipTransferred, error) {
	event := new(TransferValidatorForSolanaOwnershipTransferred)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorForSolanaPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaPauseIterator struct {
	Event *TransferValidatorForSolanaPause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaPause)
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
		it.Event = new(TransferValidatorForSolanaPause)
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
func (it *TransferValidatorForSolanaPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaPause represents a Pause event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterPause(opts *bind.FilterOpts) (*TransferValidatorForSolanaPauseIterator, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaPauseIterator{contract: _TransferValidatorForSolana.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaPause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaPause)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParsePause(log types.Log) (*TransferValidatorForSolanaPause, error) {
	event := new(TransferValidatorForSolanaPause)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorForSolanaReceiverAddedIterator is returned from FilterReceiverAdded and is used to iterate over the raw logs and unpacked data for ReceiverAdded events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaReceiverAddedIterator struct {
	Event *TransferValidatorForSolanaReceiverAdded // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaReceiverAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaReceiverAdded)
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
		it.Event = new(TransferValidatorForSolanaReceiverAdded)
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
func (it *TransferValidatorForSolanaReceiverAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaReceiverAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaReceiverAdded represents a ReceiverAdded event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaReceiverAdded struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverAdded is a free log retrieval operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterReceiverAdded(opts *bind.FilterOpts) (*TransferValidatorForSolanaReceiverAddedIterator, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaReceiverAddedIterator{contract: _TransferValidatorForSolana.contract, event: "ReceiverAdded", logs: logs, sub: sub}, nil
}

// WatchReceiverAdded is a free log subscription operation binding the contract event 0xbec1e1ee82037bd0301ab4218c8c148e3be5be35bdf180546d4ff862df359f35.
//
// Solidity: event ReceiverAdded(address receiver)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchReceiverAdded(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaReceiverAdded) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "ReceiverAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaReceiverAdded)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParseReceiverAdded(log types.Log) (*TransferValidatorForSolanaReceiverAdded, error) {
	event := new(TransferValidatorForSolanaReceiverAdded)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "ReceiverAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorForSolanaReceiverRemovedIterator is returned from FilterReceiverRemoved and is used to iterate over the raw logs and unpacked data for ReceiverRemoved events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaReceiverRemovedIterator struct {
	Event *TransferValidatorForSolanaReceiverRemoved // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaReceiverRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaReceiverRemoved)
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
		it.Event = new(TransferValidatorForSolanaReceiverRemoved)
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
func (it *TransferValidatorForSolanaReceiverRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaReceiverRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaReceiverRemoved represents a ReceiverRemoved event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaReceiverRemoved struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReceiverRemoved is a free log retrieval operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterReceiverRemoved(opts *bind.FilterOpts) (*TransferValidatorForSolanaReceiverRemovedIterator, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaReceiverRemovedIterator{contract: _TransferValidatorForSolana.contract, event: "ReceiverRemoved", logs: logs, sub: sub}, nil
}

// WatchReceiverRemoved is a free log subscription operation binding the contract event 0x2771977f239a332de92ab37b7275685268f164e51cda8f1356692695f4708f2f.
//
// Solidity: event ReceiverRemoved(address receiver)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchReceiverRemoved(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaReceiverRemoved) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "ReceiverRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaReceiverRemoved)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParseReceiverRemoved(log types.Log) (*TransferValidatorForSolanaReceiverRemoved, error) {
	event := new(TransferValidatorForSolanaReceiverRemoved)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "ReceiverRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorForSolanaSettledIterator is returned from FilterSettled and is used to iterate over the raw logs and unpacked data for Settled events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaSettledIterator struct {
	Event *TransferValidatorForSolanaSettled // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaSettled)
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
		it.Event = new(TransferValidatorForSolanaSettled)
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
func (it *TransferValidatorForSolanaSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaSettled represents a Settled event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaSettled struct {
	Key       [32]byte
	Witnesses []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterSettled(opts *bind.FilterOpts, key [][32]byte) (*TransferValidatorForSolanaSettledIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaSettledIterator{contract: _TransferValidatorForSolana.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaSettled, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaSettled)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Settled", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParseSettled(log types.Log) (*TransferValidatorForSolanaSettled, error) {
	event := new(TransferValidatorForSolanaSettled)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorForSolanaUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaUnpauseIterator struct {
	Event *TransferValidatorForSolanaUnpause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorForSolanaUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorForSolanaUnpause)
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
		it.Event = new(TransferValidatorForSolanaUnpause)
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
func (it *TransferValidatorForSolanaUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorForSolanaUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorForSolanaUnpause represents a Unpause event raised by the TransferValidatorForSolana contract.
type TransferValidatorForSolanaUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) FilterUnpause(opts *bind.FilterOpts) (*TransferValidatorForSolanaUnpauseIterator, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorForSolanaUnpauseIterator{contract: _TransferValidatorForSolana.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TransferValidatorForSolanaUnpause) (event.Subscription, error) {

	logs, sub, err := _TransferValidatorForSolana.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorForSolanaUnpause)
				if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TransferValidatorForSolana *TransferValidatorForSolanaFilterer) ParseUnpause(log types.Log) (*TransferValidatorForSolanaUnpause, error) {
	event := new(TransferValidatorForSolanaUnpause)
	if err := _TransferValidatorForSolana.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
