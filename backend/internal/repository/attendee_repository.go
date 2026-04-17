package repository

import (
	"fmt"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type attendeeRepository struct{ db *sqlx.DB }

func NewAttendeeRepository(db *sqlx.DB) domain.AttendeeRepository {
	return &attendeeRepository{db: db}
}

func (r *attendeeRepository) ListAttendees(eventID int64) ([]*domain.EventAttendee, error) {
	var attendees []*domain.EventAttendee
	err := r.db.Select(&attendees, `
		SELECT ea.*, fm.full_name AS member_name
		FROM event_attendees ea
		JOIN family_members fm ON ea.family_member_id = fm.id AND fm.is_deleted = FALSE
		WHERE ea.event_id = $1
		ORDER BY fm.full_name ASC`, eventID)
	if err != nil {
		return nil, fmt.Errorf("list attendees: %w", err)
	}
	return attendees, nil
}

// UpsertAttendee inserts a new RSVP or updates an existing one.
func (r *attendeeRepository) UpsertAttendee(eventID, memberID int64, rsvp string) error {
	_, err := r.db.Exec(`
		INSERT INTO event_attendees (event_id, family_member_id, rsvp)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, family_member_id)
		DO UPDATE SET rsvp = EXCLUDED.rsvp`,
		eventID, memberID, rsvp)
	if err != nil {
		return fmt.Errorf("upsert attendee: %w", err)
	}
	return nil
}

func (r *attendeeRepository) RemoveAttendee(eventID, memberID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM event_attendees WHERE event_id=$1 AND family_member_id=$2`,
		eventID, memberID)
	return err
}
