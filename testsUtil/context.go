package testUtil

import (
	"github.com/pokt-network/posmint/baseapp"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func GetNewContext(bapp *baseapp.BaseApp) sdk.Context {
header := abci.Header{Height: 0}
newContext := bapp.NewContext(false, header)
return newContext
}