package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCoinsFromUnstakedToStaked(t *testing.T){
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	tests := []struct{
		name string
		expected string
		validator types.Validator
		amount sdk.Int
		panics bool
	}{
		{
			name: "unstake coins from pool",
			validator: types.Validator{Address: validatorAddress},
			amount: sdk.NewInt(10),
			panics: false,
		},
		{
			name: 	"panics if negative ammount",
			validator: types.Validator{Address: validatorAddress},
			amount: sdk.NewInt(-1),
			expected: fmt.Sprintf("negative coin amount: -1"),
			panics: true,
		},
		{   name : "panics if no supply is set",
			validator: types.Validator{Address: validatorAddress},
			expected: fmt.Sprintf("insufficient account funds"),
			amount: sdk.NewInt(10),
		    panics: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			switch test.panics{
			case true:
				defer func () {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected)
				}()
				if strings.Contains(test.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100) )
				}
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100) )
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				staked := keeper.GetStakedTokens(context)
				assert.Equal(t, test.amount, staked, "values do not match")
			}
		})
	}
}

func TestCoinsFromStakedToUnstaked(t *testing.T){
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	tests := []struct{
		name string
		amount sdk.Int
		expected string
		validator types.Validator
		panics bool
	}{
		{
			name: "unstake coins from pool",
			validator: types.Validator{Address: validatorAddress, StakedTokens: sdk.NewInt(10) },
			amount: sdk.NewInt(110),
			panics: false,
		},
		{
			name: 	"panics if negative ammount",
			validator: types.Validator{Address: validatorAddress, StakedTokens: sdk.NewInt(-1) },
			amount: sdk.NewInt(-1),
			expected: fmt.Sprintf("negative coin amount: -1"),
			panics: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			switch test.panics{
			case true:
				defer func () {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected)
				}()
				if strings.Contains(test.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100) )
				}
				keeper.coinsFromStakedToUnstaked(context, test.validator)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100) )
				keeper.coinsFromStakedToUnstaked(context, test.validator)
				unstaked := keeper.GetUnstakedTokens(context)
				assert.Equal(t, test.amount, unstaked, "values do not match")
			}
		})
	}
}

func TestBurnStakedTokens(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	supplySize := sdk.NewInt(100)
	if err != nil {
		panic(err)
	}
	tests := []struct{
		name string
		expected string
		validator types.Validator
		burnAmount sdk.Int
		amount sdk.Int
		errs bool
	}{
		//{
		//	name: "burn coins from pool",
		//	validator: types.Validator{Address: validatorAddress},
		//	burnAmount: sdk.NewInt(5),
		//	amount: sdk.NewInt(10),
		//	errs: false,
		//},
		{
			name: "burn coins from pool",
			validator: types.Validator{Address: validatorAddress},
			burnAmount: sdk.NewInt(-1),
			amount: sdk.NewInt(10),
			errs: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			switch test.errs{
			case true:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, supplySize)
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				err := keeper.burnStakedTokens(context, test.burnAmount)
				assert.Nil(t, err,"error is not nil")
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, supplySize )
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				err := keeper.burnStakedTokens(context, test.burnAmount)
				if err != nil {
					t.Fail()
				}
				staked := keeper.GetStakedTokens(context)
				assert.Equal(t, test.amount.Sub(test.burnAmount).Add((supplySize)), staked, "values do not match")
			}
		})
	}
}
