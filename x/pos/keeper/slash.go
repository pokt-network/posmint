package keeper

import (
	"fmt"
	"time"

	posCrypto "github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) BurnValidator(ctx sdk.Ctx, address sdk.Address, severityPercentage sdk.Dec) {
	curBurn, _ := k.getValidatorBurn(ctx, address)
	newSeverity := curBurn.Add(severityPercentage)
	k.setValidatorBurn(ctx, newSeverity, address)
}

// slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) slash(ctx sdk.Ctx, address sdk.Address, infractionHeight, power int64, slashFactor sdk.Dec) sdk.Error {
	// error check slash
	validator, err := k.validateSlash(ctx, address, infractionHeight, power, slashFactor)
	if err != nil {
		return sdk.ErrInvalidSlash(err.Error()) // invalid slash
	}
	if validator.Address == nil {
		return sdk.ErrInternal(fmt.Sprint("Cant slash nil address")) // invalid slash
	}
	logger := k.Logger(ctx)
	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := sdk.TokensFromConsensusPower(power)
	slashAmount := amount.ToDec().Mul(slashFactor).TruncateInt()
	k.BeforeValidatorSlashed(ctx, validator.Address, slashFactor)
	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(slashAmount, validator.StakedTokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.
	// Deduct from validator's staked tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	validator = k.removeValidatorTokens(ctx, validator, tokensToBurn)
	err = k.burnStakedTokens(ctx, tokensToBurn)
	if err != nil {
		return sdk.ErrBurnStakedTokens(err.Error())
	}
	// if falls below minimum force burn all of the stake
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		err := k.ForceValidatorUnstake(ctx, validator)
		if err != nil {
			return sdk.ErrForceValidatorUnstake(err.Error())
		}
	}
	// Log that a slash occurred
	logger.Info(fmt.Sprintf("validator %s slashed by slash factor of %s; burned %v tokens",
		validator.GetAddress(), slashFactor.String(), tokensToBurn))
	k.AfterValidatorSlashed(ctx, validator.Address, slashFactor)
	return nil
}

func (k Keeper) validateSlash(ctx sdk.Ctx, address sdk.Address, infractionHeight int64, power int64, slashFactor sdk.Dec) (types.Validator, error) {
	logger := k.Logger(ctx)
	if slashFactor.LT(sdk.ZeroDec()) {
		return types.Validator{}, fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor)
	}
	if infractionHeight > ctx.BlockHeight() {
		return types.Validator{}, fmt.Errorf( // Can't slash infractions in the future
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight())
	}
	// see if infraction height is outside of unstaking time
	blockTime := ctx.BlockTime()
	infractionTime := ctx.WithBlockHeight(infractionHeight).BlockTime()
	if blockTime.After(infractionTime.Add(k.UnStakingTime(ctx))) {
		logger.Info(fmt.Sprintf( // could've been overslashed and removed
			"INFO: tried to slash with expired evidence: %s %s", infractionTime, blockTime))
		return types.Validator{}, nil
	}
	validator, found := k.GetValidator(ctx, address)
	if !found {
		logger.Error(fmt.Sprintf( // could've been overslashed and removed
			"WARNING: Ignored attempt to slash a nonexistent validator with address %s, we recommend you investigate immediately",
			address))
		return types.Validator{}, nil
	}
	// should not be slashing an unstaked validator
	if validator.IsUnstaked() {
		return types.Validator{}, fmt.Errorf("should not be slashing unstaked validator: %s", validator.GetAddress())
	}
	return validator, nil
}

// handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infraction
func (k Keeper) handleDoubleSign(ctx sdk.Ctx, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	address, signInfo, validator, err := k.validateDoubleSign(ctx, addr, infractionHeight, timestamp)
	if err != nil {
		panic(err)
	}
	// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height.
	// Note that this *can* result in a negative "distributionHeight", up to -ValidatorUpdateDelay,
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay

	// get the percentage slash penalty fraction
	fraction := k.SlashFractionDoubleSign(ctx)

	// slash validator
	// `power` is the int64 power of the validator as provided to/by Tendermint. This value is validator.StakedTokens as
	// sent to Tendermint via ABCI, and now received as evidence. The fraction is passed in to separately to slash
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
			sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
		),
	)
	err = k.slash(ctx, address, distributionHeight, power, fraction)
	if err != nil {
		ctx.Logger().Error(err.Error())
	}

	// JailValidator validator if not already jailed
	if !validator.IsJailed() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSlash,
				sdk.NewAttribute(types.AttributeKeyJailed, address.String()),
			),
		)
		k.JailValidator(ctx, address)
	}
	// force the validator to unstake if isn't already
	v, found := k.GetValidator(ctx, validator.GetAddress())
	if !found {
		panic(types.ErrNoValidatorFound(k.codespace))
	}
	err = k.ForceValidatorUnstake(ctx, v)
	if err != nil {
		panic(err)
	}
	// Set tombstoned to be true
	signInfo.Tombstoned = true
	// Set jailed until to be forever (max t)
	signInfo.JailedUntil = types.DoubleSignJailEndTime
	// Set validator signing info
	k.SetValidatorSigningInfo(ctx, address, signInfo)
}

