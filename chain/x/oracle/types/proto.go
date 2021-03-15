package types

import (
	"github.com/GeoDB-Limited/odincore/chain/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CoinDecProto sdk.DecCoin

func NewCoinDecProto(denom types.Denom) CoinDecProto {
	return CoinDecProto(sdk.NewInt64DecCoin(denom.String(), 0))
}

func (p CoinDecProto) Value() sdk.DecCoin {
	return sdk.DecCoin(p)
}

func (p CoinDecProto) Size() int {
	return len(ModuleCdc.MustMarshalJSON(p))
}

func (p CoinDecProto) MarshalTo(dst []byte) (int, error) {
	res, err := ModuleCdc.MarshalJSON(p)
	if err != nil {
		return 0, err
	}
	copy(dst, res)
	return len(dst), nil
}

func (p CoinDecProto) Unmarshal(src []byte) error {
	return ModuleCdc.UnmarshalJSON(src, &p)
}

func (p CoinDecProto) Equal(other CoinDecProto) bool {
	return p.Denom == other.Denom && p.Amount.Equal(other.Amount)
}
