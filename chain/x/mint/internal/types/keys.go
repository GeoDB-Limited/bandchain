package types

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
)

var (
	GlobalStoreKeyPrefix = []byte{0x00}

	// AccountsPoolStoreKey is used for the eligible accounts store
	AccountsPoolStoreKey = append(GlobalStoreKeyPrefix, []byte("AccountsPool")...)

	// MinterKey is used for the keeper store
	MinterKey = append(GlobalStoreKeyPrefix, []byte("Minter")...)
)
