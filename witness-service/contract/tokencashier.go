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

// TokenCashierABI is the input ABI used to generate the binding from.
const TokenCashierABI = "[{\"inputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"_wrappedCoin\",\"type\":\"address\"},{\"internalType\":\"contractITokenList[]\",\"name\":\"_tokenLists\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_tokenSafes\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"Receipt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"counts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractITokenList\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenSafes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wrappedCoin\",\"outputs\":[{\"internalType\":\"contractIWrappedCoin\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenCashierBin is the compiled bytecode used for deploying new contracts.
var TokenCashierBin = "0x60806040526000805460ff60a01b191690553480156200001e57600080fd5b50604051620014b5380380620014b5833981810160405260608110156200004457600080fd5b8151602083018051604051929492938301929190846401000000008211156200006c57600080fd5b9083019060208201858111156200008257600080fd5b8251866020820283011164010000000082111715620000a057600080fd5b82525081516020918201928201910280838360005b83811015620000cf578181015183820152602001620000b5565b5050505090500160405260200180516040519392919084640100000000821115620000f957600080fd5b9083019060208201858111156200010f57600080fd5b82518660208202830111640100000000821117156200012d57600080fd5b82525081516020918201928201910280838360005b838110156200015c57818101518382015260200162000142565b50505050919091016040525050600080546001600160a01b0319163317905550508051825114620001bf5760405162461bcd60e51b815260040180806020018281038252602b8152602001806200148a602b913960400191505060405180910390fd5b600580546001600160a01b0319166001600160a01b0385161790558151620001ef9060019060208501906200020f565b508051620002059060029060208401906200020f565b50505050620002a3565b82805482825590600052602060002090810192821562000267579160200282015b828111156200026757825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019062000230565b506200027592915062000279565b5090565b620002a091905b80821115620002755780546001600160a01b031916815560010162000280565b90565b6111d780620002b36000396000f3fe6080604052600436106100f35760003560e01c80635c975abb1161008a578063894760691161005957806389476069146102d75780638da5cb5b1461030a578063f213159c1461031f578063f2fde38b14610355576100f3565b80635c975abb1461025a57806367a527931461028357806384378ec6146102985780638456cb59146102c2576100f3565b80633f4ba83a116100c65780633f4ba83a146101da57806347e7ef24146101ef578063490ae2101461021b578063527ba9af14610245576100f3565b80630568e65e1461010557806305d85eda1461014a5780631cb928a91461017d5780633ccfd60b146101c3575b3480156100ff57600080fd5b50600080fd5b34801561011157600080fd5b506101386004803603602081101561012857600080fd5b50356001600160a01b0316610388565b60408051918252519081900360200190f35b34801561015657600080fd5b506101386004803603602081101561016d57600080fd5b50356001600160a01b031661039a565b34801561018957600080fd5b506101a7600480360360208110156101a057600080fd5b50356103b5565b604080516001600160a01b039092168252519081900360200190f35b3480156101cf57600080fd5b506101d86103dc565b005b3480156101e657600080fd5b506101d8610422565b6101d86004803603604081101561020557600080fd5b506001600160a01b038135169060200135610485565b34801561022757600080fd5b506101d86004803603602081101561023e57600080fd5b5035610494565b34801561025157600080fd5b506101a76104b0565b34801561026657600080fd5b5061026f6104bf565b604080519115158252519081900360200190f35b34801561028f57600080fd5b506101386104cf565b3480156102a457600080fd5b506101a7600480360360208110156102bb57600080fd5b50356104d5565b3480156102ce57600080fd5b506101d86104e2565b3480156102e357600080fd5b506101d8600480360360208110156102fa57600080fd5b50356001600160a01b031661054c565b34801561031657600080fd5b506101a761071b565b6101d86004803603606081101561033557600080fd5b506001600160a01b0381358116916020810135909116906040013561072a565b34801561036157600080fd5b506101d86004803603602081101561037857600080fd5b50356001600160a01b0316610eb7565b60036020526000908152604090205481565b6001600160a01b031660009081526003602052604090205490565b600181815481106103c257fe5b6000918252602090912001546001600160a01b0316905081565b6000546001600160a01b031633146103f357600080fd5b60405133904780156108fc02916000818181858888f1935050505015801561041f573d6000803e3d6000fd5b50565b6000546001600160a01b0316331461043957600080fd5b600054600160a01b900460ff1661044f57600080fd5b6000805460ff60a01b191681556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b339190a1565b61049082338361072a565b5050565b6000546001600160a01b031633146104ab57600080fd5b600455565b6005546001600160a01b031681565b600054600160a01b900460ff1681565b60045481565b600281815481106103c257fe5b6000546001600160a01b031633146104f957600080fd5b600054600160a01b900460ff161561051057600080fd5b6000805460ff60a01b1916600160a01b1781556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff6259190a1565b6000546001600160a01b0316331461056357600080fd5b604080513060248083019190915282518083039091018152604490910182526020810180516001600160e01b03166370a0823160e01b178152915181516000936060936001600160a01b038716939092909182918083835b602083106105da5780518252601f1990920191602091820191016105bb565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461063c576040519150601f19603f3d011682016040523d82523d6000602084013e610641565b606091505b509150915081610698576040805162461bcd60e51b815260206004820152601860248201527f6661696c656420746f2063616c6c2062616c616e63654f660000000000000000604482015290519081900360640190fd5b60008180602001905160208110156106af57600080fd5b505190508015610715576106c4843383610f3c565b610715576040805162461bcd60e51b815260206004820152601860248201527f6661696c656420746f20776974686472617720746f6b656e0000000000000000604482015290519081900360640190fd5b50505050565b6000546001600160a01b031681565b600054600160a01b900460ff161561074157600080fd5b6001600160a01b038216610792576040805162461bcd60e51b815260206004820152601360248201527234b73b30b634b2103232b9ba34b730ba34b7b760691b604482015290519081900360640190fd5b6000346001600160a01b03851661087257823410156107f1576040805162461bcd60e51b8152602060048201526016602482015275696e73756666696369656e74206d73672e76616c756560501b604482015290519081900360640190fd5b8234039050600560009054906101000a90046001600160a01b03166001600160a01b031663d0e30db0846040518263ffffffff1660e01b81526004016000604051808303818588803b15801561084657600080fd5b505af115801561085a573d6000803e3d6000fd5b50506005546001600160a01b03169750600194505050505b6004548110156108bc576040805162461bcd60e51b815260206004820152601060248201526f696e73756666696369656e742066656560801b604482015290519081900360640190fd5b60005b600154811015610e6457600181815481106108d657fe5b60009182526020808320909101546040805163babcc53960e01b81526001600160a01b038b811660048301529151919092169363babcc53993602480850194919392918390030190829087803b15801561092f57600080fd5b505af1158015610943573d6000803e3d6000fd5b505050506040513d602081101561095957600080fd5b505115610e5c576001818154811061096d57fe5b600091825260208083209091015460408051634d0a32db60e01b81526001600160a01b038b8116600483015291519190921693634d0a32db93602480850194919392918390030190829087803b1580156109c657600080fd5b505af11580156109da573d6000803e3d6000fd5b505050506040513d60208110156109f057600080fd5b5051841015610a37576040805162461bcd60e51b815260206004820152600e60248201526d616d6f756e7420746f6f206c6f7760901b604482015290519081900360640190fd5b60018181548110610a4457fe5b600091825260208083209091015460408051632537b82960e21b81526001600160a01b038b81166004830152915191909216936394dee0a493602480850194919392918390030190829087803b158015610a9d57600080fd5b505af1158015610ab1573d6000803e3d6000fd5b505050506040513d6020811015610ac757600080fd5b5051841115610b0f576040805162461bcd60e51b815260206004820152600f60248201526e0c2dadeeadce840e8dede40d0d2ced608b1b604482015290519081900360640190fd5b60006001600160a01b031660028281548110610b2757fe5b6000918252602090912001546001600160a01b03161415610ce95782158015610b575750610b578633308761105a565b610b925760405162461bcd60e51b81526004018080602001828103825260218152602001806111826021913960400191505060405180910390fd5b60408051602480820187905282518083039091018152604490910182526020810180516001600160e01b0316630852cd8d60e31b178152915181516000936060936001600160a01b038c16939092909182918083835b60208310610c075780518252601f199092019160209182019101610be8565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610c69576040519150601f19603f3d011682016040523d82523d6000602084013e610c6e565b606091505b5091509150818015610c9c575080511580610c9c5750808060200190516020811015610c9957600080fd5b50515b610ce2576040805162461bcd60e51b81526020600482015260126024820152713330b4b6103a3790313ab937103a37b5b2b760711b604482015290519081900360640190fd5b5050610de3565b8215610d6b57610d1b8660028381548110610d0057fe5b6000918252602090912001546001600160a01b031686610f3c565b610d66576040805162461bcd60e51b81526020600482015260176024820152766661696c656420746f2070757420696e746f207361666560481b604482015290519081900360640190fd5b610de3565b610d98863360028481548110610d7d57fe5b6000918252602090912001546001600160a01b03168761105a565b610de3576040805162461bcd60e51b81526020600482015260176024820152766661696c656420746f2070757420696e746f207361666560481b604482015290519081900360640190fd5b6001600160a01b038087166000818152600360209081526040918290208054600101908190558251338152948a169185019190915283820188905260608401869052905190927f85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919919081900360800190a3505050610eb2565b6001016108bf565b506040805162461bcd60e51b815260206004820152601760248201527f6e6f7420612077686974656c697374656420746f6b656e000000000000000000604482015290519081900360640190fd5b505050565b6000546001600160a01b03163314610ece57600080fd5b6001600160a01b038116610ee157600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663a9059cbb60e01b1781529251825160009485946060948a16939092909182918083835b60208310610fbb5780518252601f199092019160209182019101610f9c565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461101d576040519150601f19603f3d011682016040523d82523d6000602084013e611022565b606091505b5091509150818015611050575080511580611050575080806020019051602081101561104d57600080fd5b50515b9695505050505050565b604080516001600160a01b0385811660248301528481166044830152606480830185905283518084039091018152608490920183526020820180516001600160e01b03166323b872dd60e01b1781529251825160009485946060948b16939092909182918083835b602083106110e15780518252601f1990920191602091820191016110c2565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611143576040519150601f19603f3d011682016040523d82523d6000602084013e611148565b606091505b5091509150818015611176575080511580611176575080806020019051602081101561117357600080fd5b50515b97965050505050505056fe6661696c20746f207472616e7366657220746f6b656e20746f2063617368696572a265627a7a72315820ef548558871498992fc6393f52a77cd45975f8bb4729b3279ff641a12f8a82bd64736f6c6343000511003223206f6620746f6b656e206c69737473206973206e6f7420657175616c20746f2023206f66207361666573"

