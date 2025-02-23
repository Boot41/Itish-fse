package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate runs the database migrations using SQL files.
func Migrate() error {
	log.Println("Starting database migration...")

	dsn := os.Getenv("DIRECT_URL")
	if dsn == "" {
		return fmt.Errorf("DIRECT_URL environment variable not set")
	}

	// Get the underlying *sql.DB from our GORM DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting underlying *sql.DB: %w", err)
	}
	// No need to close as this is managed by GORM

	// Create the postgres driver with our existing connection
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: "schema_migrations", // Optional: specify migrations table name
	})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Migration completed successfully!")
	return nil
}
