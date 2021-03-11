package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultFromExchange = "geo"
	DefaultToExchange   = "odin"
)

type Params struct {
	RateMultiplier sdk.Dec        `json:"rate_multiplier" yaml:"rate_multiplier"`
	ValidExchanges ValidExchanges `json:"valid_exchanges" yaml:"valid_exchanges"`
}

// nolint
var (
	KeyRateMultiplier = []byte("RateMultiplier")
	KeyValidExchanges = []byte("ValidExchanges")
)

// ParamTable for coinswap module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyRateMultiplier, &p.RateMultiplier, validateRateMultiplier),
		params.NewParamSetPair(KeyValidExchanges, &p.ValidExchanges, validatePossibleExchanges),
	}
}

func DefaultParams() Params {
	return Params{
		RateMultiplier: sdk.NewDec(1),
		ValidExchanges: ValidExchanges{
			DefaultFromExchange: []string{DefaultToExchange},
		},
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

func validatePossibleExchanges(i interface{}) error {
	exchanges, ok := i.(ValidExchanges)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for k, valid := range exchanges {
		for _, v := range valid {
			if k == "" || v == "" {
				return fmt.Errorf("one or both denoms are empty. From: %s, To: %s", k, v)
			}
		}
	}
	return nil
}
