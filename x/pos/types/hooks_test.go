package types

import (
	"github.com/pokt-network/posmint/types"
	"testing"
)

//TODO This is just a test to see if hooks are working?

func TestMultiPOSHooks_AfterValidatorBeginUnstaked(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		valAddr  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestMultiPOSHooks_AfterValidatorBeginUnstaking(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		valAddr  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_AfterValidatorRegistered(t *testing.T) {
	type args struct {
		ctx     types.Context
		valAddr types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_AfterValidatorRemoved(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		valAddr  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_AfterValidatorSlashed(t *testing.T) {
	type args struct {
		ctx      types.Context
		valAddr  types.ValAddress
		fraction types.Dec
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_AfterValidatorStaked(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		valAddr  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorBeginUnstaked(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		address  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorBeginUnstaking(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		address  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorRegistered(t *testing.T) {
	type args struct {
		ctx     types.Context
		valAddr types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorSlashed(t *testing.T) {
	type args struct {
		ctx      types.Context
		valAddr  types.ValAddress
		fraction types.Dec
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorStaked(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		address  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestMultiPOSHooks_BeforeValidatorRemoved(t *testing.T) {
	type args struct {
		ctx      types.Context
		consAddr types.ConsAddress
		address  types.ValAddress
	}
	tests := []struct {
		name string
		h    MultiPOSHooks
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
