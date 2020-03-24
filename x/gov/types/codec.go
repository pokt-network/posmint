package types

import (
	"github.com/pokt-network/posmint/codec"
)

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*ACL)(nil), nil)
	cdc.RegisterConcrete(MsgChangeParam{}, "gov/msg_change_param", nil)
	cdc.RegisterConcrete(MsgDAOTransfer{}, "gov/msg_dao_transfer", nil)
	cdc.RegisterInterface((*interface{})(nil), nil)
	cdc.RegisterConcrete(BaseACL{}, "gov/base_acl", nil)
	cdc.RegisterConcrete(NonMapACL{}, "gov/non_map_acl", nil)
	cdc.RegisterConcrete(Upgrade{}, "gov/upgrade", nil)
	cdc.RegisterConcrete(MsgUpgrade{}, "gov/msg_upgrade", nil)
}
