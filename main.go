package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/mariajdab/txns-email-report/models"
	"io"
	"log"
	"os"
)

var baseReport = models.BaseReport{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 {
		log.Fatalf("Wrong number of arguments: Usage %s <path-to-file>", os.Args[0])
	}

	filePath := os.Args[1]

	err := ProcessFile(filePath)
	if err != nil {
		log.Fatalln("an error happened processing the file")
	}

	fmt.Println(baseReport)
}

func ProcessFile(file string) error {

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	csvFile := csv.NewReader(f)

	fl, err := ReadFirstLine(csvFile)
	if err != nil {
		return err
	}

	if !ValidateFirstLine(fl) {
		err = errors.New("mismatch fields expected:(Id, Date, Transaction)")
		return err
	}

	var data []models.AccountTxn

	for {
		row, err := csvFile.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("an error happened reading the csv")
		}

		rowData, err := ParseData(row)
		if err != nil {
			return err
		}

		ProcessBaseReport(rowData.Transaction)

		data = append(data, models.AccountTxn{
			Id:          rowData.Id,
			Date:        rowData.Date,
			Transaction: rowData.Transaction,
		})
	}

	return nil
}
