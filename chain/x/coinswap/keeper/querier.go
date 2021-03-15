package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	commontypes "github.com/GeoDB-Limited/odincore/chain/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryParams:
			return queryParameters(ctx, keeper)
		case types.QueryRate:
			return queryRate(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown coinswap query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	return commontypes.QueryOK(types.ModuleCdc, k.GetParams(ctx))
}

func queryRate(ctx sdk.Context, k Keeper) ([]byte, error) {
	initialRate := k.GetInitialRate(ctx)
	rateMultiplier := k.GetRateMultiplier(ctx)
	return commontypes.QueryOK(types.ModuleCdc, types.QueryRateResult{
		Rate:        initialRate.Mul(rateMultiplier),
		InitialRate: initialRate,
	})
}
