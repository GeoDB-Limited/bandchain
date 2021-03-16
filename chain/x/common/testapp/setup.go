package testapp

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/GeoDB-Limited/odincore/chain/app"
	"github.com/GeoDB-Limited/odincore/chain/pkg/filecache"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle"
	me "github.com/GeoDB-Limited/odincore/chain/x/oracle/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// Account is a data structure to store key of test account.
type Account struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

// nolint
var (
	OraclePoolProvider Account
	FeePoolProvider    Account
	Owner              Account
	Alice              Account
	Bob                Account
	Carol              Account
	Validator1         Account
	Validator2         Account
	Validator3         Account
	DataSources        []types.DataSource
	OracleScripts      []types.OracleScript
)

// nolint
var (
	Coin1geo                 = sdk.NewInt64Coin("geo", 1)
	Coin10odin               = sdk.NewInt64Coin("odin", 10)
	Coin100000000geo         = sdk.NewInt64Coin("geo", 1000000)
	Coins1000000odin         = sdk.NewCoins(sdk.NewInt64Coin("odin", 1000000))
	Coins99999999odin        = sdk.NewCoins(sdk.NewInt64Coin("odin", 99999999))
	Coin100000000odin        = sdk.NewInt64Coin("odin", 100000000)
	Coins100000000odin       = sdk.NewCoins(Coin100000000odin)
	DefaultDataProvidersPool = sdk.NewCoins(Coin100000000odin)
	DefaultCommunityPool     = sdk.NewCoins(Coin100000000geo, Coin100000000odin)
)

func init() {
	bandapp.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	Owner = createArbitraryAccount(r)
	Alice = createArbitraryAccount(r)
	Bob = createArbitraryAccount(r)
	Carol = createArbitraryAccount(r)
	Validator1 = createArbitraryAccount(r)
	Validator2 = createArbitraryAccount(r)
	Validator3 = createArbitraryAccount(r)
	OraclePoolProvider = createArbitraryAccount(r)
	FeePoolProvider = createArbitraryAccount(r)
}

func createArbitraryAccount(r *rand.Rand) Account {
	privkeySeed := make([]byte, 12)
	r.Read(privkeySeed)
	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	return Account{
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}

func getGenesisDataSources() []types.DataSource {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	fc := filecache.New(dir)
	DataSources = []types.DataSource{{}} // 0th index should be ignored
	for idx := 0; idx < 5; idx++ {
		idxStr := fmt.Sprintf("%d", idx+1)
		hash := fc.AddFile([]byte("code" + idxStr))
		DataSources = append(DataSources, types.NewDataSource(
			Owner.Address, "name"+idxStr, "desc"+idxStr, hash,
		))
	}
	return DataSources[1:]
}

func getGenesisOracleScripts() []types.OracleScript {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	fc := filecache.New(dir)
	OracleScripts = []types.OracleScript{{}} // 0th index should be ignored
	wasms := [][]byte{
		Wasm1, Wasm2, Wasm3, Wasm4, Wasm56(10), Wasm56(10000000), Wasm78(10), Wasm78(2000), Wasm9,
	}
	for idx := 0; idx < len(wasms); idx++ {
		idxStr := fmt.Sprintf("%d", idx+1)
		hash := fc.AddFile(compile(wasms[idx]))
		OracleScripts = append(OracleScripts, types.NewOracleScript(
			Owner.Address, "name"+idxStr, "desc"+idxStr, hash, "schema"+idxStr, "url"+idxStr,
		))
	}
	return OracleScripts[1:]
}

func createValidatorTx(chainID string, acc Account, moniker string, selfDelegation sdk.Coin) authtypes.StdTx {
	msg := staking.NewMsgCreateValidator(
		acc.ValAddress, acc.PubKey, selfDelegation,
		staking.NewDescription(moniker, "", "", "", ""),
		staking.NewCommissionRates(sdk.MustNewDecFromStr("0.125"), sdk.MustNewDecFromStr("0.3"), sdk.MustNewDecFromStr("0.01")),
		sdk.NewInt(1),
	)
	txMsg := authtypes.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: 0,
		Sequence:      0,
		Fee:           auth.NewStdFee(200000, sdk.Coins{}),
		Msgs:          []sdk.Msg{msg},
		Memo:          "",
	}
	sigBytes, err := acc.PrivKey.Sign(txMsg.Bytes())
	if err != nil {
		panic(err)
	}
	sigs := []authtypes.StdSignature{{
		PubKey:    acc.PubKey,
		Signature: sigBytes,
	}}
	return authtypes.NewStdTx([]sdk.Msg{msg}, auth.NewStdFee(200000, sdk.Coins{}), sigs, "")
}

