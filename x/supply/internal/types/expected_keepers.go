package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Ctx, process func(exported.Account) (stop bool))
	GetAccount(sdk.Ctx, sdk.Address) exported.Account
	SetAccount(sdk.Ctx, exported.Account)
	NewAccount(sdk.Ctx, exported.Account) exported.Account
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(ctx sdk.Ctx, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
}
