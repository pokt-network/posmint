package config

import (
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
)

const (
	stepNone int8 = 0 // Used to distinguish the initial state
)

// LoadOrGenFilePV loads a FilePV from the given filePaths
// or else generates a new one and saves it to the filePaths.
func LoadOrGenFilePV(keyFilePath, stateFilePath string) *privval.FilePV {
	var pv *privval.FilePV
	if cmn.FileExists(keyFilePath) {
		pv = LoadFilePV(keyFilePath, stateFilePath)
	} else {
		pv = GenFilePV(keyFilePath, stateFilePath)
		pv.Save()
	}
	return pv
}

// GenFilePV generates a new validator with randomly generated private key
// and sets the filePaths, but does not call Save().
func GenFilePV(keyFilePath, stateFilePath string) *privval.FilePV {
	privKey := ed25519.GenPrivKey()
	return &privval.FilePV{
		Key: privval.FilePVKey{
			privKey.PubKey().Address(),
			privKey.PubKey(),
			privKey,
			keyFilePath,
		},
		LastSignState: privval.FilePVLastSignState{
			0,
			0,
			stepNone,
			nil,
			nil,
			stateFilePath,
		},
	}
}

func LoadFilePV(privValKeyFile, privValStateFile string) *privval.FilePV {
	return privval.LoadFilePV(privValKeyFile, privValStateFile)
}
