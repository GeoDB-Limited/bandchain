package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

func CalculateReward(data []byte, pricePerByte sdk.Dec) sdk.Dec {
	return pricePerByte.MulInt64(int64(len(data)))
}
