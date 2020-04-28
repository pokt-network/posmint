package gov

import (
	"fmt"

	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/gov/types"
	"github.com/tendermint/tendermint/rpc/client"
)

func ChangeParamsTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddress sdk.Address, aclKey string, paramValue interface{}, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgChangeParam{
		FromAddress: fromAddress,
		ParamKey:    aclKey,
		ParamVal:    paramValue,
	}
	txBuilder, cliCtx := newTx(cdc, msg, fromAddress, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func DAOTransferTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddress, toAddress sdk.Address, amount sdk.Int, action, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgDAOTransfer{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
		Action:      action,
	}
	txBuilder, cliCtx := newTx(cdc, msg, fromAddress, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UpgradeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddress sdk.Address, upgrade types.Upgrade, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgUpgrade{
		Address: fromAddress,
		Upgrade: upgrade,
	}
	txBuilder, cliCtx := newTx(cdc, msg, fromAddress, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func newTx(cdc *codec.Codec, msg sdk.Msg, fromAddr sdk.Address, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	// Retrieve chain ID
	chainID := genDoc.Genesis.ChainID

	cliCtx = util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	cliCtx.BroadcastMode = util.BroadcastSync

	// retrieve private key from keypair
	account, err := cliCtx.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	fee := sdk.NewInt(types.GovFeeMap[msg.Type()])
	if account.GetCoins().AmountOf(sdk.DefaultStakeDenom).LTE(fee) { // todo get stake denom
		panic(fmt.Sprintf("insufficient funds: the fee needed is %v", fee))
	}
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, fee))).WithKeybase(keybase)
	return
}

func BuildAndSignMulti(cdc *codec.Codec, address sdk.Address, publicKey crypto.PublicKeyMultiSig, msg sdk.Msg, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBytes []byte, err error) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		genDoc.Genesis.ChainID,
		"", nil).WithKeybase(keybase)
	return txBuilder.BuildAndSignMultisigTransaction(address, publicKey, msg, passphrase)
}

func SignMulti(cdc *codec.Codec, fromAddr sdk.Address, tx []byte, keys []crypto.PublicKey, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBytes []byte, err error) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		genDoc.Genesis.ChainID,
		"", nil).WithKeybase(keybase)
	return txBuilder.SignMultisigTransaction(fromAddr, keys, passphrase, tx)
}
