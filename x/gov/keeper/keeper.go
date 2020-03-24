package keeper

import (
	"fmt"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the global paramstore
type Keeper struct {
	cdc          *codec.Codec
	key          sdk.StoreKey
	tkey         sdk.StoreKey
	codespace    sdk.CodespaceType
	paramstore   sdk.Subspace
	SupplyKeeper types.SupplyKeeper
	spaces       map[string]sdk.Subspace
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey, codespace sdk.CodespaceType, supplyKeeper types.SupplyKeeper, subspaces ...sdk.Subspace) (k Keeper) {
	k = Keeper{
		cdc:          cdc,
		key:          key,
		tkey:         tkey,
		codespace:    codespace,
		SupplyKeeper: supplyKeeper,
		spaces:       make(map[string]sdk.Subspace),
	}
	k.paramstore = sdk.NewSubspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())
	k.paramstore.SetCodec(k.cdc)
	subspaces = append(subspaces, k.paramstore)
	k.AddSubspaces(subspaces...)
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
