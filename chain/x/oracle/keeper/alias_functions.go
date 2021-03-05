package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetOracleAccount returns the oracle ModuleAccount
func (k Keeper) GetOracleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
