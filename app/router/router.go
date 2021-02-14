package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func RegisterRoutes(router *gin.Engine) {
	// 使用中间件.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 404 Handler.
	router.NoRoute(NotFound())

	router.GET("/", rootHandler)

	// 其它路由组
}
