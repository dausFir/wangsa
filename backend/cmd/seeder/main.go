package main

import (
	"flag"
	"log"
	"os"

	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/database"
	"github.com/wangsa/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var createSuperAdmin = flag.Bool("create-admin", false, "Create default superadmin account")
	flag.Parse()

	if !*createSuperAdmin {
		log.Println("Usage: go run cmd/seeder/main.go -create-admin")
		os.Exit(1)
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations first
	if err := database.RunMigrations(db, ""); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Create superadmin
	if err := createSuperAdminAccount(userRepo); err != nil {
		log.Fatalf("Failed to create superadmin account: %v", err)
	}

	log.Println("✅ Superadmin account created successfully!")
	log.Println("📧 Email: admin@wangsa.internal")
	log.Println("🔑 Password: WangsaAdmin2024!")
	log.Println("⚠️  Please change the password after first login")
}

func createSuperAdminAccount(userRepo domain.UserRepository) error {
	// Check if admin already exists
	existingAdmin, err := userRepo.FindByEmail("admin@wangsa.internal")
	if err != nil {
		return err
	}

	if existingAdmin != nil {
		log.Println("⚠️  Superadmin account already exists, skipping creation")
		return nil
	}

	// Hash password
	password := "WangsaAdmin2024!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create superadmin user
	admin := &domain.User{
		Name:     "Super Administrator",
		Email:    "admin@wangsa.internal",
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	return userRepo.Create(admin)
}
