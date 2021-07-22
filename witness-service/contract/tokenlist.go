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

// TokenListABI is the input ABI used to generate the binding from.
const TokenListABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"}],\"name\":\"TokenUpdated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"limit\",\"type\":\"uint8\"}],\"name\":\"getActiveItems\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"items_\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_item\",\"type\":\"address\"}],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_item\",\"type\":\"address\"}],\"name\":\"isExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numOfActive\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"isAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"}],\"name\":\"addToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success_\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_mins\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_maxs\",\"type\":\"uint256[]\"}],\"name\":\"addTokens\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success_\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success_\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"}],\"name\":\"setMinAmount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_maxAmount\",\"type\":\"uint256\"}],\"name\":\"setMaxAmount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"minAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"minAmount_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"maxAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxAmount_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenListBin is the compiled bytecode used for deploying new contracts.
var TokenListBin = "0x6080604052600080546001600160a01b03191633179055610cc1806100256000396000f3fe608060405234801561001057600080fd5b50600436106100f45760003560e01c806394dee0a411610097578063dee1f2af11610066578063dee1f2af14610410578063f2fde38b14610442578063f7cb131214610468578063fb52b065146104e9576100f4565b806394dee0a4146103705780639f8a13d714610396578063babcc539146103bc578063d5708d5a146103e2576100f4565b8063593f6969116100d3578063593f6969146101735780635fa7b5841461017b57806373d432ce146101a15780638da5cb5b1461034c576100f4565b806213eb4b146100f957806306661abd146101335780634d0a32db1461014d575b600080fd5b61011f6004803603602081101561010f57600080fd5b50356001600160a01b0316610515565b604080519115158252519081900360200190f35b61013b610537565b60408051918252519081900360200190f35b61013b6004803603602081101561016357600080fd5b50356001600160a01b031661053d565b61013b610569565b61011f6004803603602081101561019157600080fd5b50356001600160a01b031661056f565b61011f600480360360608110156101b757600080fd5b8101906020810181356401000000008111156101d257600080fd5b8201836020820111156101e457600080fd5b8035906020019184602083028401116401000000008311171561020657600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929594936020810193503591505064010000000081111561025657600080fd5b82018360208201111561026857600080fd5b8035906020019184602083028401116401000000008311171561028a57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092959493602081019350359150506401000000008111156102da57600080fd5b8201836020820111156102ec57600080fd5b8035906020019184602083028401116401000000008311171561030e57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506105d1945050505050565b6103546106ab565b604080516001600160a01b039092168252519081900360200190f35b61013b6004803603602081101561038657600080fd5b50356001600160a01b03166106ba565b61011f600480360360208110156103ac57600080fd5b50356001600160a01b03166106e9565b61011f600480360360208110156103d257600080fd5b50356001600160a01b031661070c565b61040e600480360360408110156103f857600080fd5b506001600160a01b03813516906020013561071d565b005b61011f6004803603606081101561042657600080fd5b506001600160a01b0381351690602081013590604001356107d1565b61040e6004803603602081101561045857600080fd5b50356001600160a01b03166108c8565b61048e6004803603604081101561047e57600080fd5b508035906020013560ff1661094d565b6040518083815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156104d45781810151838201526020016104bc565b50505050905001935050505060405180910390f35b61040e600480360360408110156104ff57600080fd5b506001600160a01b038135169060200135610a6b565b6001600160a01b03811660009081526003602052604090205460ff165b919050565b60025490565b600061054882610515565b1561053257506001600160a01b031660009081526004602052604090205490565b60015490565b600080546001600160a01b0316331461058757600080fd5b61059082610b12565b15610532576040516001600160a01b038316907f4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd390600090a2506001919050565b600080546001600160a01b031633146105e957600080fd5b825184511480156105fb575081518351145b610641576040805162461bcd60e51b8152602060048201526012602482015271696e76616c696420706172616d657465727360701b604482015290519081900360640190fd5b60005b84518110156106a35761069185828151811061065c57fe5b602002602001015185838151811061067057fe5b602002602001015185848151811061068457fe5b60200260200101516107d1565b1561069b57600191505b600101610644565b509392505050565b6000546001600160a01b031681565b60006106c582610515565b1561053257506001600160a01b031660009081526004602052604090206001015490565b6001600160a01b0316600090815260036020526040902054610100900460ff1690565b6000610717826106e9565b92915050565b6000546001600160a01b0316331461073457600080fd5b61073d82610515565b610780576040805162461bcd60e51b815260206004820152600f60248201526e1d1bdad95b881b9bdd081859191959608a1b604482015290519081900360640190fd5b6001600160a01b0382166000908152600460205260409020600101548111156107a857600080fd5b600081116107b557600080fd5b6001600160a01b03909116600090815260046020526040902055565b600080546001600160a01b031633146107e957600080fd5b6107f284610b95565b156108c15760008311801561080657508282115b61084c576040805162461bcd60e51b8152602060048201526012602482015271696e76616c696420706172616d657465727360701b604482015290519081900360640190fd5b60408051808201825284815260208082018581526001600160a01b0388166000818152600484528590209351845590516001909301929092558251868152908101859052825191927fa818a22273fc309f0a3682b642c74c5b5c25c0615ff378d07767cd231e19fffc92918290030190a25060015b9392505050565b6000546001600160a01b031633146108df57600080fd5b6001600160a01b0381166108f257600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b60025460009060609084108015610966575060ff831615155b61096f57600080fd5b8260ff1660405190808252806020026020018201604052801561099c578160200160208202803883390190505b50905060005b8360ff16811015610a6357600254858201106109bd57610a63565b600360006002838801815481106109d057fe5b60009182526020808320909101546001600160a01b0316835282019290925260400190205460ff6101009091041615610a5b57600281860181548110610a1257fe5b9060005260206000200160009054906101000a90046001600160a01b0316828481518110610a3c57fe5b6001600160a01b03909216602092830291909101909101526001909201915b6001016109a2565b509250929050565b6000546001600160a01b03163314610a8257600080fd5b610a8b82610515565b610ace576040805162461bcd60e51b815260206004820152600f60248201526e1d1bdad95b881b9bdd081859191959608a1b604482015290519081900360640190fd5b6001600160a01b038216600090815260046020526040902054811015610af357600080fd5b6001600160a01b03909116600090815260046020526040902060010155565b6001600160a01b03811660009081526003602052604081205460ff168015610b5757506001600160a01b038216600090815260036020526040902054610100900460ff165b15610b8d5750600180546000190181556001600160a01b0382166000908152600360205260409020805461ff0019169055610532565b506000919050565b6001600160a01b038116600090815260036020526040812054610100900460ff1615610bc357506000610532565b6001600160a01b03821660009081526003602052604090205460ff16610c2f57600280546001810182556000919091527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace0180546001600160a01b0319166001600160a01b0384161790555b50600180548101815560408051808201825282815260208082018481526001600160a01b0386166000908152600390925292902090518154925115156101000261ff001991151560ff19909416939093171691909117905591905056fea265627a7a723158204b5c202045580f17099379d431a94f021547a70de58e83550a43c177b23d165964736f6c63430005110032"

