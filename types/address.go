package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

const (
	// Constants defined here are the defaults value for address.
	// You can use the specific values for your project.
	// Add the follow lines to the `main()` of your server.
	//
	//	config := sdk.GetConfig()
	//	config.SetBech32PrefixForAccount(yourBech32PrefixAccAddr, yourBech32PrefixAccPub)
	//	config.SetBech32PrefixForValidator(yourBech32PrefixValAddr, yourBech32PrefixValPub)
	//	config.SetBech32PrefixForConsensusNode(yourBech32PrefixConsAddr, yourBech32PrefixConsPub)
	//	config.SetCoinType(yourCoinType)
	//	config.SetFullFundraiserPath(yourFullFundraiserPath)
	//	config.Seal()

	// AddrLen defines a valid address length
	AddrLen = 20

	// Atom in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	CoinType = 118

	// BIP44Prefix is the parts of the BIP44 HD path that are fixed by
	// what we used during the fundraiser.
	FullFundraiserPath = "44'/118'/0'/0/0"
)

// Address is a common interface for different types of addresses used by the SDK
type AddressI interface {
	Equals(Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var _ AddressI = Address{}

var _ yaml.Marshaler = Address{}

// VerifyAddressFormat verifies that the provided bytes form a valid address
// according to the default address rules or a custom address verifier set by
// GetConfig().SetAddressVerifier()
func VerifyAddressFormat(bz []byte) error {
	verifier := GetConfig().GetAddressVerifier()
	if verifier != nil {
		return verifier(bz)
	}
	if len(bz) != AddrLen {
		return errors.New("Incorrect address length")
	}
	return nil
}

// Address a wrapper around bytes meant to represent an address.
// When marshaled to a string or JSON, it uses Bech32.
type Address []byte

// AddressFromHex creates an Address from a hex string.
func AddressFromHex(address string) (addr Address, err error) {
	if len(address) == 0 {
		return Address{}, nil
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return Address(bz), nil
}

// Returns boolean for whether two Addresses are Equal
func (aa Address) Equals(aa2 Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Returns boolean for whether an Address is empty
func (aa Address) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := Address{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa Address) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *Address) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa Address) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *Address) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AddressFromHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *Address) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AddressFromHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa Address) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa Address) String() string {
	if aa.Empty() {
		return ""
	}

	str := hex.EncodeToString(aa.Bytes())

	return str
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa Address) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

// get ConsAddress from pubkey
func GetAddress(pubkey crypto.PubKey) Address {
	return Address(pubkey.Address())
}

// ----------------------------------------------------------------------------
// auxiliary
// ----------------------------------------------------------------------------

func HexAddressPubKey(pub crypto.PubKey) string {
	return hex.EncodeToString(pub.Bytes())
}

func GetAddressPubKeyFromHex(pubkey string) (pk crypto.PubKey, err error) {

	bz, err := GetFromHex(pubkey)
	if err != nil {
		return nil, err
	}
	pk, err = cryptoAmino.PubKeyFromBytes(bz)

	if err != nil {
		return nil, err
	}

	return pk, nil
}

// GetFromHex decodes a bytestring from a Hex encoded string.
func GetFromHex(hexString string) ([]byte, error) {
	if len(hexString) == 0 {
		return nil, errors.New("decoding hex address failed: must provide an address")
	}

	bz, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// ----------------------------------------------------------------------------
// Backward Compatibility
// ----------------------------------------------------------------------------

func ConsAddressFromHex(address string) (addr Address, err error) {
	return AddressFromHex(address)
}
func ValAddressFromHex(address string) (addr Address, err error) {
	return AddressFromHex(address)
}
func AccAddressFromHex(address string) (addr Address, err error) {
	return AddressFromHex(address)
}
func GetAccPubKeyHex(pubkey string) (pk crypto.PubKey, err error) {
	return GetAddressPubKeyFromHex(pubkey)
}

func GetValPubKeyHex(pubkey string) (pk crypto.PubKey, err error) {
	return GetAddressPubKeyFromHex(pubkey)
}

func GetConsPubKeyHex(pubkey string) (pk crypto.PubKey, err error) {
	return GetAddressPubKeyFromHex(pubkey)
}
