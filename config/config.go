package config

import (
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

type Config struct {
	TmConfig    *cfg.Config
	Logger      log.Logger
	TraceWriter string
}

func NewDefaultConfig() *Config {
	return NewConfig(
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		"", // todo broken
	)
}

func NewConfig(config *cfg.Config, logger log.Logger, traceWriterPath string) *Config {
	return &Config{config, logger, traceWriterPath}
}
