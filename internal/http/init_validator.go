package http

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"reflect"
)

func initValidator() *validator.Validate (
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
)