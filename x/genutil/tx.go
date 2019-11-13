package genutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/context"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/genutil/types"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// todo broken
func (am AppModule) GenesisTx(ctx *context.Context, cdc *codec.Codec, mbm module.BasicManager, genAccIterator types.GenesisAccountsIterator, homeDir, fromAddr, amountStaked, nodeIDString, valPubKeyString, keybaseDirectory string) error {
	config := ctx.Config
	config.SetRoot(homeDir)
	nodeID, valPubKey, err := InitializeNodeValidatorFiles(ctx.Config)
	if err != nil {
		return err
	}
	// Read --nodeID, if empty take it from priv_validator.json
	if nodeIDString != "" {
		nodeID = nodeIDString
	}
	// Read --pubkey, if empty take it from priv_validator.json
	if valPubKeyString != "" {
		valPubKey, err = sdk.GetConsPubKeyBech32(valPubKeyString)
		if err != nil {
			return err
		}
	}

	genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		return err
	}

	var genesisState map[string]json.RawMessage
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return err
	}

	if err = mbm.ValidateGenesis(genesisState); err != nil {
		return err
	}

	kb, err := keys.NewKeyBaseFromDir(keybaseDirectory)
	if err != nil {
		return err
	}

	key, err := kb.GetFromAddress(fromAddr)
	if err != nil {
		return err
	}

	// Set flags for creating gentx
	viper.Set(context.FlagHome, viper.GetString(flagClientHome))
	PrepareFlagsForTxCreateValidator(config, nodeID, genDoc.ChainID, valPubKey)

	// Fetch the amount of coins staked
	coins, err := sdk.ParseCoins(amountStaked)
	if err != nil {
		return err
	}

	err = ValidateAccountInGenesis(genesisState, genAccIterator, key.GetAddress(), coins, cdc)
	if err != nil {
		return err
	}

	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(util.GetTxEncoder(cdc))
	cliCtx := util.NewCLIContext(am.node).WithCodec(cdc)

	// create a 'create-validator' message
	txBldr, msg, err := smbh.BuildCreateValidatorMsg(cliCtx, txBldr)
	if err != nil {
		return err
	}

	info, err := txBldr.Keybase().Get(name)
	if err != nil {
		return err
	}

	if info.GetType() == keys.TypeOffline || info.GetType() == keys.TypeMulti {
		fmt.Println("Offline key passed in. Use `tx sign` command to sign:")
		return util.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
	}

	// write the unsigned transaction to the buffer
	w := bytes.NewBuffer([]byte{})
	cliCtx = cliCtx.WithOutput(w)

	if err = util.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}); err != nil {
		return err
	}

	// read the transaction
	stdTx, err := readUnsignedGenTxFile(cdc, w)
	if err != nil {
		return err
	}

	// sign the transaction and write it to the output file
	signedTx, err := util.SignStdTx(txBldr, cliCtx, name, stdTx, false, true)
	if err != nil {
		return err
	}

	// Fetch output file name
	outputDocument := viper.GetString(context.FlagOutputDocument)
	if outputDocument == "" {
		outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
		if err != nil {
			return err
		}
	}

	if err := writeSignedGenTx(cdc, outputDocument, signedTx); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Genesis transaction written to %q\n", outputDocument)
	return nil
}

func Init(ctx *context.Context, cdc *codec.Codec, mbm module.BasicManager,
	homeDirectory, chainID string, overwrite bool) error {
	config := ctx.Config
	config.SetRoot(homeDirectory)

	if chainID == "" {
		chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
	}

	nodeID, _, err := InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	genFile := config.GenesisFile()
	if !overwrite && common.FileExists(genFile) {
		return fmt.Errorf("genesis.json file already exists: %v", genFile)
	}
	appState, err := codec.MarshalJSONIndent(cdc, mbm.DefaultGenesis())
	if err != nil {
		return err
	}

	genDoc := &tmtypes.GenesisDoc{}
	if _, err := os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		genDoc, err = tmtypes.GenesisDocFromFile(genFile)
		if err != nil {
			return err
		}
	}

	genDoc.ChainID = chainID
	genDoc.Validators = nil
	genDoc.AppState = appState
	if err = ExportGenesisFile(genDoc, genFile); err != nil {
		return err
	}

	toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	return displayInfo(cdc, toPrint)
}

func ValidateGen(ctx *context.Context, cdc *codec.Codec, mbm module.BasicManager) error {
	// Load default if passed no args, otherwise load passed file
	var genesis string
	genesis = ctx.Config.GenesisFile()

	_, _ = fmt.Fprintf(os.Stderr, "validating genesis file at %s\n", genesis)

	var genDoc *tmtypes.GenesisDoc
	var err error
	if genDoc, err = tmtypes.GenesisDocFromFile(genesis); err != nil {
		return fmt.Errorf("error loading genesis doc from %s: %s", genesis, err.Error())
	}

	var genState map[string]json.RawMessage
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
		return fmt.Errorf("error unmarshaling genesis doc %s: %s", genesis, err.Error())
	}

	if err = mbm.ValidateGenesis(genState); err != nil {
		return fmt.Errorf("error validating genesis file %s: %s", genesis, err.Error())
	}

	// TODO test to make sure initchain doesn't panic

	fmt.Printf("File at %s is a valid genesis file\n", genesis)
	return nil
}

func CollectGenTx(ctx *context.Context, cdc *codec.Codec,
	genAccIterator types.GenesisAccountsIterator, defaultNodeHome, genTxsDir string) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))
	name := viper.GetString(context.FlagName)
	nodeID, valPubKey, err := InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		return err
	}

	if genTxsDir == "" {
		genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
	}

	toPrint := newPrintInfo(config.Moniker, genDoc.ChainID, nodeID, genTxsDir, json.RawMessage(""))
	initCfg := types.NewInitConfig(genDoc.ChainID, genTxsDir, name, nodeID, valPubKey)

	appMessage, err := GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
	if err != nil {
		return err
	}

	toPrint.AppMessage = appMessage

	// print out some key information
	return displayInfo(cdc, toPrint)
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := common.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}
	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(cdc *codec.Codec, r io.Reader) (auth.StdTx, error) {
	var stdTx auth.StdTx
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return stdTx, err
	}
	err = cdc.UnmarshalJSON(bytes, &stdTx)
	return stdTx, err
}

func writeSignedGenTx(cdc *codec.Codec, outputDocument string, tx auth.StdTx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}
