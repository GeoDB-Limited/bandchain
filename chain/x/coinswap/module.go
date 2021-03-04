package coinswap

import (
	"encoding/json"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/client/cli"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/client/rest"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is Band Oracle's module basic object.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}

func (b AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (b AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	return ModuleCdc.UnmarshalJSON(bz, &data)
}

func (b AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, router *mux.Router) {
	rest.RegisterRoutes(ctx, router, StoreKey)
}

func (b AppModuleBasic) GetTxCmd(codec *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, codec)
}

func (b AppModuleBasic) GetQueryCmd(codec *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, codec)
}

// Name returns this module's name - "coinswap" (SDK AppModuleBasic interface).
func (AppModuleBasic) Name() string { return ModuleName }

// AppModule represents the AppModule for this module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

func (am AppModule) InitGenesis(ctx sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(message, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	// todo maybe need one for uniswap
}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {

}

func (am AppModule) EndBlock(s sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
