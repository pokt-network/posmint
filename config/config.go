package config

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

func NewDefaultConfig() *Context {
	return NewConfig(
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		"", // todo broken
	)
}

func NewConfig(config *cfg.Config, logger log.Logger, traceWriterPath string) *Context {
	return &Context{config, logger, traceWriterPath}
}
