package keeper

import (
	"github.com/pokt-network/posmint/x/bank"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifyParam(t *testing.T) {
	var aclKey = "bank/sendenabled"
	ctx, k := createTestKeeperAndContext(t, false)
	res := k.ModifyParam(ctx, aclKey, false, k.GetACL(ctx).GetOwner(aclKey))
	assert.Zero(t, res.Code)
	s, ok := k.GetSubspace(bank.DefaultParamspace)
	assert.True(t, ok)
	var b bool
	b = true
	s.Get(ctx, []byte("sendenabled"), &b)
	assert.False(t, b)
}
