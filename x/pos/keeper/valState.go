package keeper

import (
	"bytes"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"sort"
	"time"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
)

// Apply and return accumulated updates to the staked validator set. Also,
// * Updates the active valset as keyed by LastValidatorsPowerKey.
// * Updates the total power as keyed by PrevStateTotalPower.
// * Updates validator status' according to updated powers.
// * Updates the pool staked vs not-staked tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with positive power, unstakedValidators who were previously staked,
// or were removed from the validator set entirely // todo when would a validator be completely removed?
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate) {
	// get the world state
	store := ctx.KVStore(k.storeKey)
	maxValidators := k.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorsPowerKey).
	prevStateValidators := k.getPreviousStateValidatorsByAddr(ctx) // get all of the valdidators from the last state

	// Iterate over validators, highest power to lowest.
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsByPowerKey)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

		// everything that is iterated in this loop is `
		// becoming a staked validator` -> need to update the old state and tell tendermint
		// or
		// `already part of the staked validator set` -> need to update the power index of everyone a

		// get the validator address
		valAddr := sdk.ValAddress(iterator.Value())
		// return the validator from the current store
		validator := k.mustGetValidator(ctx, valAddr)
		// sanity check for no jailed validators
		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}
		// sanity check for no 0 power validators in the staked set
		if validator.PotentialConsensusPower() == 0 {
			panic("should never have a zero consensus power validator in the staked set")
		}

		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		// get the old power (if was a validator in the previous state)
		oldPowerBytes, found := prevStateValidators[valAddrBytes]
		// calculate the new power
		newPower := validator.ConsensusPower()
		newPowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(newPower)

		// if not found or the power has changed -> add this validator to the updated list
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate())

			// update the previous state as this will soon be the previous state
			k.SetPrevStateValPower(ctx, valAddr, newPower)
		}

		// remove the validator from the previous state validators
		// as this structure will now be used to keep track of who is no longer
		// staked
		delete(prevStateValidators, valAddrBytes)

		// keep count of the number of validators?
		count++
		// update the total power
		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	// sort the no-longer-staked validators
	noLongerStaked := sortNoLongerStakedValidators(prevStateValidators)

	// iterate through the sorted no-longer-staked validators
	for _, valAddrBytes := range noLongerStaked {

		// fetch the validator
		validator := k.mustGetValidator(ctx, valAddrBytes)

		// delete from the stake validator index
		k.DeletePrevStateValPower(ctx, validator.GetAddress())

		// update the validator set
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	return updates
}

// register the validator in the necessary stores in the world state
func (k Keeper) RegisterValidator(ctx sdk.Context, validator types.Validator) {
	k.BeforeValidatorCreated(ctx, validator.Address)    // call before hook
	k.SetValidator(ctx, validator)                      // add validator to global registry (all validators who have staked go here)
	k.SetValidatorByConsAddr(ctx, validator)            // add validator to 2nd global registry (all validators who have staked go here)
	k.SetNewValInStakedPoolByPowerIndex(ctx, validator) // add validator to staked
	k.AddPubKeyRelation(ctx, validator.GetConsPubKey()) // add public key to validator address / public key relation map (used in slashing)
	k.AfterValidatorCreated(ctx, validator.Address)     // call after hook
}

