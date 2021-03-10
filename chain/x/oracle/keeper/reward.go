package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

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

// todo optimize store queries
// sends rewards from oracle pool to data providers, that have given data for the passed request
func (k Keeper) AllocateRewardsToDataProviders(ctx sdk.Context, rid types.RequestID) {
	logger := k.Logger(ctx)
	request := k.MustGetRequest(ctx, rid)

	for _, rawReq := range request.RawRequests {
		ds := k.MustGetDataSource(ctx, rawReq.GetDataSourceID())
		if !k.HasDataProviderReward(ctx, ds.Owner) {
			continue
		}
		reward := k.GetDataProviderAccumulatedReward(ctx, ds.Owner)
		dataProviderRewardDenom := k.GetOracleDataProviderRewardDenom(ctx)
		rewardCoinDec := sdk.NewDecCoinFromDec(dataProviderRewardDenom, reward)

		// rewards are lying in the distribution fee pool
		feePool := k.distrKeeper.GetFeePool(ctx)
		diff, hasNeg := feePool.CommunityPool.SafeSub(sdk.NewDecCoins(rewardCoinDec))
		if hasNeg {
			logger.With("lack", diff, "denom", dataProviderRewardDenom).Error("oracle pool does not have enough coins to reward data providers")
			// not return because maybe still enough coins to pay someone
			continue
		}
		feePool.CommunityPool = diff

		rewardCoin, _ := rewardCoinDec.TruncateDecimal()
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, distr.ModuleName, ds.Owner, sdk.NewCoins(rewardCoin))
		if err != nil {
			panic(err)
		}

		k.distrKeeper.SetFeePool(ctx, feePool)
		k.ClearDataProviderAccumulatedReward(ctx, ds.Owner)
	}
}
