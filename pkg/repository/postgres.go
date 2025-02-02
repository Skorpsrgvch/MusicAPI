package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	SSLMode  string
}

// Подключение к PostgreSQL
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	logrus.Info("Initializing PostgreSQL connection...")

	// Формируем строку подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DbName, cfg.Password, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logrus.Errorf("Failed to open database connection: %v", err)
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		logrus.Errorf("Failed to ping database: %v", err)
		return nil, err
	}

	logrus.Info("Database connection established successfully")

	// Запуск миграций
	if err := runMigrations(db.DB); err != nil {
		logrus.Errorf("Migration error: %v", err)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

// Выполнение миграций
func runMigrations(db *sql.DB) error {
	migrationsPath := "./"

	logrus.Info("Starting database migrations...")

	goose.SetLogger(log.New(log.Writer(), "", log.LstdFlags))

	if err := goose.Up(db, migrationsPath); err != nil {
		logrus.Errorf("Failed to run migrations: %v", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logrus.Info("Database migrations applied successfully")
	return nil
}
