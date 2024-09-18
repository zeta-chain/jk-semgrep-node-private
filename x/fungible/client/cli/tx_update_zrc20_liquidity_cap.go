package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/fungible/types"
)

func CmdUpdateZRC20LiquidityCap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-zrc20-liquidity-cap [zrc20] [liquidity-cap]",
		Short: "Broadcast message UpdateZRC20LiquidityCap",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			newCap := math.NewUintFromString(args[1])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateZRC20LiquidityCap(
				clientCtx.GetFromAddress().String(),
				args[0],
				newCap,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
