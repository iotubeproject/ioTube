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

// TransferValidatorMetaData contains all meta data concerning the TransferValidator contract.
var TransferValidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessList\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_tokenList\",\"type\":\"address\"},{\"internalType\":\"contractIMinter\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"addPair\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keys\",\"type\":\"bytes32[]\"}],\"name\":\"concatKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"extractWitnesses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"cashier\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getTokenGroup\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"minters\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numOfPairs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"cashier\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"witnessList\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TransferValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use TransferValidatorMetaData.ABI instead.
var TransferValidatorABI = TransferValidatorMetaData.ABI

// TransferValidator is an auto generated Go binding around an Ethereum contract.
type TransferValidator struct {
	TransferValidatorCaller     // Read-only binding to the contract
	TransferValidatorTransactor // Write-only binding to the contract
	TransferValidatorFilterer   // Log filterer for contract events
}

// TransferValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransferValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferValidatorSession struct {
	Contract     *TransferValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// TransferValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferValidatorCallerSession struct {
	Contract *TransferValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// TransferValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferValidatorTransactorSession struct {
	Contract     *TransferValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// TransferValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransferValidatorRaw struct {
	Contract *TransferValidator // Generic contract binding to access the raw methods on
}

// TransferValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferValidatorCallerRaw struct {
	Contract *TransferValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// TransferValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferValidatorTransactorRaw struct {
	Contract *TransferValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransferValidator creates a new instance of TransferValidator, bound to a specific deployed contract.
func NewTransferValidator(address common.Address, backend bind.ContractBackend) (*TransferValidator, error) {
	contract, err := bindTransferValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransferValidator{TransferValidatorCaller: TransferValidatorCaller{contract: contract}, TransferValidatorTransactor: TransferValidatorTransactor{contract: contract}, TransferValidatorFilterer: TransferValidatorFilterer{contract: contract}}, nil
}

// NewTransferValidatorCaller creates a new read-only instance of TransferValidator, bound to a specific deployed contract.
func NewTransferValidatorCaller(address common.Address, caller bind.ContractCaller) (*TransferValidatorCaller, error) {
	contract, err := bindTransferValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorCaller{contract: contract}, nil
}

// NewTransferValidatorTransactor creates a new write-only instance of TransferValidator, bound to a specific deployed contract.
func NewTransferValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TransferValidatorTransactor, error) {
	contract, err := bindTransferValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorTransactor{contract: contract}, nil
}

// NewTransferValidatorFilterer creates a new log filterer instance of TransferValidator, bound to a specific deployed contract.
func NewTransferValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TransferValidatorFilterer, error) {
	contract, err := bindTransferValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorFilterer{contract: contract}, nil
}

// bindTransferValidator binds a generic wrapper to an already deployed contract.
func bindTransferValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransferValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidator *TransferValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidator.Contract.TransferValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidator *TransferValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidator.Contract.TransferValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidator *TransferValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidator.Contract.TransferValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidator *TransferValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransferValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransferValidator *TransferValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransferValidator *TransferValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransferValidator.Contract.contract.Transact(opts, method, params...)
}

