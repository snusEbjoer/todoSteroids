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
	CREATE TABLE IF NOT EXISTS "todo" (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		done INTEGER NOT NULL DEFAULT 0,
		description VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatad_at TIMESTAMP,
		done_at TIMESTAMP,
		user_id INTEGER NOT NULL
	);
	`)
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "history" (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		todo_id INTEGER NOT NULL,
		changet_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		done INTEGER NOT NULL DEFAULT 0,
	);
	`) // TODO add enum
	return err
}
