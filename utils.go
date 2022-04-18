package main

import (
	"github.com/mariajdab/txns-email-report/models"
	"strconv"
	"time"
)

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