// ConcatKeys is a free data retrieval call binding the contract method 0xc836fef0.
//
// Solidity: function concatKeys(bytes32[] keys) pure returns(bytes32)
func (_TransferValidator *TransferValidatorCaller) ConcatKeys(opts *bind.CallOpts, keys [][32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "concatKeys", keys)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ConcatKeys is a free data retrieval call binding the contract method 0xc836fef0.
//
// Solidity: function concatKeys(bytes32[] keys) pure returns(bytes32)
func (_TransferValidator *TransferValidatorSession) ConcatKeys(keys [][32]byte) ([32]byte, error) {
	return _TransferValidator.Contract.ConcatKeys(&_TransferValidator.CallOpts, keys)
}

// ConcatKeys is a free data retrieval call binding the contract method 0xc836fef0.
//
// Solidity: function concatKeys(bytes32[] keys) pure returns(bytes32)
func (_TransferValidator *TransferValidatorCallerSession) ConcatKeys(keys [][32]byte) ([32]byte, error) {
	return _TransferValidator.Contract.ConcatKeys(&_TransferValidator.CallOpts, keys)
}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidator *TransferValidatorCaller) ExtractWitnesses(opts *bind.CallOpts, key [32]byte, signatures []byte) ([]common.Address, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "extractWitnesses", key, signatures)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidator *TransferValidatorSession) ExtractWitnesses(key [32]byte, signatures []byte) ([]common.Address, error) {
	return _TransferValidator.Contract.ExtractWitnesses(&_TransferValidator.CallOpts, key, signatures)
}

// ExtractWitnesses is a free data retrieval call binding the contract method 0xba390a64.
//
// Solidity: function extractWitnesses(bytes32 key, bytes signatures) view returns(address[] witnesses)
func (_TransferValidator *TransferValidatorCallerSession) ExtractWitnesses(key [32]byte, signatures []byte) ([]common.Address, error) {
	return _TransferValidator.Contract.ExtractWitnesses(&_TransferValidator.CallOpts, key, signatures)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidator *TransferValidatorCaller) GenerateKey(opts *bind.CallOpts, cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "generateKey", cashier, tokenAddr, index, from, to, amount, payload)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidator *TransferValidatorSession) GenerateKey(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3e6882aa.
//
// Solidity: function generateKey(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes payload) view returns(bytes32)
func (_TransferValidator *TransferValidatorCallerSession) GenerateKey(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, payload []byte) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, cashier, tokenAddr, index, from, to, amount, payload)
}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidator *TransferValidatorCaller) GetTokenGroup(opts *bind.CallOpts, tokenAddr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "getTokenGroup", tokenAddr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidator *TransferValidatorSession) GetTokenGroup(tokenAddr common.Address) (*big.Int, error) {
	return _TransferValidator.Contract.GetTokenGroup(&_TransferValidator.CallOpts, tokenAddr)
}

// GetTokenGroup is a free data retrieval call binding the contract method 0xe01eba71.
//
// Solidity: function getTokenGroup(address tokenAddr) view returns(uint256)
func (_TransferValidator *TransferValidatorCallerSession) GetTokenGroup(tokenAddr common.Address) (*big.Int, error) {
	return _TransferValidator.Contract.GetTokenGroup(&_TransferValidator.CallOpts, tokenAddr)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorCaller) Minters(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "minters", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidator.Contract.Minters(&_TransferValidator.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0x8623ec7b.
//
// Solidity: function minters(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) Minters(arg0 *big.Int) (common.Address, error) {
	return _TransferValidator.Contract.Minters(&_TransferValidator.CallOpts, arg0)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidator *TransferValidatorCaller) NumOfPairs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "numOfPairs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidator *TransferValidatorSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidator.Contract.NumOfPairs(&_TransferValidator.CallOpts)
}

