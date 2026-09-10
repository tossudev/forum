package handlers

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"reflect"
)

func initValidator() *validator.Validate (
	validator := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validator
)