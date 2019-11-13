package config

import (
	"encoding/json"
	"github.com/pokt-network/posmint/codec"
	"io/ioutil"
	"os"
)

type GenesisState map[string]json.RawMessage

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
