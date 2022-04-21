package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/mariajdab/txns-email-report/models"
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

//InsertTxns insert multiples values
func InsertTxns(db *sql.DB, data []models.AccountTxn) error {
	sqlStm := "INSERT INTO register_transactions(id, date_at, txn) VALUES "
	var v []interface{}

	i := 0
	for _, row := range data {
		sqlStm += fmt.Sprintf("($%d, $%d, $%d),", i+1, i+2, i+3)
		v = append(v, row.Id, row.Date, row.Transaction)
		i += 3
	}
	// remove the last ','
	sqlStm = sqlStm[0 : len(sqlStm)-1]

	//prepare the statement
	stmt, err := db.Prepare(sqlStm)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(v...)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
