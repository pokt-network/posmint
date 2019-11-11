package pos

import (
	"github.com/pokt-network/posmint/client/context"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/keeper"
	"github.com/pokt-network/posmint/x/pos/types"
)

func StakeTx(cdc *codec.Codec, keeper keeper.Keeper, ctx sdk.Context, address sdk.ValAddress, amount sdk.Int) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	kb, err := cliCtx.Keybase.GetByAddress(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	msg := types.MsgStake{
		Address: address,
		PubKey:  kb.GetPubKey(), // needed for validator creation
		Value:   sdk.NewCoin(keeper.StakeDenom(ctx), amount),
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	msg := types.MsgBeginUnstake{Address: address}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func UnjailTx(cdc *codec.Codec, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	msg := types.MsgUnjail{ValidatorAddr: address}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func Send(cdc *codec.Codec, fromAddr, toAddr sdk.ValAddress, amount sdk.Int) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
