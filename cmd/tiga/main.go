package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"moocss.com/tiga/cmd/tiga/create"
	"moocss.com/tiga/cmd/tiga/server"
)

var usageStr = `
_____________________________ 
___  __/___  _/_  ____/__    |
__  /   __  / _  / __ __  /| |
_  /   __/ /  / /_/ / _  ___ |
/_/    /___/  \____/  /_/  |_|
`

func run() {
	app := cli.NewApp()
	app.Name = "Tiga"
	app.Version = Version
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "moocss",
			Email: "moocss@gmail.com",
		},
	}
	app.Copyright = "(c) 2021 moocss"
	app.Usage = "一个轻量级的应用服务"
	app.UsageText = usageStr
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Commands = cli.Commands{
		server.Cmd,
		create.Cmd,
	}
	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "brace for impact\n")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	run()
}
