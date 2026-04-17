package handler

import (
	"strconv"

	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/wangsa/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FamilyHandler struct {
	uc   *usecase.FamilyUsecase
	repo domain.FamilyRepository
}

func NewFamilyHandler(uc *usecase.FamilyUsecase, repo domain.FamilyRepository) *FamilyHandler {
	return &FamilyHandler{uc: uc, repo: repo}
}

// GET /api/family/tree
func (h *FamilyHandler) GetTree(c *gin.Context) {
	tree, err := h.uc.BuildFamilyTree()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, tree)
}

// GET /api/family/members
func (h *FamilyHandler) ListMembers(c *gin.Context) {
	members, err := h.repo.FindAllMembers()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, members)
}

// GET /api/family/members/:id
func (h *FamilyHandler) GetMember(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	member, err := h.repo.FindMemberByID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if member == nil {
		response.NotFound(c, "family member")
		return
	}
	response.OK(c, member)
}

// POST /api/family/members
func (h *FamilyHandler) CreateMember(c *gin.Context) {
	var req domain.CreateFamilyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var createdBy int64
	if id := contextUserID(c); id != nil { createdBy = *id }
	member, err := h.uc.CreateMember(&req, createdBy)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Created(c, member, "Anggota keluarga berhasil ditambahkan")
}

// PUT /api/family/members/:id
func (h *FamilyHandler) UpdateMember(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	var req domain.CreateFamilyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	member, err := h.uc.UpdateMember(id, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, member, "Data anggota diperbarui")
}

// DELETE /api/family/members/:id
func (h *FamilyHandler) DeleteMember(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.repo.DeleteMember(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Anggota dihapus")
}

// POST /api/family/marriages
func (h *FamilyHandler) CreateMarriage(c *gin.Context) {
	var req domain.CreateMarriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.HusbandID == req.WifeID {
		response.BadRequest(c, "husband_id dan wife_id tidak boleh sama")
		return
	}
	marriage := &domain.Marriage{
		HusbandID:    req.HusbandID,
		WifeID:       req.WifeID,
		MarriageDate: req.MarriageDate,
		DivorceDate:  req.DivorceDate,
		Notes:        req.Notes,
	}
	if err := h.repo.CreateMarriage(marriage); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, marriage, "Data pernikahan ditambahkan")
}

// DELETE /api/family/marriages/:id
func (h *FamilyHandler) DeleteMarriage(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.repo.DeleteMarriage(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Data pernikahan dihapus")
}

// GET /api/family/members/:id/marriages
func (h *FamilyHandler) GetMemberMarriages(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	marriages, err := h.repo.FindMarriagesByMemberID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, marriages)
}

// helper: parse int64 path param and write BadRequest on failure
func parseID(c *gin.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid "+param)
		return 0, err
	}
	return id, nil
}
