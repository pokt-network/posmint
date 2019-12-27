package types

import (
	"github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
	"time"
)

func TestDefaultParams(t *testing.T) {
	tests := []struct {
		name string
		want Params
	}{
		{"Default Test", Params{
			UnstakingTime:            DefaultUnstakingTime,
			MaxValidators:            DefaultMaxValidators,
			StakeMinimum:             DefaultMinStake,
			StakeDenom:               types.DefaultBondDenom,
			ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
			MaxEvidenceAge:           DefaultMaxEvidenceAge,
			SignedBlocksWindow:       DefaultSignedBlocksWindow,
			MinSignedPerWindow:       DefaultMinSignedPerWindow,
			DowntimeJailDuration:     DefaultDowntimeJailDuration,
			SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
			SlashFractionDowntime:    DefaultSlashFractionDowntime,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Equal(t *testing.T) {
	type fields struct {
		UnstakingTime            time.Duration
		MaxValidators            uint64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.Dec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.Dec
		SlashFractionDowntime    types.Dec
	}
	type args struct {
		p2 Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Default Test Equal", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, args{Params{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{}}}, true},
		{"Default Test False", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, args{Params{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           1,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if got := p.Equal(tt.args.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	type fields struct {
		UnstakingTime            time.Duration
		MaxValidators            uint64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.Dec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.Dec
		SlashFractionDowntime    types.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Default Validation Test / Wrong All Parameters", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, true},
		{"Default Validation Test / Wrong StakeDenom", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "",
			StakeMinimum:             2,
			ProposerRewardPercentage: 1,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, true},
		{"Default Validation Test / Valid", fields{
			UnstakingTime:            0,
			MaxValidators:            1000,
			StakeDenom:               "3",
			StakeMinimum:             1,
			ProposerRewardPercentage: 100,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
