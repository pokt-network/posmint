package crypto

import "testing"

func getRandomPrivateKey(t *testing.T) Ed25519PrivateKey {
	return Ed25519PrivateKey{}.GenPrivateKey().(Ed25519PrivateKey)
}

func getRandomPubKey(t *testing.T) Ed25519PublicKey {
	pk := getRandomPrivateKey(t)
	return pk.PublicKey().(Ed25519PublicKey)
}
