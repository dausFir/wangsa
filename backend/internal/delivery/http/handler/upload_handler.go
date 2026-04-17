package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	// TODO: Enable when Go version supports webp
	// _ "image/webp"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
)

const (
	maxUploadSize = 5 << 20 // 5 MB
	targetSize    = 400     // max dimension in pixels
	jpegQuality   = 82
)

type UploadHandler struct {
	familyRepo domain.FamilyRepository
}

func NewUploadHandler(repo domain.FamilyRepository) *UploadHandler {
	return &UploadHandler{familyRepo: repo}
}

// POST /api/family/members/:id/photo
// Accepts multipart/form-data field "photo" (JPEG/PNG/GIF, max 5MB).
// Resizes to fit 400×400, re-encodes as JPEG, stores as base64 data URL in photo_url.
func (h *UploadHandler) UploadMemberPhoto(c *gin.Context) {
	memberID, err := parseID(c, "id")
	if err != nil {
		return
	}

	member, err := h.familyRepo.FindMemberByID(memberID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if member == nil {
		response.NotFound(c, "family member")
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	file, _, err := c.Request.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "file 'photo' wajib diisi (max 5MB)")
		return
	}
	defer file.Close()

	// Sniff first 512 bytes to detect real MIME type — cannot be spoofed via headers
	sniff := make([]byte, 512)
	n, _ := file.Read(sniff)
	detected := http.DetectContentType(sniff[:n])
	allowed := map[string]bool{
		"image/jpeg": true, "image/png": true,
		"image/gif": true, "image/webp": true,
	}
	if !allowed[detected] {
		response.BadRequest(c, "hanya file gambar yang diizinkan (JPEG, PNG, GIF, WEBP)")
		return
	}
	// Seek back to beginning so image.Decode sees the full file
	if seeker, ok := file.(interface {
		Seek(int64, int) (int64, error)
	}); ok {
		seeker.Seek(0, 0)
	}

	// Decode image — standard library handles JPEG, PNG, GIF via blank imports
	img, _, err := image.Decode(file)
	if err != nil {
		response.BadRequest(c, "file gambar tidak valid atau rusak")
		return
	}

	// Resize proportionally so neither dimension exceeds 400px
	resized := resizeImage(img, targetSize, targetSize)

	// Encode resized image → JPEG bytes
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: jpegQuality}); err != nil {
		response.InternalError(c, fmt.Errorf("encode image: %w", err))
		return
	}

	// Build data URL
	dataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	// Persist to DB
	member.PhotoURL = &dataURL
	if err := h.familyRepo.UpdateMember(member); err != nil {
		response.InternalError(c, fmt.Errorf("save photo: %w", err))
		return
	}

	response.OK(c, gin.H{"photo_url": dataURL}, "Foto berhasil diperbarui")
}

// DELETE /api/family/members/:id/photo
func (h *UploadHandler) DeleteMemberPhoto(c *gin.Context) {
	memberID, err := parseID(c, "id")
	if err != nil {
		return
	}
	member, err := h.familyRepo.FindMemberByID(memberID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if member == nil {
		response.NotFound(c, "family member")
		return
	}
	member.PhotoURL = nil
	if err := h.familyRepo.UpdateMember(member); err != nil {
		response.InternalError(c, fmt.Errorf("remove photo: %w", err))
		return
	}
	response.OK(c, nil, "Foto berhasil dihapus")
}

// resizeImage resizes img so neither dimension exceeds maxW×maxH,
// preserving aspect ratio via nearest-neighbour (no third-party deps).
func resizeImage(img image.Image, maxW, maxH int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= maxW && srcH <= maxH {
		return img
	}

	scaleW := float64(maxW) / float64(srcW)
	scaleH := float64(maxH) / float64(srcH)
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}

	dstW := int(float64(srcW) * scale)
	dstH := int(float64(srcH) * scale)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	return dst
}
