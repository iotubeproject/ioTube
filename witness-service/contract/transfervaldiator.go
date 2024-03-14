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
const TransferValidatorABI = "[{\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_witnessList\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"name\":\"Settled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"minters\",\"outputs\":[{\"internalType\":\"contractIMinter\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"settles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenLists\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"witnessList\",\"outputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"generateKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"keys\",\"type\":\"bytes32[]\"}],\"name\":\"concatKeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"cashiers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indexes\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"submitMulti\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"getTokenGroup\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"extractWitnesses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"witnesses\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"cashier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numOfPairs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIAllowlist\",\"name\":\"_tokenList\",\"type\":\"address\"},{\"internalType\":\"contractIMinter\",\"name\":\"_minter\",\"type\":\"address\"}],\"name\":\"addPair\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TransferValidatorBin is the compiled bytecode used for deploying new contracts.
var TransferValidatorBin = "0x60806040526000805460ff60a01b1916905534801561001d57600080fd5b50604051611b2e380380611b2e8339818101604052602081101561004057600080fd5b5051600080546001600160a01b03199081163317909155600480546001600160a01b0390931692909116919091179055611aaf8061007f6000396000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063ba390a6411610071578063ba390a64146106f1578063c836fef0146107ec578063e01eba711461088d578063f2fde38b146108b3578063f98b2332146108d957610116565b80638da5cb5b146102315780639f8c11e314610239578063a9013dce146105e9578063b6f3e087146106c357610116565b80635c975abb116100e95780635c975abb1461018c5780636b6bc862146101a85780638356b148146102045780638456cb591461020c5780638623ec7b1461021457610116565b80630900f0101461011b5780631cb928a914610143578063373f0d491461017c5780633f4ba83a14610184575b600080fd5b6101416004803603602081101561013157600080fd5b50356001600160a01b03166108f6565b005b6101606004803603602081101561015957600080fd5b5035610a47565b604080516001600160a01b039092168252519081900360200190f35b610160610a6e565b610141610a7d565b610194610ae0565b604080519115158252519081900360200190f35b6101f2600480360360c08110156101be57600080fd5b506001600160a01b038135811691602081013582169160408201359160608101358216916080820135169060a00135610af0565b60408051918252519081900360200190f35b6101f2610b64565b610141610b6a565b6101606004803603602081101561022a57600080fd5b5035610bd4565b610160610be1565b610141600480360360e081101561024f57600080fd5b810190602081018135600160201b81111561026957600080fd5b82018360208201111561027b57600080fd5b803590602001918460208302840111600160201b8311171561029c57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156102eb57600080fd5b8201836020820111156102fd57600080fd5b803590602001918460208302840111600160201b8311171561031e57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561036d57600080fd5b82018360208201111561037f57600080fd5b803590602001918460208302840111600160201b831117156103a057600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156103ef57600080fd5b82018360208201111561040157600080fd5b803590602001918460208302840111600160201b8311171561042257600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561047157600080fd5b82018360208201111561048357600080fd5b803590602001918460208302840111600160201b831117156104a457600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156104f357600080fd5b82018360208201111561050557600080fd5b803590602001918460208302840111600160201b8311171561052657600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561057557600080fd5b82018360208201111561058757600080fd5b803590602001918460018302840111600160201b831117156105a857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610bf0945050505050565b610141600480360360e08110156105ff57600080fd5b6001600160a01b0382358116926020810135821692604082013592606083013581169260808101359091169160a0820135919081019060e0810160c0820135600160201b81111561064f57600080fd5b82018360208201111561066157600080fd5b803590602001918460018302840111600160201b8311171561068257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611175945050505050565b610141600480360360408110156106d957600080fd5b506001600160a01b0381358116916020013516611562565b61079c6004803603604081101561070757600080fd5b81359190810190604081016020820135600160201b81111561072857600080fd5b82018360208201111561073a57600080fd5b803590602001918460018302840111600160201b8311171561075b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611601945050505050565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156107d85781810151838201526020016107c0565b505050509050019250505060405180910390f35b6101f26004803603602081101561080257600080fd5b810190602081018135600160201b81111561081c57600080fd5b82018360208201111561082e57600080fd5b803590602001918460208302840111600160201b8311171561084f57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506117d7945050505050565b6101f2600480360360208110156108a357600080fd5b50356001600160a01b0316611832565b610141600480360360208110156108c957600080fd5b50356001600160a01b0316611926565b6101f2600480360360208110156108ef57600080fd5b50356119ab565b6000546001600160a01b0316331461090d57600080fd5b3060005b600254811015610a425760006002828154811061092a57fe5b9060005260206000200160009054906101000a90046001600160a01b03169050826001600160a01b0316816001600160a01b0316638da5cb5b6040518163ffffffff1660e01b815260040160206040518083038186803b15801561098d57600080fd5b505afa1580156109a1573d6000803e3d6000fd5b505050506040513d60208110156109b757600080fd5b50516001600160a01b03161415610a3957806001600160a01b031663f2fde38b856040518263ffffffff1660e01b815260040180826001600160a01b03166001600160a01b03168152602001915050600060405180830381600087803b158015610a2057600080fd5b505af1158015610a34573d6000803e3d6000fd5b505050505b50600101610911565b505050565b60038181548110610a5457fe5b6000918252602090912001546001600160a01b0316905081565b6004546001600160a01b031681565b6000546001600160a01b03163314610a9457600080fd5b600054600160a01b900460ff16610aaa57600080fd5b6000805460ff60a01b191681556040517f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b339190a1565b600054600160a01b900460ff1681565b6040805130606090811b6020808401919091526bffffffffffffffffffffffff1999821b8a16603484015297811b89166048830152605c82019690965293851b8716607c8501529190931b909416609082015260a4808201929092528351808203909201825260c401909252815191012090565b60035490565b6000546001600160a01b03163314610b8157600080fd5b600054600160a01b900460ff1615610b9857600080fd5b6000805460ff60a01b1916600160a01b1781556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff6259190a1565b60028181548110610a5457fe5b6000546001600160a01b031681565b600054600160a01b900460ff1615610c0757600080fd5b85518751148015610c19575084518651145b8015610c26575083518551145b8015610c33575082518451145b8015610c40575081518351145b610c86576040805162461bcd60e51b8152602060048201526012602482015271696e76616c696420706172616d657465727360701b604482015290519081900360640190fd5b60608751604051908082528060200260200182016040528015610cb3578160200160208202803883390190505b50905060005b8851811015610e5857610d42898281518110610cd157fe5b6020026020010151898381518110610ce557fe5b6020026020010151898481518110610cf957fe5b6020026020010151898581518110610d0d57fe5b6020026020010151898681518110610d2157fe5b6020026020010151898781518110610d3557fe5b6020026020010151610af0565b828281518110610d4e57fe5b60200260200101818152505060016000838381518110610d6a57fe5b6020026020010151815260200190815260200160002054600014610dd1576040805162461bcd60e51b81526020600482015260196024820152781d1c985b9cd9995c881a185cc81899595b881cd95d1d1b1959603a1b604482015290519081900360640190fd5b60005b81811015610e4f57828181518110610de857fe5b6020026020010151838381518110610dfc57fe5b60200260200101511415610e47576040805162461bcd60e51b815260206004820152600d60248201526c6475706c6963617465206b657960981b604482015290519081900360640190fd5b600101610dd4565b50600101610cb9565b506060610e6d610e67836117d7565b84611601565b905060008151118015610ef75750600480546040805163593f696960e01b815290516001600160a01b039092169263593f6969928282019260209290829003018186803b158015610ebd57600080fd5b505afa158015610ed1573d6000803e3d6000fd5b505050506040513d6020811015610ee757600080fd5b5051815160029091026003909102115b610f41576040805162461bcd60e51b8152602060048201526016602482015275696e73756666696369656e74207769746e657373657360501b604482015290519081900360640190fd5b60005b8951811015611169574360016000858481518110610f5e57fe5b60200260200101518152602001908152602001600020819055506002610f968a8381518110610f8957fe5b6020026020010151611832565b81548110610fa057fe5b60009182526020909120015489516001600160a01b039091169063c6c3bbe6908b9084908110610fcc57fe5b6020026020010151888481518110610fe057fe5b6020026020010151888581518110610ff457fe5b60200260200101516040518463ffffffff1660e01b815260040180846001600160a01b03166001600160a01b03168152602001836001600160a01b03166001600160a01b031681526020018281526020019350505050602060405180830381600087803b15801561106457600080fd5b505af1158015611078573d6000803e3d6000fd5b505050506040513d602081101561108e57600080fd5b50516110d8576040805162461bcd60e51b81526020600482015260146024820152733330b4b632b2103a379036b4b73a103a37b5b2b760611b604482015290519081900360640190fd5b8281815181106110e457fe5b60200260200101517fe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0836040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561114e578181015183820152602001611136565b505050509050019250505060405180910390a2600101610f44565b50505050505050505050565b600054600160a01b900460ff161561118c57600080fd5b816111d6576040805162461bcd60e51b8152602060048201526015602482015274616d6f756e742063616e6e6f74206265207a65726f60581b604482015290519081900360640190fd5b6001600160a01b038316611231576040805162461bcd60e51b815260206004820152601860248201527f726563697069656e742063616e6e6f74206265207a65726f0000000000000000604482015290519081900360640190fd5b604181518161123c57fe5b061561128f576040805162461bcd60e51b815260206004820152601860248201527f696e76616c6964207369676e6174757265206c656e6774680000000000000000604482015290519081900360640190fd5b600061129f888888888888610af0565b600081815260016020526040902054909150156112ff576040805162461bcd60e51b81526020600482015260196024820152781d1c985b9cd9995c881a185cc81899595b881cd95d1d1b1959603a1b604482015290519081900360640190fd5b606061130b8284611601565b9050600081511180156113955750600480546040805163593f696960e01b815290516001600160a01b039092169263593f6969928282019260209290829003018186803b15801561135b57600080fd5b505afa15801561136f573d6000803e3d6000fd5b505050506040513d602081101561138557600080fd5b5051815160029091026003909102115b6113df576040805162461bcd60e51b8152602060048201526016602482015275696e73756666696369656e74207769746e657373657360501b604482015290519081900360640190fd5b600082815260016020526040902043905560026113fb89611832565b8154811061140557fe5b600091825260208083209091015460408051636361ddf360e11b81526001600160a01b038d811660048301528a81166024830152604482018a90529151919092169363c6c3bbe693606480850194919392918390030190829087803b15801561146d57600080fd5b505af1158015611481573d6000803e3d6000fd5b505050506040513d602081101561149757600080fd5b50516114e1576040805162461bcd60e51b81526020600482015260146024820152733330b4b632b2103a379036b4b73a103a37b5b2b760611b604482015290519081900360640190fd5b817fe24922ac8cf2a1430fb8c7ce6a9125d7db5076a1eb2cefced90e82d6fcb24db0826040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561154457818101518382015260200161152c565b505050509050019250505060405180910390a2505050505050505050565b6000546001600160a01b0316331461157957600080fd5b6003805460018181019092557fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b0180546001600160a01b039485166001600160a01b0319918216179091556002805492830181556000527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9091018054929093169116179055565b60606000604183518161161057fe5b0490508060405190808252806020026020018201604052801561163d578160200160208202803883390190505b50915060005b818110156117cf57600061165b8686846041026119bd565b600480546040805163babcc53960e01b81526001600160a01b0380861694820194909452905193945091169163babcc53991602480820192602092909190829003018186803b1580156116ad57600080fd5b505afa1580156116c1573d6000803e3d6000fd5b505050506040513d60208110156116d757600080fd5b505161171f576040805162461bcd60e51b8152602060048201526012602482015271696e76616c6964207369676e61747572657360701b604482015290519081900360640190fd5b60005b828110156117a05784818151811061173657fe5b60200260200101516001600160a01b0316826001600160a01b03161415611798576040805162461bcd60e51b81526020600482015260116024820152706475706c6963617465207769746e65737360781b604482015290519081900360640190fd5b600101611722565b50808483815181106117ae57fe5b6001600160a01b039092166020928302919091019091015250600101611643565b505092915050565b60008160405160200180828051906020019060200280838360005b8381101561180a5781810151838201526020016117f2565b505050509050019150506040516020818303038152906040528051906020012090505b919050565b6000805b6003548110156118e0576003818154811061184d57fe5b600091825260209182902001546040805163babcc53960e01b81526001600160a01b0387811660048301529151919092169263babcc5399260248082019391829003018186803b1580156118a057600080fd5b505afa1580156118b4573d6000803e3d6000fd5b505050506040513d60208110156118ca57600080fd5b5051156118d857905061182d565b600101611836565b506040805162461bcd60e51b8152602060048201526015602482015274696e76616c696420746f6b656e206164647265737360581b604482015290519081900360640190fd5b6000546001600160a01b0316331461193d57600080fd5b6001600160a01b03811661195057600080fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b60016020526000908152604090205481565b8181016020810151604082015160609092015160009290831a601b8110156119e357601b015b8060ff16601b141580156119fb57508060ff16601c14155b15611a0c5760009350505050611a73565b604080516000815260208082018084528a905260ff8416828401526060820186905260808201859052915160019260a0808401939192601f1981019281900390910190855afa158015611a63573d6000803e3d6000fd5b5050506020604051035193505050505b939250505056fea265627a7a723158205a32069e2932e5e1355e1e353eafc28f6ff72a6f3b90478470af548bd72f6c5364736f6c63430005110032"

