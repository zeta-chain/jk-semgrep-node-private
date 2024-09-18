package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgWhitelistERC20_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgWhitelistERC20
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgWhitelistERC20(
				"invalid_address",
				"0x0",
				1,
				"name",
				"symbol",
				6,
				10,
			),
			error: true,
		},
		{
			name: "invalid erc20",
			msg: types.NewMsgWhitelistERC20(
				sample.AccAddress(),
				"0x0",
				1,
				"name",
				"symbol",
				6,
				10,
			),
			error: true,
		},
		{
			name: "invalid decimals",
			msg: types.NewMsgWhitelistERC20(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				1,
				"name",
				"symbol",
				130,
				10,
			),
			error: true,
		},
		{
			name: "invalid gas limit",
			msg: types.NewMsgWhitelistERC20(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				1,
				"name",
				"symbol",
				6,
				-10,
			),
			error: true,
		},
		{
			name: "valid",
			msg: types.NewMsgWhitelistERC20(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				1,
				"name",
				"symbol",
				6,
				10,
			),
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgWhitelistERC20_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgWhitelistERC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgWhitelistERC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgWhitelistERC20{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgWhitelistERC20_Type(t *testing.T) {
	msg := types.MsgWhitelistERC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgWhitelistERC20, msg.Type())
}

func TestMsgWhitelistERC20_Route(t *testing.T) {
	msg := types.MsgWhitelistERC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgWhitelistERC20_GetSignBytes(t *testing.T) {
	msg := types.MsgWhitelistERC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
