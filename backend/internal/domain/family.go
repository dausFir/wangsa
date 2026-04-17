package domain

import "time"

type FamilyMember struct {
	ID         int64     `db:"id"          json:"id"`
	UserID     *int64    `db:"user_id"     json:"user_id,omitempty"`
	FullName   string    `db:"full_name"   json:"full_name"`
	Nickname   *string   `db:"nickname"    json:"nickname,omitempty"`
	Gender     string    `db:"gender"      json:"gender"`
	BirthDate  *string   `db:"birth_date"  json:"birth_date,omitempty"`
	BirthPlace *string   `db:"birth_place" json:"birth_place,omitempty"`
	DeathDate  *string   `db:"death_date"  json:"death_date,omitempty"`
	PhotoURL   *string   `db:"photo_url"   json:"photo_url,omitempty"`
	ParentID   *int64    `db:"parent_id"   json:"parent_id,omitempty"`
	Notes      *string   `db:"notes"       json:"notes,omitempty"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`

	// Virtual fields — populated during tree-build, not stored in DB
	Children  []*FamilyMember `db:"-" json:"children,omitempty"`
	Spouses   []*FamilyMember `db:"-" json:"spouses,omitempty"`
	Marriages []*Marriage     `db:"-" json:"marriages,omitempty"`
}

type Marriage struct {
	ID           int64   `db:"id"            json:"id"`
	HusbandID    int64   `db:"husband_id"    json:"husband_id"`
	WifeID       int64   `db:"wife_id"       json:"wife_id"`
	MarriageDate *string `db:"marriage_date" json:"marriage_date,omitempty"`
	DivorceDate  *string `db:"divorce_date"  json:"divorce_date,omitempty"`
	Notes        *string `db:"notes"         json:"notes,omitempty"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`
}

type CreateFamilyMemberRequest struct {
	FullName   string  `json:"full_name"   binding:"required,min=2"`
	Nickname   *string `json:"nickname"`
	Gender     string  `json:"gender"      binding:"required,oneof=male female"`
	BirthDate  *string `json:"birth_date"`
	BirthPlace *string `json:"birth_place"`
	DeathDate  *string `json:"death_date"`
	PhotoURL   *string `json:"photo_url"`
	ParentID   *int64  `json:"parent_id"`
	Notes      *string `json:"notes"`
}

type CreateMarriageRequest struct {
	HusbandID    int64   `json:"husband_id"    binding:"required"`
	WifeID       int64   `json:"wife_id"       binding:"required"`
	MarriageDate *string `json:"marriage_date"`
	DivorceDate  *string `json:"divorce_date"`
	Notes        *string `json:"notes"`
}

type FamilyRepository interface {
	CreateMember(m *FamilyMember) error
	UpdateMember(m *FamilyMember) error
	DeleteMember(id int64) error
	FindMemberByID(id int64) (*FamilyMember, error)
	FindAllMembers() ([]*FamilyMember, error)

	CreateMarriage(m *Marriage) error
	FindAllMarriages() ([]*Marriage, error)
	FindMarriagesByMemberID(memberID int64) ([]*Marriage, error)
	DeleteMarriage(id int64) error
}