// DeployTokenCashier deploys a new Ethereum contract, binding an instance of TokenCashier to it.
func DeployTokenCashier(auth *bind.TransactOpts, backend bind.ContractBackend, _wrappedCoin common.Address, _tokenLists []common.Address, _tokenSafes []common.Address) (common.Address, *types.Transaction, *TokenCashier, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenCashierABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenCashierBin), backend, _wrappedCoin, _tokenLists, _tokenSafes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenCashier{TokenCashierCaller: TokenCashierCaller{contract: contract}, TokenCashierTransactor: TokenCashierTransactor{contract: contract}, TokenCashierFilterer: TokenCashierFilterer{contract: contract}}, nil
}

// TokenCashier is an auto generated Go binding around an Ethereum contract.
type TokenCashier struct {
	TokenCashierCaller     // Read-only binding to the contract
	TokenCashierTransactor // Write-only binding to the contract
	TokenCashierFilterer   // Log filterer for contract events
}

// TokenCashierCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCashierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenCashierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenCashierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenCashierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenCashierSession struct {
	Contract     *TokenCashier     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCashierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCashierCallerSession struct {
	Contract *TokenCashierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenCashierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenCashierTransactorSession struct {
	Contract     *TokenCashierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenCashierRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenCashierRaw struct {
	Contract *TokenCashier // Generic contract binding to access the raw methods on
}

// TokenCashierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCashierCallerRaw struct {
	Contract *TokenCashierCaller // Generic read-only contract binding to access the raw methods on
}

// TokenCashierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenCashierTransactorRaw struct {
	Contract *TokenCashierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCashier creates a new instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashier(address common.Address, backend bind.ContractBackend) (*TokenCashier, error) {
	contract, err := bindTokenCashier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenCashier{TokenCashierCaller: TokenCashierCaller{contract: contract}, TokenCashierTransactor: TokenCashierTransactor{contract: contract}, TokenCashierFilterer: TokenCashierFilterer{contract: contract}}, nil
}

// NewTokenCashierCaller creates a new read-only instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierCaller(address common.Address, caller bind.ContractCaller) (*TokenCashierCaller, error) {
	contract, err := bindTokenCashier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierCaller{contract: contract}, nil
}

