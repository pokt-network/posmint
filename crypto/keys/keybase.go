package keys

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pokt-network/posmint/crypto/keys/mintkey"
	"github.com/pokt-network/posmint/types"

	"github.com/cosmos/go-bip39"

	"github.com/tendermint/crypto/ed25519"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmed "github.com/tendermint/tendermint/crypto/ed25519"

	dbm "github.com/tendermint/tm-db"
)

var _ Keybase = dbKeybase{}

// Language is a language to create the BIP 39 mnemonic in.
// Currently, only english is supported though.
// Find a list of all supported languages in the BIP 39 spec (word lists).
type Language int

//noinspection ALL
const (
	// English is the default language to create a mnemonic.
	// It is the only supported language by this package.
	English Language = iota + 1
	// Japanese is currently not supported.
	Japanese
	// Korean is currently not supported.
	Korean
	// Spanish is currently not supported.
	Spanish
	// ChineseSimplified is currently not supported.
	ChineseSimplified
	// ChineseTraditional is currently not supported.
	ChineseTraditional
	// French is currently not supported.
	French
	// Italian is currently not supported.
	Italian
	//addressSuffix = "address"
	//infoSuffix    = "info"
)

const (
	// used for deriving seed from mnemonic
	DefaultBIP39Passphrase = ""

	// bits of entropy to draw when creating a mnemonic
	defaultEntropySize = 256
)

var (
	// ErrUnsupportedSigningAlgo is raised when the caller tries to use a
	// different signing scheme than ed25519.
	ErrUnsupportedSigningAlgo = errors.New("unsupported signing algo: only ed25519 is supported")

	// ErrUnsupportedLanguage is raised when the caller tries to use a
	// different language than english for creating a mnemonic sentence.
	ErrUnsupportedLanguage = errors.New("unsupported language: only english is supported")
)

// dbKeybase combines encryption and storage implementation to provide
// a full-featured key manager
type dbKeybase struct {
	db dbm.DB
}

// newDbKeybase creates a new keybase instance using the passed DB for reading and writing keys.
func newDbKeybase(db dbm.DB) Keybase {
	return dbKeybase{
		db: db,
	}
}

// NewInMemory creates a transient keybase on top of in-memory storage
// instance useful for testing purposes and on-the-fly key generation.
func NewInMemory() Keybase { return dbKeybase{dbm.NewMemDB()} }

// List returns the keys from storage in alphabetical order.
func (kb dbKeybase) List() ([]KeyPair, error) {
	var res []KeyPair
	iter := kb.db.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		kp, err := readKeyPair(iter.Value())
		if err != nil {
			return nil, err
		}
		res = append(res, kp)
	}
	return res, nil
}

// Get returns the public information about one key.
func (kb dbKeybase) Get(address types.AccAddress) (KeyPair, error) {
	ik := kb.db.Get(addrKey(address))
	if len(ik) == 0 {
		return KeyPair{}, fmt.Errorf("key with address %s not found", address)
	}
	bs := kb.db.Get(ik)
	return readKeyPair(bs)
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security).
// It returns an error if the key doesn't exist or
// passphrases don't match.
// Passphrase is ignored when deleting references to
// offline and Ledger / HW wallet keys.
func (kb dbKeybase) Delete(address types.AccAddress, passphrase string, skipPass bool) error {
	// verify we have the proper password before deleting
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}
	if !skipPass {
		if _, err = mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase); err != nil {
			return err
		}
	}
	kb.db.DeleteSync(addrKey(kp.GetAddress()))
	return nil
}

// Update changes the passphrase with which an already stored key is
// encrypted.
//
// oldpass must be the current passphrase used for encryption,
// getNewpass is a function to get the passphrase to permanently replace
// the current passphrase
func (kb dbKeybase) Update(address types.AccAddress, oldpass string, newpass string) error {
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}

	privKey, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, oldpass)
	if err != nil {
		return err
	}

	kb.writeLocalKeyPair(privKey, newpass)
	return nil
}

// Sign signs the msg with the named key.
// It returns an error if the key doesn't exist or the decryption fails.
func (kb dbKeybase) Sign(address types.AccAddress, passphrase string, msg []byte) ([]byte, tmcrypto.PubKey, error) {
	kp, err := kb.Get(address)
	if err != nil {
		return nil, nil, err
	}

	if kp.PrivKeyArmor == "" {
		err = fmt.Errorf("private key not available")
		return nil, nil, err
	}

	priv, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	if err != nil {
		return nil, nil, err
	}

	sig, err := priv.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	pub := priv.PubKey()
	return sig, pub, nil
}

