package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"api-go/internal/models"

	// Autoloads .env file
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service represents a service that interacts with a database using GORM.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetDB returns the underlying GORM DB instance for other parts of the app to use.
	GetDB() *gorm.DB
}

// service struct now holds a *gorm.DB instance.
type service struct {
	db *gorm.DB
}

var (
	// Environment variables for database connection
	database = os.Getenv("BLUEPRINT_DB_DATABASE")
	password = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username = os.Getenv("BLUEPRINT_DB_USERNAME")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	schema   = os.Getenv("BLUEPRINT_DB_SCHEMA")

	// dbInstance holds the singleton service instance.
	dbInstance *service
)

// New initializes a new database service using GORM.
// It uses a singleton pattern to ensure only one connection pool is created.
func New() Service {
	// Reuse Connection if it already exists
	if dbInstance != nil {
		return dbInstance
	}

	// Construct the Data Source Name (DSN) for PostgreSQL.
	// Note the different format required by the GORM driver.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema)

	// Open a new GORM database connection.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Configure logger
	})
	if err != nil {
		// If connection fails, log the error and exit.
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// AutoMigrate will create or update the tables for the given models.
	// It will only add missing fields, and won't delete/change existing ones.
	log.Println("Running database migrations...")
	err = db.AutoMigrate(&models.User{}, &models.Room{}, &models.Reservation{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed.")

	// Create the singleton instance.
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Get the underlying *sql.DB instance from GORM to perform low-level checks.
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to get underlying sql.DB: %v", err)
		log.Printf("Error getting sql.DB from GORM: %v", err)
		return stats
	}

	// Ping the database to check for connectivity.
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		// Using log.Printf instead of Fatalf to avoid exiting during a health check.
		log.Printf("db down: %v", err)
		return stats
	}

	// Database is up, let's gather more statistics.
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database connection pool stats.
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a more descriptive health message.
	if dbStats.OpenConnections > 40 { // Example threshold
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider revising usage patterns."
	}

	return stats
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Printf("Disconnecting from database: %s", database)
	// Get the underlying sql.DB instance to close the connection pool.
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Printf("Error getting sql.DB from GORM for closing: %v", err)
		return err
	}
	return sqlDB.Close()
}

// GetDB provides access to the GORM DB instance for other packages to use.
func (s *service) GetDB() *gorm.DB {
	return s.db
}
