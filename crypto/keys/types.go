package keys

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/pokt-network/posmint/types"
)

// SigningAlgo defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type SigningAlgo string

const (
	// Secp256k1 uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1 = SigningAlgo("secp256k1")
	// Ed25519 represents the Ed25519 signature system.
	// It is currently not supported for end-user keys (wallets/ledgers).
	Ed25519 = SigningAlgo("ed25519")
)

// Keybase exposes operations on a generic keystore
type Keybase interface {
	// CRUD on the keystore
	List() ([]KeyPair, error)
	Get(address types.AccAddress) (KeyPair, error)
	Delete(address types.AccAddress, passphrase string, skipPass bool) error
	Update(address types.AccAddress, oldpass string, newpass string) error

	// Sign some bytes, looking up the private key to use
	Sign(address types.AccAddress, passphrase string, msg []byte) ([]byte, crypto.PubKey, error)

	// CreateMnemonic creates a new mnemonic and persists the keypair to disk encrypted using passwd
	CreateMnemonic(bip39Passwd string, passwd string) (kp KeyPair, mnemonic string, err error)

	// Derive computes a BIP39 seed from th mnemonic and bip39Passwd.
	// Encrypt the key to disk using encryptPasswd.
	DeriveFromMnemonic(mnemonic, bip39Passwd, encryptPasswd string) (KeyPair, error)

	// ImportPrivKey using Armored private key string. Decrypts armor with decryptPassphrase, and stores locally using encryptPassphrase
	ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) error

	// ExportPrivKeyArmor using Armored private key string. Decrypts armor with decryptPassphrase, and encrypts result armor using the encryptPassphrase
	ExportPrivKeyEncryptedArmor(address types.AccAddress, decryptPassphrase, encryptPassphrase string) (armor string, err error)

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
	return cdc.MustMarshalBinaryLengthPrefixed(i)
}

// decoding info
func readKeyPair(bz []byte) (kp KeyPair, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(bz, &kp)
	return
}
