package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Params struct {
	RateMultiplier sdk.Dec `json:"rate_multiplier" yaml:"rate_multiplier"`
}

// nolint
var (
	KeyRateMultiplier = []byte("RateMultiplier")
)

// ParamTable for coinswap module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyRateMultiplier, &p.RateMultiplier, validateRateMultiplier),
	}
}

func DefaultParams() Params {
	return Params{
		RateMultiplier: sdk.NewDec(1),
	}
}

func validateRateMultiplier(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsPositive() && !v.IsZero() {
		return fmt.Errorf("rate multiplier %s must be positive or zero", v)
	}
	return nil
}
