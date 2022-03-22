package actions

import (
	"api/modules/account/services"
	"api/pkg/apperror"
	"api/pkg/userip"
	"api/providers/binding"
	"api/providers/vm"
	"errors"
	"net/http"

	"github.com/go-flow/flow/v2"
)

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginAction struct {
	vm             vm.Transformer
	binder         binding.Binder
	accountService services.AccountService
}

func NewLoginAction(vm vm.Transformer, binder binding.Binder, accountService services.AccountService) *LoginAction {
	return &LoginAction{
		vm:             vm,
		binder:         binder,
		accountService: accountService,
	}
}

func (a *LoginAction) Method() string {
	return http.MethodPost
}

func (a *LoginAction) Path() string {
	return "/login"
}

func (a *LoginAction) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

// Handle provides authentication tokens for user
// @Summary Login user and provides accesToken and refreshToken pair
// @Produce json
// @Tags account
// @Param req body Login true "Account Login Request"
// @Success 200 {object} models.Auth
// @Failure 400 {object} vm.ResponseError
// @Router /account/login [post]
func (a *LoginAction) Handle(r *http.Request) flow.Response {
	var reqObj Login
	if err := a.binder.Bind(r, &reqObj); err != nil {
		return a.vm.Error(http.StatusBadRequest, apperror.New("400", errors.New("validation error"), err))
	}

	ip := userip.Get(r)
	auth, err := a.accountService.Login(r.Context(), reqObj.Email, reqObj.Password, ip)
	if err != nil {
		return a.vm.Error(http.StatusBadRequest, err)
	}

	return a.vm.Success(http.StatusOK, auth)
}
