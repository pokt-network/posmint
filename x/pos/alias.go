// nolint
package staking

import (
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/keeper"
	"github.com/pokt-network/posmint/x/pos/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	// functions aliases
	RegisterInvariants                 = keeper.RegisterInvariants
	NewQuerier                         = keeper.NewQuerier
	RegisterCodec                      = types.RegisterCodec
	ErrNoValidatorFound                = types.ErrNoValidatorFound
	ErrValidatorPubKeyExists           = types.ErrValidatorPubKeyExists
	ErrValidatorPubKeyTypeNotSupported = types.ErrValidatorPubKeyTypeNotSupported
	ErrValidatorJailed                 = types.ErrValidatorJailed
	ErrBadDenom                        = types.ErrBadDenom
	DefaultGenesisState                = types.DefaultGenesisState
	NewValidator                       = types.NewValidator
	ModuleCdc                          = types.ModuleCdc
)

type (
	Keeper                  = keeper.Keeper
	CodeType                = types.CodeType
	GenesisState            = types.GenesisState
	PrevStateValidatorPower = types.PrevStatePowerMapping
	MultiStakingHooks       = types.MultiPOSHooks
	MsgCreateValidator      = types.MsgStake
	Params                  = types.Params
	Pool                    = types.StakingPool
	QueryValidatorParams    = types.QueryValidatorParams
	QueryValidatorsParams   = types.QueryValidatorsParams
	Validator               = types.Validator
	Validators              = types.Validators
	ValidatorI              = exported.ValidatorI
)
