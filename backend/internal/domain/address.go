package domain

import "time"

type Address struct {
	ID             int64     `db:"id"               json:"id"`
	FamilyMemberID *int64    `db:"family_member_id" json:"family_member_id,omitempty"`
	Label          string    `db:"label"            json:"label"`
	Street         *string   `db:"street"           json:"street,omitempty"`
	City           string    `db:"city"             json:"city"`
	Province       *string   `db:"province"         json:"province,omitempty"`
	PostalCode     *string   `db:"postal_code"      json:"postal_code,omitempty"`
	Country        string    `db:"country"          json:"country"`
	Latitude       *float64  `db:"latitude"         json:"latitude,omitempty"`
	Longitude      *float64  `db:"longitude"        json:"longitude,omitempty"`
	IsCurrent      bool      `db:"is_current"       json:"is_current"`

	// audit
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version"    json:"version"`
	IsDeleted bool      `db:"is_deleted" json:"-"`

	// joined
	MemberName *string `db:"member_name" json:"member_name,omitempty"`
}

type CreateAddressRequest struct {
	FamilyMemberID *int64   `json:"family_member_id"`
	Label          string   `json:"label"    binding:"required"`
	Street         *string  `json:"street"`
	City           string   `json:"city"     binding:"required"`
	Province       *string  `json:"province"`
	PostalCode     *string  `json:"postal_code"`
	Country        string   `json:"country"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	IsCurrent      bool     `json:"is_current"`
}

type AddressRepository interface {
	Create(a *Address) error
	FindAll() ([]*Address, error)
	FindByMemberID(memberID int64) ([]*Address, error)
	FindByID(id int64) (*Address, error)
	Update(a *Address) error
	Delete(id int64) error
}
