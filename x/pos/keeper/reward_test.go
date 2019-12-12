package keeper

import (
	"github.com/pokt-network/posmint/testsUtil"
	"encoding/hex"
	"github.com/pokt-network/posmint/baseapp"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)


var capKey = sdk.NewKVStoreKey(baseapp.MainStoreKey)

type args struct {
	ctx     sdk.Context
	amount  sdk.Int
	address sdk.ValAddress
}

func TestSetValidatorAward(t *testing.T) {
	addressBytes := []byte("abcdefghijklmnopqrst")
	validatorAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		panic(err)
	}
	minGasPrices := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 5000)}
	options := baseapp.SetMinGasPrices(minGasPrices.String())
	cdc := testUtil.MakeCodec()
	bapp := testUtil.GetNewApp(capKey, cdc, options)

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		expectedFind bool
		app           *baseapp.BaseApp
		keeper        Keeper
	}{
		{
			name:          "can set Value",
			expectedCoins: sdk.NewInt(1),
			expectedFind:	true,
			app:           bapp,
			args:          args{ctx: testUtil.GetNewContext(bapp), amount: sdk.NewInt(int64(1)), address: validatorAddress},
			keeper:        getNewKeeper(capKey, cdc),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.keeper.setValidatorAward(test.args.ctx, test.args.amount, test.args.address)
			coins, found := test.keeper.getValidatorAward(test.args.ctx, test.args.address)
			assert.Equal(t, test.expectedCoins, coins, "coins don't match")
			assert.Equal(t, test.expectedFind, found, "finds don't match")

		})
	}
}

func getNewKeeper(key sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{storeKey: key, cdc: cdc}
}
