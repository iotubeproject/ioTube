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

// TransferValidatorABI is the input ABI used to generate the binding from.
const TransferValidatorABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"expireHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"submissions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"witness\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transfers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"settleHeight\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"whitelistedTokens\",\"outputs\":[{\"internalType\":\"contractAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"whitelistedWitnesses\",\"outputs\":[{\"internalType\":\"contractAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_expireHeight\",\"type\":\"uint256\"}],\"name\":\"setExpireHeight\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"settled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"submit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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
	parsed, err := abi.JSON(strings.NewReader(TransferValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransferValidator *TransferValidatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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
func (_TransferValidator *TransferValidatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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

// ExpireHeight is a free data retrieval call binding the contract method 0xc219a0e3.
//
// Solidity: function expireHeight() view returns(uint256)
func (_TransferValidator *TransferValidatorCaller) ExpireHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "expireHeight")
	return *ret0, err
}

// ExpireHeight is a free data retrieval call binding the contract method 0xc219a0e3.
//
// Solidity: function expireHeight() view returns(uint256)
func (_TransferValidator *TransferValidatorSession) ExpireHeight() (*big.Int, error) {
	return _TransferValidator.Contract.ExpireHeight(&_TransferValidator.CallOpts)
}

// ExpireHeight is a free data retrieval call binding the contract method 0xc219a0e3.
//
// Solidity: function expireHeight() view returns(uint256)
func (_TransferValidator *TransferValidatorCallerSession) ExpireHeight() (*big.Int, error) {
	return _TransferValidator.Contract.ExpireHeight(&_TransferValidator.CallOpts)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3a9b56cc.
//
// Solidity: function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) pure returns(bytes32)
func (_TransferValidator *TransferValidatorCaller) GenerateKey(opts *bind.CallOpts, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "generateKey", tokenAddr, index, from, to, amount)
	return *ret0, err
}

// GenerateKey is a free data retrieval call binding the contract method 0x3a9b56cc.
//
// Solidity: function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) pure returns(bytes32)
func (_TransferValidator *TransferValidatorSession) GenerateKey(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, tokenAddr, index, from, to, amount)
}

// GenerateKey is a free data retrieval call binding the contract method 0x3a9b56cc.
//
// Solidity: function generateKey(address tokenAddr, uint256 index, address from, address to, uint256 amount) pure returns(bytes32)
func (_TransferValidator *TransferValidatorCallerSession) GenerateKey(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, tokenAddr, index, from, to, amount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TransferValidator *TransferValidatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "owner")
	return *ret0, err
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
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "paused")
	return *ret0, err
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

// Settled is a free data retrieval call binding the contract method 0x4bd7a8c1.
//
// Solidity: function settled(address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bool)
func (_TransferValidator *TransferValidatorCaller) Settled(opts *bind.CallOpts, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "settled", tokenAddr, index, from, to, amount)
	return *ret0, err
}

