// +build !test node

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
	app.Name = "node"
	app.Usage = "node"
	app.Version = version
	app.EnableBashCompletion = true
	app.Flags = nodeInitializers.FLAGS
	app.Action = func(c *cli.Context) error { return nodeApp.Main(c, node.ActionMain) }

	cli.VersionPrinter = nodeApp.ShowVersionAndExit

	app.Run(os.Args)
}
