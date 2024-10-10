package middleware

import (
	"net/http"
	"service/constant"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientID := ctx.Request.Header.Get("clientId")
		if clientID != "tC0ND8ar26Jk9L5b" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code": constant.FORBIDDEN,
				"msg":  "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