// NewTokenCashierTransactor creates a new write-only instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenCashierTransactor, error) {
	contract, err := bindTokenCashier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCashierTransactor{contract: contract}, nil
}

// NewTokenCashierFilterer creates a new log filterer instance of TokenCashier, bound to a specific deployed contract.
func NewTokenCashierFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenCashierFilterer, error) {
	contract, err := bindTokenCashier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenCashierFilterer{contract: contract}, nil
}

// bindTokenCashier binds a generic wrapper to an already deployed contract.
func bindTokenCashier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenCashierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashier *TokenCashierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashier.Contract.TokenCashierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashier *TokenCashierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.Contract.TokenCashierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashier *TokenCashierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashier.Contract.TokenCashierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenCashier *TokenCashierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenCashier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenCashier *TokenCashierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenCashier *TokenCashierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenCashier.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierCaller) Count(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "count", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Count(&_TokenCashier.CallOpts, _token)
}

// Count is a free data retrieval call binding the contract method 0x05d85eda.
//
// Solidity: function count(address _token) view returns(uint256)
func (_TokenCashier *TokenCashierCallerSession) Count(_token common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Count(&_TokenCashier.CallOpts, _token)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashier *TokenCashierCaller) Counts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "counts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashier *TokenCashierSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Counts(&_TokenCashier.CallOpts, arg0)
}

