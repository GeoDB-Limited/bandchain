package types

const (
	ExchangeDenomSep = "/"
)

type ExchangeDenom struct {
	From Denom
	To   Denom
}

func NewExchangeDenom(from, to Denom) ExchangeDenom {
	return ExchangeDenom{
		From: from,
		To:   to,
	}
}
