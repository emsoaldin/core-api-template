package actions

import (
	"net/http"

	"github.com/go-flow/flow/v2"
)

type IndexAction struct {
}

func NewIndexAction() *IndexAction {
	return &IndexAction{}
}

func (a *IndexAction) Method() string {
	return http.MethodGet
}

func (a *IndexAction) Path() string {
	return "/"
}

func (a *IndexAction) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

func (a *IndexAction) Handle(r *http.Request) flow.Response {
	return flow.ResponseJSON(http.StatusOK, flow.Map{
		"msg:": "hello"})
}
