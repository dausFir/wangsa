package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/domain"
	"github.com/wangsa/backend/internal/pkg/database"
	"github.com/wangsa/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var createSuperAdmin = flag.Bool("create-admin", false, "Create default superadmin account")
	var createSampleEvents = flag.Bool("create-events", false, "Create sample events for testing")
	flag.Parse()

	if !*createSuperAdmin && !*createSampleEvents {
		log.Println("Usage:")
		log.Println("  go run cmd/seeder/main.go -create-admin    # Create superadmin account")
		log.Println("  go run cmd/seeder/main.go -create-events   # Create sample events")
		log.Println("  go run cmd/seeder/main.go -create-admin -create-events   # Create both")
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Create superadmin
	if *createSuperAdmin {
		if err := createSuperAdminAccount(userRepo); err != nil {
			log.Fatalf("Failed to create superadmin account: %v", err)
		}
		log.Println("✅ Superadmin account created successfully!")
		log.Println("📧 Email: admin@wangsa.internal")
		log.Println("🔑 Password: WangsaAdmin2024!")
		log.Println("⚠️  Please change the password after first login")
	}

	// Create sample events
	if *createSampleEvents {
		if err := createSampleEventsData(eventRepo, userRepo); err != nil {
			log.Fatalf("Failed to create sample events: %v", err)
		}
		log.Println("✅ Sample events created successfully!")
	}
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

func createSampleEventsData(eventRepo domain.EventRepository, userRepo domain.UserRepository) error {
	// Get admin user to use as creator
	admin, err := userRepo.FindByEmail("admin@wangsa.internal")
	if err != nil {
		return err
	}
	if admin == nil {
		log.Println("⚠️  Superadmin not found, creating sample events with user ID 1")
		admin = &domain.User{ID: 1} // Fallback to ID 1
	}

	// Get current date for creating events around this month
	now := time.Now()

	sampleEvents := []*domain.Event{
		{
			Title:       "Reuni Keluarga Besar",
			Description: strPtr("Acara reuni tahunan keluarga besar di rumah nenek"),
			Location:    strPtr("Rumah Nenek, Yogyakarta"),
			StartAt:     time.Date(now.Year(), now.Month(), 15, 10, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 15, 17, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#CC6649",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Ulang Tahun Mama",
			Description: strPtr("Perayaan ulang tahun mama yang ke-65"),
			Location:    strPtr("Rumah Mama"),
			StartAt:     time.Date(now.Year(), now.Month(), 20, 18, 30, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 20, 22, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#10B981",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Arisan RT",
			Description: strPtr("Arisan bulanan RT 05/RW 02"),
			Location:    strPtr("Balai RT"),
			StartAt:     time.Date(now.Year(), now.Month(), 25, 19, 0, 0, 0, now.Location()),
			IsRecurring: false,
			Color:       "#3B82F6",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Pengajian Keluarga",
			Description: strPtr("Pengajian rutin keluarga setiap akhir bulan"),
			Location:    strPtr("Masjid Al-Ikhlas"),
			StartAt:     time.Date(now.Year(), now.Month(), 28, 19, 30, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month(), 28, 21, 0, 0, 0, now.Location())),
			IsRecurring: true,
			Color:       "#8B5CF6",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Rapat Keluarga",
			Description: strPtr("Rapat koordinasi acara akhir tahun"),
			Location:    strPtr("Rumah Pak Lurah"),
			StartAt:     time.Date(now.Year(), now.Month(), 30, 14, 0, 0, 0, now.Location()),
			IsRecurring: false,
			Color:       "#F59E0B",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		// Events for next month
		{
			Title:       "Piknik Keluarga",
			Description: strPtr("Piknik bersama di pantai parangtritis"),
			Location:    strPtr("Pantai Parangtritis"),
			StartAt:     time.Date(now.Year(), now.Month()+1, 5, 8, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month()+1, 5, 18, 0, 0, 0, now.Location())),
			IsRecurring: false,
			Color:       "#06D6A0",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
		{
			Title:       "Workshop Memasak",
			Description: strPtr("Belajar masak gudeg bareng-bareng"),
			Location:    strPtr("Rumah Tante Sari"),
			StartAt:     time.Date(now.Year(), now.Month()+1, 12, 9, 0, 0, 0, now.Location()),
			EndAt:       timePtr(time.Date(now.Year(), now.Month()+1, 12, 14, 0, 0, 0, now.Location())),
			IsRecurring: false,
			Color:       "#FF6B6B",
			CreatedBy:   int64Ptr(admin.ID),
			UpdatedBy:   int64Ptr(admin.ID),
		},
	}

	for _, event := range sampleEvents {
		if err := eventRepo.Create(event); err != nil {
			log.Printf("⚠️  Failed to create event %s: %v", event.Title, err)
		} else {
			log.Printf("✅ Created event: %s", event.Title)
		}
	}

	return nil
}

// Helper functions for pointers
func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func int64Ptr(i int64) *int64 {
	return &i
}