// DeployTransferValidator deploys a new Ethereum contract, binding an instance of TransferValidator to it.
func DeployTransferValidator(auth *bind.TransactOpts, backend bind.ContractBackend, _witnessList common.Address) (common.Address, *types.Transaction, *TransferValidator, error) {
	parsed, err := abi.JSON(strings.NewReader(TransferValidatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TransferValidatorBin), backend, _witnessList)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransferValidator{TransferValidatorCaller: TransferValidatorCaller{contract: contract}, TransferValidatorTransactor: TransferValidatorTransactor{contract: contract}, TransferValidatorFilterer: TransferValidatorFilterer{contract: contract}}, nil
}

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

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidator *TransferValidatorCaller) GenerateKey(opts *bind.CallOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TransferValidator.contract.Call(opts, &out, "generateKey", cashier, tokenAddr, index, from, to, amount)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidator *TransferValidatorSession) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, cashier, tokenAddr, index, from, to, amount)
}

// GenerateKey is a free data retrieval call binding the contract method 0x6b6bc862.
//
// Solidity: function generateKey(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount) view returns(bytes32)
func (_TransferValidator *TransferValidatorCallerSession) GenerateKey(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int) ([32]byte, error) {
	return _TransferValidator.Contract.GenerateKey(&_TransferValidator.CallOpts, cashier, tokenAddr, index, from, to, amount)
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

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidator *TransferValidatorTransactor) Submit(opts *bind.TransactOpts, cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "submit", cashier, tokenAddr, index, from, to, amount, signatures)
}

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidator *TransferValidatorSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures)
}

