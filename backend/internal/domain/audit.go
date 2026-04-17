package domain

import "time"

// AuditFields adalah embed struct yang dipakai oleh semua entitas utama.
// Tambahkan ke setiap struct dengan cara embed atau deklarasi eksplisit.
type AuditFields struct {
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"` // tidak di-expose ke API
}

// AuditLog adalah record immutable dari setiap perubahan data.
type AuditLog struct {
	ID        int64     `db:"id"         json:"id"`
	TableName string    `db:"table_name" json:"table_name"`
	RecordID  int64     `db:"record_id"  json:"record_id"`
	Action    string    `db:"action"     json:"action"` // create|update|delete|restore
	ChangedBy *int64    `db:"changed_by" json:"changed_by,omitempty"`
	ChangedAt time.Time `db:"changed_at" json:"changed_at"`
	Diff      *string   `db:"diff"       json:"diff,omitempty"`
}

// AuditLogRepository adalah interface untuk menulis audit log.
type AuditLogRepository interface {
	Write(log *AuditLog) error
	FindByRecord(tableName string, recordID int64) ([]*AuditLog, error)
}
