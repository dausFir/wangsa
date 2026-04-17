package response

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// isProd is evaluated once at startup — avoids parsing env on every request.
var isProd = os.Getenv("PRODUCTION") == "true"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *gin.Context, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusOK, Response{Success: true, Message: msg, Data: data})
}

func Created(c *gin.Context, data interface{}, message ...string) {
	msg := "created"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusCreated, Response{Success: true, Message: msg, Data: data})
}

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Error: err})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{Success: false, Error: "unauthorized"})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{Success: false, Error: "access forbidden"})
}

func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, Response{Success: false, Error: resource + " not found"})
}

// InternalError logs the full error server-side via slog (structured JSON in prod),
// but only returns details to the client in development mode.
func InternalError(c *gin.Context, err error) {
	slog.Error("internal server error",
		"method", c.Request.Method,
		"path",   c.Request.URL.Path,
		"error",  err.Error(),
	)

	msg := err.Error()
	if isProd {
		msg = "Terjadi kesalahan pada server. Silakan coba lagi."
	}
	c.JSON(http.StatusInternalServerError, Response{Success: false, Error: msg})
}
