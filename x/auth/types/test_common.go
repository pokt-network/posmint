// nolint noalias
package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

func NewTestTx(ctx sdk.Ctx, msgs []sdk.Msg, privs []crypto.PrivateKey, entropy int64, fee sdk.Coins) sdk.Tx {
	sigs := make([]StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := StdSignBytes(ctx.ChainID(), entropy, fee, msgs, "")
		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i] = StdSignature{PublicKey: priv.PublicKey(), Signature: sig}
	}
	tx := NewStdTx(msgs, fee, sigs, "", entropy)
	return tx
}
