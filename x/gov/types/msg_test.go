package types

import (
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsgChangeParam_ValidateBasic(t *testing.T) {
	m := MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "bank/sendenabled",
		ParamVal:    false,
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "",
		ParamVal:    false,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgChangeParam{
		ParamKey: "bank/sendenabled",
		ParamVal: false,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgChangeParam{
		FromAddress: getRandomValidatorAddress(),
		ParamKey:    "bank/sendenabled",
	}
	assert.NotNil(t, m.ValidateBasic())
}

func TestMsgDAOTransfer_ValidateBasic(t *testing.T) {
	m := MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.OneInt(),
		Action:      DAOTransferString,
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.OneInt(),
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		ToAddress:   getRandomValidatorAddress(),
		Amount:      types.ZeroInt(),
		Action:      DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		FromAddress: getRandomValidatorAddress(),
		Amount:      types.ZeroInt(),
		Action:      DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgDAOTransfer{
		ToAddress: getRandomValidatorAddress(),
		Amount:    types.ZeroInt(),
		Action:    DAOTransferString,
	}
	assert.NotNil(t, m.ValidateBasic())
}

func TestMsgUpgrade_ValidateBasic(t *testing.T) {
	m := MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  100,
			Version: "2.0.0",
		},
	}
	assert.Nil(t, m.ValidateBasic())
	m = MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  0,
			Version: "2.0.0",
		},
	}
	assert.NotNil(t, m.ValidateBasic())
	m = MsgUpgrade{
		Address: getRandomValidatorAddress(),
		Upgrade: Upgrade{
			Height:  100,
			Version: "",
		},
	}
	assert.NotNil(t, m.ValidateBasic())
}
