package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
)

var (
	//go:embed schema.sql
	schemaSQL string
)

func OpenDB() (*sql.DB, error) {
	// Connect to database
	dns := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))

	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully created connection to database")

	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return err
	}

	fmt.Println("Successfully created the table")
	return err
}
