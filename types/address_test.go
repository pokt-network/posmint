package types_test

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/pokt-network/posmint/types"
)

var invalidStrs = []string{
	"hello, world!",
	"0xAA",
	"AAA",
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

func TestEmptyAddresses(t *testing.T) {
	require.Equal(t, (types.Address{}).String(), "")
	require.Equal(t, (types.Address{}).String(), "")
	require.Equal(t, (types.Address{}).String(), "")

	accAddr, err := types.AddressFromHex("")
	require.True(t, accAddr.Empty())
	require.Nil(t, err)

	valAddr, err := types.AddressFromHex("")
	require.True(t, valAddr.Empty())
	require.Nil(t, err)

	consAddr, err := types.AddressFromHex("")
	require.True(t, consAddr.Empty())
	require.Nil(t, err)
}

func TestYAMLMarshalers(t *testing.T) {
	addr := secp256k1.GenPrivKey().PubKey().Address()

	acc := types.Address(addr)
	val := types.Address(addr)
	cons := types.Address(addr)

	got, _ := yaml.Marshal(&acc)
	require.Equal(t, acc.String()+"\n", string(got))

	got, _ = yaml.Marshal(&val)
	require.Equal(t, val.String()+"\n", string(got))

	got, _ = yaml.Marshal(&cons)
	require.Equal(t, cons.String()+"\n", string(got))
}

func TestRandBech32AccAddrConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		rand.Read(pub[:])

		acc := types.Address(pub.Address())
		res := types.Address{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AddressFromHex(str)
		require.NotNil(t, err)

		err = (*types.Address)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestValAddr(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.Address(pub.Address())
		res := types.Address{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AddressFromHex(str)
		require.NotNil(t, err)

		err = (*types.Address)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestAddress(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.Address(pub.Address())
		res := types.Address{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AddressFromHex(str)
		require.NotNil(t, err)

		err = (*types.Address)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestAddressInterface(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	addrs := []types.AddressI{
		types.Address(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case types.Address:
			_, err := types.AddressFromHex(addr.String())
			require.Nil(t, err)
		default:
			t.Fail()
		}
	}

}

func TestPubKeyInterfaceAssertion(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	var pub2 secp256k1.PubKeySecp256k1
	rand.Read(pub2[:])

	values := []crypto.PubKey{
		pub, pub2,
	}

	for _, v := range values {
		switch v := v.(type) {
		case ed25519.PubKeyEd25519:
			fmt.Println(v)
			s := types.HexAddressPubKey(v)
			as := types.HexAddressPubKeyAmino(v)
			fmt.Println(s)
			fmt.Println(as)
			require.NotNil(t, s)
			require.NotEqual(t, s, as)

		case secp256k1.PubKeySecp256k1:
			fmt.Println(v)
			s := types.HexAddressPubKey(v)
			as := types.HexAddressPubKeyAmino(v)
			fmt.Println(s)
			fmt.Println(as)
			require.NotNil(t, s)
			require.NotEqual(t, s, as)
		default:
			t.Fail()
		}
	}

}

func TestCustomAddressVerifier(t *testing.T) {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	accBech := types.Address(addr).String()
	valBech := types.Address(addr).String()
	consBech := types.Address(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := types.VerifyAddressFormat(addr)
	require.NotNil(t, err)
	_, err = types.AddressFromHex(accBech)
	require.NotNil(t, err)
	_, err = types.AddressFromHex(valBech)
	require.NotNil(t, err)
	_, err = types.AddressFromHex(consBech)
	require.NotNil(t, err)

	// Set a custom address verifier that accepts 10 or 20 byte addresses
	types.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == 10 || n == types.AddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic accepts this 10 byte address
	err = types.VerifyAddressFormat(addr)
	require.Nil(t, err)
	_, err = types.AddressFromHex(accBech)
	require.Nil(t, err)
	_, err = types.AddressFromHex(valBech)
	require.Nil(t, err)
	_, err = types.AddressFromHex(consBech)
	require.Nil(t, err)
}
