package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want GenesisState
	}{{"defaultState", GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
		DAO:          DAOPool(NewPool(sdk.ZeroInt())),
	}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultGenesisState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGenesisState(t *testing.T) {
	type args struct {
		params           Params
		validators       []Validator
		dao              DAOPool
		previousProposer sdk.ConsAddress
		signingInfos     map[string]ValidatorSigningInfo
		missedBlocks     map[string][]MissedBlock
	}
	tests := []struct {
		name string
		args args
		want GenesisState
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenesisState(tt.args.params, tt.args.validators, tt.args.dao, tt.args.previousProposer, tt.args.signingInfos, tt.args.missedBlocks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}
