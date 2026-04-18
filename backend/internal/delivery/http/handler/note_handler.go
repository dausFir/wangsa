package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
)

type NoteHandler struct {
	repo domain.NoteRepository
}

func NewNoteHandler(repo domain.NoteRepository) *NoteHandler {
	return &NoteHandler{repo: repo}
}

// GET /api/notes?category=string
func (h *NoteHandler) List(c *gin.Context) {
	category := c.Query("category")
	notes, err := h.repo.FindAll(category)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, notes)
}

// GET /api/notes/categories
func (h *NoteHandler) ListCategories(c *gin.Context) {
	categories, err := h.repo.FindCategories()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, categories)
}

// GET /api/notes/:id
func (h *NoteHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid note ID")
		return
	}

	note, err := h.repo.FindByID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if note == nil {
		response.NotFound(c, "note")
		return
	}
	response.OK(c, note)
}

// POST /api/notes
func (h *NoteHandler) Create(c *gin.Context) {
	var req domain.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := contextUserID(c)
	note := &domain.Note{
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		IsPinned:  req.IsPinned,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := h.repo.Create(note); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Created(c, note, "Catatan berhasil dibuat")
}

// PUT /api/notes/:id
func (h *NoteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid note ID")
		return
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if existing == nil {
		response.NotFound(c, "note")
		return
	}

	var req domain.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	existing.Title = req.Title
	existing.Content = req.Content
	existing.Category = req.Category
	existing.IsPinned = req.IsPinned
	existing.UpdatedBy = contextUserID(c)

	if err := h.repo.Update(existing); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, existing, "Catatan berhasil diperbarui")
}

// DELETE /api/notes/:id
func (h *NoteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid note ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Catatan berhasil dihapus")
}
