// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package prototype

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

// IPrototypeMetaData contains all meta data concerning the IPrototype contract.
var IPrototypeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32ToHexAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"prefix\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"bech32ify\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int64\",\"name\":\"chainID\",\"type\":\"int64\"}],\"name\":\"getGasStabilityPoolBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IPrototypeABI is the input ABI used to generate the binding from.
// Deprecated: Use IPrototypeMetaData.ABI instead.
var IPrototypeABI = IPrototypeMetaData.ABI

// IPrototype is an auto generated Go binding around an Ethereum contract.
type IPrototype struct {
	IPrototypeCaller     // Read-only binding to the contract
	IPrototypeTransactor // Write-only binding to the contract
	IPrototypeFilterer   // Log filterer for contract events
}

// IPrototypeCaller is an auto generated read-only Go binding around an Ethereum contract.
type IPrototypeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPrototypeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IPrototypeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPrototypeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPrototypeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPrototypeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPrototypeSession struct {
	Contract     *IPrototype       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPrototypeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPrototypeCallerSession struct {
	Contract *IPrototypeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// IPrototypeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPrototypeTransactorSession struct {
	Contract     *IPrototypeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// IPrototypeRaw is an auto generated low-level Go binding around an Ethereum contract.
type IPrototypeRaw struct {
	Contract *IPrototype // Generic contract binding to access the raw methods on
}

// IPrototypeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPrototypeCallerRaw struct {
	Contract *IPrototypeCaller // Generic read-only contract binding to access the raw methods on
}

// IPrototypeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPrototypeTransactorRaw struct {
	Contract *IPrototypeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPrototype creates a new instance of IPrototype, bound to a specific deployed contract.
func NewIPrototype(address common.Address, backend bind.ContractBackend) (*IPrototype, error) {
	contract, err := bindIPrototype(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPrototype{IPrototypeCaller: IPrototypeCaller{contract: contract}, IPrototypeTransactor: IPrototypeTransactor{contract: contract}, IPrototypeFilterer: IPrototypeFilterer{contract: contract}}, nil
}

// NewIPrototypeCaller creates a new read-only instance of IPrototype, bound to a specific deployed contract.
func NewIPrototypeCaller(address common.Address, caller bind.ContractCaller) (*IPrototypeCaller, error) {
	contract, err := bindIPrototype(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPrototypeCaller{contract: contract}, nil
}

// NewIPrototypeTransactor creates a new write-only instance of IPrototype, bound to a specific deployed contract.
func NewIPrototypeTransactor(address common.Address, transactor bind.ContractTransactor) (*IPrototypeTransactor, error) {
	contract, err := bindIPrototype(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPrototypeTransactor{contract: contract}, nil
}

// NewIPrototypeFilterer creates a new log filterer instance of IPrototype, bound to a specific deployed contract.
func NewIPrototypeFilterer(address common.Address, filterer bind.ContractFilterer) (*IPrototypeFilterer, error) {
	contract, err := bindIPrototype(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPrototypeFilterer{contract: contract}, nil
}

// bindIPrototype binds a generic wrapper to an already deployed contract.
func bindIPrototype(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPrototypeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPrototype *IPrototypeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPrototype.Contract.IPrototypeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPrototype *IPrototypeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPrototype.Contract.IPrototypeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPrototype *IPrototypeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPrototype.Contract.IPrototypeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPrototype *IPrototypeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPrototype.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPrototype *IPrototypeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPrototype.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPrototype *IPrototypeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPrototype.Contract.contract.Transact(opts, method, params...)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_IPrototype *IPrototypeCaller) Bech32ToHexAddr(opts *bind.CallOpts, bech32 string) (common.Address, error) {
	var out []interface{}
	err := _IPrototype.contract.Call(opts, &out, "bech32ToHexAddr", bech32)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_IPrototype *IPrototypeSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _IPrototype.Contract.Bech32ToHexAddr(&_IPrototype.CallOpts, bech32)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_IPrototype *IPrototypeCallerSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _IPrototype.Contract.Bech32ToHexAddr(&_IPrototype.CallOpts, bech32)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_IPrototype *IPrototypeCaller) Bech32ify(opts *bind.CallOpts, prefix string, addr common.Address) (string, error) {
	var out []interface{}
	err := _IPrototype.contract.Call(opts, &out, "bech32ify", prefix, addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_IPrototype *IPrototypeSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _IPrototype.Contract.Bech32ify(&_IPrototype.CallOpts, prefix, addr)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_IPrototype *IPrototypeCallerSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _IPrototype.Contract.Bech32ify(&_IPrototype.CallOpts, prefix, addr)
}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_IPrototype *IPrototypeCaller) GetGasStabilityPoolBalance(opts *bind.CallOpts, chainID int64) (*big.Int, error) {
	var out []interface{}
	err := _IPrototype.contract.Call(opts, &out, "getGasStabilityPoolBalance", chainID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_IPrototype *IPrototypeSession) GetGasStabilityPoolBalance(chainID int64) (*big.Int, error) {
	return _IPrototype.Contract.GetGasStabilityPoolBalance(&_IPrototype.CallOpts, chainID)
}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_IPrototype *IPrototypeCallerSession) GetGasStabilityPoolBalance(chainID int64) (*big.Int, error) {
	return _IPrototype.Contract.GetGasStabilityPoolBalance(&_IPrototype.CallOpts, chainID)
}
