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

// CrosschainERC20ABI is the input ABI used to generate the binding from.
const CrosschainERC20ABI = "[{\"inputs\":[{\"internalType\":\"contractERC20\",\"name\":\"_coToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_minter\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MinterSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coToken\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMinter\",\"type\":\"address\"}],\"name\":\"transferMintership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// CrosschainERC20Bin is the compiled bytecode used for deploying new contracts.
var CrosschainERC20Bin = "0x60806040523480156200001157600080fd5b50604051620019d8380380620019d8833981810160405260a08110156200003757600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200006357600080fd5b9083019060208201858111156200007957600080fd5b82516401000000008111828201881017156200009457600080fd5b82525081516020918201929091019080838360005b83811015620000c3578181015183820152602001620000a9565b50505050905090810190601f168015620000f15780820380516001836020036101000a031916815260200191505b50604052602001805160405193929190846401000000008211156200011557600080fd5b9083019060208201858111156200012b57600080fd5b82516401000000008111828201881017156200014657600080fd5b82525081516020918201929091019080838360005b83811015620001755781810151838201526020016200015b565b50505050905090810190601f168015620001a35780820380516001836020036101000a031916815260200191505b5060405260209081015185519093508592508491620001c89160039185019062000289565b508051620001de90600490602084019062000289565b505060058054601260ff1990911617610100600160a81b0319166101006001600160a01b038981169190910291909117909155600680546001600160a01b03191691871691909117905550620002348162000273565b6040516001600160a01b038516907f726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b690600090a2505050505062000335565b6005805460ff191660ff92909216919091179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620002c157600085556200030c565b82601f10620002dc57805160ff19168380011785556200030c565b828001600101855582156200030c579182015b828111156200030c578251825591602001919060010190620002ef565b506200031a9291506200031e565b5090565b5b808211156200031a57600081556001016200031f565b61169380620003456000396000f3fe608060405234801561001057600080fd5b50600436106101775760003560e01c806342966c68116100d8578063a457c2d71161008c578063cf86a95a11610066578063cf86a95a14610462578063dd62ed3e14610488578063ffaad6a5146104b657610177565b8063a457c2d7146103ed578063a9059cbb14610419578063b6b55f251461044557610177565b806379cc6790116100bd57806379cc6790146103b15780637f9864f7146103dd57806395d89b41146103e557610177565b806342966c681461036e57806370a082311461038b57610177565b806323b872dd1161012f578063313ce56711610114578063313ce567146102f8578063395093511461031657806340c10f191461034257610177565b806323b872dd146102a55780632e1a7d4d146102db57610177565b8063095ea7b311610160578063095ea7b31461021d57806318160ddd1461025d578063205c28781461027757610177565b806306fdde031461017c57806307546172146101f9575b600080fd5b6101846104e2565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101be5781810151838201526020016101a6565b50505050905090810190601f1680156101eb5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610201610596565b604080516001600160a01b039092168252519081900360200190f35b6102496004803603604081101561023357600080fd5b506001600160a01b0381351690602001356105a5565b604080519115158252519081900360200190f35b6102656105c2565b60408051918252519081900360200190f35b6102a36004803603604081101561028d57600080fd5b506001600160a01b0381351690602001356105c8565b005b610249600480360360608110156102bb57600080fd5b506001600160a01b038135811691602081013590911690604001356106a6565b6102a3600480360360208110156102f157600080fd5b503561072e565b61030061073b565b6040805160ff9092168252519081900360200190f35b6102496004803603604081101561032c57600080fd5b506001600160a01b038135169060200135610744565b6102496004803603604081101561035857600080fd5b506001600160a01b038135169060200135610792565b6102a36004803603602081101561038457600080fd5b5035610850565b610265600480360360208110156103a157600080fd5b50356001600160a01b0316610861565b6102a3600480360360408110156103c757600080fd5b506001600160a01b03813516906020013561087c565b6102016108d6565b6101846108ea565b6102496004803603604081101561040357600080fd5b506001600160a01b038135169060200135610969565b6102496004803603604081101561042f57600080fd5b506001600160a01b0381351690602001356109d1565b6102a36004803603602081101561045b57600080fd5b50356109e5565b6102a36004803603602081101561047857600080fd5b50356001600160a01b03166109ef565b6102656004803603604081101561049e57600080fd5b506001600160a01b0381358116916020013516610ab0565b6102a3600480360360408110156104cc57600080fd5b506001600160a01b038135169060200135610adb565b60038054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561058c5780601f106105615761010080835404028352916020019161058c565b820191906000526020600020905b81548152906001019060200180831161056f57829003601f168201915b5050505050905090565b6006546001600160a01b031681565b60006105b96105b2610b64565b8484610b68565b50600192915050565b60025490565b60055461010090046001600160a01b031661062a576040805162461bcd60e51b815260206004820152600b60248201527f6e6f20636f2d746f6b656e000000000000000000000000000000000000000000604482015290519081900360640190fd5b8061067c576040805162461bcd60e51b815260206004820152600b60248201527f616d6f756e742069732030000000000000000000000000000000000000000000604482015290519081900360640190fd5b6106863382610c54565b6005546106a29061010090046001600160a01b03168383610d50565b5050565b60006106b3848484610dd0565b610723846106bf610b64565b61071e85604051806060016040528060288152602001611559602891396001600160a01b038a166000908152600160205260408120906106fd610b64565b6001600160a01b031681526020810191909152604001600020549190610f2b565b610b68565b5060015b9392505050565b61073833826105c8565b50565b60055460ff1690565b60006105b9610751610b64565b8461071e8560016000610762610b64565b6001600160a01b03908116825260208083019390935260409182016000908120918c168152925290205490610fc2565b6006546000906001600160a01b031633146107f4576040805162461bcd60e51b815260206004820152600e60248201527f6e6f7420746865206d696e746572000000000000000000000000000000000000604482015290519081900360640190fd5b81610846576040805162461bcd60e51b815260206004820152600b60248201527f616d6f756e742069732030000000000000000000000000000000000000000000604482015290519081900360640190fd5b6105b9838361101c565b61073861085b610b64565b82610c54565b6001600160a01b031660009081526020819052604090205490565b60006108b382604051806060016040528060248152602001611581602491396108ac866108a7610b64565b610ab0565b9190610f2b565b90506108c7836108c1610b64565b83610b68565b6108d18383610c54565b505050565b60055461010090046001600160a01b031681565b60048054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561058c5780601f106105615761010080835404028352916020019161058c565b60006105b9610976610b64565b8461071e8560405180606001604052806025815260200161163960259139600160006109a0610b64565b6001600160a01b03908116825260208083019390935260409182016000908120918d16815292529020549190610f2b565b60006105b96109de610b64565b8484610dd0565b6107383382610adb565b6006546001600160a01b03163314610a4e576040805162461bcd60e51b815260206004820152600e60248201527f6e6f7420746865206d696e746572000000000000000000000000000000000000604482015290519081900360640190fd5b600680547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383169081179091556040517f726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b690600090a250565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60055461010090046001600160a01b0316610b3d576040805162461bcd60e51b815260206004820152600b60248201527f6e6f20636f2d746f6b656e000000000000000000000000000000000000000000604482015290519081900360640190fd5b600554610b5a9061010090046001600160a01b031633308461110c565b6106a2828261101c565b3390565b6001600160a01b038316610bad5760405162461bcd60e51b81526004018080602001828103825260248152602001806115eb6024913960400191505060405180910390fd5b6001600160a01b038216610bf25760405162461bcd60e51b81526004018080602001828103825260228152602001806114eb6022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b038216610c995760405162461bcd60e51b81526004018080602001828103825260218152602001806115a56021913960400191505060405180910390fd5b610ca5826000836108d1565b610ce2816040518060600160405280602281526020016114c9602291396001600160a01b0385166000908152602081905260409020549190610f2b565b6001600160a01b038316600090815260208190526040902055600254610d08908261119a565b6002556040805182815290516000916001600160a01b038516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35050565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001790526108d19084906111f7565b6001600160a01b038316610e155760405162461bcd60e51b81526004018080602001828103825260258152602001806115c66025913960400191505060405180910390fd5b6001600160a01b038216610e5a5760405162461bcd60e51b81526004018080602001828103825260238152602001806114a66023913960400191505060405180910390fd5b610e658383836108d1565b610ea28160405180606001604052806026815260200161150d602691396001600160a01b0386166000908152602081905260409020549190610f2b565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610ed19082610fc2565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008184841115610fba5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610f7f578181015183820152602001610f67565b50505050905090810190601f168015610fac5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b600082820183811015610727576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b6001600160a01b038216611077576040805162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b611083600083836108d1565b6002546110909082610fc2565b6002556001600160a01b0382166000908152602081905260409020546110b69082610fc2565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b604080516001600160a01b0380861660248301528416604482015260648082018490528251808303909101815260849091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd000000000000000000000000000000000000000000000000000000001790526111949085906111f7565b50505050565b6000828211156111f1576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b600061124c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166112a89092919063ffffffff16565b8051909150156108d15780806020019051602081101561126b57600080fd5b50516108d15760405162461bcd60e51b815260040180806020018281038252602a81526020018061160f602a913960400191505060405180910390fd5b60606112b784846000856112bf565b949350505050565b606030318311156113015760405162461bcd60e51b81526004018080602001828103825260268152602001806115336026913960400191505060405180910390fd5b61130a85611439565b61135b576040805162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015290519081900360640190fd5b600080866001600160a01b031685876040518082805190602001908083835b602083106113b757805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161137a565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d8060008114611419576040519150601f19603f3d011682016040523d82523d6000602084013e61141e565b606091505b509150915061142e82828661143f565b979650505050505050565b3b151590565b6060831561144e575081610727565b82511561145e5782518084602001fd5b60405162461bcd60e51b8152602060048201818152845160248401528451859391928392604401919085019080838360008315610f7f578181015183820152602001610f6756fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a206275726e20616d6f756e7420657863656564732062616c616e636545524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e6365416464726573733a20696e73756666696369656e742062616c616e636520666f722063616c6c45524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a206275726e20616d6f756e74206578636565647320616c6c6f77616e636545524332303a206275726e2066726f6d20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f20616464726573735361666545524332303a204552433230206f7065726174696f6e20646964206e6f74207375636365656445524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa264697066735822122011f14534b53bd4289a85b68c08b54b87febc4dab1fc30df955c69582ee7be94864736f6c63430007060033"

