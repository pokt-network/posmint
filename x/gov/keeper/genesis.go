package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis - Init store state from genesis data
func (k Keeper) InitGenesis(ctx sdk.Ctx, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)
	// validate acl
	if err := k.GetACL(ctx).Validate(k.GetAllParamNames(ctx)); err != nil {
		panic(err)
	}
	dao := k.GetDAOAccount(ctx)
	if dao == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DAOAccountName))
	}
	err := dao.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, data.DAOTokens)))
	if err != nil {
		panic(fmt.Sprintf("unable to set dao tokens: %s", err.Error()))
	}
	k.SupplyKeeper.SetModuleAccount(ctx, dao)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Ctx) types.GenesisState {
	return types.NewGenesisState(k.GetParams(ctx), k.GetDAOTokens(ctx))
}