// Settled is a free data retrieval call binding the contract method 0x4bd7a8c1.
//
// Solidity: function settled(address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bool)
func (_TransferValidator *TransferValidatorSession) Settled(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (bool, error) {
	return _TransferValidator.Contract.Settled(&_TransferValidator.CallOpts, tokenAddr, index, from, to, amount)
}

// Settled is a free data retrieval call binding the contract method 0x4bd7a8c1.
//
// Solidity: function settled(address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bool)
func (_TransferValidator *TransferValidatorCallerSession) Settled(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (bool, error) {
	return _TransferValidator.Contract.Settled(&_TransferValidator.CallOpts, tokenAddr, index, from, to, amount)
}

// Submissions is a free data retrieval call binding the contract method 0x65aefe73.
//
// Solidity: function submissions(bytes32 , uint256 ) view returns(address witness, uint256 blockNumber)
func (_TransferValidator *TransferValidatorCaller) Submissions(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (struct {
	Witness     common.Address
	BlockNumber *big.Int
}, error) {
	ret := new(struct {
		Witness     common.Address
		BlockNumber *big.Int
	})
	out := ret
	err := _TransferValidator.contract.Call(opts, out, "submissions", arg0, arg1)
	return *ret, err
}

// Submissions is a free data retrieval call binding the contract method 0x65aefe73.
//
// Solidity: function submissions(bytes32 , uint256 ) view returns(address witness, uint256 blockNumber)
func (_TransferValidator *TransferValidatorSession) Submissions(arg0 [32]byte, arg1 *big.Int) (struct {
	Witness     common.Address
	BlockNumber *big.Int
}, error) {
	return _TransferValidator.Contract.Submissions(&_TransferValidator.CallOpts, arg0, arg1)
}

// Submissions is a free data retrieval call binding the contract method 0x65aefe73.
//
// Solidity: function submissions(bytes32 , uint256 ) view returns(address witness, uint256 blockNumber)
func (_TransferValidator *TransferValidatorCallerSession) Submissions(arg0 [32]byte, arg1 *big.Int) (struct {
	Witness     common.Address
	BlockNumber *big.Int
}, error) {
	return _TransferValidator.Contract.Submissions(&_TransferValidator.CallOpts, arg0, arg1)
}

// Transfers is a free data retrieval call binding the contract method 0x3c64f04b.
//
// Solidity: function transfers(bytes32 ) view returns(address tokenAddr, uint256 index, address from, address to, uint256 amount, uint256 settleHeight, bool flag)
func (_TransferValidator *TransferValidatorCaller) Transfers(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TokenAddr    common.Address
	Index        *big.Int
	From         common.Address
	To           common.Address
	Amount       *big.Int
	SettleHeight *big.Int
	Flag         bool
}, error) {
	ret := new(struct {
		TokenAddr    common.Address
		Index        *big.Int
		From         common.Address
		To           common.Address
		Amount       *big.Int
		SettleHeight *big.Int
		Flag         bool
	})
	out := ret
	err := _TransferValidator.contract.Call(opts, out, "transfers", arg0)
	return *ret, err
}

// Transfers is a free data retrieval call binding the contract method 0x3c64f04b.
//
// Solidity: function transfers(bytes32 ) view returns(address tokenAddr, uint256 index, address from, address to, uint256 amount, uint256 settleHeight, bool flag)
func (_TransferValidator *TransferValidatorSession) Transfers(arg0 [32]byte) (struct {
	TokenAddr    common.Address
	Index        *big.Int
	From         common.Address
	To           common.Address
	Amount       *big.Int
	SettleHeight *big.Int
	Flag         bool
}, error) {
	return _TransferValidator.Contract.Transfers(&_TransferValidator.CallOpts, arg0)
}

// Transfers is a free data retrieval call binding the contract method 0x3c64f04b.
//
// Solidity: function transfers(bytes32 ) view returns(address tokenAddr, uint256 index, address from, address to, uint256 amount, uint256 settleHeight, bool flag)
func (_TransferValidator *TransferValidatorCallerSession) Transfers(arg0 [32]byte) (struct {
	TokenAddr    common.Address
	Index        *big.Int
	From         common.Address
	To           common.Address
	Amount       *big.Int
	SettleHeight *big.Int
	Flag         bool
}, error) {
	return _TransferValidator.Contract.Transfers(&_TransferValidator.CallOpts, arg0)
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidator *TransferValidatorCaller) WhitelistedTokens(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "whitelistedTokens")
	return *ret0, err
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidator *TransferValidatorSession) WhitelistedTokens() (common.Address, error) {
	return _TransferValidator.Contract.WhitelistedTokens(&_TransferValidator.CallOpts)
}

// WhitelistedTokens is a free data retrieval call binding the contract method 0x5e1762a0.
//
// Solidity: function whitelistedTokens() view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) WhitelistedTokens() (common.Address, error) {
	return _TransferValidator.Contract.WhitelistedTokens(&_TransferValidator.CallOpts)
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidator *TransferValidatorCaller) WhitelistedWitnesses(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TransferValidator.contract.Call(opts, out, "whitelistedWitnesses")
	return *ret0, err
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidator *TransferValidatorSession) WhitelistedWitnesses() (common.Address, error) {
	return _TransferValidator.Contract.WhitelistedWitnesses(&_TransferValidator.CallOpts)
}

