package sanitize

import (
	"fmt"
	"net/url"
	"strings"
)

// PhotoURL validates that a photo_url value is safe to store and render.
//
// Allowed schemes: https, http, data:image/...;base64,...
// Blocked: javascript:, vbscript:, data:text/html, blob:, file:, etc.
func PhotoURL(raw string) error {
	if raw == "" {
		return nil
	}

	lower := strings.ToLower(strings.TrimSpace(raw))

	// Allow base64 data URLs for safe image MIME types only
	if strings.HasPrefix(lower, "data:") {
		allowed := []string{
			"data:image/jpeg;base64,",
			"data:image/jpg;base64,",
			"data:image/png;base64,",
			"data:image/gif;base64,",
			"data:image/webp;base64,",
		}
		for _, prefix := range allowed {
			if strings.HasPrefix(lower, prefix) {
				if len(raw) <= len(prefix) {
					return fmt.Errorf("data URL memiliki payload kosong")
				}
				return nil
			}
		}
		return fmt.Errorf("tipe data URL tidak diizinkan — hanya data:image/(jpeg|png|gif|webp);base64")
	}

	// For regular URLs: parse and validate scheme
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("photo_url bukan URL yang valid")
	}

	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "https", "http":
		if u.Host == "" {
			return fmt.Errorf("photo_url harus memiliki host")
		}
		return nil
	default:
		return fmt.Errorf("skema URL '%s' tidak diizinkan pada photo_url", u.Scheme)
	}
}
