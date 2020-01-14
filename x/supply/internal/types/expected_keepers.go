package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(exported.Account) (stop bool))
	GetAccount(sdk.Context, sdk.Address) exported.Account
	SetAccount(sdk.Context, exported.Account)
	NewAccount(sdk.Context, exported.Account) exported.Account
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
}
