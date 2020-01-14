package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address sdk.Address
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr sdk.Address) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
