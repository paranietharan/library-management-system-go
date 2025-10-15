package main

import (
	"library-management-system-go/internal/config"
	"library-management-system-go/internal/database"
	"library-management-system-go/internal/database/migration"
	"log"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := migration.RunMigrations(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "rollback" {
		if err := migration.RollbackMigration(cfg); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("Rollback completed")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := database.SeedDefaultUsers(cfg); err != nil {
			log.Printf("Warning: Failed to seed default users: %v", err)
		}
		log.Println("Seeding completed ....")
		return
	}
}
