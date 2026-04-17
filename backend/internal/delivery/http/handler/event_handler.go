package handler

import (
	"fmt"
	"time"

	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	repo domain.EventRepository
}

func NewEventHandler(repo domain.EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

// parseDateTime tries RFC3339 first, then the HTML datetime-local format.
// Returns an error only when StartAt fails — EndAt failures are silently ignored
// because the field is optional and the frontend may send an empty string.
func parseDateTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("format waktu tidak valid, gunakan RFC3339 atau YYYY-MM-DDTHH:MM")
}

func parseDateTimeOptional(s string) *time.Time {
	t, err := parseDateTime(s)
	if err != nil {
		return nil
	}
	return &t
}

// GET /api/events?from=RFC3339&to=RFC3339
func (h *EventHandler) List(c *gin.Context) {
	from := c.Query("from")
	to   := c.Query("to")
	events, err := h.repo.FindAll(from, to)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, events)
}

// POST /api/events
func (h *EventHandler) Create(c *gin.Context) {
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	startAt, err := parseDateTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	createdBy := contextUserID(c)
	color := req.Color
	if color == "" {
		color = "#CC6649"
	}

	event := &domain.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartAt:     startAt,
		IsRecurring: req.IsRecurring,
		RecurRule:   req.RecurRule,
		Color:       color,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if req.EndAt != nil && *req.EndAt != "" {
		event.EndAt = parseDateTimeOptional(*req.EndAt)
	}

	if err := h.repo.Create(event); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Created(c, event, "Acara berhasil dibuat")
}

// PUT /api/events/:id
func (h *EventHandler) Update(c *gin.Context) {
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
		response.NotFound(c, "event")
		return
	}

	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	startAt, err := parseDateTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	existing.Title       = req.Title
	existing.Description = req.Description
	existing.Location    = req.Location
	existing.StartAt     = startAt
	existing.IsRecurring = req.IsRecurring
	existing.RecurRule   = req.RecurRule
	existing.UpdatedBy   = contextUserID(c)
	if req.Color != "" {
		existing.Color = req.Color
	}
	if req.EndAt != nil && *req.EndAt != "" {
		existing.EndAt = parseDateTimeOptional(*req.EndAt)
	}

	if err := h.repo.Update(existing); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, existing, "Acara diperbarui")
}

// DELETE /api/events/:id
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Acara dihapus")
}