// NewSimApp creates instance of our app using in test.
func NewSimApp(chainID string, logger log.Logger) *bandapp.BandApp {
	// Set HomeFlag to a temp folder for simulation run.
	dir, err := ioutil.TempDir("", "bandd")
	if err != nil {
		panic(err)
	}
	viper.Set(cli.HomeFlag, dir)
	db := dbm.NewMemDB()
	app := bandapp.NewBandApp(logger, db, nil, true, 0, map[int64]bool{}, "", false)
	genesis := bandapp.NewDefaultGenesisState()
	// Fund seed accounts and validators with 1000000odin and 100000000odin initially.
	authGenesis := auth.NewGenesisState(auth.DefaultParams(), []authexported.GenesisAccount{
		&auth.BaseAccount{Address: Owner.Address, Coins: Coins1000000odin},
		&auth.BaseAccount{Address: Alice.Address, Coins: Coins1000000odin.Add(Coin100000000geo)},
		&auth.BaseAccount{Address: Bob.Address, Coins: Coins1000000odin},
		&auth.BaseAccount{Address: Carol.Address, Coins: Coins1000000odin},
		&auth.BaseAccount{Address: Validator1.Address, Coins: Coins100000000odin},
		&auth.BaseAccount{Address: Validator2.Address, Coins: Coins100000000odin},
		&auth.BaseAccount{Address: Validator3.Address, Coins: Coins100000000odin},
		&auth.BaseAccount{Address: OraclePoolProvider.Address, Coins: DefaultDataProvidersPool},
		&auth.BaseAccount{Address: FeePoolProvider.Address, Coins: DefaultCommunityPool},
	})
	genesis[auth.ModuleName] = app.Codec().MustMarshalJSON(authGenesis)
	// Add genesis transactions to create 3 validators during chain genesis.
	genutilGenesis := genutil.NewGenesisStateFromStdTx([]authtypes.StdTx{
		createValidatorTx(chainID, Validator1, "validator1", Coins100000000odin[0]),
		createValidatorTx(chainID, Validator2, "validator2", Coins1000000odin[0]),
		createValidatorTx(chainID, Validator3, "validator3", Coins99999999odin[0]),
	})
	genesis[genutil.ModuleName] = app.Codec().MustMarshalJSON(genutilGenesis)
	// Add genesis data sources and oracle scripts
	oracleGenesis := oracle.DefaultGenesisState()

	oracleGenesis.Params.DataRequesterBasicFee = types.CoinProto(Coin10odin)
	oracleGenesis.Params.DataProviderRewardPerByte = types.CoinDecProto(sdk.NewDecCoinFromCoin(Coin1geo))

	oracleGenesis.DataSources = getGenesisDataSources()
	oracleGenesis.OracleScripts = getGenesisOracleScripts()

	genesis[oracle.ModuleName] = app.Codec().MustMarshalJSON(oracleGenesis)
	// Initialize the sim blockchain. We are ready for testing!
	app.InitChain(abci.RequestInitChain{
		ChainId:       chainID,
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), genesis),
	})
	return app
}

// CreateTestInput creates a new test environment for unit tests.
// params[0] - activate;
// params[1] - fund pools;
func CreateTestInput(params ...bool) (*bandapp.BandApp, sdk.Context, me.Keeper) {
	app := NewSimApp("ODINCHAIN", log.NewNopLogger())
	ctx := app.NewContext(false, abci.Header{})
	if len(params) > 0 && params[0] {
		app.OracleKeeper.Activate(ctx, Validator1.ValAddress)
		app.OracleKeeper.Activate(ctx, Validator2.ValAddress)
		app.OracleKeeper.Activate(ctx, Validator3.ValAddress)
	}

	if len(params) > 1 && params[1] {
		app.DistrKeeper.FundCommunityPool(ctx, DefaultCommunityPool, FeePoolProvider.Address)
		app.OracleKeeper.FundOraclePool(ctx, DefaultDataProvidersPool, OraclePoolProvider.Address)

		ctx = app.NewContext(false, abci.Header{})
	}
	return app, ctx, app.OracleKeeper
}
