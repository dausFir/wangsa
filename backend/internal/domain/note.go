package domain

import "time"

type Note struct {
	ID       int64   `db:"id"        json:"id"`
	Title    string  `db:"title"     json:"title"`
	Content  string  `db:"content"   json:"content"`
	Category *string `db:"category"  json:"category,omitempty"`
	IsPinned bool    `db:"is_pinned" json:"is_pinned"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`

	// joined fields
	CreatedByName *string `db:"created_by_name" json:"created_by_name,omitempty"`
	UpdatedByName *string `db:"updated_by_name" json:"updated_by_name,omitempty"`
}

type CreateNoteRequest struct {
	Title    string  `json:"title"    binding:"required,min=2,max=200"`
	Content  string  `json:"content"  binding:"required"`
	Category *string `json:"category"`
	IsPinned bool    `json:"is_pinned"`
}

type UpdateNoteRequest struct {
	Title    string  `json:"title"    binding:"required,min=2,max=200"`
	Content  string  `json:"content"  binding:"required"`
	Category *string `json:"category"`
	IsPinned bool    `json:"is_pinned"`
}

type NoteRepository interface {
	Create(n *Note) error
	FindAll(category string) ([]*Note, error)
	FindByID(id int64) (*Note, error)
	Update(n *Note) error
	Delete(id int64) error
	FindCategories() ([]string, error)
}
