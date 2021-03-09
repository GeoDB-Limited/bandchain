package coinswap

import (
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
)

var (
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper      = keeper.Keeper
	MsgExchange = types.MsgExchange
)
