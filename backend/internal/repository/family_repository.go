package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type familyRepository struct{ db *sqlx.DB }

func NewFamilyRepository(db *sqlx.DB) domain.FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) CreateMember(m *domain.FamilyMember) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.Version   = 1
	return r.db.QueryRow(
		`INSERT INTO family_members
		 (user_id, full_name, nickname, gender, birth_date, birth_place,
		  death_date, photo_url, parent_id, notes,
		  created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,FALSE)
		 RETURNING id`,
		m.UserID, m.FullName, m.Nickname, m.Gender, m.BirthDate, m.BirthPlace,
		m.DeathDate, m.PhotoURL, m.ParentID, m.Notes,
		m.CreatedBy, m.CreatedAt, m.CreatedBy, m.UpdatedAt, m.Version,
	).Scan(&m.ID)
}

func (r *familyRepository) UpdateMember(m *domain.FamilyMember) error {
	m.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE family_members
		 SET full_name=$1, nickname=$2, gender=$3, birth_date=$4, birth_place=$5,
		     death_date=$6, photo_url=$7, parent_id=$8, notes=$9,
		     updated_by=$10, updated_at=$11, version=version+1
		 WHERE id=$12 AND is_deleted=FALSE`,
		m.FullName, m.Nickname, m.Gender, m.BirthDate, m.BirthPlace,
		m.DeathDate, m.PhotoURL, m.ParentID, m.Notes,
		m.UpdatedBy, m.UpdatedAt, m.ID,
	)
	return err
}

func (r *familyRepository) DeleteMember(id int64) error {
	_, err := r.db.Exec(
		`UPDATE family_members
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}

func (r *familyRepository) FindMemberByID(id int64) (*domain.FamilyMember, error) {
	var m domain.FamilyMember
	err := r.db.Get(&m,
		`SELECT * FROM family_members WHERE id=$1 AND is_deleted=FALSE LIMIT 1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find member by id: %w", err)
	}
	return &m, nil
}

func (r *familyRepository) FindAllMembers() ([]*domain.FamilyMember, error) {
	var members []*domain.FamilyMember
	err := r.db.Select(&members,
		`SELECT * FROM family_members WHERE is_deleted=FALSE ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("find all members: %w", err)
	}
	return members, nil
}

func (r *familyRepository) CreateMarriage(m *domain.Marriage) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.Version   = 1
	return r.db.QueryRow(
		`INSERT INTO marriages
		 (husband_id, wife_id, marriage_date, divorce_date, notes,
		  created_by, created_at, updated_by, updated_at, version, is_deleted)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,FALSE)
		 RETURNING id`,
		m.HusbandID, m.WifeID, m.MarriageDate, m.DivorceDate, m.Notes,
		m.CreatedBy, m.CreatedAt, m.CreatedBy, m.UpdatedAt, m.Version,
	).Scan(&m.ID)
}

func (r *familyRepository) FindAllMarriages() ([]*domain.Marriage, error) {
	var marriages []*domain.Marriage
	err := r.db.Select(&marriages,
		`SELECT * FROM marriages WHERE is_deleted=FALSE ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("find all marriages: %w", err)
	}
	return marriages, nil
}

func (r *familyRepository) FindMarriagesByMemberID(memberID int64) ([]*domain.Marriage, error) {
	var marriages []*domain.Marriage
	err := r.db.Select(&marriages,
		`SELECT * FROM marriages WHERE (husband_id=$1 OR wife_id=$2) AND is_deleted=FALSE`,
		memberID, memberID)
	return marriages, err
}

func (r *familyRepository) DeleteMarriage(id int64) error {
	_, err := r.db.Exec(
		`UPDATE marriages
		 SET is_deleted=TRUE, updated_at=NOW(), version=version+1
		 WHERE id=$1 AND is_deleted=FALSE`, id)
	return err
}
