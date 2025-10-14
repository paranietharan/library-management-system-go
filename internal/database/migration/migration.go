package migration

import (
	"database/sql"
	"fmt"
	"library-management-system-go/internal/config"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationPath, err := getMigrationPath()
	if err != nil {
		return err
	}

	log.Printf("Migration Path : %s", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

func RollbackMigration(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationPath, err := getMigrationPath()
	if err != nil {
		return err
	}

	log.Printf("Migration Path : %s", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

func getMigrationPath() (string, error) {
	cwd, _ := os.Getwd()
	migrationPath := filepath.Join(cwd, "internal", "database", "migration")

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return "", fmt.Errorf("migration directory does not exist: %s", migrationPath)
	}

	path := filepath.ToSlash(migrationPath)
	if os.PathSeparator == '\\' {
		// Windows
		return "file://" + path, nil
	}
	// Unix/Linux
	return "file://" + path, nil
}
