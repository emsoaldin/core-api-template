package services

import (
	"context"
	"errors"

	"api/modules/account/models"
	"api/modules/auth"
	"api/modules/roles"
	"api/modules/tokens"
	"api/modules/users/services"
	"api/pkg/apperror"

	"api/providers/jwt"
)

var (
	// ErrRegisterUser error is returned when user could not be registered
	ErrRegisterUser = errors.New("unable to register user")

	// ErrLoginUser error is returned when user could not be logged in
	ErrLoginUser = errors.New("unable to login user")

	// ErrRefreshTokens error is returned when auth tokens could not be refreshed
	ErrRefreshTokens = errors.New("unable to refresh tokens")
)

// AccountService interface
type AccountService interface {
	// AccountService returns interface implementation signature
	AccountService() string

	// Register new user to system using email and password combination
	Register(ctx context.Context, email string, password string, firstName string, lasatName string, clientIP string) (*models.Auth, error)

	// Login user to system using email and password combination
	Login(ctx context.Context, email string, password string, clientIP string) (*models.Auth, error)

	// RefreshToken issues new Auth Tokens based on given refreshToken
	RefreshToken(ctx context.Context, token string) (*models.Auth, error)
}

// NewAccountService creates AccountService Implementation
func NewAccountService(
	rolesService roles.RolesService,
	usersService services.UsersService,
	authService auth.AuthService,
	tokensService tokens.TokensService,
	jwt jwt.TokenAuth) AccountService {
	return &accountService{
		rolesService:  rolesService,
		usersService:  usersService,
		authService:   authService,
		tokensService: tokensService,
		jwt:           jwt,
	}
}

type accountService struct {
	rolesService  roles.RolesService
	usersService  services.UsersService
	authService   auth.AuthService
	tokensService tokens.TokensService
	jwt           jwt.TokenAuth
}

// AccountService returns Interface implementation signature
func (svc *accountService) AccountService() string {
	return "accountService"
}

func (svc *accountService) authenticate(ctx context.Context, userID uint64, roles ...string) (*models.Auth, error) {
	// create access token
	accessToken, err := svc.jwt.GenerateAccessToken(userID, roles...)
	if err != nil {
		return nil, err
	}

	// create and store refresh token
	token, err := svc.tokensService.CreateRefreshToken(ctx, userID, "")
	if err != nil {
		return nil, err
	}

	//generate refresh token
	refreshToken, err := svc.jwt.GenerateRefreshToken(token.Token)
	if err != nil {
		return nil, err
	}
	return &models.Auth{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Register new user to system using email and password combination
func (svc *accountService) Register(ctx context.Context, email string, password string, firstName string, lastName string, clientIP string) (*models.Auth, error) {
	defaultRole := uint64(roles.UserRoleUser)

	// create user
	user, err := svc.usersService.Create(ctx, firstName, lastName, email)
	if err != nil {
		return nil, apperror.New("ACCOUNT.000", ErrRegisterUser, err)
	}

	//assign default role
	if err := svc.rolesService.Assign(ctx, user.ID, defaultRole); err != nil {
		return nil, apperror.New("ACCOUNT.001", ErrRegisterUser, err)
	}

	// create local auth provider
	if err = svc.authService.CreateLocal(ctx, user.ID, password); err != nil {
		return nil, apperror.New("ACCOUNT.002", ErrRegisterUser, err)
	}

	// get user roles
	roles, err := svc.rolesService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.New("ACCOUNT.003", ErrRegisterUser, err)
	}

	// by default users do not have MFA enabled and access token should be authorized
	rolesArr := []string{"Authorized"}
	for _, role := range roles {
		rolesArr = append(rolesArr, role.Name)
	}

	// authenticate user
	auth, err := svc.authenticate(ctx, user.ID, rolesArr...)
	if err != nil {
		return nil, apperror.New("ACCOUNT.004", ErrRegisterUser, err)
	}

	return auth, nil
}

// Login user to system using email and password combination
func (svc *accountService) Login(ctx context.Context, email string, password string, clientIP string) (*models.Auth, error) {

	// get user by email
	user, err := svc.usersService.GetByEmail(ctx, email)

	if err != nil {
		return nil, apperror.New("ACCOUNT.010", ErrLoginUser, err)
	}

	// verify password
	if err := svc.authService.AuthenticateLocal(ctx, user.ID, password); err != nil {
		return nil, apperror.New("ACCOUNT.011", ErrLoginUser, err)
	}

	// get user roles
	roles, err := svc.rolesService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.New("ACCOUNT.012", ErrLoginUser, err)
	}

	var rolesArr []string
	for _, role := range roles {
		rolesArr = append(rolesArr, role.Name)
	}

	// authenticate user
	return svc.authenticate(ctx, user.ID, rolesArr...)
}

// RefreshToken issues new Auth Tokens based on given refreshToken
func (svc *accountService) RefreshToken(ctx context.Context, token string) (*models.Auth, error) {
	tokenObj, err := svc.tokensService.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, apperror.New("ACCOUNT.020", ErrRefreshTokens, err)
	}

	user, err := svc.usersService.GetByID(ctx, tokenObj.UserID)
	if err != nil {
		return nil, apperror.New("ACCOUNT.021", ErrRefreshTokens, err)
	}

	// get user roles
	roles, err := svc.rolesService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.New("ACCOUNT.022", ErrRefreshTokens, err)
	}

	var rolesArr []string
	for _, role := range roles {
		rolesArr = append(rolesArr, role.Name)
	}

	rolesArr = append(rolesArr, "Authorized")

	// delete old token
	if err := svc.tokensService.Delete(ctx, tokenObj); err != nil {
		return nil, apperror.New("ACCOUNT.023", ErrRefreshTokens, err)
	}

	// authenticate user
	return svc.authenticate(ctx, user.ID, rolesArr...)
}
