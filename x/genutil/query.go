package genutil

import (
	"github.com/pokt-network/posmint/client/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/genutil/types"
)

func QueryGenesisTxs(cliCtx context.CLIContext) ([]sdk.Tx, error) {
	resultGenesis, err := cliCtx.Client.Genesis()
	if err != nil {
		return nil, err
	}

	appState, err := types.GenesisStateFromGenDoc(cliCtx.Codec, *resultGenesis.Genesis)
	if err != nil {
		return nil, err
	}

	genState := types.GetGenesisStateFromAppState(cliCtx.Codec, appState)
	genTxs := make([]sdk.Tx, len(genState.GenTxs))
	for i, tx := range genState.GenTxs {
		err := cliCtx.Codec.UnmarshalJSON(tx, &genTxs[i])
		if err != nil {
			return nil, err
		}
	}
	return genTxs, nil
}
