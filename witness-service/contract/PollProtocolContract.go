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

// IPollProtocolContractCandidate is an auto generated low-level Go binding around an user-defined struct.
type IPollProtocolContractCandidate struct {
	Id              common.Address
	Name            string
	OperatorAddress common.Address
	RewardAddress   common.Address
	BlsPubKey       []byte
	Votes           *big.Int
}

// IPollProtocolContractProbationInfo is an auto generated low-level Go binding around an user-defined struct.
type IPollProtocolContractProbationInfo struct {
	Candidate       common.Address
	OperatorAddress common.Address
	Count           uint32
}

// IPollProtocolContractProbationList is an auto generated low-level Go binding around an user-defined struct.
type IPollProtocolContractProbationList struct {
	ProbationInfo []IPollProtocolContractProbationInfo
	IntensityRate uint32
}

// PollProtocolContractMetaData contains all meta data concerning the PollProtocolContract contract.
var PollProtocolContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"ActiveBlockProducersByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ActiveBlockProducersByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BlockProducersByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"BlockProducersByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CandidatesByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"CandidatesByEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"id\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"blsPubKey\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structIPollProtocolContract.Candidate[]\",\"name\":\"candidates\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"ProbationListByEpoch\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"count\",\"type\":\"uint32\"}],\"internalType\":\"structIPollProtocolContract.ProbationInfo[]\",\"name\":\"probationInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32\",\"name\":\"intensityRate\",\"type\":\"uint32\"}],\"internalType\":\"structIPollProtocolContract.ProbationList\",\"name\":\"probation\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ProbationListByEpoch\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"count\",\"type\":\"uint32\"}],\"internalType\":\"structIPollProtocolContract.ProbationInfo[]\",\"name\":\"probationInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint32\",\"name\":\"intensityRate\",\"type\":\"uint32\"}],\"internalType\":\"structIPollProtocolContract.ProbationList\",\"name\":\"probation\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PollProtocolContractABI is the input ABI used to generate the binding from.
// Deprecated: Use PollProtocolContractMetaData.ABI instead.
var PollProtocolContractABI = PollProtocolContractMetaData.ABI

// PollProtocolContract is an auto generated Go binding around an Ethereum contract.
type PollProtocolContract struct {
	PollProtocolContractCaller     // Read-only binding to the contract
	PollProtocolContractTransactor // Write-only binding to the contract
	PollProtocolContractFilterer   // Log filterer for contract events
}

// PollProtocolContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type PollProtocolContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PollProtocolContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PollProtocolContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PollProtocolContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PollProtocolContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PollProtocolContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PollProtocolContractSession struct {
	Contract     *PollProtocolContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PollProtocolContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PollProtocolContractCallerSession struct {
	Contract *PollProtocolContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// PollProtocolContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PollProtocolContractTransactorSession struct {
	Contract     *PollProtocolContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// PollProtocolContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type PollProtocolContractRaw struct {
	Contract *PollProtocolContract // Generic contract binding to access the raw methods on
}

// PollProtocolContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PollProtocolContractCallerRaw struct {
	Contract *PollProtocolContractCaller // Generic read-only contract binding to access the raw methods on
}

// PollProtocolContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PollProtocolContractTransactorRaw struct {
	Contract *PollProtocolContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPollProtocolContract creates a new instance of PollProtocolContract, bound to a specific deployed contract.
func NewPollProtocolContract(address common.Address, backend bind.ContractBackend) (*PollProtocolContract, error) {
	contract, err := bindPollProtocolContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PollProtocolContract{PollProtocolContractCaller: PollProtocolContractCaller{contract: contract}, PollProtocolContractTransactor: PollProtocolContractTransactor{contract: contract}, PollProtocolContractFilterer: PollProtocolContractFilterer{contract: contract}}, nil
}

// NewPollProtocolContractCaller creates a new read-only instance of PollProtocolContract, bound to a specific deployed contract.
func NewPollProtocolContractCaller(address common.Address, caller bind.ContractCaller) (*PollProtocolContractCaller, error) {
	contract, err := bindPollProtocolContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PollProtocolContractCaller{contract: contract}, nil
}

// NewPollProtocolContractTransactor creates a new write-only instance of PollProtocolContract, bound to a specific deployed contract.
func NewPollProtocolContractTransactor(address common.Address, transactor bind.ContractTransactor) (*PollProtocolContractTransactor, error) {
	contract, err := bindPollProtocolContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PollProtocolContractTransactor{contract: contract}, nil
}

