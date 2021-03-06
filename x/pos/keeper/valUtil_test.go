package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustGetValidator(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   error
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets validator",
			panics:   false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator},
		},
		{
			name:     "panics if no validator",
			panics:   true,
			args:     args{validator: stakedValidator},
			expected: expected{message: fmt.Errorf("validator record not found for address: %X", stakedValidator.Address)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.panics {
			case true:
				defer func() {
					err := recover()
					assert.Equal(t, test.expected.message, err, "does not cointain error message")
				}()
				_ = keeper.mustGetValidator(context, test.args.validator.Address)
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator := keeper.mustGetValidator(context, test.args.validator.Address)
				assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
			}
		})
	}

}

func TestValidatorByAddress(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
		null      bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets validator",
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator, null: false},
		},
		{
			name:     "nil if not found",
			args:     args{validator: stakedValidator},
			expected: expected{null: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.expected.null {
			case true:
				validator := keeper.Validator(context, test.args.validator.GetAddress())
				assert.Nil(t, validator)
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator := keeper.Validator(context, test.args.validator.GetAddress())
				assert.Equal(t, validator, test.expected.validator, "validator does not match")
			}
		})
	}
}

func TestValidatorCaching(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		bz        []byte
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets validator",
			panics:   false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.args.validator)
			keeper.SetStakedValidator(context, test.args.validator)
			store := context.KVStore(keeper.storeKey)
			bz := store.Get(types.KeyForValByAllVals(test.args.validator.Address))
			validator := keeper.validatorCaching(bz, test.args.validator.Address)
			assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
		})
	}

}

func TestNewValidatorCaching(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		bz        []byte
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
		length    int
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "getPrevStatePowerMap",
			panics:   false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator, length: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.args.validator)
			keeper.SetStakedValidator(context, test.args.validator)
			store := context.KVStore(keeper.storeKey)
			key := types.KeyForValidatorPrevStateStateByPower(test.args.validator.Address)
			store.Set(key, test.args.validator.Address)
			powermap := keeper.getPrevStatePowerMap(context)
			assert.Len(t, powermap, test.expected.length, "does not have correct length")
			var valAddr [sdk.AddrLen]byte
			copy(valAddr[:], key[1:])

			for mapKey, value := range powermap {
				assert.Equal(t, valAddr, mapKey, "key is not correct")
				bz := make([]byte, len(test.args.validator.Address))
				copy(bz, test.args.validator.Address)
				assert.Equal(t, bz, value, "key is not correct")
			}
		})
	}
}
