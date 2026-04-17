package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type addressRepository struct{ db *sqlx.DB }

func NewAddressRepository(db *sqlx.DB) domain.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(a *domain.Address) error {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Version   = 1
	return r.db.QueryRow(
		`INSERT INTO addresses
		 (family_member_id, label, street, city, province, postal_code, country,
		  latitude, longitude, is_current,
		  created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,FALSE)
		 RETURNING id`,
		a.FamilyMemberID, a.Label, a.Street, a.City, a.Province, a.PostalCode, a.Country,
		a.Latitude, a.Longitude, a.IsCurrent,
		a.CreatedBy, a.CreatedAt, a.CreatedBy, a.UpdatedAt, a.Version,
	).Scan(&a.ID)
}

func (r *addressRepository) FindAll() ([]*domain.Address, error) {
	var addrs []*domain.Address
	err := r.db.Select(&addrs, `
		SELECT a.*, fm.full_name AS member_name
		FROM addresses a
		LEFT JOIN family_members fm ON a.family_member_id = fm.id
		WHERE a.is_deleted = FALSE
		ORDER BY a.city ASC, a.id ASC`)
	if err != nil {
		return nil, fmt.Errorf("find all addresses: %w", err)
	}
	return addrs, nil
}

func (r *addressRepository) FindByMemberID(memberID int64) ([]*domain.Address, error) {
	var addrs []*domain.Address
	err := r.db.Select(&addrs,
		`SELECT * FROM addresses WHERE family_member_id=$1 AND is_deleted=FALSE ORDER BY is_current DESC`,
		memberID)
	return addrs, err
}

func (r *addressRepository) FindByID(id int64) (*domain.Address, error) {
	var a domain.Address
	err := r.db.Get(&a,
		`SELECT * FROM addresses WHERE id=$1 AND is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find address by id: %w", err)
	}
	return &a, nil
}

func (r *addressRepository) Update(a *domain.Address) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE addresses
		 SET label=$1, street=$2, city=$3, province=$4, postal_code=$5,
		     country=$6, latitude=$7, longitude=$8, is_current=$9,
		     updated_by=$10, updated_at=$11, version=version+1
		 WHERE id=$12 AND is_deleted=FALSE`,
		a.Label, a.Street, a.City, a.Province, a.PostalCode,
		a.Country, a.Latitude, a.Longitude, a.IsCurrent,
		a.UpdatedBy, a.UpdatedAt, a.ID,
	)
	return err
}

func (r *addressRepository) Delete(id int64) error {
	_, err := r.db.Exec(
		`UPDATE addresses
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}
