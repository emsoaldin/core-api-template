package actions

import (
	"net/http"

	"github.com/go-flow/flow/v2"
)

type HealthAction struct {
}

func NewHealthAction() *HealthAction {
	return &HealthAction{}
}

func (a *HealthAction) Method() string {
	return http.MethodGet
}

func (a *HealthAction) Path() string {
	return "/health"
}

func (a *HealthAction) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

// Handle dislays application health status
// @Summary Dislays application health status
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} vm.ResponseError
// @Router /health [get]
func (a *HealthAction) Handle(r *http.Request) flow.Response {
	return flow.ResponseJSON(http.StatusOK, map[string]interface{}{
		"msg:": "Health"})
}
