package baseapp

import (
	sdk "github.com/pokt-network/posmint/types"
)

var testHandler = func(_ sdk.Ctx, _ sdk.Msg) sdk.Result {
	return sdk.Result{}
}
