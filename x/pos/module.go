package staking

import (
	"encoding/json"
	"github.com/pokt-network/posmint/x/pos/keeper"

	"github.com/gorilla/mux"
	"github.com/pokt-network/posmint/client/context"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/pos/client/cli"
	"github.com/pokt-network/posmint/x/pos/client/rest"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the staking module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the staking
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic

	keeper        Keeper
	accountKeeper types.AccountKeeper
	supplyKeeper  types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper types.AccountKeeper, supplyKeeper types.SupplyKeeper) AppModule {

	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns the staking module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the staking module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the staking module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the pos module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	keeper.BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return keeper.EndBlocker(ctx, am.keeper)
}
