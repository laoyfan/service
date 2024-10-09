package translator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	znTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator

func InitTranslator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		uni := ut.New(zhT, zhT)

		var o bool
		Trans, o = uni.GetTranslator("zh")
		if !o {
			return fmt.Errorf("翻译器获取失败")
		}

		if err := znTranslations.RegisterDefaultTranslations(v, Trans); err != nil {
			return fmt.Errorf("翻译器注册失败:%w", err)
		}
	}
	return nil
}
