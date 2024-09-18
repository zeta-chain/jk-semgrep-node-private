package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/fungible/types"
)

// SetForeignCoins set a specific foreignCoins in the store from its index
func (k Keeper) SetForeignCoins(ctx sdk.Context, foreignCoins types.ForeignCoins) {
	p := types.KeyPrefix(types.ForeignCoinsKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&foreignCoins)
	store.Set(types.ForeignCoinsKey(
		foreignCoins.Zrc20ContractAddress,
	), b)
}

// GetForeignCoins returns a foreignCoins from its index
func (k Keeper) GetForeignCoins(
	ctx sdk.Context,
	zrc20Addr string,
) (val types.ForeignCoins, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))

	b := store.Get(types.ForeignCoinsKey(
		zrc20Addr,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveForeignCoins removes a foreignCoins from the store
func (k Keeper) RemoveForeignCoins(
	ctx sdk.Context,
	zrc20Addr string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	store.Delete(types.ForeignCoinsKey(
		zrc20Addr,
	))
}

// GetAllForeignCoinsForChain returns all foreignCoins on a given chain
func (k Keeper) GetAllForeignCoinsForChain(ctx sdk.Context, foreignChainID int64) (list []types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ForeignCoins
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.ForeignChainId == foreignChainID {
			list = append(list, val)
		}
	}
	return
}

// GetAllForeignCoins returns all foreignCoins
func (k Keeper) GetAllForeignCoins(ctx sdk.Context) (list []types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ForeignCoins
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// GetAllForeignCoinMap returns all foreign ERC20 coins in a map of chainID -> asset -> coin
// Note: DO NOT use this method outside of gRPC queries
func (k Keeper) GetAllForeignCoinMap(ctx sdk.Context) map[int64]map[string]types.ForeignCoins {
	allForeignCoins := k.GetAllForeignCoins(ctx)

	foreignCoinMap := make(map[int64]map[string]types.ForeignCoins)
	for _, c := range allForeignCoins {
		if _, found := foreignCoinMap[c.ForeignChainId]; !found {
			foreignCoinMap[c.ForeignChainId] = make(map[string]types.ForeignCoins)
		}
		foreignCoinMap[c.ForeignChainId][strings.ToLower(c.Asset)] = c
	}
	return foreignCoinMap
}

// GetGasCoinForForeignCoin returns the gas coin for a given chain
func (k Keeper) GetGasCoinForForeignCoin(ctx sdk.Context, chainID int64) (types.ForeignCoins, bool) {
	foreignCoinList := k.GetAllForeignCoinsForChain(ctx, chainID)
	for _, c := range foreignCoinList {
		if c.CoinType == coin.CoinType_Gas {
			return c, true
		}
	}
	return types.ForeignCoins{}, false
}

// GetForeignCoinFromAsset returns the foreign coin for a given asset for a given chain
func (k Keeper) GetForeignCoinFromAsset(ctx sdk.Context, asset string, chainID int64) (types.ForeignCoins, bool) {
	if !ethcommon.IsHexAddress(asset) {
		return types.ForeignCoins{}, false
	}
	assetAddr := ethcommon.HexToAddress(asset)

	foreignCoinList := k.GetAllForeignCoinsForChain(ctx, chainID)
	for _, coin := range foreignCoinList {
		coinAssetAddr := ethcommon.HexToAddress(coin.Asset)
		if coinAssetAddr == assetAddr && coin.ForeignChainId == chainID {
			return coin, true
		}
	}
	return types.ForeignCoins{}, false
}
