package crisis

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
)

func (am AppModule) InvariantBroken(cdc *codec.Codec, moduleName, route string, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.node, sdk.AccAddress(address), passphrase).WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
