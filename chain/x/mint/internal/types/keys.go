package types

// MinterKey is used for the keeper store
var MinterKey = []byte{0x00}

// nolint
const (
	// ModuleName
	ModuleName = "mint"

	WrappedModuleName = "odin" + ModuleName

	// DefaultParamspace params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
	EthIntegrationAddress = "eth_integration_address"
)
