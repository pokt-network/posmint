package keeper

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/bank/internal/types"
	"github.com/pokt-network/posmint/x/params"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper

	ak         types.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(ak types.AccountKeeper,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType, blacklistedAddrs map[string]bool) BaseKeeper {

	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return BaseKeeper{
		BaseSendKeeper: NewBaseSendKeeper(ak, ps, codespace, blacklistedAddrs),
		ak:             ak,
		paramSpace:     ps,
	}
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	InputOutputCoins(ctx sdk.Ctx, inputs []types.Input, outputs []types.Output) sdk.Error
	SendCoins(ctx sdk.Ctx, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error

	SubtractCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SetCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) sdk.Error

	GetSendEnabled(ctx sdk.Ctx) bool
	SetSendEnabled(ctx sdk.Ctx, enabled bool)

	BlacklistedAddr(addr sdk.Address) bool
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak         types.AccountKeeper
	paramSpace params.Subspace

	// list of addresses that are restricted from receiving transactions
	blacklistedAddrs map[string]bool
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak types.AccountKeeper,
	paramSpace params.Subspace, codespace sdk.CodespaceType, blacklistedAddrs map[string]bool) BaseSendKeeper {

	return BaseSendKeeper{
		BaseViewKeeper:   NewBaseViewKeeper(ak, codespace),
		ak:               ak,
		paramSpace:       paramSpace,
		blacklistedAddrs: blacklistedAddrs,
	}
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseSendKeeper) InputOutputCoins(ctx sdk.Ctx, inputs []types.Input, outputs []types.Output) sdk.Error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	for _, in := range inputs {
		_, err := keeper.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}

	for _, out := range outputs {
		_, err := keeper.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeyRecipient, out.Address.String()),
			),
		)
	}

	return nil
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(ctx sdk.Ctx, fromAddr sdk.Address, toAddr sdk.Address, amt sdk.Coins) sdk.Error {
	_, err := keeper.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = keeper.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
		),
	})

	return nil
}

// SubtractCoins subtracts amt from the coins at the addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func (keeper BaseSendKeeper) SubtractCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := keeper.SetCoins(ctx, addr, newCoins)

	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseSendKeeper) AddCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins := keeper.GetCoins(ctx, addr)
	newCoins := oldCoins.Add(amt)

	if newCoins.IsAnyNegative() {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := keeper.SetCoins(ctx, addr, newCoins)
	return newCoins, err
}

// SetCoins sets the coins at the addr.
func (keeper BaseSendKeeper) SetCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) sdk.Error {

	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		panic(err)
	}

	keeper.ak.SetAccount(ctx, acc)
	return nil
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper BaseSendKeeper) GetSendEnabled(ctx sdk.Ctx) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx sdk.Ctx, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

// BlacklistedAddr checks if a given address is blacklisted (i.e restricted from
// receiving funds)
func (keeper BaseSendKeeper) BlacklistedAddr(addr sdk.Address) bool {
	return keeper.blacklistedAddrs[addr.String()]
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Ctx, addr sdk.Address) sdk.Coins
	HasCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) bool

	Codespace() sdk.CodespaceType
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak        types.AccountKeeper
	codespace sdk.CodespaceType
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak types.AccountKeeper, codespace sdk.CodespaceType) BaseViewKeeper {
	return BaseViewKeeper{ak: ak, codespace: codespace}
}

// Logger returns a module-specific logger.
func (keeper BaseViewKeeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Ctx, addr sdk.Address) sdk.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Ctx, addr sdk.Address, amt sdk.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

// Codespace returns the keeper's codespace.
func (keeper BaseViewKeeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}
