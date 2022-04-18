package models

import "time"

type AccountTxn struct {
	Id          uint64    `json:"id"`
	Date        time.Time `json:"date"`
	Transaction float64   `json:"transaction"`
}
