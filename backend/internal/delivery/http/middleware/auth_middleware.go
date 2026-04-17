package middleware

import (
	"strings"

	jwtutil "github.com/wangsa/backend/internal/pkg/jwt"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "user_id"
	ContextRole      = "user_role"
	CookieName       = "wangsa_token"         // access token cookie
	RefreshCookieName = "wangsa_refresh"      // refresh token cookie
)

// Auth validates the JWT access token from HttpOnly cookie (primary) or
// Authorization header (fallback for API clients).
func Auth(jm *jwtutil.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		claims, err := jm.Validate(token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

// RequireRole restricts access to users with one of the given roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextRole)
		if !exists {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		role, _ := roleVal.(string)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		response.Forbidden(c)
		c.Abort()
	}
}

func extractToken(c *gin.Context) string {
	if cookie, err := c.Cookie(CookieName); err == nil && cookie != "" {
		return cookie
	}
	if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}
