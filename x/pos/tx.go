package pos

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/tendermint/rpc/client"
)

func StakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, amount sdk.Int, kp keys.KeyPair, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx := newTx(cdc, tmNode, keybase, passphrase)
	msg := types.MsgStake{
		Address: sdk.Address(kp.GetAddress()),
		PubKey:  kp.PubKey,
		Value:   amount,
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx := newTx(cdc, tmNode, keybase, passphrase)
	msg := types.MsgBeginUnstake{Address: address}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnjailTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx := newTx(cdc, tmNode, keybase, passphrase)
	msg := types.MsgUnjail{ValidatorAddr: address}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func Send(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddr, toAddr sdk.Address, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	txBuilder, cliCtx := newTx(cdc, tmNode, keybase, passphrase)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func RawTx(cdc *codec.Codec, tmNode client.Client, fromAddr sdk.Address, txBytes []byte) (sdk.TxResponse, error) {
	return util.CLIContext{
		Codec:       cdc,
		Client:      tmNode,
		FromAddress: sdk.Address(fromAddr),
	}.BroadcastTx(txBytes)
}

func newTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	kp, err := keybase.List()
	if err != nil {
		panic(err)
	}
	chainID := genDoc.Genesis.ChainID
	fromAddr := kp[0].GetAddress()
	cliCtx = util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	accGetter := auth.NewAccountRetriever(cliCtx)
	err = accGetter.EnsureExists(fromAddr)
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		account.GetAccountNumber(),
		account.GetSequence(),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10)))).WithKeybase(keybase) // todo get stake denom
	return
}
