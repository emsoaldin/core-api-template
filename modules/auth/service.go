package auth

import (
	"context"
	"errors"

	"api/pkg/apperror"

	"golang.org/x/crypto/bcrypt"
)

const (
	// AuthLocal -
	AuthLocal = "local"

	// AuthGoogle -
	AuthGoogle = "google"

	// AuthProviderFacebook -
	AuthFacebook = "facebook"

	// AuthProviderApple -
	AuthApple = "apple"
)

var (

	// ErrUnsupportedProvider error is returned when unsupported provider name was used to get/set auth
	ErrUnsupportedProvider = errors.New("unsupported provider")

	// ErrCreateLocalAuth error is returned when local auth provider could not be created
	ErrCreateLocalAuth = errors.New("unable to create local auth provider")

	// ErrCreateSocialAuth error is returned when social auth provider could not be created
	ErrCreateSocialAuth = errors.New("unable to create social auth provider")

	// ErrResetPassword error is returned when password could not be changed
	ErrResetPassword = errors.New("unable to reset password")

	// ErrLocalAuthNotExist error is returned when local auth provider does not exist
	ErrLocalAuthNotExist = errors.New("local auth provider does not exist")

	// ErrSocialAuthNotExist error is returned when social auth provider does not exist
	ErrSocialAuthNotExist = errors.New("social auth provider does not exist")

	// ErrAuthentication error is returned when uset authentication fails
	ErrAuthentication = errors.New("unable to authenticate user")

	// ErrPasswordMatch error is returned when passwords do not match
	ErrPasswordMatch = errors.New("passwords don't match")

	// ErrCreatePassword error is returned when password hash could not be created
	ErrCreatePassword = errors.New("unable to create password hash")

	// ErrFetchAuthProviders error is returned when auth providers could not be fetched
	ErrFetchAuthProviders = errors.New("unable to get auth providers")

	// ErrDeleteAuthProviders error is returned when auth providers could not be deleted
	ErrDeleteAuthProviders = errors.New("unable to delete auth providers")
)

// AuthService interface
type AuthService interface {
	// AuthService returns service implementation signature
	AuthService() string

	// CreateLocal creates local (password) authentication strategy for user
	CreateLocal(ctx context.Context, userID uint64, password string) error

	// CreateSocial creates social authentication  strategy for given user and social provider
	CreateSocial(ctx context.Context, userID uint64, uid, providerName string) error

	// ResetLocal  sets new password  for local authentication strategy
	ResetLocal(ctx context.Context, userID uint64, password string) error

	// AuthenticateLocal authenticates user for local (password) strategy
	AuthenticateLocal(ctx context.Context, userID uint64, password string) error

	// AuthenticateSocial authenticates user for social strategy
	AuthenticateSocial(ctx context.Context, userID uint64, uid, provider string) error

	// ComparePassword compares password against given hash
	ComparePassword(ctx context.Context, passwordHash string, password string) error

	// GeneratePasswordHash generates hash for given password
	GeneratePasswordHash(ctx context.Context, password string) (string, error)

	// GetByUserID returns all authentication strategies for given user
	GetByUserID(ctx context.Context, userID uint64) ([]string, error)

	// DeleteAll removes all authentication strategies for given user
	DeleteAll(ctx context.Context, userID uint64) error
}

// NewAuthService creates AuthService interface implementation
func NewAuthService(authRepository AuthRepository) AuthService {

	return &authService{
		repo: authRepository,
	}
}

type authService struct {
	repo AuthRepository
}

// AuthService returns service implementation signature
func (authService) AuthService() string {
	return "authService"
}

// CreateLocal creates local (password) authentication strategy for user
func (svc *authService) CreateLocal(ctx context.Context, userID uint64, password string) error {
	hash, err := svc.GeneratePasswordHash(ctx, password)
	if err != nil {
		return apperror.New("AUTH.000", ErrCreateLocalAuth, err)
	}

	localProvider := &AuthProvider{
		UserID:   userID,
		Provider: AuthLocal,
		Hash:     hash,
	}

	if err = svc.repo.Create(ctx, localProvider); err != nil {
		return apperror.New("AUTH.001", ErrCreateLocalAuth, err)
	}

	return nil
}

