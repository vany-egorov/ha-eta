package node

import (
	cli "gopkg.in/urfave/cli.v1"

	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type action uint8

const (
	actionUnknown action = iota
	actionMain
)

var flagsMain []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:  "host, server-host",
		Usage: "listen on specific host",
		Value: defaultServerHost,
	},
	cli.UintFlag{
		Name:  "port, p, server-port",
		Usage: "listen on specific port",
		Value: defaultServerPort,
	},
	cli.DurationFlag{
		Name:  "period-memstats",
		Usage: "log down memstats info for monitoring",
		Value: defaultPeriodMemstats,
	},
	cli.DurationFlag{
		Name:  "timeout-wait-terminate",
		Usage: "wait duration for workers to shutdown gracefully. otherwise force shutdown",
		Value: defaultTimeoutWaitTerminate,
	},
	cli.StringFlag{
		Name:  "geo-engine-kind",
		Usage: "backend kind used for geo detection",
		Value: geoEngine.DefaultKind.String(),
	},
	cli.StringFlag{
		Name:  "wheely-url",
		Usage: "url to wheely api",
		Value: wheely.DefaultUrlRaw,
	},
}

func (a *App) CmdTrySetDefaultAction(cliApp *cli.App) {
	cliApp.Flags = flagsMain
	cliApp.Action = func(c *cli.Context) error { return a.Main(c, actionMain) }
}

func (a *App) Cmd() cli.Command {
	return cli.Command{
		Name:    "node",
		Usage:   "ha ETA min service node",
		Aliases: []string{"s", "eta-min", "server"},
		Flags:   flagsMain,
		Action: func(c *cli.Context) error {
			return a.Main(c, actionMain)
		},
		Subcommands: a.CmdSubcommands(),
	}
}

func (a *App) CmdSubcommands() []cli.Command {
	return nil
}
