package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// Access control list of parameters
func (k Keeper) GetACL(ctx sdk.Ctx) (res types.ACL) {
	k.paramstore.Get(ctx, types.ACLKey, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		ACL:      k.GetACL(ctx),
		Upgrade:  k.GetUpgrade(ctx),
		DAOOwner: k.GetDAOOwner(ctx),
	}
}

// set the params
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	if params.ACL != nil {
		// if non map acl turn into map acl for efficiency
		acl, ok := params.ACL.(*types.NonMapACL)
		if ok {
			params.ACL = types.BaseACL{M: acl.GetAll()}
		}
	}
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) GetDAOOwner(ctx sdk.Ctx) (res sdk.Address) {
	k.paramstore.Get(ctx, types.DAOOwnerKey, &res)
	return
}

func (k Keeper) GetUpgrade(ctx sdk.Ctx) (res types.Upgrade) {
	k.paramstore.Get(ctx, types.UpgradeKey, &res)
	return
}
