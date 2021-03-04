package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paramSpace   params.Subspace
	supplyKeeper types.SupplyKeeper
	distrKeeper  types.DistrKeeper
}

func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	subspace params.Subspace,
	sk types.SupplyKeeper,
	dk types.DistrKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		paramSpace:   subspace,
		supplyKeeper: sk,
		distrKeeper:  dk,
	}
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.paramSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.paramSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}
