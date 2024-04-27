package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"postgres",
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Prepare(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		hashed_password VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
	`)
	return err
}
