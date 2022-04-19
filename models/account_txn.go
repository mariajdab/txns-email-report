package models

import "time"

type AccountTxn struct {
	Id          uint64
	Date        time.Time
	Transaction float64
}