// DeployTokenList deploys a new Ethereum contract, binding an instance of TokenList to it.
func DeployTokenList(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenList, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenListABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenListBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenList{TokenListCaller: TokenListCaller{contract: contract}, TokenListTransactor: TokenListTransactor{contract: contract}, TokenListFilterer: TokenListFilterer{contract: contract}}, nil
}

// TokenList is an auto generated Go binding around an Ethereum contract.
type TokenList struct {
	TokenListCaller     // Read-only binding to the contract
	TokenListTransactor // Write-only binding to the contract
	TokenListFilterer   // Log filterer for contract events
}

// TokenListCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenListCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenListTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenListTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenListFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenListSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenListSession struct {
	Contract     *TokenList        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenListCallerSession struct {
	Contract *TokenListCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TokenListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenListTransactorSession struct {
	Contract     *TokenListTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TokenListRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenListRaw struct {
	Contract *TokenList // Generic contract binding to access the raw methods on
}

// TokenListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenListCallerRaw struct {
	Contract *TokenListCaller // Generic read-only contract binding to access the raw methods on
}

// TokenListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenListTransactorRaw struct {
	Contract *TokenListTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenList creates a new instance of TokenList, bound to a specific deployed contract.
func NewTokenList(address common.Address, backend bind.ContractBackend) (*TokenList, error) {
	contract, err := bindTokenList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenList{TokenListCaller: TokenListCaller{contract: contract}, TokenListTransactor: TokenListTransactor{contract: contract}, TokenListFilterer: TokenListFilterer{contract: contract}}, nil
}

// NewTokenListCaller creates a new read-only instance of TokenList, bound to a specific deployed contract.
func NewTokenListCaller(address common.Address, caller bind.ContractCaller) (*TokenListCaller, error) {
	contract, err := bindTokenList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenListCaller{contract: contract}, nil
}

// NewTokenListTransactor creates a new write-only instance of TokenList, bound to a specific deployed contract.
func NewTokenListTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenListTransactor, error) {
	contract, err := bindTokenList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenListTransactor{contract: contract}, nil
}

