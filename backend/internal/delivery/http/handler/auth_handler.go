package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/delivery/http/middleware"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/response"
	"github.com/wangsa/backend/internal/usecase"
)

type AuthHandler struct {
	uc       *usecase.AuthUsecase
	cfg      *config.Config
	userRepo domain.UserRepository
}

func NewAuthHandler(uc *usecase.AuthUsecase, cfg *config.Config, userRepo domain.UserRepository) *AuthHandler {
	return &AuthHandler{uc: uc, cfg: cfg, userRepo: userRepo}
}

// Register endpoint - enabled for development
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, pair, err := h.uc.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.setTokenCookies(c, pair)
	response.Created(c, domain.AuthResponse{User: *user}, "Registrasi berhasil")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, pair, err := h.uc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: err.Error()})
		return
	}
	h.setTokenCookies(c, pair)
	response.OK(c, domain.AuthResponse{User: *user}, "Login berhasil")
}

// Refresh exchanges a valid refresh token for a new access + refresh token pair.
// Called automatically by the frontend when it receives a 401.
func (h *AuthHandler) Refresh(c *gin.Context) {
	rawRefresh, err := c.Cookie(middleware.RefreshCookieName)
	if err != nil || rawRefresh == "" {
		response.Unauthorized(c)
		return
	}
	user, pair, err := h.uc.Refresh(rawRefresh)
	if err != nil {
		// Refresh token invalid or expired — clear cookies and force re-login
		h.clearTokenCookies(c)
		c.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   "Sesi habis. Silakan login kembali.",
		})
		return
	}
	h.setTokenCookies(c, pair)
	response.OK(c, domain.AuthResponse{User: *user}, "Token diperbarui")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Best-effort: revoke refresh token in DB (don't fail if user_id missing)
	if userIDVal, exists := c.Get(middleware.ContextUserID); exists {
		if userID, ok := userIDVal.(int64); ok && userID != 0 {
			_ = h.uc.Logout(userID)
		}
	}
	h.clearTokenCookies(c)
	response.OK(c, nil, "Logout berhasil")
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c)
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok || userID == 0 {
		response.Unauthorized(c)
		return
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		response.Unauthorized(c)
		return
	}
	response.OK(c, user)
}

// setTokenCookies writes both access and refresh tokens as HttpOnly cookies.
func (h *AuthHandler) setTokenCookies(c *gin.Context, pair *usecase.TokenPair) {
	c.SetSameSite(http.SameSiteLaxMode)

	// Access token — short TTL (15 min)
	accessMaxAge := int(h.cfg.JWTExpiresIn / time.Second)
	c.SetCookie(
		middleware.CookieName, pair.AccessToken, accessMaxAge,
		"/", h.cfg.CookieDomain, h.cfg.IsProduction, true,
	)

	// Refresh token — long TTL (30 days), restricted to /api/auth path
	// Path restriction limits the cookie's exposure surface
	c.SetCookie(
		middleware.RefreshCookieName, pair.RefreshToken, pair.RefreshTokenExpiry,
		"/api/auth", h.cfg.CookieDomain, h.cfg.IsProduction, true,
	)
}

// clearTokenCookies removes both cookies by setting MaxAge=-1.
func (h *AuthHandler) clearTokenCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.CookieName, "", -1, "/", h.cfg.CookieDomain, h.cfg.IsProduction, true)
	c.SetCookie(middleware.RefreshCookieName, "", -1, "/api/auth", h.cfg.CookieDomain, h.cfg.IsProduction, true)
}
