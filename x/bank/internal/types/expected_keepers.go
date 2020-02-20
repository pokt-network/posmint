package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Ctx, addr sdk.Address) exported.Account

	GetAccount(ctx sdk.Ctx, addr sdk.Address) exported.Account
	GetAllAccounts(ctx sdk.Ctx) []exported.Account
	SetAccount(ctx sdk.Ctx, acc exported.Account)

	IterateAccounts(ctx sdk.Ctx, process func(exported.Account) bool)
}
