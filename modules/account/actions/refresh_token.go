package actions

import (
	"api/modules/account/services"
	"api/pkg/apperror"
	"api/providers/binding"
	"api/providers/jwt"
	"api/providers/vm"
	"errors"
	"net/http"

	"github.com/go-flow/flow/v2"
)

type RefreshToken struct {
	Token string `json:"token" binding:"required"`
}

type RefreshTokenAction struct {
	vm             vm.Transformer
	auth           jwt.TokenAuth
	binder         binding.Binder
	accountService services.AccountService
}

func NewRefreshTokenAction(vm vm.Transformer, auth jwt.TokenAuth, binder binding.Binder, accountService services.AccountService) *RefreshTokenAction {
	return &RefreshTokenAction{
		vm:             vm,
		auth:           auth,
		binder:         binder,
		accountService: accountService,
	}
}

func (a *RefreshTokenAction) Method() string {
	return http.MethodPost
}

func (a *RefreshTokenAction) Path() string {
	return "/refresh-token"
}

func (a *RefreshTokenAction) Middlewares() []flow.MiddlewareHandlerFunc {
	return []flow.MiddlewareHandlerFunc{}
}

// Handle issues new tokens pair
// @Summary This action is used when accessToken is expired, using refresh token new token pair is generated
// and user can continue using service without interruption.
// @Produce json
// @Tags account
// @Param req body RefreshToken true "Refresh Token Request"
// @Success 200 {object} models.Auth
// @Failure 400 {object} vm.ResponseError
// @Router /account/refresh-token [post]
func (a *RefreshTokenAction) Handle(r *http.Request) flow.Response {
	var reqObj RefreshToken

	if err := a.binder.Bind(r, &reqObj); err != nil {
		return a.vm.Error(http.StatusBadRequest, apperror.New("400", errors.New("validation error"), err))
	}

	// verify refresh token
	token, err := a.auth.VerifyRefreshToken(reqObj.Token)
	if err != nil {
		return a.vm.Error(http.StatusBadRequest, apperror.New("400", errors.New("validation error"), err))
	}

	auth, err := a.accountService.RefreshToken(r.Context(), token)
	if err != nil {
		return a.vm.Error(http.StatusBadRequest, err)
	}

	return a.vm.Success(http.StatusOK, auth)
}
