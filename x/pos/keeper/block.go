package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) GetLatestBlockID(ctx sdk.Context) abci.BlockID {
	// return fixtures.GenerateBlockHash()
	header := ctx.BlockHeader()
	return header.GetLastBlockId()
}
