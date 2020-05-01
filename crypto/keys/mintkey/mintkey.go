package mintkey

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	posCrypto "github.com/pokt-network/posmint/crypto"
	"github.com/tendermint/crypto/bcrypt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/armor"
	"github.com/tendermint/tendermint/crypto/xsalsa20symmetric"
	cmn "github.com/tendermint/tendermint/libs/common"
	"strconv"
)

const (
	blockTypePrivKey = "TENDERMINT PRIVATE KEY"
	blockTypeKeyInfo = "TENDERMINT KEY INFO"
	blockTypePubKey  = "TENDERMINT PUBLIC KEY"
	defaultKDF       = "bcrypt"
)

// Make bcrypt security parameter var, so it can be changed within the lcd test
// Making the bcrypt security parameter a var shouldn't be a security issue:
// One can't verify an invalid key by maliciously changing the bcrypt
// parameter during a runtime vulnerability. The main security
// threat this then exposes would be something that changes this during
// runtime before the user creates their key. This vulnerability must
// succeed to update this to that same value before every subsequent call
// to the keys command in future startups / or the attacker must get access
// to the filesystem. However, with a similar threat model (changing
// variables in runtime), one can cause the user to sign a different tx
// than what they see, which is a significantly cheaper attack then breaking
// a bcrypt hash. (Recall that the nonce still exists to break rainbow tables)
// For further notes on security parameter choice, see README.md
var BcryptSecurityParameter = 12

//-----------------------------------------------------------------
// add armor

// Armor the InfoBytes
func ArmorInfoBytes(bz []byte) string {
	return armorBytes(bz, blockTypeKeyInfo)
}

// Armor the PubKeyBytes
func ArmorPubKeyBytes(bz []byte) string {
	return armorBytes(bz, blockTypePubKey)
}

func armorBytes(bz []byte, blockType string) string {
	header := map[string]string{
		"type":    "Info",
		"version": "0.0.0",
	}
	return armor.EncodeArmor(blockType, header, bz)
}

//-----------------------------------------------------------------
// remove armor

// Unarmor the InfoBytes
func UnarmorInfoBytes(armorStr string) (bz []byte, err error) {
	return unarmorBytes(armorStr, blockTypeKeyInfo)
}

// Unarmor the PubKeyBytes
func UnarmorPubKeyBytes(armorStr string) (bz []byte, err error) {
	return unarmorBytes(armorStr, blockTypePubKey)
}

func unarmorBytes(armorStr, blockType string) (bz []byte, err error) {
	bType, header, bz, err := armor.DecodeArmor(armorStr)
	if err != nil {
		return
	}
	if bType != blockType {
		err = fmt.Errorf("Unrecognized armor type %q, expected: %q", bType, blockType)
		return
	}
	if header["version"] != "0.0.0" {
		err = fmt.Errorf("Unrecognized version: %v", header["version"])
		return
	}
	return
}

//-----------------------------------------------------------------
// encrypt/decrypt with armor
type ArmoredJson struct {
	Kdf        string `json:"kdf" yaml:"kdf"`
	Salt       string `json:"salt" yaml:"salt"`
	SecParam   string `json:"secparam" yaml:"secparam"`
	Hint       string `json:"hint" yaml:"hint"`
	Ciphertext string `json:"ciphertext" yaml:"ciphertext"`
}

func NewArmoredJson(kdf, salt, hint, ciphertext string) ArmoredJson {
	return ArmoredJson{
		Kdf:        kdf,
		Salt:       salt,
		SecParam:   strconv.Itoa(BcryptSecurityParameter),
		Hint:       hint,
		Ciphertext: ciphertext,
	}
}

// Encrypt and armor the private key.
func EncryptArmorPrivKey(privKey posCrypto.PrivateKey, passphrase, hint string) (string, error) {
	//first  encrypt the key
	saltBytes, encBytes := encryptPrivKey(privKey, passphrase)
	//"armor" the encrypted key encoding it in base64
	armorStr := base64.StdEncoding.EncodeToString(encBytes)
	//create the ArmoredJson with the parameters to be able to decrypt it later.
	armoredJson := NewArmoredJson(defaultKDF, fmt.Sprintf("%X", saltBytes), hint, armorStr)
	//marshalling to json
	js, err := json.Marshal(armoredJson)
	if err != nil {
		return "", err
	}
	//return the json string
	return string(js), nil
}

