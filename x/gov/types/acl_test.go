package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestACLGetSetOwner(t *testing.T) {
	acl := BaseACL{M: make(map[string]sdk.Address)}
	a := getRandomValidatorAddress()
	acl.SetOwner("bank/sendenabled", a)
	assert.Equal(t, acl.GetOwner("bank/sendenabled").String(), a.String())
}

func TestValidateACL(t *testing.T) {
	acl := createTestACL()
	adjMap := createTestAdjacencyMap()
	assert.Nil(t, acl.Validate(adjMap))
	acl.SetOwner("bank/sendenabled2", getRandomValidatorAddress())
	assert.NotNil(t, acl.Validate(adjMap))
}