// NewTokenListFilterer creates a new log filterer instance of TokenList, bound to a specific deployed contract.
func NewTokenListFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenListFilterer, error) {
	contract, err := bindTokenList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenListFilterer{contract: contract}, nil
}

// bindTokenList binds a generic wrapper to an already deployed contract.
func bindTokenList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenListABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenList *TokenListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenList.Contract.TokenListCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenList *TokenListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenList.Contract.TokenListTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenList *TokenListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenList.Contract.TokenListTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenList *TokenListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenList.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenList *TokenListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenList.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenList *TokenListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenList.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenList *TokenListCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenList *TokenListSession) Count() (*big.Int, error) {
	return _TokenList.Contract.Count(&_TokenList.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TokenList *TokenListCallerSession) Count() (*big.Int, error) {
	return _TokenList.Contract.Count(&_TokenList.CallOpts)
}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_TokenList *TokenListCaller) GetActiveItems(opts *bind.CallOpts, offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "getActiveItems", offset, limit)

	outstruct := new(struct {
		Count *big.Int
		Items []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Count = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Items = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_TokenList *TokenListSession) GetActiveItems(offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	return _TokenList.Contract.GetActiveItems(&_TokenList.CallOpts, offset, limit)
}

// GetActiveItems is a free data retrieval call binding the contract method 0xf7cb1312.
//
// Solidity: function getActiveItems(uint256 offset, uint8 limit) view returns(uint256 count_, address[] items_)
func (_TokenList *TokenListCallerSession) GetActiveItems(offset *big.Int, limit uint8) (struct {
	Count *big.Int
	Items []common.Address
}, error) {
	return _TokenList.Contract.GetActiveItems(&_TokenList.CallOpts, offset, limit)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_TokenList *TokenListCaller) IsActive(opts *bind.CallOpts, _item common.Address) (bool, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "isActive", _item)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_TokenList *TokenListSession) IsActive(_item common.Address) (bool, error) {
	return _TokenList.Contract.IsActive(&_TokenList.CallOpts, _item)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address _item) view returns(bool)
func (_TokenList *TokenListCallerSession) IsActive(_item common.Address) (bool, error) {
	return _TokenList.Contract.IsActive(&_TokenList.CallOpts, _item)
}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _token) view returns(bool)
func (_TokenList *TokenListCaller) IsAllowed(opts *bind.CallOpts, _token common.Address) (bool, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "isAllowed", _token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _token) view returns(bool)
func (_TokenList *TokenListSession) IsAllowed(_token common.Address) (bool, error) {
	return _TokenList.Contract.IsAllowed(&_TokenList.CallOpts, _token)
}

