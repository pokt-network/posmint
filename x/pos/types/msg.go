package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStake{}
	_ sdk.Msg = &MsgBeginUnstake{}
	_ sdk.Msg = &MsgUnjail{}
	_ sdk.Msg = &MsgSend{}
)

//----------------------------------------------------------------------------------------------------------------------
// MsgStake - struct for staking transactions
type MsgStake struct {
	PubKey crypto.PublicKey `json:"pubkey" yaml:"pubkey"`
	Value  sdk.Int          `json:"value" yaml:"value"`
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgStake) GetSigner() sdk.Address {
	addrs := sdk.Address(msg.PubKey.Address())
	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgStake) ValidateBasic() sdk.Error {

	if msg.PubKey == nil || msg.PubKey.RawString() == "" {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	return nil
}

// Route provides router key for msg
func (msg MsgStake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgStake) Type() string { return "stake_validator" }

// GetFee get fee for msg
func (msg MsgStake) GetFee() sdk.Int {
	return sdk.NewInt(PosFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------
// MsgBeginUnstake - struct for unstaking transaciton
type MsgBeginUnstake struct {
	Address sdk.Address `json:"validator_address" yaml:"validator_address"`
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgBeginUnstake) GetSigner() sdk.Address {
	return msg.Address
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBeginUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgBeginUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

// Route provides router key for msg
func (msg MsgBeginUnstake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgBeginUnstake) Type() string { return "begin_unstaking_validator" }

// GetFee get fee for msg
func (msg MsgBeginUnstake) GetFee() sdk.Int {
	return sdk.NewInt(PosFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------
// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr sdk.Address `json:"address" yaml:"address"` // address of the validator operator
}

// Route provides router key for msg
func (msg MsgUnjail) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgUnjail) Type() string { return "unjail" }

// GetFee get fee for msg
func (msg MsgUnjail) GetFee() sdk.Int {
	return sdk.NewInt(PosFeeMap[msg.Type()])
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUnjail) GetSigner() sdk.Address {
	return msg.ValidatorAddr
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------

// MsgSend structure for sending coins
type MsgSend struct {
	FromAddress sdk.Address
	ToAddress   sdk.Address
	Amount      sdk.Int
}

// Route provides router key for msg
func (msg MsgSend) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgSend) Type() string { return "send" }

// GetFee get fee for msg
func (msg MsgSend) GetFee() sdk.Int {
	return sdk.NewInt(PosFeeMap[msg.Type()])
}

// GetSigner return address(es) that must sign over msg.GetSignBytes()
func (msg MsgSend) GetSigner() sdk.Address {
	return msg.FromAddress
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.ToAddress.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSendAmount(DefaultCodespace)
	}
	return nil
}
