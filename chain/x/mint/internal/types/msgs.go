package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var _ sdk.Msg = &MsgMintCoinToAcc{}

// MsgMintCoinToAcc defines a msg to mint some amount for receiver account
type MsgMintCoinToAcc struct {
	Amount   sdk.Coin       `json:"amount" yaml:"amount"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
}

// NewMsgMintCoinToAcc returns a new MsgMintCoinToAcc
func NewMsgMintCoinToAcc(amt sdk.Coin, receiver sdk.AccAddress, sender sdk.AccAddress) MsgMintCoinToAcc {
	return MsgMintCoinToAcc{
		Amount:   amt,
		Receiver: receiver,
		Sender:   sender,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgMintCoinToAcc) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgMintCoinToAcc) Type() string {
	return "mint_coins"
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgMintCoinToAcc) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if msg.Receiver.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "receiver: %s", msg.Sender)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}
	if !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidMintAmount, "amount: %s", msg.Amount.String())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMintCoinToAcc) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgMintCoinToAcc) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
