package checker

import (
	"fmt"
	errors "github.com/wuntsong/wterrors"
	"reflect"
)

var FieldCheckError = errors.NewClass("field check error")

func ReturnFieldError(field *reflect.StructField, msg string, args ...any) errors.WTError {
	if field == nil {
		return FieldCheckError.Errorf("check fail: %s", fmt.Sprintf(msg, args...))
	}
	return FieldCheckError.Errorf("field %s check fail: %s", field.Name, fmt.Sprintf(msg, args...))
}
