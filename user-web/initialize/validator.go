package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"shop-api/user-web/global"
	"strings"
)

func InitTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s", local)
		}
		switch local {
		case "en":
			err := en_translations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		case "zh":
			err := zh_translations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		default:
			err := en_translations.RegisterDefaultTranslations(v, global.Trans)
			if err != nil {
				return err
			}
		}
		return
	}
	return
}
