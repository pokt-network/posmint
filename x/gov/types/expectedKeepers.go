package types

import (
	sdk "github.com/pokt-network/posmint/types"
	posexported "github.com/pokt-network/posmint/x/pos/exported"
	supplyexported "github.com/pokt-network/posmint/x/supply/exported"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	// get total supply of tokens
	GetSupply(ctx sdk.Ctx) supplyexported.SupplyI
	// get the address of a module account
	GetModuleAddress(name string) sdk.Address
	// get the module account structure
	GetModuleAccount(ctx sdk.Ctx, moduleName string) supplyexported.ModuleAccountI
	// set module account structure
	SetModuleAccount(sdk.Ctx, supplyexported.ModuleAccountI)
	// send coins to/from module accounts
	SendCoinsFromModuleToModule(ctx sdk.Ctx, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	// send coins from module to validator
	SendCoinsFromModuleToAccount(ctx sdk.Ctx, senderModule string, recipientAddr sdk.Address, amt sdk.Coins) sdk.Error
	// send coins from validator to module
	SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	// mint coins
	MintCoins(ctx sdk.Ctx, moduleName string, amt sdk.Coins) sdk.Error
	// burn coins
	BurnCoins(ctx sdk.Ctx, name string, amt sdk.Coins) sdk.Error
}

type PosKeeper interface {
	RewardForRelays(ctx sdk.Ctx, relays sdk.Int, address sdk.Address)
	GetStakedTokens(ctx sdk.Ctx) sdk.Int
	Validator(ctx sdk.Ctx, addr sdk.Address) posexported.ValidatorI
	TotalTokens(ctx sdk.Ctx) sdk.Int
	BurnForChallenge(ctx sdk.Ctx, challenges sdk.Int, address sdk.Address)
	JailValidator(ctx sdk.Ctx, addr sdk.Address)
	AllValidators(ctx sdk.Ctx) (validators []posexported.ValidatorI)
	GetStakedValidators(ctx sdk.Ctx) (validators []posexported.ValidatorI)
	SessionBlockFrequency(ctx sdk.Ctx) (res int64)
	StakeDenom(ctx sdk.Ctx) (res string)
}
