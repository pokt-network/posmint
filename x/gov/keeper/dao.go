package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov/types"
	"github.com/pokt-network/posmint/x/supply/exported"
)

func (k Keeper) DAOTransferFrom(ctx sdk.Ctx, owner, to sdk.Address, amount sdk.Int) sdk.Result {
	if !k.GetDAOOwner(ctx).Equals(owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("non dao owner is trying to transfer from the dao %s", owner.String())).Result()
	}
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, amount))
	err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.DAOAccountName, to, coins)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventDAOTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) DAOBurn(ctx sdk.Ctx, owner sdk.Address, amount sdk.Int) sdk.Result {
	if !k.GetDAOOwner(ctx).Equals(owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("non dao owner is trying to burn from the dao %s", owner.String())).Result()
	}
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, amount))
	err := k.SupplyKeeper.BurnCoins(ctx, types.DAOAccountName, coins)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventDAOBurn,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) GetDAOTokens(ctx sdk.Ctx) sdk.Int {
	return k.GetDAOAccount(ctx).GetCoins().AmountOf(sdk.DefaultStakeDenom)
}

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetDAOAccount(ctx sdk.Ctx) (stakedPool exported.ModuleAccountI) {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.DAOAccountName)
}
