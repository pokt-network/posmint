package types

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"math/rand"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	return pub
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

var testACL ACL

func createTestACL() ACL {
	if testACL == nil {
		acl := BaseACL{M: make(map[string]sdk.Address)}
		acl.SetOwner("bank/sendenabled", getRandomValidatorAddress())
		acl.SetOwner("auth/MaxMemoCharacters", getRandomValidatorAddress())
		acl.SetOwner("auth/TxSigLimit", getRandomValidatorAddress())
		acl.SetOwner("auth/TxSizeCostPerByte", getRandomValidatorAddress())
		acl.SetOwner("gov/daoOwner", getRandomValidatorAddress())
		acl.SetOwner("gov/acl", getRandomValidatorAddress())
		acl.SetOwner("gov/upgrade", getRandomValidatorAddress())
		testACL = acl
	}
	return testACL
}

func createTestAdjacencyMap() map[string]bool {
	m := make(map[string]bool)
	m["bank/sendenabled"] = true       // set
	m["auth/MaxMemoCharacters"] = true // set
	m["auth/TxSigLimit"] = true        // set
	m["auth/TxSizeCostPerByte"] = true // set
	m["gov/daoOwner"] = true           // set
	m["gov/acl"] = true                // set
	m["gov/upgrade"] = true
	return m
}
