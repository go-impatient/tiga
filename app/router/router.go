package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"moocss.com/tiga/app/handler"
	"moocss.com/tiga/app/middleware"
	"moocss.com/tiga/internal/service"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "welcome to tiga app.",
	})
}

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "page not found",
		})
	}
}

// RegisterRoutes ...
func RegisterRoutes(router *gin.Engine , services *service.Services) {
	// 使用中间件.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Handler())
	router.Use(middleware.NoCache)
	router.Use(middleware.Options)
	router.Use(middleware.Secure)
	router.Use(middleware.RequestId())

	// 404 Handler.
	router.NoRoute(NotFound())

	router.GET("/", rootHandler)

	// user 路由组
	handler.MakeUserHandler(router, services)
}
