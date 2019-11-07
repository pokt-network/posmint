package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	Params                   Params                          `json:"params" yaml:"params"`
	PrevStateTotalPower      sdk.Int                         `json:"prevState_total_power" yaml:"prevState_total_power"`
	PrevStateValidatorPowers []PrevStateBlockValidatorPower       `json:"prevState_validator_powers" yaml:"prevState_validator_powers"`
	Validators               Validators                      `json:"validators" yaml:"validators"`
	Exported                 bool                            `json:"exported" yaml:"exported"`
	DAO                      DAOPool                         `json:"dao" yaml:"dao"`
	SigningInfos             map[string]ValidatorSigningInfo `json:"signing_infos" yaml:"signing_infos"`
	MissedBlocks             map[string][]MissedBlock        `json:"missed_blocks" yaml:"missed_blocks"`
}

// PrevState validator power, needed for validator set update logic
type PrevStateBlockValidatorPower struct {
	Address sdk.ValAddress
	Power   int64
}

func NewGenesisState(params Params, validators []Validator, dao DAOPool,
	signingInfos map[string]ValidatorSigningInfo, missedBlocks map[string][]MissedBlock) GenesisState {
	return GenesisState{
		Params:       params,
		Validators:   validators,
		SigningInfos: signingInfos,
		MissedBlocks: missedBlocks,
		DAO:          dao,
	}
}

// MissedBlock
type MissedBlock struct {
	Index  int64 `json:"index" yaml:"index"`
	Missed bool  `json:"missed" yaml:"missed"`
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
		DAO:          DAOPool(NewPool(sdk.ZeroInt())),
	}
}
