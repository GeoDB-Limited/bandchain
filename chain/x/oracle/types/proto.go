package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DecProto struct {
	sdk.Dec
}

func NewDecProto() DecProto {
	return DecProto{
		Dec: sdk.NewDec(0),
	}
}

func (p DecProto) Size() int {
	return len(p.Bytes())
}

func (p DecProto) MarshalTo(dst []byte) (int, error) {
	res, err := p.MarshalJSON()
	if err != nil {
		return 0, err
	}
	copy(dst, res)
	return len(dst), nil
}

func (p DecProto) Unmarshal(dst []byte) error {
	return p.Dec.UnmarshalJSON(dst)
}

func (p DecProto) Equal(other DecProto) bool {
	return p.Dec.Equal(other.Dec)
}
