package router

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) *gin.Engine {
	// 装载路由
	Api(r)
	return r
}
