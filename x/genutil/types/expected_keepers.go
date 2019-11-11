package types

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	authexported "github.com/pokt-network/posmint/x/auth/exported"
)

// PosKeeper defines the expected staking keeper (noalias)
type PosKeeper interface {
	UpdateTendermintValidators(sdk.Context) (updates []abci.ValidatorUpdate)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis map[string]json.RawMessage,
		iterateFn func(authexported.Account) (stop bool),
	)
}
