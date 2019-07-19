package node

import (
	cli "gopkg.in/urfave/cli.v1"
)

func (a *App) Cmd() cli.Command {
	return cli.Command{
		Name:    "node",
		Usage:   "ha ETA min service node",
		Aliases: []string{"eta-min", "server"},
		Flags:   Flags,
		Action: func(c *cli.Context) error {
			a.Main(c)
			return nil
		},
		Subcommands: a.CmdSubcommands(),
	}
}

func (a *App) CmdSubcommands() []cli.Command {
	return nil
}