// CreateMnemonic generates a new key and persists it to storage, encrypted
// using the provided password.
// It returns the generated mnemonic and the key Info.
// It returns an error if it fails to
// generate a key for the given algo type
func (kb dbKeybase) CreateMnemonic(bip39Passwd string, passwd string) (kp KeyPair, mnemonic string, err error) {
	// default number of words (24):
	// this generates a mnemonic directly from the number of words by reading system entropy.
	entropy, err := bip39.NewEntropy(defaultEntropySize)
	if err != nil {
		return
	}
	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return
	}

	seed := bip39.NewSeed(mnemonic, bip39Passwd)
	res := ed25519.NewKeyFromSeed(seed)
	var pk [64]byte
	copy(pk[:], res)
	kp = kb.writeLocalKeyPair(tmed.PrivKeyEd25519(pk), passwd)
	return
}

// DeriveFromMnemonic a KeyPair using bip39Passwd, and encrypts using encryptPasswd
func (kb dbKeybase) DeriveFromMnemonic(mnemonic, bip39Passwd, encryptPasswd string) (KeyPair, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passwd)
	if err != nil {
		return KeyPair{}, err
	}

	res := ed25519.NewKeyFromSeed(seed)
	var pk [64]byte
	copy(pk[:], res)
	kp = kb.writeLocalKeyPair(tmed.PrivKeyEd25519(pk), passwd)
	return kp, nil
}

// ImportPrivKey imports a private key in ASCII armor format.
// It returns an error if a key with the same address exists or a wrong decryptPassphrase is
// supplied.
func (kb dbKeybase) ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) error {
	privKey, err := mintkey.UnarmorDecryptPrivKey(armor, decryptPassphrase)
	if err != nil {
		return err
	}

	accAddress, err := types.AccAddressFromHex(privKey.PubKey().Address().String())
	if err != nil {
		return err
	}

	if _, err := kb.Get(accAddress); err == nil {
		return errors.New("Cannot overwrite key with address: " + accAddress.String())
	}

	kb.writeLocalKeyPair(privKey, encryptPassphrase)
	return nil
}

// ExportPrivKeyEncryptedArmor finds the KeyPair by the address, decrypts the armor private key,
// and returns an encrypted armored private key string
func (kb dbKeybase) ExportPrivKeyEncryptedArmor(address types.AccAddress, decryptPassphrase, encryptPassphrase string) (armor string, err error) {
	priv, err := kb.ExportPrivateKeyObject(address, decryptPassphrase)
	if err != nil {
		return "", err
	}
	return mintkey.EncryptArmorPrivKey(priv, encryptPassphrase), nil
}

// ExportPrivateKeyObject exports raw PrivKey object.
func (kb dbKeybase) ExportPrivateKeyObject(address types.AccAddress, passphrase string) (tmcrypto.PrivKey, error) {
	kp, err := kb.Get(address)
	if err != nil {
		return nil, err
	}

	if kp.PrivKeyArmor == "" {
		err = fmt.Errorf("private key not available")
		return nil, err
	}

	priv, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	if err != nil {
		return nil, err
	}
	return priv, err
}

// CloseDB releases the lock and closes the storage backend.
func (kb dbKeybase) CloseDB() {
	kb.db.Close()
}

// Private interface
func (kb dbKeybase) writeLocalKeyPair(priv tmcrypto.PrivKey, passphrase string) KeyPair {
	// encrypt private key using passphrase
	privArmor := mintkey.EncryptArmorPrivKey(priv, passphrase)
	// make Info
	pub := priv.PubKey()
	localKeyPair := NewKeyPair(pub, privArmor)
	kb.writeKeyPair(localKeyPair)
	return localKeyPair
}

func (kb dbKeybase) writeKeyPair(kp KeyPair) {
	// write the info by key
	key := addrKey(kp.GetAddress())
	serializedInfo := writeKeyPair(kp)
	kb.db.SetSync(key, serializedInfo)
}

func addrKey(address types.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s", address.String()))
}
