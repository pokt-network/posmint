package crypto

import (
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type PublicKey interface {
	NewPublicKey([]byte) (PublicKey, error)
	PubKey() crypto.PubKey
	Bytes() []byte
	RawBytes() []byte
	String() string
	RawString() string
	Address() crypto.Address
	Equals(other crypto.PubKey) bool
	VerifyBytes(msg []byte, sig []byte) bool
	PubKeyToPublicKey(crypto.PubKey) PublicKey
	Size() int
}

type PrivateKey interface {
	Bytes() []byte
	RawBytes() []byte
	String() string
	RawString() string
	PubKey() crypto.PubKey
	Equals(other crypto.PrivKey) bool
	PublicKey() PublicKey
	Sign(msg []byte) ([]byte, error)
	PrivKeyToPrivateKey(crypto.PrivKey) PrivateKey
	GenPrivateKey() PrivateKey
	Size() int
}

func NewPublicKey(hexString string) (PublicKey, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return NewPublicKeyBz(b)
}

func NewPublicKeyBz(b []byte) (PublicKey, error) {
	switch len(b) {
	case Ed25519PubKeySize:
		return Ed25519PublicKey{}.NewPublicKey(b)
	case Secp256k1PublicKeySize:
		return Secp256k1PublicKey{}.NewPublicKey(b)
	}
	panic("unsupported public key type")
}

func PubKeyToPublicKey(key crypto.PubKey) PublicKey {
	switch key.(type) {
	case secp256k1.PubKeySecp256k1:
		return Secp256k1PublicKey{}.PubKeyToPublicKey(key)
	case ed25519.PubKeyEd25519:
		return Ed25519PublicKey{}.PubKeyToPublicKey(key)
	}
	panic("unsupported private key type")
}

func PrivKeyToPrivateKey(key crypto.PrivKey) PrivateKey {
	switch key.(type) {
	case secp256k1.PrivKeySecp256k1:
		return Secp256k1PrivateKey{}.PrivKeyToPrivateKey(key)
	case ed25519.PrivKeyEd25519:
		return Ed25519PrivateKey{}.PrivKeyToPrivateKey(key)
	}
	panic("unsupported private key type")
}

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivateKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

func PubKeyFromBytes(pubKeyBytes []byte) (pubKey PublicKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}

func GenerateSecp256k1PrivKey() PrivateKey {
	return Secp256k1PrivateKey{}.GenPrivateKey()
}

func GenereateEd25519PrivKey() PrivateKey {
	return Ed25519PrivateKey{}.GenPrivateKey()
}
