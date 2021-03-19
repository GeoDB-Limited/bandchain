package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"strings"
)

// Parameter store keys
var (
	KeyMintDenom            = []byte("MintDenom")
	KeyInflationRateChange  = []byte("InflationRateChange")
	KeyInflationMax         = []byte("InflationMax")
	KeyInflationMin         = []byte("InflationMin")
	KeyGoalBonded           = []byte("GoalBonded")
	KeyBlocksPerYear        = []byte("BlocksPerYear")
	KeyMintAir              = []byte("MintAir")
	KeyMaxWithdrawalPerTime = []byte("MaxWithdrawalPerTime")
)

// Params defines a mint parameters
type Params struct {
	MintDenom            string    `json:"mint_denom" yaml:"mint_denom"`                           // type of coin to mint
	InflationRateChange  sdk.Dec   `json:"inflation_rate_change" yaml:"inflation_rate_change"`     // maximum annual change in inflation rate
	InflationMax         sdk.Dec   `json:"inflation_max" yaml:"inflation_max"`                     // maximum inflation rate
	InflationMin         sdk.Dec   `json:"inflation_min" yaml:"inflation_min"`                     // minimum inflation rate
	GoalBonded           sdk.Dec   `json:"goal_bonded" yaml:"goal_bonded"`                         // goal of percent bonded atoms
	BlocksPerYear        uint64    `json:"blocks_per_year" yaml:"blocks_per_year"`                 // expected blocks per year
	MintAir              bool      `json:"mint_air" yaml:"mint_air"`                               // flag if mint from air
	MaxWithdrawalPerTime sdk.Coins `json:"max_withdrawal_per_time" yaml:"max_withdrawal_per_time"` // max to mint in one withdraw
}

// ParamKeyTable defines a key table for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams returns a new mint params
func NewParams(
	mintDenom string,
	inflationRateChange, inflationMax, inflationMin, goalBonded sdk.Dec,
	MaxWithdrawalPerTime sdk.Coins,
	blocksPerYear uint64,
	mintAir bool,
) Params {

	return Params{
		MintDenom:            mintDenom,
		InflationRateChange:  inflationRateChange,
		InflationMax:         inflationMax,
		InflationMin:         inflationMin,
		GoalBonded:           goalBonded,
		BlocksPerYear:        blocksPerYear,
		MintAir:              mintAir,
		MaxWithdrawalPerTime: MaxWithdrawalPerTime,
	}
}

// DefaultParams returns default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:            sdk.DefaultBondDenom,
		InflationRateChange:  sdk.NewDecWithPrec(13, 2),
		InflationMax:         sdk.NewDecWithPrec(20, 2),
		InflationMin:         sdk.NewDecWithPrec(7, 2),
		GoalBonded:           sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:        uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		MintAir:              false,
		MaxWithdrawalPerTime: sdk.Coins{},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateMintAir(p.MintAir); err != nil {
		return err
	}
	if err := validateMaxWithdrawalPerTime(p.MaxWithdrawalPerTime); err != nil {
		return err
	}
	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf(
			"max inflation (%s) must be greater than or equal to min inflation (%s)",
			p.InflationMax, p.InflationMin,
		)
	}

	return nil

}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
	Mint Denom:             %s
	Inflation Rate Change:  %s
	Inflation Max:          %s
	Inflation Min:          %s
	Goal Bonded:            %s
	Blocks Per Year:        %d
	Mint Max:				%s
`,
		p.MintDenom, p.InflationRateChange, p.InflationMax,
		p.InflationMin, p.GoalBonded, p.BlocksPerYear, p.MaxWithdrawalPerTime,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		params.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRateChange),
		params.NewParamSetPair(KeyInflationMax, &p.InflationMax, validateInflationMax),
		params.NewParamSetPair(KeyInflationMin, &p.InflationMin, validateInflationMin),
		params.NewParamSetPair(KeyGoalBonded, &p.GoalBonded, validateGoalBonded),
		params.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		params.NewParamSetPair(KeyMintAir, &p.MintAir, validateMintAir),
		params.NewParamSetPair(KeyMaxWithdrawalPerTime, &p.MaxWithdrawalPerTime, validateMaxWithdrawalPerTime),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return sdkerrors.Wrap(ErrInvalidMintDenom, "mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRateChange(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateInflationMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}

func validateInflationMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min inflation too large: %s", v)
	}

	return nil
}

func validateGoalBonded(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("goal bonded cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("goal bonded too large: %s", v)
	}

	return nil
}

func validateMaxWithdrawalPerTime(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsValid() {
		return fmt.Errorf("max withdrawal per time parameter is not valid: %s", v)
	}
	if v.IsAnyNegative() {
		return fmt.Errorf("max withdrawal per time cannot be negative: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateMintAir(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
