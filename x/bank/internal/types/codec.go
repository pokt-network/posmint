package types

import (
	"github.com/pokt-network/posmint/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "posmint/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "posmint/MsgMultiSend", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
