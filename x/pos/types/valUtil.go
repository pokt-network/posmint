package types

import (
	"fmt"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
	"time"
)

// Validators is a collection of Validator
type Validators []Validator

func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// MUST return the amino encoded version of this validator
func MustMarshalValidator(cdc *codec.Codec, validator Validator) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(validator)
}

// MUST decode the validator from the bytes
func MustUnmarshalValidator(cdc *codec.Codec, valBytes []byte) Validator {
	validator, err := UnmarshalValidator(cdc, valBytes)
	if err != nil {
		panic(err)
	}
	return validator
}

// unmarshal the validator
func UnmarshalValidator(cdc *codec.Codec, valBytes []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(valBytes, &validator)
	return validator, err
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	bechConsPubKey := sdk.HexConsPub(v.ConsPubKey)

	return fmt.Sprintf(`Validator
  Address:           		  %s
  Validator Cons Pubkey: %s
  Jailed:                     %v
  Status:                     %s
  Tokens:               	  %s
  Unstakeing Completion Time:  %v`,
		v.Address, bechConsPubKey, v.Jailed, v.Status, v.StakedTokens, v.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type hexValidator struct {
	Address                 sdk.ValAddress `json:"operator_address" yaml:"operator_address"` // the hex address of the validator
	ConsPubKey              string         `json:"cons_pubkey" yaml:"cons_pubkey"`           // the hex consensus public key of the validator
	Jailed                  bool           `json:"jailed" yaml:"jailed"`                     // has the validator been jailed from staked status?
	Status                  sdk.BondStatus `json:"status" yaml:"status"`                     // validator status (bonded/unbonding/unbonded)
	StakedTokens            sdk.Int        `json:"stakedTokens" yaml:"stakedTokens"`         // how many staked tokens
	UnstakingCompletionTime time.Time      `json:"unstaking_time" yaml:"unstaking_time"`     // if unstaking, min time for the validator to complete unstaking
}

// MarshalJSON marshals the validator to JSON using Hex
func (v Validator) MarshalJSON() ([]byte, error) {
	bechConsPubKey := sdk.HexConsPub(v.ConsPubKey)
	return codec.Cdc.MarshalJSON(hexValidator{
		Address:                 v.Address,
		ConsPubKey:              bechConsPubKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Hex
func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &hexValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyHex(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*v = Validator{
		Address:                 bv.Address,
		ConsPubKey:              consPubKey,
		Jailed:                  bv.Jailed,
		StakedTokens:            bv.StakedTokens,
		Status:                  bv.Status,
		UnstakingCompletionTime: bv.UnstakingCompletionTime,
	}
	return nil
}
