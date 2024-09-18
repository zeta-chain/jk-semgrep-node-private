package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const payloadMessageZEVMCall = "this is a test ZEVM call payload"

func TestV2EVMToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.AssertTestDAppZEVMCalled(false, payloadMessageZEVMCall, big.NewInt(0))

	// perform the withdraw
	tx := r.V2EVMToZEMVCall(
		r.TestDAppV2ZEVMAddr,
		[]byte(payloadMessageZEVMCall),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payloadMessageZEVMCall, big.NewInt(0))
}
