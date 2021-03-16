package testing

import (
	"github.com/GeoDB-Limited/odincore/chain/x/common/testapp"
	odinmint "github.com/GeoDB-Limited/odincore/chain/x/mint"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrappedSupplyKeeper_MintCoinsCommPoolSuccess(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(false, true)
	err := app.WrappedSupplyKeeper.MintCoins(ctx, odinmint.ModuleName, testapp.Coins1000000odin)
	require.NoError(t, err)
	events := ctx.EventManager().Events()
	require.Equal(t, events, sdk.Events{
		sdk.NewEvent(
			bank.EventTypeTransfer,
			sdk.NewAttribute(bank.AttributeKeyRecipient, app.SupplyKeeper.GetModuleAddress(odinmint.ModuleName).String()),
			sdk.NewAttribute(bank.AttributeKeySender, app.SupplyKeeper.GetModuleAddress(distr.ModuleName).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, testapp.Coins1000000odin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(bank.AttributeKeySender, app.SupplyKeeper.GetModuleAddress(distr.ModuleName).String()),
		),
	})
	commPool := app.DistrKeeper.GetFeePool(ctx)
	res, _ := commPool.CommunityPool.TruncateDecimal()
	require.Equal(t, testapp.DefaultCommunityPool.Sub(testapp.Coins1000000odin).String(), res.String())
}
