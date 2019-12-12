package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
	"reflect"
	"testing"
)

func TestMsgBeginUnstake_GetSignBytes(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_GetSigners(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   []sdk.AccAddress
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_Route(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
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
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_Type(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
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
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBeginUnstake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBeginUnstake{
				Address: tt.fields.Address,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_GetSignBytes(t *testing.T) {
	type fields struct {
		FromAddress sdk.ValAddress
		ToAddress   sdk.ValAddress
		Amount      sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_GetSigners(t *testing.T) {
	type fields struct {
		FromAddress sdk.ValAddress
		ToAddress   sdk.ValAddress
		Amount      sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []sdk.AccAddress
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_Route(t *testing.T) {
	type fields struct {
		FromAddress sdk.ValAddress
		ToAddress   sdk.ValAddress
		Amount      sdk.Int
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
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_Type(t *testing.T) {
	type fields struct {
		FromAddress sdk.ValAddress
		ToAddress   sdk.ValAddress
		Amount      sdk.Int
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
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgSend_ValidateBasic(t *testing.T) {
	type fields struct {
		FromAddress sdk.ValAddress
		ToAddress   sdk.ValAddress
		Amount      sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSend{
				FromAddress: tt.fields.FromAddress,
				ToAddress:   tt.fields.ToAddress,
				Amount:      tt.fields.Amount,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_GetSignBytes(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
		PubKey  crypto.PubKey
		Value   sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				Address: tt.fields.Address,
				PubKey:  tt.fields.PubKey,
				Value:   tt.fields.Value,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_GetSigners(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
		PubKey  crypto.PubKey
		Value   sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []sdk.AccAddress
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				Address: tt.fields.Address,
				PubKey:  tt.fields.PubKey,
				Value:   tt.fields.Value,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_Route(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
		PubKey  crypto.PubKey
		Value   sdk.Int
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
			msg := MsgStake{
				Address: tt.fields.Address,
				PubKey:  tt.fields.PubKey,
				Value:   tt.fields.Value,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_Type(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
		PubKey  crypto.PubKey
		Value   sdk.Int
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
			msg := MsgStake{
				Address: tt.fields.Address,
				PubKey:  tt.fields.PubKey,
				Value:   tt.fields.Value,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgStake_ValidateBasic(t *testing.T) {
	type fields struct {
		Address sdk.ValAddress
		PubKey  crypto.PubKey
		Value   sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgStake{
				Address: tt.fields.Address,
				PubKey:  tt.fields.PubKey,
				Value:   tt.fields.Value,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_GetSignBytes(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_GetSigners(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   []sdk.AccAddress
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_Route(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.ValAddress
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
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_Type(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.ValAddress
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
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnjail_ValidateBasic(t *testing.T) {
	type fields struct {
		ValidatorAddr sdk.ValAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnjail{
				ValidatorAddr: tt.fields.ValidatorAddr,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
