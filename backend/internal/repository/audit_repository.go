package repository

import (
	"fmt"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type auditLogRepository struct{ db *sqlx.DB }

func NewAuditLogRepository(db *sqlx.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Write(log *domain.AuditLog) error {
	_, err := r.db.Exec(
		`INSERT INTO audit_log (table_name, record_id, action, changed_by, diff)
		 VALUES ($1,$2,$3,$4,$5)`,
		log.TableName, log.RecordID, log.Action, log.ChangedBy, log.Diff,
	)
	if err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepository) FindByRecord(tableName string, recordID int64) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	err := r.db.Select(&logs,
		`SELECT * FROM audit_log WHERE table_name=$1 AND record_id=$2 ORDER BY changed_at DESC`,
		tableName, recordID)
	if err != nil {
		return nil, fmt.Errorf("find audit log: %w", err)
	}
	return logs, nil
}
