package cleanevents
package main

import (
	"log"

	"github.com/wangsa/backend/config"
	"github.com/wangsa/backend/internal/pkg/database"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Delete existing sample events
	_, err = db.Exec(`DELETE FROM events WHERE title IN ('Reuni Keluarga Besar', 'Ulang Tahun Mama', 'Arisan RT', 'Pengajian Keluarga', 'Rapat Keluarga', 'Piknik Keluarga', 'Workshop Memasak')`)
	if err != nil {
		log.Fatalf("Failed to delete existing events: %v", err)
	}

	log.Println("✅ Deleted existing sample events")
	log.Println("💡 Now run: go run cmd/seeder/main.go -create-events")
}