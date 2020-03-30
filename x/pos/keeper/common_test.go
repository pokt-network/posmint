package keeper

import (
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth/keeper"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/pokt-network/posmint/x/pos/types"
	//"github.com/pokt-network/posmint/x/supply/internal/types"

	sdk "github.com/pokt-network/posmint/types"
)

// nolint: deadcode unused
var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		gov.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, []auth.Account, Keeper) {
	initPower := int64(100000000000)
	nAccs := int64(4)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyPOS := sdk.NewKVStoreKey(types.ModuleName)
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyPOS, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName: nil,
		types.StakedPoolName:  {auth.Burner, auth.Staking, auth.Minter},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)
	akSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	posSubspace := sdk.NewSubspace(DefaultParamspace)
	ak := auth.NewKeeper(cdc, keyAcc, akSubspace, maccPerms)
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)
	keeper := NewKeeper(cdc, keyPOS, ak, posSubspace, DefaultParamspace)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)
	return ctx, accs, keeper
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Ctx, numAccs int, initialCoins sdk.Coins, ak *keeper.Keeper) (accs []auth.Account) {
	var err error
	for i := 0; i < numAccs; i++ {
		privKey := crypto.Ed25519PrivateKey{}.GenPrivateKey()
		pubKey := privKey.PubKey()
		addr := sdk.Address(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins

		acc.PubKey, err = crypto.PubKeyToPublicKey(pubKey)
		if err != nil {
			panic(err)
		}
		ak.SetAccount(ctx, &acc)
	}
	return
}

func addMintedCoinsToModule(t *testing.T, ctx sdk.Ctx, k *Keeper, module string) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(100000000000)))
	mintErr := k.authKeeper.MintCoins(ctx, module, coins.Add(coins))
	if mintErr != nil {
		t.Fail()
	}
}

func sendFromModuleToAccount(t *testing.T, ctx sdk.Ctx, k *Keeper, module string, address sdk.Address, amount sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.authKeeper.SendCoinsFromModuleToAccount(ctx, module, sdk.Address(address), coins)
	if err != nil {
		t.Fail()
	}
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	return pub
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

func getValidator() types.Validator {
	pub := getRandomPubKey()
	return types.Validator{
		Address:      sdk.Address(pub.Address()),
		StakedTokens: sdk.NewInt(100000000000),
		PublicKey:    pub,
		Jailed:       false,
		Status:       sdk.Staked,
	}
}

func getStakedValidator() types.Validator {
	return getValidator()
}

func getUnstakedValidator() types.Validator {
	v := getValidator()
	return v.UpdateStatus(sdk.Unstaked)
}

func getUnstakingValidator() types.Validator {
	v := getValidator()
	return v.UpdateStatus(sdk.Unstaking)
}
