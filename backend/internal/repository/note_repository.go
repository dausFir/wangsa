package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wangsa/backend/internal/domain"
)

type noteRepository struct{ db *sqlx.DB }

func NewNoteRepository(db *sqlx.DB) domain.NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(n *domain.Note) error {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	n.Version = 1
	return r.db.QueryRow(
		`INSERT INTO notes
		 (title, content, category, is_pinned, created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,FALSE)
		 RETURNING id`,
		n.Title, n.Content, n.Category, n.IsPinned,
		n.CreatedBy, n.CreatedAt, n.UpdatedBy, n.UpdatedAt, n.Version,
	).Scan(&n.ID)
}

func (r *noteRepository) FindAll(category string) ([]*domain.Note, error) {
	var notes []*domain.Note
	base := `
		SELECT n.*, 
		       uc.name AS created_by_name,
		       uu.name AS updated_by_name
		FROM notes n
		LEFT JOIN users uc ON n.created_by = uc.id
		LEFT JOIN users uu ON n.updated_by = uu.id
		WHERE n.is_deleted = FALSE`

	var err error
	if category != "" {
		err = r.db.Select(&notes,
			base+` AND n.category = $1 ORDER BY n.is_pinned DESC, n.updated_at DESC`,
			category)
	} else {
		err = r.db.Select(&notes,
			base+` ORDER BY n.is_pinned DESC, n.updated_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("find all notes: %w", err)
	}
	return notes, nil
}

func (r *noteRepository) FindByID(id int64) (*domain.Note, error) {
	var n domain.Note
	err := r.db.Get(&n,
		`SELECT n.*, 
		        uc.name AS created_by_name,
		        uu.name AS updated_by_name
		 FROM notes n
		 LEFT JOIN users uc ON n.created_by = uc.id
		 LEFT JOIN users uu ON n.updated_by = uu.id
		 WHERE n.id=$1 AND n.is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find note by id: %w", err)
	}
	return &n, nil
}

func (r *noteRepository) Update(n *domain.Note) error {
	n.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE notes
		 SET title=$1, content=$2, category=$3, is_pinned=$4,
		     updated_by=$5, updated_at=$6, version=version+1
		 WHERE id=$7 AND is_deleted=FALSE`,
		n.Title, n.Content, n.Category, n.IsPinned,
		n.UpdatedBy, n.UpdatedAt, n.ID,
	)
	return err
}

func (r *noteRepository) Delete(id int64) error {
	_, err := r.db.Exec(
		`UPDATE notes
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}

func (r *noteRepository) FindCategories() ([]string, error) {
	var categories []string
	err := r.db.Select(&categories,
		`SELECT DISTINCT category 
		 FROM notes 
		 WHERE category IS NOT NULL AND category != '' AND is_deleted = FALSE
		 ORDER BY category`)
	if err != nil {
		return nil, fmt.Errorf("find categories: %w", err)
	}
	return categories, nil
}
