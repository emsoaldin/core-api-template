package tokens

import (
	"context"
	"errors"
	"time"

	"api/pkg/apperror"

	"github.com/google/uuid"
)

const (
	// TokenTypePasswordReset holds db ID value for Password Reset
	TokenTypePasswordReset = 1

	// TokenTypeEmailConfirmation holds db ID value for Email configuration
	TokenTypeEmailConfirmation = 2

	// TokenTypeInvitation holds db ID value for Invitation
	TokenTypeInvitation = 3

	// TokenTypeRefresh holds db ID value for Refresh token
	TokenTypeRefresh = 4

	// TokenTypeDeleteAccount holds db ID value for Delete Account token
	TokenTypeDeleteAccount = 5

	// PasswordResetTokenDuration holds duration value in minutes for pasword reset token
	PasswordResetTokenDuration = 30

	// EmailConfirmationTokenDuration holds duration value in minutes for email confirmation token
	EmailConfirmationTokenDuration = 21600 // 15 days

	// InvitationTokenDuration holds duration value in minutes for invitation token
	InvitationTokenDuration = 43200 //30 days

	// RefreshTokenDuration holds duration value in minutes for refresh token
	RefreshTokenDuration = 43200 // 30 days

	// DeleteAccountTokenDuration holds duration value in minutes for Delete account token
	DeleteAccountTokenDuration = 15
)

var (

	// ErrWrongTkenType error is returned when token type missmatch is detected while fetching token
	ErrWrongTkenType = errors.New("wrong token type")

	// ErrCreatePasswordResetToken error is returned when new password reset token could not be created
	ErrCreatePasswordResetToken = errors.New("unable to create password reset token")

	// ErrFetchPasswordresetToken error is returned when password reset token could not be retrieved
	ErrFetchPasswordresetToken = errors.New("unable to fetch password reset token")

	// ErrCreateEmailConfirmToken error is returned when new email confirmation token could not be created
	ErrCreateEmailConfirmToken = errors.New("unable to create email confirm token")

	// ErrFetchEmailConfirmationToken error is returned when email confirmation  token could not be retrieved
	ErrFetchEmailConfirmationToken = errors.New("unable to fetch email confirmation token")

	// ErrCreateInvitationToken error is returned when new invitation token could not be created
	ErrCreateInvitationToken = errors.New("unable to create invitation token")

	// ErrFetchInviteToken error is returned when invite token could not be retrieved
	ErrFetchInviteToken = errors.New("unable to fetch invite token")

	// ErrCreateRefreshToken error is returned when refrech token could not be created
	ErrCreateRefreshToken = errors.New("unable to create refresh token")

	// ErrFetchRefreshToken error is returned when refresh token could not be retrieved
	ErrFetchRefreshToken = errors.New("unable to fetch refresh token")

	// ErrCreateDeleteAccountToken error is returned when delete account token could not be created
	ErrCreateDeleteAccountToken = errors.New("unable to create delete account token")

	// ErrFetchDeleteAccountToken error is returned when delete account token could not be retrieved
	ErrFetchDeleteAccountToken = errors.New("unable to fetch delete account token")

	// ErrFetchToken error is returned when token could not be retrieved
	ErrFetchToken = errors.New("unable to fetch token")

	// ErrDeleteToken error is returned when token could not be deleted
	ErrDeleteToken = errors.New("unable to delete token")

	// ErrTokenNotExist error is returned when token can not be fetched from database for given criteria
	ErrTokenNotExist = errors.New("token does not exist")

	// ErrExpiredToken error is returned when fetched token is expired
	ErrExpiredToken = errors.New("token expired")

	// ErrDeleteUserTokens error is returned when user tokens could not be deleted
	ErrDeleteUserTokens = errors.New("unable to delete user tokens")

	// ErrDeleteExpiredTokens error is returned when expired tokens could not be deleted
	ErrDeleteExpiredTokens = errors.New("unable to delete expired tokens")
)

