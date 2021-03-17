package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidAmountToMint  = sdkerrors.Register(ModuleName, 1, "The given amount to mint is invalid")
	ErrAccountIsNotEligible = sdkerrors.Register(ModuleName, 2, "The given account is not eligible to mint")
)
