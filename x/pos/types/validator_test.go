package types

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewValidator(t *testing.T) {
	type args struct {
		addr          sdk.ValAddress
		consPubKey    crypto.PubKey
		tokensToStake sdk.Int
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want Validator
	}{
		{"defaultValidator", args{sdk.ValAddress(pub.Address()), pub, sdk.ZeroInt()},
			Validator{
				Address:                 sdk.ValAddress(pub.Address()),
				ConsPubKey:              pub,
				Jailed:                  false,
				Status:                  sdk.Bonded,
				StakedTokens:            sdk.ZeroInt(),
				UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: bonded
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValidator(tt.args.addr, tt.args.consPubKey, tt.args.tokensToStake); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ABCIValidatorUpdate(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.ValidatorUpdate
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
			if got := v.ABCIValidatorUpdate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ABCIValidatorUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ABCIValidatorUpdateZero(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.ValidatorUpdate
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
			if got := v.ABCIValidatorUpdateZero(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ABCIValidatorUpdateZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_AddStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
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
			if got := v.AddStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
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
			if got := v.ConsAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
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
			if got := v.ConsensusPower(); got != tt.want {
				t.Errorf("ConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_Equals(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	type args struct {
		v2 Validator
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
			v := Validator{
				Address:                 tt.fields.Address,
				ConsPubKey:              tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.Equals(tt.args.v2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.ValAddress
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
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsAddr(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.ConsAddress
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
			if got := v.GetConsAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsPubKey(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   crypto.PubKey
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
			if got := v.GetConsPubKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
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
			if got := v.GetConsensusPower(); got != tt.want {
				t.Errorf("GetConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.BondStatus
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
			if got := v.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Int
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
			if got := v.GetTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsJailed(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
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
			if got := v.IsJailed(); got != tt.want {
				t.Errorf("IsJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsStaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
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
			if got := v.IsStaked(); got != tt.want {
				t.Errorf("IsStaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaked(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
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
			if got := v.IsUnstaked(); got != tt.want {
				t.Errorf("IsUnstaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaking(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
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
			if got := v.IsUnstaking(); got != tt.want {
				t.Errorf("IsUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_PotentialConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
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
			if got := v.PotentialConsensusPower(); got != tt.want {
				t.Errorf("PotentialConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_RemoveStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
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
			if got := v.RemoveStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_UpdateStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.ValAddress
		ConsPubKey              crypto.PubKey
		Jailed                  bool
		Status                  sdk.BondStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	type args struct {
		newStatus sdk.BondStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
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
			if got := v.UpdateStatus(tt.args.newStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