// encrypt the given privKey with the passphrase using a randomly
// generated salt and the AES-256 GCM cipher. returns the salt and the
// encrypted priv key.
func encryptPrivKey(privKey posCrypto.PrivateKey, passphrase string) (saltBytes []byte, encBytes []byte) {
	saltBytes = crypto.CRandBytes(16)
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), BcryptSecurityParameter)
	if err != nil {
		cmn.Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // get 32 bytes
	privKeyBytes := privKey.RawBytes()
	//encrypt using AES
	encBytes, err = EncryptAESGCM(key, privKeyBytes)
	if err != nil {
		cmn.Exit("Error encrypting bytes: " + err.Error())
	}
	return saltBytes, encBytes
}

// Unarmor and decrypt the private key.
func UnarmorDecryptPrivKey(armorStr string, passphrase string) (posCrypto.PrivateKey, error) {
	var privKey posCrypto.PrivateKey
	armoredJson := ArmoredJson{}
	//trying to unmarshal to ArmoredJson Struct
	err := json.Unmarshal([]byte(armorStr), &armoredJson)
	if err != nil {
		//if the unmarshal fails that could mean the armor in the keybase could be in the old format.
		//[RBM]
		fmt.Println("Checking if Pre RC 0.3.0 armored")
		//check the pre RC 0.3.0 unarmor
		return compatibilityUnarmor(armorStr, passphrase)
	}
	// check the ArmoredJson for the correct parameters on kdf and salt
	if armoredJson.Kdf != "bcrypt" {
		return privKey, fmt.Errorf("Unrecognized KDF type: %v", armoredJson.Kdf)
	}
	if armoredJson.Salt == "" {
		return privKey, fmt.Errorf("Missing salt bytes")
	}
	//decoding the salt
	saltBytes, err := hex.DecodeString(armoredJson.Salt)
	if err != nil {
		return privKey, fmt.Errorf("Error decoding salt: %v", err.Error())
	}
	//decoding the "armored" ciphertext stored in base64
	encBytes, err := base64.StdEncoding.DecodeString(armoredJson.Ciphertext)
	if err != nil {
		return privKey, fmt.Errorf("Error decoding ciphertext: %v", err.Error())
	}
	//decrypt the actual privkey with the parameters
	privKey, err = decryptPrivKey(saltBytes, encBytes, passphrase)
	return privKey, err
}

//compatibilityUnarmor - used to unarmor pre RC 0.3.0 keys in keybase or exported, this was the old UnarmorDecryptPrivKey
//[RBM]
func compatibilityUnarmor(armorStr string, passphrase string) (posCrypto.PrivateKey, error) {
	var privKey posCrypto.PrivateKey
	blockType, header, encBytes, err := armor.DecodeArmor(armorStr)
	if err != nil {
		return privKey, err
	}
	if blockType != blockTypePrivKey {
		return privKey, fmt.Errorf("Unrecognized armor type: %v", blockType)
	}
	if header["kdf"] != "bcrypt" {
		return privKey, fmt.Errorf("Unrecognized KDF type: %v", header["KDF"])
	}
	if header["salt"] == "" {
		return privKey, fmt.Errorf("Missing salt bytes")
	}
	saltBytes, err := hex.DecodeString(header["salt"])
	if err != nil {
		return privKey, fmt.Errorf("Error decoding salt: %v", err.Error())
	}
	privKey, err = decryptPrivKey(saltBytes, encBytes, passphrase)
	return privKey, err
}

func decryptPrivKey(saltBytes []byte, encBytes []byte, passphrase string) (privKey posCrypto.PrivateKey, err error) {
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), BcryptSecurityParameter)
	if err != nil {
		cmn.Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // Get 32 bytes
	//decrypt using AES
	privKeyBytes, err := DecryptAESGCM(key, encBytes)
	if err != nil {
		fmt.Println("Checking if Pre RC 0.3.0 encrypted")
		//[RBM]
		//Compatibility with pre RC 0.3.0 keys encrypted using salsa20
		privKeyBytes, err = xsalsa20symmetric.DecryptSymmetric(encBytes, key)
		if err != nil {
			return
		}
	}
	pk, err := posCrypto.NewPrivateKeyBz(privKeyBytes)
	if err != nil {
		fmt.Println("Checking if Pre RC 0.3.0 Key")
		//[RBM]
		//Compatibility with pre RC 0.3.0 keys using amino bytes
		pk, err2 := posCrypto.PrivKeyFromBytes(privKeyBytes)
		if err2 != nil {
			return
		}
		return pk, err2
	}
	return pk, err
}

func EncryptAESGCM(key []byte, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	nonce := key[:12]
	out := gcm.Seal(nil, nonce, src, nil)
	return out, nil
}

func DecryptAESGCM(key []byte, enBytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	nonce := key[:12]
	result, err := gcm.Open(nil, nonce, enBytes, nil)
	if err != nil {
		fmt.Printf("Can't Decrypt Using AES : %s", err)
		return nil, err
	}
	return result, nil
}
