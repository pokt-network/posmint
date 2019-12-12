package types

import (
	"github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestNewQuerySigningInfoParams(t *testing.T) {
	type args struct {
		consAddr types.ConsAddress
	}
	tests := []struct {
		name string
		args args
		want QuerySigningInfoParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuerySigningInfoParams(tt.args.consAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuerySigningInfoParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQuerySigningInfosParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QuerySigningInfosParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuerySigningInfosParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuerySigningInfosParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryStakedValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryStakedValidatorsParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryStakedValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryStakedValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakedValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakedValidatorsParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakedValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakedValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryUnstakingValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryUnstakingValidatorsParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryUnstakingValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryUnstakingValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryValidatorParams(t *testing.T) {
	type args struct {
		validatorAddr types.ValAddress
	}
	tests := []struct {
		name string
		args args
		want QueryValidatorParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryValidatorParams(tt.args.validatorAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryValidatorParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueryValidatorsParams(t *testing.T) {
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want QueryValidatorsParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryValidatorsParams(tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryValidatorsParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
