package domain

import "time"

type User struct {
	ID        int64     `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"`
	Role      string    `db:"role"       json:"role"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`
}

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User User `json:"user"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
	CountAll() (int, error)
	ClaimRoleForFirstUser() (string, error)
	Update(user *User) error
}
