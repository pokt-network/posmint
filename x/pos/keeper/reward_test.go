package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type args struct {
	amount      sdk.Int
	valAddress  sdk.ValAddress
	consAddress sdk.ConsAddress
}

func TestSetandGetValidatorAward(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		expectedFind  bool
	}{
		{
			name:          "can set award",
			expectedCoins: sdk.NewInt(1),
			expectedFind:  true,
			args:          args{amount: sdk.NewInt(int64(1)), valAddress: validatorAddress},
		},
		{
			name:          "can get award",
			expectedCoins: sdk.NewInt(2),
			expectedFind:  true,
			args:          args{amount: sdk.NewInt(int64(2)), valAddress: validatorAddress},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			keeper.setValidatorAward(context, test.args.amount, test.args.valAddress)
			coins, found := keeper.getValidatorAward(context, test.args.valAddress)
			assert.Equal(t, test.expectedCoins, coins, "coins don't match")
			assert.Equal(t, test.expectedFind, found, "finds don't match")

		})
	}
}

func TestSetAndGetProposer(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	consAddress, err := sdk.ConsAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name            string
		args            args
		expectedAddress sdk.ConsAddress
	}{
		{
			name:            "can set the preivous proposer",
			args:            args{consAddress: consAddress},
			expectedAddress: consAddress,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			keeper.SetPreviousProposer(context, test.args.consAddress)
			receivedAddress := keeper.GetPreviousProposer(context)
			assert.Equal(t, test.expectedAddress, receivedAddress, "addresses do not match ")
		})
	}
}

func TestDeleteValidatorAward(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		expectedFind  bool
	}{
		{
			name:          "can delete award",
			expectedCoins: sdk.NewInt(0),
			expectedFind:  false,
			args:          args{amount: sdk.NewInt(int64(1)), valAddress: validatorAddress},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			keeper.setValidatorAward(context, test.args.amount, test.args.valAddress)
			keeper.deleteValidatorAward(context, test.args.valAddress)
			_, found := keeper.getValidatorAward(context, test.args.valAddress)
			assert.Equal(t, test.expectedFind, found, "finds do not match")

		})
	}
}

func TestGetProposerRewardPercentage(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	tests := []struct {
		name               string
		expectedPercentage sdk.Int
	}{
		{
			name:               "get reward percentage",
			expectedPercentage: sdk.NewInt(90),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			percentage := keeper.getProposerRewardPercentage(context) // TODO: replace with  sdk.Dec isntead of sdk.Int
			assert.Equal(t, test.expectedPercentage, percentage, "percentages do not match")
		})
	}
}

func TestMint(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name   string
		amount sdk.Int
		expected string
		address  sdk.ValAddress
		panics   bool
	}{
		{
			name:   "mints a coin",
			amount: sdk.NewInt(90),
			expected: fmt.Sprintf("was successfully minted to %s", validatorAddress.String()),
			address:  validatorAddress,
			panics:   false,
		},
		{
			name:     "panics invalid ammount of coins",
			amount:   sdk.NewInt(-1),
			expected: fmt.Sprintf("negative coin amount: -1"),
			address:  validatorAddress,
			panics:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected, "error does not match")
				}()
				_ = keeper.mint(context, test.amount, test.address)
			default:
				result := keeper.mint(context, test.amount, test.address)
				assert.Contains(t, result.Log, test.expected, "does not contain message")
				coins := keeper.coinKeeper.GetCoins(context, sdk.AccAddress(test.address))
				assert.Equal(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)), coins, "coins should match")
			}
		})
	}
}

func TestMintValidatorAwards(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	for i := 0; i < 3; i++ {
		addressBytes = append(addressBytes, []byte("abcdefghijklmnopqrst")...)
	}
	addressBytes = append(addressBytes, []byte("a")...)
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name     string
		amount   sdk.Int
		expected string
		address  sdk.ValAddress
		panics   bool
	}{
		{
			name:     "mints a coin",
			amount:   sdk.NewInt(90),
			expected: fmt.Sprintf("was successfully minted to %s", validatorAddress.String()),
			address:  validatorAddress,
			panics:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)
			keeper.setValidatorAward(context, test.amount, test.address)

			keeper.mintValidatorAwards(context)
			coins := keeper.coinKeeper.GetCoins(context, sdk.AccAddress(test.address))
			assert.Equal(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)), coins, "coins should match")
		})
	}
}
