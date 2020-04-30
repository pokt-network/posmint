package types

import (
	"bytes"
	"fmt"
	"strings"

	sdk "github.com/pokt-network/posmint/types"
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
)

var _ sdk.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	MaxMemoCharacters uint64 `json:"max_memo_characters" yaml:"max_memo_characters"`
	TxSigLimit        uint64 `json:"tx_sig_limit" yaml:"tx_sig_limit"`
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
		{Key: KeyMaxMemoCharacters, Value: &p.MaxMemoCharacters},
		{Key: KeyTxSigLimit, Value: &p.TxSigLimit},
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
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("ACLKey: %d\n", p.MaxMemoCharacters))
	sb.WriteString(fmt.Sprintf("TxSigLimit: %d\n", p.TxSigLimit))
	return sb.String()
}
