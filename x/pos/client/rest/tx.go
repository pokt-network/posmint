package rest

import (
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/tendermint/crypto"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pokt-network/posmint/client/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/rest"
	"github.com/pokt-network/posmint/x/auth/client/utils"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/pos/{validatorAddr}/stake",
		postStakeHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/pos/{validatorAddr}/unstake",
		postUnstakeHandlerFn(cliCtx),
	).Methods("POST")
}

type (
	// StakeRequest defines the properties of a stake request's body.
	StakeRequest struct {
		BaseReq rest.BaseReq  `json:"base_req" yaml:"base_req"`
		PubKey  crypto.PubKey `json:"pubkey" yaml:"pubkey"`
		Amount  sdk.Coin      `json:"amount" yaml:"amount"`
	}

	// UnstakeRequest defines the properties of a unstake request's body.
	UnstakeRequest struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
)

func postStakeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StakeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.ValAddressFromBech32(req.BaseReq.From)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		pubKey, err := cliCtx.Keybase.GetByAddress(sdk.AccAddress(address))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// todo may have to pass in the public key cause tendermint doesn't use the same pub key for consensus
		msg := types.NewMsgStake(address, pubKey.GetPubKey(), req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postUnstakeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnstakeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		address, err := sdk.ValAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgBeginUnstake(address)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
