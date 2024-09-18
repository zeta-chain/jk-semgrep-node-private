// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
)

// DepositorMetaData contains all meta data concerning the Depositor contract.
var DepositorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"custody_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"runDeposits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161073c38038061073c833981810160405281019061003291906100cf565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050506100fc565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061009c82610071565b9050919050565b6100ac81610091565b81146100b757600080fd5b50565b6000815190506100c9816100a3565b92915050565b6000602082840312156100e5576100e461006c565b5b60006100f3848285016100ba565b91505092915050565b60805161061f61011d6000396000818160e4015261016b015261061f6000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80633d496c9314610030575b600080fd5b61004a6004803603810190610045919061032a565b61004c565b005b8473ffffffffffffffffffffffffffffffffffffffff166323b872dd333084886100769190610415565b6040518463ffffffff1660e01b815260040161009493929190610475565b600060405180830381600087803b1580156100ae57600080fd5b505af11580156100c2573d6000803e3d6000fd5b505050508473ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000838761010f9190610415565b6040518363ffffffff1660e01b815260040161012c9291906104ac565b600060405180830381600087803b15801561014657600080fd5b505af115801561015a573d6000803e3d6000fd5b5050505060005b8181101561020b577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e609055e8989898989896040518763ffffffff1660e01b81526004016101cc96959493929190610592565b600060405180830381600087803b1580156101e657600080fd5b505af11580156101fa573d6000803e3d6000fd5b505050508080600101915050610161565b5050505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126102445761024361021f565b5b8235905067ffffffffffffffff81111561026157610260610224565b5b60208301915083600182028301111561027d5761027c610229565b5b9250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006102af82610284565b9050919050565b60006102c1826102a4565b9050919050565b6102d1816102b6565b81146102dc57600080fd5b50565b6000813590506102ee816102c8565b92915050565b6000819050919050565b610307816102f4565b811461031257600080fd5b50565b600081359050610324816102fe565b92915050565b600080600080600080600060a0888a03121561034957610348610215565b5b600088013567ffffffffffffffff8111156103675761036661021a565b5b6103738a828b0161022e565b975097505060206103868a828b016102df565b95505060406103978a828b01610315565b945050606088013567ffffffffffffffff8111156103b8576103b761021a565b5b6103c48a828b0161022e565b935093505060806103d78a828b01610315565b91505092959891949750929550565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610420826102f4565b915061042b836102f4565b9250828202610439816102f4565b915082820484148315176104505761044f6103e6565b5b5092915050565b610460816102a4565b82525050565b61046f816102f4565b82525050565b600060608201905061048a6000830186610457565b6104976020830185610457565b6104a46040830184610466565b949350505050565b60006040820190506104c16000830185610457565b6104ce6020830184610466565b9392505050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b600061051283856104d5565b935061051f8385846104e6565b610528836104f5565b840190509392505050565b6000819050919050565b600061055861055361054e84610284565b610533565b610284565b9050919050565b600061056a8261053d565b9050919050565b600061057c8261055f565b9050919050565b61058c81610571565b82525050565b600060808201905081810360008301526105ad81888a610506565b90506105bc6020830187610583565b6105c96040830186610466565b81810360608301526105dc818486610506565b905097965050505050505056fea2646970667358221220a5a92c8f323a5a31ee98d3f743848cbaba73616e237bfd5f78b15847c3c0abb864736f6c63430008190033",
}

// DepositorABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositorMetaData.ABI instead.
var DepositorABI = DepositorMetaData.ABI

// DepositorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositorMetaData.Bin instead.
var DepositorBin = DepositorMetaData.Bin

// DeployDepositor deploys a new Ethereum contract, binding an instance of Depositor to it.
func DeployDepositor(auth *bind.TransactOpts, backend bind.ContractBackend, custody_ common.Address) (common.Address, *types.Transaction, *Depositor, error) {
	parsed, err := DepositorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositorBin), backend, custody_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// Depositor is an auto generated Go binding around an Ethereum contract.
type Depositor struct {
	DepositorCaller     // Read-only binding to the contract
	DepositorTransactor // Write-only binding to the contract
	DepositorFilterer   // Log filterer for contract events
}

// DepositorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositorSession struct {
	Contract     *Depositor        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositorCallerSession struct {
	Contract *DepositorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DepositorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositorTransactorSession struct {
	Contract     *DepositorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DepositorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositorRaw struct {
	Contract *Depositor // Generic contract binding to access the raw methods on
}

// DepositorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositorCallerRaw struct {
	Contract *DepositorCaller // Generic read-only contract binding to access the raw methods on
}

// DepositorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositorTransactorRaw struct {
	Contract *DepositorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositor creates a new instance of Depositor, bound to a specific deployed contract.
func NewDepositor(address common.Address, backend bind.ContractBackend) (*Depositor, error) {
	contract, err := bindDepositor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// NewDepositorCaller creates a new read-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorCaller(address common.Address, caller bind.ContractCaller) (*DepositorCaller, error) {
	contract, err := bindDepositor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorCaller{contract: contract}, nil
}

// NewDepositorTransactor creates a new write-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositorTransactor, error) {
	contract, err := bindDepositor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorTransactor{contract: contract}, nil
}

// NewDepositorFilterer creates a new log filterer instance of Depositor, bound to a specific deployed contract.
func NewDepositorFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositorFilterer, error) {
	contract, err := bindDepositor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositorFilterer{contract: contract}, nil
}

// bindDepositor binds a generic wrapper to an already deployed contract.
func bindDepositor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.DepositorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transact(opts, method, params...)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorTransactor) RunDeposits(opts *bind.TransactOpts, recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.contract.Transact(opts, "runDeposits", recipient, asset, amount, message, count)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorSession) RunDeposits(recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.Contract.RunDeposits(&_Depositor.TransactOpts, recipient, asset, amount, message, count)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorTransactorSession) RunDeposits(recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.Contract.RunDeposits(&_Depositor.TransactOpts, recipient, asset, amount, message, count)
}
