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
		case MsgMintCoinsToAcc:
			return handleMintCoinsToAcc(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMintCoinsToAcc handles MsgMintCoinsToAcc
func handleMintCoinsToAcc(ctx sdk.Context, k Keeper, msg MsgMintCoinsToAcc) (*sdk.Result, error) {
	if !k.IsEligibleAccount(ctx, msg.Sender) {
		return nil, sdkerrors.Wrapf(types.ErrAccountIsNotEligible, "account: %s", msg.Sender)
	}

	err := k.MintCoinsToAcc(ctx, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to mint %s coins to account %s", msg.Amount, msg.Depositor)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeMint,
		sdk.NewAttribute(types.AttributeKeyMintAmount, msg.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDepositor, msg.Depositor.String()),
		sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
