package keys

import (
	"github.com/pokt-network/posmint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/pokt-network/posmint/codec"
)

var cdc *codec.Codec

func init() {
	cdc = codec.New()
	cryptoAmino.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	cdc.RegisterConcrete(KeyPair{}, "crypto/keys/keypair", nil)
	cdc.Seal()
}
