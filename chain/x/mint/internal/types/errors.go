package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidMintDenom     = sdkerrors.Register(ModuleName, 1, "The given mint denom is invalid")
	ErrAccountIsNotEligible = sdkerrors.Register(ModuleName, 2, "The given account is not eligible to mint")
	ErrExceedsMintLimit     = sdkerrors.Register(ModuleName, 3, "The given amount exceeds the limit to mint")
	ErrInvalidMintAmount    = sdkerrors.Register(ModuleName, 4, "The given amount to mint is invalid")
)
