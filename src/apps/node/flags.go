package node

var Flags []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:  "host, addr",
		Usage: "listen on specific host",
	},
	cli.StringFlag{
		Name:  "port, p",
		Usage: "listen on specific port",
	},
}
