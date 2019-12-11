package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"os"
	"testing"
)

const (
	routeMsgCounter  = "msgCounter"
	routeMsgCounter2 = "msgCounter2"
)

var (
	capKey1      = sdk.NewKVStoreKey("key1")
	capKey2      = sdk.NewKVStoreKey("key2")
	ModuleBasics = module.NewBasicManager(
		bank.AppModuleBasic{},
	)
)

type args struct {
	ctx     sdk.Context
	amount  sdk.Int
	address sdk.ValAddress
}

func TestSetValidatorAward(t *testing.T) {
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	minGasPrices := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 5000)}
	options := baseapp.SetMinGasPrices(minGasPrices.String())
	cdc := MakeCodec()
	bapp := getNewApp(cdc, options)

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		app           *baseapp.BaseApp
		keeper        Keeper
	}{
		{
			name:          "can set Value",
			expectedCoins: sdk.NewInt(1),
			app:           bapp,
			args:          args{ctx: getNewContext(bapp), amount: sdk.NewInt(int64(1)), address: validatorAddress},
			keeper:        Keeper{cdc: cdc},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Fail()
		})
	}
}

func defaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}

func testTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx txTest
		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("").TraceSDK(err.Error())
		}
		return tx, nil
	}
}
func getNewApp(cdc *codec.Codec, options ...func(*baseapp.BaseApp)) *baseapp.BaseApp {
	logger := defaultLogger()
	db := dbm.NewMemDB()
	registerTestCodec(cdc)

	app := baseapp.NewBaseApp("test-app", logger, db, testTxDecoder(cdc), options...)
	app.MountStores(capKey1, capKey2)
	app.InitChain(abci.RequestInitChain{})

	return app
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func getNewContext(bapp *baseapp.BaseApp) sdk.Context {
	header := abci.Header{Height: 0}
	newContext := bapp.NewContext(false, header)
	return newContext
}

func registerTestCodec(cdc *codec.Codec) {
	// register test types
	cdc.RegisterConcrete(&txTest{}, "posmint/baseapp/txTest", nil)
	cdc.RegisterConcrete(&msgCounter{}, "posmint/baseapp/msgCounter", nil)
	cdc.RegisterConcrete(&msgCounter2{}, "posmint/baseapp/msgCounter2", nil)
	cdc.RegisterConcrete(&msgNoRoute{}, "posmint/baseapp/msgNoRoute", nil)
}

type msgCounter struct {
	Counter       int64
	FailOnHandler bool
}

// Implements Msg
func (msg msgCounter) Route() string                { return routeMsgCounter }
func (msg msgCounter) Type() string                 { return "counter1" }
func (msg msgCounter) GetSignBytes() []byte         { return nil }
func (msg msgCounter) GetSigners() []sdk.AccAddress { return nil }
func (msg msgCounter) ValidateBasic() sdk.Error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdk.ErrInvalidSequence("counter should be a non-negative integer.")
}

type msgCounter2 struct {
	Counter int64
}

// Implements Msg
func (msg msgCounter2) Route() string                { return routeMsgCounter2 }
func (msg msgCounter2) Type() string                 { return "counter2" }
func (msg msgCounter2) GetSignBytes() []byte         { return nil }
func (msg msgCounter2) GetSigners() []sdk.AccAddress { return nil }
func (msg msgCounter2) ValidateBasic() sdk.Error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdk.ErrInvalidSequence("counter should be a non-negative integer.")
}

type msgNoRoute struct {
	msgCounter
}

func (tx msgNoRoute) Route() string { return "noroute" }

// a msg we dont know how to decode
type msgNoDecode struct {
	msgCounter
}

func (tx msgNoDecode) Route() string { return routeMsgCounter }

type txTest struct {
	Msgs       []sdk.Msg
	Counter    int64
	FailOnAnte bool
}

func (tx *txTest) setFailOnAnte(fail bool) {
	tx.FailOnAnte = fail
}

func (tx *txTest) setFailOnHandler(fail bool) {
	for i, msg := range tx.Msgs {
		tx.Msgs[i] = msgCounter{msg.(msgCounter).Counter, fail}
	}
}

// Implements Tx
func (tx txTest) GetMsgs() []sdk.Msg       { return tx.Msgs }
func (tx txTest) ValidateBasic() sdk.Error { return nil }
