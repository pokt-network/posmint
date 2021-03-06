package types

import (
	"bytes"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
)

// DefaultCodespace defines the default auth module parameter subspace
const DefaultParamspace = ModuleName

// Default parameter values
const ()

// Parameter keys
var (
	ACLKey      = []byte("acl")
	DAOOwnerKey = []byte("daoOwner")
	UpgradeKey  = []byte("upgrade")
)

var _ sdk.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	ACL      ACL         `json:"acl"`
	DAOOwner sdk.Address `json:"dao_owner"`
	Upgrade  Upgrade     `json:"upgrade"`
}

// NewParams creates a new Params object
func NewParams(acl ACL, daoOwner sdk.Address) Params {
	return Params{
		ACL:      acl,
		DAOOwner: daoOwner,
	}
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
		{ACLKey, &p.ACL},
		{DAOOwnerKey, &p.DAOOwner},
		{UpgradeKey, &p.Upgrade},
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
	acl := ACL(make([]ACLPair, 0))
	return Params{
		ACL:      acl,
		DAOOwner: sdk.Address{},
		Upgrade:  NewUpgrade(0, ""),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("ACLKey: %v\n", p.ACL))
	sb.WriteString(fmt.Sprintf("DAOOwnerKey: %s\n", p.DAOOwner))
	sb.WriteString(fmt.Sprintf("UpgradeKey: %v\n", p.Upgrade))
	return sb.String()
}
