package crisis

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
)

func (am AppModule) InvariantBroken(cdc *codec.Codec, address sdk.ValAddress)error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.node).WithCodec(cdc)

	senderAddr := cliCtx.GetFromAddress()
	moduleName, route := args[0], args[1] // todo
	msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
	return auth.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
