package context

import (
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

type Context struct {
	Config      *cfg.Config
	Logger      log.Logger
	TraceWriter string
}

func NewDefaultContext() *Context {
	return NewContext(
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		"", // todo broken
	)
}

func NewContext(config *cfg.Config, logger log.Logger, traceWriterPath string) *Context {
	return &Context{config, logger, traceWriterPath}
}
