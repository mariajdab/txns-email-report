package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mariajdab/txns-email-report/database"
	"github.com/mariajdab/txns-email-report/models"
	"io"
	"log"
	"os"
)

var baseReport = models.BaseReport{}
var txnsMonth = make(map[uint64]float64)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 {
		log.Fatalf("Wrong number of arguments: Usage %s <path-to-file>", os.Args[0])
	}

	filePath := os.Args[1]

	db, err := database.OpenDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := database.CreateTable(db); err != nil {
		log.Fatal(err)
	}

	data, err := ProcessFile(filePath)
	if err != nil {
		log.Fatalln("an error happened processing the file")
	}

	database.InsertTxns(db, data)

	fmt.Println(txnsMonth)
	fmt.Println(baseReport)
}

func ProcessFile(file string) ([]models.AccountTxn, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvFile := csv.NewReader(f)

	fl, err := ReadFirstLine(csvFile)
	if err != nil {
		return nil, err
	}

	if !ValidateFirstLine(fl) {
		err = errors.New("mismatch fields expected:(Id, Date, Transaction)")
		return nil, err
	}

	var data []models.AccountTxn

	for {
		row, err := csvFile.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.New("an error happened reading the csv")
		}

		rowData, err := ParseData(row)
		if err != nil {
			return nil, err
		}

		ProcessBaseReport(rowData.Transaction)

		TnxByMonth(row[1])

		data = append(data, models.AccountTxn{
			Id:          rowData.Id,
			Date:        rowData.Date,
			Transaction: rowData.Transaction,
		})
	}
	return data, nil
}
