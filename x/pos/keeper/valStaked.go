package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/types"
)

// set staked validator
func (k Keeper) SetStakedValidator(ctx sdk.Ctx, validator types.Validator) {
	if validator.Jailed {
		return // jailed validators are not kept in the power index
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorInStakingSet(validator), validator.Address)
}

// delete validator from staked set
func (k Keeper) deleteValidatorFromStakingSet(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorInStakingSet(validator))
}

// Update the staked tokens of an existing validator, update the validators power index key
func (k Keeper) removeValidatorTokens(ctx sdk.Ctx, v types.Validator, tokensToRemove sdk.Int) types.Validator {
	k.deleteValidatorFromStakingSet(ctx, v)
	v = v.RemoveStakedTokens(tokensToRemove)
	k.SetValidator(ctx, v)
	k.SetStakedValidator(ctx, v)
	return v
}

// get the current staked validators sorted by power-rank
func (k Keeper) getStakedValidators(ctx sdk.Ctx) types.Validators {
	maxValidators := k.MaxValidators(ctx)
	validators := make([]types.Validator, maxValidators)
	iterator := k.stakedValsIterator(ctx)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		address := iterator.Value()
		validator := k.mustGetValidator(ctx, address)
		if validator.IsStaked() {
			validators[i] = validator
			i++
		}
	}
	return validators[:i] // trim
}

// returns an iterator for the current staked validators
func (k Keeper) stakedValsIterator(ctx sdk.Ctx) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
}

// iterate through the staked validator set and perform the provided function
func (k Keeper) IterateAndExecuteOverStakedVals(
	ctx sdk.Ctx, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	maxValidators := k.MaxValidators(ctx)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid() && i < int64(maxValidators); iterator.Next() {
		address := iterator.Value()
		validator := k.mustGetValidator(ctx, address)
		if validator.IsStaked() {
			stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}
