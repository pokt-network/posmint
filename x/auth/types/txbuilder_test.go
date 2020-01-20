package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
)

func TestTxBuilderBuild(t *testing.T) {
	type fields struct {
		TxEncoder     sdk.TxEncoder
		AccountNumber uint64
		Sequence      uint64
		ChainID       string
		Memo          string
		Fees          sdk.Coins
	}
	defaultMsg := []sdk.Msg{sdk.NewTestMsg(addr)}
	tests := []struct {
		name    string
		fields  fields
		msgs    []sdk.Msg
		want    StdSignMsg
		wantErr bool
	}{
		{
			"builder with fees",
			fields{
				TxEncoder:     DefaultTxEncoder(codec.New()),
				AccountNumber: 1,
				Sequence:      1,
				ChainID:       "test-chain",
				Memo:          "hello from Voyager 1!",
				Fees:          sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))),
			},
			defaultMsg,
			StdSignMsg{
				ChainID:       "test-chain",
				AccountNumber: 1,
				Sequence:      1,
				Memo:          "hello from Voyager 1!",
				Msgs:          defaultMsg,
				Fee:           sdk.Coins{sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))},
			},
			false,
		},
		{
			"no chain-id supplied",
			fields{
				TxEncoder:     DefaultTxEncoder(codec.New()),
				AccountNumber: 1,
				Sequence:      1,
				ChainID:       "",
				Memo:          "hello from Voyager 1!",
				Fees:          sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))),
			},
			defaultMsg,
			StdSignMsg{
				ChainID:       "test-chain",
				AccountNumber: 1,
				Sequence:      1,
				Memo:          "hello from Voyager 1!",
				Msgs:          defaultMsg,
				Fee:           sdk.Coins{sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))},
			},
			true,
		},
		{
			"builder w/ fees",
			fields{
				TxEncoder:     DefaultTxEncoder(codec.New()),
				AccountNumber: 1,
				Sequence:      1,
				ChainID:       "test-chain",
				Memo:          "hello from Voyager 1!",
				Fees:          sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))),
			},
			defaultMsg,
			StdSignMsg{
				ChainID:       "test-chain",
				AccountNumber: 1,
				Sequence:      1,
				Memo:          "hello from Voyager 1!",
				Msgs:          defaultMsg,
				Fee:           sdk.Coins{sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1))},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bldr := NewTxBuilder(
				tt.fields.TxEncoder,
				tt.fields.AccountNumber,
				tt.fields.Sequence,
				tt.fields.ChainID,
				tt.fields.Memo,
				tt.fields.Fees,
			)
			got, err := bldr.BuildSignMsg(tt.msgs)
			require.Equal(t, tt.wantErr, (err != nil))
			if err == nil {
				require.True(t, reflect.DeepEqual(tt.want, got))
			}
		})
	}
}
