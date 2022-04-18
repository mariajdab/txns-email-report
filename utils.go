package main

import (
	"encoding/csv"
	"errors"
	"github.com/mariajdab/txns-email-report/models"
	"io"
	"strconv"
	"time"
)

var fieldsCSV = []string{"Id", "Date", "Transaction"}

func ReadFirstLine(csvFile *csv.Reader) ([]string, error) {
	// get the first line
	firstLine, err := csvFile.Read()

	// error in the first line
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("an error happened file seems empty")
		} else {
			return nil, errors.New("an error happened Unprocessable Entity")
		}
	}
	return firstLine, nil
}

func ValidateFirstLine(row1 []string) bool {
	if len(fieldsCSV) != len(row1) {
		return false
	}
	for i, f := range fieldsCSV {
		if f != row1[i] {
			return false
		}
	}
	return true
}

func ParseData(row []string) (models.AccountTxn, error) {
	id, err := ParseId(row[0])
	if err != nil {
		return models.AccountTxn{}, err
	}

	date, err := ParseDate(row[1])

	if err != nil {
		return models.AccountTxn{}, err
	}

	txn, err := ParseTxn(row[2])
	if err != nil {
		return models.AccountTxn{}, err
	}

	return models.AccountTxn{
		Id:          *id,
		Date:        *date,
		Transaction: *txn,
	}, nil
}

func ParseId(rowId string) (*uint64, error) {
	id, err := strconv.ParseUint(rowId, 10, 64)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func ParseDate(rowDate string) (*time.Time, error) {
	date, err := time.Parse("1/02", rowDate)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func ParseTxn(rowTxn string) (*float64, error) {
	txn, err := strconv.ParseFloat(rowTxn, 64)
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func ProcessBaseReport(amount float64) {
	if amount < 0 {
		baseReport.AverageDebit += amount
	} else {
		baseReport.AverageCredit += amount
	}
	baseReport.TotalBalance += amount
}
