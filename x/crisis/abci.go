package crisis

import (
	sdk "github.com/pokt-network/posmint/types"
)

// check all registered invariants
func EndBlocker(ctx sdk.Ctx, k Keeper) {
	if k.InvCheckPeriod() == 0 || ctx.BlockHeight()%int64(k.InvCheckPeriod()) != 0 {
		// skip running the invariant check
		return
	}
	k.AssertInvariants(ctx)
}
