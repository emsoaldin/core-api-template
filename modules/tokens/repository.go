package tokens

import (
	"context"
	"database/sql"
	"time"

	"api/providers/db"
)

// TokensRepository interface
type TokensRepository interface {
	//TokensRepository interface implementation signature
	TokensRepository() string

	// GetByID returns Token from database with provided id
	GetByID(ctx context.Context, id uint64) (*Token, error)

	// Save token
	Save(ctx context.Context, token *Token) error

	// Create token
	Create(ctx context.Context, token *Token) error

	// Update token
	Update(ctx context.Context, token *Token) error

	// Delete token
	Delete(ctx context.Context, token *Token) error

	// DeleteByID removes token with provided id
	DeleteByID(ctx context.Context, id uint64) error

	// GetByUserID returns all tokens for provided userID
	GetByUserID(ctx context.Context, userID uint64) ([]*Token, error)

	// GetByUserAndTokenID returns tokens for provided userID and tokenTypeID
	GetByUserAndTokenID(ctx context.Context, userID uint64, tokenTypeID uint64) ([]*Token, error)

	// GetByToken returns token object for provided token string
	GetByToken(ctx context.Context, token string) (*Token, error)

	// DeleteByUserID removes all tokens for given userID
	DeleteByUserID(ctx context.Context, userID uint64) error

	// DeleteByUserAndTokenTypeID removes all tokens for given userID and tokenTypeID
	DeleteByUserAndTokenTypeID(ctx context.Context, userID uint64, tokenTypeID uint64) error

	// DeleteExpiredTokens removes all tokens that are expired
	DeleteExpiredTokens(ctx context.Context) error
}

// NewTokensRepository creates TokensRepository interface implementation
func NewTokensRepository(store db.Store) TokensRepository {
	return &tokensRepository{
		store: store,
	}
}

type tokensRepository struct {
	store db.Store
}

func (r *tokensRepository) getTx(ctx context.Context) (*sql.Tx, bool, error) {
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

func (tokensRepository) closeTx(tx *sql.Tx, shouldCommit bool, hasError bool) {
	if shouldCommit {
		if hasError {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}

func (r *tokensRepository) TokensRepository() string {
	return "tokensRepository"
}

// GetByID returns Token from database with provided id
func (r *tokensRepository) GetByID(ctx context.Context, id uint64) (*Token, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT id, user_id, token, meta, token_type_id, expires_at, created_at, updated_at FROM tokens WHERE id = ? "

	// create empty model object
	model := new(Token)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, id).Scan(&model.ID, &model.UserID, &model.Token, &model.Meta, &model.TokenTypeID, &model.ExpiresAt, &model.CreatedAt, &model.UpdatedAt)

	if err != nil && err == sql.ErrNoRows {
		err = nil
		return nil, nil
	}

	return model, err
}

// Save token
func (r *tokensRepository) Save(ctx context.Context, token *Token) error {
	if token.ID > 0 {
		return r.Update(ctx, token)
	}
	return r.Create(ctx, token)
}

// Create token
func (r *tokensRepository) Create(ctx context.Context, token *Token) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "INSERT INTO tokens (user_id, token, meta, token_type_id, expires_at) VALUES(?,?,?,?,?)"

	result, err := tx.Exec(query, token.UserID, token.Token, token.Meta, token.TokenTypeID, token.ExpiresAt)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()

	if lastID > 0 {
		token.ID = uint64(lastID)
		token.CreatedAt = time.Now()
		token.UpdatedAt = token.CreatedAt
	}

	return err
}

// Update token
func (r *tokensRepository) Update(ctx context.Context, token *Token) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "UPDATE tokens SET token = ?, meta = ?, token_type_id = ?, expires_at = ? WHERE id = ?"
	_, err = tx.Exec(query, token.Token, token.Meta, token.TokenTypeID, token.ExpiresAt, token.ID)
	token.UpdatedAt = time.Now()
	return err
}

// Delete token
func (r *tokensRepository) Delete(ctx context.Context, token *Token) error {
	return r.DeleteByID(ctx, token.ID)
}

// DeleteByID removes token with provided id
func (r *tokensRepository) DeleteByID(ctx context.Context, id uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM tokens WHERE id = ? "
	_, err = tx.Exec(query, id)
	return err
}

// GetByUserID returns all tokens for provided userID
func (r *tokensRepository) GetByUserID(ctx context.Context, userID uint64) ([]*Token, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT id, user_id, token, meta, token_type_id, expires_at, created_at, updated_at FROM tokens WHERE user_id = ?"

	// create empty model object
	tokens := make([]*Token, 0)

	// execute query statement
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// loop over results
	for rows.Next() {
		model := new(Token)
		// scan row to model
		if err = rows.Scan(&model.ID, &model.UserID, &model.Token, &model.Meta, &model.TokenTypeID, &model.ExpiresAt, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}

		tokens = append(tokens, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetByUserAndTokenID returns tokens for provided userID and tokenTypeID
func (r *tokensRepository) GetByUserAndTokenID(ctx context.Context, userID uint64, tokenTypeID uint64) ([]*Token, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT id, user_id, token, meta, token_type_id, expires_at, created_at, updated_at FROM tokens WHERE user_id = ? AND token_type_id = ?"

	// create empty model object
	tokens := make([]*Token, 0)

	// execute query statement
	rows, err := tx.Query(query, userID, tokenTypeID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// loop over results
	for rows.Next() {
		model := new(Token)
		// scan row to model
		if err = rows.Scan(&model.ID, &model.UserID, &model.Token, &model.Meta, &model.TokenTypeID, &model.ExpiresAt, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}

		tokens = append(tokens, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetByToken returns token object for provided token string
func (r *tokensRepository) GetByToken(ctx context.Context, token string) (*Token, error) {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "SELECT id, user_id, token, meta, token_type_id, expires_at, created_at, updated_at FROM tokens WHERE token = ?"

	// create empty model object
	model := new(Token)

	// execute query statement and scan row to model
	err = tx.QueryRow(query, token).Scan(&model.ID, &model.UserID, &model.Token, &model.Meta, &model.TokenTypeID, &model.ExpiresAt, &model.CreatedAt, &model.UpdatedAt)

	if err != nil && err == sql.ErrNoRows {
		err = nil
		return nil, nil
	}

	return model, err
}

// DeleteByUserID removes all tokens for given userID
func (r *tokensRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM tokens WHERE user_id = ? "
	_, err = tx.Exec(query, userID)

	return err
}

// DeleteByUserAndTokenTypeID removes all tokens for given userID and tokenTypeID
func (r *tokensRepository) DeleteByUserAndTokenTypeID(ctx context.Context, userID uint64, tokenTypeID uint64) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM tokens WHERE user_id = ? AND token_type_id = ?"
	_, err = tx.Exec(query, userID, tokenTypeID)

	return err
}

// DeleteExpiredTokens removes all tokens that are expired
func (r *tokensRepository) DeleteExpiredTokens(ctx context.Context) error {
	tx, shouldCommit, err := r.getTx(ctx)
	if err != nil {
		return err
	}

	defer r.closeTx(tx, shouldCommit, err != nil)

	query := "DELETE FROM tokens WHERE expires_at < NOW()"
	_, err = tx.Exec(query)
	return err
}