// TokensService interface
type TokensService interface {

	// TokensService interface implementation signature
	TokensService() string

	// CreatePasswordRessetToken creates password reset token for given user
	// all previous password reset tokens are deleted when new token is created
	CreatePasswordRessetToken(ctx context.Context, userID uint64, meta string) (*Token, error)

	// GetPasswordRessetToken retrieves password reset token
	GetPasswordRessetToken(ctx context.Context, token string) (*Token, error)

	// CreateEmailConfirmToken creates email confirmation token for given user
	// all previous email confirmation tokens are deleted when new token is created
	CreateEmailConfirmToken(ctx context.Context, userID uint64, meta string) (*Token, error)

	// GetEmailConfirmationToken retrieves email confirmation token
	GetEmailConfirmToken(ctx context.Context, token string) (*Token, error)

	// CreateInviteToken creates user invitation token for given user
	// all previous invitation tokens are removed for that user
	CreateInviteToken(ctx context.Context, userID uint64, meta string) (*Token, error)

	// GetInviteToken retrieves user invitation token
	GetInviteToken(ctx context.Context, token string) (*Token, error)

	// CreateRefreshToken creates refresh token for given user
	// all previous refresh tokens for given user are removed
	CreateRefreshToken(ctx context.Context, userID uint64, meta string) (*Token, error)

	// GetrefreshToken retrieves refresh token
	GetRefreshToken(ctx context.Context, token string) (*Token, error)

	// CreateDeleteAccountToken creates delete account token for given user
	// all previous delete account tokens for given user are removed
	CreateDeleteAccountToken(ctx context.Context, userID uint64, meta string) (*Token, error)

	// GetDeleteAccountToken retrieves delete account
	GetDeleteAccountToken(ctx context.Context, token string) (*Token, error)

	// GetByID returns Token from database with provided id
	GetByID(ctx context.Context, id uint64) (*Token, error)

	// Delete token
	Delete(ctx context.Context, token *Token) error

	// DeleteByID removes token with provided id
	DeleteByID(ctx context.Context, id uint64) error

	// GetByToken returns token object for provided token string
	GetByToken(ctx context.Context, token string) (*Token, error)

	// DeleteByUserID removes all tokens for given userID
	DeleteByUserID(ctx context.Context, userID uint64) error

	// DeleteExpiredTokens removes all tokens that are expired
	DeleteExpiredTokens(ctx context.Context) error
}

// NewTokensService creates TokensService interface implementation
func NewTokensService(tokensRepository TokensRepository) TokensService {
	return &tokensService{
		repo: tokensRepository,
	}
}

type tokensService struct {
	repo TokensRepository
}

func (tokensService) TokensService() string {
	return "tokensService"
}

func (svc *tokensService) create(ctx context.Context, userID uint64, tokenTypeID uint64, token string, meta string, expiresAt time.Time) (*Token, error) {
	t := &Token{
		UserID:      userID,
		TokenTypeID: tokenTypeID,
		Token:       token,
		Meta:        meta,
		ExpiresAt:   expiresAt,
	}

	if token == "" {
		t.Token = uuid.New().String()
	}

	if meta == "" {
		t.Meta = "{}"
	}

	return t, svc.repo.Create(ctx, t)
}

// CreatePasswordRessetToken creates password reset token for given user
// all previous password reset tokens are deleted when new token is created
func (svc *tokensService) CreatePasswordRessetToken(ctx context.Context, userID uint64, meta string) (*Token, error) {

	if err := svc.repo.DeleteByUserAndTokenTypeID(ctx, userID, TokenTypePasswordReset); err != nil {
		return nil, apperror.New("TOKENS.000", ErrCreatePasswordResetToken, err)
	}

	exp := time.Now().Add(time.Minute * time.Duration(PasswordResetTokenDuration))
	t, err := svc.create(ctx, userID, TokenTypePasswordReset, "", "", exp)
	if err != nil {
		return nil, apperror.New("TOKENS.001", ErrCreatePasswordResetToken, err)
	}
	return t, nil
}

// GetPasswordRessetToken retrieves password reset token
func (svc *tokensService) GetPasswordRessetToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.TokenTypeID != TokenTypePasswordReset {
		return nil, apperror.New("TOKENS.010", ErrFetchPasswordresetToken, ErrWrongTkenType)
	}

	return t, nil
}

// CreateEmailConfirmToken creates email confirmation token for given user
// all previous email confirmation tokens are deleted when new token is created
func (svc *tokensService) CreateEmailConfirmToken(ctx context.Context, userID uint64, meta string) (*Token, error) {

	if err := svc.repo.DeleteByUserAndTokenTypeID(ctx, userID, TokenTypeEmailConfirmation); err != nil {
		return nil, apperror.New("TOKENS.020", ErrCreateEmailConfirmToken, err)
	}

	exp := time.Now().Add(time.Minute * time.Duration(EmailConfirmationTokenDuration))
	t, err := svc.create(ctx, userID, TokenTypeEmailConfirmation, "", "", exp)
	if err != nil {
		return nil, apperror.New("TOKENS.021", ErrCreateEmailConfirmToken, err)
	}
	return t, nil
}

// GetEmailConfirmationToken retrieves email confirmation token
func (svc *tokensService) GetEmailConfirmToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.TokenTypeID != TokenTypeEmailConfirmation {
		return nil, apperror.New("TOKENS.030", ErrFetchEmailConfirmationToken, ErrWrongTkenType)
	}

	return t, nil
}

