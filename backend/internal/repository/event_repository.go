package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type eventRepository struct{ db *sqlx.DB }

func NewEventRepository(db *sqlx.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(e *domain.Event) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	e.Version   = 1
	return r.db.QueryRow(
		`INSERT INTO events
		 (title, description, location, start_at, end_at, is_recurring, recur_rule, color,
		  created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,FALSE)
		 RETURNING id`,
		e.Title, e.Description, e.Location, e.StartAt, e.EndAt,
		e.IsRecurring, e.RecurRule, e.Color,
		e.CreatedBy, e.CreatedAt, e.CreatedBy, e.UpdatedAt, e.Version,
	).Scan(&e.ID)
}

func (r *eventRepository) FindAll(from, to string) ([]*domain.Event, error) {
	var events []*domain.Event
	base := `
		SELECT e.*, COUNT(ea.family_member_id) AS attendee_count
		FROM events e
		LEFT JOIN event_attendees ea ON e.id = ea.event_id
		WHERE e.is_deleted = FALSE`

	var err error
	if from != "" && to != "" {
		err = r.db.Select(&events,
			base+` AND e.start_at BETWEEN $1 AND $2 GROUP BY e.id ORDER BY e.start_at ASC`,
			from, to)
	} else {
		err = r.db.Select(&events,
			base+` GROUP BY e.id ORDER BY e.start_at ASC`)
	}
	if err != nil {
		return nil, fmt.Errorf("find all events: %w", err)
	}
	return events, nil
}

func (r *eventRepository) FindByID(id int64) (*domain.Event, error) {
	var e domain.Event
	err := r.db.Get(&e,
		`SELECT * FROM events WHERE id=$1 AND is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find event by id: %w", err)
	}
	return &e, nil
}

func (r *eventRepository) Update(e *domain.Event) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE events
		 SET title=$1, description=$2, location=$3, start_at=$4, end_at=$5,
		     is_recurring=$6, recur_rule=$7, color=$8,
		     updated_by=$9, updated_at=$10, version=version+1
		 WHERE id=$11 AND is_deleted=FALSE`,
		e.Title, e.Description, e.Location, e.StartAt, e.EndAt,
		e.IsRecurring, e.RecurRule, e.Color,
		e.UpdatedBy, e.UpdatedAt, e.ID,
	)
	return err
}

func (r *eventRepository) Delete(id int64) error {
	_, err := r.db.Exec(
		`UPDATE events
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}
