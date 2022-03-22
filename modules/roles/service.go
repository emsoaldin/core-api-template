package roles

import (
	"context"
	"errors"
	"math"

	"api/pkg/apperror"
)

// UserRole enum
type UserRole uint64

const (
	// UserRoleUnknown unknown
	UserRoleUnknown UserRole = math.MaxUint64
	// UserRoleAdmin admin role
	UserRoleAdmin UserRole = 1
	// UserRoleClient client role
	UserRoleClient UserRole = 2
	// UserRoleUser user role
	UserRoleUser UserRole = 3
)

var (
	// ErrSaveRole error is returned when role object could not be saved
	ErrSaveRole = errors.New("unable to save role")

	// ErrNilRole error is returned when role object is nil where it should not be
	ErrNilRole = errors.New("role object is nil")

	// ErrUpdateRole error is returned when role object could not be updated
	ErrUpdateRole = errors.New("unable to update role")

	// ErrCreateRole error is returned when new role object could not be created
	ErrCreateRole = errors.New("unable to create role")

	// ErrFetchRole error is returned when role object could not be retrieved
	ErrFetchRole = errors.New("unable to fetch role")

	// ErrFetchRoles error is returned when roles collection could not be retrieved
	ErrFetchRoles = errors.New("unable to fetch roles")

	// ErrDeleteRole error is returned when role could not be deleted
	ErrDeleteRole = errors.New("unable to delete role")

	// ErrAssignRole error is returned when role could not be assigned
	ErrAssignRole = errors.New("unable to assign role")

	// ErrUnassignRole error is returned when role could not be un-assigned
	ErrUnassignRole = errors.New("unable to un-assign role")
)

// RolesService interface
type RolesService interface {
	// RolesService returns service interface
	RolesService() string

	// Count returns number of records in database
	Count(ctx context.Context) (int, error)

	// Save saves given role object
	Save(ctx context.Context, role *Role) error

	// Update updates role object
	Update(ctx context.Context, role *Role) error

	// Create role
	Create(ctx context.Context, role *Role) error

	// GetByID returns role object from database for given id
	GetByID(ctx context.Context, id uint64) (*Role, error)

	// GetAll returns all Role objects for given params
	GetAll(ctx context.Context, page int, perPage int, orderBy string, orderDir string) ([]*Role, error)

	// Delete removes role from database
	Delete(ctx context.Context, role *Role) error

	// GetByUserID returns all roles assigned to user
	GetByUserID(ctx context.Context, userID uint64) ([]*Role, error)

	// Assign assigns role to user
	Assign(ctx context.Context, userID uint64, roleID uint64) error

	// Unassign unnasigns role from user
	Unassign(ctx context.Context, userID uint64, roleID uint64) error

	// GetUsersRole returns user role. If user has more roles, will return role with most privilegies
	GetUserRole(ctx context.Context, userID uint64) (UserRole, error)
}

// NewRolesService creates RolesService interface implementation
func NewRolesService(repository RolesRepository) RolesService {
	return &rolesService{
		repo: repository,
	}
}

type rolesService struct {
	repo RolesRepository
}

func (svc *rolesService) RolesService() string {
	return "rolesService"
}

func (svc *rolesService) Count(ctx context.Context) (int, error) {
	return svc.repo.Count(ctx)
}

func (svc *rolesService) Save(ctx context.Context, role *Role) error {
	if role == nil {
		return apperror.New("ROLES.000", ErrSaveRole, ErrNilRole)
	}
	if err := svc.repo.Save(ctx, role); err != nil {
		return apperror.New("ROLES.001", ErrSaveRole, err)
	}
	return nil
}

func (svc *rolesService) Update(ctx context.Context, role *Role) error {
	if role == nil {
		return apperror.New("ROLES.010", ErrUpdateRole, ErrNilRole)
	}
	if err := svc.repo.Update(ctx, role); err != nil {
		return apperror.New("ROLES.011", ErrUpdateRole, err)
	}

	return nil
}

func (svc *rolesService) Create(ctx context.Context, role *Role) error {
	if role == nil {
		return apperror.New("ROLES.020", ErrCreateRole, ErrNilRole)
	}
	if err := svc.repo.Create(ctx, role); err != nil {
		return apperror.New("ROLES.021", ErrCreateRole, err)
	}
	return nil
}

func (svc *rolesService) GetByID(ctx context.Context, id uint64) (*Role, error) {
	role, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.New("ROLES.030", ErrFetchRole, err)
	}

	return role, nil
}

func (svc *rolesService) GetAll(ctx context.Context, page int, perPage int, orderBy string, orderDir string) ([]*Role, error) {
	roles, err := svc.repo.GetAll(ctx, page, perPage, orderBy, orderDir)
	if err != nil {
		return nil, apperror.New("ROLES.040", ErrFetchRoles, err)
	}
	return roles, nil
}

func (svc *rolesService) Delete(ctx context.Context, role *Role) error {
	if role == nil {
		return apperror.New("ROLES.050", ErrDeleteRole, ErrNilRole)
	}
	if err := svc.repo.Delete(ctx, role); err != nil {
		return apperror.New("ROLES.051", ErrDeleteRole, err)
	}
	return nil
}

func (svc *rolesService) GetByUserID(ctx context.Context, userID uint64) ([]*Role, error) {
	roles, err := svc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.New("ROLES.060", ErrFetchRoles, err)
	}
	return roles, nil
}

func (svc *rolesService) Assign(ctx context.Context, userID uint64, roleID uint64) error {
	if err := svc.repo.Assign(ctx, userID, roleID); err != nil {
		return apperror.New("ROLES.070", ErrAssignRole, err)
	}
	return nil
}

func (svc *rolesService) Unassign(ctx context.Context, userID uint64, roleID uint64) error {
	if err := svc.repo.Unassign(ctx, userID, roleID); err != nil {
		return apperror.New("ROLES.080", ErrUnassignRole, err)
	}
	return nil
}

func (svc *rolesService) GetUserRole(ctx context.Context, userID uint64) (UserRole, error) {
	roles, err := svc.GetByUserID(ctx, userID)
	if err != nil {
		return UserRoleUnknown, err
	}

	// figure out user's role
	minRole := uint64(math.MaxUint64)
	for _, role := range roles {
		if role.ID < minRole {
			minRole = role.ID
		}
	}

	return UserRole(minRole), nil
}