// Submit is a paid mutator transaction binding the contract method 0xa9013dce.
//
// Solidity: function submit(address cashier, address tokenAddr, uint256 index, address from, address to, uint256 amount, bytes signatures) returns()
func (_TransferValidator *TransferValidatorTransactorSession) Submit(cashier common.Address, tokenAddr common.Address, index *big.Int, from common.Address, to common.Address, amount *big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.Submit(&_TransferValidator.TransactOpts, cashier, tokenAddr, index, from, to, amount, signatures)
}

// SubmitMulti is a paid mutator transaction binding the contract method 0x9f8c11e3.
//
// Solidity: function submitMulti(address[] cashiers, address[] tokenAddrs, uint256[] indexes, address[] senders, address[] recipients, uint256[] amounts, bytes signatures) returns()
func (_TransferValidator *TransferValidatorTransactor) SubmitMulti(opts *bind.TransactOpts, cashiers []common.Address, tokenAddrs []common.Address, indexes []*big.Int, senders []common.Address, recipients []common.Address, amounts []*big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.contract.Transact(opts, "submitMulti", cashiers, tokenAddrs, indexes, senders, recipients, amounts, signatures)
}

// SubmitMulti is a paid mutator transaction binding the contract method 0x9f8c11e3.
//
// Solidity: function submitMulti(address[] cashiers, address[] tokenAddrs, uint256[] indexes, address[] senders, address[] recipients, uint256[] amounts, bytes signatures) returns()
func (_TransferValidator *TransferValidatorSession) SubmitMulti(cashiers []common.Address, tokenAddrs []common.Address, indexes []*big.Int, senders []common.Address, recipients []common.Address, amounts []*big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.SubmitMulti(&_TransferValidator.TransactOpts, cashiers, tokenAddrs, indexes, senders, recipients, amounts, signatures)
}

// SubmitMulti is a paid mutator transaction binding the contract method 0x9f8c11e3.
//
// Solidity: function submitMulti(address[] cashiers, address[] tokenAddrs, uint256[] indexes, address[] senders, address[] recipients, uint256[] amounts, bytes signatures) returns()
func (_TransferValidator *TransferValidatorTransactorSession) SubmitMulti(cashiers []common.Address, tokenAddrs []common.Address, indexes []*big.Int, senders []common.Address, recipients []common.Address, amounts []*big.Int, signatures []byte) (*types.Transaction, error) {
	return _TransferValidator.Contract.SubmitMulti(&_TransferValidator.TransactOpts, cashiers, tokenAddrs, indexes, senders, recipients, amounts, signatures)
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
