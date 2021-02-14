package main

import (
	"flag"
	"os"

	"moocss.com/tiga/app"
	"moocss.com/tiga/app/router"
	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/conf/file"
	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/log/stdlog"
	"moocss.com/tiga/pkg/server"

	"github.com/gin-gonic/gin"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./config/", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	logger := stdlog.NewLogger(stdlog.Writer(os.Stdout))
	defer logger.Close()

	c := conf.New(
		conf.WithSource(
			file.NewFile(flagconf),
		),
		conf.WithLogger(logger),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	l := log.NewHelper("main", logger)
	l.Info()

	engine := gin.New()

	srv := server.NewServer(
		server.Address(":8080"),
		server.HttpHandler(engine),
		server.Logger(logger),
	)

	switch conf.Get("app.mode") {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("unknown mode")
	}

	// 初始化路由和业务服务
	router.RegisterRoutes(engine)

	// 初始化 app 生命周期
	a := app.NewApp(
		app.Version("v1.0.0"),
		app.Server(srv),
	)

	// 启动服务并等待停止信号
	if err := a.Run(); err != nil {
		l.Errorf("start failed: %v\n", err)
	}
}