// validate check called before staking
func (k Keeper) ValidateValidatorStaking(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	if !validator.IsUnstaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.coinKeeper.HasCoins(ctx, sdk.AccAddress(validator.Address), coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// perform all the store operations for when a validator status becomes staked
func (k Keeper) StakeValidator(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	// call the before hook
	k.BeforeValidatorStaked(ctx, validator.ConsAddress(), validator.Address)
	// create the coins structure
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	// sanity check on amount of coins
	if !k.coinKeeper.HasCoins(ctx, sdk.AccAddress(validator.Address), coin) {
		panic("should not happen: trying to stake amount with less coins than in account")
	}
	// delete the validator by power index (because it's about to be updated)
	if _, found := k.GetValidator(ctx, validator.Address); found {
		k.DeleteValidatorFromStakingSet(ctx, validator)
	}
	// add coins to the staked field
	validator.AddTokens(amount)
	// send the coins from address to module account
	k.coinsFromUnstakedToStaked(ctx, validator, amount)
	// set the status to staked
	validator = validator.UpdateStatus(sdk.Bonded)
	// save in the validator store
	k.SetValidator(ctx, validator)
	// save in the staked store
	k.SetStakedValByPower(ctx, validator)
	// ensure there's a signing info entry for the validator
	_, found := k.GetValidatorSigningInfo(ctx, validator.ConsAddress())
	if !found {
		signingInfo := types.ValidatorSigningInfo{
			Address:     validator.ConsAddress(),
			StartHeight: ctx.BlockHeight(),
			JailedUntil: time.Unix(0, 0),
		}
		k.SetValidatorSigningInfo(ctx, validator.ConsAddress(), signingInfo)
	}
	// call the after hook
	k.AfterValidatorStaked(ctx, validator.ConsAddress(), validator.Address)
	return nil
}

func (k Keeper) ValidateValidatorBeginUnstaking(ctx sdk.Context, validator types.Validator) sdk.Error {
	// must be staked to begin unstaking
	if !validator.IsStaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	// sanity check
	if validator.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: validator trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// perform all the store operations for when a validator begins uning
func (k Keeper) BeginUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call before unstaking hook
	k.BeforeValidatorBeginUnstaking(ctx, validator.ConsAddress(), validator.Address)
	// get params
	params := k.GetParams(ctx)
	// delete the validator by power index, as the key will change
	k.DeleteValidatorFromStakingSet(ctx, validator)
	// sanity check
	if validator.Status != sdk.Bonded {
		panic(fmt.Sprintf("should not already be unstaked or unstaking, validator: %v\n", validator))
	}
	// set the status
	validator = validator.UpdateStatus(sdk.Unbonding)
	// set the unstaking completion time and completion height appropriately
	validator.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	// save the now unstaked validator record and power index
	k.SetValidator(ctx, validator)
	//// still technically staked -> participates in all burning
	//k.SetStakedValByPower(ctx, validator) todo removed for now
	// Adds to unstaking validator queue
	k.InsertUnstakingValidator(ctx, validator)
	// call after hook
	k.AfterValidatorBeginUnstaking(ctx, validator.ConsAddress(), validator.Address)
	return nil
}

// force unstake
func (k Keeper) ForceValidatorUnstake(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call the before unstaked hook
	k.BeforeValidatorUnstaked(ctx, validator.ConsAddress(), validator.Address)
	// delete the validator from staking set as they are unstaked
	k.DeleteValidatorFromStakingSet(ctx, validator)
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, validator.StakedTokens)
	if err != nil {
		return err
	}
	// remove their tokens from the field
	validator = validator.RemoveTokens(validator.StakedTokens)
	// update their status to unstaked
	validator = validator.UpdateStatus(sdk.Unbonded)
	// set the validator in store
	k.SetValidator(ctx, validator)
	// call after hook
	k.AfterValidatorUnstaked(ctx, validator.ConsAddress(), validator.Address)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
	})
	return nil
}

// perform all the store operations for when a validator status becomes unstaked
func (k Keeper) CompleteUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call the before hook
	k.BeforeValidatorUnstaked(ctx, validator.ConsAddress(), validator.Address)
	// delete the validator from staking set as they are unstaked
	k.DeleteValidatorFromStakingSet(ctx, validator)
	// delete the validator from the unstaking queue
	k.DeleteUnstakingValidator(ctx, validator)
	// amount unstaked = stakedTokens
	amount := sdk.NewInt(validator.StakedTokens.Int64())
	// removed the staked tokens field from validator structure
	validator = validator.RemoveTokens(amount)
	// send the tokens from staking module account to validator account
	k.coinsFromStakedToUnstaked(ctx, validator)
	// update the status to unstaked
	validator = validator.UpdateStatus(sdk.Unbonded)
	// set the validator in the store
	k.SetValidator(ctx, validator)
	// call the after hook
	k.AfterValidatorUnstaked(ctx, validator.ConsAddress(), validator.Address)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
	})
	return nil
}

// send a validator to jail
func (k Keeper) jailValidator(ctx sdk.Context, validator types.Validator) {
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}
	validator.Jailed = true
	k.SetValidator(ctx, validator)
	k.DeleteValidatorFromStakingSet(ctx, validator)
}

// remove a validator from jail
func (k Keeper) unjailValidator(ctx sdk.Context, validator types.Validator) {
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
	}
	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetStakedValByPower(ctx, validator)
}

// map of validator addresses to serialized power
type validatorsByAddr map[[sdk.AddrLen]byte][]byte

// get the last validator set
func (k Keeper) getPreviousStateValidatorsByAddr(ctx sdk.Context) validatorsByAddr {
	last := make(validatorsByAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LastValidatorsPowerKey)
	defer iterator.Close()
	// iterate over the last validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr], powerBytes)
	}
	return last
}

// given a map of remaining validators to previous staked power
// returns the list of validators to be unbstaked, sorted by operator address
func sortNoLongerStakedValidators(last validatorsByAddr) [][]byte {
	// sort the map keys for determinism
	noLongerStaked := make([][]byte, len(last))
	index := 0
	for valAddrBytes := range last {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr, valAddrBytes[:])
		noLongerStaked[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerStaked, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerStaked[i], noLongerStaked[j]) == -1
	})
	return noLongerStaked
}