// IsAllowed is a free data retrieval call binding the contract method 0xbabcc539.
//
// Solidity: function isAllowed(address _token) view returns(bool)
func (_TokenList *TokenListCallerSession) IsAllowed(_token common.Address) (bool, error) {
	return _TokenList.Contract.IsAllowed(&_TokenList.CallOpts, _token)
}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_TokenList *TokenListCaller) IsExist(opts *bind.CallOpts, _item common.Address) (bool, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "isExist", _item)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_TokenList *TokenListSession) IsExist(_item common.Address) (bool, error) {
	return _TokenList.Contract.IsExist(&_TokenList.CallOpts, _item)
}

// IsExist is a free data retrieval call binding the contract method 0x0013eb4b.
//
// Solidity: function isExist(address _item) view returns(bool)
func (_TokenList *TokenListCallerSession) IsExist(_item common.Address) (bool, error) {
	return _TokenList.Contract.IsExist(&_TokenList.CallOpts, _item)
}

// MaxAmount is a free data retrieval call binding the contract method 0x94dee0a4.
//
// Solidity: function maxAmount(address _token) view returns(uint256 maxAmount_)
func (_TokenList *TokenListCaller) MaxAmount(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "maxAmount", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAmount is a free data retrieval call binding the contract method 0x94dee0a4.
//
// Solidity: function maxAmount(address _token) view returns(uint256 maxAmount_)
func (_TokenList *TokenListSession) MaxAmount(_token common.Address) (*big.Int, error) {
	return _TokenList.Contract.MaxAmount(&_TokenList.CallOpts, _token)
}

// MaxAmount is a free data retrieval call binding the contract method 0x94dee0a4.
//
// Solidity: function maxAmount(address _token) view returns(uint256 maxAmount_)
func (_TokenList *TokenListCallerSession) MaxAmount(_token common.Address) (*big.Int, error) {
	return _TokenList.Contract.MaxAmount(&_TokenList.CallOpts, _token)
}

// MinAmount is a free data retrieval call binding the contract method 0x4d0a32db.
//
// Solidity: function minAmount(address _token) view returns(uint256 minAmount_)
func (_TokenList *TokenListCaller) MinAmount(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "minAmount", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAmount is a free data retrieval call binding the contract method 0x4d0a32db.
//
// Solidity: function minAmount(address _token) view returns(uint256 minAmount_)
func (_TokenList *TokenListSession) MinAmount(_token common.Address) (*big.Int, error) {
	return _TokenList.Contract.MinAmount(&_TokenList.CallOpts, _token)
}

// MinAmount is a free data retrieval call binding the contract method 0x4d0a32db.
//
// Solidity: function minAmount(address _token) view returns(uint256 minAmount_)
func (_TokenList *TokenListCallerSession) MinAmount(_token common.Address) (*big.Int, error) {
	return _TokenList.Contract.MinAmount(&_TokenList.CallOpts, _token)
}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_TokenList *TokenListCaller) NumOfActive(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "numOfActive")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_TokenList *TokenListSession) NumOfActive() (*big.Int, error) {
	return _TokenList.Contract.NumOfActive(&_TokenList.CallOpts)
}