// DeployCrosschainERC20 deploys a new Ethereum contract, binding an instance of CrosschainERC20 to it.
func DeployCrosschainERC20(auth *bind.TransactOpts, backend bind.ContractBackend, _coToken common.Address, _minter common.Address, _name string, _symbol string, _decimals uint8) (common.Address, *types.Transaction, *CrosschainERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(CrosschainERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CrosschainERC20Bin), backend, _coToken, _minter, _name, _symbol, _decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrosschainERC20{CrosschainERC20Caller: CrosschainERC20Caller{contract: contract}, CrosschainERC20Transactor: CrosschainERC20Transactor{contract: contract}, CrosschainERC20Filterer: CrosschainERC20Filterer{contract: contract}}, nil
}

// CrosschainERC20 is an auto generated Go binding around an Ethereum contract.
type CrosschainERC20 struct {
	CrosschainERC20Caller     // Read-only binding to the contract
	CrosschainERC20Transactor // Write-only binding to the contract
	CrosschainERC20Filterer   // Log filterer for contract events
}

// CrosschainERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type CrosschainERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type CrosschainERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrosschainERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrosschainERC20Session struct {
	Contract     *CrosschainERC20  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrosschainERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrosschainERC20CallerSession struct {
	Contract *CrosschainERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CrosschainERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrosschainERC20TransactorSession struct {
	Contract     *CrosschainERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CrosschainERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type CrosschainERC20Raw struct {
	Contract *CrosschainERC20 // Generic contract binding to access the raw methods on
}

// CrosschainERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrosschainERC20CallerRaw struct {
	Contract *CrosschainERC20Caller // Generic read-only contract binding to access the raw methods on
}

// CrosschainERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrosschainERC20TransactorRaw struct {
	Contract *CrosschainERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCrosschainERC20 creates a new instance of CrosschainERC20, bound to a specific deployed contract.
func NewCrosschainERC20(address common.Address, backend bind.ContractBackend) (*CrosschainERC20, error) {
	contract, err := bindCrosschainERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20{CrosschainERC20Caller: CrosschainERC20Caller{contract: contract}, CrosschainERC20Transactor: CrosschainERC20Transactor{contract: contract}, CrosschainERC20Filterer: CrosschainERC20Filterer{contract: contract}}, nil
}

// NewCrosschainERC20Caller creates a new read-only instance of CrosschainERC20, bound to a specific deployed contract.
func NewCrosschainERC20Caller(address common.Address, caller bind.ContractCaller) (*CrosschainERC20Caller, error) {
	contract, err := bindCrosschainERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20Caller{contract: contract}, nil
}

// NewCrosschainERC20Transactor creates a new write-only instance of CrosschainERC20, bound to a specific deployed contract.
func NewCrosschainERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*CrosschainERC20Transactor, error) {
	contract, err := bindCrosschainERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20Transactor{contract: contract}, nil
}

