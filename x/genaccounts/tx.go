package genaccounts

import (
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/server"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/genaccounts/types"
	"github.com/pokt-network/posmint/x/genutil"
)

func AddGenesisAccount(ctx *server.Context, cdc *codec.Codec, addr sdk.AccAddress, coins, vestingAmt sdk.Coins,
	dataDirectory string, vestingStart, vestingEnd int64) error {
	config := ctx.Config
	config.SetRoot(dataDirectory)
	genAcc := types.NewGenesisAccountRaw(addr, coins, vestingAmt, vestingStart, vestingEnd, "", "")
	if err := genAcc.Validate(); err != nil {
		return err
	}

	// retrieve the app state
	genFile := config.GenesisFile()
	appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
	if err != nil {
		return err
	}

	// add genesis account to the app state
	var genesisAccounts types.GenesisAccounts
	cdc.MustUnmarshalJSON(appState[types.ModuleName], &genesisAccounts)
	if genesisAccounts.Contains(addr) {
		return fmt.Errorf("cannot add account at existing address %v", addr)
	}

	genesisAccounts = append(genesisAccounts, genAcc)
	genesisStateBz := cdc.MustMarshalJSON(types.GenesisState(genesisAccounts))
	appState[types.ModuleName] = genesisStateBz
	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	// export app state
	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}
