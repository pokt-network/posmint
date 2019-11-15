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

func NewDefaultConfig(rootDirectory string) *Config {
	return &Config{
		cfg.DefaultConfig().SetRoot(rootDirectory),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		"",
	}
}

func NewConfig(rootDirectory string, MaxNumberInboundPeers, MaxNumberOutboundPeers int, logger log.Logger, traceWriterPath string) *Config {
	// setup tendermint node config
	newTMConfig := cfg.DefaultConfig()
	newTMConfig.SetRoot(rootDirectory)
	newTMConfig.P2P.MaxNumInboundPeers = MaxNumberInboundPeers
	newTMConfig.P2P.MaxNumOutboundPeers = MaxNumberOutboundPeers

	return &Config{newTMConfig, logger, traceWriterPath}
}
