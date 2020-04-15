package auth

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/keeper"
	"github.com/pokt-network/posmint/x/auth/types"
)

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Ctx, ak keeper.Keeper) types.GenesisState {
	params := ak.GetParams(ctx)
	accounts := ak.GetAllAccountsExport(ctx)

	return types.NewGenesisState(params, accounts)
}

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Ctx, k keeper.Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	for _, account := range data.Accounts {
		k.SetAccount(ctx, account)
	}
	// manually set the total supply based on accounts if not provided
	if data.Supply.Empty() {
		var totalSupply sdk.Coins
		k.IterateAccounts(ctx,
			func(acc Account) (stop bool) {
				totalSupply = totalSupply.Add(acc.GetCoins())
				return false
			},
		)

		data.Supply = totalSupply
	}
	k.SetSupply(ctx, types.NewSupply(data.Supply))
}
