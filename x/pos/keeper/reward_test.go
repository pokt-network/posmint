package keeper

import (
	"encoding/hex"
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
	}{
		{
			name:   "mints a coin",
			amount: sdk.NewInt(90),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			result := keeper.mint(context, test.amount, validatorAddress)
			assert.Contains(t, result.Log, "was successfully minted", "does not contain message")
		})
	}
}
