package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/exported"
	"github.com/pokt-network/posmint/x/auth/types"
)

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Ctx) (supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKeyPrefix)
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Ctx, supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKeyPrefix, b)
}
