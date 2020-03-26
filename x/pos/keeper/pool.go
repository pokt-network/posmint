package keeper

import (
	"fmt"
	exported2 "github.com/pokt-network/posmint/x/auth/exported"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/types"
)

// StakedRatio the fraction of the staking tokens which are currently staked
func (k Keeper) StakedRatio(ctx sdk.Ctx) sdk.Dec {
	stakedPool := k.GetStakedPool(ctx)

	stakeSupply := k.TotalTokens(ctx)
	if stakeSupply.IsPositive() {
		return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx)).ToDec().QuoInt(stakeSupply)
	}
	return sdk.ZeroDec()
}

// GetStakedTokens total staking tokens supply which is staked
func (k Keeper) GetStakedTokens(ctx sdk.Ctx) sdk.Int {
	stakedPool := k.GetStakedPool(ctx)
	return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx))
}

// GetUnstakedTokens returns the amount of not staked tokens
func (k Keeper) GetUnstakedTokens(ctx sdk.Ctx) (unstakedTokens sdk.Int) {
	return k.TotalTokens(ctx).Sub(k.GetStakedPool(ctx).GetCoins().AmountOf(k.StakeDenom(ctx)))
}

// TotalTokens staking tokens from the total supply
func (k Keeper) TotalTokens(ctx sdk.Ctx) sdk.Int {
	return k.authKeeper.GetSupply(ctx).GetTotal().AmountOf(k.StakeDenom(ctx))
}

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetStakedPool(ctx sdk.Ctx) (stakedPool exported2.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, types.StakedPoolName)
}

// moves coins from the module account to the validator -> used in unstaking
func (k Keeper) coinsFromStakedToUnstaked(ctx sdk.Ctx, validator types.Validator) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), validator.StakedTokens))
	err := k.authKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, sdk.Address(validator.Address), coins)
	if err != nil {
		panic(err)
	}
}

// moves coins from the module account to validator -> used in staking
func (k Keeper) coinsFromUnstakedToStaked(ctx sdk.Ctx, validator types.Validator, amount sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.authKeeper.SendCoinsFromAccountToModule(ctx, sdk.Address(validator.Address), types.StakedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// burnStakedTokens removes coins from the staked pool module account
func (k Keeper) burnStakedTokens(ctx sdk.Ctx, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		return sdk.ErrNegativeAmount(fmt.Sprintf("%d is not positive", amt.ToDec()))
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amt))
	return k.authKeeper.BurnCoins(ctx, types.StakedPoolName, coins)
}

func (k Keeper) getFeePool(ctx sdk.Ctx) (feePool exported2.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
}
