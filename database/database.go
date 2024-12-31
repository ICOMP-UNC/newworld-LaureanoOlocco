package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

//var dbPool *pgxpool.Pool

// ConnectDB establishes a connection to the PostgreSQL database.
func ConnectDB() (*pgxpool.Pool, error) {
	// Construir la URL de la base de datos usando variables de entorno
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var pool *pgxpool.Pool
	var err error

	for retries := 5; retries > 0; retries-- {
		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return nil, fmt.Errorf("unable to parse database url: %v", err)
		}

		pool, err = pgxpool.ConnectConfig(context.Background(), config)
		if err == nil {
			return pool, nil
		}

		fmt.Printf("Failed to connect to the database: %v. Retrying in 5 seconds...\n", err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("unable to create connection pool: %v", err)
}

// InitDB initializes the database with the required tables.
func InitDB(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		jwt TEXT
	);
	CREATE TABLE IF NOT EXISTS offers (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INT NOT NULL,
		price FLOAT NOT NULL,
		category TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		items JSONB NOT NULL,
		status TEXT NOT NULL
	);
	
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	return nil
}
