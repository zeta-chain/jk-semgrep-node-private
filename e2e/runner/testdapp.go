package runner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/contracts/testdappv2"
)

// AssertTestDAppZEVMCalled is a function that asserts the values of the test dapp on the ZEVM
// this function uses TestDAppV2 for the assertions, in the future we should only use this contracts for all tests
// https://github.com/zeta-chain/node/issues/2655
func (r *E2ERunner) AssertTestDAppZEVMCalled(expectedCalled bool, message string, amount *big.Int) {
	r.assertTestDAppCalled(r.TestDAppV2ZEVM, message, expectedCalled, amount)
}

// AssertTestDAppEVMCalled is a function that asserts the values of the test dapp on the external EVM
func (r *E2ERunner) AssertTestDAppEVMCalled(expectedCalled bool, message string, amount *big.Int) {
	r.assertTestDAppCalled(r.TestDAppV2EVM, message, expectedCalled, amount)
}

func (r *E2ERunner) assertTestDAppCalled(
	testDApp *testdappv2.TestDAppV2,
	message string,
	expectedCalled bool,
	expectedAmount *big.Int,
) {
	// check the payload was received on the contract
	called, err := testDApp.GetCalledWithMessage(&bind.CallOpts{}, message)
	require.NoError(r, err)
	require.EqualValues(r, expectedCalled, called)

	if expectedCalled {
		amount, err := testDApp.GetAmountWithMessage(&bind.CallOpts{}, message)
		require.NoError(r, err)
		require.EqualValues(
			r,
			expectedAmount.Uint64(),
			amount.Uint64(),
			"Amounts do not match, expected %s, actual %s",
			expectedAmount.String(),
			amount.String(),
		)
	}
}

// EncodeGasCall encodes the payload for the gasCall function
func (r *E2ERunner) EncodeGasCall(message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("gasCall", message)
	require.NoError(r, err)
	return encoded
}

// EncodeGasCallRevert encodes the payload for the gasCall function that reverts
func (r *E2ERunner) EncodeGasCallRevert() []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("gasCall", "revert")
	require.NoError(r, err)
	return encoded
}

// EncodeERC20Call encodes the payload for the erc20Call function
func (r *E2ERunner) EncodeERC20Call(erc20Addr ethcommon.Address, amount *big.Int, message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("erc20Call", erc20Addr, amount, message)
	require.NoError(r, err)
	return encoded
}

// EncodeERC20CallRevert encodes the payload for the erc20Call function that reverts
func (r *E2ERunner) EncodeERC20CallRevert(erc20Addr ethcommon.Address, amount *big.Int) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("erc20Call", erc20Addr, amount, "revert")
	require.NoError(r, err)
	return encoded
}

// EncodeSimpleCall encodes the payload for the simpleCall function
func (r *E2ERunner) EncodeSimpleCall(message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("simpleCall", message)
	require.NoError(r, err)
	return encoded
}
