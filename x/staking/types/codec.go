package types

import (
	"github.com/pokt-network/posmint/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "posmint/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "posmint/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "posmint/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "posmint/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "posmint/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