// NewPollProtocolContractFilterer creates a new log filterer instance of PollProtocolContract, bound to a specific deployed contract.
func NewPollProtocolContractFilterer(address common.Address, filterer bind.ContractFilterer) (*PollProtocolContractFilterer, error) {
	contract, err := bindPollProtocolContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PollProtocolContractFilterer{contract: contract}, nil
}

// bindPollProtocolContract binds a generic wrapper to an already deployed contract.
func bindPollProtocolContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PollProtocolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PollProtocolContract *PollProtocolContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PollProtocolContract.Contract.PollProtocolContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PollProtocolContract *PollProtocolContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PollProtocolContract.Contract.PollProtocolContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PollProtocolContract *PollProtocolContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PollProtocolContract.Contract.PollProtocolContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PollProtocolContract *PollProtocolContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PollProtocolContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PollProtocolContract *PollProtocolContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PollProtocolContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PollProtocolContract *PollProtocolContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PollProtocolContract.Contract.contract.Transact(opts, method, params...)
}

// ActiveBlockProducersByEpoch is a free data retrieval call binding the contract method 0x33310089.
//
// Solidity: function ActiveBlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) ActiveBlockProducersByEpoch(opts *bind.CallOpts, epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "ActiveBlockProducersByEpoch", epoch)

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ActiveBlockProducersByEpoch is a free data retrieval call binding the contract method 0x33310089.
//
// Solidity: function ActiveBlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) ActiveBlockProducersByEpoch(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.ActiveBlockProducersByEpoch(&_PollProtocolContract.CallOpts, epoch)
}

// ActiveBlockProducersByEpoch is a free data retrieval call binding the contract method 0x33310089.
//
// Solidity: function ActiveBlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) ActiveBlockProducersByEpoch(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.ActiveBlockProducersByEpoch(&_PollProtocolContract.CallOpts, epoch)
}

// ActiveBlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x8f9e0b08.
//
// Solidity: function ActiveBlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) ActiveBlockProducersByEpoch0(opts *bind.CallOpts) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "ActiveBlockProducersByEpoch0")

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ActiveBlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x8f9e0b08.
//
// Solidity: function ActiveBlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) ActiveBlockProducersByEpoch0() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.ActiveBlockProducersByEpoch0(&_PollProtocolContract.CallOpts)
}

// ActiveBlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x8f9e0b08.
//
// Solidity: function ActiveBlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) ActiveBlockProducersByEpoch0() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.ActiveBlockProducersByEpoch0(&_PollProtocolContract.CallOpts)
}

// BlockProducersByEpoch is a free data retrieval call binding the contract method 0x68cfedce.
//
// Solidity: function BlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) BlockProducersByEpoch(opts *bind.CallOpts) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "BlockProducersByEpoch")

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// BlockProducersByEpoch is a free data retrieval call binding the contract method 0x68cfedce.
//
// Solidity: function BlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) BlockProducersByEpoch() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.BlockProducersByEpoch(&_PollProtocolContract.CallOpts)
}

// BlockProducersByEpoch is a free data retrieval call binding the contract method 0x68cfedce.
//
// Solidity: function BlockProducersByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) BlockProducersByEpoch() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.BlockProducersByEpoch(&_PollProtocolContract.CallOpts)
}

// BlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x809b1a49.
//
// Solidity: function BlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) BlockProducersByEpoch0(opts *bind.CallOpts, epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "BlockProducersByEpoch0", epoch)

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// BlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x809b1a49.
//
// Solidity: function BlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) BlockProducersByEpoch0(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.BlockProducersByEpoch0(&_PollProtocolContract.CallOpts, epoch)
}

// BlockProducersByEpoch0 is a free data retrieval call binding the contract method 0x809b1a49.
//
// Solidity: function BlockProducersByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) BlockProducersByEpoch0(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.BlockProducersByEpoch0(&_PollProtocolContract.CallOpts, epoch)
}

// CandidatesByEpoch is a free data retrieval call binding the contract method 0x2283cd98.
//
// Solidity: function CandidatesByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) CandidatesByEpoch(opts *bind.CallOpts) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "CandidatesByEpoch")

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CandidatesByEpoch is a free data retrieval call binding the contract method 0x2283cd98.
//
// Solidity: function CandidatesByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) CandidatesByEpoch() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.CandidatesByEpoch(&_PollProtocolContract.CallOpts)
}

