package types

import (
	"github.com/pokt-network/posmint/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "posmint/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorCommission{}, "posmint/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "posmint/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(CommunityPoolSpendProposal{}, "posmint/CommunityPoolSpendProposal", nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
