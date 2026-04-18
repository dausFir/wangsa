package domain

import "time"

type Event struct {
	ID          int64      `db:"id"          json:"id"`
	Title       string     `db:"title"       json:"title"`
	Description *string    `db:"description" json:"description,omitempty"`
	Location    *string    `db:"location"    json:"location,omitempty"`
	StartAt     time.Time  `db:"start_at"    json:"start_at"`
	EndAt       *time.Time `db:"end_at"      json:"end_at,omitempty"`
	IsRecurring bool       `db:"is_recurring" json:"is_recurring"`
	RecurRule   *string    `db:"recur_rule"  json:"recur_rule,omitempty"`
	Color       string     `db:"color"       json:"color"`
	Notes       *string    `db:"notes"       json:"notes,omitempty"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`

	// joined
	AttendeeCount int `db:"attendee_count" json:"attendee_count"`
}

type CreateEventRequest struct {
	Title       string  `json:"title"    binding:"required,min=2"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	StartAt     string  `json:"start_at" binding:"required"`
	EndAt       *string `json:"end_at"`
	IsRecurring bool    `json:"is_recurring"`
	RecurRule   *string `json:"recur_rule"`
	Color       string  `json:"color"`
	Notes       *string `json:"notes"`
}

type EventRepository interface {
	Create(e *Event) error
	FindAll(from, to string) ([]*Event, error)
	FindByID(id int64) (*Event, error)
	Update(e *Event) error
	Delete(id int64) error
}

// EventAttendee represents a family member's RSVP for an event.
type EventAttendee struct {
	EventID        int64  `db:"event_id"         json:"event_id"`
	FamilyMemberID int64  `db:"family_member_id" json:"family_member_id"`
	RSVP           string `db:"rsvp"             json:"rsvp"` // pending|yes|no
	// joined fields
	MemberName *string `db:"member_name" json:"member_name,omitempty"`
}

type RSVPRequest struct {
	RSVP string `json:"rsvp" binding:"required,oneof=pending yes no"`
}

// extend EventRepository interface
type AttendeeRepository interface {
	ListAttendees(eventID int64) ([]*EventAttendee, error)
	UpsertAttendee(eventID, memberID int64, rsvp string) error
	RemoveAttendee(eventID, memberID int64) error
}