// WhitelistedWitnesses is a free data retrieval call binding the contract method 0x92072052.
//
// Solidity: function whitelistedWitnesses() view returns(address)
func (_TransferValidator *TransferValidatorCallerSession) WhitelistedWitnesses() (common.Address, error) {
	return _TransferValidator.Contract.WhitelistedWitnesses(&_TransferValidator.CallOpts)
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

// SetExpireHeight is a paid mutator transaction binding the contract method 0x4a46e739.
//
// Solidity: function setExpireHeight(uint256 _expireHeight) returns()
func (_TransferValidator *TransferValidatorTransactor) SetExpireHeight(opts *bind.TransactOpts, _expireHeight *big.Int) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "setExpireHeight", _expireHeight)
}

// SetExpireHeight is a paid mutator transaction binding the contract method 0x4a46e739.
//
// Solidity: function setExpireHeight(uint256 _expireHeight) returns()
func (_TransferValidator *TransferValidatorSession) SetExpireHeight(_expireHeight *big.Int) (*types.Transaction, error) {
	return _TransferValidator.Contract.SetExpireHeight(&_TransferValidator.TransactOpts, _expireHeight)
}

// SetExpireHeight is a paid mutator transaction binding the contract method 0x4a46e739.
//
// Solidity: function setExpireHeight(uint256 _expireHeight) returns()
func (_TransferValidator *TransferValidatorTransactorSession) SetExpireHeight(_expireHeight *big.Int) (*types.Transaction, error) {
	return _TransferValidator.Contract.SetExpireHeight(&_TransferValidator.TransactOpts, _expireHeight)
}

// Submit is a paid mutator transaction binding the contract method 0x2145c88c.
//
// Solidity: function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) returns()
func (_TransferValidator *TransferValidatorTransactor) Submit(opts *bind.TransactOpts, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "submit", tokenAddr, index, from, to, amount)
}

// Submit is a paid mutator transaction binding the contract method 0x2145c88c.
//
// Solidity: function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) returns()
func (_TransferValidator *TransferValidatorSession) Submit(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, tokenAddr, index, from, to, amount)
}

// Submit is a paid mutator transaction binding the contract method 0x2145c88c.
//
// Solidity: function submit(address tokenAddr, uint256 index, address from, address to, uint256 amount) returns()
func (_TransferValidator *TransferValidatorTransactorSession) Submit(tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, tokenAddr, index, from, to, amount)
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
	Token       common.Address
	Index       *big.Int
	From        common.Address
	To          common.Address
	Amount      *big.Int
	BlockNumber *big.Int
	Witnesses   []common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSettled is a free log retrieval operation binding the contract event 0xadf6b97a37e85478f91d8443b8caf744c91d2adad56befc7b15a04e23e2ddc11.
//
// Solidity: event Settled(address indexed token, uint256 indexed index, address indexed from, address to, uint256 amount, uint256 blockNumber, address[] witnesses)
func (_TransferValidator *TransferValidatorFilterer) FilterSettled(opts *bind.FilterOpts, token []common.Address, index []*big.Int, from []common.Address) (*TransferValidatorSettledIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TransferValidator.contract.FilterLogs(opts, "Settled", tokenRule, indexRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &TransferValidatorSettledIterator{contract: _TransferValidator.contract, event: "Settled", logs: logs, sub: sub}, nil
}

// WatchSettled is a free log subscription operation binding the contract event 0xadf6b97a37e85478f91d8443b8caf744c91d2adad56befc7b15a04e23e2ddc11.
//
// Solidity: event Settled(address indexed token, uint256 indexed index, address indexed from, address to, uint256 amount, uint256 blockNumber, address[] witnesses)
func (_TransferValidator *TransferValidatorFilterer) WatchSettled(opts *bind.WatchOpts, sink chan<- *TransferValidatorSettled, token []common.Address, index []*big.Int, from []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TransferValidator.contract.WatchLogs(opts, "Settled", tokenRule, indexRule, fromRule)
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

// ParseSettled is a log parse operation binding the contract event 0xadf6b97a37e85478f91d8443b8caf744c91d2adad56befc7b15a04e23e2ddc11.
//
// Solidity: event Settled(address indexed token, uint256 indexed index, address indexed from, address to, uint256 amount, uint256 blockNumber, address[] witnesses)
func (_TransferValidator *TransferValidatorFilterer) ParseSettled(log types.Log) (*TransferValidatorSettled, error) {
	event := new(TransferValidatorSettled)
	if err := _TransferValidator.contract.UnpackLog(event, "Settled", log); err != nil {
		return nil, err
	}
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
	return event, nil
}
