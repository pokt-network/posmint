package types

import (
	"encoding/binary"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestAddressFromPrevStateValidatorPowerKey(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{{"sampleByteArray", args{key: []byte{0x51, 0x41, 0x33}}, []byte{0x41, 0x33}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddressFromKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddressFromKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAddrPubkeyRelationKey(t *testing.T) {
	type args struct {
		address []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{{"sampleByteArray", args{address: []byte{0x51, 0x51, 0x51}}, []byte{0x13, 0x51, 0x51, 0x51}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAddrPubkeyRelationKey(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddrPubkeyRelationKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValMissedBlockKey(t *testing.T) {
	type args struct {
		v types.Address
		i int64
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1))

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca, int64(1)}, append(append([]byte{0x12}, ca.Bytes()...), []byte{1, 0, 0, 0, 0, 0, 0, 0}...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValMissedBlockKey(tt.args.v, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValMissedBlockKey() = %v, want %v", got, tt.want)
				t.FailNow()
			}
		})
	}
}

func TestGetValMissedBlockPrefixKey(t *testing.T) {
	type args struct {
		v types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x12}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValMissedBlockPrefixKey(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValMissedBlockPrefixKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValidatorSigningInfoAddress(t *testing.T) {
	type args struct {
		key []byte
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	ca := types.Address(pub.Address())

	tests := []struct {
		name  string
		args  args
		wantV types.Address
	}{
		{"sampleByteArray", args{append([]byte{0x11}, ca.Bytes()...)}, append(ca.Bytes())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := GetValidatorSigningInfoAddress(tt.args.key); !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("GetValidatorSigningInfoAddress() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestGetValidatorSigningInfoKey(t *testing.T) {
	type args struct {
		v types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x11}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValidatorSigningInfoKey(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValidatorSigningInfoKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForUnstakingValidators(t *testing.T) {
	type args struct {
		unstakingTime time.Time
	}
	ut := time.Now()

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ut}, append([]byte{0x41}, types.FormatTimeBytes(ut)...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForUnstakingValidators(tt.args.unstakingTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForUnstakingValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValByAllVals(t *testing.T) {
	type args struct {
		addr types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x21}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValByAllVals(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValByAllVals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorAward(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x51}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorAward(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorAward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorBurn(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x52}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorBurn(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorBurn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorInStakingSet(t *testing.T) {
	type args struct {
		validator Validator
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	operAddrInvr := types.CopyBytes(pub.Address())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"NewValidator", args{validator: NewValidator(types.Address(pub.Address()), pub, types.ZeroInt())}, append([]byte{0x23, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorInStakingSet(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorInStakingSet() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestKeyForValidatorPrevStateStateByPower(t *testing.T) {
	type args struct {
		address types.Address
	}
	ca, _ := types.AddressFromHex("22a3ecfff14962f93614d225828cad8bdb188279")

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"sampleByteArray", args{ca}, append([]byte{0x31}, ca.Bytes()...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyForValidatorPrevStateStateByPower(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyForValidatorPrevStateStateByPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseValidatorPowerRankKey(t *testing.T) {
	type args struct {
		key []byte
	}

	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	operAddrInvr := types.CopyBytes(pub.Address())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	prk := append([]byte{0x23, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)

	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	tests := []struct {
		name         string
		args         args
		wantOperAddr []byte
	}{
		{"samplepowerrankKey", args{key: prk}, operAddrInvr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOperAddr := ParseValidatorPowerRankKey(tt.args.key); !reflect.DeepEqual(gotOperAddr, tt.wantOperAddr) {
				t.Errorf("ParseValidatorPowerRankKey() = %v, want %v", gotOperAddr, tt.wantOperAddr)
			}
		})
	}
}

func Test_getStakedValPowerRankKey(t *testing.T) {
	type args struct {
		validator Validator
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	operAddrInvr := types.CopyBytes(pub.Address())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"NewValidator", args{validator: NewValidator(types.Address(pub.Address()), pub, types.ZeroInt())}, append([]byte{0x23, 0, 0, 0, 0, 0, 0, 0, 0}, operAddrInvr...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStakedValPowerRankKey(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStakedValPowerRankKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
