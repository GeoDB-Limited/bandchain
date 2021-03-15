package types

import (
	"github.com/GeoDB-Limited/odincore/chain/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CoinProto sdk.Coin

func NewCoinProto(denom types.Denom) CoinProto {
	return CoinProto(sdk.NewInt64Coin(denom.String(), 0))
}

func (p CoinProto) Value() sdk.Coin {
	return sdk.Coin(p)
}

func (p CoinProto) Size() int {
	return len(ModuleCdc.MustMarshalJSON(p))
}

func (p CoinProto) MarshalTo(dst []byte) (int, error) {
	res, err := ModuleCdc.MarshalJSON(p)
	if err != nil {
		return 0, err
	}
	copy(dst, res)
	return len(dst), nil
}

func (p CoinProto) Unmarshal(src []byte) error {
	return ModuleCdc.UnmarshalJSON(src, &p)
}

func (p CoinProto) Equal(other CoinProto) bool {
	return p.Denom == other.Denom && p.Amount.Equal(other.Amount)
}
