package types

const (
	DAOTransferFee    = 100000
	MsgChangeParamFee = 100000
	MsgUpgradeFee     = 100000
)

var (
	GovFeeMap = map[string]int64{
		MsgDAOTransferName: DAOTransferFee,
		MsgChangeParamName: MsgChangeParamFee,
		MsgUpgradeName:     MsgUpgradeFee,
	}
)
