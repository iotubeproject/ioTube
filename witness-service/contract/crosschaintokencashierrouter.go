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

// CrosschainTokenCashierRouterABI is the input ABI used to generate the binding from.
const CrosschainTokenCashierRouterABI = "[{\"inputs\":[{\"internalType\":\"contractICashier\",\"name\":\"_cashier\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"cashier\",\"outputs\":[{\"internalType\":\"contractICashier\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_crosschainToken\",\"type\":\"address\"}],\"name\":\"approveCrosschainToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_crosschainToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// CrosschainTokenCashierRouter is an auto generated Go binding around an Ethereum contract.
type CrosschainTokenCashierRouter struct {
	CrosschainTokenCashierRouterCaller     // Read-only binding to the contract
	CrosschainTokenCashierRouterTransactor // Write-only binding to the contract
	CrosschainTokenCashierRouterFilterer   // Log filterer for contract events
}

// CrosschainTokenCashierRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrosschainTokenCashierRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTokenCashierRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrosschainTokenCashierRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTokenCashierRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrosschainTokenCashierRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTokenCashierRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrosschainTokenCashierRouterSession struct {
	Contract     *CrosschainTokenCashierRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// CrosschainTokenCashierRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrosschainTokenCashierRouterCallerSession struct {
	Contract *CrosschainTokenCashierRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// CrosschainTokenCashierRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrosschainTokenCashierRouterTransactorSession struct {
	Contract     *CrosschainTokenCashierRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// CrosschainTokenCashierRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrosschainTokenCashierRouterRaw struct {
	Contract *CrosschainTokenCashierRouter // Generic contract binding to access the raw methods on
}

// CrosschainTokenCashierRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrosschainTokenCashierRouterCallerRaw struct {
	Contract *CrosschainTokenCashierRouterCaller // Generic read-only contract binding to access the raw methods on
}

// CrosschainTokenCashierRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrosschainTokenCashierRouterTransactorRaw struct {
	Contract *CrosschainTokenCashierRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrosschainTokenCashierRouter creates a new instance of CrosschainTokenCashierRouter, bound to a specific deployed contract.
func NewCrosschainTokenCashierRouter(address common.Address, backend bind.ContractBackend) (*CrosschainTokenCashierRouter, error) {
	contract, err := bindCrosschainTokenCashierRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrosschainTokenCashierRouter{CrosschainTokenCashierRouterCaller: CrosschainTokenCashierRouterCaller{contract: contract}, CrosschainTokenCashierRouterTransactor: CrosschainTokenCashierRouterTransactor{contract: contract}, CrosschainTokenCashierRouterFilterer: CrosschainTokenCashierRouterFilterer{contract: contract}}, nil
}

// NewCrosschainTokenCashierRouterCaller creates a new read-only instance of CrosschainTokenCashierRouter, bound to a specific deployed contract.
func NewCrosschainTokenCashierRouterCaller(address common.Address, caller bind.ContractCaller) (*CrosschainTokenCashierRouterCaller, error) {
	contract, err := bindCrosschainTokenCashierRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTokenCashierRouterCaller{contract: contract}, nil
}

// NewCrosschainTokenCashierRouterTransactor creates a new write-only instance of CrosschainTokenCashierRouter, bound to a specific deployed contract.
func NewCrosschainTokenCashierRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*CrosschainTokenCashierRouterTransactor, error) {
	contract, err := bindCrosschainTokenCashierRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTokenCashierRouterTransactor{contract: contract}, nil
}

// NewCrosschainTokenCashierRouterFilterer creates a new log filterer instance of CrosschainTokenCashierRouter, bound to a specific deployed contract.
func NewCrosschainTokenCashierRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*CrosschainTokenCashierRouterFilterer, error) {
	contract, err := bindCrosschainTokenCashierRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrosschainTokenCashierRouterFilterer{contract: contract}, nil
}

// bindCrosschainTokenCashierRouter binds a generic wrapper to an already deployed contract.
func bindCrosschainTokenCashierRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CrosschainTokenCashierRouterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTokenCashierRouter.Contract.CrosschainTokenCashierRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.CrosschainTokenCashierRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.CrosschainTokenCashierRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTokenCashierRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.contract.Transact(opts, method, params...)
}

// Cashier is a free data retrieval call binding the contract method 0xed740e97.
//
// Solidity: function cashier() view returns(address)
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterCaller) Cashier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrosschainTokenCashierRouter.contract.Call(opts, &out, "cashier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Cashier is a free data retrieval call binding the contract method 0xed740e97.
//
// Solidity: function cashier() view returns(address)
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterSession) Cashier() (common.Address, error) {
	return _CrosschainTokenCashierRouter.Contract.Cashier(&_CrosschainTokenCashierRouter.CallOpts)
}

// Cashier is a free data retrieval call binding the contract method 0xed740e97.
//
// Solidity: function cashier() view returns(address)
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterCallerSession) Cashier() (common.Address, error) {
	return _CrosschainTokenCashierRouter.Contract.Cashier(&_CrosschainTokenCashierRouter.CallOpts)
}

// ApproveCrosschainToken is a paid mutator transaction binding the contract method 0x0d71f7ca.
//
// Solidity: function approveCrosschainToken(address _crosschainToken) returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactor) ApproveCrosschainToken(opts *bind.TransactOpts, _crosschainToken common.Address) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.contract.Transact(opts, "approveCrosschainToken", _crosschainToken)
}

// ApproveCrosschainToken is a paid mutator transaction binding the contract method 0x0d71f7ca.
//
// Solidity: function approveCrosschainToken(address _crosschainToken) returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterSession) ApproveCrosschainToken(_crosschainToken common.Address) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.ApproveCrosschainToken(&_CrosschainTokenCashierRouter.TransactOpts, _crosschainToken)
}

// ApproveCrosschainToken is a paid mutator transaction binding the contract method 0x0d71f7ca.
//
// Solidity: function approveCrosschainToken(address _crosschainToken) returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactorSession) ApproveCrosschainToken(_crosschainToken common.Address) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.ApproveCrosschainToken(&_CrosschainTokenCashierRouter.TransactOpts, _crosschainToken)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _crosschainToken, address _to, uint256 _amount) payable returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactor) DepositTo(opts *bind.TransactOpts, _crosschainToken common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.contract.Transact(opts, "depositTo", _crosschainToken, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _crosschainToken, address _to, uint256 _amount) payable returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterSession) DepositTo(_crosschainToken common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.DepositTo(&_CrosschainTokenCashierRouter.TransactOpts, _crosschainToken, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _crosschainToken, address _to, uint256 _amount) payable returns()
func (_CrosschainTokenCashierRouter *CrosschainTokenCashierRouterTransactorSession) DepositTo(_crosschainToken common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainTokenCashierRouter.Contract.DepositTo(&_CrosschainTokenCashierRouter.TransactOpts, _crosschainToken, _to, _amount)
}
