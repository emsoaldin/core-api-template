package vm

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"api/pkg/apperror"

	"github.com/go-flow/flow/v2"
	"github.com/go-playground/validator/v10"
)

// NewJson creates JSON response transformer
func NewJson() Transformer {
	return &jsonTransformer{}
}

type jsonTransformer struct{}

// Error transforms given error to standardised JSON response object
func (t *jsonTransformer) Error(code int, err error) flow.Response {

	vm := t.buildError(err)

	return t.response(false, code, nil, vm)
}

// Success transforms given data to standardised JSON response object
func (t *jsonTransformer) Success(code int, data interface{}) flow.Response {
	return t.response(true, code, data, nil)
}

func (jsonTransformer) response(success bool, code int, data interface{}, err interface{}) flow.Response {
	return flow.ResponseJSON(code, &Response{
		Success: success,
		Data:    data,
		Error:   err,
	})
}

func (t *jsonTransformer) buildError(err error) *ResponseError {
	vm := &ResponseError{
		Message: err.Error(),
	}

	// check if error is apperror
	if aErr, ok := err.(*apperror.Error); ok {
		vm.Code = aErr.Code()
		if aErr.Cause() != nil {
			if tErr, ok := aErr.Cause().(*apperror.Error); ok {
				vm.Cause = t.buildError(tErr)
			} else {
				vm.Cause = &ResponseError{
					Message: aErr.Cause().Error(),
				}
			}

		}
		if os.Getenv("ENV") != "production" {
			vm.Stack = aErr.Trace()
		}

	}

	if wErr := errors.Unwrap(err); wErr != nil {
		// check if httpError is caused by validation
		if verrs, ok := wErr.(validator.ValidationErrors); ok {
			m := map[string]string{}

			for _, verr := range verrs {
				m[verr.Field()] = fmt.Sprintf("%s_%s", strings.ToLower(verr.Field()), verr.Tag())
			}

			vm.Validation = m
		}
	}

	return vm
}