// Counts is a free data retrieval call binding the contract method 0x0568e65e.
//
// Solidity: function counts(address ) view returns(uint256)
func (_TokenCashier *TokenCashierCallerSession) Counts(arg0 common.Address) (*big.Int, error) {
	return _TokenCashier.Contract.Counts(&_TokenCashier.CallOpts, arg0)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "depositFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierSession) DepositFee() (*big.Int, error) {
	return _TokenCashier.Contract.DepositFee(&_TokenCashier.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() view returns(uint256)
func (_TokenCashier *TokenCashierCallerSession) DepositFee() (*big.Int, error) {
	return _TokenCashier.Contract.DepositFee(&_TokenCashier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierSession) Owner() (common.Address, error) {
	return _TokenCashier.Contract.Owner(&_TokenCashier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenCashier *TokenCashierCallerSession) Owner() (common.Address, error) {
	return _TokenCashier.Contract.Owner(&_TokenCashier.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierSession) Paused() (bool, error) {
	return _TokenCashier.Contract.Paused(&_TokenCashier.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_TokenCashier *TokenCashierCallerSession) Paused() (bool, error) {
	return _TokenCashier.Contract.Paused(&_TokenCashier.CallOpts)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierCaller) TokenLists(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "tokenLists", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashier.Contract.TokenLists(&_TokenCashier.CallOpts, arg0)
}

// TokenLists is a free data retrieval call binding the contract method 0x1cb928a9.
//
// Solidity: function tokenLists(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierCallerSession) TokenLists(arg0 *big.Int) (common.Address, error) {
	return _TokenCashier.Contract.TokenLists(&_TokenCashier.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierCaller) TokenSafes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "tokenSafes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashier.Contract.TokenSafes(&_TokenCashier.CallOpts, arg0)
}

// TokenSafes is a free data retrieval call binding the contract method 0x84378ec6.
//
// Solidity: function tokenSafes(uint256 ) view returns(address)
func (_TokenCashier *TokenCashierCallerSession) TokenSafes(arg0 *big.Int) (common.Address, error) {
	return _TokenCashier.Contract.TokenSafes(&_TokenCashier.CallOpts, arg0)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashier *TokenCashierCaller) WrappedCoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenCashier.contract.Call(opts, &out, "wrappedCoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashier *TokenCashierSession) WrappedCoin() (common.Address, error) {
	return _TokenCashier.Contract.WrappedCoin(&_TokenCashier.CallOpts)
}

// WrappedCoin is a free data retrieval call binding the contract method 0x527ba9af.
//
// Solidity: function wrappedCoin() view returns(address)
func (_TokenCashier *TokenCashierCallerSession) WrappedCoin() (common.Address, error) {
	return _TokenCashier.Contract.WrappedCoin(&_TokenCashier.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactor) Deposit(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "deposit", _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierSession) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.Deposit(&_TokenCashier.TransactOpts, _token, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address _token, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactorSession) Deposit(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.Deposit(&_TokenCashier.TransactOpts, _token, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactor) DepositTo(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "depositTo", _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.DepositTo(&_TokenCashier.TransactOpts, _token, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xf213159c.
//
// Solidity: function depositTo(address _token, address _to, uint256 _amount) payable returns()
func (_TokenCashier *TokenCashierTransactorSession) DepositTo(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.DepositTo(&_TokenCashier.TransactOpts, _token, _to, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierSession) Pause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Pause(&_TokenCashier.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_TokenCashier *TokenCashierTransactorSession) Pause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Pause(&_TokenCashier.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.SetDepositFee(&_TokenCashier.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(uint256 _fee) returns()
func (_TokenCashier *TokenCashierTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _TokenCashier.Contract.SetDepositFee(&_TokenCashier.TransactOpts, _fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.TransferOwnership(&_TokenCashier.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenCashier *TokenCashierTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.TransferOwnership(&_TokenCashier.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierSession) Unpause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Unpause(&_TokenCashier.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_TokenCashier *TokenCashierTransactorSession) Unpause() (*types.Transaction, error) {
	return _TokenCashier.Contract.Unpause(&_TokenCashier.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashier *TokenCashierTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashier *TokenCashierSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashier.Contract.Withdraw(&_TokenCashier.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_TokenCashier *TokenCashierTransactorSession) Withdraw() (*types.Transaction, error) {
	return _TokenCashier.Contract.Withdraw(&_TokenCashier.TransactOpts)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashier *TokenCashierTransactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _TokenCashier.contract.Transact(opts, "withdrawToken", _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashier *TokenCashierSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.WithdrawToken(&_TokenCashier.TransactOpts, _token)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x89476069.
//
// Solidity: function withdrawToken(address _token) returns()
func (_TokenCashier *TokenCashierTransactorSession) WithdrawToken(_token common.Address) (*types.Transaction, error) {
	return _TokenCashier.Contract.WithdrawToken(&_TokenCashier.TransactOpts, _token)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.Contract.Fallback(&_TokenCashier.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_TokenCashier *TokenCashierTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenCashier.Contract.Fallback(&_TokenCashier.TransactOpts, calldata)
}

// TokenCashierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenCashier contract.
type TokenCashierOwnershipTransferredIterator struct {
	Event *TokenCashierOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenCashierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierOwnershipTransferred)
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
		it.Event = new(TokenCashierOwnershipTransferred)
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
func (it *TokenCashierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierOwnershipTransferred represents a OwnershipTransferred event raised by the TokenCashier contract.
type TokenCashierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashier *TokenCashierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenCashierOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierOwnershipTransferredIterator{contract: _TokenCashier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenCashier *TokenCashierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenCashierOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierOwnershipTransferred)
				if err := _TokenCashier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParseOwnershipTransferred(log types.Log) (*TokenCashierOwnershipTransferred, error) {
	event := new(TokenCashierOwnershipTransferred)
	if err := _TokenCashier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the TokenCashier contract.
type TokenCashierPauseIterator struct {
	Event *TokenCashierPause // Event containing the contract specifics and raw log

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
func (it *TokenCashierPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierPause)
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
		it.Event = new(TokenCashierPause)
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
func (it *TokenCashierPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierPause represents a Pause event raised by the TokenCashier contract.
type TokenCashierPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashier *TokenCashierFilterer) FilterPause(opts *bind.FilterOpts) (*TokenCashierPauseIterator, error) {

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierPauseIterator{contract: _TokenCashier.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_TokenCashier *TokenCashierFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *TokenCashierPause) (event.Subscription, error) {

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierPause)
				if err := _TokenCashier.contract.UnpackLog(event, "Pause", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParsePause(log types.Log) (*TokenCashierPause, error) {
	event := new(TokenCashierPause)
	if err := _TokenCashier.contract.UnpackLog(event, "Pause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierReceiptIterator is returned from FilterReceipt and is used to iterate over the raw logs and unpacked data for Receipt events raised by the TokenCashier contract.
type TokenCashierReceiptIterator struct {
	Event *TokenCashierReceipt // Event containing the contract specifics and raw log

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
func (it *TokenCashierReceiptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierReceipt)
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
		it.Event = new(TokenCashierReceipt)
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
func (it *TokenCashierReceiptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierReceiptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierReceipt represents a Receipt event raised by the TokenCashier contract.
type TokenCashierReceipt struct {
	Token     common.Address
	Id        *big.Int
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Fee       *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReceipt is a free log retrieval operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) FilterReceipt(opts *bind.FilterOpts, token []common.Address, id []*big.Int) (*TokenCashierReceiptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return &TokenCashierReceiptIterator{contract: _TokenCashier.contract, event: "Receipt", logs: logs, sub: sub}, nil
}

// WatchReceipt is a free log subscription operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) WatchReceipt(opts *bind.WatchOpts, sink chan<- *TokenCashierReceipt, token []common.Address, id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Receipt", tokenRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierReceipt)
				if err := _TokenCashier.contract.UnpackLog(event, "Receipt", log); err != nil {
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

// ParseReceipt is a log parse operation binding the contract event 0x85425e130ee5cbf9eea6de0d309f1fdd5f7a343aeb20ad4263f3e1305fd5b919.
//
// Solidity: event Receipt(address indexed token, uint256 indexed id, address sender, address recipient, uint256 amount, uint256 fee)
func (_TokenCashier *TokenCashierFilterer) ParseReceipt(log types.Log) (*TokenCashierReceipt, error) {
	event := new(TokenCashierReceipt)
	if err := _TokenCashier.contract.UnpackLog(event, "Receipt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenCashierUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the TokenCashier contract.
type TokenCashierUnpauseIterator struct {
	Event *TokenCashierUnpause // Event containing the contract specifics and raw log

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
func (it *TokenCashierUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenCashierUnpause)
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
		it.Event = new(TokenCashierUnpause)
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
func (it *TokenCashierUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenCashierUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenCashierUnpause represents a Unpause event raised by the TokenCashier contract.
type TokenCashierUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashier *TokenCashierFilterer) FilterUnpause(opts *bind.FilterOpts) (*TokenCashierUnpauseIterator, error) {

	logs, sub, err := _TokenCashier.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &TokenCashierUnpauseIterator{contract: _TokenCashier.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_TokenCashier *TokenCashierFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *TokenCashierUnpause) (event.Subscription, error) {

	logs, sub, err := _TokenCashier.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenCashierUnpause)
				if err := _TokenCashier.contract.UnpackLog(event, "Unpause", log); err != nil {
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
func (_TokenCashier *TokenCashierFilterer) ParseUnpause(log types.Log) (*TokenCashierUnpause, error) {
	event := new(TokenCashierUnpause)
	if err := _TokenCashier.contract.UnpackLog(event, "Unpause", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}