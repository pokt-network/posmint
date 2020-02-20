package auth

import (
	sdk "github.com/pokt-network/posmint/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func InitGenesis(ctx sdk.Ctx, ak AccountKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
	for _, account := range data.Accounts {
		ak.SetAccount(ctx, account)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Ctx, ak AccountKeeper) GenesisState {
	params := ak.GetParams(ctx)
	accounts := ak.GetAllAccounts(ctx)
	return NewGenesisState(params, accounts)
}
