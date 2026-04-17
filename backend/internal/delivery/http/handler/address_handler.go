package handler

import (
	"github.com/wangsa/backend/internal/delivery/http/middleware"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	repo domain.AddressRepository
}

func NewAddressHandler(repo domain.AddressRepository) *AddressHandler {
	return &AddressHandler{repo: repo}
}

func (h *AddressHandler) List(c *gin.Context) {
	addrs, err := h.repo.FindAll()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, addrs)
}

func (h *AddressHandler) Create(c *gin.Context) {
	var req domain.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := contextUserID(c)
	country := req.Country
	if country == "" {
		country = "Indonesia"
	}
	a := &domain.Address{
		FamilyMemberID: req.FamilyMemberID,
		Label:          req.Label,
		Street:         req.Street,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
		Country:        country,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		IsCurrent:      req.IsCurrent,
		CreatedBy:      userID,
		UpdatedBy:      userID,
	}
	if err := h.repo.Create(a); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Created(c, a, "Alamat berhasil ditambahkan")
}

func (h *AddressHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	existing, err := h.repo.FindByID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if existing == nil {
		response.NotFound(c, "address")
		return
	}
	var req domain.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := contextUserID(c)
	existing.Label      = req.Label
	existing.Street     = req.Street
	existing.City       = req.City
	existing.Province   = req.Province
	existing.PostalCode = req.PostalCode
	existing.Latitude   = req.Latitude
	existing.Longitude  = req.Longitude
	existing.IsCurrent  = req.IsCurrent
	existing.UpdatedBy  = userID
	if req.Country != "" {
		existing.Country = req.Country
	}
	if err := h.repo.Update(existing); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, existing, "Alamat diperbarui")
}

func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Alamat dihapus")
}

// contextUserID extracts the authenticated user's ID from gin context.
// Returns nil if not present (safe — audit fields are nullable).
func contextUserID(c *gin.Context) *int64 {
	v, exists := c.Get(middleware.ContextUserID)
	if !exists {
		return nil
	}
	id, ok := v.(int64)
	if !ok || id == 0 {
		return nil
	}
	return &id
}
