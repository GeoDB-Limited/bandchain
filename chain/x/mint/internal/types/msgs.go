package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawCoinsToAccFromTreasury = "withdraw_coins_from_treasury"

// ensure Msg interface compliance at compile time
var _ sdk.Msg = &MsgMsgWithdrawCoinsToAccFromTreasury{}

// MsgMsgWithdrawCoinsToAccFromTreasury defines a msg to mint some amount for receiver account
type MsgMsgWithdrawCoinsToAccFromTreasury struct {
	Amount   sdk.Coins      `json:"amount" yaml:"amount"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
}

// NewMsgMsgWithdrawCoinsToAccFromTreasury returns a new MsgMsgWithdrawCoinsToAccFromTreasury
func NewMsgMsgWithdrawCoinsToAccFromTreasury(
	amt sdk.Coins,
	receiver sdk.AccAddress,
	sender sdk.AccAddress,
) MsgMsgWithdrawCoinsToAccFromTreasury {
	return MsgMsgWithdrawCoinsToAccFromTreasury{
		Amount:   amt,
		Receiver: receiver,
		Sender:   sender,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgMsgWithdrawCoinsToAccFromTreasury) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgMsgWithdrawCoinsToAccFromTreasury) Type() string {
	return TypeMsgWithdrawCoinsToAccFromTreasury
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgMsgWithdrawCoinsToAccFromTreasury) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if msg.Receiver.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "receiver: %s", msg.Sender)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrapf(ErrInvalidMintAmount, "amount: %s", msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMsgWithdrawCoinsToAccFromTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgMsgWithdrawCoinsToAccFromTreasury) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
