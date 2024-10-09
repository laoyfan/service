package middleware

import (
	"net/http"
	"service/constant"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Request.Header.Get("clientId")
		if clientID != "tC0ND8ar26Jk9L5b" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": constant.FORBIDDEN,
				"msg":  "无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
