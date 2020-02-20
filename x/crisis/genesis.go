package crisis

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/crisis/internal/keeper"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
)

// new crisis genesis
func InitGenesis(ctx sdk.Ctx, keeper keeper.Keeper, data types.GenesisState) {
	keeper.SetConstantFee(ctx, data.ConstantFee)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Ctx, keeper keeper.Keeper) types.GenesisState {
	constantFee := keeper.GetConstantFee(ctx)
	return types.NewGenesisState(constantFee)
}
