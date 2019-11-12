package pos

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/types"
)

func (am AppModule) StakeTx(cdc *codec.Codec, address sdk.ValAddress, amount sdk.Int) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.GetTendermintNode()).WithCodec(cdc).WithFromAddress(address)
	kb, err := cliCtx.Keybase.GetByAddress(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	msg := types.MsgStake{
		Address: address,
		PubKey:  kb.GetPubKey(), // needed for validator creation
		Value:   amount,         // todo convert to sdk.Int
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, name string, height int64, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.GetTendermintNode()).WithCodec(cdc).WithFromAddress(address)
	msg := types.MsgBeginUnstake{Address: address}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (am AppModule) UnjailTx(cdc *codec.Codec, name string, height int64, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.GetTendermintNode()).WithCodec(cdc).WithFromAddress(address)
	msg := types.MsgUnjail{ValidatorAddr: address}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (am AppModule) Send(cdc *codec.Codec, fromAddr, toAddr sdk.ValAddress, amount sdk.Int, name string, height int64) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.GetTendermintNode()).WithCodec(cdc).WithFromAddress(fromAddr)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
