package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

func (k Keeper) ExchangeDenom(ctx sdk.Context, from, to types.Denom, amt sdk.Coin, requester sdk.AccAddress) error {

	// convert source amount to destination amount according to rate
	convertedAmt, err := k.convertToRate(ctx, from, to, amt)
	if err != nil {
		return sdkerrors.Wrap(err, "converting rate")
	}

	// first send source tokens to module
	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, requester, distr.ModuleName, sdk.NewCoins(amt))
	if err != nil {
		return sdkerrors.Wrapf(err, "sending coins from account: %s, to module: %s", requester.String(), distr.ModuleName)
	}

	toSend, remainder := convertedAmt.TruncateDecimal()
	if !remainder.IsZero() {
		k.Logger(ctx).With("coins", remainder.String()).Info("performing exchange according to limited precision some coins are lost")
	}

	feePool := k.distrKeeper.GetFeePool(ctx)

	// first add received tokens to fee pool
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(sdk.NewCoins(amt)...)...)

	k.distrKeeper.SetFeePool(ctx, feePool)

	oraclePool := k.oracleKeeper.GetOraclePool(ctx)

	// then subtract requested tokens from
	diff, hasNeg := oraclePool.DataProvidersPool.SafeSub(sdk.NewDecCoins(convertedAmt))
	if hasNeg {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "data providers pool does not have enough funds")
	}

	oraclePool.DataProvidersPool = diff
	k.oracleKeeper.SetOraclePool(ctx, oraclePool)

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, oracle.ModuleName, requester, sdk.NewCoins(toSend))
	if err != nil {
		return sdkerrors.Wrapf(err, "sending coins from module: %s, to account: %s", distr.ModuleName, requester.String())
	}

	return nil
}

func (k Keeper) GetRate(ctx sdk.Context, from, to types.Denom) sdk.Dec {
	totalSupply := k.supplyKeeper.GetSupply(ctx).GetTotal()
	fromSupply := totalSupply.AmountOf(from.String())
	toSupply := totalSupply.AmountOf(to.String())
	return fromSupply.ToDec().QuoRoundUp(toSupply.ToDec())
}

// returns the converted amount according to current rate
func (k Keeper) convertToRate(ctx sdk.Context, from, to types.Denom, amt sdk.Coin) (sdk.DecCoin, error) {
	rate := k.GetRate(ctx, from, to)
	if rate.GT(amt.Amount.ToDec()) {
		return sdk.DecCoin{}, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "current rate: %s is higher then amount provided: %s", rate.String(), amt.String())
	}
	convertedAmt := amt.Amount.ToDec().QuoRoundUp(rate)
	return sdk.NewDecCoinFromDec(to.String(), convertedAmt), nil
}
