package types

import (
	sdk "github.com/pokt-network/posmint/types"
	authexported "github.com/pokt-network/posmint/x/auth/exported"
	posexported "github.com/pokt-network/posmint/x/pos/exported"
	supplyexported "github.com/pokt-network/posmint/x/supply/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Ctx, process func(authexported.Account) (stop bool))
}

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

// ValidatorSet expected properties for the set of all validators (noalias)
// todo this is here so other modules can conform to this interface
type ValidatorSet interface {
	// iterate through validators by address, execute func for each validator
	IterateAndExecuteOverVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through staked validators by address, execute func for each validator
	IterateAndExecuteOverStakedVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// iterate through the validator set of the prevState block by address, execute func for each validator
	IterateAndExecuteOverPrevStateVals(sdk.Ctx, func(index int64, validator posexported.ValidatorI) (stop bool))
	// get a particular validator by address
	Validator(sdk.Ctx, sdk.Address) posexported.ValidatorI
	// total staked tokens within the validator set
	TotalTokens(sdk.Ctx) sdk.Int
	// jail a validator
	JailValidator(sdk.Ctx, sdk.Address)
	// unjail a validator
	UnjailValidator(sdk.Ctx, sdk.Address)
	// MaxValidators returns the maximum amount of staked validators
	MaxValidators(sdk.Ctx) uint64
}

//_______________________________________________________________________________
// Event Hooks
// These can be utilized to communicate between the pos keeper and another
// keeper which must take particular actions when validators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// POSHooks event hooks for staking validator object (noalias)
type POSHooks interface {
	BeforeValidatorRegistered(ctx sdk.Ctx, valAddr sdk.Address)
	AfterValidatorRegistered(ctx sdk.Ctx, valAddr sdk.Address)
	BeforeValidatorRemoved(ctx sdk.Ctx, valAddr sdk.Address)
	AfterValidatorRemoved(ctx sdk.Ctx, valAddr sdk.Address)
	BeforeValidatorStaked(ctx sdk.Ctx, valAddr sdk.Address)
	AfterValidatorStaked(ctx sdk.Ctx, valAddr sdk.Address)
	BeforeValidatorBeginUnstaking(ctx sdk.Ctx, valAddr sdk.Address)
	AfterValidatorBeginUnstaking(ctx sdk.Ctx, valAddr sdk.Address)
	BeforeValidatorUnstaked(ctx sdk.Ctx, valAddr sdk.Address)
	AfterValidatorUnstaked(ctx sdk.Ctx, valAddr sdk.Address)
	BeforeValidatorSlashed(ctx sdk.Ctx, valAddr sdk.Address, fraction sdk.Dec)
	AfterValidatorSlashed(ctx sdk.Ctx, valAddr sdk.Address, fraction sdk.Dec)
}
