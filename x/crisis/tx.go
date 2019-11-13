package crisis

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
)

func (am AppModule) InvariantBroken(cdc *codec.Codec, moduleName, route string, address sdk.ValAddress) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(util.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext(am.node).WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
	return util.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
