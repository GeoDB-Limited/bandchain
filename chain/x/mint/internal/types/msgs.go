package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var _ sdk.Msg = &MsgMintTokens{}

// MsgMintTokens defines a msg to mint some amount for depositor account
type MsgMintTokens struct {
	Amount    sdk.Int        `json:"amount" yaml:"amount"`
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"`
	Sender    sdk.AccAddress `json:"sender" yaml:"sender"`
}

// NewMsgMintTokens returns a new MsgMintTokens
func NewMsgMintTokens(Amount sdk.Int, Depositor sdk.AccAddress, Sender sdk.AccAddress) MsgMintTokens {
	return MsgMintTokens{
		Amount:    Amount,
		Depositor: Depositor,
		Sender:    Sender,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgMintTokens) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgMintTokens) Type() string {
	return "mint_tokens"
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgMintTokens) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if msg.Depositor.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "depositor: %s", msg.Sender)
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidAmountToMint, "amount: %s", msg.Amount)
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMintTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgMintTokens) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
