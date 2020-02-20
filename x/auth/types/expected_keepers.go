package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/supply/exported"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Ctx, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Ctx, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.Address
}
