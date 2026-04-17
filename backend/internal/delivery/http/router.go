package http

import (
	"net/http"

	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/delivery/http/handler"
	"github.com/wangsa/backend/internal/delivery/http/middleware"
	jwtutil "github.com/wangsa/backend/internal/pkg/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine         *gin.Engine
	cfg            *config.Config
	jm             *jwtutil.Manager
	db             *sqlx.DB
	authHandler    *handler.AuthHandler
	familyHandler  *handler.FamilyHandler
	kasHandler     *handler.KasHandler
	addressHandler *handler.AddressHandler
	eventHandler   *handler.EventHandler
	uploadHandler    *handler.UploadHandler
	attendeeHandler *handler.AttendeeHandler
}

func NewRouter(
	cfg *config.Config,
	jm *jwtutil.Manager,
	db *sqlx.DB,
	auth *handler.AuthHandler,
	family *handler.FamilyHandler,
	kas *handler.KasHandler,
	address *handler.AddressHandler,
	event *handler.EventHandler,
	upload *handler.UploadHandler,
	attendee *handler.AttendeeHandler,
) *Router {
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())      // panic recovery — always keep
	engine.Use(middleware.RequestLogger()) // structured JSON/text logger

	r := &Router{
		engine:         engine,
		cfg:            cfg,
		jm:             jm,
		db:             db,
		authHandler:    auth,
		familyHandler:  family,
		kasHandler:     kas,
		addressHandler: address,
		eventHandler:   event,
		uploadHandler:   upload,
		attendeeHandler: attendee,
	}
	r.setupCORS()
	r.registerRoutes()
	return r
}

func (r *Router) setupCORS() {
	r.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", r.cfg.FrontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
}

func (r *Router) registerRoutes() {
	api := r.engine.Group("/api")

	// ── Public ────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimitRegister(), r.authHandler.Register)
		auth.POST("/login",    middleware.RateLimitLogin(),    r.authHandler.Login)
		auth.POST("/refresh",  r.authHandler.Refresh)
	}

	// ── Protected (JWT required) ───────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.Auth(r.jm))
	{
		protected.GET("/auth/me", r.authHandler.Me)
		protected.POST("/auth/logout", r.authHandler.Logout)

		// Silsilah
		family := protected.Group("/family")
		{
			family.GET("/tree", r.familyHandler.GetTree)
			family.GET("/members", r.familyHandler.ListMembers)
			family.GET("/members/:id", r.familyHandler.GetMember)
			family.POST("/members", r.familyHandler.CreateMember)
			family.PUT("/members/:id", r.familyHandler.UpdateMember)
			family.DELETE("/members/:id",
				middleware.RequireRole("super_admin"),
				r.familyHandler.DeleteMember,
			)
			family.GET("/members/:id/marriages", r.familyHandler.GetMemberMarriages)
			family.POST("/marriages", r.familyHandler.CreateMarriage)
			family.DELETE("/marriages/:id",
				middleware.RequireRole("super_admin"),
				r.familyHandler.DeleteMarriage,
			)
			family.POST("/members/:id/photo", r.uploadHandler.UploadMemberPhoto)
			family.DELETE("/members/:id/photo", r.uploadHandler.DeleteMemberPhoto)
		}

		// Kas
		kas := protected.Group("/kas")
		{
			kas.GET("/summary", r.kasHandler.GetSummary)
			kas.GET("/categories", r.kasHandler.ListCategories)
			kas.GET("/transactions", r.kasHandler.ListTransactions)
			kas.POST("/transactions", r.kasHandler.CreateTransaction)
			kas.PUT("/transactions/:id", r.kasHandler.UpdateTransaction)
			kas.DELETE("/transactions/:id",
				middleware.RequireRole("super_admin"),
				r.kasHandler.DeleteTransaction,
			)
		}

		// Peta Domisili
		addr := protected.Group("/addresses")
		{
			addr.GET("", r.addressHandler.List)
			addr.POST("", r.addressHandler.Create)
			addr.PUT("/:id", r.addressHandler.Update)
			addr.DELETE("/:id",
				middleware.RequireRole("super_admin"),
				r.addressHandler.Delete,
			)
		}

		// Kalender
		events := protected.Group("/events")
		{
			events.GET("", r.eventHandler.List)
			events.POST("", r.eventHandler.Create)
			events.PUT("/:id", r.eventHandler.Update)
			events.DELETE("/:id",
				middleware.RequireRole("super_admin"),
				r.eventHandler.Delete,
			)
			// Attendees & RSVP
			events.GET("/:id/attendees", r.attendeeHandler.List)
			events.PUT("/:id/attendees/:member_id", r.attendeeHandler.Upsert)
			events.DELETE("/:id/attendees/:member_id", r.attendeeHandler.Remove)
		}
	}

	// ── Health check — actually tests the DB connection ───────
	r.engine.GET("/health", func(c *gin.Context) {
		if err := r.db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"service": "wangsa",
				"error":   "database unreachable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "wangsa",
			"db":      "connected",
		})
	})
}

func (r *Router) Run() error {
	return r.engine.Run(":" + r.cfg.ServerPort)
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
