package keys

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/pokt-network/posmint/types"
)

// SigningAlgo defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type SigningAlgo string

// const (
// 	// Ed25519 represents the Ed25519 signature system.
// 	// It is currently not supported for end-user keys (wallets/ledgers).
// 	Ed25519 = SigningAlgo("ed25519")
// )

// Keybase exposes operations on a generic keystore
// Keybase only supports Ed25519 key pairs
type Keybase interface {
	// CRUD on the keystore
	List() ([]KeyPair, error)
	Get(address types.AccAddress) (KeyPair, error)
	Delete(address types.AccAddress, passphrase string) error
	Update(address types.AccAddress, oldpass string, newpass string) error

	// Sign some bytes, looking up the private key to use
	Sign(address types.AccAddress, passphrase string, msg []byte) ([]byte, crypto.PubKey, error)

	// Create a new KeyPair and encrypt it to disk using encryptPassphrase
	Create(encryptPassphrase string) (KeyPair, error)

	// ImportPrivKey using Armored private key string. Decrypts armor with decryptPassphrase, and stores locally using encryptPassphrase
	ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) (KeyPair, error)

	// ExportPrivKeyArmor using Armored private key string. Decrypts armor with decryptPassphrase, and encrypts result armor using the encryptPassphrase
	ExportPrivKeyEncryptedArmor(address types.AccAddress, decryptPassphrase, encryptPassphrase string) (armor string, err error)

	// ImportPrivateKeyObject using the raw unencrypted privateKey string and encrypts it to disk using encryptPassphrase
	ImportPrivateKeyObject(privateKey [64]byte, encryptPassphrase string) (KeyPair, error)

	// ExportPrivateKeyObject exports raw PrivKey object.
	ExportPrivateKeyObject(address types.AccAddress, passphrase string) (crypto.PrivKey, error)

	// CloseDB closes the database.
	CloseDB()
}

// KeyPair is the public information about a locally stored key
type KeyPair struct {
	PubKey       crypto.PubKey `json:"pubkey"`
	PrivKeyArmor string        `json:"privkey.armor"`
}

// NewKeyPair with the given public key and priv armor key
func NewKeyPair(pub crypto.PubKey, privArmor string) KeyPair {
	return KeyPair{
		PubKey:       pub,
		PrivKeyArmor: privArmor,
	}
}

// GetAddress for the given KeyPair
func (kp KeyPair) GetAddress() types.AccAddress {
	return kp.PubKey.Address().Bytes()
}

// encoding info
func writeKeyPair(kp KeyPair) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(kp)
}

// decoding info
func readKeyPair(bz []byte) (kp KeyPair, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(bz, &kp)
	return
}
