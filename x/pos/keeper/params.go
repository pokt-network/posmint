package keeper

import (
	"time"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() sdk.KeyTable {
	return sdk.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnstakingTime
func (k Keeper) UnStakingTime(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyUnstakingTime, &res)
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Ctx) (res uint64) {
	k.Paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// StakeDenom - Stakeable coin denomination
func (k Keeper) StakeDenom(ctx sdk.Ctx) (res string) {
	k.Paramstore.Get(ctx, types.KeyStakeDenom, &res)
	return
}

func (k Keeper) MinimumStake(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyStakeMinimum, &res)
	return
}

func (k Keeper) ProposerRewardPercentage(ctx sdk.Ctx) (res int8) {
	k.Paramstore.Get(ctx, types.KeyProposerRewardPercentage, &res)
	return
}

// MaxEvidenceAge - max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

// Downtime slashing threshold
func (k Keeper) MinSignedPerWindow(ctx sdk.Ctx) (res int64) {
	var minSignedPerWindow sdk.Dec
	k.Paramstore.Get(ctx, types.KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64() // todo may have to be int64 .RoundInt64()
}

// Downtime jail duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Ctx) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime
func (k Keeper) SlashFractionDowntime(ctx sdk.Ctx) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		UnstakingTime:            k.UnStakingTime(ctx),
		MaxValidators:            k.MaxValidators(ctx),
		StakeDenom:               k.StakeDenom(ctx),
		StakeMinimum:             k.MinimumStake(ctx),
		ProposerRewardPercentage: k.ProposerRewardPercentage(ctx),
		MaxEvidenceAge:           k.MaxEvidenceAge(ctx),
		SignedBlocksWindow:       k.SignedBlocksWindow(ctx),
		MinSignedPerWindow:       sdk.NewDec(k.MinSignedPerWindow(ctx)),
		DowntimeJailDuration:     k.DowntimeJailDuration(ctx),
		SlashFractionDoubleSign:  k.SlashFractionDoubleSign(ctx),
		SlashFractionDowntime:    k.SlashFractionDowntime(ctx),
	}
}

// set the params
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
