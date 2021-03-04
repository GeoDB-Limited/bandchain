package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

const (
	geo        = "geo"
	odin       = "odin"
	geoSupply  = 100
	odinSupply = 10
)

type testSupplyKeeper struct {
	testSupply
}

func (k testSupplyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (k testSupplyKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

type testSupply struct {
	total sdk.Coins
}

func (s testSupply) GetTotal() sdk.Coins {
	return s.total
}
func (s testSupply) SetTotal(total sdk.Coins) exported.SupplyI {
	s.total = total
	return s
}

func (s testSupply) Inflate(amount sdk.Coins) exported.SupplyI {
	return s
}
func (s testSupply) Deflate(amount sdk.Coins) exported.SupplyI {
	return s
}

func (s testSupply) String() string {
	return ""
}

func (s testSupply) ValidateBasic() error {
	return nil
}

func (k testSupplyKeeper) GetSupply(ctx sdk.Context) exported.SupplyI {
	return k.testSupply
}

func (k testSupplyKeeper) GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI {
	return nil
}

func (k testSupplyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (k testSupplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return nil
}

type testDistrKeeper struct {
	feePool distr.FeePool
}

func (k testDistrKeeper) GetFeePool(ctx sdk.Context) (feePool distr.FeePool) {
	return k.feePool
}

func (k testDistrKeeper) SetFeePool(ctx sdk.Context, feePool distr.FeePool) {
	k.feePool = feePool
}

func TestKeeper_ExchangeDenom(t *testing.T) {
	cdc := codec.New()
	key := types.NewKVStoreKey("test")
	k := NewKeeper(
		cdc,
		key,
		params.Subspace{},
		&testSupplyKeeper{
			testSupply{
				total: sdk.NewCoins(sdk.NewInt64Coin(geo, geoSupply), sdk.NewInt64Coin(odin, odinSupply)),
			},
		},
		&testDistrKeeper{
			feePool: distr.FeePool{
				CommunityPool: sdk.NewDecCoins(sdk.NewInt64DecCoin(geo, geoSupply), sdk.NewInt64DecCoin(odin, odinSupply)),
			},
		},
	)
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, types.StoreTypeIAVL, db)
	cms.LoadLatestVersion()

	addr, _ := types.AccAddressFromBech32("odin12983g7jhxyynse2jmnjy54ukjene837wcncysg")
	ctx := types.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	err := k.ExchangeDenom(ctx, geo, odin, sdk.NewInt64Coin(geo, 10), addr)

	assert.NoError(t, err, "exchange denom failed")
}