// CreateSocial creates social authentication  strategy for given user and social provider
func (svc *authService) CreateSocial(ctx context.Context, userID uint64, hash, providerName string) error {
	if providerName != AuthApple && providerName != AuthFacebook && providerName != AuthGoogle {
		return apperror.New("AUTH.010", ErrCreateSocialAuth, ErrUnsupportedProvider)
	}

	provider := &AuthProvider{
		UserID:   userID,
		Provider: providerName,
		Hash:     hash,
	}

	if err := svc.repo.Create(ctx, provider); err != nil {
		return apperror.New("AUTH.011", ErrCreateSocialAuth, err)
	}

	return nil
}

// ResetLocal  sets new password  for local authentication strategy
func (svc *authService) ResetLocal(ctx context.Context, userID uint64, password string) error {
	provider, err := svc.repo.GetByID(ctx, AuthLocal, userID)
	if err != nil {
		return apperror.New("AUTH.020", ErrResetPassword, err)
	}

	if provider == nil {
		return apperror.New("AUTH.021", ErrResetPassword, ErrLocalAuthNotExist)
	}

	// generate new Hash
	provider.Hash, err = svc.GeneratePasswordHash(ctx, password)
	if err != nil {
		return apperror.New("AUTH.022", ErrResetPassword, err)
	}

	if err = svc.repo.Update(ctx, provider); err != nil {
		return apperror.New("AUTH.023", ErrResetPassword, err)
	}

	return nil
}

// AuthenticateLocal authenticates user for local (password) strategy
func (svc *authService) AuthenticateLocal(ctx context.Context, userID uint64, password string) error {
	provider, err := svc.repo.GetByID(ctx, AuthLocal, userID)
	if err != nil {
		return apperror.New("AUTH.080", ErrAuthentication, err)
	}

	if provider == nil {
		return apperror.New("AUTH.090", ErrAuthentication, ErrLocalAuthNotExist)
	}

	if err := svc.ComparePassword(ctx, provider.Hash, password); err != nil {
		return apperror.New("AUTH.100", ErrAuthentication, err)
	}

	return err
}

// AuthenticateSocial authenticates user for social strategy
func (svc *authService) AuthenticateSocial(ctx context.Context, userID uint64, hash, providerName string) error {
	if providerName != AuthApple && providerName != AuthFacebook && providerName != AuthGoogle {
		return apperror.New("AUTH.110", ErrAuthentication, ErrUnsupportedProvider)
	}

	provider, err := svc.repo.GetByID(ctx, providerName, userID)
	if err != nil {
		return apperror.New("AUTH.120", ErrAuthentication, err)
	}

	if provider == nil {
		return apperror.New("AUTH.130", ErrAuthentication, ErrSocialAuthNotExist)
	}

	if err = svc.ComparePassword(ctx, provider.Hash, hash); err != nil {
		return apperror.New("AUTH.140", ErrAuthentication, err)
	}

	return err
}

// ComparePassword compares password against given hash
func (svc *authService) ComparePassword(ctx context.Context, passwordHash string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return apperror.New("AUTH.150", ErrPasswordMatch, err)
	}
	return nil
}

// GeneratePasswordHash generates hash for given password
func (svc *authService) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperror.New("AUTH.160", ErrCreatePassword, err)
	}
	return string(hash), nil
}

// GetByUserID returns all authentication strategies for given user
func (svc *authService) GetByUserID(ctx context.Context, userID uint64) ([]string, error) {
	providers, err := svc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.New("AUTH.170", ErrFetchAuthProviders, err)
	}
	str := []string{}

	for _, p := range providers {
		str = append(str, p.Provider)
	}

	return str, nil
}

// DeleteAll removes all authentication strategies for given user
func (svc *authService) DeleteAll(ctx context.Context, userID uint64) error {
	if err := svc.repo.DeleteByUserID(ctx, userID); err != nil {
		return apperror.New("AUTH.180", ErrDeleteAuthProviders, err)
	}
	return nil
}
