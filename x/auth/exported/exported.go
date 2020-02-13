package exported

import (
	"github.com/pokt-network/posmint/crypto"
	"time"

	sdk "github.com/pokt-network/posmint/types"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() sdk.Address
	SetAddress(sdk.Address) error // errors if already set.

	GetPubKey() crypto.PublicKey // can return nil.
	SetPubKey(crypto.PublicKey) error

	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error

	// Calculates the amount of coins that can be sent to other accounts given
	// the current time.
	SpendableCoins(blockTime time.Time) sdk.Coins

	// Ensure that account implements stringer
	String() string
}
