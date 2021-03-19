package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MintPool
type MintPool struct {
	TreasuryPool sdk.DecCoins `json:"treasury_pool" yaml:"treasury_pool"`
	EligiblePool AddrPool     `json:"eligible_accounts_pool" yaml:"eligible_accounts_pool"` // eligible to mint accounts
}

// InitialMintPool returns the initial state of MintPool
func InitialMintPool() MintPool {
	return MintPool{
		TreasuryPool: sdk.DecCoins{},
		EligiblePool: AddrPool{},
	}
}

// ValidateGenesis validates the mint pool for a genesis state
func (m MintPool) ValidateGenesis() error {
	if m.TreasuryPool.IsAnyNegative() {
		return fmt.Errorf("negative TreasuryPool in mint pool, is %v", m.TreasuryPool)
	}

	return nil
}

// AddrPool defines a pool of addresses
type AddrPool []sdk.AccAddress

// Contains checks id addr exists in the slice
func (p *AddrPool) Contains(addr sdk.AccAddress) bool {
	for _, item := range *p {
		if item.Equals(addr) {
			return true
		}
	}
	return false
}
