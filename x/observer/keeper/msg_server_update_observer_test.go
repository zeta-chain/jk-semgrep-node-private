package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_UpdateObserver(t *testing.T) {
	t.Run("successfully update tombstoned observer", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		count := uint64(0)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count = 1

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.NoError(t, err)
		acc, found := k.GetNodeAccount(ctx, newOperatorAddress.String())
		require.True(t, found)
		require.Equal(t, newOperatorAddress.String(), acc.Operator)
	})

	t.Run(
		"unable to update a tombstoned observer if the new address already exists in the observer set",
		func(t *testing.T) {
			//ARRANGE
			k, ctx, _, _ := keepertest.ObserverKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			// #nosec G404 test purpose - weak randomness is not an issue here
			r := rand.New(rand.NewSource(9))
			// Set validator in the store
			validator := sample.Validator(t, r)
			validatorNew := sample.Validator(t, r)
			validatorNew.Status = stakingtypes.Bonded
			k.GetStakingKeeper().SetValidator(ctx, validatorNew)
			k.GetStakingKeeper().SetValidator(ctx, validator)

			consAddress, err := validator.GetConsAddr()
			require.NoError(t, err)
			k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
				Address:             consAddress.String(),
				StartHeight:         0,
				JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
				Tombstoned:          true,
				MissedBlocksCounter: 1,
			})

			accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
			require.NoError(t, err)

			newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
			require.NoError(t, err)

			observerList := []string{accAddressOfValidator.String(), newOperatorAddress.String()}
			k.SetObserverSet(ctx, types.ObserverSet{
				ObserverList: observerList,
			})
			k.SetNodeAccount(ctx, types.NodeAccount{
				Operator: accAddressOfValidator.String(),
			})
			k.SetLastObserverCount(ctx, &types.LastObserverCount{
				Count: uint64(len(observerList)),
			})

			//ACT
			_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
				Creator:            accAddressOfValidator.String(),
				OldObserverAddress: accAddressOfValidator.String(),
				NewObserverAddress: newOperatorAddress.String(),
				UpdateReason:       types.ObserverUpdateReason_Tombstoned,
			})

			// ASSERT
			require.ErrorContains(t, err, types.ErrDuplicateObserver.Error())
		},
	)

	t.Run("unable to update to a non validator address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count = 1
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.ErrorIs(t, err, types.ErrUpdateObserver)
	})

	t.Run("unable to update tombstoned validator with with non operator account", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            sample.AccAddress(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.ErrorIs(t, err, types.ErrUpdateObserver)
	})

	t.Run("unable to update non-tombstoned observer with update reason tombstoned", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.ErrorIs(t, err, types.ErrUpdateObserver)
	})

	t.Run("unable to update observer with no node account", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.ErrorIs(t, err, types.ErrNodeAccountNotFound)
	})

	t.Run("unable to update observer when last observer count is missing", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		require.ErrorIs(t, err, types.ErrLastObserverCountNotFound)
	})

	t.Run("update observer using admin policy", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateObserver{
			Creator:            admin,
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		acc, found := k.GetNodeAccount(ctx, newOperatorAddress.String())
		require.True(t, found)
		require.Equal(t, newOperatorAddress.String(), acc.Operator)
	})

	t.Run("fail to update observer using regular account and update type admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		count := uint64(0)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		count += 1

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		require.NoError(t, err)

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            sample.AccAddress(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
		})
		require.ErrorIs(t, err, types.ErrUpdateObserver)
	})
}

func TestUpdateObserverList(t *testing.T) {
	t.Run("update observer list", func(t *testing.T) {
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		list := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress(), oldObserverAddress}
		require.Equal(t, oldObserverAddress, list[3])
		keeper.UpdateObserverList(list, oldObserverAddress, newObserverAddress)
		require.Equal(t, 4, len(list))
		require.Equal(t, newObserverAddress, list[3])
	})
}
