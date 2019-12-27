package crisis

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
	"github.com/tendermint/tendermint/rpc/client"
)

func InvariantBroken(cdc *codec.Codec, tmNode client.Client, moduleName, route string, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, sdk.AccAddress(address), passphrase).WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg := types.NewMsgVerifyInvariant(senderAddr, moduleName, route)
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
