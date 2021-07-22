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

// ShadowTokenABI is the input ABI used to generate the binding from.
const ShadowTokenABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_coToken\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"coToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ShadowTokenBin is the compiled bytecode used for deploying new contracts.
var ShadowTokenBin = "0x60806040526003805460ff60a01b191690553480156200001e57600080fd5b506040516200108138038062001081833981810160405260a08110156200004457600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200007057600080fd5b9083019060208201858111156200008657600080fd5b8251640100000000811182820188101715620000a157600080fd5b82525081516020918201929091019080838360005b83811015620000d0578181015183820152602001620000b6565b50505050905090810190601f168015620000fe5780820380516001836020036101000a031916815260200191505b50604052602001805160405193929190846401000000008211156200012257600080fd5b9083019060208201858111156200013857600080fd5b82516401000000008111828201881017156200015357600080fd5b82525081516020918201929091019080838360005b838110156200018257818101518382015260200162000168565b50505050905090810190601f168015620001b05780820380516001836020036101000a031916815260200191505b5060405260209081015160038054336001600160a01b0319918216179091556005805482166001600160a01b038b8116919091179091556004805490921690891617905585519093506200020b925060069186019062000270565b5081516200022190600790602085019062000270565b506008805460ff191660ff83161790556040516001600160a01b038616907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a2505050505062000315565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620002b357805160ff1916838001178555620002e3565b82800160010185558215620002e3579182015b82811115620002e3578251825591602001919060010190620002c6565b50620002f1929150620002f5565b5090565b6200031291905b80821115620002f15760008155600101620002fc565b90565b610d5c80620003256000396000f3fe608060405234801561001057600080fd5b506004361061012c5760003560e01c806366188463116100ad57806395d89b411161007157806395d89b4114610345578063a9059cbb1461034d578063d73dd62314610379578063dd62ed3e146103a5578063f2fde38b146103d35761012c565b806366188463146102db57806370a08231146103075780637f9864f71461032d5780638456cb59146103355780638da5cb5b1461033d5761012c565b8063313ce567116100f4578063313ce567146102625780633f4ba83a1461028057806340c10f191461028a57806342966c68146102b65780635c975abb146102d35761012c565b806306fdde031461013157806307546172146101ae578063095ea7b3146101d257806318160ddd1461021257806323b872dd1461022c575b600080fd5b6101396103f9565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561017357818101518382015260200161015b565b50505050905090810190601f1680156101a05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101b6610487565b604080516001600160a01b039092168252519081900360200190f35b6101fe600480360360408110156101e857600080fd5b506001600160a01b038135169060200135610496565b604080519115158252519081900360200190f35b61021a6104fc565b60408051918252519081900360200190f35b6101fe6004803603606081101561024257600080fd5b506001600160a01b03813581169160208101359091169060400135610502565b61026a610665565b6040805160ff9092168252519081900360200190f35b61028861066e565b005b6101fe600480360360408110156102a057600080fd5b506001600160a01b0381351690602001356106d3565b6101fe600480360360208110156102cc57600080fd5b5035610804565b6101fe6108cb565b6101fe600480360360408110156102f157600080fd5b506001600160a01b0381351690602001356108db565b61021a6004803603602081101561031d57600080fd5b50356001600160a01b03166109cb565b6101b66109e6565b6102886109f5565b6101b6610a61565b610139610a70565b6101fe6004803603604081101561036357600080fd5b506001600160a01b038135169060200135610acb565b6101fe6004803603604081101561038f57600080fd5b506001600160a01b038135169060200135610b98565b61021a600480360360408110156103bb57600080fd5b506001600160a01b0381358116916020013516610c31565b610288600480360360208110156103e957600080fd5b50356001600160a01b0316610c5c565b6006805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561047f5780601f106104545761010080835404028352916020019161047f565b820191906000526020600020905b81548152906001019060200180831161046257829003601f168201915b505050505081565b6005546001600160a01b031681565b3360008181526002602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60015490565b60006001600160a01b03831661051757600080fd5b6001600160a01b03841660009081526020819052604090205482111561053c57600080fd5b6001600160a01b038416600090815260026020908152604080832033845290915290205482111561056c57600080fd5b6001600160a01b038416600090815260208190526040902054610595908363ffffffff610ce216565b6001600160a01b0380861660009081526020819052604080822093909355908516815220546105ca908363ffffffff610cf416565b6001600160a01b0380851660009081526020818152604080832094909455918716815260028252828120338252909152205461060c908363ffffffff610ce216565b6001600160a01b0380861660008181526002602090815260408083203384528252918290209490945580518681529051928716939192600080516020610d08833981519152929181900390910190a35060019392505050565b60085460ff1681565b6003546001600160a01b0316331461068557600080fd5b600354600160a01b900460ff1661069b57600080fd5b6003805460ff60a01b191690556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3390600090a1565b6005546000906001600160a01b03163314610726576040805162461bcd60e51b815260206004820152600e60248201526d3737ba103a34329036b4b73a32b960911b604482015290519081900360640190fd5b600354600160a01b900460ff161561073d57600080fd5b600154610750908363ffffffff610cf416565b6001556001600160a01b03831660009081526020819052604090205461077c908363ffffffff610cf416565b6001600160a01b03841660008181526020818152604091829020939093558051858152905191927f30385c845b448a36257a6a1716e6ad2e1bc2cbe333cde1e69fe849ad6511adfe92918290030190a26040805183815290516001600160a01b03851691600091600080516020610d088339815191529181900360200190a350600192915050565b3360009081526020819052604081205482111561082057600080fd5b600154610833908363ffffffff610ce216565b60015533600090815260208190526040902054610856908363ffffffff610ce216565b3360008181526020818152604091829020939093558051858152905191927f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df792918290030190a26040805183815290516000913391600080516020610d088339815191529181900360200190a3506001919050565b600354600160a01b900460ff1681565b3360009081526002602090815260408083206001600160a01b038616845290915281205480831115610930573360009081526002602090815260408083206001600160a01b0388168452909152812055610965565b610940818463ffffffff610ce216565b3360009081526002602090815260408083206001600160a01b03891684529091529020555b3360008181526002602090815260408083206001600160a01b0389168085529083529281902054815190815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a35060019392505050565b6001600160a01b031660009081526020819052604090205490565b6004546001600160a01b031681565b6003546001600160a01b03163314610a0c57600080fd5b600354600160a01b900460ff1615610a2357600080fd5b6003805460ff60a01b1916600160a01b1790556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62590600090a1565b6003546001600160a01b031681565b6007805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561047f5780601f106104545761010080835404028352916020019161047f565b60006001600160a01b038316610ae057600080fd5b33600090815260208190526040902054821115610afc57600080fd5b33600090815260208190526040902054610b1c908363ffffffff610ce216565b33600090815260208190526040808220929092556001600160a01b03851681522054610b4e908363ffffffff610cf416565b6001600160a01b03841660008181526020818152604091829020939093558051858152905191923392600080516020610d088339815191529281900390910190a350600192915050565b3360009081526002602090815260408083206001600160a01b0386168452909152812054610bcc908363ffffffff610cf416565b3360008181526002602090815260408083206001600160a01b0389168085529083529281902085905580519485525191937f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929081900390910190a350600192915050565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205490565b6003546001600160a01b03163314610c7357600080fd5b6001600160a01b038116610c8657600080fd5b6003546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600380546001600160a01b0319166001600160a01b0392909216919091179055565b600082821115610cee57fe5b50900390565b81810182811015610d0157fe5b9291505056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa265627a7a723158201516428878d48881d042bc960bba7953208f63aa81f6a482cd8527538a45cc0264736f6c63430005110032"

