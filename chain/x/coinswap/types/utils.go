package types

const (
	validExchangeFrom = "geo"
	validExchangeTo   = "loki"
)

// rewriting this simple logic allows to restrict specific denominations changes
func ValidExchangeDenom(from, to Denom) bool {
	if from.Base() == validExchangeFrom && to.Base() == validExchangeTo {
		return true
	}
	return false
}
