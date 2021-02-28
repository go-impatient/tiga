package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"

	"moocss.com/tiga/app"
	"moocss.com/tiga/app/router"
	"moocss.com/tiga/internal/service"
	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/conf/file"
	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/server"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./config/", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, services *service.Services) *app.App {
	// 初始化web
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

	// 初始化web路由
	router.RegisterRoutes(engine, services)
	return app.New(
		app.Version("v1.0.0"),
		app.Server(srv),
	)
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	// 1. 初始化配置文件
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

	// 2. 初始化 app 生命周期
	app, err := InitApp(logger)
	if err != nil {
		panic(err)
	}
	// 3. 启动服务并等待停止信号
	if err := app.Run(); err != nil {
		l.Errorf("start failed: %v\n", err)
	}
}
