package keys

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pokt-network/posmint/crypto/keys/mintkey"
	"github.com/pokt-network/posmint/types"

	//tmed "github.com/tendermint/crypto/ed25519"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"

	dbm "github.com/tendermint/tm-db"
)

var _ Keybase = dbKeybase{}

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
	return readKeyPair(ik)
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security).
// It returns an error if the key doesn't exist or
// passphrases don't match.
func (kb dbKeybase) Delete(address types.AccAddress, passphrase string) error {
	// verify we have the key in the keybase
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}

	// Verify passphrase matches
	if _, err = mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase); err != nil {
		return err
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

// Create a new KeyPair and encrypt it to disk using encryptPassphrase
func (kb dbKeybase) Create(encryptPassphrase string) (KeyPair, error) {
	//_, privKey, err := ed25519.GenerateKey(nil)
	privKey := tmed25519.GenPrivKey()
	//if err != nil {
	//	return KeyPair{}, err
	//}
	//var privKeyBytes [64]byte
	//copy(privKeyBytes[:], privKey.Bytes())
	kp := kb.writeLocalKeyPair(privKey, encryptPassphrase)
	return kp, nil
}

// ImportPrivKey imports a private key in ASCII armor format.
// It returns an error if a key with the same address exists or a wrong decryptPassphrase is
// supplied.
func (kb dbKeybase) ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) (KeyPair, error) {
	privKey, err := mintkey.UnarmorDecryptPrivKey(armor, decryptPassphrase)
	if err != nil {
		return KeyPair{}, err
	}

	accAddress, err := types.AccAddressFromHex(privKey.PubKey().Address().String())
	if err != nil {
		return KeyPair{}, err
	}

	if _, err := kb.Get(accAddress); err == nil {
		return KeyPair{}, errors.New("Cannot overwrite key with address: " + accAddress.String())
	}

	return kb.writeLocalKeyPair(privKey, encryptPassphrase), nil
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

// ImportPrivateKeyObject using the raw unencrypted privateKey string and encrypts it to disk using encryptPassphrase
func (kb dbKeybase) ImportPrivateKeyObject(privateKey [64]byte, encryptPassphrase string) (KeyPair, error) {
	ed25519PK := tmed25519.PrivKeyEd25519(privateKey)
	fmt.Println(ed25519PK.PubKey().Address().String())
	accAddress, err := types.AccAddressFromHex(ed25519PK.PubKey().Address().String())
	if err != nil {
		return KeyPair{}, err
	}
	if _, err := kb.Get(accAddress); err == nil {
		return KeyPair{}, errors.New("Cannot overwrite key with address: " + accAddress.String())
	}
	return kb.writeLocalKeyPair(ed25519PK, encryptPassphrase), nil
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
