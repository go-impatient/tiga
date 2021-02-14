package create

import (
	"log"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:     "create",
	Usage:    "Tiga 工具",
	HideHelp: false,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "the output api file",
		},
	},
	Action: func(c *cli.Context) error {
		dir := c.String("output")
		log.Fatalf("dir: %v", dir)
		return nil
	},
	Subcommands: cli.Commands{
		&cli.Command{
			Name:    "doc",
			Aliases: []string{"d"},
			Usage:   "generate doc files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "the target dir",
				},
			},
			Action: func(c *cli.Context) error {
				dir := c.String("dir")
				log.Fatalf("dir: %v", dir)
				return nil
			},
		},
		&cli.Command{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "generate model code",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "the target dir",
				},
			},
			Action: func(c *cli.Context) error {
				dir := c.String("dir")
				log.Fatalf("dir: %v", dir)
				return nil
			},
		},
	},
}
