package types

import (
	"fmt"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params   Params `json:"params" yaml:"params"`
	Accounts Accounts
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, accounts Accounts) GenesisState {
	return GenesisState{
		Params:   params,
		Accounts: accounts,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), nil)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSizeCostPerByte == 0 {
		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
	}
	return nil
}
