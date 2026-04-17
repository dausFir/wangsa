package handler

import (
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AttendeeHandler struct {
	repo      domain.AttendeeRepository
	eventRepo domain.EventRepository
}

func NewAttendeeHandler(repo domain.AttendeeRepository, eventRepo domain.EventRepository) *AttendeeHandler {
	return &AttendeeHandler{repo: repo, eventRepo: eventRepo}
}

// GET /api/events/:id/attendees
func (h *AttendeeHandler) List(c *gin.Context) {
	eventID, err := parseID(c, "id")
	if err != nil {
		return
	}
	attendees, err := h.repo.ListAttendees(eventID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, attendees)
}

// PUT /api/events/:id/attendees/:member_id
// Body: { "rsvp": "yes" | "no" | "pending" }
func (h *AttendeeHandler) Upsert(c *gin.Context) {
	eventID, err := parseID(c, "id")
	if err != nil {
		return
	}
	memberID, err := parseID(c, "member_id")
	if err != nil {
		return
	}

	// Verify event exists
	event, err := h.eventRepo.FindByID(eventID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if event == nil {
		response.NotFound(c, "event")
		return
	}

	var req domain.RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.repo.UpsertAttendee(eventID, memberID, req.RSVP); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, gin.H{"event_id": eventID, "family_member_id": memberID, "rsvp": req.RSVP}, "RSVP berhasil diperbarui")
}

// DELETE /api/events/:id/attendees/:member_id
func (h *AttendeeHandler) Remove(c *gin.Context) {
	eventID, err := parseID(c, "id")
	if err != nil {
		return
	}
	memberID, err := parseID(c, "member_id")
	if err != nil {
		return
	}
	if err := h.repo.RemoveAttendee(eventID, memberID); err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, nil, "Peserta berhasil dihapus dari acara")
}
