package keeper

// DONTCOVER

import (
	"github.com/pokt-network/posmint/x/gov/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	authTypes "github.com/pokt-network/posmint/x/auth/types"
	"github.com/pokt-network/posmint/x/bank/internal/types"
	"github.com/pokt-network/posmint/x/gov"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	k   Keeper
	ak  auth.AccountKeeper
	pk  keeper.Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.ParamsKey
	tkeyParams := sdk.ParamsTKey

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.Address([]byte("moduleAcc")).String()] = true
	akSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	ak := auth.NewAccountKeeper(
		cdc, authCapKey, akSubspace, auth.ProtoBaseAccount,
	)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak.SetParams(ctx, authTypes.DefaultParams())
	bkSubspace := sdk.NewSubspace(types.DefaultParamspace)
	bankKeeper := NewBaseKeeper(ak, bkSubspace, types.DefaultCodespace, blacklistedAddrs)
	bankKeeper.SetSendEnabled(ctx, true)
	pk := keeper.NewKeeper(cdc, keyParams, tkeyParams, gov.DefaultCodespace, nil, akSubspace, bkSubspace)
	return testInput{cdc: cdc, ctx: ctx, k: bankKeeper, ak: ak, pk: pk}
}
