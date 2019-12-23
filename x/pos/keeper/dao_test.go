package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	types "github.com/pokt-network/posmint/x/pos/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinsFromDAOToValidator(t *testing.T) {
	initialPower := int64(100)
	nAccs := int64(4)
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name               string
		amount sdk.Int
		expected string
		address sdk.ValAddress
		panics bool
	}{
		{
			name:               "sends coin to account",
			amount: sdk.NewInt(90),
			expected: fmt.Sprintf("was successfully minted to %s", validatorAddress.String()),
			address: validatorAddress,
			panics: false,
		},
		{
			name:               "panics invalid ammount of coins",
			amount: sdk.NewInt(-1),
			expected: fmt.Sprintf("negative coin amount: -1"),
			address: validatorAddress,
			panics: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true, initialPower, nAccs)

			switch test.panics{
			case true:
				defer func () {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected, "error does not match")
				}()
				keeper.coinsFromDAOToValidator(context, types.Validator{Address: test.address}, test.amount)
			default:
				AddMintedCoins(t, context, &keeper)
				keeper.coinsFromDAOToValidator(context, types.Validator{Address: test.address}, test.amount)
				coins := keeper.coinKeeper.GetCoins(context, sdk.AccAddress(test.address))
				assert.Equal(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)), coins, "coins should match")
			}
		})
	}
}
func AddMintedCoins(t *testing.T, ctx sdk.Context, k *Keeper) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), sdk.NewInt(100)))
	mintErr := k.supplyKeeper.MintCoins(ctx, types.DAOPoolName, coins.Add(coins))
	if mintErr != nil {
		t.Fail()
	}
}
