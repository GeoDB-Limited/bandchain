package yoda

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/app"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestGetSignBytesVerificationMessage(t *testing.T) {
	app.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	validator, _ := sdk.ValAddressFromBech32("odinvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgj3avjp9")
	vmsg := NewVerificationMessage("odinchain", validator, types.RequestID(1), types.ExternalID(1))
	expected := []byte("{\"chain_id\":\"odinchain\",\"external_id\":\"1\",\"request_id\":\"1\",\"validator\":\"odinvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgj3avjp9\"}")
	require.Equal(t, expected, vmsg.GetSignBytes())
}
