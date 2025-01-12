package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	BasicName          = "BASIC_NAME"
	BasicDesc          = "BASIC_DESCRIPTION"
	BasicSchema        = "BASIC_SCHEMA"
	BasicSourceCodeURL = "BASIC_SOURCE_CODE_URL"
	BasicFilename      = "BASIC_FILENAME"
	BasicCalldata      = []byte("BASIC_CALLDATA")
	BasicClientID      = "BASIC_CLIENT_ID"
	BasicReport        = []byte("BASIC_REPORT")
	BasicResult        = []byte("BASIC_RESULT")
	CoinsZero          = sdk.NewCoins()
	Coins10odin        = sdk.NewCoins(sdk.NewInt64Coin("odin", 10))
	Coins20odin        = sdk.NewCoins(sdk.NewInt64Coin("odin", 20))
	Coins1000000odin   = sdk.NewCoins(sdk.NewInt64Coin("odin", 1000000))
)
