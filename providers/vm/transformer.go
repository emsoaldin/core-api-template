package vm

import "github.com/go-flow/flow/v2"

type Transformer interface {
	Error(code int, err error) flow.Response
	Success(code int, data interface{}) flow.Response
}
