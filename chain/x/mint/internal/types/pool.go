package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// AccPool defines a pool of accounts
type AccPool []sdk.AccAddress

// Contains checks id addr exists in the slice
func (p *AccPool) Contains(addr sdk.AccAddress) bool {
	for _, item := range *p {
		if item.Equals(addr) {
			return true
		}
	}
	return false
}
