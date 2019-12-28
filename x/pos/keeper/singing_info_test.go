package keeper

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMissedArray(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	consAddr, err := sdk.ConsAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name     string
		expected bool
		address  sdk.ConsAddress
	}{
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: true,
		},
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: false,
		},
	}

	for _, test := range tests {
		context, _, keeper := createTestInput(t, true, initialPower, nAccs)
		keeper.SetMissedBlockArray(context, test.address, 1, test.expected)
		missed := keeper.getMissedBlockArray(context, test.address, 1)
		assert.Equal(t, missed, test.expected, "found does not match")
	}
}

func TestClearMissedArray(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	consAddr, err := sdk.ConsAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name     string
		expected bool
		address  sdk.ConsAddress
	}{
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: false,
		},
	}

	for _, test := range tests {
		context, _, keeper := createTestInput(t, true, initialPower, nAccs)
		keeper.SetMissedBlockArray(context, test.address, 1, true)
		keeper.clearMissedArray(context, test.address)
		missed := keeper.getMissedBlockArray(context, test.address, 1)
		assert.Equal(t, missed, test.expected, "found does not match")
	}
}
