package commands

import (
	"github.com/urfave/cli/v2"
)

var (
	InitCommand = &cli.Command{
		Name:      "init",
		Action:    initGenesis,
		Aliases:   []string{"i"},
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags:     nil,
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
)

func initGenesis(ctx *cli.Context) error {
	log.Debug("initGenesis cmd")
	return nil
}
