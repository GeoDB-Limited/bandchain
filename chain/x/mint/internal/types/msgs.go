package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var _ sdk.Msg = &MsgMintCoinsToAcc{}

// MsgMintCoinsToAcc defines a msg to mint some amount for depositor account
type MsgMintCoinsToAcc struct {
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"`
	Sender    sdk.AccAddress `json:"sender" yaml:"sender"`
}

// NewMsgMintCoinsToAcc returns a new MsgMintCoinsToAcc
func NewMsgMintCoinsToAcc(amount sdk.Coins, depositor sdk.AccAddress, sender sdk.AccAddress) MsgMintCoinsToAcc {
	return MsgMintCoinsToAcc{
		Amount:    amount,
		Depositor: depositor,
		Sender:    sender,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgMintCoinsToAcc) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgMintCoinsToAcc) Type() string {
	return "mint_coins"
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgMintCoinsToAcc) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if msg.Depositor.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "depositor: %s", msg.Sender)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMintCoinsToAcc) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgMintCoinsToAcc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
