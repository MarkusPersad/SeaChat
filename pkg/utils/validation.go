package utils

import (
	"SeaChat/pkg/constants"
	"SeaChat/pkg/exception"
	"errors"
	"fmt"
	"reflect"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

var rules = map[string]func(fl validator.FieldLevel) bool{}

func init() {
	rules["pass"] = checkPassword
}

func checkPassword(fl validator.FieldLevel) bool {
	if isMatch,_ := regexp2.MustCompile(constants.PASSWORD_REGEX,regexp2.None).MatchString(fl.Field().String()); !isMatch {
		return false
	}
	return true
}

func processValidationErrors(entity any,err error) string {
	if err == nil {
		return ""
	}
	var invalid *validator.InvalidValidationError
	if errors.As(err, &invalid) {
		return fmt.Sprintf("输入参数错误: %s", invalid.Error())
	}
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, validationErr := range validationErrs{
			fieldName := validationErr.Field()
			typeOf := reflect.TypeOf(entity)
			if typeOf.Kind() == reflect.Pointer {
				typeOf = typeOf.Elem()
			}
			if field, ok := typeOf.FieldByName(fieldName); ok {
				errorInfo := field.Tag.Get(constants.FIELD_ERROR_INFO)
				return fmt.Sprintf("%s : %s", fieldName, errorInfo)
			} else {
				return "缺失字段错误信息"
			}
		}
	}
	return ""
}

func Validate(data any) error {
	valid := validator.New()
	for key,rule := range rules {
		if err :=valid.RegisterValidation(key, rule); err != nil {
			return err
		}
	}
	if errs := valid.Struct(data); errs != nil {
		return exception.New(constants.VALIDATION_ERROR,processValidationErrors(data, errs))
	}
	return nil
}
