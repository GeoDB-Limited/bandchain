package types

// GenesisState - minter state
type GenesisState struct {
	Minter  Minter  `json:"minter" yaml:"minter"`                       // minter object
	Params  Params  `json:"params" yaml:"params"`                       // inflation params
	AccPool AccPool `json:"eligible_accounts" yaml:"eligible_accounts"` // pool of eligible accounts
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params, accPool AccPool) GenesisState {
	return GenesisState{
		Minter:  minter,
		Params:  params,
		AccPool: accPool,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter:  DefaultInitialMinter(),
		Params:  DefaultParams(),
		AccPool: AccPool{},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return ValidateMinter(data.Minter)
}
