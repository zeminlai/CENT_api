package storage

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
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
	// Define command line flags for database operations
	migrateUp := flag.Bool("migrate", false, "Run database migrations")
	migrateDown := flag.Bool("migrate-down", false, "Revert all migrations")
	seedDB := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

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

	// Only proceed with migrations if flags are set
	if *migrateUp || *migrateDown {
		log.Println("Starting migration process...")
		
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Printf("Driver creation error: %v", err)
			log.Fatalf("Could not create the postgres driver: %v", err)
		}
		log.Println("Database driver created successfully")

		// Check if migrations directory exists and list files
		files, err := os.ReadDir("cmd/storage/migrations")
		if err != nil {
			log.Printf("Error reading migrations directory: %v", err)
		} else {
			log.Println("Migration files found:")
			for _, file := range files {
				log.Printf("- %s", file.Name())
			}
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://cmd/storage/migrations",
			"postgres",
			driver,
		)
		if err != nil {
			log.Printf("Migration initialization error: %v", err)
			log.Fatalf("Migration initialization failed: %v", err)
		}
		log.Println("Migration instance created successfully")

		if *migrateDown {
			log.Println("Reverting all migrations...")
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Printf("Migration down error: %v", err)
				log.Fatalf("Failed to revert migrations: %v", err)
			}
			log.Println("Successfully reverted all migrations")
		}

		if *migrateUp {
			log.Println("Running migrations...")
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Printf("Migration up error: %v", err)
				log.Fatalf("Failed to apply migrations: %v", err)
			}
			log.Println("Successfully applied migrations")
		}

		// Verify tables after migration
		var tableCount int
		err = db.QueryRow(`
			SELECT COUNT(*) FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name IN ('users', 'directories', 'notes')`).Scan(&tableCount)
		if err != nil {
			log.Printf("Error checking tables: %v", err)
		} else {
			log.Printf("Found %d of 3 expected tables", tableCount)
		}
	}

	// Handle seeding if flag is set
	if *seedDB {
		log.Println("Starting database seeding...")
		
		// Read the seed file
		seedSQL, err := ioutil.ReadFile("cmd/storage/seeds/seed_data.sql")
		if err != nil {
			log.Printf("Error reading seed file: %v", err)
			return
		}

		// Execute the seed SQL
		_, err = db.Exec(string(seedSQL))
		if err != nil {
			log.Printf("Error seeding database: %v", err)
			return
		}

		log.Println("Database seeding completed successfully")

		// Verify seeded data
		var counts struct {
			users       int
			directories int
			notes       int
		}

		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&counts.users)
		err = db.QueryRow("SELECT COUNT(*) FROM directories").Scan(&counts.directories)
		err = db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&counts.notes)

		log.Printf("Seeded data counts - Users: %d, Directories: %d, Notes: %d",
			counts.users, counts.directories, counts.notes)
	}
}

func GetDB() *sql.DB {
	return db
}
