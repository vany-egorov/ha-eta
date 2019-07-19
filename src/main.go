package main

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/apps/node"
	nodeInitializers "github.com/vany-egorov/ha-eta/apps/node/initializers"
)

func main() {
	MustInitialize()

	nodeApp := node.App{}
	nodeApp.SetVersion(version)
	nodeApp.SetDoneAt(doneAt)

	app := cli.NewApp()
	app.Name = "ha-eta"
	app.Usage = "ha-eta control daemons, services, utils, tools, clis"
	app.Version = version

	app.Flags = nodeInitializers.FLAGS
	app.Action = func(c *cli.Context) { nodeApp.Main(c, node.ActionMain) }
	cli.VersionPrinter = nodeApp.ShowVersionAndExit

	app.Commands = []cli.Command{
		nodeApp.Cmd(),
	}

	app.Run(os.Args)
}
