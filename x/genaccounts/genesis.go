/*
Package genaccounts contains specialized functionality for initializing
accounts from genesis including:
 - genesis account validation,
 - initchain processing of genesis accounts,
 - export processing (to genesis) of accounts,
 - server command for adding accounts to the genesis file.
*/

package genaccounts

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	authexported "github.com/pokt-network/posmint/x/auth/exported"
	"github.com/pokt-network/posmint/x/genaccounts/types"
)

// InitGenesis initializes accounts and deliver genesis transactions
func InitGenesis(ctx sdk.Context, _ *codec.Codec, accountKeeper types.AccountKeeper, genesisState types.GenesisState) {
	genesisState.Sanitize()

	// load the accounts
	for _, gacc := range genesisState {
		acc := gacc.ToAccount()
		acc = accountKeeper.NewAccount(ctx, acc) // set account number
		accountKeeper.SetAccount(ctx, acc)
	}
}

// ExportGenesis exports genesis for all accounts
func ExportGenesis(ctx sdk.Context, accountKeeper types.AccountKeeper) types.GenesisState {

	// iterate to get the accounts
	var accounts []types.GenesisAccount
	accountKeeper.IterateAccounts(ctx,
		func(acc authexported.Account) (stop bool) {
			account, err := types.NewGenesisAccountI(acc)
			if err != nil {
				panic(err)
			}
			accounts = append(accounts, account)
			return false
		},
	)

	return accounts
}
