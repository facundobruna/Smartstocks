package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smartstocks/backend/internal/config"
)

type MySQLDatabase struct {
	DB *sql.DB
}

func NewMySQL(cfg *config.DatabaseConfig) (*MySQLDatabase, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Connected to MySQL database successfully")

	return &MySQLDatabase{DB: db}, nil
}

func (m *MySQLDatabase) Close() error {
	log.Println("Closing MySQL database connection...")
	return m.DB.Close()
}

func (m *MySQLDatabase) HealthCheck() error {
	return m.DB.Ping()
}

func (m *MySQLDatabase) RunMigrations(migrationSQL string) error {
	_, err := m.DB.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	log.Println("✅ Database migrations completed")
	return nil
}
