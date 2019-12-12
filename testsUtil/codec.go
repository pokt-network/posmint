package testUtil

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/bank"
)

var ModuleBasics = module.NewBasicManager(
	bank.AppModule{},
	)

func registerTestCodec(cdc *codec.Codec) {
	// register test types
	cdc.RegisterConcrete(&txTest{}, "posmint/baseapp/txTest", nil)
	cdc.RegisterConcrete(&msgCounter{}, "posmint/baseapp/msgCounter", nil)
	cdc.RegisterConcrete(&msgCounter2{}, "posmint/baseapp/msgCounter2", nil)
	cdc.RegisterConcrete(&msgNoRoute{}, "posmint/baseapp/msgNoRoute", nil)
}
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