// CandidatesByEpoch is a free data retrieval call binding the contract method 0x2283cd98.
//
// Solidity: function CandidatesByEpoch() view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) CandidatesByEpoch() (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.CandidatesByEpoch(&_PollProtocolContract.CallOpts)
}

// CandidatesByEpoch0 is a free data retrieval call binding the contract method 0xe5d7dbba.
//
// Solidity: function CandidatesByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) CandidatesByEpoch0(opts *bind.CallOpts, epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "CandidatesByEpoch0", epoch)

	outstruct := new(struct {
		Candidates []IPollProtocolContractCandidate
		Height     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Candidates = *abi.ConvertType(out[0], new([]IPollProtocolContractCandidate)).(*[]IPollProtocolContractCandidate)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CandidatesByEpoch0 is a free data retrieval call binding the contract method 0xe5d7dbba.
//
// Solidity: function CandidatesByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) CandidatesByEpoch0(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.CandidatesByEpoch0(&_PollProtocolContract.CallOpts, epoch)
}

// CandidatesByEpoch0 is a free data retrieval call binding the contract method 0xe5d7dbba.
//
// Solidity: function CandidatesByEpoch(uint256 epoch) view returns((address,string,address,address,bytes,uint256)[] candidates, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) CandidatesByEpoch0(epoch *big.Int) (struct {
	Candidates []IPollProtocolContractCandidate
	Height     *big.Int
}, error) {
	return _PollProtocolContract.Contract.CandidatesByEpoch0(&_PollProtocolContract.CallOpts, epoch)
}

// ProbationListByEpoch is a free data retrieval call binding the contract method 0x5ac5cb17.
//
// Solidity: function ProbationListByEpoch(uint256 epoch) view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) ProbationListByEpoch(opts *bind.CallOpts, epoch *big.Int) (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "ProbationListByEpoch", epoch)

	outstruct := new(struct {
		Probation IPollProtocolContractProbationList
		Height    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Probation = *abi.ConvertType(out[0], new(IPollProtocolContractProbationList)).(*IPollProtocolContractProbationList)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ProbationListByEpoch is a free data retrieval call binding the contract method 0x5ac5cb17.
//
// Solidity: function ProbationListByEpoch(uint256 epoch) view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) ProbationListByEpoch(epoch *big.Int) (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	return _PollProtocolContract.Contract.ProbationListByEpoch(&_PollProtocolContract.CallOpts, epoch)
}

// ProbationListByEpoch is a free data retrieval call binding the contract method 0x5ac5cb17.
//
// Solidity: function ProbationListByEpoch(uint256 epoch) view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) ProbationListByEpoch(epoch *big.Int) (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	return _PollProtocolContract.Contract.ProbationListByEpoch(&_PollProtocolContract.CallOpts, epoch)
}

// ProbationListByEpoch0 is a free data retrieval call binding the contract method 0xc2bf76e9.
//
// Solidity: function ProbationListByEpoch() view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractCaller) ProbationListByEpoch0(opts *bind.CallOpts) (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	var out []interface{}
	err := _PollProtocolContract.contract.Call(opts, &out, "ProbationListByEpoch0")

	outstruct := new(struct {
		Probation IPollProtocolContractProbationList
		Height    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Probation = *abi.ConvertType(out[0], new(IPollProtocolContractProbationList)).(*IPollProtocolContractProbationList)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ProbationListByEpoch0 is a free data retrieval call binding the contract method 0xc2bf76e9.
//
// Solidity: function ProbationListByEpoch() view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractSession) ProbationListByEpoch0() (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	return _PollProtocolContract.Contract.ProbationListByEpoch0(&_PollProtocolContract.CallOpts)
}

// ProbationListByEpoch0 is a free data retrieval call binding the contract method 0xc2bf76e9.
//
// Solidity: function ProbationListByEpoch() view returns(((address,address,uint32)[],uint32) probation, uint256 height)
func (_PollProtocolContract *PollProtocolContractCallerSession) ProbationListByEpoch0() (struct {
	Probation IPollProtocolContractProbationList
	Height    *big.Int
}, error) {
	return _PollProtocolContract.Contract.ProbationListByEpoch0(&_PollProtocolContract.CallOpts)
}
