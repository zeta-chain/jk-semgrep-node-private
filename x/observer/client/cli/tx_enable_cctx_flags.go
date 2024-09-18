package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

func CmdEnableCCTX() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable-cctx [enable-inbound] [enable-outbound]",
		Short: "Enable inbound and outbound for CCTX",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			enableInbound, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			enableOutbound, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgEnableCCTX(clientCtx.GetFromAddress().String(), enableInbound, enableOutbound)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
