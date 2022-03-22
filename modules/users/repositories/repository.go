package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"api/modules/users/models"
	"api/providers/db"
)

type UsersRepository interface {
	UsersRepository() string

	// Count returns number of records in database
	Count(ctx context.Context, filter string) (int, error)

	// Save saves given user object
	Save(ctx context.Context, user *models.User) error

	// Update updates user object
	Update(ctx context.Context, user *models.User) error

	// Create user
	Create(ctx context.Context, user *models.User) error

	// GetByID returns user object from database for given id
	GetByID(ctx context.Context, id uint64) (*models.User, error)

	// GetByEmail returns user object from database for given email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetAll returns all User objects for given params
	GetAll(ctx context.Context, filter string, page int, perPage int, orderBy string, orderDir string) ([]*models.User, error)

	// Delete user from database
	Delete(ctx context.Context, user *models.User) error

	// DeleteByID user from database
	DeleteByID(ctx context.Context, id uint64) error
}

// NewUsersRepository creates UsersRepository interface implementation
func NewUsersRepository(store db.Store) UsersRepository {
	return &usersRepository{
		store: store,
	}
}

type usersRepository struct {
	store db.Store
}

func (r *usersRepository) getTx(ctx context.Context) (*sql.Tx, bool, error) {
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

func (usersRepository) closeTx(tx *sql.Tx, shouldCommit bool, hasError bool) {
	if shouldCommit {
		if hasError {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}

func (r *usersRepository) UsersRepository() string {
	return "usersRepository"
}

// Count returns number of records in database
func (r *usersRepository) Count(ctx context.Context, filter string) (int, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return 0, err
	}

	query := "SELECT COUNT(id) as count FROM users"
	wc := r.whereClause(filter)
	if wc != "" {
		query += fmt.Sprintf(" WHERE %s", wc)
	}

	var count int
	err = tx.QueryRow(query).Scan(&count)

	r.closeTx(tx, shouldCommit, err != nil)

	return count, err
}

// Save saves given user object
func (r *usersRepository) Save(ctx context.Context, user *models.User) error {
	if user.ID > 0 {
		return r.Update(ctx, user)
	}
	return r.Create(ctx, user)
}

// Update updates user object
func (r *usersRepository) Update(ctx context.Context, user *models.User) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := `
		UPDATE 
			users 
		SET 
			first_name = ?, 
			last_name = ?
		WHERE id = ?`

	_, err = tx.Exec(query,
		user.FirstName,
		user.LastName,
		user.ID)

	user.UpdatedAt = time.Now()
	return err
}

// Create user
func (r *usersRepository) Create(ctx context.Context, user *models.User) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := `
		INSERT INTO users 
		(first_name, last_name, email) 
		VALUES(?, ?, ?, ?)`

	result, err := tx.Exec(query,
		user.FirstName,
		user.LastName,
		user.Email)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()

	if lastID > 0 {
		user.ID = uint64(lastID)
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
	}

	return err
}

// GetByID returns user object from database for given id
func (r *usersRepository) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := `
	SELECT 
		id, 
		first_name, 
		last_name, 
		email, 
		created_at, 
		updated_at
	FROM 
		users 
	WHERE 
		id = ? AND deleted_at IS NULL`

	// create empty model object
	model := new(models.User)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, id).Scan(
		&model.ID,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return model, err
}

// GetByEmail returns user object from database for given email
func (r *usersRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := `
		SELECT 
			id, 
			first_name, 
			last_name, 
			email, 
			created_at, 
			updated_at
		FROM 
			users 
		WHERE 
			email = ? AND deleted_at IS NULL`

	// create empty model object
	model := new(models.User)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, email).Scan(
		&model.ID,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.CreatedAt,
		&model.UpdatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return model, err
}

// GetAll returns all User objects for given params
func (r *usersRepository) GetAll(ctx context.Context, filter string, page int, perPage int, orderBy string, orderDir string) ([]*models.User, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	offset := (page - 1) * perPage
	if orderBy == "" {
		orderBy = "id"
	}
	if orderDir == "" {
		orderDir = "ASC"
	}

	wc := r.whereClause(filter)
	if wc != "" {
		wc = fmt.Sprintf(" WHERE %s", wc)
	}

	query := fmt.Sprintf(`
		SELECT 
			id, 
			first_name, 
			last_name, 
			email, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM 
			users 
		%s
		ORDER BY %s %s
		LIMIT %d OFFSET %d`,
		wc, orderBy, orderDir, perPage, offset)

	// execute query statement
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*models.User, 0, perPage)
	// loop over results
	for rows.Next() {
		model := new(models.User)
		// scan row to model
		if err := rows.Scan(
			&model.ID,
			&model.FirstName,
			&model.LastName,
			&model.Email,
			&model.CreatedAt,
			&model.UpdatedAt,
			&model.DeletedAt); err != nil {
			return nil, err
		}
		roles = append(roles, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, err
}

// Delete user from database
func (r *usersRepository) Delete(ctx context.Context, user *models.User) error {
	return r.DeleteByID(ctx, user.ID)
}

// DeleteByID user from database
func (r *usersRepository) DeleteByID(ctx context.Context, id uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM users WHERE id = ?"
	_, err = tx.Exec(query, id)
	return err
}

func (usersRepository) whereClause(filter string) string {
	wc := "deleted_at IS NULL"

	if len(filter) > 0 {
		filter = strings.ReplaceAll(filter, "'", "''")
		wc += fmt.Sprintf(" AND (name LIKE '%s%%' OR email LIKE '%s%%')", filter, filter)
	}

	return wc
}
