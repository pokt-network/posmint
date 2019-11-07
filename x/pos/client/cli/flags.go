package cli

import (
	flag "github.com/spf13/pflag"
)

// nolint
const (
	FlagAddressValidator = "validator"
	FlagPubKey           = "pubkey"
	FlagAmount           = "amount"
)

// common flagsets to add to various functions
var (
	FsPk        = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount    = flag.NewFlagSet("", flag.ContinueOnError)
	fsValidator = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded ConsPubKey of the validator")
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
	fsValidator.String(FlagAddressValidator, "", "The Bech32 address of the validator")
}
