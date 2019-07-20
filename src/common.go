package main

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v1"
)

var (
	buildDate string
	version   string
)

func versionPrinter(c *cli.Context) {
	fmt.Printf("version: %s\n", version)
	fmt.Printf("build-date: %s\n", buildDate)
}
