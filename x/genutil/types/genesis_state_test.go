package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/pokt-network/posmint/types"
	authtypes "github.com/pokt-network/posmint/x/auth/types"
	stakingtypes "github.com/pokt-network/posmint/x/pos/types"
)

var (
	pk1 = ed25519.GenPrivKey().PubKey()
	pk2 = ed25519.GenPrivKey().PubKey()
)

func TestValidateGenesisMultipleMessages(t *testing.T) {

	msg1 := stakingtypes.MsgStake{
		Address: sdk.ValAddress(pk1.Address()),
		PubKey:  pk1,
		Value:   sdk.OneInt(),
	}

	msg2 := stakingtypes.MsgStake{
		Address: sdk.ValAddress(pk2.Address()),
		PubKey:  pk2,
		Value:   sdk.OneInt(),
	}

	genTxs := authtypes.NewStdTx([]sdk.Msg{msg1, msg2}, authtypes.StdFee{}, nil, "")
	genesisState := NewGenesisStateFromStdTx([]authtypes.StdTx{genTxs})

	err := ValidateGenesis(genesisState)
	require.Error(t, err)
}

func TestValidateGenesisBadMessage(t *testing.T) {
	msg1 := stakingtypes.MsgStake{
		Address: sdk.ValAddress(pk1.Address()),
		PubKey:  pk1,
		Value:   sdk.ZeroInt(),
	}

	genTxs := authtypes.NewStdTx([]sdk.Msg{msg1}, authtypes.StdFee{}, nil, "")
	genesisState := NewGenesisStateFromStdTx([]authtypes.StdTx{genTxs})

	err := ValidateGenesis(genesisState)
	require.Error(t, err)
}
