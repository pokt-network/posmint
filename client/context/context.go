package context

import (
	"github.com/pokt-network/posmint/codec"
	cryptokeys "github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	tndmt "github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// CLIContext implements a typical CLI context created in SDK modules for
// transaction handling and queries.
type CLIContext struct {
	Codec         *codec.Codec
	Client        rpcclient.Client
	Keybase       cryptokeys.Keybase
	From          string
	FromAddress   sdk.AccAddress
	FromName      string
	Height        int64
	BroadcastMode BroadcastType
}

// NewCLIContext returns a new initialized CLIContext with parameters from the
// command line using Viper. It takes a key name or address and populates the FromName and
// FromAddress field accordingly.
func NewCLIContext(fromAddress sdk.AccAddress, fromName string, height int64, node *tndmt.Node) CLIContext { // todo figure out node input
	var rpc rpcclient.Client
	rpc = rpcclient.NewLocal(node)
	return CLIContext{
		Client:      rpc,
		FromAddress: fromAddress,
		FromName:    fromName,
		Height:      height,
	}
}

// WithCodec returns a copy of the context with an updated codec.
func (ctx CLIContext) WithCodec(cdc *codec.Codec) CLIContext {
	ctx.Codec = cdc
	return ctx
}

// WithFrom returns a copy of the context with an updated from address or name.
func (ctx CLIContext) WithFrom(from string) CLIContext {
	ctx.From = from
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx CLIContext) WithClient(client rpcclient.Client) CLIContext {
	ctx.Client = client
	return ctx
}

// WithFromName returns a copy of the context with an updated from account name.
func (ctx CLIContext) WithFromName(name string) CLIContext {
	ctx.FromName = name
	return ctx
}

// WithFromAddress returns a copy of the context with an updated from account
// address.
func (ctx CLIContext) WithFromAddress(addr sdk.AccAddress) CLIContext {
	ctx.FromAddress = addr
	return ctx
}

// WithHeight returns a copy of the context with an updated height.
func (ctx CLIContext) WithHeight(height int64) CLIContext {
	ctx.Height = height
	return ctx
}
