package config

import (
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/p2p"
	"io/ioutil"
)

func LoadOrGenerateNodeKeyFile(rootPath string) error {
	config.DefaultConfig().SetRoot(rootPath)
	nodeKeyFile := config.DefaultConfig().NodeKeyFile()
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFile)
	if err != nil {
		return err
	}
	fmt.Println(nodeKey.ID())
	return nil
}

func GenerateNodeKeyFile(cdc *codec.Codec, filePath string) (*p2p.NodeKey, error) {
	privKey := ed25519.GenPrivKey()
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}

	jsonBytes, err := cdc.MarshalJSON(nodeKey)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}

func LoadNodeKeyFile(cdc *codec.Codec, filePath string) (*p2p.NodeKey, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := new(p2p.NodeKey)
	err = cdc.UnmarshalJSON(jsonBytes, nodeKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading NodeKey from %v: %v", filePath, err)
	}
	return nodeKey, nil
}
