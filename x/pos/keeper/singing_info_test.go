package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMissedArray(t *testing.T) {
	validator := getStakedValidator()
	address := validator.GetAddress()

	tests := []struct {
		name     string
		expected bool
		address  sdk.Address
	}{
		{
			name:     "gets missed block array",
			address:  address,
			expected: true,
		},
		{
			name:     "gets missed block array",
			address:  address,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetMissedBlockArray(context, test.address, 1, test.expected)
			missed := keeper.getMissedBlockArray(context, test.address, 1)
			assert.Equal(t, missed, test.expected, "found does not match")
		})
	}
}

func TestClearMissedArray(t *testing.T) {
	validator := getStakedValidator()
	address := validator.GetAddress()

	tests := []struct {
		name     string
		expected bool
		address  sdk.Address
	}{
		{
			name:     "gets missed block array",
			address:  address,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetMissedBlockArray(context, test.address, 1, true)
			keeper.clearMissedArray(context, test.address)
			missed := keeper.getMissedBlockArray(context, test.address, 1)
			assert.Equal(t, missed, test.expected, "found does not match")
		})
	}
}
