// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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

// TimeLedgerTokenMetaData contains all meta data concerning the TimeLedgerToken contract.
var TimeLedgerTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b5060405161188f38038061188f833981810160405281019061003191906102ed565b33828281600390816100439190610573565b5080600490816100539190610573565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100c6575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100bd9190610681565b60405180910390fd5b6100d5816100dd60201b60201c565b50505061069a565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6101ff826101b9565b810181811067ffffffffffffffff8211171561021e5761021d6101c9565b5b80604052505050565b5f6102306101a0565b905061023c82826101f6565b919050565b5f67ffffffffffffffff82111561025b5761025a6101c9565b5b610264826101b9565b9050602081019050919050565b8281835e5f83830152505050565b5f61029161028c84610241565b610227565b9050828152602081018484840111156102ad576102ac6101b5565b5b6102b8848285610271565b509392505050565b5f82601f8301126102d4576102d36101b1565b5b81516102e484826020860161027f565b91505092915050565b5f5f60408385031215610303576103026101a9565b5b5f83015167ffffffffffffffff8111156103205761031f6101ad565b5b61032c858286016102c0565b925050602083015167ffffffffffffffff81111561034d5761034c6101ad565b5b610359858286016102c0565b9150509250929050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806103b157607f821691505b6020821081036103c4576103c361036d565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026104267fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826103eb565b61043086836103eb565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61047461046f61046a84610448565b610451565b610448565b9050919050565b5f819050919050565b61048d8361045a565b6104a16104998261047b565b8484546103f7565b825550505050565b5f5f905090565b6104b86104a9565b6104c3818484610484565b505050565b5b818110156104e6576104db5f826104b0565b6001810190506104c9565b5050565b601f82111561052b576104fc816103ca565b610505846103dc565b81016020851015610514578190505b610528610520856103dc565b8301826104c8565b50505b505050565b5f82821c905092915050565b5f61054b5f1984600802610530565b1980831691505092915050565b5f610563838361053c565b9150826002028217905092915050565b61057c82610363565b67ffffffffffffffff811115610595576105946101c9565b5b61059f825461039a565b6105aa8282856104ea565b5f60209050601f8311600181146105db575f84156105c9578287015190505b6105d38582610558565b86555061063a565b601f1984166105e9866103ca565b5f5b82811015610610578489015182556001820191506020850194506020810190506105eb565b8683101561062d5784890151610629601f89168261053c565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61066b82610642565b9050919050565b61067b81610661565b82525050565b5f6020820190506106945f830184610672565b92915050565b6111e8806106a75f395ff3fe608060405234801561000f575f5ffd5b50600436106100e8575f3560e01c8063715018a61161008a5780639dc29fac116100645780639dc29fac14610238578063a9059cbb14610254578063dd62ed3e14610284578063f2fde38b146102b4576100e8565b8063715018a6146101f25780638da5cb5b146101fc57806395d89b411461021a576100e8565b806323b872dd116100c657806323b872dd14610158578063313ce5671461018857806340c10f19146101a657806370a08231146101c2576100e8565b806306fdde03146100ec578063095ea7b31461010a57806318160ddd1461013a575b5f5ffd5b6100f46102d0565b6040516101019190610e61565b60405180910390f35b610124600480360381019061011f9190610f12565b610360565b6040516101319190610f6a565b60405180910390f35b610142610382565b60405161014f9190610f92565b60405180910390f35b610172600480360381019061016d9190610fab565b61038b565b60405161017f9190610f6a565b60405180910390f35b6101906103b9565b60405161019d9190611016565b60405180910390f35b6101c060048036038101906101bb9190610f12565b6103c1565b005b6101dc60048036038101906101d7919061102f565b6103d7565b6040516101e99190610f92565b60405180910390f35b6101fa61041c565b005b61020461042f565b6040516102119190611069565b60405180910390f35b610222610457565b60405161022f9190610e61565b60405180910390f35b610252600480360381019061024d9190610f12565b6104e7565b005b61026e60048036038101906102699190610f12565b6104fd565b60405161027b9190610f6a565b60405180910390f35b61029e60048036038101906102999190611082565b61051f565b6040516102ab9190610f92565b60405180910390f35b6102ce60048036038101906102c9919061102f565b6105a1565b005b6060600380546102df906110ed565b80601f016020809104026020016040519081016040528092919081815260200182805461030b906110ed565b80156103565780601f1061032d57610100808354040283529160200191610356565b820191905f5260205f20905b81548152906001019060200180831161033957829003601f168201915b5050505050905090565b5f5f61036a610625565b905061037781858561062c565b600191505092915050565b5f600254905090565b5f5f610395610625565b90506103a285828561063e565b6103ad8585856106d1565b60019150509392505050565b5f6012905090565b6103c96107c1565b6103d38282610848565b5050565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6104246107c1565b61042d5f6108c7565b565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b606060048054610466906110ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610492906110ed565b80156104dd5780601f106104b4576101008083540402835291602001916104dd565b820191905f5260205f20905b8154815290600101906020018083116104c057829003601f168201915b5050505050905090565b6104ef6107c1565b6104f9828261098a565b5050565b5f5f610507610625565b90506105148185856106d1565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b6105a96107c1565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610619575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016106109190611069565b60405180910390fd5b610622816108c7565b50565b5f33905090565b6106398383836001610a09565b505050565b5f610649848461051f565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8110156106cb57818110156106bc578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016106b39392919061111d565b60405180910390fd5b6106ca84848484035f610a09565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610741575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016107389190611069565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036107b1575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016107a89190611069565b60405180910390fd5b6107bc838383610bd8565b505050565b6107c9610625565b73ffffffffffffffffffffffffffffffffffffffff166107e761042f565b73ffffffffffffffffffffffffffffffffffffffff16146108465761080a610625565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161083d9190611069565b60405180910390fd5b565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108b8575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016108af9190611069565b60405180910390fd5b6108c35f8383610bd8565b5050565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036109fa575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016109f19190611069565b60405180910390fd5b610a05825f83610bd8565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610a79575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610a709190611069565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610ae9575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610ae09190611069565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610bd2578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610bc99190610f92565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610c28578060025f828254610c1c919061117f565b92505081905550610cf6565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610cb1578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610ca89392919061111d565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610d3d578060025f8282540392505081905550610d87565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610de49190610f92565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610e3382610df1565b610e3d8185610dfb565b9350610e4d818560208601610e0b565b610e5681610e19565b840191505092915050565b5f6020820190508181035f830152610e798184610e29565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610eae82610e85565b9050919050565b610ebe81610ea4565b8114610ec8575f5ffd5b50565b5f81359050610ed981610eb5565b92915050565b5f819050919050565b610ef181610edf565b8114610efb575f5ffd5b50565b5f81359050610f0c81610ee8565b92915050565b5f5f60408385031215610f2857610f27610e81565b5b5f610f3585828601610ecb565b9250506020610f4685828601610efe565b9150509250929050565b5f8115159050919050565b610f6481610f50565b82525050565b5f602082019050610f7d5f830184610f5b565b92915050565b610f8c81610edf565b82525050565b5f602082019050610fa55f830184610f83565b92915050565b5f5f5f60608486031215610fc257610fc1610e81565b5b5f610fcf86828701610ecb565b9350506020610fe086828701610ecb565b9250506040610ff186828701610efe565b9150509250925092565b5f60ff82169050919050565b61101081610ffb565b82525050565b5f6020820190506110295f830184611007565b92915050565b5f6020828403121561104457611043610e81565b5b5f61105184828501610ecb565b91505092915050565b61106381610ea4565b82525050565b5f60208201905061107c5f83018461105a565b92915050565b5f5f6040838503121561109857611097610e81565b5b5f6110a585828601610ecb565b92505060206110b685828601610ecb565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061110457607f821691505b602082108103611117576111166110c0565b5b50919050565b5f6060820190506111305f83018661105a565b61113d6020830185610f83565b61114a6040830184610f83565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61118982610edf565b915061119483610edf565b92508282019050808211156111ac576111ab611152565b5b9291505056fea26469706673582212204230db4778eabb3ce42550935529740cb7deb46798d1024881c5f286c3e810c464736f6c634300081e0033",
}

// TimeLedgerTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use TimeLedgerTokenMetaData.ABI instead.
var TimeLedgerTokenABI = TimeLedgerTokenMetaData.ABI

// TimeLedgerTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TimeLedgerTokenMetaData.Bin instead.
var TimeLedgerTokenBin = TimeLedgerTokenMetaData.Bin

// DeployTimeLedgerToken deploys a new Ethereum contract, binding an instance of TimeLedgerToken to it.
func DeployTimeLedgerToken(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string) (common.Address, *types.Transaction, *TimeLedgerToken, error) {
	parsed, err := TimeLedgerTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TimeLedgerTokenBin), backend, name_, symbol_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TimeLedgerToken{TimeLedgerTokenCaller: TimeLedgerTokenCaller{contract: contract}, TimeLedgerTokenTransactor: TimeLedgerTokenTransactor{contract: contract}, TimeLedgerTokenFilterer: TimeLedgerTokenFilterer{contract: contract}}, nil
}

// TimeLedgerToken is an auto generated Go binding around an Ethereum contract.
type TimeLedgerToken struct {
	TimeLedgerTokenCaller     // Read-only binding to the contract
	TimeLedgerTokenTransactor // Write-only binding to the contract
	TimeLedgerTokenFilterer   // Log filterer for contract events
}

// TimeLedgerTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TimeLedgerTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TimeLedgerTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TimeLedgerTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TimeLedgerTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TimeLedgerTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TimeLedgerTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TimeLedgerTokenSession struct {
	Contract     *TimeLedgerToken  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TimeLedgerTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TimeLedgerTokenCallerSession struct {
	Contract *TimeLedgerTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// TimeLedgerTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TimeLedgerTokenTransactorSession struct {
	Contract     *TimeLedgerTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// TimeLedgerTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TimeLedgerTokenRaw struct {
	Contract *TimeLedgerToken // Generic contract binding to access the raw methods on
}

// TimeLedgerTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TimeLedgerTokenCallerRaw struct {
	Contract *TimeLedgerTokenCaller // Generic read-only contract binding to access the raw methods on
}

// TimeLedgerTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TimeLedgerTokenTransactorRaw struct {
	Contract *TimeLedgerTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTimeLedgerToken creates a new instance of TimeLedgerToken, bound to a specific deployed contract.
func NewTimeLedgerToken(address common.Address, backend bind.ContractBackend) (*TimeLedgerToken, error) {
	contract, err := bindTimeLedgerToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerToken{TimeLedgerTokenCaller: TimeLedgerTokenCaller{contract: contract}, TimeLedgerTokenTransactor: TimeLedgerTokenTransactor{contract: contract}, TimeLedgerTokenFilterer: TimeLedgerTokenFilterer{contract: contract}}, nil
}

// NewTimeLedgerTokenCaller creates a new read-only instance of TimeLedgerToken, bound to a specific deployed contract.
func NewTimeLedgerTokenCaller(address common.Address, caller bind.ContractCaller) (*TimeLedgerTokenCaller, error) {
	contract, err := bindTimeLedgerToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenCaller{contract: contract}, nil
}

// NewTimeLedgerTokenTransactor creates a new write-only instance of TimeLedgerToken, bound to a specific deployed contract.
func NewTimeLedgerTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TimeLedgerTokenTransactor, error) {
	contract, err := bindTimeLedgerToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenTransactor{contract: contract}, nil
}

// NewTimeLedgerTokenFilterer creates a new log filterer instance of TimeLedgerToken, bound to a specific deployed contract.
func NewTimeLedgerTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TimeLedgerTokenFilterer, error) {
	contract, err := bindTimeLedgerToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenFilterer{contract: contract}, nil
}