// DeployShadowToken deploys a new Ethereum contract, binding an instance of ShadowToken to it.
func DeployShadowToken(auth *bind.TransactOpts, backend bind.ContractBackend, _minter common.Address, _coToken common.Address, _name string, _symbol string, _decimals uint8) (common.Address, *types.Transaction, *ShadowToken, error) {
	parsed, err := abi.JSON(strings.NewReader(ShadowTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ShadowTokenBin), backend, _minter, _coToken, _name, _symbol, _decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ShadowToken{ShadowTokenCaller: ShadowTokenCaller{contract: contract}, ShadowTokenTransactor: ShadowTokenTransactor{contract: contract}, ShadowTokenFilterer: ShadowTokenFilterer{contract: contract}}, nil
}

// ShadowToken is an auto generated Go binding around an Ethereum contract.
type ShadowToken struct {
	ShadowTokenCaller     // Read-only binding to the contract
	ShadowTokenTransactor // Write-only binding to the contract
	ShadowTokenFilterer   // Log filterer for contract events
}

// ShadowTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type ShadowTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShadowTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ShadowTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShadowTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ShadowTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShadowTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ShadowTokenSession struct {
	Contract     *ShadowToken      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ShadowTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ShadowTokenCallerSession struct {
	Contract *ShadowTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ShadowTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ShadowTokenTransactorSession struct {
	Contract     *ShadowTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ShadowTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type ShadowTokenRaw struct {
	Contract *ShadowToken // Generic contract binding to access the raw methods on
}

// ShadowTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ShadowTokenCallerRaw struct {
	Contract *ShadowTokenCaller // Generic read-only contract binding to access the raw methods on
}

// ShadowTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ShadowTokenTransactorRaw struct {
	Contract *ShadowTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewShadowToken creates a new instance of ShadowToken, bound to a specific deployed contract.
func NewShadowToken(address common.Address, backend bind.ContractBackend) (*ShadowToken, error) {
	contract, err := bindShadowToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ShadowToken{ShadowTokenCaller: ShadowTokenCaller{contract: contract}, ShadowTokenTransactor: ShadowTokenTransactor{contract: contract}, ShadowTokenFilterer: ShadowTokenFilterer{contract: contract}}, nil
}

// NewShadowTokenCaller creates a new read-only instance of ShadowToken, bound to a specific deployed contract.
func NewShadowTokenCaller(address common.Address, caller bind.ContractCaller) (*ShadowTokenCaller, error) {
	contract, err := bindShadowToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenCaller{contract: contract}, nil
}

// NewShadowTokenTransactor creates a new write-only instance of ShadowToken, bound to a specific deployed contract.
func NewShadowTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*ShadowTokenTransactor, error) {
	contract, err := bindShadowToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenTransactor{contract: contract}, nil
}

// NewShadowTokenFilterer creates a new log filterer instance of ShadowToken, bound to a specific deployed contract.
func NewShadowTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*ShadowTokenFilterer, error) {
	contract, err := bindShadowToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenFilterer{contract: contract}, nil
}

