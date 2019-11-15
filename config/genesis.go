package config

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

type GenesisState map[string]json.RawMessage

// expected usage
//func (app *nameServiceApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
//	genesisState := GetGensisFromFile(app.cdc, "genesis.go")
//	return app.mm.InitGenesis(ctx, genesisState)
//}
func GetGensisFromFile(cdc *codec.Codec, genFile string) GenesisState {
	if !common.FileExists(genFile) {
		panic(fmt.Errorf("%s does not exist, run `init` first", genFile))
	}
	genDoc, err := tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		panic(err)
	}
	genesisState, err := GenesisStateFromGenDoc(cdc, *genDoc)
	if err != nil {
		panic(err)
	}
	return genesisState
}

func GenesisStateFromGenDoc(cdc *codec.Codec, genDoc tmtypes.GenesisDoc) (genesisState map[string]json.RawMessage, err error) {
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}
	return genesisState, nil
}

// InitConfig common config options for init
type InitConfig struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey crypto.PubKey
}

// NewInitConfig creates a new InitConfig object
func NewInitConfig(chainID, genTxsDir, name, nodeID string, valPubKey crypto.PubKey) InitConfig {
	return InitConfig{
		ChainID:   chainID,
		GenTxsDir: genTxsDir,
		Name:      name,
		NodeID:    nodeID,
		ValPubKey: valPubKey,
	}
}
