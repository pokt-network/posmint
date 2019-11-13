package config

import (
	"encoding/json"
	"github.com/pokt-network/posmint/codec"
	"io/ioutil"
	"os"
)

type GenesisState map[string]json.RawMessage

// expected usage
//func (app *nameServiceApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
//	genesisState := GetGensisFromFile(app.cdc, "genesis.go")
//	return app.mm.InitGenesis(ctx, genesisState)
//}
func GetGensisFromFile(cdc *codec.Codec, genFilePath string) GenesisState {
	var genesisState GenesisState
	jsonFile, err := os.Open(genFilePath)
	if err != nil {
		panic("unable to open genesisFile: " + err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = cdc.UnmarshalJSON(byteValue, &genesisState)
	if err != nil {
		panic(err)
	}
	return genesisState
}
