package types

import "github.com/GeoDB-Limited/odincore/chain/x/common/types"

type ValidExchanges map[string][]string

func (v ValidExchanges) Contains(from types.Denom, to types.Denom) bool {
	exchanges, ok := v[from.String()]
	if !ok {
		return false
	}
	for _, e := range exchanges {
		if types.Denom(e).Equal(to) {
			return true
		}
	}
	return false
}