// NewCrosschainERC20Filterer creates a new log filterer instance of CrosschainERC20, bound to a specific deployed contract.
func NewCrosschainERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*CrosschainERC20Filterer, error) {
	contract, err := bindCrosschainERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20Filterer{contract: contract}, nil
}

// bindCrosschainERC20 binds a generic wrapper to an already deployed contract.
func bindCrosschainERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CrosschainERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainERC20 *CrosschainERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainERC20.Contract.CrosschainERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainERC20 *CrosschainERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.CrosschainERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainERC20 *CrosschainERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.CrosschainERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainERC20 *CrosschainERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainERC20 *CrosschainERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainERC20 *CrosschainERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CrosschainERC20.Contract.Allowance(&_CrosschainERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CrosschainERC20.Contract.Allowance(&_CrosschainERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _CrosschainERC20.Contract.BalanceOf(&_CrosschainERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _CrosschainERC20.Contract.BalanceOf(&_CrosschainERC20.CallOpts, account)
}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_CrosschainERC20 *CrosschainERC20Caller) CoToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "coToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_CrosschainERC20 *CrosschainERC20Session) CoToken() (common.Address, error) {
	return _CrosschainERC20.Contract.CoToken(&_CrosschainERC20.CallOpts)
}

