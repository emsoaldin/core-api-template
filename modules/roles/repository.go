package roles

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api/providers/db"
)

// RolesRepository interface
type RolesRepository interface {
	// RolesRepository returns service interface
	RolesRepository() string

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

	// Delete role from database
	Delete(ctx context.Context, role *Role) error

	// Delete role from database
	DeleteByID(ctx context.Context, id uint64) error

	// GetByUserID returns all roles assigned to user
	GetByUserID(ctx context.Context, userID uint64) ([]*Role, error)

	// Assign assignes role to user
	Assign(ctx context.Context, userID uint64, roleID uint64) error

	// Unassign removes role from user
	Unassign(ctx context.Context, userID uint64, roleID uint64) error
}

// NewRolesRepository creates RolesRepository interface implementation
func NewRolesRepository(store db.Store) RolesRepository {
	return &rolesRepository{
		store: store,
	}
}

type rolesRepository struct {
	store db.Store
}

func (r *rolesRepository) getTx(ctx context.Context) (*sql.Tx, bool, error) {
	var err error
	// get transaction from context
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		// create new transaction
		tx, err = r.store.Begin()
		if err != nil {
			return nil, false, err
		}
		return tx, true, nil
	}
	return tx, false, nil
}

func (rolesRepository) closeTx(tx *sql.Tx, shouldCommit bool, hasError bool) {
	if shouldCommit {
		if hasError {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}

func (r *rolesRepository) RolesRepository() string {
	return "rolesRepository"
}

func (r *rolesRepository) Count(ctx context.Context) (int, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return 0, err
	}
	query := "SELECT COUNT(id) as count FROM roles"

	var count int
	err = tx.QueryRow(query).Scan(&count)

	r.closeTx(tx, shouldCommit, err != nil)

	return count, err
}

func (r *rolesRepository) Save(ctx context.Context, role *Role) error {
	if role.ID > 0 {
		return r.Update(ctx, role)
	}
	return r.Create(ctx, role)
}

func (r *rolesRepository) Update(ctx context.Context, role *Role) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "UPDATE roles SET name = ?, description = ? WHERE id = ?"
	_, err = tx.Exec(query, role.Name, role.Description, role.ID)
	role.UpdatedAt = time.Now()
	return err
}

func (r *rolesRepository) Create(ctx context.Context, role *Role) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "INSERT INTO roles (name, description) VALUES(?,?)"

	result, err := tx.Exec(query, role.Name, role.Description)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()

	if lastID > 0 {
		role.ID = uint64(lastID)
		role.CreatedAt = time.Now()
		role.UpdatedAt = role.CreatedAt
	}

	return err
}

func (r *rolesRepository) GetByID(ctx context.Context, id uint64) (*Role, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT id, name, description, created_at, updated_at FROM roles WHERE id = ? "

	// create empty model object
	model := new(Role)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, id).Scan(&model.ID, &model.Name, &model.Description, &model.CreatedAt, &model.UpdatedAt)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return model, err
}

func (r *rolesRepository) GetAll(ctx context.Context, page int, perPage int, orderBy string, orderDir string) ([]*Role, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	offset := page * perPage
	if orderBy == "" {
		orderBy = "id"
	}
	if orderDir == "" {
		orderDir = "ASC"
	}
	query := fmt.Sprintf(`
		SELECT 
			id, name, description, created_at, updated_at 
		FROM roles 
		LIMIT %d OFFSET %d 
		ORDER BY %s %s`,
		perPage, offset, orderBy, orderDir)

	// execute query statement
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*Role, 0, perPage)
	// loop over results
	for rows.Next() {
		model := new(Role)
		// scan row to model
		if err = rows.Scan(&model.ID, &model.Name, &model.Description, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, err
}

func (r *rolesRepository) Delete(ctx context.Context, role *Role) error {
	return r.DeleteByID(ctx, role.ID)
}

func (r *rolesRepository) DeleteByID(ctx context.Context, id uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM roles WHERE id = ? "
	_, err = tx.Exec(query, id)
	return err
}

func (r *rolesRepository) GetByUserID(ctx context.Context, userID uint64) ([]*Role, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := `
		SELECT 
			r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles AS r 
		INNER JOIN users_roles AS ur ON ur.role_id = r.id 
		WHERE ur.user_id = ?`

	// create empty model object
	roles := make([]*Role, 0)

	// execute query statement
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// loop over results
	for rows.Next() {
		model := new(Role)
		// scan row to model
		if err = rows.Scan(&model.ID, &model.Name, &model.Description, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, err
}

func (r *rolesRepository) Assign(ctx context.Context, userID uint64, roleID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "INSERT INTO users_roles (user_id, role_id) VALUES(?,?)"
	_, err = tx.Exec(query, userID, roleID)

	return err
}

func (r *rolesRepository) Unassign(ctx context.Context, userID uint64, roleID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	// close tx
	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM users_roles WHERE user_id = ? AND role_id = ?"
	_, err = tx.Exec(query, userID, roleID)

	return err
}
