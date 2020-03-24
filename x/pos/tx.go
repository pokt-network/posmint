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
	txBuilder, cliCtx, err := newTx(cdc, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	msg := types.MsgStake{
		PubKey: kp.PublicKey,
		Value:  amount,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx, err := newTx(cdc, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	msg := types.MsgBeginUnstake{Address: address}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnjailTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx, err := newTx(cdc, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	msg := types.MsgUnjail{ValidatorAddr: address}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func Send(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddr, toAddr sdk.Address, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	txBuilder, cliCtx, err := newTx(cdc, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func RawTx(cdc *codec.Codec, tmNode client.Client, fromAddr sdk.Address, txBytes []byte) (sdk.TxResponse, error) {
	return util.CLIContext{
		Codec:         cdc,
		Client:        tmNode,
		FromAddress:   fromAddr,
		BroadcastMode: util.BroadcastSync,
	}.BroadcastTx(txBytes)
}

func newTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, passphrase string) (auth.TxBuilder, util.CLIContext, error) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		return auth.TxBuilder{}, util.CLIContext{}, err
	}
	kp, err := keybase.List()
	if err != nil {
		return auth.TxBuilder{}, util.CLIContext{}, err
	}
	chainID := genDoc.Genesis.ChainID
	fromAddr := kp[0].GetAddress()
	cliCtx := util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	accGetter := auth.NewAccountRetriever(cliCtx)
	err = accGetter.EnsureExists(fromAddr)
	if err != nil {
		return auth.TxBuilder{}, cliCtx, err
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin("upokt", sdk.NewInt(10)))).WithKeybase(keybase) // todo get stake denom
	return txBuilder, cliCtx, nil
}
