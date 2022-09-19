package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	zap.S().Debugf("验证手机号码是否符合格式")
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	return ok
}
