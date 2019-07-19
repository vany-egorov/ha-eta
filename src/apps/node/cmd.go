package node

import (
	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/apps/node/initializers"
)

type Action uint8

const (
	ActionMain Action = iota
)

var actionString = map[Action]string{
	ActionMain: "main",
}

func (it Action) String() string { return actionString[it] }

func (a *App) Cmd() cli.Command {
	return cli.Command{
		Name:    "node",
		Usage:   "ha ETA min service node",
		Aliases: []string{"eta-min", "server"},
		Flags:   initializers.FLAGS,
		Action: func(c *cli.Context) error {
			return a.Main(c, ActionMain)
		},
		Subcommands: a.CmdSubcommands(),
	}
}

func (a *App) CmdSubcommands() []cli.Command {
	return nil
}
