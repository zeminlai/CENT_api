package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Add retry mechanism
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
			db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
					dbHost, dbUser, dbPass, dbName, dbPort))
			if err != nil {
					log.Printf("Failed to open database connection: %v", err)
					time.Sleep(time.Second * 2)
					continue
			}

			err = db.Ping()
			if err == nil {
					break
			}
			log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Second * 2)
	}

	if err != nil {
			log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
	}

	fmt.Println("Successfully connected to database")

	// Clean up existing tables if they exist
	_, _ = db.Exec(`DROP TABLE IF EXISTS notes, directories, users CASCADE;`)

	// Create a new migrate instance using the database driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
			log.Fatalf("Could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
			"file://cmd/storage/migrations",
			"postgres",
			driver,
	)
	if err != nil {
			log.Fatalf("Migration initialization failed: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Migration error: %v", err)
			// Try migrations again
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
					log.Fatalf("Failed to apply migrations: %v", err)
			}
	}

	fmt.Println("Database migrations applied successfully")
}

func GetDB() *sql.DB {
	return db
}
