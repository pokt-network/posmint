package keeper

import (
	"bytes"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
	"time"
)

// get the set of all unstaking validators with no limits, used during genesis dump
func (k Keeper) GetAllUnstakingValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakingValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// retrieve all unstaked validators
func (k Keeper) GetAllUnstakedValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakedValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// gets all of the validators who will be unstaked at exactly this time
func (k Keeper) GetUnstakingValidators(ctx sdk.Context, unstakingTime time.Time) (valAddrs []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForUnstakingValidators(unstakingTime))
	if bz == nil {
		return []sdk.ValAddress{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return valAddrs
}

// Sets validators in unstaking queue at a certain unstaking time
func (k Keeper) SetUnstakingValidators(ctx sdk.Context, unstakingTime time.Time, keys []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.KeyForUnstakingValidators(unstakingTime), bz)
}

// Deletes all the validators for a specific unstaking time
func (k Keeper) DeleteUnstakingValidators(ctx sdk.Context, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForUnstakingValidators(unstakingTime))
}

// Insert a validator address to the appropriate position in the unstaking queue
func (k Keeper) InsertUnstakingValidator(ctx sdk.Context, val types.Validator) {
	validators := k.GetUnstakingValidators(ctx, val.UnstakingCompletionTime)
	validators = append(validators, val.Address)
	k.SetUnstakingValidators(ctx, val.UnstakingCompletionTime, validators)
}

// Delete a validator address from the unstaking queue
func (k Keeper) DeleteUnstakingValidator(ctx sdk.Context, val types.Validator) {
	validators := k.GetUnstakingValidators(ctx, val.UnstakingCompletionTime)
	var newValidators []sdk.ValAddress
	for _, addr := range validators {
		if !bytes.Equal(addr, val.Address) {
			newValidators = append(newValidators, addr)
		}
	}
	if len(newValidators) == 0 {
		k.DeleteUnstakingValidators(ctx, val.UnstakingCompletionTime)
	} else {
		k.SetUnstakingValidators(ctx, val.UnstakingCompletionTime, newValidators)
	}
}

// iterator for all unstaking validators up to a certain time
func (k Keeper) UnstakingValidatorsIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingValidatorsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingValidators(endTime)))
}

// Returns a concatenated list of all the mature validators
func (k Keeper) GetMatureValidators(ctx sdk.Context) (matureValsAddrs []sdk.ValAddress) {
	// gets an iterator for all validators from time 0 until the current Blockheader time
	unstakingValsIterator := k.UnstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()

	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var validators []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &validators)
		matureValsAddrs = append(matureValsAddrs, validators...)
	}

	return matureValsAddrs
}

// Unstakes all the unstaking validators that have finished their unstaking period
func (k Keeper) UnstakeAllMatureValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	unstakingValidatorsIterator := k.UnstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValidatorsIterator.Close()

	for ; unstakingValidatorsIterator.Valid(); unstakingValidatorsIterator.Next() {
		var unstakingVals []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValidatorsIterator.Value(), &unstakingVals)

		for _, valAddr := range unstakingVals {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic("validator in the unstakeing queue was not found")
			}
			if !val.IsUnstaking() {
				panic("unexpected validator in unstakeing queue; status was not unstakeing")
			}
			err := k.CompleteUnstakingValidator(ctx, val)
			if err != nil {
				panic(err)
			}
		}
		store.Delete(unstakingValidatorsIterator.Key())
	}
}
