package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestV2ETHWithdrawAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHWithdrawAndCall")

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// use a random address to get the revert amount
	revertAddress := sample.EthAddress()
	balance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, int64(0), balance.Int64())

	// perform the withdraw
	tx := r.V2ETHWithdrawAndCall(r.TestDAppV2EVMAddr, amount, r.EncodeGasCall("revert"), gatewayzevm.RevertOptions{
		RevertAddress:    revertAddress,
		OnRevertGasLimit: big.NewInt(0),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)

	// check the balance is more than 0
	balance, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.True(r, balance.Cmp(big.NewInt(0)) > 0)
}
