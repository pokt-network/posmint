package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
	// DefaultSendEnabled enabled
	DefaultSendEnabled = true
)

// ParamStoreKeySendEnabled is store's key for SendEnabled
var ParamStoreKeySendEnabled = []byte("sendenabled")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() sdk.KeyTable {
	return sdk.NewKeyTable(
		ParamStoreKeySendEnabled, false,
	)
}
