package testUtil

import (
	"github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)


func GetNewApp(capKey *sdk.KVStoreKey, cdc *codec.Codec, options ...func(*baseapp.BaseApp)) *baseapp.BaseApp {
	logger := defaultLogger()
	db := dbm.NewMemDB()
	registerTestCodec(cdc)

	app := baseapp.NewBaseApp("test-app", logger, db, TestTxDecoder(cdc), options...)
	//app.LoadLatestVersion(capKey1)

	// make a cap key and mount the store
	app.MountStores(capKey)
	err := app.LoadLatestVersion(capKey) // needed to make stores non-nil

	if err != nil {
		panic(err)
	}
	app.InitChain(abci.RequestInitChain{})

	return app
}
