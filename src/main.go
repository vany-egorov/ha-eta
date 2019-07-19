package main

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/apps/node"
)

func main() {
	nodeApp := node.App{}

	app := cli.NewApp()
	app.Name = "ha-eta"
	app.Usage = "ha-eta control daemons, services, utils, tools, clis"
	app.Version = version

	// specify command by default
	app.Flags = node.FLAGS
	app.Action = func(c *cli.Context) { nodeApp.Main(c, node.ActionMain) }
	cli.VersionPrinter = nodeApp.ShowVersionAndExit

	app.Commands = []cli.Command{
		nodeApp.Cmd(),
		// ... add other commands here if any
	}

	app.Run(os.Args)
}
