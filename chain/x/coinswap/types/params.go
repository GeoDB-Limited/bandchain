package types

import (
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

type Params struct {
	// todo work on params
}

func (p Params) ParamSetPairs() subspace.ParamSetPairs {
	panic("implement me!")
	return nil
}

func DefaultParams() Params {
	return Params{}
}
