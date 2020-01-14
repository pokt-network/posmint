package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// SupplyKeeper defines the expected supply keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.Address, recipientModule string, amt sdk.Coins) sdk.Error
}
