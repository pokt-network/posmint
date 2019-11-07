package keeper

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// blockKeeper handles access/modifiers of blocks
type blockKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// newBlockKeeper creates a new instance of block keeper
func newBlockKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) blockKeeper {
	return blockKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (bk blockKeeper) GetLatestBlockID(ctx sdk.Context) abci.BlockID {
	// return fixtures.GenerateBlockHash()
	header := ctx.BlockHeader()
	return header.GetLastBlockId()
}
