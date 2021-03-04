package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is the name of the coinswap module
const RouterKey = ModuleName

type MsgExchange struct {
	From      Denom
	To        Denom
	Amount    sdk.Coin       // amount of coins to change
	Requester sdk.AccAddress // address output coin is being sent to
}

// Route returns the route of MsgExchange - "coinswap" (sdk.Msg interface).
func (msg MsgExchange) Route() string { return RouterKey }

// Type returns the message type of MsgExchange (sdk.Msg interface).
func (msg MsgExchange) Type() string { return "exchange" }

// ValidateBasic checks whether the given MsgExchange instance (sdk.Msg interface).
func (msg MsgExchange) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Requester); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "requester: %s", msg.Requester)
	}
	if ok := msg.Amount.IsValid(); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount)
	}
	if ok := msg.Amount.IsPositive(); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount)
	}
	if ok := ValidExchangeDenom(msg.From, msg.To); !ok {
		return sdkerrors.Wrapf(ErrInvalidExchangeDenom, "denominations: %s:%s", msg.From, msg.To)
	}
	if !msg.From.Equal(Denom(msg.Amount.Denom)) {
		return sdkerrors.Wrapf(ErrExchangeDenomMissmatch, "denominations: %s:%s", msg.From, msg.To)
	}
	return nil
}

func (msg MsgExchange) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgExchange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Requester}
}
