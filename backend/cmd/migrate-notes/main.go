package migratenotes
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

	// Add notes column to events table if it doesn't exist
	_, err = db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'events' AND column_name = 'notes'
			) THEN
				ALTER TABLE events ADD COLUMN notes TEXT;
			END IF;
		END $$;
	`)
	if err != nil {
		log.Fatalf("Failed to add notes column: %v", err)
	}

	log.Println("✅ Successfully added notes column to events table")
}