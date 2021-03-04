package coinswap

import (
	"encoding/json"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the oracle state that must be provided at genesis.
type GenesisState struct {
	Params types.Params `json:"params" yaml:"params"`
}

// DefaultGenesisState returns the default oracle genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	// todo set params
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params: k.GetParams(ctx),
	}
}

// GetGenesisStateFromAppState returns x/coinswap GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
