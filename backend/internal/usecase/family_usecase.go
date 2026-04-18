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

// TreeNode represents a family tree node for JSON serialization without circular references
type TreeNode struct {
	ID           int64      `json:"id"`
	FullName     string     `json:"full_name"`
	Gender       string     `json:"gender"`
	BirthDate    *string    `json:"birth_date,omitempty"`
	DeathDate    *string    `json:"death_date,omitempty"`
	ParentID     *int64     `json:"parent_id,omitempty"`
	SpouseIDs    []int64    `json:"spouse_ids,omitempty"`
	Children     []TreeNode `json:"children,omitempty"`
	MarriageInfo []struct {
		SpouseID     int64   `json:"spouse_id"`
		SpouseName   string  `json:"spouse_name"`
		MarriageDate *string `json:"marriage_date,omitempty"`
	} `json:"marriage_info,omitempty"`
}

// BuildFamilyTree constructs a tree using safe TreeNode structure
func (u *FamilyUsecase) BuildFamilyTree() ([]TreeNode, error) {
	allMembers, err := u.repo.FindAllMembers()
	if err != nil {
		return nil, fmt.Errorf("fetch members: %w", err)
	}
	fmt.Printf("DEBUG: Found %d members\n", len(allMembers))

	allMarriages, err := u.repo.FindAllMarriages()
	if err != nil {
		return nil, fmt.Errorf("fetch marriages: %w", err)
	}
	fmt.Printf("DEBUG: Found %d marriages\n", len(allMarriages))

	// Create member map and node map
	memberMap := make(map[int64]*domain.FamilyMember)
	nodeMap := make(map[int64]*TreeNode)

	// First pass: create all nodes
	for _, m := range allMembers {
		memberMap[m.ID] = m
		nodeMap[m.ID] = &TreeNode{
			ID:        m.ID,
			FullName:  m.FullName,
			Gender:    m.Gender,
			BirthDate: m.BirthDate,
			DeathDate: m.DeathDate,
			ParentID:  m.ParentID,
			SpouseIDs: []int64{},
			Children:  []TreeNode{},
			MarriageInfo: []struct {
				SpouseID     int64   `json:"spouse_id"`
				SpouseName   string  `json:"spouse_name"`
				MarriageDate *string `json:"marriage_date,omitempty"`
			}{},
		}
	}

	// Second pass: Add marriage relationships
	for _, marriage := range allMarriages {
		if hNode, hOk := nodeMap[marriage.HusbandID]; hOk {
			if wNode, wOk := nodeMap[marriage.WifeID]; wOk {
				// Add spouse IDs
				hNode.SpouseIDs = append(hNode.SpouseIDs, marriage.WifeID)
				wNode.SpouseIDs = append(wNode.SpouseIDs, marriage.HusbandID)

				// Add marriage info
				hNode.MarriageInfo = append(hNode.MarriageInfo, struct {
					SpouseID     int64   `json:"spouse_id"`
					SpouseName   string  `json:"spouse_name"`
					MarriageDate *string `json:"marriage_date,omitempty"`
				}{
					SpouseID:     marriage.WifeID,
					SpouseName:   wNode.FullName,
					MarriageDate: marriage.MarriageDate,
				})

				wNode.MarriageInfo = append(wNode.MarriageInfo, struct {
					SpouseID     int64   `json:"spouse_id"`
					SpouseName   string  `json:"spouse_name"`
					MarriageDate *string `json:"marriage_date,omitempty"`
				}{
					SpouseID:     marriage.HusbandID,
					SpouseName:   hNode.FullName,
					MarriageDate: marriage.MarriageDate,
				})
			}
		}
	}

	// Third pass: Build parent-child relationships and identify roots
	var rootIDs []int64
	for _, member := range allMembers {
		if member.ParentID != nil {
			// This member has a parent - add to parent's children
			if parentNode, ok := nodeMap[*member.ParentID]; ok {
				if childNode, ok := nodeMap[member.ID]; ok {
					parentNode.Children = append(parentNode.Children, *childNode)
				}
			}
		} else {
			// This member has no parent - it's a root
			rootIDs = append(rootIDs, member.ID)
		}
	}

	// Fourth pass: Recursively build full tree structure for roots
	var buildNodeTree func(*TreeNode) TreeNode
	buildNodeTree = func(node *TreeNode) TreeNode {
		// Create a copy of the node
		result := *node
		result.Children = []TreeNode{}

		// Find all children of this node
		for _, member := range allMembers {
			if member.ParentID != nil && *member.ParentID == node.ID {
				if childNode, ok := nodeMap[member.ID]; ok {
					childTree := buildNodeTree(childNode)
					result.Children = append(result.Children, childTree)
				}
			}
		}
		return result
	}

	// Build the final root trees
	var finalRoots []TreeNode
	for _, rootID := range rootIDs {
		if rootNode, ok := nodeMap[rootID]; ok {
			finalTree := buildNodeTree(rootNode)
			finalRoots = append(finalRoots, finalTree)
		}
	}

	fmt.Printf("DEBUG: Found %d root members\n", len(finalRoots))

	// Debug: print tree structure for first root
	if len(finalRoots) > 0 {
		fmt.Printf("DEBUG: First root '%s' has %d children\n", finalRoots[0].FullName, len(finalRoots[0].Children))
	}

	return finalRoots, nil
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