// bindTimeLedgerToken binds a generic wrapper to an already deployed contract.
func bindTimeLedgerToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TimeLedgerTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TimeLedgerToken *TimeLedgerTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TimeLedgerToken.Contract.TimeLedgerTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TimeLedgerToken *TimeLedgerTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TimeLedgerTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TimeLedgerToken *TimeLedgerTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TimeLedgerTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TimeLedgerToken *TimeLedgerTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TimeLedgerToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TimeLedgerToken *TimeLedgerTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TimeLedgerToken *TimeLedgerTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TimeLedgerToken.Contract.Allowance(&_TimeLedgerToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TimeLedgerToken.Contract.Allowance(&_TimeLedgerToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TimeLedgerToken.Contract.BalanceOf(&_TimeLedgerToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TimeLedgerToken.Contract.BalanceOf(&_TimeLedgerToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TimeLedgerToken *TimeLedgerTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TimeLedgerToken *TimeLedgerTokenSession) Decimals() (uint8, error) {
	return _TimeLedgerToken.Contract.Decimals(&_TimeLedgerToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) Decimals() (uint8, error) {
	return _TimeLedgerToken.Contract.Decimals(&_TimeLedgerToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenSession) Name() (string, error) {
	return _TimeLedgerToken.Contract.Name(&_TimeLedgerToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) Name() (string, error) {
	return _TimeLedgerToken.Contract.Name(&_TimeLedgerToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TimeLedgerToken *TimeLedgerTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TimeLedgerToken *TimeLedgerTokenSession) Owner() (common.Address, error) {
	return _TimeLedgerToken.Contract.Owner(&_TimeLedgerToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) Owner() (common.Address, error) {
	return _TimeLedgerToken.Contract.Owner(&_TimeLedgerToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenSession) Symbol() (string, error) {
	return _TimeLedgerToken.Contract.Symbol(&_TimeLedgerToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) Symbol() (string, error) {
	return _TimeLedgerToken.Contract.Symbol(&_TimeLedgerToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TimeLedgerToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenSession) TotalSupply() (*big.Int, error) {
	return _TimeLedgerToken.Contract.TotalSupply(&_TimeLedgerToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TimeLedgerToken *TimeLedgerTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _TimeLedgerToken.Contract.TotalSupply(&_TimeLedgerToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Approve(&_TimeLedgerToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Approve(&_TimeLedgerToken.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "burn", from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Burn(&_TimeLedgerToken.TransactOpts, from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Burn(&_TimeLedgerToken.TransactOpts, from, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Mint(&_TimeLedgerToken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Mint(&_TimeLedgerToken.TransactOpts, to, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TimeLedgerToken *TimeLedgerTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.RenounceOwnership(&_TimeLedgerToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.RenounceOwnership(&_TimeLedgerToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Transfer(&_TimeLedgerToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.Transfer(&_TimeLedgerToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TransferFrom(&_TimeLedgerToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TransferFrom(&_TimeLedgerToken.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TimeLedgerToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TimeLedgerToken *TimeLedgerTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TransferOwnership(&_TimeLedgerToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TimeLedgerToken *TimeLedgerTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TimeLedgerToken.Contract.TransferOwnership(&_TimeLedgerToken.TransactOpts, newOwner)
}

// TimeLedgerTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TimeLedgerToken contract.
type TimeLedgerTokenApprovalIterator struct {
	Event *TimeLedgerTokenApproval // Event containing the contract specifics and raw log

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
func (it *TimeLedgerTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TimeLedgerTokenApproval)
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
		it.Event = new(TimeLedgerTokenApproval)
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
func (it *TimeLedgerTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TimeLedgerTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TimeLedgerTokenApproval represents a Approval event raised by the TimeLedgerToken contract.
type TimeLedgerTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TimeLedgerTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenApprovalIterator{contract: _TimeLedgerToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TimeLedgerTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TimeLedgerTokenApproval)
				if err := _TimeLedgerToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_TimeLedgerToken *TimeLedgerTokenFilterer) ParseApproval(log types.Log) (*TimeLedgerTokenApproval, error) {
	event := new(TimeLedgerTokenApproval)
	if err := _TimeLedgerToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TimeLedgerTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TimeLedgerToken contract.
type TimeLedgerTokenOwnershipTransferredIterator struct {
	Event *TimeLedgerTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TimeLedgerTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TimeLedgerTokenOwnershipTransferred)
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
		it.Event = new(TimeLedgerTokenOwnershipTransferred)
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
func (it *TimeLedgerTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TimeLedgerTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TimeLedgerTokenOwnershipTransferred represents a OwnershipTransferred event raised by the TimeLedgerToken contract.
type TimeLedgerTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TimeLedgerTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenOwnershipTransferredIterator{contract: _TimeLedgerToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TimeLedgerTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TimeLedgerTokenOwnershipTransferred)
				if err := _TimeLedgerToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TimeLedgerToken *TimeLedgerTokenFilterer) ParseOwnershipTransferred(log types.Log) (*TimeLedgerTokenOwnershipTransferred, error) {
	event := new(TimeLedgerTokenOwnershipTransferred)
	if err := _TimeLedgerToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TimeLedgerTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TimeLedgerToken contract.
type TimeLedgerTokenTransferIterator struct {
	Event *TimeLedgerTokenTransfer // Event containing the contract specifics and raw log

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
func (it *TimeLedgerTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TimeLedgerTokenTransfer)
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
		it.Event = new(TimeLedgerTokenTransfer)
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
func (it *TimeLedgerTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TimeLedgerTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TimeLedgerTokenTransfer represents a Transfer event raised by the TimeLedgerToken contract.
type TimeLedgerTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TimeLedgerTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TimeLedgerTokenTransferIterator{contract: _TimeLedgerToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TimeLedgerToken *TimeLedgerTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TimeLedgerTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TimeLedgerToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TimeLedgerTokenTransfer)
				if err := _TimeLedgerToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_TimeLedgerToken *TimeLedgerTokenFilterer) ParseTransfer(log types.Log) (*TimeLedgerTokenTransfer, error) {
	event := new(TimeLedgerTokenTransfer)
	if err := _TimeLedgerToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
