package keys

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/pokt-network/posmint/types"
	sdk "github.com/pokt-network/posmint/types"
)

var _ Keybase = &lazyKeybase{}

type lazyKeybase struct {
	name     string
	dir      string
	coinbase KeyPair
}

// New creates a new instance of a lazy keybase.
func New(name, dir string) Keybase {
	if err := cmn.EnsureDir(dir, 0700); err != nil {
		panic(fmt.Sprintf("failed to create Keybase directory: %s", err))
	}

	return &lazyKeybase{name: name, dir: dir}
}

func (kb *lazyKeybase) GetCoinbase() (KeyPair, error) {
	if kb.coinbase.PrivKeyArmor == "" {
		kps, err := kb.List()
		if err != nil {
			return KeyPair{}, err
		}
		if len(kps) == 0 {
			return KeyPair{}, fmt.Errorf("0 keypairs in the keybase, so could not get a coinbase")
		}
		kb.coinbase = kps[0]
	}
	return kb.coinbase, nil
}

func (kb *lazyKeybase) SetCoinbase(address types.AccAddress) error {
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}
	kb.coinbase = kp
	return nil
}

func (lkb lazyKeybase) List() ([]KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return newDbKeybase(db).List()
}

func (lkb lazyKeybase) Get(address types.AccAddress) (KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return KeyPair{}, err
	}
	defer db.Close()

	return newDbKeybase(db).Get(address)
}

func (lkb lazyKeybase) Delete(address types.AccAddress, passphrase string) error {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).Delete(address, passphrase)
}

func (lkb lazyKeybase) Update(address types.AccAddress, oldpass string, newpass string) error {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).Update(address, oldpass, newpass)
}

func (lkb lazyKeybase) Sign(address types.AccAddress, passphrase string, msg []byte) ([]byte, crypto.PubKey, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	return newDbKeybase(db).Sign(address, passphrase, msg)
}

func (lkb lazyKeybase) Create(encryptPassphrase string) (KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return KeyPair{}, err
	}
	defer db.Close()

	return newDbKeybase(db).Create(encryptPassphrase)
}

func (lkb lazyKeybase) ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) (KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return KeyPair{}, err
	}
	defer db.Close()

	return newDbKeybase(db).ImportPrivKey(armor, decryptPassphrase, encryptPassphrase)
}

func (lkb lazyKeybase) ExportPrivKeyEncryptedArmor(address types.AccAddress, decryptPassphrase, encryptPassphrase string) (armor string, err error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return "", err
	}
	defer db.Close()

	return newDbKeybase(db).ExportPrivKeyEncryptedArmor(address, decryptPassphrase, encryptPassphrase)
}

func (lkb lazyKeybase) ImportPrivateKeyObject(privateKey [64]byte, encryptPassphrase string) (KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return KeyPair{}, err
	}
	defer db.Close()

	return newDbKeybase(db).ImportPrivateKeyObject(privateKey, encryptPassphrase)
}

func (lkb lazyKeybase) ExportPrivateKeyObject(address types.AccAddress, passphrase string) (crypto.PrivKey, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return newDbKeybase(db).ExportPrivateKeyObject(address, passphrase)
}

func (lkb lazyKeybase) CloseDB() {}
