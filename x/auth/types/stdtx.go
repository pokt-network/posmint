package types

import (
	"encoding/json"
	"fmt"

	"github.com/pokt-network/posmint/codec"
	posCrypto "github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"gopkg.in/yaml.v2"
)

var (
	_ sdk.Tx = (*StdTx)(nil)

	//maxGasWanted = uint64((1 << 63) - 1)
)

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msgs       []sdk.Msg      `json:"msg" yaml:"msg"`
	Fee        sdk.Coins      `json:"fee" yaml:"fee"`
	Signatures []StdSignature `json:"signatures" yaml:"signatures"`
	Memo       string         `json:"memo" yaml:"memo"`
}

func NewStdTx(msgs []sdk.Msg, fee sdk.Coins, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsgs() []sdk.Msg { return tx.Msgs }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	stdSigs := tx.GetSignatures()
	if tx.Fee.IsValid() == false {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.String()))
	}
	if len(stdSigs) == 0 {
		return sdk.ErrNoSignatures("no signers")
	}
	if len(stdSigs) != len(tx.GetSigners()) {
		return sdk.ErrUnauthorized("wrong number of signers")
	}

	return nil
}

// CountSubKeys counts the total number of keys for a multi-sig public key.
func CountSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}

	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += CountSubKeys(subkey)
	}

	return numKeys
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []sdk.Address {
	seen := map[string]bool{}
	var signers []sdk.Address
	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}
	return signers
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignatures returns the signature of signers who signed the Msg.
// GetSignatures returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }

//__________________________________________________________

// StdFee includes the amount of coins paid in fees and the maximum
// gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
// type StdFee struct {
// 	Amount sdk.Coins `json:"amount" yaml:"amount"`
// 	Gas    uint64    `json:"gas" yaml:"gas"`
// }

// NewStdFee returns a new instance of StdFee
// func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
// 	return StdFee{
// 		Amount: amount,
// 		Gas:    gas,
// 	}
// }

// RawBytes for signing later
// func (fee StdFee) RawBytes() []byte {
// 	// normalize. XXX
// 	// this is a sign of something ugly
// 	// (in the lcd_test, client side its null,
// 	// server side its [])
// 	if len(fee.Amount) == 0 {
// 		fee.Amount = sdk.NewCoins()
// 	}
// 	bz, err := ModuleCdc.MarshalJSON(fee) // TODO
// 	if err != nil {
// 		panic(err)
// 	}
// 	return bz
// }

// GasPrices returns the gas prices for a StdFee.
//
// NOTE: The gas prices returned are not the true gas prices that were
// originally part of the submitted transaction because the fee is computed
// as fee = ceil(gasWanted * gasPrices).
// func (fee StdFee) GasPrices() sdk.DecCoins {
// 	return sdk.NewDecCoins(fee.Amount).QuoDec(sdk.NewDec(int64(fee.Gas)))
// }

//__________________________________________________________

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Memo          string            `json:"memo" yaml:"memo"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee sdk.Coins, msgs []sdk.Msg, memo string) []byte {
	var msgsBytes []json.RawMessage
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	var feeBytes json.RawMessage
	feeBytes, err := fee.MarshalJSON()
	if err != nil {
		panic(err)
	}
	bz, err := ModuleCdc.MarshalJSON(StdSignDoc{
		AccountNumber: accnum,
		ChainID:       chainID,
		Fee:           json.RawMessage(feeBytes),
		Memo:          memo,
		Msgs:          msgsBytes,
		Sequence:      sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// StdSignature represents a sig
type StdSignature struct {
	posCrypto.PublicKey `json:"pub_key" yaml:"pub_key"` // optional
	Signature           []byte `json:"signature" yaml:"signature"`
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = StdTx{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}

		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}

// MarshalYAML returns the YAML representation of the signature.
func (ss StdSignature) MarshalYAML() (interface{}, error) {
	var (
		bz     []byte
		pubkey string
		err    error
	)

	if ss.PublicKey != nil {
		pubkey = ss.PublicKey.RawString()
	}

	bz, err = yaml.Marshal(struct {
		PubKey    string
		Signature string
	}{
		PubKey:    pubkey,
		Signature: fmt.Sprintf("%s", ss.Signature),
	})
	if err != nil {
		return nil, err
	}

	return string(bz), err
}
