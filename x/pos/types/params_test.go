package types

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
	"reflect"
	"testing"
	"time"
)

func TestDefaultParams(t *testing.T) {
	tests := []struct {
		name string
		want Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustUnmarshalParams(t *testing.T) {
	type args struct {
		cdc   *codec.Codec
		value []byte
	}
	tests := []struct {
		name string
		args args
		want Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustUnmarshalParams(tt.args.cdc, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustUnmarshalParams() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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

func TestParams_ParamSetPairs(t *testing.T) {
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
		name   string
		fields fields
		want   params.ParamSetPairs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
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
			if got := p.ParamSetPairs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParamSetPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_String(t *testing.T) {
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
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
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
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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

func TestUnmarshalParams(t *testing.T) {
	type args struct {
		cdc   *codec.Codec
		value []byte
	}
	tests := []struct {
		name       string
		args       args
		wantParams Params
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, err := UnmarshalParams(tt.args.cdc, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("UnmarshalParams() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}
