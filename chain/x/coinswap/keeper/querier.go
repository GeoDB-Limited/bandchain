package keeper

import (
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
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
			return queryRate(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown coinswap query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	return types.QueryOK(k.GetParams(ctx))
}

func queryRate(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 2 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "not all the arguments are specified")
	}
	from, err := types.ParseDenom(path[0])
	if err != nil {
		return types.QueryBadRequest(fmt.Sprintf("%s - is not a valid denom", path[0]))
	}
	to, err := types.ParseDenom(path[1])
	if err != nil {
		return types.QueryBadRequest(fmt.Sprintf("%s - is not a valid denom", path[1]))
	}

	return types.QueryOK(types.QueryRateResult{
		Rate: k.GetRate(ctx, from, to),
		From: from,
		To:   to,
	})
}
