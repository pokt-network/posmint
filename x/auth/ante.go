package auth

import (
	"bytes"
	"fmt"
	posCrypto "github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewAnteHandler returns an AnteHandler that checks signatures and deducts fees from the first signer.
func NewAnteHandler(ak AccountKeeper, supplyKeeper types.SupplyKeeper) sdk.AnteHandler {
	return func(ctx sdk.Ctx, tx sdk.Tx, txBz []byte, tmNode *node.Node, simulate bool, ) (newCtx sdk.Ctx, res sdk.Result, abort bool) {
		if addr := supplyKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
		}
		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(StdTx)
		if !ok {
			return newCtx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}
		//Default Fee
		defaultFee := sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(100000))
		//Check Tx Amount and Denomination.
		if !stdTx.Fee.AmountOf(sdk.DefaultStakeDenom).Equal(defaultFee.Amount) {
			return newCtx, sdk.ErrInsufficientFee(
				fmt.Sprintf(
					"insufficient fees; got: %q required: %q", stdTx.Fee, defaultFee,
				),
			).Result(), true
		}
		params := ak.GetParams(ctx)
		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			res := EnsureSufficientMempoolFees(ctx, stdTx.Fee)
			if !res.IsOK() {
				return newCtx, res, true
			}
		}
		if res := ValidateSigCount(stdTx, params); !res.IsOK() {
			return newCtx, res, true
		}
		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), true
		}
		if res := ValidateMemo(stdTx, params); !res.IsOK() {
			return newCtx, res, true
		}
		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		signerAddrs := stdTx.GetSigners()
		signerAccs := make([]Account, len(signerAddrs))
		// fetch first signer, who's going to pay the fees
		signerAccs[0], res = GetSignerAcc(ctx, ak, signerAddrs[0])
		if !res.IsOK() {
			return newCtx, res, true
		}
		// deduct the fees
		if !stdTx.Fee.IsZero() {
			res = DeductFees(supplyKeeper, ctx, signerAccs[0], stdTx.Fee)
			if !res.IsOK() {
				return newCtx, res, true
			}

			// reload the account as fees have been deducted
			signerAccs[0] = ak.GetAccount(ctx, signerAccs[0].GetAddress())
		}
		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		stdSigs := stdTx.GetSignatures()
		for i := 0; i < len(stdSigs); i++ {
			// skip the fee payer, account is cached and fees were deducted already
			if i != 0 {
				signerAccs[i], res = GetSignerAcc(ctx, ak, signerAddrs[i])
				if !res.IsOK() {
					return newCtx, res, true
				}
			}
			// check signature, return account with incremented nonce
			signBytes := GetSignBytes(ctx.ChainID(), stdTx)
			signerAccs[i], res = processSig(signerAccs[i], stdSigs[i], signBytes, simulate)
			// check for duplicate transaction not in cache todo added
			// todo when editing tendermint, pass txIndexer so no http
			c := client.NewHTTP(tmNode.Config().RPC.ListenAddress, "/websocket")
			_, err := c.Tx(tmTypes.Tx(txBz).Hash(), false)
			if err == nil {
				return newCtx,
					sdk.ErrUnauthorized(fmt.Sprint("transaction with this hash already found, possible replay attack, please try a different entropy")).Result(),
					true
			}
			if !res.IsOK() {
				return newCtx, res, true
			}
			ak.SetAccount(ctx, signerAccs[i])
		}
		return ctx, sdk.Result{}, false // continue...
	}
}

// GetSignerAcc returns an account for a given address that is expected to sign
// a transaction.
func GetSignerAcc(ctx sdk.Ctx, ak AccountKeeper, addr sdk.Address) (Account, sdk.Result) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, sdk.Result{}
	}
	return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr)).Result()
}

// ValidateSigCount validates that the transaction has a valid cumulative total
// amount of signatures.
func ValidateSigCount(stdTx StdTx, params Params) sdk.Result {
	stdSigs := stdTx.GetSignatures()

	sigCount := 0
	for i := 0; i < len(stdSigs); i++ {
		sigCount += CountSubKeys(stdSigs[i].PublicKey)
		if uint64(sigCount) > params.TxSigLimit {
			return sdk.ErrTooManySignatures(
				fmt.Sprintf("signatures: %d, limit: %d", sigCount, params.TxSigLimit),
			).Result()
		}
	}

	return sdk.Result{}
}

// ValidateMemo validates the memo size.
func ValidateMemo(stdTx StdTx, params Params) sdk.Result {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf(
				"maximum number of characters is %d but received %d characters",
				params.MaxMemoCharacters, memoLength,
			),
		).Result()
	}

	return sdk.Result{}
}

// verify the signature and increment the sequence. If the account doesn't have
// a pubkey, set it.
func processSig(acc Account, sig StdSignature, signBytes []byte, simulate bool) (updatedAcc Account, res sdk.Result) {
	pubKey, res := ProcessPubKey(acc, sig)
	if !res.IsOK() {
		return nil, res
	}
	err := acc.SetPubKey(posCrypto.PubKeyToPublicKey(pubKey))
	if err != nil {
		return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
	}

	if !simulate && !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, sdk.ErrUnauthorized("signature verification failed; verify correct account sequence and chain-id").Result()
	}

	//if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
	//	panic(err)
	//}
	// TODO removed

	return acc, res
}

// ProcessPubKey verifies that the given account address matches that of the
// StdSignature. In addition, it will set the public key of the account if it
// has not been set.
func ProcessPubKey(acc Account, sig StdSignature) (crypto.PubKey, sdk.Result) {
	// If pubkey is not known for account, set it from the StdSignature.
	pubKey := acc.GetPubKey()

	if pubKey == nil {
		pubKey = sig.PublicKey
		if pubKey == nil {
			return nil, sdk.ErrInvalidPubKey("PubKey not found").Result()
		}

		if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
			return nil, sdk.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", acc.GetAddress())).Result()
		}
	}

	return pubKey, sdk.Result{}
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the CoinKeeper (in addition to the AccountKeeper, because
// the CoinKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(supplyKeeper types.SupplyKeeper, ctx sdk.Ctx, acc Account, fees sdk.Coins) sdk.Result {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee amount: %s", fees)).Result()
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", coins, fees),
		).Result()
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", spendableCoins, fees),
		).Result()
	}

	err := supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// EnsureSufficientMempoolFees verifies that the given transaction has supplied
// enough fees to cover a proposer's minimum fees. A result object is returned
// indicating success or failure.
//
// Contract: This should only be called during CheckTx as it cannot be part of
// consensus.
func EnsureSufficientMempoolFees(ctx sdk.Ctx, stdFee sdk.Coins) sdk.Result {
	minGasPrices := ctx.MinGasPrices()
	if !minGasPrices.IsZero() {
		requiredFees := make(sdk.Coins, len(minGasPrices))

		if !stdFee.IsAnyGTE(requiredFees) {
			return sdk.ErrInsufficientFee(
				fmt.Sprintf(
					"insufficient fees; got: %q required: %q", stdFee, requiredFees,
				),
			).Result()
		}
	}
	return sdk.Result{}
}

// GetSignBytes returns a slice of bytes to sign over for a given transaction
// and an account.
func GetSignBytes(chainID string, stdTx StdTx) []byte {
	return StdSignBytes(
		chainID, stdTx.Entropy, stdTx.Fee, stdTx.Msgs, stdTx.Memo,
	)
}
