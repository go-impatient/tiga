package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/jwt"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if s, exist := c.GetQuery("Authorization"); exist {
			token = s
		} else {
			token = c.GetHeader("Authorization")
		}

		if token == "" {
			c.JSON(401, map[string]interface{}{
				"code":    0,
				"message": "无权限",
			})
			c.Abort()
			return
		}

		var t string
		fmt.Sscanf(token, "Bearer %s", &t)

		secret := conf.Get("app.jwt.secret")
		j := jwt.New(jwt.WithSigningKey([]byte(secret)))
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == jwt.TokenExpired {
				c.JSON(402, map[string]interface{}{
					"code":    0,
					"message": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(403, map[string]interface{}{
				"code":    0,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)

		c.Next()
	}
}
