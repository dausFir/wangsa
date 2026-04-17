package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type kasRepository struct{ db *sqlx.DB }

func NewKasRepository(db *sqlx.DB) domain.KasRepository {
	return &kasRepository{db: db}
}

func (r *kasRepository) FindAllCategories() ([]*domain.KasCategory, error) {
	var cats []*domain.KasCategory
	err := r.db.Select(&cats, `SELECT * FROM kas_categories ORDER BY name ASC`)
	return cats, err
}

func (r *kasRepository) CreateTransaction(t *domain.KasTransaction) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.Version   = 1
	return r.db.QueryRow(
		`INSERT INTO kas_transactions
		 (category_id, type, amount, description, date,
		  created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,FALSE)
		 RETURNING id`,
		t.CategoryID, t.Type, t.Amount, t.Description, t.Date,
		t.CreatedBy, t.CreatedAt, t.CreatedBy, t.UpdatedAt, t.Version,
	).Scan(&t.ID)
}

func (r *kasRepository) FindAllTransactions(limit, offset int) ([]*domain.KasTransaction, error) {
	var txs []*domain.KasTransaction
	err := r.db.Select(&txs, `
		SELECT t.*, c.name AS category_name, u.name AS creator_name
		FROM kas_transactions t
		LEFT JOIN kas_categories c ON t.category_id = c.id
		LEFT JOIN users u ON t.created_by = u.id
		WHERE t.is_deleted = FALSE
		ORDER BY t.date DESC, t.id DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find transactions: %w", err)
	}
	return txs, nil
}

func (r *kasRepository) FindTransactionByID(id int64) (*domain.KasTransaction, error) {
	var t domain.KasTransaction
	err := r.db.Get(&t,
		`SELECT * FROM kas_transactions WHERE id=$1 AND is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find transaction by id: %w", err)
	}
	return &t, nil
}

func (r *kasRepository) UpdateTransaction(t *domain.KasTransaction) error {
	t.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE kas_transactions
		 SET category_id=$1, type=$2, amount=$3, description=$4, date=$5,
		     updated_by=$6, updated_at=$7, version=version+1
		 WHERE id=$8 AND is_deleted=FALSE`,
		t.CategoryID, t.Type, t.Amount, t.Description, t.Date,
		t.UpdatedBy, t.UpdatedAt, t.ID,
	)
	return err
}

func (r *kasRepository) DeleteTransaction(id int64) error {
	_, err := r.db.Exec(
		`UPDATE kas_transactions
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}

func (r *kasRepository) GetSummary() (*domain.KasSummary, error) {
	row := r.db.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0)
		FROM kas_transactions
		WHERE is_deleted = FALSE`)
	s := &domain.KasSummary{}
	if err := row.Scan(&s.Income, &s.Expense); err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}
	s.Balance = s.Income - s.Expense
	return s, nil
}
