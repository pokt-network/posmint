package keys

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/pokt-network/posmint/types"
	sdk "github.com/pokt-network/posmint/types"
)

var _ Keybase = lazyKeybase{}

type lazyKeybase struct {
	name string
	dir  string
}

// New creates a new instance of a lazy keybase.
func New(name, dir string) Keybase {
	if err := cmn.EnsureDir(dir, 0700); err != nil {
		panic(fmt.Sprintf("failed to create Keybase directory: %s", err))
	}

	return lazyKeybase{name: name, dir: dir}
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
		return nil, err
	}
	defer db.Close()

	return newDbKeybase(db).Get(address)
}

func (lkb lazyKeybase) Delete(address types.AccAddress, passphrase string, skipPass bool) error {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).Delete(address, passphrase, skipPass)
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

func (lkb lazyKeybase) CreateMnemonic(bip39Passwd string, passwd string) (kp KeyPair, mnemonic string, err error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return nil, "", err
	}
	defer db.Close()

	return newDbKeybase(db).CreateMnemonic(bip39Passwd, passwd)
}

func (lkb lazyKeybase) DeriveFromMnemonic(mnemonic, bip39Passwd, encryptPasswd string) (KeyPair, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return newDbKeybase(db).Derive(mnemonic, bip39Passwd, encryptPasswd)
}

func (lkb lazyKeybase) ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) error {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).ImportPrivKey(armor, decryptPassphrase, encryptPassphrase)
}

func (lkb lazyKeybase) ExportPrivKeyEncryptedArmor(address types.AccAddress, decryptPassphrase, encryptPassphrase string) (armor string, err error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).ExportPrivKeyEncryptedArmor(address, decryptPassphrase, encryptPassphrase)
}

func (lkb lazyKeybase) ExportPrivateKeyObject(address types.AccAddress, passphrase string) (crypto.PrivKey, error) {
	db, err := sdk.NewLevelDB(lkb.name, lkb.dir)
	if err != nil {
		return err
	}
	defer db.Close()

	return newDbKeybase(db).ExportPrivateKeyObject(address, passphrase)
}

func (lkb lazyKeybase) CloseDB() {}
