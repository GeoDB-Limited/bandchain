package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidMintDenom                     = sdkerrors.Register(ModuleName, 111, "The given mint denom is invalid")
	ErrAccountIsNotEligible                 = sdkerrors.Register(ModuleName, 112, "The given account is not eligible to mint")
	ErrInvalidWithdrawalAmount              = sdkerrors.Register(ModuleName, 113, "The given withdrawal amount is invalid")
	ErrExceedsWithdrawalLimitPerTime        = sdkerrors.Register(ModuleName, 114, "The given amount exceeds the withdrawal limit per time")
	ErrWithdrawalAmountExceedsModuleBalance = sdkerrors.Register(ModuleName, 115, "The given amount to withdraw exceeds module balance")
)
