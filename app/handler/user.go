package handler

import (
	"github.com/gin-gonic/gin"
	"moocss.com/tiga/app/middleware"
	"moocss.com/tiga/internal/service"
	"net/http"
)

// UserHandler ...
type userHandler struct {
	UserService *service.UserService
}

// Login 用户登录
func (handler *userHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"text": "登录成功.",
		})
	}
}

// Register 注册
func (handler *userHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"text": "注册成功.",
		})
	}
}

// MakeUserHandler ...
func MakeUserHandler(r *gin.Engine, services *service.Services) {
	handler := &userHandler{UserService: services.UserService}

	userGroup := r.Group("/user")
	userGroup.Use(middleware.TokenAuth())
	{
		userGroup.GET("/login", handler.Login())
		userGroup.POST("/register", handler.Register())
	}
}
