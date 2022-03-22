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

// Register request object
type Register struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=3"`
	FirstName string `json:"firstName" binding:"required,min=3"`
	LastName  string `json:"lastName" binding:"required,min=3"`
}

type RegisterAction struct {
	vm             vm.Transformer
	binder         binding.Binder
	accountService services.AccountService
}

func NewRegisterAction(vm vm.Transformer, binder binding.Binder, accountService services.AccountService) *RegisterAction {
	return &RegisterAction{
		vm:             vm,
		binder:         binder,
		accountService: accountService,
	}
}

func (a *RegisterAction) Method() string {
	return http.MethodPost
}

func (a *RegisterAction) Path() string {
	return "/register"
}

func (a *RegisterAction) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

// Handle Registers user to database and provide authentication tokens
// @Summary Register user to database and provides acces_token and refresh_token pair
// @Produce json
// @Tags account
// @Param req body Register true "Account Register Request"
// @Success 200 {object} models.Auth
// @Failure 400 {object} vm.ResponseError
// @Router /account/register [post]
func (a *RegisterAction) Handle(r *http.Request) flow.Response {
	var reqObj Register
	if err := a.binder.Bind(r, &reqObj); err != nil {
		return a.vm.Error(http.StatusBadRequest, apperror.New("400", errors.New("validation error"), err))
	}

	ip := userip.Get(r)

	auth, err := a.accountService.Register(r.Context(), reqObj.Email, reqObj.Password, reqObj.FirstName, reqObj.LastName, ip)
	if err != nil {
		return a.vm.Error(http.StatusBadRequest, err)
	}

	return a.vm.Success(http.StatusOK, auth)
}