// NumOfActive is a free data retrieval call binding the contract method 0x593f6969.
//
// Solidity: function numOfActive() view returns(uint256)
func (_TokenList *TokenListCallerSession) NumOfActive() (*big.Int, error) {
	return _TokenList.Contract.NumOfActive(&_TokenList.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenList *TokenListCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenList.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenList *TokenListSession) Owner() (common.Address, error) {
	return _TokenList.Contract.Owner(&_TokenList.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenList *TokenListCallerSession) Owner() (common.Address, error) {
	return _TokenList.Contract.Owner(&_TokenList.CallOpts)
}

// AddToken is a paid mutator transaction binding the contract method 0xdee1f2af.
//
// Solidity: function addToken(address _token, uint256 _min, uint256 _max) returns(bool success_)
func (_TokenList *TokenListTransactor) AddToken(opts *bind.TransactOpts, _token common.Address, _min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "addToken", _token, _min, _max)
}

// AddToken is a paid mutator transaction binding the contract method 0xdee1f2af.
//
// Solidity: function addToken(address _token, uint256 _min, uint256 _max) returns(bool success_)
func (_TokenList *TokenListSession) AddToken(_token common.Address, _min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.AddToken(&_TokenList.TransactOpts, _token, _min, _max)
}

// AddToken is a paid mutator transaction binding the contract method 0xdee1f2af.
//
// Solidity: function addToken(address _token, uint256 _min, uint256 _max) returns(bool success_)
func (_TokenList *TokenListTransactorSession) AddToken(_token common.Address, _min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.AddToken(&_TokenList.TransactOpts, _token, _min, _max)
}

// AddTokens is a paid mutator transaction binding the contract method 0x73d432ce.
//
// Solidity: function addTokens(address[] _tokens, uint256[] _mins, uint256[] _maxs) returns(bool success_)
func (_TokenList *TokenListTransactor) AddTokens(opts *bind.TransactOpts, _tokens []common.Address, _mins []*big.Int, _maxs []*big.Int) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "addTokens", _tokens, _mins, _maxs)
}

// AddTokens is a paid mutator transaction binding the contract method 0x73d432ce.
//
// Solidity: function addTokens(address[] _tokens, uint256[] _mins, uint256[] _maxs) returns(bool success_)
func (_TokenList *TokenListSession) AddTokens(_tokens []common.Address, _mins []*big.Int, _maxs []*big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.AddTokens(&_TokenList.TransactOpts, _tokens, _mins, _maxs)
}

// AddTokens is a paid mutator transaction binding the contract method 0x73d432ce.
//
// Solidity: function addTokens(address[] _tokens, uint256[] _mins, uint256[] _maxs) returns(bool success_)
func (_TokenList *TokenListTransactorSession) AddTokens(_tokens []common.Address, _mins []*big.Int, _maxs []*big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.AddTokens(&_TokenList.TransactOpts, _tokens, _mins, _maxs)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address _token) returns(bool success_)
func (_TokenList *TokenListTransactor) RemoveToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "removeToken", _token)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address _token) returns(bool success_)
func (_TokenList *TokenListSession) RemoveToken(_token common.Address) (*types.Transaction, error) {
	return _TokenList.Contract.RemoveToken(&_TokenList.TransactOpts, _token)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address _token) returns(bool success_)
func (_TokenList *TokenListTransactorSession) RemoveToken(_token common.Address) (*types.Transaction, error) {
	return _TokenList.Contract.RemoveToken(&_TokenList.TransactOpts, _token)
}

// SetMaxAmount is a paid mutator transaction binding the contract method 0xfb52b065.
//
// Solidity: function setMaxAmount(address _token, uint256 _maxAmount) returns()
func (_TokenList *TokenListTransactor) SetMaxAmount(opts *bind.TransactOpts, _token common.Address, _maxAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "setMaxAmount", _token, _maxAmount)
}

// SetMaxAmount is a paid mutator transaction binding the contract method 0xfb52b065.
//
// Solidity: function setMaxAmount(address _token, uint256 _maxAmount) returns()
func (_TokenList *TokenListSession) SetMaxAmount(_token common.Address, _maxAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.SetMaxAmount(&_TokenList.TransactOpts, _token, _maxAmount)
}

