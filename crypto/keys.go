package crypto

import (
	"encoding/hex"
	"errors"

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
	PrivKey() crypto.PrivKey
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
	return nil, errors.New("unsupported public key type")
}

func PubKeyToPublicKey(key crypto.PubKey) (PublicKey, error) {
	switch key.(type) {
	case secp256k1.PubKeySecp256k1:
		return Secp256k1PublicKey{}.PubKeyToPublicKey(key), nil
	case ed25519.PubKeyEd25519:
		return Ed25519PublicKey(key.(ed25519.PubKeyEd25519)), nil
	case Ed25519PublicKey:
		return key.(Ed25519PublicKey), nil
	case Secp256k1PublicKey:
		return key.(Secp256k1PublicKey), nil
	default:
		return nil, errors.New("error converting pubkey to public key -> unsupported public key type")
	}
}

func PrivKeyToPrivateKey(key crypto.PrivKey) (PrivateKey, error) {
	switch key.(type) {
	case secp256k1.PrivKeySecp256k1:
		return Secp256k1PrivateKey{}.PrivKeyToPrivateKey(key), nil
	case ed25519.PrivKeyEd25519:
		return Ed25519PrivateKey{}.PrivKeyToPrivateKey(key), nil
	case Secp256k1PrivateKey:
		return key.(Secp256k1PrivateKey), nil
	case Ed25519PrivateKey:
		return key.(Ed25519PrivateKey), nil
	default:
		return nil, errors.New("error converting privkey to private key -> unsupported private key type")
	}
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

func GenerateEd25519PrivKey() PrivateKey {
	return Ed25519PrivateKey{}.GenPrivateKey()
}
