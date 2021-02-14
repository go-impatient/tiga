package server

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		EnvVars: []string{"TIGA_CONFIG"},
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "设置配置文件",
		Value:   "",
	},
}

var Cmd = &cli.Command{
	Name:     "server",
	Usage:    "Tiga 应用管理",
	HideHelp: false,
	Subcommands: cli.Commands{
		&cli.Command{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "运行Tiga服务",
			Action:  start,
			Flags:   flags,
		},
		&cli.Command{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "部署Tiga服务",
			Action:  deploy,
			Flags:   flags,
		},
	},
}

// 运行服务
func start(c *cli.Context) error {
	// var group errgroup.Group
	fileName := c.String("config")
	if len(fileName) == 0 {
		return errors.New("server s -c ./../../config/config.json 或者 server start -config ./../../config/config.json")
	}

	// 启动APP
	// app.Init(fileName);

	return nil
}

// 编译和部署
func deploy(c *cli.Context) error {
	return nil
}
