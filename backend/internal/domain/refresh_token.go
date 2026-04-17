package domain

import "time"

// RefreshToken is stored in DB so we can revoke on logout or suspicious activity.
type RefreshToken struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	TokenHash string     `db:"token_hash"` // SHA-256 hex of the raw token
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"` // nil = still valid
}

type RefreshTokenRepository interface {
	// Store saves a new refresh token (hashed). Any previously active token
	// for the same user is revoked atomically (single-session policy).
	Store(rt *RefreshToken) error
	// FindByHash looks up an active (non-revoked, non-expired) token by its hash.
	FindByHash(tokenHash string) (*RefreshToken, error)
	// Revoke marks a single token as revoked (called on logout or rotation).
	Revoke(id int64) error
	// RevokeAllForUser revokes every active token for a user (security reset).
	RevokeAllForUser(userID int64) error
}
