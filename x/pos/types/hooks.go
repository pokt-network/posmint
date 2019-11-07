package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// combine multiple staking hooks, all hook functions are run in array sequence
type MultiPOSHooks []POSHooks

func NewMultiStakingHooks(hooks ...POSHooks) MultiPOSHooks {
	return hooks
}

func (h MultiPOSHooks) BeforeValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeValidatorCreated(ctx, valAddr)
	}
}

// nolint
func (h MultiPOSHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterValidatorCreated(ctx, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidtorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeValidatorRemoved(ctx, consAddr, address)
	}
}

func (h MultiPOSHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorStaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeValidatorStaked(ctx, consAddr, address)
	}
}

func (h MultiPOSHooks) AfterValidatorStaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterValidatorStaked(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeValidatorBeginUnstaking(ctx, consAddr, address)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaking(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterValidatorBeginUnstaking(ctx, consAddr, valAddr)
	}
}

func (h MultiPOSHooks) BeforeValidatorBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, address sdk.ValAddress) {
	for i := range h {
		h[i].BeforeValidatorUnstaked(ctx, consAddr, address)
	}
}
func (h MultiPOSHooks) AfterValidatorBeginUnstaked(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterValidatorUnstaked(ctx, consAddr, valAddr)
	}
}
func (h MultiPOSHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}

func (h MultiPOSHooks) AfterValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	for i := range h {
		h[i].AfterValidatorSlashed(ctx, valAddr, fraction)
	}
}
