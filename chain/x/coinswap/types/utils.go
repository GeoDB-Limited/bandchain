package types

import "github.com/GeoDB-Limited/odincore/chain/x/common/types"

const (
	validExchangeFrom = "geo"
	validExchangeTo   = "loki"
)

// rewriting this simple logic allows to restrict specific denominations changes
func ValidExchangeDenom(from, to types.Denom) bool {
	if from.Base() == validExchangeFrom && to.Base() == validExchangeTo {
		return true
	}
	return false
}