func (k Keeper) validateDoubleSign(ctx sdk.Ctx, addr crypto.Address, infractionHeight int64, timestamp time.Time) (address sdk.Address, signInfo types.ValidatorSigningInfo, validator exported.ValidatorI, err sdk.Error) {
	logger := k.Logger(ctx)
	// fetch the validator public key
	address = sdk.Address(addr)
	pubkey, er := k.getPubKeyRelation(ctx, addr)
	if er != nil {
		// Ignore evidence that cannot be handled.
		err = types.ErrCantHandleEvidence(k.Codespace())
		return
	}
	// calculate the age of the evidence
	t := ctx.BlockHeader().Time
	age := t.Sub(timestamp)
	// Reject evidence if the double-sign is too old
	if age > k.MaxEvidenceAge(ctx) {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, age of %d past max age of %d",
			sdk.Address(pubkey.Address()), infractionHeight, age, k.MaxEvidenceAge(ctx)))
		return
	}
	// Get validator and signing info
	validator = k.Validator(ctx, address)
	if validator == nil || validator.IsUnstaked() {
		err = types.ErrNoValidatorFound(k.Codespace())
		return
	}
	// fetch the validator signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, address)
	if !found {
		logger.Info(fmt.Sprintf("Expected signing info for validator %s but not found", address))
		err = types.ErrNoSigningInfoFound(k.Codespace(), address)
		return
	}
	// validator is already tombstoned
	if signInfo.Tombstoned {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already tombstoned", sdk.Address(pubkey.Address()), infractionHeight))
		err = types.ErrValidatorTombstoned(k.Codespace())
		return
	}
	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", sdk.Address(pubkey.Address()), infractionHeight, age))
	return
}

// handle a validator signature, must be called once per validator per block
func (k Keeper) handleValidatorSignature(ctx sdk.Ctx, addr crypto.Address, power int64, signed bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()
	address := sdk.Address(addr)
	pubkey, err := k.getPubKeyRelation(ctx, addr)
	if err != nil {
		panic(fmt.Errorf("Validator consensus-address %s not found", address))
	}
	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, address)
	if !found {
		panic(fmt.Errorf("Expected signing info for validator %s but not found", address))
	}
	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := signInfo.IndexOffset % k.SignedBlocksWindow(ctx)
	signInfo.IndexOffset++
	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.getMissedBlockArray(ctx, address, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetMissedBlockArray(ctx, address, index, true)
		signInfo.MissedBlocksCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetMissedBlockArray(ctx, address, index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}
	if missed {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
				sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
			),
		)
		logger.Info(
			fmt.Sprintf("Absent validator %s (%s) at height %d, %d missed, threshold %d", address, pubkey, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}
	minHeight := signInfo.StartHeight + k.SignedBlocksWindow(ctx)
	maxMissed := k.SignedBlocksWindow(ctx) - k.MinSignedPerWindow(ctx)
	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator := k.Validator(ctx, address)
		if validator != nil && !validator.IsJailed() {
			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				address, minHeight, k.MinSignedPerWindow(ctx)))
			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the PrevStateCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSlash,
					sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
					sdk.NewAttribute(types.AttributeKeyJailed, address.String()),
				),
			)
			err := k.slash(ctx, address, distributionHeight, power, k.SlashFractionDowntime(ctx))
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
			k.JailValidator(ctx, address)
			signInfo.JailedUntil = ctx.BlockHeader().Time.Add(k.DowntimeJailDuration(ctx))
			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon restaking.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearMissedArray(ctx, address)
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed", address),
			)
		}
	}
	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, address, signInfo)
}

func (k Keeper) AddPubKeyRelation(ctx sdk.Ctx, pubkey posCrypto.PublicKey) {
	addr := pubkey.Address()
	k.setAddrPubkeyRelation(ctx, addr, pubkey)
}

func (k Keeper) getPubKeyRelation(ctx sdk.Ctx, address crypto.Address) (posCrypto.PublicKey, error) {
	store := ctx.KVStore(k.storeKey)
	var pubkey posCrypto.PublicKey
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.GetAddrPubkeyRelationKey(address)), &pubkey)
	if err != nil {
		return nil, fmt.Errorf("address %s not found", sdk.Address(address))
	}
	return pubkey, nil
}

func (k Keeper) setAddrPubkeyRelation(ctx sdk.Ctx, addr crypto.Address, pubkey posCrypto.PublicKey) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pubkey)
	store.Set(types.GetAddrPubkeyRelationKey(addr), bz)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Ctx, addr crypto.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddrPubkeyRelationKey(addr))
}

func (k Keeper) getBurnFromSeverity(ctx sdk.Ctx, address sdk.Address, severityPercentage sdk.Dec) sdk.Int {
	val := k.mustGetValidator(ctx, address)
	amount := sdk.TokensFromConsensusPower(val.ConsensusPower())
	slashAmount := amount.ToDec().Mul(severityPercentage).TruncateInt()
	return slashAmount
}

// called on begin blocker
func (k Keeper) burnValidators(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BurnValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		severity := sdk.Dec{}
		address := sdk.Address(types.AddressFromKey(iterator.Key()))
		amino.MustUnmarshalBinaryBare(iterator.Value(), &severity)
		val := k.mustGetValidator(ctx, address)
		err := k.slash(ctx, sdk.Address(address), ctx.BlockHeight(), val.ConsensusPower(), severity)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
		// remove from the burn store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a validator burn
func (k Keeper) setValidatorBurn(ctx sdk.Ctx, amount sdk.Dec, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorBurn(address), amino.MustMarshalBinaryBare(amount))
}

func (k Keeper) getValidatorBurn(ctx sdk.Ctx, address sdk.Address) (coins sdk.Dec, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorBurn(address))
	if value == nil {
		return coins, false
	}
	found = true
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	return
}

func (k Keeper) deleteValidatorBurn(ctx sdk.Ctx, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorBurn(address))
}
