package mint

import (
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper types.SupplyKeeper, data GenesisState) {
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	keeper.SetAccPool(ctx, data.AccPool)

	moduleAcc := keeper.GetMinterAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
	}

	supplyKeeper.SetModuleAccount(ctx, moduleAcc)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinter(ctx)
	params := keeper.GetParams(ctx)
	accPool := keeper.GetAccPool(ctx)
	return NewGenesisState(minter, params, accPool)
}
