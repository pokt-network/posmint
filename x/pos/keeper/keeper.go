package keeper

import (
	"container/list"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/tendermint/libs/log"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                *codec.Codec
	authKeeper         types.AuthKeeper
	hooks              types.POSHooks
	Paramstore         sdk.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, authKeeper types.AuthKeeper, paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure staked module accounts are set
	if addr := authKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.StakedPoolName))
	}
	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		authKeeper:         authKeeper,
		Paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		hooks:              nil,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
		codespace:          codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.POSHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
