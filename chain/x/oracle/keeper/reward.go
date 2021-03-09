package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDataSourceID(ctx sdk.Context, rid types.RequestID, eid types.ExternalID, did types.DataSourceID) {
	key := types.DataSourceByExternalIDPrefixKey(rid, eid)
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshalBinaryBare(did))
}

func (k Keeper) GetDataSourceID(ctx sdk.Context, rid types.RequestID, eid types.ExternalID) types.DataSourceID {
	key := types.DataSourceByExternalIDPrefixKey(rid, eid)
	bz := ctx.KVStore(k.storeKey).Get(key)
	var did types.DataSourceID
	k.cdc.MustUnmarshalBinaryBare(bz, &did)
	return did
}

func (k Keeper) SetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress, reward sdk.Dec) {
	key := types.DataProviderRewardsPrefixKey(acc)
	if !k.HasDataProviderReward(ctx, acc) {
		ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshalBinaryBare(reward))
		return
	}
	oldReward := k.GetDataProviderAccumulatedReward(ctx, acc)
	newReward := oldReward.Add(reward)
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshalBinaryBare(newReward))
}

func (k Keeper) ClearDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.DataProviderRewardsPrefixKey(acc))
}

func (k Keeper) GetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) (reward sdk.Dec) {
	key := types.DataProviderRewardsPrefixKey(acc)
	bz := ctx.KVStore(k.storeKey).Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &reward)
	return reward
}

func (k Keeper) HasDataProviderReward(ctx sdk.Context, acc sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.DataProviderRewardsPrefixKey(acc))
}

// todo refactor to common/Denom from coinswap
func (k Keeper) SetOracleDataProviderRewardDenom(ctx sdk.Context, denom string) {
	ctx.KVStore(k.storeKey).Set(types.OracleDataProviderRewardDenomStoreKey, k.cdc.MustMarshalBinaryBare(denom))
}

func (k Keeper) GetOracleDataProviderRewardDenom(ctx sdk.Context) (denom string) {
	bz := ctx.KVStore(k.storeKey).Get(types.OracleDataProviderRewardDenomStoreKey)
	k.cdc.MustUnmarshalBinaryBare(bz, &denom)
	return denom
}

func (k Keeper) AllocateRewardsToDataProviders(ctx sdk.Context, rid types.RequestID) {
	request := k.MustGetRequest(ctx, rid)

	for _, rawReq := range request.RawRequests {
		ds := k.MustGetDataSource(ctx, rawReq.GetDataSourceID())
		if !k.HasDataProviderReward(ctx, ds.Owner) {
			continue
		}
		//reward := k.GetDataProviderAccumulatedReward(ctx, ds.Owner)
		//dataProviderRewardDenom := k.GetOracleDataProviderRewardDenom(ctx)
		//rewardCoin := sdk.NewDecCoinFromDec(dataProviderRewardDenom, reward)

		//k.supplyKeeper.SendCoinsFromModuleToAccount()
	}
}
