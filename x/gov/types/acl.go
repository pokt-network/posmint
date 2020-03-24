package types

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
)

type ACL interface {
	Validate(adjacencyMap map[string]bool) error
	GetOwner(permKey string) sdk.Address
	SetOwner(permKey string, ownerValue sdk.Address)
	GetAll() map[string]sdk.Address
	String() string
}

const (
	ACLKeySep = "/"
)

func NewACLKey(subspaceName, paramName string) string {
	return subspaceName + ACLKeySep + paramName
}

func SplitACLKey(aclKey string) (subspaceName, paramName string) {
	s := strings.Split(aclKey, ACLKeySep)
	subspaceName = s[0]
	paramName = s[1]
	return
}

var _ ACL = BaseACL{}

type BaseACL struct {
	M map[string]sdk.Address
}

func (b BaseACL) Validate(adjacencyMap map[string]bool) error {
	for key, val := range b.M {
		_, ok := adjacencyMap[key]
		if !ok {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the key: %s is not a recognized parameter", key))
		}
		adjacencyMap[key] = true
		if val == nil {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the address provided for: %s is nil", key))
		}
	}
	var unOwnedParams []string
	for key, val := range adjacencyMap {
		if !val {
			unOwnedParams = append(unOwnedParams, key)
		}
	}
	if len(unOwnedParams) != 0 {
		return ErrInvalidACL(ModuleName, fmt.Errorf("the following params have no owner: %v", unOwnedParams))
	}
	return nil
}

func (b BaseACL) GetOwner(permKey string) sdk.Address {
	return b.M[permKey]
}

func (b BaseACL) SetOwner(permKey string, ownerValue sdk.Address) {
	b.M[permKey] = ownerValue
}

func (b BaseACL) GetAll() map[string]sdk.Address {
	return b.M
}

func (b BaseACL) String() string {
	return fmt.Sprintf(`ACL:
%v`, b.M)
}

var _ ACL = &NonMapACL{}

type NonMapACL []ACLPair // cant use map cause of amino concrete marshal in tx

func (b NonMapACL) Validate(adjacencyMap map[string]bool) error {
	for _, aclPair := range b {
		key := aclPair.Key
		val := aclPair.Addr
		_, ok := adjacencyMap[key]
		if !ok {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the key: %s is not a recognized parameter", key))
		}
		adjacencyMap[key] = true
		if val == nil {
			return ErrInvalidACL(ModuleName, fmt.Errorf("the address provided for: %s is nil", key))
		}
	}
	var unOwnedParams []string
	for key, val := range adjacencyMap {
		if !val {
			unOwnedParams = append(unOwnedParams, key)
		}
	}
	if len(unOwnedParams) != 0 {
		return ErrInvalidACL(ModuleName, fmt.Errorf("the following params have no owner: %v aclPair", unOwnedParams))
	}
	return nil
}

func (b NonMapACL) GetOwner(permKey string) sdk.Address {
	for _, aclPair := range b {
		if aclPair.Key == permKey {
			return aclPair.Addr
		}
	}
	return nil
}

func (b *NonMapACL) SetOwner(permKey string, ownerValue sdk.Address) {
	for i, aclPair := range *b {
		if aclPair.Key == permKey {
			aclPair.Addr = ownerValue
			(*b)[i] = aclPair
			return
		}
	}
	temp := append(*b, ACLPair{
		Key:  permKey,
		Addr: ownerValue,
	})
	*b = temp
}

func (b NonMapACL) GetAll() map[string]sdk.Address {
	m := make(map[string]sdk.Address)
	for _, aclPair := range b {
		m[aclPair.Key] = aclPair.Addr
	}
	return m
}

func (b NonMapACL) String() string {
	return fmt.Sprintf(`ACL:
%v`, b.GetAll())
}

type ACLPair struct {
	Key  string      `json:"acl_key"`
	Addr sdk.Address `json:"address"`
}
