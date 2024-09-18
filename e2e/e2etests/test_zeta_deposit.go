package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestZetaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount
	amount := parseBigInt(r, args[0])

	hash := r.DepositZetaWithAmount(r.EVMAddress(), amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}
