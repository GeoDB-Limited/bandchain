package mint

import (
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgWithdrawCoinsToAccFromTreasury:
			return handleWithdrawCoinsToAccFromTreasury(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleWithdrawCoinsToAccFromTreasury handles MsgWithdrawCoins
func handleWithdrawCoinsToAccFromTreasury(
	ctx sdk.Context,
	k Keeper,
	msg MsgWithdrawCoinsToAccFromTreasury,
) (*sdk.Result, error) {
	if !k.IsEligibleAccount(ctx, msg.Sender) {
		return nil, sdkerrors.Wrapf(types.ErrAccountIsNotEligible, "account: %s", msg.Sender)
	}
	if k.LimitExceeded(ctx, msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrExceedsWithdrawalLimitPerTime, "amount: %s", msg.Amount)
	}

	err := k.WithdrawCoinsToAccFromTreasury(ctx, msg.Receiver, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to mint %s coins to account %s", msg.Amount, msg.Receiver)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeWithdrawal,
		sdk.NewAttribute(types.AttributeKeyWithdrawalAmount, msg.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver.String()),
		sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
