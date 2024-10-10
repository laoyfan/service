package controller

import (
	"errors"
	"net/http"
	"service/constant"
	"service/logger"
	"service/model"
	"service/translator"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Controller 基础控制器
type Controller struct{}

// Result 基础封装
func (c *Controller) Result(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, model.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功响应
func (c *Controller) Success(ctx *gin.Context, data interface{}) {
	c.Result(ctx, constant.SUCCESS, "请求成功", data)
}

// Error 失败响应
func (c *Controller) Error(ctx *gin.Context, msg string, data interface{}) {
	c.Result(ctx, constant.ERROR, msg, data)
}

// Valid 参数校验
func (c *Controller) Valid(ctx *gin.Context, valid interface{}) error {
	if err := ctx.ShouldBind(valid); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			logger.Error(ctx, "参数检验失败",
				zap.String("url", ctx.Request.URL.Path),
				zap.Any("validationErrors", errs.Translate(translator.Trans)),
			)
			c.Result(ctx, constant.VALID, "请求参数校验失败", c.removeTopStruct(errs.Translate(translator.Trans)))
		} else {
			logger.Error(ctx, "请求解析失败",
				zap.String("url", ctx.Request.URL.Path),
				zap.Any("error", err),
			)
			c.Result(ctx, http.StatusBadRequest, err.Error(), nil)
		}
		return err
	}
	return nil
}

func (c *Controller) removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
