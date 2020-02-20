package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
)

// GetConstantFee get's the constant fee from the paramSpace
func (k Keeper) GetConstantFee(ctx sdk.Ctx) (constantFee sdk.Coin) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyConstantFee, &constantFee)
	return
}

// GetConstantFee set's the constant fee in the paramSpace
func (k Keeper) SetConstantFee(ctx sdk.Ctx, constantFee sdk.Coin) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyConstantFee, constantFee)
}
