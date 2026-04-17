package handler

import (
	"strconv"

	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type KasHandler struct {
	repo domain.KasRepository
}

func NewKasHandler(repo domain.KasRepository) *KasHandler {
	return &KasHandler{repo: repo}
}

// GET /api/kas/summary
func (h *KasHandler) GetSummary(c *gin.Context) {
	summary, err := h.repo.GetSummary()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, summary)
}

// GET /api/kas/categories
func (h *KasHandler) ListCategories(c *gin.Context) {
	cats, err := h.repo.FindAllCategories()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, cats)
}

// GET /api/kas/transactions?limit=20&offset=0
func (h *KasHandler) ListTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	txs, err := h.repo.FindAllTransactions(limit, offset)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, txs)
}

// POST /api/kas/transactions
func (h *KasHandler) CreateTransaction(c *gin.Context) {
	var req domain.CreateKasTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	createdBy := contextUserID(c)
	tx := &domain.KasTransaction{
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := h.repo.CreateTransaction(tx); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Created(c, tx, "Transaksi berhasil dicatat")
}

// PUT /api/kas/transactions/:id
func (h *KasHandler) UpdateTransaction(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	existing, err := h.repo.FindTransactionByID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if existing == nil {
		response.NotFound(c, "transaction")
		return
	}
	var req domain.CreateKasTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	existing.CategoryID = req.CategoryID
	existing.Type = req.Type
	existing.Amount = req.Amount
	existing.Description = req.Description
	existing.Date = req.Date
	existing.UpdatedBy = contextUserID(c)
	if err := h.repo.UpdateTransaction(existing); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, existing, "Transaksi diperbarui")
}

// DELETE /api/kas/transactions/:id
func (h *KasHandler) DeleteTransaction(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.repo.DeleteTransaction(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Transaksi dihapus")
}
