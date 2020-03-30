package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgChangeParam{}
	_ sdk.Msg = &MsgDAOTransfer{}
)

const (
	MsgDAOTransferName = "dao_tranfer"
	MsgChangeParamName = "change_param"
	MsgUpgradeName     = "upgrade"
)

//----------------------------------------------------------------------------------------------------------------------
// MsgChangeParam structure for changing governance parameters
type MsgChangeParam struct {
	FromAddress sdk.Address `json:"address"`
	ParamKey    string      `json:"param_key"`
	ParamVal    interface{} `json:"param_value"`
}

//nolint
func (msg MsgChangeParam) Route() string { return RouterKey }
func (msg MsgChangeParam) Type() string  { return MsgChangeParamName }
func (msg MsgChangeParam) GetSigners() []sdk.Address {
	return []sdk.Address{msg.FromAddress}
}

func (msg MsgChangeParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeParam) ValidateBasic() sdk.Error {
	if msg.FromAddress == nil {
		return sdk.ErrInvalidAddress("nil address")
	}
	if msg.ParamKey == "" {
		return ErrEmptyKey(ModuleName)
	}
	if msg.ParamVal == nil {
		return ErrEmptyValue(ModuleName)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// MsgChangeParam structure for changing governance parameters
type MsgDAOTransfer struct {
	FromAddress sdk.Address `json:"from_address"`
	ToAddress   sdk.Address `json:"to_address"`
	Amount      sdk.Int     `json:"amount"`
	Action      string      `json:"action"`
}

//nolint
func (msg MsgDAOTransfer) Route() string { return RouterKey }
func (msg MsgDAOTransfer) Type() string  { return MsgDAOTransferName }
func (msg MsgDAOTransfer) GetSigners() []sdk.Address {
	return []sdk.Address{msg.FromAddress}
}

func (msg MsgDAOTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDAOTransfer) ValidateBasic() sdk.Error {
	if msg.FromAddress == nil {
		return sdk.ErrInvalidAddress("nil from address")
	}
	if msg.Amount.Int64() == 0 {
		return ErrZeroValueDAOAction(ModuleName)
	}
	daoAction, err := DAOActionFromString(msg.Action)
	if err != nil {
		return err
	}
	if daoAction == DAOTransfer && msg.ToAddress == nil {
		return sdk.ErrInvalidAddress("nil to address")
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// MsgUpgrade structure for changing governance parameters
type MsgUpgrade struct {
	Address sdk.Address `json:"address"`
	Upgrade Upgrade     `json:"upgrade"`
}

//nolint
func (msg MsgUpgrade) Route() string { return RouterKey }
func (msg MsgUpgrade) Type() string  { return MsgUpgradeName }
func (msg MsgUpgrade) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

func (msg MsgUpgrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUpgrade) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("nil from address")
	}
	if msg.Upgrade.UpgradeHeight() == 0 {
		return ErrZeroHeightUpgrade(ModuleName)
	}
	if msg.Upgrade.UpgradeVersion() == "" {
		return ErrZeroHeightUpgrade(ModuleName)
	}
	return nil
}
