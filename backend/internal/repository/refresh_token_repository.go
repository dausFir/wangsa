package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct{ db *sqlx.DB }

func NewRefreshTokenRepository(db *sqlx.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Store saves a new refresh token and atomically revokes any previous active
// tokens for the same user — enforces single-session-per-user policy.
func (r *refreshTokenRepository) Store(rt *domain.RefreshToken) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Revoke all existing active tokens for this user
	_, err = tx.Exec(
		`UPDATE refresh_tokens
		 SET revoked_at = NOW()
		 WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()`,
		rt.UserID,
	)
	if err != nil {
		return fmt.Errorf("revoke old tokens: %w", err)
	}

	// Insert the new token
	_, err = tx.Exec(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)`,
		rt.UserID, rt.TokenHash, rt.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}

	return tx.Commit()
}

// FindByHash returns an active (non-revoked, non-expired) token by hash.
func (r *refreshTokenRepository) FindByHash(tokenHash string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.Get(&rt,
		`SELECT * FROM refresh_tokens
		 WHERE token_hash = $1
		   AND revoked_at IS NULL
		   AND expires_at > NOW()
		 LIMIT 1`,
		tokenHash,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return &rt, nil
}

func (r *refreshTokenRepository) Revoke(id int64) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`, id,
	)
	return err
}

func (r *refreshTokenRepository) RevokeAllForUser(userID int64) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens
		 SET revoked_at = NOW()
		 WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	)
	return err
}
