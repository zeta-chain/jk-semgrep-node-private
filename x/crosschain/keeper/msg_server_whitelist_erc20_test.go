package keeper_test

import (
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_WhitelistERC20(t *testing.T) {
	t.Run("can deploy and whitelist an erc20", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)

		admin := sample.AccAddress()
		erc20Address := sample.EthAddress().Hex()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")
		k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, sample.Tss())
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{1},
		})

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		res, err := msgServer.WhitelistERC20(ctx, &msg)
		require.NoError(t, err)
		require.NotNil(t, res)
		zrc20 := res.Zrc20Address
		cctxIndex := res.CctxIndex

		// check zrc20 and cctx created
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(zrc20))
		fc, found := zk.FungibleKeeper.GetForeignCoins(ctx, zrc20)
		require.True(t, found)
		require.EqualValues(t, "foo", fc.Name)
		require.EqualValues(t, erc20Address, fc.Asset)
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		require.True(t, found)
		require.EqualValues(t, fmt.Sprintf("%s:%s", constant.CmdWhitelistERC20, erc20Address), cctx.RelayedMessage)

		// check gas limit is set
		gasLimit, err := zk.FungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(zrc20))
		require.NoError(t, err)
		require.Equal(t, uint64(100000), gasLimit.Uint64())

		msgNew := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: sample.EthAddress().Hex(),
			ChainId:      chainID,
			Name:         "bar",
			Symbol:       "BAR",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msgNew, nil)

		// Ensure that whitelist a new erc20 create a cctx with a different index
		res, err = msgServer.WhitelistERC20(ctx, &msgNew)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEqual(t, cctxIndex, res.CctxIndex)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: sample.EthAddress().Hex(),
			ChainId:      getValidEthChainID(),
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.WhitelistERC20(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if invalid erc20 address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: "invalid",
			ChainId:      getValidEthChainID(),
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		_, err := msgServer.WhitelistERC20(ctx, &msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if foreign coin already exists for the asset", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		asset := sample.EthAddress().Hex()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		fc := sample.ForeignCoins(t, sample.EthAddress().Hex())
		fc.Asset = asset
		fc.ForeignChainId = chainID
		zk.FungibleKeeper.SetForeignCoins(ctx, fc)

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: asset,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.WhitelistERC20(ctx, &msg)
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinAlreadyExist)
	})

	t.Run("should fail if no tss set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		erc20Address := sample.EthAddress().Hex()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.WhitelistERC20(ctx, &msg)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("should fail if nox valid chain ID", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		erc20Address := sample.EthAddress().Hex()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, sample.Tss())

		msg := types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      10000,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.WhitelistERC20(ctx, &msg)
		require.ErrorIs(t, err, types.ErrInvalidChainID)
	})
}
