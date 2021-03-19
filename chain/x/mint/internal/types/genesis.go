package types

// GenesisState - minter state
type GenesisState struct {
	Minter   Minter   `json:"minter" yaml:"minter"`       // minter object
	Params   Params   `json:"params" yaml:"params"`       // inflation params
	MintPool MintPool `json:"mint_pool" yaml:"mint_pool"` // module pool
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params, mintPool MintPool) GenesisState {
	return GenesisState{
		Minter:   minter,
		Params:   params,
		MintPool: mintPool,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter:   DefaultInitialMinter(),
		Params:   DefaultParams(),
		MintPool: InitialMintPool(),
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
