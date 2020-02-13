package types

import (
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/libs/common"
	"strings"

	crkeys "github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder sdk.TxEncoder
	keybase   crkeys.Keybase
	chainID   string
	memo      string
	fees      sdk.Coins
}

// NewTxBuilder returns a new initialized TxBuilder.
func NewTxBuilder(txEncoder sdk.TxEncoder, chainID, memo string, fees sdk.Coins) TxBuilder {
	return TxBuilder{
		txEncoder: txEncoder,
		keybase:   nil,
		chainID:   chainID,
		memo:      memo,
		fees:      fees,
	}
}

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxEncoder() sdk.TxEncoder { return bldr.txEncoder }

// Keybase returns the keybase
func (bldr TxBuilder) Keybase() crkeys.Keybase { return bldr.keybase }

// ChainID returns the chain id
func (bldr TxBuilder) ChainID() string { return bldr.chainID }

// Memo returns the memo message
func (bldr TxBuilder) Memo() string { return bldr.memo }

// Fees returns the fees for the transaction
func (bldr TxBuilder) Fees() sdk.Coins { return bldr.fees }

// WithTxEncoder returns a copy of the context with an updated codec.
func (bldr TxBuilder) WithTxEncoder(txEncoder sdk.TxEncoder) TxBuilder {
	bldr.txEncoder = txEncoder
	return bldr
}

// WithChainID returns a copy of the context with an updated chainID.
func (bldr TxBuilder) WithChainID(chainID string) TxBuilder {
	bldr.chainID = chainID
	return bldr
}

// WithFees returns a copy of the context with an updated fee.
func (bldr TxBuilder) WithFees(fees string) TxBuilder {
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		panic(err)
	}

	bldr.fees = parsedFees
	return bldr
}

// WithKeybase returns a copy of the context with updated keybase.
func (bldr TxBuilder) WithKeybase(keybase crkeys.Keybase) TxBuilder {
	bldr.keybase = keybase
	return bldr
}

// WithMemo returns a copy of the context with an updated memo.
func (bldr TxBuilder) WithMemo(memo string) TxBuilder {
	bldr.memo = strings.TrimSpace(memo)
	return bldr
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a address, passphrase, and a set of messages.
func (bldr TxBuilder) BuildAndSign(address sdk.Address, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	if bldr.keybase == nil {
		return nil, errors.New(fmt.Sprintf("cant build and sign transaciton: the keybase is nil"))
	}
	if bldr.chainID == "" {
		return nil, errors.New(fmt.Sprintf("cant build and sign transaciton: the chainID is empty"))
	}
	entropy := common.RandInt64()
	bytesToSign := StdSignBytes(bldr.chainID, entropy, bldr.fees, msgs, bldr.memo)
	sigBytes, pubkey, err := bldr.keybase.Sign(address, passphrase, bytesToSign)
	if err != nil {
		return nil, err
	}
	sig := StdSignature{
		PublicKey: pubkey,
		Signature: sigBytes,
	}
	return bldr.txEncoder(NewStdTx(msgs, bldr.fees, []StdSignature{sig}, bldr.memo, entropy))
}
