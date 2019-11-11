package rpc

import (
	"github.com/pokt-network/posmint/client/context"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func GetNodeStatus(cliCtx context.CLIContext) (*ctypes.ResultStatus, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}

	return node.Status()
}