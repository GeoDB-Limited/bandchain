package mint

import (
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgMintTokens:
			return handleMsgMintTokens(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgMintTokens handles MsgMintTokens
func handleMsgMintTokens(ctx sdk.Context, keeper Keeper, msg MsgMintTokens) (*sdk.Result, error) {
	if !keeper.IsEligibleAccount(ctx, msg.Sender) {
		return nil, sdkerrors.Wrapf(types.ErrAccountIsNotEligible, fmt.Sprintf("account: %s", msg.Sender))
	}

	// TODO: mint from module to account

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeMint,
		sdk.NewAttribute(types.AttributeKeyMintAmount, fmt.Sprintf("%d", msg.Amount)),
		sdk.NewAttribute(types.AttributeKeyDepositor, fmt.Sprintf("%d", msg.Depositor)),
		sdk.NewAttribute(types.AttributeKeySender, fmt.Sprintf("%d", msg.Sender)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
