package types

import (
	"reflect"
	"testing"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/types"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want types.GenesisState
	}{{"defaultState", types.GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
		DAO:          DAOPool(NewPool(sdk.ZeroInt())),
	}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.DefaultGenesisState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGenesisState(t *testing.T) {
	type args struct {
		params           types.Params
		validators       []types.Validator
		dao              types.DAOPool
		previousProposer sdk.ConsAddress
		signingInfos     map[string]types.ValidatorSigningInfo
		missedBlocks     map[string][]types.MissedBlock
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
