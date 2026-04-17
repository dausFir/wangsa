package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.QueryRow(
		`INSERT INTO users (name, email, password, role, avatar_url, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		user.Name, user.Email, user.Password, user.Role,
		user.AvatarURL, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u,
		`SELECT * FROM users WHERE email=$1 AND is_deleted=FALSE LIMIT 1`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u,
		`SELECT * FROM users WHERE id=$1 AND is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *userRepository) CountAll() (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM users WHERE is_deleted=FALSE`)
	return count, err
}

// ClaimRoleForFirstUser atomically determines the role for a new registrant.
// Uses a transaction with SERIALIZABLE isolation to prevent the first-user race condition.
func (r *userRepository) ClaimRoleForFirstUser() (string, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// SERIALIZABLE prevents concurrent transactions from both seeing count=0
	if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return "", fmt.Errorf("set isolation: %w", err)
	}

	var count int
	if err := tx.Get(&count, `SELECT COUNT(*) FROM users WHERE is_deleted=FALSE`); err != nil {
		return "", fmt.Errorf("count users: %w", err)
	}

	role := "member"
	if count == 0 {
		role = "super_admin"
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return role, nil
}

func (r *userRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET name=$1, avatar_url=$2, updated_at=$3 WHERE id=$4`,
		user.Name, user.AvatarURL, user.UpdatedAt, user.ID,
	)
	return err
}
