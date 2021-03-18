package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidMintDenom     = sdkerrors.Register(ModuleName, 1, "The given mint denom is invalid")
	ErrAccountIsNotEligible = sdkerrors.Register(ModuleName, 2, "The given account is not eligible to mint")
)
