package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}
	// retrieve the previous proposer from store
	// and allocate tokens accordingly
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the new proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)

	// slashing begin blocker below todo

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unstake any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}

	// Iterate through any newly discovered evidence of infraction
	// Slash any validators (and since-unstaked stake within the unstaking period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			k.HandleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}
}

// AllocateTokens handles distribution of the collected fees
func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, previousVotes []abci.VoteInfo,
) { // todo may be able to remove some params

	logger := k.Logger(ctx)

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.GetFeePool(ctx)
	feesCollected := feeCollector.GetCoins()

	// transfer collected fees to the pos module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}
	bonusProposerReward := sdk.NewInt(0) // TODO from % of relays
	// calculate the total reward by adding relays to the
	totalReward := feesCollected.AmountOf(k.StakeDenom(ctx)).Add(bonusProposerReward)
	// calculate previous proposer reward
	baseProposerRewardPercentage := k.GetBaseProposerRewardPercentage(ctx)
	// divide up the reward from the proposer reward and the dao reward
	proposerReward := baseProposerRewardPercentage.Mul(totalReward)
	daoReward := totalReward.Sub(proposerReward)
	// get the validator structure
	proposerValidator := k.ValidatorByConsAddr(ctx, previousProposer)

	if proposerValidator != nil {
		propRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), proposerReward))
		daoRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), daoReward))
		// send to validator
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
			sdk.AccAddress(proposerValidator.GetAddress()), propRewardCoins); err != nil {
			panic(err)
		}
		// send to rest dao
		if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.DAOPoolName, daoRewardCoins); err != nil {
			panic(err)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetAddress().String()),
			),
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDAOAllocation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, daoReward.String()),
			),
		)
	} else {
		// previous proposer can be unknown if say, the unstaking period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unstaked completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}
}

// Mints sdk.Coins
func (k Keeper) Mint(ctx sdk.Context, amount sdk.Int, address sdk.ValAddress) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins.Add(coins))
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), coins)
	if sendErr != nil {
		return sendErr.Result()
	}

	logString := amount.String() + " was successfully minted to " + address.String()
	return sdk.Result{
		Log: logString,
	}
}

// get the proposer public key for this block
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) (consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(types.ProposerKey, b)
}

// returns the current BaseProposerReward rate from the global param store
// nolint: errcheck
func (k Keeper) GetBaseProposerRewardPercentage(ctx sdk.Context) sdk.Int {
	return sdk.NewInt(int64(k.BaseProposerAward(ctx))).Quo(sdk.NewInt(100))
}
