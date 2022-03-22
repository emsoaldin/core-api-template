package auth

import (
	"context"
	"database/sql"
	"time"

	"api/providers/db"
)

// AuthRepository interface
type AuthRepository interface {

	// AuthRepository interface implementation signature
	AuthRepository() string

	//GetByID returns Auth for given user
	GetByID(ctx context.Context, provider string, userID uint64) (*AuthProvider, error)

	// Create creates auth object
	Create(ctx context.Context, auth *AuthProvider) error

	// Update auth object
	Update(ctx context.Context, auth *AuthProvider) error

	// Delete auth object
	Delete(ctx context.Context, auth *AuthProvider) error

	// DeleteByID removes auth for given id
	DeleteByID(ctx context.Context, provider string, userID uint64) error

	// DeleteByUserID removes all auth providers for given user
	DeleteByUserID(ctx context.Context, userID uint64) error

	// GetByUserID returns all Auth strategies for given user
	GetByUserID(ctx context.Context, userID uint64) ([]*AuthProvider, error)
}

// NewAuthRepository creates AuthRepository interface implementation
func NewAuthRepository(store db.Store) AuthRepository {
	return &authRepository{
		store: store,
	}
}

type authRepository struct {
	store db.Store
}

func (r *authRepository) AuthRepository() string {
	return "authRepository"
}

func (r *authRepository) getTx(ctx context.Context) (*sql.Tx, bool, error) {
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

func (authRepository) closeTx(tx *sql.Tx, shouldCommit bool, hasError bool) {
	if shouldCommit {
		if hasError {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}

//GetByID returns Auth for given user
func (r *authRepository) GetByID(ctx context.Context, provider string, userID uint64) (*AuthProvider, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT provider, user_id, uid, created_at, updated_at FROM auth_providers WHERE provider = ? AND user_id = ? "

	// create empty model object
	model := new(AuthProvider)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, provider, userID).Scan(&model.Provider, &model.UserID, &model.Hash, &model.CreatedAt, &model.UpdatedAt)

	if err != nil && err == sql.ErrNoRows {
		err = nil // set err to nil
		return nil, nil
	}

	return model, err
}

// Create creates auth object
func (r *authRepository) Create(ctx context.Context, auth *AuthProvider) error {

	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "INSERT INTO auth_providers (provider, user_id, uid) VALUES(?,?,?)"

	_, err = tx.Exec(query, auth.Provider, auth.UserID, auth.Hash)
	if err != nil {
		return err
	}

	auth.CreatedAt = time.Now()
	auth.UpdatedAt = auth.CreatedAt

	return err
}

// Update auth object
func (r *authRepository) Update(ctx context.Context, auth *AuthProvider) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	q := "UPDATE auth_providers SET uid = ? WHERE provider = ? AND user_id = ?"
	_, err = tx.Exec(q, auth.Hash, auth.Provider, auth.UserID)
	auth.UpdatedAt = time.Now()
	return err
}

// Delete auth object
func (r *authRepository) Delete(ctx context.Context, auth *AuthProvider) error {
	return r.DeleteByID(ctx, auth.Provider, auth.UserID)
}

// DeleteByID removes auth for given id
func (r *authRepository) DeleteByID(ctx context.Context, provider string, userID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	q := "DELETE FROM auth_providers WHERE provider = ? AND user_id = ? "
	_, err = tx.Exec(q, provider, userID)
	return err
}

// DeleteByUserID AuthProvider record from database
func (r *authRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	q := "DELETE FROM auth_providers WHERE user_id = ? "
	_, err = tx.Exec(q, userID)
	return err
}

// GetByUserID returns all Auth strategies for given user
func (r *authRepository) GetByUserID(ctx context.Context, userID uint64) ([]*AuthProvider, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT provider, user_id, uid, created_at, updated_at FROM auth_providers WHERE user_id = ?"

	// create empty model object
	authProviders := make([]*AuthProvider, 0)

	// execute query statement
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// loop over results
	for rows.Next() {
		model := new(AuthProvider)
		// scan row to model
		if err = rows.Scan(&model.Provider, &model.UserID, &model.Hash, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		authProviders = append(authProviders, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authProviders, nil
}
