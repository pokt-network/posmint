package crypto

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type PublicKey ed25519.PubKeyEd25519
type PrivateKey ed25519.PrivKeyEd25519

const (
	PubKeySize = ed25519.PubKeyEd25519Size
)

func NewPublicKey(hexString string) (*PublicKey, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	if len(b) != PubKeySize {
		return nil, err
	}
	var temp [PubKeySize]byte
	copy(temp[:], b)
	pk := PublicKey(ed25519.PubKeyEd25519(temp))
	return &pk, nil
}

func (pub PublicKey) Bytes() []byte {
	return ed25519.PubKeyEd25519(pub).Bytes()
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
