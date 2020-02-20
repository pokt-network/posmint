package keeper

import (
	"fmt"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/supply/internal/types"
)

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an Address
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Ctx, senderModule string,
	recipientAddr sdk.Address, amt sdk.Coins) sdk.Error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Ctx, senderModule, recipientModule string, amt sdk.Coins) sdk.Error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an Address to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address,
	recipientModule string, amt sdk.Coins) sdk.Error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Minter) {
		panic(fmt.Sprintf("module account %s does not have permissions to mint tokens", moduleName))
	}

	_, err := k.bk.AddCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply = supply.Inflate(amt)

	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Burner) {
		panic(fmt.Sprintf("module account %s does not have permissions to burn tokens", moduleName))
	}

	_, err := k.bk.SubtractCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply = supply.Deflate(amt)
	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))

	return nil
}
