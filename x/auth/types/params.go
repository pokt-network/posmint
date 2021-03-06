package types

import (
	"bytes"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
)

// Default parameter values
const (
	DefaultMaxMemoCharacters uint64 = 256
	DefaultTxSigLimit        uint64 = 7
)

// Parameter keys
var (
	KeyMaxMemoCharacters = []byte("MaxMemoCharacters")
	KeyTxSigLimit        = []byte("TxSigLimit")
	KeyFeeMultiplier     = []byte("FeeMultipliers")
	DefaultFeeMultiplier = FeeMultipliers{
		FeeMultis: nil,
		Default:   1,
	}
)

var _ sdk.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	MaxMemoCharacters uint64         `json:"max_memo_characters" yaml:"max_memo_characters"`
	TxSigLimit        uint64         `json:"tx_sig_limit" yaml:"tx_sig_limit"`
	FeeMultiplier     FeeMultipliers `json:"fee_multipliers"`
}

// ParamKeyTable for auth module
func ParamKeyTable() sdk.KeyTable {
	return sdk.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() sdk.ParamSetPairs {
	return sdk.ParamSetPairs{
		{KeyMaxMemoCharacters, &p.MaxMemoCharacters},
		{KeyTxSigLimit, &p.TxSigLimit},
		{KeyFeeMultiplier, &p.FeeMultiplier},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxMemoCharacters: DefaultMaxMemoCharacters,
		TxSigLimit:        DefaultTxSigLimit,
		FeeMultiplier:     DefaultFeeMultiplier,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("ACLKey: %d\n", p.MaxMemoCharacters))
	sb.WriteString(fmt.Sprintf("TxSigLimit: %d\n", p.TxSigLimit))
	sb.WriteString(fmt.Sprintf("FeeMultiplier: %v\n", p.FeeMultiplier))
	return sb.String()
}
