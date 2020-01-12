package crypto

import (
	ed255192 "crypto/ed25519"
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type PublicKey ed25519.PubKeyEd25519
type PrivateKey ed25519.PrivKeyEd25519

const (
	PrivKeySize   = ed255192.PrivateKeySize
	PubKeySize    = ed25519.PubKeyEd25519Size
	SignatureSize = ed25519.SignatureSize
)

func NewPublicKey(hexString string) (PublicKey, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return PublicKey{}, err
	}
	var bz [PubKeySize]byte
	copy(bz[:], b)
	pubkey := ed25519.PubKeyEd25519(bz)
	pk := PublicKey(pubkey)
	return pk, nil
}

func (pub PublicKey) AminoBytes() []byte {
	return ed25519.PubKeyEd25519(pub).Bytes()
}

func (pub PublicKey) Bytes() []byte {
	pkBytes := [PubKeySize]byte(pub)
	return pkBytes[:]
}

func (pub PublicKey) String() string {
	return hex.EncodeToString(pub.Bytes())
}

func (pub PublicKey) Address() sdk.ValAddress {
	return sdk.ValAddress(ed25519.PubKeyEd25519(pub).Address())
}

func (pub PublicKey) VerifySignature(msg []byte, sig []byte) bool {
	return ed25519.PubKeyEd25519(pub).VerifyBytes(msg, sig)
}

func (priv PrivateKey) Bytes() []byte {
	pkBytes := [PrivKeySize]byte(priv)
	return pkBytes[:]
}

func (priv PrivateKey) AminoBytes() []byte {
	return ed25519.PrivKeyEd25519(priv).Bytes()
}

func (priv PrivateKey) String() string {
	return hex.EncodeToString(ed25519.PrivKeyEd25519(priv).Bytes())
}

func (priv PrivateKey) Public() PublicKey {
	return PublicKey(ed25519.PrivKeyEd25519(priv).PubKey().(ed25519.PubKeyEd25519))
}

func (priv PrivateKey) Sign(msg []byte) ([]byte, error) {
	return ed25519.PrivKeyEd25519(priv).Sign(msg)
}
