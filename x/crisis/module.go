package crisis

import (
	"encoding/json"
	"github.com/pokt-network/posmint/crypto/keys"
	"github.com/tendermint/tendermint/node"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/crisis/internal/keeper"
	"github.com/pokt-network/posmint/x/crisis/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the crisis module.
type AppModuleBasic struct{}

// Name returns the crisis module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the crisis module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the crisis
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the crisis module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// AppModule implements an application module for the crisis module.
type AppModule struct {
	AppModuleBasic
	// NOTE: We store a reference to the keeper here so that after a module
	// manager is created, the invariants can be properly registered and
	// executed.
	keeper *keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper *keeper.Keeper, node *node.Node, keybase keys.Keybase) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name returns the crisis module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the crisis module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the crisis module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(*am.keeper)
}

// QuerierRoute returns no querier route.
func (AppModule) QuerierRoute() string { return "" }

// NewQuerierHandler returns no sdk.Querier.
func (AppModule) NewQuerierHandler() sdk.Querier { return nil }

// InitGenesis performs genesis initialization for the crisis module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Ctx, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	if data == nil {
		genesisState = DefaultGenesisState()
	} else {
		ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	}
	InitGenesis(ctx, *am.keeper, genesisState)

	am.keeper.AssertInvariants(ctx)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the crisis
// module.
func (am AppModule) ExportGenesis(ctx sdk.Ctx) json.RawMessage {
	gs := ExportGenesis(ctx, *am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the crisis module.
func (AppModule) BeginBlock(_ sdk.Ctx, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the crisis module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Ctx, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, *am.keeper)
	return []abci.ValidatorUpdate{}
}
