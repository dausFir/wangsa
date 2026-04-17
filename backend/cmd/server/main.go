package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wangsa/backend/config"
	delivery "github.com/wangsa/backend/internal/delivery/http"
	"github.com/wangsa/backend/internal/delivery/http/handler"
	"github.com/wangsa/backend/internal/pkg/database"
	jwtutil "github.com/wangsa/backend/internal/pkg/jwt"
	"github.com/wangsa/backend/internal/repository"
	"github.com/wangsa/backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	// ── Startup validation ────────────────────────────────────
	if err := cfg.ValidateDatabaseURL(); err != nil {
		log.Fatalf("❌  Invalid DATABASE_URL: %v", err)
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("❌  JWT_SECRET must be at least 32 characters. Generate: openssl rand -hex 32")
	}
	if cfg.IsProduction && cfg.JWTSecret == "super-secret-jwt-key-change-in-production-please" {
		log.Fatal("❌  Change JWT_SECRET before running in production")
	}

	// ── Database ──────────────────────────────────────────────
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌  Database connection failed: %v", err)
	}
	defer db.Close()

	schemaPath := os.Getenv("SCHEMA_PATH")
	if err := database.RunMigrations(db, schemaPath); err != nil {
		log.Fatalf("❌  Migration failed: %v", err)
	}
	log.Println("✅  Database migrations applied")

	// ── JWT ───────────────────────────────────────────────────
	jm := jwtutil.NewManager(cfg.JWTSecret, cfg.JWTExpiresIn)

	// ── Repositories ─────────────────────────────────────────
	userRepo     := repository.NewUserRepository(db)
	refreshRepo  := repository.NewRefreshTokenRepository(db)
	familyRepo  := repository.NewFamilyRepository(db)
	kasRepo     := repository.NewKasRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	eventRepo    := repository.NewEventRepository(db)
	attendeeRepo := repository.NewAttendeeRepository(db)

	// ── Use Cases ────────────────────────────────────────────
	authUC   := usecase.NewAuthUsecase(userRepo, refreshRepo, jm)
	familyUC := usecase.NewFamilyUsecase(familyRepo)

	// ── Handlers ─────────────────────────────────────────────
	authHandler    := handler.NewAuthHandler(authUC, cfg, userRepo)
	familyHandler  := handler.NewFamilyHandler(familyUC, familyRepo)
	kasHandler     := handler.NewKasHandler(kasRepo)
	addressHandler := handler.NewAddressHandler(addressRepo)
	eventHandler   := handler.NewEventHandler(eventRepo)
	uploadHandler   := handler.NewUploadHandler(familyRepo)
	attendeeHandler := handler.NewAttendeeHandler(attendeeRepo, eventRepo)

	// ── HTTP Server ──────────────────────────────────────────
	// db is passed to router so /health can do a real DB ping
	router := delivery.NewRouter(cfg, jm, db,
		authHandler, familyHandler, kasHandler,
		addressHandler, eventHandler, uploadHandler, attendeeHandler,
	)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router.Engine(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀  Wangsa API listening on http://localhost:%s", cfg.ServerPort)
		log.Printf("🌐  CORS allowed for: %s", cfg.FrontendURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌  Server error: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏳  Shutting down gracefully (up to 10s)...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️   Forced shutdown: %v", err)
	}
	log.Println("✅  Server stopped cleanly")
}
