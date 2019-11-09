package keeper

import (
	"github.com/pokt-network/posmint/client/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/pos/types"
)

func (k Keeper) StakeTx(ctx sdk.Context, address sdk.ValAddress, amount sdk.Int) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(k.cdc))
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	kb, err := cliCtx.Keybase.GetByAddress(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	msg := types.MsgStake{
		Address: address,
		PubKey:  kb.GetPubKey(), // needed for validator creation
		Value:   sdk.NewCoin(k.StakeDenom(ctx), amount),
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (k Keeper) UnstakeTx(address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(k.cdc))
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	msg := types.MsgBeginUnstake{Address: address,}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (k Keeper) UnjailTx(ctx sdk.Context, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(k.cdc))
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	msg := types.MsgUnjail{ValidatorAddr: address,}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func (k Keeper) Send(ctx sdk.Context, fromAddr, toAddr sdk.ValAddress, amount sdk.Int) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(k.cdc))
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
