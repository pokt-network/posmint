package genaccounts

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth/exported"
	"github.com/pokt-network/posmint/x/genaccounts/types"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.GenesisState{})
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// extra function from sdk.AppModuleBasic
// iterate the genesis accounts and perform an operation at each of them
// - to used by other modules
func (AppModuleBasic) IterateGenesisAccounts(cdc *codec.Codec, appGenesis map[string]json.RawMessage,
	iterateFn func(exported.Account) (stop bool)) {

	genesisState := types.GetGenesisStateFromAppState(cdc, appGenesis)
	for _, genAcc := range genesisState {
		acc := genAcc.ToAccount()
		if iterateFn(acc) {
			break
		}
	}
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
	})
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, types.ModuleCdc, am.accountKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.accountKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}