// NumOfPairs is a free data retrieval call binding the contract method 0x8356b148.
//
// Solidity: function numOfPairs() view returns(uint256)
func (_TransferValidator *TransferValidatorCallerSession) NumOfPairs() (*big.Int, error) {
	return _TransferValidator.Contract.NumOfPairs(&_TransferValidator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidator *TransferValidatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidator *TransferValidatorSession) Owner() (common.Address, error) {
	return _TransferValidator.Contract.Owner(&_TransferValidator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) Owner() (common.Address, error) {
	return _TransferValidator.Contract.Owner(&_TransferValidator.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidator *TransferValidatorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidator *TransferValidatorSession) Paused() (bool, error) {
	return _TransferValidator.Contract.Paused(&_TransferValidator.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TransferValidator *TransferValidatorCallerSession) Paused() (bool, error) {
	return _TransferValidator.Contract.Paused(&_TransferValidator.CallOpts)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidator *TransferValidatorCaller) Settles(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "settles", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidator *TransferValidatorSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidator.Contract.Settles(&_TransferValidator.CallOpts, arg0)
}

// Settles is a free data retrieval call binding the contract method 0xf98b2332.
//
// Solidity: function settles(bytes32 ) view returns(uint256)
func (_TransferValidator *TransferValidatorCallerSession) Settles(arg0 [32]byte) (*big.Int, error) {
	return _TransferValidator.Contract.Settles(&_TransferValidator.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidator.Contract.TokenLists(&_TransferValidator.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TransferValidator.Contract.TokenLists(&_TransferValidator.CallOpts, arg0)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidator *TransferValidatorCaller) WitnessList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "witnessList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidator *TransferValidatorSession) WitnessList() (common.Address, error) {
	return _TransferValidator.Contract.WitnessList(&_TransferValidator.CallOpts)
}

// WitnessList is a free data retrieval call binding the contract method 0x373f0d49.
//
// Solidity: function witnessList() view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) WitnessList() (common.Address, error) {
	return _TransferValidator.Contract.WitnessList(&_TransferValidator.CallOpts)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidator *TransferValidatorTransactor) AddPair(opts *bind.TransactOpts, _tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "addPair", _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidator *TransferValidatorSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.AddPair(&_TransferValidator.TransactOpts, _tokenList, _minter)
}

// AddPair is a paid mutator transaction binding the contract method 0xb6f3e087.
//
// Solidity: function addPair(address _tokenList, address _minter) returns()
func (_TransferValidator *TransferValidatorTransactorSession) AddPair(_tokenList common.Address, _minter common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.AddPair(&_TransferValidator.TransactOpts, _tokenList, _minter)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidator *TransferValidatorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidator *TransferValidatorSession) Pause() (*types.Transaction, error) {
	return _TransferValidator.Contract.Pause(&_TransferValidator.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TransferValidator *TransferValidatorTransactorSession) Pause() (*types.Transaction, error) {
	return _TransferValidator.Contract.Pause(&_TransferValidator.TransactOpts)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidator *TransferValidatorTransactor) Submit(opts *bind.TransactOpts, cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidator *TransferValidatorSession) Submit(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// Submit is a paid mutator transaction binding the contract method 0x87726554.
//
// Solidity: function submit(bytes cashier, address tokenAddr, uint256 index, bytes from, address to, uint256 amount, bytes signatures, bytes payload) returns()
func (_TransferValidator *TransferValidatorTransactorSession) Submit(cashier []byte, tokenAddr common.Address, index *big.Int, from []byte, to common.Address, amount *big.Int, signatures []byte, payload []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures, payload)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidator *TransferValidatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidator *TransferValidatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.TransferOwnership(&_TransferValidator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TransferValidator *TransferValidatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.TransferOwnership(&_TransferValidator.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidator *TransferValidatorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidator *TransferValidatorSession) Unpause() (*types.Transaction, error) {
	return _TransferValidator.Contract.Unpause(&_TransferValidator.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TransferValidator *TransferValidatorTransactorSession) Unpause() (*types.Transaction, error) {
	return _TransferValidator.Contract.Unpause(&_TransferValidator.TransactOpts)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidator *TransferValidatorTransactor) Upgrade(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "upgrade", _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidator *TransferValidatorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.Upgrade(&_TransferValidator.TransactOpts, _newValidator)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address _newValidator) returns()
func (_TransferValidator *TransferValidatorTransactorSession) Upgrade(_newValidator common.Address) (*types.Transaction, error) {
	return _TransferValidator.Contract.Upgrade(&_TransferValidator.TransactOpts, _newValidator)
}

// TransferValidatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TransferValidator contract.
type TransferValidatorOwnershipTransferredIterator struct {
	Event *TransferValidatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TransferValidatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorOwnershipTransferred)
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
		it.Event = new(TransferValidatorOwnershipTransferred)
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
func (it *TransferValidatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorOwnershipTransferred represents a OwnershipTransferred event raised by the TransferValidator contract.
type TransferValidatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidator *TransferValidatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TransferValidatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorOwnershipTransferredIterator{contract: _TransferValidator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TransferValidator *TransferValidatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TransferValidatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TransferValidator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorOwnershipTransferred)
				if err := _TransferValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TransferValidator *TransferValidatorFilterer) ParseOwnershipTransferred(log types.Log) (*TransferValidatorOwnershipTransferred, error) {
	event := new(TransferValidatorOwnershipTransferred)
	if err := _TransferValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TransferValidator contract.
type TransferValidatorPauseIterator struct {
	Event *TransferValidatorPause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorPause)
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
		it.Event = new(TransferValidatorPause)
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
func (it *TransferValidatorPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorPause represents a Pause event raised by the TransferValidator contract.
type TransferValidatorPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidator *TransferValidatorFilterer) FilterPause(opts *bind.FilterOpts) (*TransferValidatorPauseIterator, error) {

	logs, sub, err := _TransferValidator.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorPauseIterator{contract: _TransferValidator.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TransferValidator *TransferValidatorFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TransferValidatorPause) (event.Subscription, error) {

	logs, sub, err := _TransferValidator.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorPause)
				if err := _TransferValidator.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TransferValidator *TransferValidatorFilterer) ParsePause(log types.Log) (*TransferValidatorPause, error) {
	event := new(TransferValidatorPause)
	if err := _TransferValidator.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorSettledIterator is returned from FilterSettled and is used to iterate over the raw logs and unpacked data for Settled events raised by the TransferValidator contract.
type TransferValidatorSettledIterator struct {
	Event *TransferValidatorSettled // Event containing the contract specifics and raw log

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
func (it *TransferValidatorSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorSettled)
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
		it.Event = new(TransferValidatorSettled)
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
func (it *TransferValidatorSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorSettled represents a Settled event raised by the TransferValidator contract.
type TransferValidatorSettled struct {
	Key       [32]byte
	Witnesses []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidator *TransferValidatorFilterer) FilterSettled(opts *bind.FilterOpts, key [][32]byte) (*TransferValidatorSettledIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidator.contract.FilterLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorSettledIterator{contract: _TransferValidator.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0.
//
// Solidity: event Settled(bytes32 indexed key, address[] witnesses)
func (_TransferValidator *TransferValidatorFilterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorSettled, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _TransferValidator.contract.WatchLogs(opts, "Settled", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorSettled)
				if err := _TransferValidator.contract.UnpackLog(event, "Settled", log); err != nil {
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
func (_TransferValidator *TransferValidatorFilterer) ParseSettled(log types.Log) (*TransferValidatorSettled, error) {
	event := new(TransferValidatorSettled)
	if err := _TransferValidator.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransferValidatorUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TransferValidator contract.
type TransferValidatorUnpauseIterator struct {
	Event *TransferValidatorUnpause // Event containing the contract specifics and raw log

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
func (it *TransferValidatorUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransferValidatorUnpause)
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
		it.Event = new(TransferValidatorUnpause)
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
func (it *TransferValidatorUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransferValidatorUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransferValidatorUnpause represents a Unpause event raised by the TransferValidator contract.
type TransferValidatorUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidator *TransferValidatorFilterer) FilterUnpause(opts *bind.FilterOpts) (*TransferValidatorUnpauseIterator, error) {

	logs, sub, err := _TransferValidator.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TransferValidatorUnpauseIterator{contract: _TransferValidator.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TransferValidator *TransferValidatorFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TransferValidatorUnpause) (event.Subscription, error) {

	logs, sub, err := _TransferValidator.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransferValidatorUnpause)
				if err := _TransferValidator.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TransferValidator *TransferValidatorFilterer) ParseUnpause(log types.Log) (*TransferValidatorUnpause, error) {
	event := new(TransferValidatorUnpause)
	if err := _TransferValidator.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
