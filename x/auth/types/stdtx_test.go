package types

import (
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"gopkg.in/yaml.v2"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
)

var (
	priv = crypto.GenerateSecp256k1PrivKey()
	addr = sdk.Address(priv.PubKey().Address())
)

func TestStdTx(t *testing.T) {
	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
	sigs := []StdSignature{}
	fee := sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(5)))

	tx := NewStdTx(msgs, fee, sigs, "")
	require.Equal(t, msgs, tx.GetMsgs())
	require.Equal(t, sigs, tx.GetSignatures())

	feePayer := tx.GetSigners()[0]
	require.Equal(t, addr, feePayer)
}

func TestStdSignBytes(t *testing.T) {
	type args struct {
		chainID  string
		accnum   uint64
		sequence uint64
		fee      sdk.Coins
		msgs     []sdk.Msg
		memo     string
	}
	defaultFee := NewTestCoins()
	tests := []struct {
		args args
		want string
	}{
		{
			args{"1234", 3, 6, defaultFee, []sdk.Msg{sdk.NewTestMsg(addr)}, "memo"},
			fmt.Sprintf("{\"account_number\":\"3\",\"chain_id\":\"1234\",\"fee\":[{\"amount\":\"10000000\",\"denom\":\"atom\"}],\"memo\":\"memo\",\"msgs\":[[\"%s\"]],\"sequence\":\"6\"}", addr),
		},
	}
	for i, tc := range tests {
		got := string(StdSignBytes(tc.args.chainID, tc.args.accnum, tc.args.sequence, tc.args.fee, tc.args.msgs, tc.args.memo))
		require.Equal(t, tc.want, got, "Got unexpected result on test case i: %d", i)
	}
}

func TestTxValidateBasic(t *testing.T) {
	ctx := sdk.NewContext(nil, abci.Header{ChainID: "mychainid"}, false, log.NewNopLogger())

	// keys and addresses
	priv1, _, addr1 := KeyTestPubAddr()
	priv2, _, addr2 := KeyTestPubAddr()

	// msg and signatures
	msg1 := NewTestMsg(addr1, addr2)
	fee := NewTestCoins()

	msgs := []sdk.Msg{msg1}

	// require to fail validation upon invalid fee
	badFee := NewTestCoins()
	badFee[0].Amount = sdk.NewInt(-5)
	tx := NewTestTx(ctx, nil, nil, nil, nil, badFee)

	err := tx.ValidateBasic()
	require.Error(t, err)
	require.Equal(t, sdk.CodeInsufficientFee, err.Result().Code)

	// require to fail validation when no signatures exist
	privs, accNums, seqs := []crypto.PrivateKey{}, []uint64{}, []uint64{}
	tx = NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	err = tx.ValidateBasic()
	require.Error(t, err)
	require.Equal(t, sdk.CodeNoSignatures, err.Result().Code)

	// require to fail validation when signatures do not match expected signers
	privs, accNums, seqs = []crypto.PrivateKey{priv1}, []uint64{0, 1}, []uint64{0, 0}
	tx = NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	err = tx.ValidateBasic()
	require.Error(t, err)
	require.Equal(t, sdk.CodeUnauthorized, err.Result().Code)

	// require to pass when above criteria are matched
	privs, accNums, seqs = []crypto.PrivateKey{priv1, priv2}, []uint64{0, 1}, []uint64{0, 0}
	tx = NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	err = tx.ValidateBasic()
	require.NoError(t, err)
}

func TestDefaultTxEncoder(t *testing.T) {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "posmint/Test", nil)
	encoder := DefaultTxEncoder(cdc)

	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
	fee := NewTestCoins()
	sigs := []StdSignature{}

	tx := NewStdTx(msgs, fee, sigs, "")

	cdcBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)

	require.NoError(t, err)
	encoderBytes, err := encoder(tx)

	require.NoError(t, err)
	require.Equal(t, cdcBytes, encoderBytes)
}

func TestStdSignatureMarshalYAML(t *testing.T) {
	_, pubKey, _ := KeyTestPubAddr()

	testCases := []struct {
		sig    StdSignature
		output string
	}{
		{
			StdSignature{},
			"|\n  pubkey: \"\"\n  signature: \"\"\n",
		},
		{
			StdSignature{PublicKey: pubKey, Signature: []byte("dummySig")},
			fmt.Sprintf("|\n  pubkey: %s\n  signature: dummySig\n", pubKey.RawString()),
		},
		{
			StdSignature{PublicKey: pubKey, Signature: nil},
			fmt.Sprintf("|\n  pubkey: %s\n  signature: \"\"\n", pubKey.RawString()),
		},
	}

	for i, tc := range testCases {
		bz, err := yaml.Marshal(tc.sig)
		require.NoError(t, err)
		require.Equal(t, tc.output, string(bz), "test case #%d", i)
	}
}