// CoToken is a free data retrieval call binding the contract method 0x7f9864f7.
//
// Solidity: function coToken() view returns(address)
func (_CrosschainERC20 *CrosschainERC20CallerSession) CoToken() (common.Address, error) {
	return _CrosschainERC20.Contract.CoToken(&_CrosschainERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_CrosschainERC20 *CrosschainERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_CrosschainERC20 *CrosschainERC20Session) Decimals() (uint8, error) {
	return _CrosschainERC20.Contract.Decimals(&_CrosschainERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_CrosschainERC20 *CrosschainERC20CallerSession) Decimals() (uint8, error) {
	return _CrosschainERC20.Contract.Decimals(&_CrosschainERC20.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CrosschainERC20 *CrosschainERC20Caller) Minter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "minter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CrosschainERC20 *CrosschainERC20Session) Minter() (common.Address, error) {
	return _CrosschainERC20.Contract.Minter(&_CrosschainERC20.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CrosschainERC20 *CrosschainERC20CallerSession) Minter() (common.Address, error) {
	return _CrosschainERC20.Contract.Minter(&_CrosschainERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CrosschainERC20 *CrosschainERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CrosschainERC20 *CrosschainERC20Session) Name() (string, error) {
	return _CrosschainERC20.Contract.Name(&_CrosschainERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CrosschainERC20 *CrosschainERC20CallerSession) Name() (string, error) {
	return _CrosschainERC20.Contract.Name(&_CrosschainERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CrosschainERC20 *CrosschainERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CrosschainERC20 *CrosschainERC20Session) Symbol() (string, error) {
	return _CrosschainERC20.Contract.Symbol(&_CrosschainERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CrosschainERC20 *CrosschainERC20CallerSession) Symbol() (string, error) {
	return _CrosschainERC20.Contract.Symbol(&_CrosschainERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrosschainERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20Session) TotalSupply() (*big.Int, error) {
	return _CrosschainERC20.Contract.TotalSupply(&_CrosschainERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CrosschainERC20 *CrosschainERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _CrosschainERC20.Contract.TotalSupply(&_CrosschainERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Approve(&_CrosschainERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Approve(&_CrosschainERC20.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Burn(&_CrosschainERC20.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Burn(&_CrosschainERC20.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.BurnFrom(&_CrosschainERC20.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.BurnFrom(&_CrosschainERC20.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.DecreaseAllowance(&_CrosschainERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.DecreaseAllowance(&_CrosschainERC20.TransactOpts, spender, subtractedValue)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) Deposit(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Deposit(&_CrosschainERC20.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Deposit(&_CrosschainERC20.TransactOpts, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xffaad6a5.
//
// Solidity: function depositTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) DepositTo(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "depositTo", _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xffaad6a5.
//
// Solidity: function depositTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) DepositTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.DepositTo(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xffaad6a5.
//
// Solidity: function depositTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) DepositTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.DepositTo(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.IncreaseAllowance(&_CrosschainERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.IncreaseAllowance(&_CrosschainERC20.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) Mint(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "mint", _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Mint(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Mint(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Transfer(&_CrosschainERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Transfer(&_CrosschainERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.TransferFrom(&_CrosschainERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_CrosschainERC20 *CrosschainERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.TransferFrom(&_CrosschainERC20.TransactOpts, sender, recipient, amount)
}

// TransferMintership is a paid mutator transaction binding the contract method 0xcf86a95a.
//
// Solidity: function transferMintership(address _newMinter) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) TransferMintership(opts *bind.TransactOpts, _newMinter common.Address) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "transferMintership", _newMinter)
}

// TransferMintership is a paid mutator transaction binding the contract method 0xcf86a95a.
//
// Solidity: function transferMintership(address _newMinter) returns()
func (_CrosschainERC20 *CrosschainERC20Session) TransferMintership(_newMinter common.Address) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.TransferMintership(&_CrosschainERC20.TransactOpts, _newMinter)
}

// TransferMintership is a paid mutator transaction binding the contract method 0xcf86a95a.
//
// Solidity: function transferMintership(address _newMinter) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) TransferMintership(_newMinter common.Address) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.TransferMintership(&_CrosschainERC20.TransactOpts, _newMinter)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Withdraw(&_CrosschainERC20.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.Withdraw(&_CrosschainERC20.TransactOpts, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Transactor) WithdrawTo(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.contract.Transact(opts, "withdrawTo", _to, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20Session) WithdrawTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.WithdrawTo(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address _to, uint256 _amount) returns()
func (_CrosschainERC20 *CrosschainERC20TransactorSession) WithdrawTo(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CrosschainERC20.Contract.WithdrawTo(&_CrosschainERC20.TransactOpts, _to, _amount)
}

// CrosschainERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CrosschainERC20 contract.
type CrosschainERC20ApprovalIterator struct {
	Event *CrosschainERC20Approval // Event containing the contract specifics and raw log

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
func (it *CrosschainERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainERC20Approval)
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
		it.Event = new(CrosschainERC20Approval)
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
func (it *CrosschainERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainERC20Approval represents a Approval event raised by the CrosschainERC20 contract.
type CrosschainERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_CrosschainERC20 *CrosschainERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CrosschainERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CrosschainERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20ApprovalIterator{contract: _CrosschainERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_CrosschainERC20 *CrosschainERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CrosschainERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CrosschainERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainERC20Approval)
				if err := _CrosschainERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_CrosschainERC20 *CrosschainERC20Filterer) ParseApproval(log types.Log) (*CrosschainERC20Approval, error) {
	event := new(CrosschainERC20Approval)
	if err := _CrosschainERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainERC20MinterSetIterator is returned from FilterMinterSet and is used to iterate over the raw logs and unpacked data for MinterSet events raised by the CrosschainERC20 contract.
type CrosschainERC20MinterSetIterator struct {
	Event *CrosschainERC20MinterSet // Event containing the contract specifics and raw log

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
func (it *CrosschainERC20MinterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainERC20MinterSet)
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
		it.Event = new(CrosschainERC20MinterSet)
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
func (it *CrosschainERC20MinterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainERC20MinterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainERC20MinterSet represents a MinterSet event raised by the CrosschainERC20 contract.
type CrosschainERC20MinterSet struct {
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinterSet is a free log retrieval operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed minter)
func (_CrosschainERC20 *CrosschainERC20Filterer) FilterMinterSet(opts *bind.FilterOpts, minter []common.Address) (*CrosschainERC20MinterSetIterator, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CrosschainERC20.contract.FilterLogs(opts, "MinterSet", minterRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20MinterSetIterator{contract: _CrosschainERC20.contract, event: "MinterSet", logs: logs, sub: sub}, nil
}

// WatchMinterSet is a free log subscription operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed minter)
func (_CrosschainERC20 *CrosschainERC20Filterer) WatchMinterSet(opts *bind.WatchOpts, sink chan<- *CrosschainERC20MinterSet, minter []common.Address) (event.Subscription, error) {

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _CrosschainERC20.contract.WatchLogs(opts, "MinterSet", minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainERC20MinterSet)
				if err := _CrosschainERC20.contract.UnpackLog(event, "MinterSet", log); err != nil {
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

// ParseMinterSet is a log parse operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed minter)
func (_CrosschainERC20 *CrosschainERC20Filterer) ParseMinterSet(log types.Log) (*CrosschainERC20MinterSet, error) {
	event := new(CrosschainERC20MinterSet)
	if err := _CrosschainERC20.contract.UnpackLog(event, "MinterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CrosschainERC20 contract.
type CrosschainERC20TransferIterator struct {
	Event *CrosschainERC20Transfer // Event containing the contract specifics and raw log

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
func (it *CrosschainERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainERC20Transfer)
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
		it.Event = new(CrosschainERC20Transfer)
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
func (it *CrosschainERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainERC20Transfer represents a Transfer event raised by the CrosschainERC20 contract.
type CrosschainERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_CrosschainERC20 *CrosschainERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CrosschainERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrosschainERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainERC20TransferIterator{contract: _CrosschainERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_CrosschainERC20 *CrosschainERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CrosschainERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CrosschainERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainERC20Transfer)
				if err := _CrosschainERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_CrosschainERC20 *CrosschainERC20Filterer) ParseTransfer(log types.Log) (*CrosschainERC20Transfer, error) {
	event := new(CrosschainERC20Transfer)
	if err := _CrosschainERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
