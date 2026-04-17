package usecase

import (
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/sanitize"
)

type FamilyUsecase struct {
	repo domain.FamilyRepository
}

// validateDates checks that birth_date < death_date when both are provided.
func validateDates(req *domain.CreateFamilyMemberRequest) error {
	if req.BirthDate == nil || req.DeathDate == nil {
		return nil
	}
	if *req.BirthDate == "" || *req.DeathDate == "" {
		return nil
	}
	birth, err := time.Parse("2006-01-02", *req.BirthDate)
	if err != nil {
		return fmt.Errorf("format birth_date tidak valid, gunakan YYYY-MM-DD")
	}
	death, err := time.Parse("2006-01-02", *req.DeathDate)
	if err != nil {
		return fmt.Errorf("format death_date tidak valid, gunakan YYYY-MM-DD")
	}
	if !death.After(birth) {
		return fmt.Errorf("tanggal meninggal harus setelah tanggal lahir")
	}
	return nil
}

func NewFamilyUsecase(repo domain.FamilyRepository) *FamilyUsecase {
	return &FamilyUsecase{repo: repo}
}

// BuildFamilyTree constructs a fully-nested JSON tree in O(n) time.
func (u *FamilyUsecase) BuildFamilyTree() ([]*domain.FamilyMember, error) {
	allMembers, err := u.repo.FindAllMembers()
	if err != nil {
		return nil, fmt.Errorf("fetch members: %w", err)
	}
	allMarriages, err := u.repo.FindAllMarriages()
	if err != nil {
		return nil, fmt.Errorf("fetch marriages: %w", err)
	}

	memberMap := make(map[int64]*domain.FamilyMember, len(allMembers))
	for _, m := range allMembers {
		m.Children = []*domain.FamilyMember{}
		m.Spouses = []*domain.FamilyMember{}
		m.Marriages = []*domain.Marriage{}
		memberMap[m.ID] = m
	}

	for _, marriage := range allMarriages {
		husband, hOk := memberMap[marriage.HusbandID]
		wife, wOk := memberMap[marriage.WifeID]
		if hOk && wOk {
			husband.Spouses = append(husband.Spouses, wife)
			wife.Spouses = append(wife.Spouses, husband)
			husband.Marriages = append(husband.Marriages, marriage)
			wife.Marriages = append(wife.Marriages, marriage)
		}
	}

	roots := make([]*domain.FamilyMember, 0)
	for _, m := range allMembers {
		if m.ParentID != nil {
			if parent, ok := memberMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, m)
				continue
			}
		}
		roots = append(roots, m)
	}

	return roots, nil
}

func (u *FamilyUsecase) CreateMember(req *domain.CreateFamilyMemberRequest, createdBy int64) (*domain.FamilyMember, error) {
	if err := validateDates(req); err != nil {
		return nil, err
	}
	if req.PhotoURL != nil && *req.PhotoURL != "" {
		if err := sanitize.PhotoURL(*req.PhotoURL); err != nil {
			return nil, fmt.Errorf("photo_url tidak valid: %w", err)
		}
	}

	m := &domain.FamilyMember{
		FullName:   req.FullName,
		Nickname:   req.Nickname,
		Gender:     req.Gender,
		BirthDate:  req.BirthDate,
		BirthPlace: req.BirthPlace,
		DeathDate:  req.DeathDate,
		PhotoURL:   req.PhotoURL,
		ParentID:   req.ParentID,
		Notes:      req.Notes,
		CreatedBy:  &createdBy,
	}
	if err := u.repo.CreateMember(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (u *FamilyUsecase) UpdateMember(id int64, req *domain.CreateFamilyMemberRequest) (*domain.FamilyMember, error) {
	if err := validateDates(req); err != nil {
		return nil, err
	}
	if req.PhotoURL != nil && *req.PhotoURL != "" {
		if err := sanitize.PhotoURL(*req.PhotoURL); err != nil {
			return nil, fmt.Errorf("photo_url tidak valid: %w", err)
		}
	}

	m, err := u.repo.FindMemberByID(id)
	if err != nil {
		return nil, fmt.Errorf("find member: %w", err)
	}
	if m == nil {
		return nil, fmt.Errorf("member not found")
	}

	// Update non-photo fields
	m.FullName = req.FullName
	m.Nickname = req.Nickname
	m.Gender = req.Gender
	m.BirthDate = req.BirthDate
	m.BirthPlace = req.BirthPlace
	m.DeathDate = req.DeathDate
	m.ParentID = req.ParentID
	m.Notes = req.Notes

	// Only update PhotoURL if explicitly provided in request
	// This prevents accidentally clearing the photo_url when it's not included in the request
	if req.PhotoURL != nil {
		m.PhotoURL = req.PhotoURL
	}

	if err := u.repo.UpdateMember(m); err != nil {
		return nil, err
	}
	return m, nil
}
