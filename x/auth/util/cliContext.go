package util

import (
	"errors"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/node"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// CLIContext implements a typical CLI context created in SDK modules for
// transaction handling and queries.
type CLIContext struct {
	Codec         *codec.Codec
	Client        rpcclient.Client
	FromAddress   sdk.AccAddress
	Passphrase    string
	Height        int64
	BroadcastMode BroadcastType
}

// NewCLIContext returns a new initialized CLIContext with parameters from the
// command line using Viper. It takes a key name or address and populates the FromName and
// FromAddress field accordingly.
func NewCLIContext(node *node.Node, fromAddress sdk.AccAddress, passphrase string) CLIContext {
	return CLIContext{
		Client:      rpcclient.NewLocal(node),
		Passphrase:  passphrase,
		FromAddress: fromAddress,
	}
}

// WithCodec returns a copy of the context with an updated codec.
func (ctx CLIContext) WithCodec(cdc *codec.Codec) CLIContext {
	ctx.Codec = cdc
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx CLIContext) WithClient(client rpcclient.Client) CLIContext {
	ctx.Client = client
	return ctx
}

// WithFromAddress returns a copy of the context with an updated from account
// address.
func (ctx CLIContext) WithFromAddress(addr sdk.ValAddress) CLIContext {
	ctx.FromAddress = sdk.AccAddress(addr)
	return ctx
}

// WithHeight returns a copy of the context with an updated height.
func (ctx CLIContext) WithHeight(height int64) CLIContext {
	ctx.Height = height
	return ctx
}

// GetFromAddress returns the from address from the context's name.
func (ctx CLIContext) GetFromAddress() sdk.AccAddress {
	return ctx.FromAddress
}

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (ctx CLIContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("no client defined")
	}

	return ctx.Client, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// BroadcastTx broadcasts a transactions either synchronously or asynchronously
// based on the context parameters. The result of the broadcast is parsed into
// an intermediate structure which is logged if the context has a logger
// defined.
func (ctx CLIContext) BroadcastTx(txBytes []byte) (res sdk.TxResponse, err error) {
	switch ctx.BroadcastMode {
	case BroadcastSync:
		res, err = ctx.BroadcastTxSync(txBytes)

	case BroadcastAsync:
		res, err = ctx.BroadcastTxAsync(txBytes)

	case BroadcastBlock:
		res, err = ctx.BroadcastTxCommit(txBytes)

	default:
		return sdk.TxResponse{}, fmt.Errorf("unsupported return type %s; supported types: sync, async, block", ctx.BroadcastMode)
	}

	return res, err
}

// BroadcastTxCommit broadcasts transaction bytes to a Tendermint node and
// waits for a commit. An error is only returned if there is no RPC node
// connection or if broadcasting fails.
//
// NOTE: This should ideally not be used as the request may timeout but the tx
// may still be included in a block. Use BroadcastTxAsync or BroadcastTxSync
// instead.
func (ctx CLIContext) BroadcastTxCommit(txBytes []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	res, err := node.BroadcastTxCommit(txBytes)
	if err != nil {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}

	if !res.CheckTx.IsOK() {
		return sdk.NewResponseFormatBroadcastTxCommit(res), nil
	}

	if !res.DeliverTx.IsOK() {
		return sdk.NewResponseFormatBroadcastTxCommit(res), nil
	}

	return sdk.NewResponseFormatBroadcastTxCommit(res), nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously (i.e. returns after CheckTx execution).
func (ctx CLIContext) BroadcastTxSync(txBytes []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	res, err := node.BroadcastTxSync(txBytes)
	return sdk.NewResponseFormatBroadcastTx(res), err
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously (i.e. returns immediately).
func (ctx CLIContext) BroadcastTxAsync(txBytes []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	res, err := node.BroadcastTxAsync(txBytes)
	return sdk.NewResponseFormatBroadcastTx(res), err
}

type BroadcastType int

const (
	BroadcastSync BroadcastType = iota + 1
	BroadcastAsync
	BroadcastBlock
)

// ---------------------------------------------------------------------------------------------------------------------
// Query performs a query to a Tendermint node with the provided path.
// It returns the result and height of the query upon success or an error if
// the query fails.
func (ctx CLIContext) Query(path string) ([]byte, int64, error) {
	return ctx.query(path, nil)
}

// QueryWithData performs a query to a Tendermint node with the provided path
// and a data payload. It returns the result and height of the query upon success
// or an error if the query fails.
func (ctx CLIContext) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	return ctx.query(path, data)
}

// QueryStore performs a query to a Tendermint node with the provided key and
// store name. It returns the result and height of the query upon success
// or an error if the query fails.
func (ctx CLIContext) QueryStore(key cmn.HexBytes, storeName string) ([]byte, int64, error) {
	return ctx.queryStore(key, storeName, "key")
}

// QuerySubspace performs a query to a Tendermint node with the provided
// store name and subspace. It returns key value pair and height of the query
// upon success or an error if the query fails.
func (ctx CLIContext) QuerySubspace(subspace []byte, storeName string) (res []sdk.KVPair, height int64, err error) {
	resRaw, height, err := ctx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, height, err
	}

	ctx.Codec.MustUnmarshalBinaryLengthPrefixed(resRaw, &res)
	return
}

// query performs a query to a Tendermint node with the provided store name
// and path. It returns the result and height of the query upon success
// or an error if the query fails. If query height is invalid, an error will be returned.
func (ctx CLIContext) query(path string, key cmn.HexBytes) (res []byte, height int64, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
		Prove:  false,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, height, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, resp.Height, errors.New(resp.Log)
	}

	return resp.Value, resp.Height, nil
}

// queryStore performs a query to a Tendermint node with the provided a store
// name and path. It returns the result and height of the query upon success
// or an error if the query fails.
func (ctx CLIContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, int64, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}