// CreateInviteToken creates user invitation token for given user
// all previous invitation tokens are removed for that user
func (svc *tokensService) CreateInviteToken(ctx context.Context, userID uint64, meta string) (*Token, error) {
	if err := svc.repo.DeleteByUserAndTokenTypeID(ctx, userID, TokenTypeInvitation); err != nil {
		return nil, apperror.New("TOKENS.040", ErrCreateInvitationToken, err)
	}

	exp := time.Now().Add(time.Minute * time.Duration(InvitationTokenDuration))
	t, err := svc.create(ctx, userID, TokenTypeInvitation, "", "", exp)
	if err != nil {
		return nil, apperror.New("TOKENS.041", ErrCreateInvitationToken, err)
	}
	return t, nil
}

// GetInviteToken retrieves user invitation token
func (svc *tokensService) GetInviteToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.TokenTypeID != TokenTypeInvitation {
		return nil, apperror.New("TOKENS.050", ErrFetchInviteToken, ErrWrongTkenType)
	}

	return t, nil
}

// CreateRefreshToken creates refresh token for given user
// all previous refresh tokens for given user are removed
func (svc *tokensService) CreateRefreshToken(ctx context.Context, userID uint64, meta string) (*Token, error) {
	// if err := svc.repo.DeleteByUserAndTokenTypeID(ctx, userID, TokenTypeRefresh); err != nil {
	//	 return nil, apperror.New("040.090", ErrCreateRefreshToken, err)
	// }

	exp := time.Now().Add(time.Minute * time.Duration(RefreshTokenDuration))
	t, err := svc.create(ctx, userID, TokenTypeRefresh, "", "", exp)
	if err != nil {
		return nil, apperror.New("TOKENS.060", ErrCreateRefreshToken, err)
	}
	return t, nil
}

// GetrefreshToken retrieves refresh token
func (svc *tokensService) GetRefreshToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.TokenTypeID != TokenTypeRefresh {
		return nil, apperror.New("TOKENS.070", ErrFetchRefreshToken, ErrWrongTkenType)
	}

	return t, nil
}

// CreateDeleteAccountToken creates delete account token for given user
// all previous delete account tokens for given user are removed
func (svc *tokensService) CreateDeleteAccountToken(ctx context.Context, userID uint64, meta string) (*Token, error) {

	if err := svc.repo.DeleteByUserAndTokenTypeID(ctx, userID, TokenTypeDeleteAccount); err != nil {
		return nil, apperror.New("TOKENS.080", ErrCreateDeleteAccountToken, err)
	}

	exp := time.Now().Add(time.Minute * time.Duration(DeleteAccountTokenDuration))
	t, err := svc.create(ctx, userID, TokenTypeDeleteAccount, "", "", exp)
	if err != nil {
		return nil, apperror.New("TOKENS.081", ErrCreateDeleteAccountToken, err)
	}
	return t, nil
}

// GetDeleteAccountToken retrieves delete account token
func (svc *tokensService) GetDeleteAccountToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.TokenTypeID != TokenTypeDeleteAccount {
		return nil, apperror.New("TOKENS.090", ErrFetchDeleteAccountToken, ErrWrongTkenType)
	}

	return t, nil
}

// GetByID returns Token from database with provided id
func (svc *tokensService) GetByID(ctx context.Context, id uint64) (*Token, error) {
	t, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.New("TOKENS.100", ErrFetchToken, err)
	}

	return t, nil
}

// Delete token
func (svc *tokensService) Delete(ctx context.Context, token *Token) error {
	if err := svc.repo.Delete(ctx, token); err != nil {
		return apperror.New("TOKENS.110", ErrDeleteToken, err)
	}

	return nil
}

// DeleteByID removes token with provided id
func (svc *tokensService) DeleteByID(ctx context.Context, id uint64) error {
	if err := svc.repo.DeleteByID(ctx, id); err != nil {
		return apperror.New("TOKENS.120", ErrDeleteToken, err)
	}

	return nil
}

// GetByToken returns token object for provided token string
func (svc *tokensService) GetByToken(ctx context.Context, token string) (*Token, error) {
	t, err := svc.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, apperror.New("TOKENS.130", ErrFetchToken, err)
	}

	if t == nil {
		return nil, apperror.New("TOKENS.131", ErrFetchToken, ErrTokenNotExist)
	}

	if time.Now().After(t.ExpiresAt) {
		return nil, apperror.New("TOKENS.132", ErrFetchToken, ErrExpiredToken)
	}

	return t, nil
}

// DeleteByUserID removes all tokens for given userID
func (svc *tokensService) DeleteByUserID(ctx context.Context, userID uint64) error {
	if err := svc.repo.DeleteByUserID(ctx, userID); err != nil {
		return apperror.New("TOKENS.140", ErrDeleteUserTokens, err)
	}
	return nil
}

// DeleteExpiredTokens removes all tokens that are expired
func (svc *tokensService) DeleteExpiredTokens(ctx context.Context) error {
	if err := svc.repo.DeleteExpiredTokens(ctx); err != nil {
		return apperror.New("TOKENS.150", ErrDeleteExpiredTokens, err)
	}
	return nil
}
