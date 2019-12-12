package testUtil

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)


func defaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}