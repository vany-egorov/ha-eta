package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/apps/node"
)

func main() {
	nodeApp := node.App{}

	app := cli.NewApp()

	cli.VersionPrinter = versionPrinter

	app.Version = version
	app.Name = "ha-eta"
	app.Usage = "ha-eta control daemons, services, utils, tools, clis"

	// specify command by default
	nodeApp.CmdTrySetDefaultAction(app)

	app.Commands = []cli.Command{
		nodeApp.Cmd(),
		// ... add other commands here if any
	}

	if e := app.Run(os.Args); e != nil {
		fmt.Fprint(os.Stderr, e.Error())
	}
}
