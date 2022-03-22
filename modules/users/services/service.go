package services

import (
	"context"
	"errors"

	"api/modules/users/models"
	"api/modules/users/repositories"
	"api/pkg/apperror"
	"api/pkg/paging"
)

var (
	// ErrUserExists error is returned when we want to create new user, but the user already exists
	// user existence is defined by unique fields in database for the user - email
	ErrUserExists = errors.New("user already exists")

	// ErrUserNotExist error is returned when user can not be fetched from database for given criteria
	ErrUserNotExist = errors.New("user does not exist")

	// ErrCreateUser error is returned when new user can not be created in database
	ErrCreateUser = errors.New("unable to create user")

	// ErrFetchUser error is returned when user can not be retrieved from database
	ErrFetchUser = errors.New("unable to fetch user")

	// ErrUpdateUser error is returned when user can not be updated in database
	ErrUpdateUser = errors.New("unable to update user")

	// ErrFetchUsers error is returned when users can not be retrieved from database
	ErrFetchUsers = errors.New("unable to fetch users")

	// ErrInvalidUserState error is returned when user status is not valid enumeration
	ErrInvalidUserState = errors.New("invalid user state")

	ErrDeleteUser = errors.New("unable to delete user")
)

// UsersService interface
type UsersService interface {
	// UsersService implementation signature
	UsersService() string

	// Create user
	Create(ctx context.Context, firstName string, lastName string, email string) (*models.User, error)

	// GetByEmail returns user by given email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByID returns user by given user id
	GetByID(ctx context.Context, id uint64) (*models.User, error)

	// Find retrieves all users for given filter and pagination params
	Find(ctx context.Context, userType string, paginator *paging.Paginator) ([]*models.User, error)

	// Update updates user profile information
	Update(ctx context.Context, id uint64, firstName *string, lastName *string) error
}

// NewUsersService creates UsersService interface implementation
func NewUsersService(usersRepository repositories.UsersRepository) UsersService {
	return &usersService{
		repo: usersRepository,
	}
}

type usersService struct {
	repo repositories.UsersRepository
}

func (usersService) UsersService() string {
	return "usersService"
}

// Create user
func (svc *usersService) Create(ctx context.Context, firstName string, lastName string, email string) (*models.User, error) {
	user, err := svc.GetByEmail(ctx, email)

	if err != nil && !errors.Is(err, ErrUserNotExist) {
		return nil, err
	}

	if user != nil {
		return nil, apperror.New("USERS.000", ErrUserExists)
	}

	user = &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err = svc.repo.Create(ctx, user)
	if err != nil {
		return nil, apperror.New("USERS.001", ErrCreateUser, err)
	}

	return user, err
}

// GetByEmail returns user by given email
func (svc *usersService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := svc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.New("USERS.010", ErrFetchUser, err)
	}

	if user == nil {
		return nil, apperror.New("USERS.011", ErrFetchUser, ErrUserNotExist)
	}

	return user, nil
}

// GetByID returns user by given user id
func (svc *usersService) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	user, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.New("USERS.020", ErrFetchUser, err)
	}

	if user == nil {
		return nil, apperror.New("USERS.021", ErrFetchUser, ErrUserNotExist)
	}

	return user, nil
}

// Find retrieves all users for given filter and pagination params
func (svc *usersService) Find(ctx context.Context, userState string, paginator *paging.Paginator) ([]*models.User, error) {
	users, err := svc.repo.GetAll(ctx, paginator.Filter, paginator.Page, paginator.PerPage, paginator.OrderBy, paginator.OrderDir)
	if err != nil {
		return nil, apperror.New("USERS.030", ErrFetchUsers, err)
	}

	count, err := svc.repo.Count(ctx, paginator.Filter)
	if err != nil {
		return nil, apperror.New("USERS.031", ErrFetchUsers, err)
	}

	paginator.TotalEntriesSize = count
	paginator.CurrentEntriesSize = len(users)
	paginator.TotalPages = paginator.TotalEntriesSize / paginator.PerPage
	if paginator.TotalEntriesSize%paginator.PerPage > 0 {
		paginator.TotalPages = paginator.TotalPages + 1
	}
	return users, nil
}

// Update updates user profile information
func (svc *usersService) Update(ctx context.Context, id uint64, firstName *string, lastName *string) error {
	user, err := svc.GetByID(ctx, id)
	if err != nil {
		return apperror.New("USERS.040", ErrUpdateUser, err)
	}

	if firstName != nil {
		user.FirstName = *firstName
	}

	if lastName != nil {
		user.LastName = *lastName
	}

	if err := svc.repo.Update(ctx, user); err != nil {
		return apperror.New("USERS.041", ErrUpdateUser, err)
	}
	return nil
}
