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

func StakeTx(cdc *codec.Codec, tmNode client.Client, txBuilder auth.TxBuilder, kp keys.KeyPair, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, kp.GetAddress(), passphrase).WithCodec(cdc)
	msg := types.MsgStake{
		Address: sdk.ValAddress(kp.GetAddress()),
		PubKey:  kp.PubKey,
		Value:   amount,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgBeginUnstake{Address: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnjailTx(cdc *codec.Codec, tmNode client.Client, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgUnjail{ValidatorAddr: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func Send(cdc *codec.Codec, tmNode client.Client, fromAddr, toAddr sdk.ValAddress, txBuilder auth.TxBuilder, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, sdk.AccAddress(fromAddr), passphrase).WithCodec(cdc)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
