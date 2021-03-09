package types

import "github.com/cosmos/cosmos-sdk/codec"

// ModuleCdc is the codec for the module.
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers the module's concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgExchange{}, "coinswap/Exchange", nil)
}
