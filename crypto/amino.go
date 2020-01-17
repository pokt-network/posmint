package crypto

import (
	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
	cryptoAmino.RegisterAmino(cdc)
}

// RegisterAmino registers all go-crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*PublicKey)(nil), nil)
	cdc.RegisterInterface((*PrivateKey)(nil), nil)
	cdc.RegisterConcrete(Ed25519PublicKey{}, "crypto/ed25519_public_key", nil)
	cdc.RegisterConcrete(Ed25519PrivateKey{}, "crypto/ed25519_private_key", nil)
	cdc.RegisterConcrete(Secp256k1PublicKey{}, "crypto/secp256k1_public_key", nil)
	cdc.RegisterConcrete(Secp256k1PrivateKey{}, "crypto/secp256k1_private_key", nil)
}
