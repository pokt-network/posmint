package keeper

import (
	"fmt"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/go-amino"
)

// award coins to an address (will be called at the beginning of the next block)
func (k Keeper) AwardCoinsTo(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) {
	award, _ := k.getValidatorAward(ctx, address)
	k.setValidatorAward(ctx, award.Add(amount), address)
}

// rewardFromFees handles distribution of the collected fees
func (k Keeper) rewardFromFees(ctx sdk.Ctx, previousProposer sdk.Address) {
	logger := k.Logger(ctx)
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.getFeePool(ctx)
	feesCollected := feeCollector.GetCoins()
	// transfer collected fees to the pos module account
	err := k.authKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}
	// calculate the total reward by adding relays to the
	totalReward := feesCollected.AmountOf(k.StakeDenom(ctx))
	// get the validator structure
	proposerValidator := k.Validator(ctx, previousProposer)
	if proposerValidator != nil {
		propRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), totalReward))
		// send to validator
		if err := k.authKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
			sdk.Address(proposerValidator.GetAddress()), propRewardCoins); err != nil {
			panic(err)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, totalReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetAddress().String()),
			),
		)
	} else {
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unstaked completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}
}

// called on begin blocker
func (k Keeper) mintValidatorAwards(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AwardValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.Int{}
		address := sdk.Address(types.AddressFromKey(iterator.Key()))
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &amount)
		k.mint(ctx, amount, address)
		// remove from the award store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a validator award
func (k Keeper) setValidatorAward(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyForValidatorAward(address)
	val := amino.MustMarshalBinaryBare(amount)
	store.Set(key, val)
}

func (k Keeper) getValidatorAward(ctx sdk.Ctx, address sdk.Address) (coins sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorAward(address))
	if value == nil {
		return sdk.ZeroInt(), false
	}
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	found = true
	return
}

func (k Keeper) deleteValidatorAward(ctx sdk.Ctx, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorAward(address))
}

// Mints sdk.Coins and sends them onto an address
func (k Keeper) mint(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.authKeeper.MintCoins(ctx, types.StakedPoolName, coins.Add(coins))
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := k.authKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, sdk.Address(address), coins)
	if sendErr != nil {
		return sendErr.Result()
	}

	logString := amount.String() + " was successfully minted to " + address.String()
	return sdk.Result{
		Log: logString,
	}
}

// get the proposer public key for this block
func (k Keeper) GetPreviousProposer(ctx sdk.Ctx) (address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &address)
	return
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Ctx, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(address)
	store.Set(types.ProposerKey, b)
}

// returns the current BaseProposerReward rate from the global param store
// nolint: errcheck
func (k Keeper) getProposerRewardPercentage(ctx sdk.Ctx) sdk.Int {
	return sdk.NewInt(int64(k.ProposerRewardPercentage(ctx)))
}
