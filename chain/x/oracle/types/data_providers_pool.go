package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// pool of <odin> tokens bought by data providers for <geo> tokens
type OraclePool struct {
	DataProvidersPool sdk.DecCoins `json:"data_providers_pool" yaml:"data_providers_pool"`
}

// zero oracle pool
func InitialOraclePool() OraclePool {
	return OraclePool{
		DataProvidersPool: sdk.DecCoins{},
	}
}

// ValidateGenesis validates the oracle pool for a genesis state
func (f OraclePool) ValidateGenesis() error {
	if f.DataProvidersPool.IsAnyNegative() {
		return fmt.Errorf("negative DataProvidersPool in oracle pool, is %v",
			f.DataProvidersPool)
	}

	return nil
}
