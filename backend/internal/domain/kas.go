package domain

import "time"

type KasCategory struct {
	ID   int64  `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
	Icon string `db:"icon" json:"icon"`
}

type KasTransaction struct {
	ID          int64     `db:"id"          json:"id"`
	CategoryID  *int64    `db:"category_id" json:"category_id,omitempty"`
	Type        string    `db:"type"        json:"type"`
	Amount      float64   `db:"amount"      json:"amount"`
	Description *string   `db:"description" json:"description,omitempty"`
	Date        string    `db:"date"        json:"date"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`

	// joined
	CategoryName *string `db:"category_name" json:"category_name,omitempty"`
	CreatorName  *string `db:"creator_name"  json:"creator_name,omitempty"`
}

type CreateKasTransactionRequest struct {
	CategoryID  *int64  `json:"category_id"`
	Type        string  `json:"type"   binding:"required,oneof=income expense"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description *string `json:"description"`
	Date        string  `json:"date"   binding:"required"`
}

type KasSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type KasRepository interface {
	FindAllCategories() ([]*KasCategory, error)
	CreateTransaction(t *KasTransaction) error
	FindAllTransactions(limit, offset int) ([]*KasTransaction, error)
	FindTransactionByID(id int64) (*KasTransaction, error)
	UpdateTransaction(t *KasTransaction) error
	DeleteTransaction(id int64) error
	GetSummary() (*KasSummary, error)
}
