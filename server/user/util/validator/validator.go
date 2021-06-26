/**
 * @Author fzh
 * @Date 2021/6/25
 */
package validator

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// 初始化验证器
func InitValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		logging.Error("验证器初始化失败")
		return
	}

	var err error
	add := func(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
		if err != nil {
			return
		}
		err = v.RegisterValidation(tag, fn, callValidationEvenIfNull...)
	}

	// 添加验证器
	add("phoneNumber", phoneNumber)
	add("certificateNumber", certificateNumber)

	if err != nil {
		logging.Error("验证器加载失败")
	}
}

func phoneNumber(fl validator.FieldLevel) bool {
	if number, ok := fl.Field().Interface().(string); ok {
		if matched, _ := regexp.MatchString(`^1[3456789]\d{9}$`, number); matched {
			return true
		}
	}
	return false
}

func certificateNumber(fl validator.FieldLevel) bool {
	if number, ok := fl.Field().Interface().(string); ok {
		if matched, _ := regexp.MatchString(`^\d{17}[0-9xX]$`, number); matched {
			return true
		}
	}
	return false
}
