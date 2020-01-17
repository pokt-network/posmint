package exported

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ValidatorI expected validator functions
type ValidatorI interface {
	IsStaked() bool                 // check if has a bonded status
	IsUnstaked() bool               // check if has status unbonded
	IsUnstaking() bool              // check if has status unbonding
	IsJailed() bool                 // whether the validator is jailed
	GetStatus() sdk.BondStatus      // status of the validator
	GetAddress() sdk.Address        // operator address to receive/return validators coins
	GetPublicKey() crypto.PublicKey // validation consensus pubkey
	GetTokens() sdk.Int             // validation tokens
	GetConsensusPower() int64       // validation power in tendermint
}