// SetMaxAmount is a paid mutator transaction binding the contract method 0xfb52b065.
//
// Solidity: function setMaxAmount(address _token, uint256 _maxAmount) returns()
func (_TokenList *TokenListTransactorSession) SetMaxAmount(_token common.Address, _maxAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.SetMaxAmount(&_TokenList.TransactOpts, _token, _maxAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0xd5708d5a.
//
// Solidity: function setMinAmount(address _token, uint256 _minAmount) returns()
func (_TokenList *TokenListTransactor) SetMinAmount(opts *bind.TransactOpts, _token common.Address, _minAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "setMinAmount", _token, _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0xd5708d5a.
//
// Solidity: function setMinAmount(address _token, uint256 _minAmount) returns()
func (_TokenList *TokenListSession) SetMinAmount(_token common.Address, _minAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.SetMinAmount(&_TokenList.TransactOpts, _token, _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0xd5708d5a.
//
// Solidity: function setMinAmount(address _token, uint256 _minAmount) returns()
func (_TokenList *TokenListTransactorSession) SetMinAmount(_token common.Address, _minAmount *big.Int) (*types.Transaction, error) {
	return _TokenList.Contract.SetMinAmount(&_TokenList.TransactOpts, _token, _minAmount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenList *TokenListTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenList.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenList *TokenListSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenList.Contract.TransferOwnership(&_TokenList.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenList *TokenListTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenList.Contract.TransferOwnership(&_TokenList.TransactOpts, newOwner)
}

// TokenListOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenList contract.
type TokenListOwnershipTransferredIterator struct {
	Event *TokenListOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenListOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenListOwnershipTransferred)
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
		it.Event = new(TokenListOwnershipTransferred)
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
func (it *TokenListOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenListOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenListOwnershipTransferred represents a OwnershipTransferred event raised by the TokenList contract.
type TokenListOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenList *TokenListFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenListOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenList.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenListOwnershipTransferredIterator{contract: _TokenList.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenList *TokenListFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenListOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenList.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenListOwnershipTransferred)
				if err := _TokenList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenList *TokenListFilterer) ParseOwnershipTransferred(log types.Log) (*TokenListOwnershipTransferred, error) {
	event := new(TokenListOwnershipTransferred)
	if err := _TokenList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenListTokenAddedIterator is returned from FilterTokenAdded and is used to iterate over the raw logs and unpacked data for TokenAdded events raised by the TokenList contract.
type TokenListTokenAddedIterator struct {
	Event *TokenListTokenAdded // Event containing the contract specifics and raw log

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
func (it *TokenListTokenAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenListTokenAdded)
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
		it.Event = new(TokenListTokenAdded)
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
func (it *TokenListTokenAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenListTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenListTokenAdded represents a TokenAdded event raised by the TokenList contract.
type TokenListTokenAdded struct {
	Token     common.Address
	MinAmount *big.Int
	MaxAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokenAdded is a free log retrieval operation binding the contract event 0xa818a22273fc309f0a3682b642c74c5b5c25c0615ff378d07767cd231e19fffc.
//
// Solidity: event TokenAdded(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) FilterTokenAdded(opts *bind.FilterOpts, token []common.Address) (*TokenListTokenAddedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.FilterLogs(opts, "TokenAdded", tokenRule)
	if err != nil {
		return nil, err
	}
	return &TokenListTokenAddedIterator{contract: _TokenList.contract, event: "TokenAdded", logs: logs, sub: sub}, nil
}

// WatchTokenAdded is a free log subscription operation binding the contract event 0xa818a22273fc309f0a3682b642c74c5b5c25c0615ff378d07767cd231e19fffc.
//
// Solidity: event TokenAdded(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) WatchTokenAdded(opts *bind.WatchOpts, sink chan<- *TokenListTokenAdded, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.WatchLogs(opts, "TokenAdded", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenListTokenAdded)
				if err := _TokenList.contract.UnpackLog(event, "TokenAdded", log); err != nil {
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

// ParseTokenAdded is a log parse operation binding the contract event 0xa818a22273fc309f0a3682b642c74c5b5c25c0615ff378d07767cd231e19fffc.
//
// Solidity: event TokenAdded(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) ParseTokenAdded(log types.Log) (*TokenListTokenAdded, error) {
	event := new(TokenListTokenAdded)
	if err := _TokenList.contract.UnpackLog(event, "TokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenListTokenRemovedIterator is returned from FilterTokenRemoved and is used to iterate over the raw logs and unpacked data for TokenRemoved events raised by the TokenList contract.
type TokenListTokenRemovedIterator struct {
	Event *TokenListTokenRemoved // Event containing the contract specifics and raw log

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
func (it *TokenListTokenRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenListTokenRemoved)
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
		it.Event = new(TokenListTokenRemoved)
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
func (it *TokenListTokenRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenListTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenListTokenRemoved represents a TokenRemoved event raised by the TokenList contract.
type TokenListTokenRemoved struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokenRemoved is a free log retrieval operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_TokenList *TokenListFilterer) FilterTokenRemoved(opts *bind.FilterOpts, token []common.Address) (*TokenListTokenRemovedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.FilterLogs(opts, "TokenRemoved", tokenRule)
	if err != nil {
		return nil, err
	}
	return &TokenListTokenRemovedIterator{contract: _TokenList.contract, event: "TokenRemoved", logs: logs, sub: sub}, nil
}

// WatchTokenRemoved is a free log subscription operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_TokenList *TokenListFilterer) WatchTokenRemoved(opts *bind.WatchOpts, sink chan<- *TokenListTokenRemoved, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.WatchLogs(opts, "TokenRemoved", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenListTokenRemoved)
				if err := _TokenList.contract.UnpackLog(event, "TokenRemoved", log); err != nil {
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

// ParseTokenRemoved is a log parse operation binding the contract event 0x4c910b69fe65a61f7531b9c5042b2329ca7179c77290aa7e2eb3afa3c8511fd3.
//
// Solidity: event TokenRemoved(address indexed token)
func (_TokenList *TokenListFilterer) ParseTokenRemoved(log types.Log) (*TokenListTokenRemoved, error) {
	event := new(TokenListTokenRemoved)
	if err := _TokenList.contract.UnpackLog(event, "TokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenListTokenUpdatedIterator is returned from FilterTokenUpdated and is used to iterate over the raw logs and unpacked data for TokenUpdated events raised by the TokenList contract.
type TokenListTokenUpdatedIterator struct {
	Event *TokenListTokenUpdated // Event containing the contract specifics and raw log

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
func (it *TokenListTokenUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenListTokenUpdated)
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
		it.Event = new(TokenListTokenUpdated)
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
func (it *TokenListTokenUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenListTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenListTokenUpdated represents a TokenUpdated event raised by the TokenList contract.
type TokenListTokenUpdated struct {
	Token     common.Address
	MinAmount *big.Int
	MaxAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokenUpdated is a free log retrieval operation binding the contract event 0x5d4804fe0ec949f552f757bfb400c951422d44c10c004077ecd19a9d1f503562.
//
// Solidity: event TokenUpdated(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) FilterTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*TokenListTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.FilterLogs(opts, "TokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &TokenListTokenUpdatedIterator{contract: _TokenList.contract, event: "TokenUpdated", logs: logs, sub: sub}, nil
}

// WatchTokenUpdated is a free log subscription operation binding the contract event 0x5d4804fe0ec949f552f757bfb400c951422d44c10c004077ecd19a9d1f503562.
//
// Solidity: event TokenUpdated(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) WatchTokenUpdated(opts *bind.WatchOpts, sink chan<- *TokenListTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _TokenList.contract.WatchLogs(opts, "TokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenListTokenUpdated)
				if err := _TokenList.contract.UnpackLog(event, "TokenUpdated", log); err != nil {
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

// ParseTokenUpdated is a log parse operation binding the contract event 0x5d4804fe0ec949f552f757bfb400c951422d44c10c004077ecd19a9d1f503562.
//
// Solidity: event TokenUpdated(address indexed token, uint256 minAmount, uint256 maxAmount)
func (_TokenList *TokenListFilterer) ParseTokenUpdated(log types.Log) (*TokenListTokenUpdated, error) {
	event := new(TokenListTokenUpdated)
	if err := _TokenList.contract.UnpackLog(event, "TokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
