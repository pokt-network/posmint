package types

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
	"reflect"
	"testing"
	"time"
)

func TestMustMarshalValidator(t *testing.T) {
	type args struct {
		cdc       *codec.Codec
		validator Validator
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustMarshalValidator(tt.args.cdc, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustMarshalValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustUnmarshalValidator(t *testing.T) {
	type args struct {
		cdc      *codec.Codec
		valBytes []byte
	}
	tests := []struct {
		name string
		args args
		want Validator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustUnmarshalValidator(tt.args.cdc, tt.args.valBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustUnmarshalValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalValidator(t *testing.T) {
	type args struct {
		cdc      *codec.Codec
		valBytes []byte
	}
	tests := []struct {
		name          string
		args          args
		wantValidator Validator
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValidator, err := UnmarshalValidator(tt.args.cdc, tt.args.valBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValidator, tt.wantValidator) {
				t.Errorf("UnmarshalValidator() gotValidator = %v, want %v", gotValidator, tt.wantValidator)
			}
		})
	}
}

func TestValidator_MarshalJSON(t *testing.T) {
	type fields struct {
		Address                 types.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  types.BondStatus
		StakedTokens            types.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			got, err := v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_String(t *testing.T) {
	type fields struct {
		Address                 types.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  types.BondStatus
		StakedTokens            types.Int
		UnstakingCompletionTime time.Time
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
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Address                 types.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  types.BondStatus
		StakedTokens            types.Int
		UnstakingCompletionTime time.Time
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if err := v.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidators_String(t *testing.T) {
	tests := []struct {
		name    string
		v       Validators
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tt.v.String(); gotOut != tt.wantOut {
				t.Errorf("String() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
