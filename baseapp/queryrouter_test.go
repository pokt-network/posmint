package baseapp

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/pokt-network/posmint/types"
)

var testQuerier = func(_ sdk.Ctx, _ []string, _ abci.RequestQuery) (res []byte, err sdk.Error) {
	return nil, nil
}

func TestQueryRouter(t *testing.T) {
	qr := NewQueryRouter()

	// require panic on invalid route
	require.Panics(t, func() {
		qr.AddRoute("*", testQuerier)
	})

	qr.AddRoute("testRoute", testQuerier)
	q := qr.Route("testRoute")
	require.NotNil(t, q)

	// require panic on duplicate route
	require.Panics(t, func() {
		qr.AddRoute("testRoute", testQuerier)
	})
}
