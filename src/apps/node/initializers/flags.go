package initializers

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/lib/environment"
	"github.com/vany-egorov/ha-eta/lib/helpers"
)

var (
	FLAGS []cli.Flag = []cli.Flag{
		cli.StringFlag{
			Name:   "environment, env, e",
			Value:  DefaultEnvironment.String(),
			Usage:  fmt.Sprintf("Application environment. Possible environments are number of: %s;", environment.EnvironmentValidListAsString()),
			EnvVar: EnvEnvironment,
		},
		cli.StringFlag{
			Name:   "config, conf, c",
			Usage:  "path to main configuration file",
			EnvVar: EnvConfig,
		},
		cli.StringFlag{
			Name:   "log, l",
			Usage:  "path to log directory",
			EnvVar: EnvLog,
		},
		cli.BoolFlag{
			Name:  "print-config",
			Usage: "print configuration and exit",
		},
		cli.BoolFlag{
			Name:  "foreground",
			Usage: "force application NOT to daemonize, even if in config said so",
		},
		cli.BoolFlag{
			Name:  "vv",
			Usage: "all loggers log level forced to debug",
		},
		cli.BoolFlag{
			Name:  "vvv",
			Usage: "all loggers log level forced to trace",
		},
	}
)

type Flags struct {
	Environment environment.Environment
	config      string
	log         string
	version     bool
	ctx         *cli.Context
}

func (it *Flags) Config() string           { return it.config }
func (it *Flags) Log() string              { return it.log }
func (it *Flags) ShowVersionAndExit() bool { return it.version }

func NewFlags(ctx *cli.Context) *Flags {
	it := new(Flags)

	if ctx == nil {
		app := cli.NewApp()
		app.Commands = []cli.Command{{
			Name:   "node",
			Flags:  FLAGS,
			Action: func(c *cli.Context) { ctx = c },
		}}
		app.Run([]string{"ha-eta", "node"})
	}

	it.ctx = ctx
	it.Environment = environment.NewEnvironment(ctx.String("environment"))
	it.config = ctx.String("config")
	it.log = ctx.String("log")
	it.version = ctx.Bool("version")

	if it.Environment.IsUnknown() {
		it.Environment = DefaultEnvironment
	}

	if it.config == "" {
		it.config = DefaultPathConfig
	}
	if it.log == "" {
		it.log = DefaultPathLog
	}

	helpers.PathsAbsolutize([]*string{&it.config, &it.log})

	return it
}