// bindShadowToken binds a generic wrapper to an already deployed contract.
func bindShadowToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ShadowTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ShadowToken *ShadowTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ShadowToken.Contract.ShadowTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ShadowToken *ShadowTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShadowToken.Contract.ShadowTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ShadowToken *ShadowTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ShadowToken.Contract.ShadowTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ShadowToken *ShadowTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ShadowToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ShadowToken *ShadowTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShadowToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ShadowToken *ShadowTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ShadowToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_ShadowToken *ShadowTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_ShadowToken *ShadowTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ShadowToken.Contract.Allowance(&_ShadowToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_ShadowToken *ShadowTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ShadowToken.Contract.Allowance(&_ShadowToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_ShadowToken *ShadowTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_ShadowToken *ShadowTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ShadowToken.Contract.BalanceOf(&_ShadowToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_ShadowToken *ShadowTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ShadowToken.Contract.BalanceOf(&_ShadowToken.CallOpts, _owner)
}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_ShadowToken *ShadowTokenCaller) CoToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "coToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_ShadowToken *ShadowTokenSession) CoToken() (common.Address, error) {
	return _ShadowToken.Contract.CoToken(&_ShadowToken.CallOpts)
}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_ShadowToken *ShadowTokenCallerSession) CoToken() (common.Address, error) {
	return _ShadowToken.Contract.CoToken(&_ShadowToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ShadowToken *ShadowTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ShadowToken *ShadowTokenSession) Decimals() (uint8, error) {
	return _ShadowToken.Contract.Decimals(&_ShadowToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ShadowToken *ShadowTokenCallerSession) Decimals() (uint8, error) {
	return _ShadowToken.Contract.Decimals(&_ShadowToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_ShadowToken *ShadowTokenCaller) Minter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "minter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_ShadowToken *ShadowTokenSession) Minter() (common.Address, error) {
	return _ShadowToken.Contract.Minter(&_ShadowToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_ShadowToken *ShadowTokenCallerSession) Minter() (common.Address, error) {
	return _ShadowToken.Contract.Minter(&_ShadowToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ShadowToken *ShadowTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ShadowToken *ShadowTokenSession) Name() (string, error) {
	return _ShadowToken.Contract.Name(&_ShadowToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ShadowToken *ShadowTokenCallerSession) Name() (string, error) {
	return _ShadowToken.Contract.Name(&_ShadowToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ShadowToken *ShadowTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ShadowToken *ShadowTokenSession) Owner() (common.Address, error) {
	return _ShadowToken.Contract.Owner(&_ShadowToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ShadowToken *ShadowTokenCallerSession) Owner() (common.Address, error) {
	return _ShadowToken.Contract.Owner(&_ShadowToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ShadowToken *ShadowTokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ShadowToken *ShadowTokenSession) Paused() (bool, error) {
	return _ShadowToken.Contract.Paused(&_ShadowToken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ShadowToken *ShadowTokenCallerSession) Paused() (bool, error) {
	return _ShadowToken.Contract.Paused(&_ShadowToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ShadowToken *ShadowTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ShadowToken *ShadowTokenSession) Symbol() (string, error) {
	return _ShadowToken.Contract.Symbol(&_ShadowToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ShadowToken *ShadowTokenCallerSession) Symbol() (string, error) {
	return _ShadowToken.Contract.Symbol(&_ShadowToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ShadowToken *ShadowTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ShadowToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ShadowToken *ShadowTokenSession) TotalSupply() (*big.Int, error) {
	return _ShadowToken.Contract.TotalSupply(&_ShadowToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ShadowToken *ShadowTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _ShadowToken.Contract.TotalSupply(&_ShadowToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Approve(&_ShadowToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Approve(&_ShadowToken.TransactOpts, _spender, _value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) Burn(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "burn", _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenSession) Burn(_amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Burn(&_ShadowToken.TransactOpts, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) Burn(_amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Burn(&_ShadowToken.TransactOpts, _amount)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address _spender, uint256 _subtractedValue) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address _spender, uint256 _subtractedValue) returns(bool)
func (_ShadowToken *ShadowTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.DecreaseApproval(&_ShadowToken.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address _spender, uint256 _subtractedValue) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.DecreaseApproval(&_ShadowToken.TransactOpts, _spender, _subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address _spender, uint256 _addedValue) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address _spender, uint256 _addedValue) returns(bool)
func (_ShadowToken *ShadowTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.IncreaseApproval(&_ShadowToken.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address _spender, uint256 _addedValue) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.IncreaseApproval(&_ShadowToken.TransactOpts, _spender, _addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "mint", _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Mint(&_ShadowToken.TransactOpts, _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Mint(&_ShadowToken.TransactOpts, _to, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ShadowToken *ShadowTokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ShadowToken *ShadowTokenSession) Pause() (*types.Transaction, error) {
	return _ShadowToken.Contract.Pause(&_ShadowToken.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ShadowToken *ShadowTokenTransactorSession) Pause() (*types.Transaction, error) {
	return _ShadowToken.Contract.Pause(&_ShadowToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Transfer(&_ShadowToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.Transfer(&_ShadowToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.TransferFrom(&_ShadowToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_ShadowToken *ShadowTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ShadowToken.Contract.TransferFrom(&_ShadowToken.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ShadowToken *ShadowTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ShadowToken *ShadowTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ShadowToken.Contract.TransferOwnership(&_ShadowToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ShadowToken *ShadowTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ShadowToken.Contract.TransferOwnership(&_ShadowToken.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ShadowToken *ShadowTokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShadowToken.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ShadowToken *ShadowTokenSession) Unpause() (*types.Transaction, error) {
	return _ShadowToken.Contract.Unpause(&_ShadowToken.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ShadowToken *ShadowTokenTransactorSession) Unpause() (*types.Transaction, error) {
	return _ShadowToken.Contract.Unpause(&_ShadowToken.TransactOpts)
}

// ShadowTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ShadowToken contract.
type ShadowTokenApprovalIterator struct {
	Event *ShadowTokenApproval // Event containing the contract specifics and raw log

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
func (it *ShadowTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenApproval)
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
		it.Event = new(ShadowTokenApproval)
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
func (it *ShadowTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenApproval represents a Approval event raised by the ShadowToken contract.
type ShadowTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ShadowTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenApprovalIterator{contract: _ShadowToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ShadowTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenApproval)
				if err := _ShadowToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) ParseApproval(log types.Log) (*ShadowTokenApproval, error) {
	event := new(ShadowTokenApproval)
	if err := _ShadowToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenBurnedIterator is returned from FilterBurned and is used to iterate over the raw logs and unpacked data for Burned events raised by the ShadowToken contract.
type ShadowTokenBurnedIterator struct {
	Event *ShadowTokenBurned // Event containing the contract specifics and raw log

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
func (it *ShadowTokenBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenBurned)
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
		it.Event = new(ShadowTokenBurned)
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
func (it *ShadowTokenBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenBurned represents a Burned event raised by the ShadowToken contract.
type ShadowTokenBurned struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurned is a free log retrieval operation binding the contract event 0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7.
//
// Solidity: event Burned(address indexed from, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) FilterBurned(opts *bind.FilterOpts, from []common.Address) (*ShadowTokenBurnedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Burned", fromRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenBurnedIterator{contract: _ShadowToken.contract, event: "Burned", logs: logs, sub: sub}, nil
}

// WatchBurned is a free log subscription operation binding the contract event 0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7.
//
// Solidity: event Burned(address indexed from, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *ShadowTokenBurned, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Burned", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenBurned)
				if err := _ShadowToken.contract.UnpackLog(event, "Burned", log); err != nil {
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

// ParseBurned is a log parse operation binding the contract event 0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7.
//
// Solidity: event Burned(address indexed from, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) ParseBurned(log types.Log) (*ShadowTokenBurned, error) {
	event := new(ShadowTokenBurned)
	if err := _ShadowToken.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the ShadowToken contract.
type ShadowTokenMintedIterator struct {
	Event *ShadowTokenMinted // Event containing the contract specifics and raw log

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
func (it *ShadowTokenMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenMinted)
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
		it.Event = new(ShadowTokenMinted)
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
func (it *ShadowTokenMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenMinted represents a Minted event raised by the ShadowToken contract.
type ShadowTokenMinted struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0x30385c845b448a36257a6a1716e6ad2e1bc2cbe333cde1e69fe849ad6511adfe.
//
// Solidity: event Minted(address indexed to, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) FilterMinted(opts *bind.FilterOpts, to []common.Address) (*ShadowTokenMintedIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Minted", toRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenMintedIterator{contract: _ShadowToken.contract, event: "Minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0x30385c845b448a36257a6a1716e6ad2e1bc2cbe333cde1e69fe849ad6511adfe.
//
// Solidity: event Minted(address indexed to, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *ShadowTokenMinted, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Minted", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenMinted)
				if err := _ShadowToken.contract.UnpackLog(event, "Minted", log); err != nil {
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

// ParseMinted is a log parse operation binding the contract event 0x30385c845b448a36257a6a1716e6ad2e1bc2cbe333cde1e69fe849ad6511adfe.
//
// Solidity: event Minted(address indexed to, uint256 amount)
func (_ShadowToken *ShadowTokenFilterer) ParseMinted(log types.Log) (*ShadowTokenMinted, error) {
	event := new(ShadowTokenMinted)
	if err := _ShadowToken.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the ShadowToken contract.
type ShadowTokenMinterAddedIterator struct {
	Event *ShadowTokenMinterAdded // Event containing the contract specifics and raw log

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
func (it *ShadowTokenMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenMinterAdded)
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
		it.Event = new(ShadowTokenMinterAdded)
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
func (it *ShadowTokenMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenMinterAdded represents a MinterAdded event raised by the ShadowToken contract.
type ShadowTokenMinterAdded struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) FilterMinterAdded(opts *bind.FilterOpts, minter []common.Address) (*ShadowTokenMinterAddedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "MinterAdded", minterRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenMinterAddedIterator{contract: _ShadowToken.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *ShadowTokenMinterAdded, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "MinterAdded", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenMinterAdded)
				if err := _ShadowToken.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// ParseMinterAdded is a log parse operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) ParseMinterAdded(log types.Log) (*ShadowTokenMinterAdded, error) {
	event := new(ShadowTokenMinterAdded)
	if err := _ShadowToken.contract.UnpackLog(event, "MinterAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the ShadowToken contract.
type ShadowTokenMinterRemovedIterator struct {
	Event *ShadowTokenMinterRemoved // Event containing the contract specifics and raw log

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
func (it *ShadowTokenMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenMinterRemoved)
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
		it.Event = new(ShadowTokenMinterRemoved)
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
func (it *ShadowTokenMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenMinterRemoved represents a MinterRemoved event raised by the ShadowToken contract.
type ShadowTokenMinterRemoved struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) FilterMinterRemoved(opts *bind.FilterOpts, minter []common.Address) (*ShadowTokenMinterRemovedIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "MinterRemoved", minterRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenMinterRemovedIterator{contract: _ShadowToken.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *ShadowTokenMinterRemoved, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "MinterRemoved", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenMinterRemoved)
				if err := _ShadowToken.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed minter)
func (_ShadowToken *ShadowTokenFilterer) ParseMinterRemoved(log types.Log) (*ShadowTokenMinterRemoved, error) {
	event := new(ShadowTokenMinterRemoved)
	if err := _ShadowToken.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ShadowToken contract.
type ShadowTokenOwnershipTransferredIterator struct {
	Event *ShadowTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ShadowTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenOwnershipTransferred)
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
		it.Event = new(ShadowTokenOwnershipTransferred)
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
func (it *ShadowTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenOwnershipTransferred represents a OwnershipTransferred event raised by the ShadowToken contract.
type ShadowTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ShadowToken *ShadowTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ShadowTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenOwnershipTransferredIterator{contract: _ShadowToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ShadowToken *ShadowTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ShadowTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenOwnershipTransferred)
				if err := _ShadowToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ShadowToken *ShadowTokenFilterer) ParseOwnershipTransferred(log types.Log) (*ShadowTokenOwnershipTransferred, error) {
	event := new(ShadowTokenOwnershipTransferred)
	if err := _ShadowToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the ShadowToken contract.
type ShadowTokenPauseIterator struct {
	Event *ShadowTokenPause // Event containing the contract specifics and raw log

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
func (it *ShadowTokenPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenPause)
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
		it.Event = new(ShadowTokenPause)
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
func (it *ShadowTokenPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenPause represents a Pause event raised by the ShadowToken contract.
type ShadowTokenPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_ShadowToken *ShadowTokenFilterer) FilterPause(opts *bind.FilterOpts) (*ShadowTokenPauseIterator, error) {

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &ShadowTokenPauseIterator{contract: _ShadowToken.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_ShadowToken *ShadowTokenFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *ShadowTokenPause) (event.Subscription, error) {

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenPause)
				if err := _ShadowToken.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_ShadowToken *ShadowTokenFilterer) ParsePause(log types.Log) (*ShadowTokenPause, error) {
	event := new(ShadowTokenPause)
	if err := _ShadowToken.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ShadowToken contract.
type ShadowTokenTransferIterator struct {
	Event *ShadowTokenTransfer // Event containing the contract specifics and raw log

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
func (it *ShadowTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenTransfer)
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
		it.Event = new(ShadowTokenTransfer)
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
func (it *ShadowTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenTransfer represents a Transfer event raised by the ShadowToken contract.
type ShadowTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ShadowTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ShadowTokenTransferIterator{contract: _ShadowToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ShadowTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenTransfer)
				if err := _ShadowToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ShadowToken *ShadowTokenFilterer) ParseTransfer(log types.Log) (*ShadowTokenTransfer, error) {
	event := new(ShadowTokenTransfer)
	if err := _ShadowToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ShadowTokenUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the ShadowToken contract.
type ShadowTokenUnpauseIterator struct {
	Event *ShadowTokenUnpause // Event containing the contract specifics and raw log

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
func (it *ShadowTokenUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ShadowTokenUnpause)
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
		it.Event = new(ShadowTokenUnpause)
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
func (it *ShadowTokenUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ShadowTokenUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ShadowTokenUnpause represents a Unpause event raised by the ShadowToken contract.
type ShadowTokenUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_ShadowToken *ShadowTokenFilterer) FilterUnpause(opts *bind.FilterOpts) (*ShadowTokenUnpauseIterator, error) {

	logs, sub, err := _ShadowToken.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &ShadowTokenUnpauseIterator{contract: _ShadowToken.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_ShadowToken *ShadowTokenFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *ShadowTokenUnpause) (event.Subscription, error) {

	logs, sub, err := _ShadowToken.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ShadowTokenUnpause)
				if err := _ShadowToken.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_ShadowToken *ShadowTokenFilterer) ParseUnpause(log types.Log) (*ShadowTokenUnpause, error) {
	event := new(ShadowTokenUnpause)
	if err := _ShadowToken.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
