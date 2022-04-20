package models

type Report struct {
	TotalBalance  float64
	AverageDebit  float64
	AverageCredit float64
	MonthlyTxn    map[string]float64
}
