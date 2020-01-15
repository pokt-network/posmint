package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

type MultiPOSHooks []POSHooks

func (h MultiPOSHooks) BeforeValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorRegistered(ctx, valAddr)
	}
}

// nolint
func (h MultiPOSHooks) AfterValidatorRegistered(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorRegistered(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorRemoved(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorRemoved(ctx, valAddr)
	}
}

func (h MultiPOSHooks) AfterValidatorRemoved(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorRemoved(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorStaked(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorStaked(ctx, valAddr)
	}
}

func (h MultiPOSHooks) AfterValidatorStaked(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorStaked(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaking(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorBeginUnstaking(ctx, valAddr)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaking(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorBeginUnstaking(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaked(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].BeforeValidatorUnstaked(ctx, valAddr)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaked(ctx sdk.Context, valAddr sdk.Address) {
	for i := range h {
		h[i].AfterValidatorUnstaked(ctx, valAddr)
	}
}
func (h MultiPOSHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	for i := range h {
		h[i].BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}

func (h MultiPOSHooks) AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.Address, fraction sdk.Dec) {
	for i := range h {
		h[i].AfterValidatorSlashed(ctx, valAddr, fraction)
	}
}
